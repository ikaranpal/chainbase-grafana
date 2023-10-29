package plugin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/kallydev/chainbase-grafana/pkg/chainbase"
)

var _ error = (*UnsupportedTypeError)(nil)

type UnsupportedTypeError struct {
	Field string
	Type  string
}

func NewUnsupportedTypeError(columnName, columnType string) UnsupportedTypeError {
	return UnsupportedTypeError{
		Field: columnName,
		Type:  columnType,
	}
}

func (u UnsupportedTypeError) Error() string {
	return fmt.Sprintf("type %s of field %s is not supported", u.Type, u.Field)
}

var _ error = (*InvalidValueError)(nil)

type InvalidValueError struct {
	Field string
	Value any
}

func NewInvalidValueError(columnName string, columnValue any) InvalidValueError {
	return InvalidValueError{
		Field: columnName,
		Value: columnValue,
	}
}

func (i InvalidValueError) Error() string {
	return fmt.Sprintf("type %T of value %s is invalid", i.Value, i.Field)
}

func AppendRow(frame *data.Frame, columns []chainbase.DataWarehouseDataMeta, rows []map[string]any) error {
	// Initialize fields if the frame is not initialized
	if len(frame.Fields) == 0 {
		for _, column := range columns {
			switch columnType := proto.ColumnType(column.Type); columnType {
			case proto.ColumnTypeInt8:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []int8{}))
			case proto.ColumnTypeInt16:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []int16{}))
			case proto.ColumnTypeInt32:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []int32{}))
			case proto.ColumnTypeInt128, proto.ColumnTypeInt256:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []string{}))
			case proto.ColumnTypeUInt8:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []uint8{}))
			case proto.ColumnTypeUInt16:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []uint16{}))
			case proto.ColumnTypeUInt32:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []uint32{}))
			case proto.ColumnTypeUInt64:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []uint64{}))
			case proto.ColumnTypeUInt128, proto.ColumnTypeUInt256:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []string{}))
			case proto.ColumnTypeFloat32:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []float32{}))
			case proto.ColumnTypeFloat64:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []float64{}))
			case proto.ColumnTypeString:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []string{}))
			case proto.ColumnTypeDateTime, proto.ColumnTypeDate:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []time.Time{}))
			case proto.ColumnTypeBool:
				frame.Fields = append(frame.Fields, data.NewField(column.Name, nil, []bool{}))
			default:
				return NewUnsupportedTypeError(column.Name, column.Type)
			}
		}
	}

	// Convert and append rows
	for index, column := range columns {
		var (
			field      = frame.Fields[index]
			columnType = proto.ColumnType(column.Type)
		)

		for _, row := range rows {
			switch columnType {
			case
				proto.ColumnTypeInt8, proto.ColumnTypeInt16, proto.ColumnTypeInt32,
				proto.ColumnTypeUInt8, proto.ColumnTypeUInt16, proto.ColumnTypeUInt32,
				proto.ColumnTypeFloat32, proto.ColumnTypeFloat64:
				value, ok := row[column.Name].(float64)
				if !ok {
					return NewInvalidValueError(column.Name, row[column.Name])
				}

				switch columnType {
				case proto.ColumnTypeInt8:
					field.Append(int8(value))
				case proto.ColumnTypeInt16:
					field.Append(int16(value))
				case proto.ColumnTypeInt32:
					field.Append(int32(value))
				case proto.ColumnTypeUInt8:
					field.Append(uint8(value))
				case proto.ColumnTypeUInt16:
					field.Append(uint16(value))
				case proto.ColumnTypeUInt32:
					field.Append(uint32(value))
				case proto.ColumnTypeFloat32:
					field.Append(float32(value))
				case proto.ColumnTypeFloat64:
					field.Append(value)
				}
			case
				proto.ColumnTypeInt64, proto.ColumnTypeInt128, proto.ColumnTypeInt256,
				proto.ColumnTypeUInt64, proto.ColumnTypeUInt128, proto.ColumnTypeUInt256,
				proto.ColumnTypeString,
				proto.ColumnTypeDateTime, proto.ColumnTypeDate:
				value, ok := row[column.Name].(string)
				if !ok {
					return NewInvalidValueError(column.Name, row[column.Name])
				}

				switch columnType {
				case proto.ColumnTypeInt64:
					value, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return NewInvalidValueError(column.Name, row[column.Name])
					}

					field.Append(value)
				case proto.ColumnTypeUInt64:
					value, err := strconv.ParseUint(value, 10, 64)
					if err != nil {
						return NewInvalidValueError(column.Name, row[column.Name])
					}

					field.Append(value)
				case proto.ColumnTypeString:
					field.Append(value)
				case proto.ColumnTypeDateTime:
					dateTime, err := time.Parse(time.DateTime, value)
					if err != nil {
						return NewInvalidValueError(column.Name, row[column.Name])
					}

					field.Append(dateTime)
				case proto.ColumnTypeDate:
					value, err := time.Parse(proto.DateLayout, value)
					if err != nil {
						return NewInvalidValueError(column.Name, row[column.Name])
					}

					field.Append(value)
				default:
					return NewUnsupportedTypeError(column.Name, column.Type)
				}
			case proto.ColumnTypeBool:
				value, ok := row[column.Name].(bool)
				if !ok {
					return NewInvalidValueError(column.Name, row[column.Name])
				}

				field.Append(value)
			default:
				return NewUnsupportedTypeError(field.Name, field.Type().String())
			}
		}
	}

	return nil
}
