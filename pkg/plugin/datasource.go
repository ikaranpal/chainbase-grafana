package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/kallydev/chainbase-grafana/pkg/chainbase"
	"github.com/samber/lo"
	"golang.org/x/time/rate"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var options Options
	if err := json.Unmarshal(settings.JSONData, &options); err != nil {
		return nil, fmt.Errorf("parse data source settings: %w", err)
	}

	instance := Datasource{
		limiter: rate.NewLimiter(rate.Every(time.Second), lo.Ternary(options.QueriesPerSecond == 0, 1, options.QueriesPerSecond)),
		options: options,
	}

	var err error
	if instance.client, err = chainbase.NewClient(chainbase.WithAPIKey(options.APIKey)); err != nil {
		return nil, fmt.Errorf("initlize client: %w", err)
	}

	return &instance, nil
}

type Datasource struct {
	client  *chainbase.Client
	limiter *rate.Limiter
	options Options
}

func (d *Datasource) QueryData(ctx context.Context, request *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, query := range request.Queries {
		var q Query

		if err := json.Unmarshal(query.JSON, &q); err != nil {
			response.Responses[query.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, err.Error())

			continue
		}

		if q.Statement != "" {
			response.Responses[query.RefID] = d.queryData(ctx, query)
		}
	}

	return response, nil
}

func (d *Datasource) queryData(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	q := Query{}
	if err := json.Unmarshal(query.JSON, &q); err != nil {
		return backend.ErrDataResponse(backend.StatusValidationFailed, err.Error())
	}

	var (
		dataResponse backend.DataResponse
		response     *chainbase.Response[*chainbase.DataWarehouseData]
		frame        = data.NewFrame(query.RefID)
	)

	for page := 0; page == 0 || response.Data.NextPage > 0; page++ {
		// Wait for rate limiter
		if err := d.limiter.Wait(ctx); err != nil {
			return backend.ErrDataResponse(backend.StatusTimeout, fmt.Sprintf("wait for rate limiter: %s", err))
		}

		var (
			httpResponse *http.Response
			err          error
		)

		if page == 0 {
			if response, httpResponse, err = d.client.DataWarehouse.Query(ctx, q.Statement); err != nil {
				return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("query chainbase api: %s", err))
			}
		} else {
			if response, httpResponse, err = d.client.DataWarehouse.Paginate(ctx, response.Data.TaskID, page); err != nil {
				return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("query chainbase api by task id %s: %s", response.Data.TaskID, err))
			}
		}

		if httpResponse.StatusCode != http.StatusOK {
			return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("http response has an error: %d %s", httpResponse.StatusCode, httpResponse.Status))
		}

		if response.Code != chainbase.CodeOK {
			return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("response has an error: %d %s", response.Code, response.Message))
		}

		if err := AppendRow(frame, response.Data.Meta, response.Data.Result); err != nil {
			return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("convert result for native format: %s", err))
		}
	}

	dataResponse.Frames = append(dataResponse.Frames, frame)

	return dataResponse
}

func (d *Datasource) CheckHealth(ctx context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	result := backend.CheckHealthResult{
		Status: backend.HealthStatusOk,
	}

	response, _, err := d.client.DataWarehouse.Query(ctx, "SELECT 1;")
	if err != nil {
		result.Status = backend.HealthStatusError
		result.Message = err.Error()

		return &result, nil
	}

	result.Message = response.Message

	if response.Code != chainbase.CodeOK {
		result.Status = backend.HealthStatusError

		return &result, nil
	}

	return &result, nil
}

func (d *Datasource) Dispose() {}
