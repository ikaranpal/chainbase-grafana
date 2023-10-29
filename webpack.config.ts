import { Configuration } from "webpack";
import * as path from "path";
import * as process from "process";
import CopyPlugin from "copy-webpack-plugin";

const config = async (env: { production: any }): Promise<Configuration> => {
    const PATH_SOURCE = "src";
    const PATH_DISTRIBUTION = "dist";

    const pluginJSON = require(path.resolve(process.cwd(), `${PATH_SOURCE}/plugin.json`));

    return {
        mode: env.production ? "production" : "development",
        context: path.join(process.cwd(), PATH_SOURCE),
        devtool: env.production ? "source-map" : "eval-source-map",
        entry: {
            module: "module.ts",
        },
        module: {
            rules: [
                {
                    exclude: /(node_modules)/,
                    test: /\.[tj]sx?$/,
                    use: {
                        loader: "swc-loader",
                        options: {
                            jsc: {
                                baseUrl: path.resolve(__dirname, "src"),
                                target: "es2022",
                                loose: false,
                                parser: {
                                    syntax: "typescript",
                                    tsx: true,
                                    decorators: false,
                                    dynamicImport: true,
                                },
                            },
                        },
                    },
                },
                {
                    test: /\.css$/,
                    use: ["style-loader", "css-loader"],
                },
                {
                    test: /\.s[ac]ss$/,
                    use: ["style-loader", "css-loader", "sass-loader"],
                },
                {
                    test: /\.(png|jpe?g|gif|svg)$/,
                    type: "asset/resource",
                    generator: {
                        publicPath: `public/plugins/${pluginJSON.id}/img/`,
                        outputPath: "images/",
                        filename: "[hash][ext]",
                    },
                },
            ],
        },
        externals: [
            "@grafana/data",
            "@grafana/runtime",
            "@grafana/schema",
            "@grafana/ui",
            "react",
            "react-dom",
        ],
        output: {
            clean: {
                keep: new RegExp("(.*?_(amd64|arm(64)?)(.exe)?|go_plugin_build_manifest)"),
            },
            filename: "[name].js",
            library: {
                type: "amd",
            },
            path: path.resolve(process.cwd(), PATH_DISTRIBUTION),
            publicPath: `public/plugins/${pluginJSON.id}/`,
            uniqueName: pluginJSON.id,
        },
        plugins: [
            new CopyPlugin({
                patterns: [
                    { from: "../README.md", to: ".", force: true },
                    { from: "plugin.json", to: "." },
                    { from: "../LICENSE", to: "." },
                    { from: "../CHANGELOG.md", to: ".", force: true },
                    { from: "**/*.json", to: "." },
                    { from: "**/*.svg", to: ".", noErrorOnMissing: true },
                    { from: "**/*.webp", to: ".", noErrorOnMissing: true },
                    { from: "**/*.html", to: ".", noErrorOnMissing: true },
                ],
            }),
        ],
        resolve: {
            extensions: [".js", ".jsx", ".ts", ".tsx"],
            modules: [path.resolve(process.cwd(), PATH_SOURCE), "node_modules"]
        },
    };
};

export default config;
