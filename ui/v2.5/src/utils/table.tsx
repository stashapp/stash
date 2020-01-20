import React from "react";
import { Form } from "react-bootstrap";
import { FilterSelect } from "src/components/Shared";

const renderEditableTextTableRow = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <Form.Control
        readOnly={!options.isEditing}
        plaintext={!options.isEditing}
        onChange={(event: React.FormEvent<HTMLInputElement>) =>
          options.onChange(event.currentTarget.value)
        }
        value={
          typeof options.value === "number"
            ? options.value.toString()
            : options.value
        }
        placeholder={options.title}
      />
    </td>
  </tr>
);

const renderTextArea = (options: {
  title: string;
  value: string | undefined;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <Form.Control
        as="textarea"
        readOnly={!options.isEditing}
        plaintext={!options.isEditing}
        onChange={(event: React.FormEvent<HTMLTextAreaElement>) =>
          options.onChange(event.currentTarget.value)
        }
        value={options.value}
      />
    </td>
  </tr>
);

const renderInputGroup = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <Form.Control
        readOnly={!options.isEditing}
        plaintext={!options.isEditing}
        defaultValue={options.value}
        placeholder={options.placeholder ?? options.title}
        onChange={(event: React.FormEvent<HTMLInputElement>) =>
          options.onChange(event.currentTarget.value)
        }
      />
    </td>
  </tr>
);

const renderHtmlSelect = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
  selectOptions: Array<string | number>;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <Form.Control
        as="select"
        readOnly={!options.isEditing}
        plaintext={!options.isEditing}
        onChange={(event: React.FormEvent<HTMLSelectElement>) =>
          options.onChange(event.currentTarget.value)
        }
      />
    </td>
  </tr>
);

// TODO: isediting
const renderFilterSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialId: string | undefined;
  onChange: (id: string | undefined) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <FilterSelect
        type={options.type}
        onSelect={items => options.onChange(items[0]?.id)}
        initialIds={options.initialId ? [options.initialId] : []}
      />
    </td>
  </tr>
);

// TODO: isediting
const renderMultiSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialIds: string[] | undefined;
  onChange: (ids: string[]) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>
      <FilterSelect
        type={options.type}
        isMulti
        onSelect={items => options.onChange(items.map(i => i.id))}
        initialIds={options.initialIds ?? []}
      />
    </td>
  </tr>
);

const Table = {
  renderEditableTextTableRow,
  renderTextArea,
  renderInputGroup,
  renderHtmlSelect,
  renderFilterSelect,
  renderMultiSelect
};
export default Table;
