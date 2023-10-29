package plugin_test

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/kallydev/chainbase-grafana/pkg/plugin"
	"github.com/stretchr/testify/require"
)

func TestExpandMacros(t *testing.T) {
	type arguments struct {
		query     backend.DataQuery
		statement string
	}

	testcases := []struct {
		name      string
		arguments arguments
		want      require.ValueAssertionFunc
		wantError require.ErrorAssertionFunc
	}{
		{
			name: "Time filter",
			arguments: arguments{
				query: backend.DataQuery{
					TimeRange: backend.TimeRange{
						From: time.Date(2023, 10, 29, 13, 00, 00, 00, time.UTC),
						To:   time.Date(2023, 10, 29, 14, 00, 00, 00, time.UTC),
					},
				},
				statement: `SELECT reinterpretAsUInt256(reverse(unhex(topic1))) / 1e8 AS price
FROM ethereum.transaction_logs
WHERE address = '0xE62B71cf983019BFf55bC83B48601ce8419650CC'
  AND $__timeFilter(block_timestamp)
  AND topic0 = '0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f'
ORDER BY block_timestamp DESC
LIMIT 1;
`,
			},
			want: func(t require.TestingT, value interface{}, msgAndArgs ...interface{}) {
				require.Equal(t, value, `SELECT reinterpretAsUInt256(reverse(unhex(topic1))) / 1e8 AS price
FROM ethereum.transaction_logs
WHERE address = '0xE62B71cf983019BFf55bC83B48601ce8419650CC'
  AND block_timestamp >= '1698584400' AND block_timestamp <= '1698588000'
  AND topic0 = '0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f'
ORDER BY block_timestamp DESC
LIMIT 1;
`)
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			statement := plugin.ExpandMacros(testcase.arguments.query, testcase.arguments.statement)
			if testcase.want != nil {
				testcase.want(t, statement)
			}
		})
	}
}
