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
        scalars: {
          Time: "string",
          Timestamp: "string",
          Map: "{ [key: string]: any }",
          Any: "any",
          Int64: "number",
          Upload: "File",
          UIConfig: "src/core/config#IUIConfig",
        },
        withRefetchFn: true,
      },
    },
  },
};

export default config;
