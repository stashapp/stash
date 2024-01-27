import { useContext } from "react";
import { useConfigureUI } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { useToast } from "./Toast";

export const useTableColumns = (
  tableName: string,
  defaultColumns: string[]
) => {
  const Toast = useToast();

  const { configuration } = useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const selectedColumns: string[] =
    configuration?.ui?.tableColumns?.[tableName] ?? defaultColumns;

  async function saveColumns(updatedColumns: readonly string[]) {
    try {
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            tableColumns: {
              ...configuration?.ui?.tableColumns,
              [tableName]: updatedColumns,
            },
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  return { selectedColumns, saveColumns };
};
