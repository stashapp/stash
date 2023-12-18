import React, { useContext } from "react";
import { useConfigureUI } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";
import { CheckBoxSelect } from "src/components/Shared/Select";

interface IColumn {
  label: string;
  value: string;
}

interface IProps {
  tableName: string;
  columns: IColumn[];
  defaultColumns: string[];
}

export const ColumnSelector: React.FC<IProps> = ({
  tableName,
  columns,
  defaultColumns,
}) => {
  const { configuration } = useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();
  const Toast = useToast();
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const selectedColumns = uiConfig?.tableColumns?.[tableName] ?? defaultColumns;
  const selectedOptions = columns.filter((col) =>
    selectedColumns.includes(col.value)
  );

  const handleChange = async (
    updatedColumns?: readonly IColumn[] | undefined
  ) => {
    if (!updatedColumns) {
      return;
    }

    try {
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            tableColumns: {
              ...configuration?.ui?.tableColumns,
              [tableName]: updatedColumns.map((col) => col.value),
            },
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  };

  return (
    <CheckBoxSelect
      options={columns}
      selectedOptions={selectedOptions}
      onChange={handleChange}
    />
  );
};
