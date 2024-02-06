import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: [
    "../../graphql/schema/**/*.graphql",
    "graphql/client-schema.graphql",
  ],
  config: {
    // makes conflicting fields override rather than error
    onFieldTypeConflict: (_existing: unknown, other: unknown) => other,
  },
  documents: "graphql/**/*.graphql",
  generates: {
    "src/core/generated-graphql.ts": {
      plugins: [
        "time",
        "typescript",
        "typescript-operations",
        "typescript-react-apollo",
      ],
      config: {
        strictScalars: true,
        scalars: {
          Time: "string",
          Timestamp: "string",
          Map: "{ [key: string]: unknown }",
          BoolMap: "{ [key: string]: boolean }",
          PluginConfigMap: "{ [id: string]: { [key: string]: unknown } }",
          Any: "unknown",
          Int64: "number",
          Upload: "File",
          UIConfig: "src/core/config#IUIConfig",
          SavedObjectFilter: "src/models/list-filter/types#SavedObjectFilter",
          SavedUIOptions: "src/models/list-filter/types#SavedUIOptions",
        },
        withRefetchFn: true,
      },
    },
  },
};

export default config;
