import React, { ChangeEvent } from "react";
import { DataSourcePluginOptionsEditorProps } from "@grafana/data";
import { InlineField, Input } from "@grafana/ui";
import { ChainbaseDataSourceOptions } from "../types";

type Props = DataSourcePluginOptionsEditorProps<ChainbaseDataSourceOptions>;

const ConfigEditor = (props: Props) => {
    const { onOptionsChange, options } = props;
    const { jsonData } = options;

    const onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const jsonData = {
            ...options.jsonData,
            apiKey: event.target.value,
        };

        onOptionsChange({ ...options, jsonData });
    };

    return (
        <div className="gf-form-group">
            <InlineField label="API Key">
                <Input
                    onChange={onAPIKeyChange}
                    value={jsonData.apiKey || ""}
                    placeholder="Enter your API key"
                    width={48}
                />
            </InlineField>
            <InlineField label="Queries Per Second">
                <Input
                    onChange={onAPIKeyChange}
                    value={jsonData.queriesPerSecond || 1}
                    placeholder="Enter seconds"
                    width={6}
                />
            </InlineField>
        </div>
    );
};

export default ConfigEditor;
