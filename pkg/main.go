package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/kallydev/chainbase-grafana/pkg/plugin"
)

const PluginID = "kallydev-chainbase-datasource"

func main() {
	if err := datasource.Manage(PluginID, plugin.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())

		os.Exit(1)
	}
}
