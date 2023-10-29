import { DataQuery, DataSourceJsonData } from "@grafana/schema";

export interface ChainbaseQuery extends DataQuery {
    statement: string;
}

export interface ChainbaseDataSourceOptions extends DataSourceJsonData {
    apiKey: string;
    queriesPerSecond: number;
}
