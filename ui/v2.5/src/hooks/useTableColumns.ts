import { useContext } from "react";
import { ConfigurationContext } from "src/hooks/Config";

export const useTableColumns = (
  tableName: string,
  defaultColumns: string[]
) => {
  const { configuration } = useContext(ConfigurationContext);
  const selectedColumns: string[] =
    configuration?.ui?.tableColumns?.[tableName] ?? defaultColumns;
  return Object.fromEntries(selectedColumns.map((col) => [col, true]));
};
