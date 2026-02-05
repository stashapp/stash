import { useConfigureUI } from "src/core/StashService";
import { useConfigurationContext } from "src/hooks/Config";
import { useToast } from "./Toast";

export const useTableColumns = (
  tableName: string,
  defaultColumns: string[]
) => {
  const Toast = useToast();

  const { configuration } = useConfigurationContext();
  const [saveUI] = useConfigureUI();

  const ui = configuration?.ui;

  const selectedColumns = ui?.tableColumns?.[tableName] ?? defaultColumns;

  async function saveColumns(updatedColumns: string[]) {
    try {
      await saveUI({
        variables: {
          input: {
            ...ui,
            tableColumns: {
              ...ui?.tableColumns,
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
