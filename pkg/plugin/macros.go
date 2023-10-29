package plugin

import (
	"fmt"
	"regexp"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var (
	timeFilter = regexp.MustCompile(`\$__timeFilter\((.*?)\)`)
)

func ExpandMacros(query backend.DataQuery, statement string) string {
	statement = timeFilter.ReplaceAllStringFunc(statement, func(match string) string {
		column := timeFilter.FindStringSubmatch(match)[1]

		return fmt.Sprintf("%s >= '%d' AND %s <= '%d'", column, query.TimeRange.From.Unix(), column, query.TimeRange.To.Unix())
	})

	return statement
}
