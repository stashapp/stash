import { EditableText, HTMLSelect, InputGroup, IOptionProps, TextArea } from "@blueprintjs/core";
import React from "react";
import { FilterMultiSelect } from "../components/select/FilterMultiSelect";
import { FilterSelect } from "../components/select/FilterSelect";

export class TableUtils {
  public static renderEditableTextTableRow(options: {
    title: string;
    value: string | number | undefined;
    isEditing: boolean;
    onChange: ((value: string) => void);
  }) {
    let stringValue = options.value;
    if (typeof stringValue === "number") {
      stringValue = stringValue.toString();
    }
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          <EditableText
            disabled={!options.isEditing}
            value={stringValue}
            placeholder={options.title}
            multiline={true}
            onChange={(newValue) => options.onChange(newValue)}
          />
        </td>
      </tr>
    );
  }

  public static renderTextArea(options: {
    title: string,
    value: string | undefined,
    isEditing: boolean,
    onChange: ((value: string) => void),
  }) {
    let element: JSX.Element;
    if (options.isEditing) {
      element = (
        <TextArea
          fill={true}
          onChange={(newValue) => options.onChange(newValue.target.value)}
          value={options.value}
        />
      );
    } else {
      element = <p className="pre">{options.value}</p>;
    }
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          {element}
        </td>
      </tr>
    );
  }

  public static renderInputGroup(options: {
    title: string,
    value: string | undefined,
    isEditing: boolean,
    onChange: ((value: string) => void),
  }) {
    let element: JSX.Element;
    if (options.isEditing) {
      element = (
        <InputGroup
          onChange={(newValue: any) => options.onChange(newValue.target.value)}
          value={options.value}
        />
      );
    } else {
      element = <span>{options.value}</span>;
    }
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          {element}
        </td>
      </tr>
    );
  }

  public static renderHtmlSelect(options: {
    title: string,
    value: string | number | undefined,
    isEditing: boolean,
    onChange: ((value: string) => void),
    selectOptions: Array<string | number | IOptionProps>,
  }) {
    let stringValue = options.value;
    if (typeof stringValue === "number") {
      stringValue = stringValue.toString();
    }

    let element: JSX.Element;
    if (options.isEditing) {
      element = (
        <HTMLSelect
          options={options.selectOptions}
          onChange={(event) => options.onChange(event.target.value)}
          value={stringValue}
        />
      );
    } else {
      element = <span>{options.value}</span>;
    }
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          {element}
        </td>
      </tr>
    );
  }

  // TODO: isediting
  public static renderFilterSelect(options: {
    title: string,
    type: "performers" | "studios" | "tags",
    initialId: string | undefined,
    onChange: ((id: string) => void),
  }) {
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          <FilterSelect
            type={options.type}
            onSelectItem={(item) => options.onChange(item.id)}
            initialId={options.initialId}
          />
        </td>
      </tr>
    );
  }

  // TODO: isediting
  public static renderMultiSelect(options: {
    title: string,
    type: "performers" | "studios" | "tags",
    initialIds: string[] | undefined,
    onChange: ((ids: string[]) => void),
  }) {
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          <FilterMultiSelect
            type={options.type}
            onUpdate={(items) => options.onChange(items.map((i) => i.id))}
            openOnKeyDown={true}
            initialIds={options.initialIds}
          />
        </td>
      </tr>
    );
  }
}
