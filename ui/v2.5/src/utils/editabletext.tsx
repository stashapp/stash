import React from "react";
import { Form } from "react-bootstrap";
import { FilterSelect, DurationInput } from "src/components/Shared";
import { DurationUtils } from ".";

const renderTextArea = (options: {
  title: string;
  value: string | undefined;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => {
  return (
    <Form.Control
      className="text-input"
      as="textarea"
      readOnly={!options.isEditing}
      plaintext={!options.isEditing}
      onChange={(event: React.FormEvent<HTMLTextAreaElement>) =>
        options.onChange(event.currentTarget.value)
      }
      value={options.value}
    />
  );
}

const renderEditableText = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => {
  return (
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
  )
}

const renderInputGroup = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  url?: string;
  onChange: (value: string) => void;
}) => {
  if (options.url && !options.isEditing) {
    return (
      <a 
        href={options.url}
        target="_blank"
        rel="noopener noreferrer"
      >
        {options.value}
      </a>
    );
  } else {
    return (
      <Form.Control
        className="text-input"
        readOnly={!options.isEditing}
        plaintext={!options.isEditing}
        defaultValue={options.value}
        placeholder={options.placeholder ?? options.title}
        onChange={(event: React.FormEvent<HTMLInputElement>) =>
          options.onChange(event.currentTarget.value)
        }
      />
    );
  }
}

const renderDurationInput = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  url?: string;
  onChange: (value: string | undefined) => void;
}) => {
  let numericValue: number | undefined = undefined;
  if (options.value) {
    try {
      numericValue = Number.parseInt(options.value, 10);
    } catch {
      // ignore
    }
  }
  
  if (!options.isEditing) {
    let durationString = undefined;
    if (numericValue !== undefined) {
      durationString = DurationUtils.secondsToString(numericValue);
    }

    return (
      <Form.Control
        className="text-input"
        readOnly={true}
        plaintext={true}
        defaultValue={durationString}
      />
    );
  } else {
    return (
      <DurationInput
        disabled={!options.isEditing}
        numericValue={numericValue}
        onValueChange={(valueAsNumber: number) => { 
          options.onChange(valueAsNumber !== undefined ? valueAsNumber.toString() : undefined);
        }}
      />
    );
  }
}

const renderHtmlSelect = (options: {
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
  selectOptions: Array<string | number>;
}) => {
  if (!options.isEditing) {
    return (
      <Form.Control
        className="text-input"
        readOnly={true}
        plaintext={true}
        defaultValue={options.value}
      />
    );
  } else {
    return (
      <Form.Control
        as="select"
        className="input-control"
        disabled={!options.isEditing}
        plaintext={!options.isEditing}
        value={options.value?.toString()}
        onChange={(event: React.FormEvent<HTMLSelectElement>) =>
          options.onChange(event.currentTarget.value)
        }
      >
        {options.selectOptions.map((opt) => (
          <option value={opt} key={opt}>
            {opt}
          </option>
        ))}
      </Form.Control>
    );
  }
}

// TODO: isediting
const renderFilterSelect = (options: {
  type: "performers" | "studios" | "tags";
  initialId: string | undefined;
  onChange: (id: string | undefined) => void;
}) => (
  <FilterSelect
    type={options.type}
    onSelect={(items) => options.onChange(items[0]?.id)}
    initialIds={options.initialId ? [options.initialId] : []}
  />
);

// TODO: isediting
const renderMultiSelect = (options: {
  type: "performers" | "studios" | "tags";
  initialIds: string[] | undefined;
  onChange: (ids: string[]) => void;
}) => (
  <FilterSelect
    type={options.type}
    isMulti
    onSelect={(items) => options.onChange(items.map((i) => i.id))}
    initialIds={options.initialIds ?? []}
  />
);

const EditableTextUtils = {
  renderTextArea,
  renderEditableText,
  renderInputGroup,
  renderDurationInput,
  renderHtmlSelect,
  renderFilterSelect,
  renderMultiSelect,
};
export default EditableTextUtils;