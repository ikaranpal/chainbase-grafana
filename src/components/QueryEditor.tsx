import { CodeEditor } from "@grafana/ui";
import React from "react";
import { QueryEditorProps } from "@grafana/data";
import { DataSource } from "../datasource";
import { ChainbaseDataSourceOptions, ChainbaseQuery } from "../types";
import "./QueryEditor.css";

type Props = QueryEditorProps<DataSource, ChainbaseQuery, ChainbaseDataSourceOptions>;

const QueryEditor = (props: Props) => {
    const { query, onChange } = props;
    const { statement } = query;

    const defaultHeight = "192px";

    const onQueryTextChange = ((value: string) => {
        onChange({ ...query, statement: value });
    });

    return (
        <div className="wrapper">
            <CodeEditor
                aria-label="SQL"
                language="sql"
                height={defaultHeight}
                showMiniMap={false}
                showLineNumbers={true}
                onBlur={onQueryTextChange}
                value={statement || ""}/>
        </div>
    );
};

export default QueryEditor;
