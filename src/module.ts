import { DataSourcePlugin } from "@grafana/data";
import { DataSource } from "./datasource";
import { ChainbaseDataSourceOptions, ChainbaseQuery } from "./types";
import ConfigEditor from "./components/ConfigEditor";
import QueryEditor from "./components/QueryEditor";

export const plugin = new DataSourcePlugin<DataSource, ChainbaseQuery, ChainbaseDataSourceOptions>(DataSource)
    .setConfigEditor(ConfigEditor)
    .setQueryEditor(QueryEditor);
