import React, { useMemo } from "react";
import { Table, Form } from "react-bootstrap";
import { CheckBoxSelect } from "../Shared/Select";
import cx from "classnames";
import { SortDirectionEnum } from "src/core/generated-graphql";
import { SortByIcon } from "./ListFilter";

export interface IColumn {
  label: string;
  value: string;
  mandatory?: boolean;
  sortable?: boolean;
}

export const ColumnSelector: React.FC<{
  selected: string[];
  allColumns: IColumn[];
  setSelected: (selected: string[]) => void;
}> = ({ selected, allColumns, setSelected }) => {
  const disableOptions = useMemo(() => {
    return allColumns.map((col) => {
      return {
        ...col,
        isDisabled: col.mandatory,
      };
    });
  }, [allColumns]);

  const selectedColumns = useMemo(() => {
    return disableOptions.filter((col) => selected.includes(col.value));
  }, [selected, disableOptions]);

  return (
    <CheckBoxSelect
      options={disableOptions}
      selectedOptions={selectedColumns}
      onChange={(v) => {
        setSelected(v.map((col) => col.value));
      }}
    />
  );
};

interface IListTableProps<T> {
  className?: string;
  items: T[];
  columns: string[];
  setColumns: (columns: string[]) => void;
  allColumns: IColumn[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  renderCell: (column: IColumn, item: T, index: number) => React.ReactNode;
  onSort?: (value: string) => void;
  sortBy?: string;
  sortDirection?: SortDirectionEnum;
}

export const ListTable = <T extends { id: string }>(
  props: IListTableProps<T>
) => {
  const {
    className,
    items,
    columns,
    setColumns,
    allColumns,
    selectedIds,
    onSelectChange,
    renderCell,
    onSort,
    sortBy,
    sortDirection,
  } = props;

  const visibleColumns = useMemo(() => {
    return allColumns.filter(
      (col) => col.mandatory || columns.includes(col.value)
    );
  }, [columns, allColumns]);

  const renderObjectRow = (item: T, index: number) => {
    let shiftKey = false;

    return (
      <tr key={item.id}>
        <td className="select-col">
          <label>
            <Form.Control
              type="checkbox"
              checked={selectedIds.has(item.id)}
              onChange={() =>
                onSelectChange(item.id, !selectedIds.has(item.id), shiftKey)
              }
              onClick={(
                event: React.MouseEvent<HTMLInputElement, MouseEvent>
              ) => {
                shiftKey = event.shiftKey;
                event.stopPropagation();
              }}
            />
          </label>
        </td>

        {visibleColumns.map((column) => (
          <td key={column.value} className={`${column.value}-data`}>
            {renderCell(column, item, index)}
          </td>
        ))}
      </tr>
    );
  };

  const columnHeaders = useMemo(() => {
    return visibleColumns.map((column) => {
      const isSortable = column.sortable !== false && onSort;
      const isActive = sortBy === column.value;

      return (
        <th
          key={column.value}
          className={`${column.value}-head${isSortable ? " sortable" : ""}`}
          onClick={isSortable ? () => onSort(column.value) : undefined}
          style={isSortable ? { cursor: "pointer" } : undefined}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            {column.label}
            {isActive && sortDirection && (
              <SortByIcon sortDirection={sortDirection} />
            )}
          </div>
        </th>
      );
    });
  }, [visibleColumns, onSort, sortBy, sortDirection]);

  return (
    <div className={cx("table-list", className)}>
      <Table striped bordered>
        <thead>
          <tr>
            <th className="select-col">
              <div
                className="d-inline-block"
                data-toggle="popover"
                data-trigger="focus"
              >
                <ColumnSelector
                  allColumns={allColumns}
                  selected={columns}
                  setSelected={setColumns}
                />
              </div>
            </th>

            {columnHeaders}
          </tr>
          <tr>
            <th className="border-row" colSpan={100}></th>
          </tr>
        </thead>
        <tbody>{items.map(renderObjectRow)}</tbody>
      </Table>
    </div>
  );
};
