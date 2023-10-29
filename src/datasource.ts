import { DataSourceInstanceSettings } from "@grafana/data";
import { ChainbaseDataSourceOptions, ChainbaseQuery } from "./types";
import { DataSourceWithBackend } from "@grafana/runtime";

export class DataSource extends DataSourceWithBackend<ChainbaseQuery, ChainbaseDataSourceOptions> {
    constructor(instanceSettings: DataSourceInstanceSettings<ChainbaseDataSourceOptions>) {
        super(instanceSettings);
    }
}
