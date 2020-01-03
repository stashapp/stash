import { EditableText, IOptionProps } from "@blueprintjs/core";
import { Form } from 'react-bootstrap';
import React from "react";
import { EditableTextUtils } from "./editabletext";
import { FilterMultiSelect } from "../components/select/FilterMultiSelect";
import { FilterSelect } from "../components/select/FilterSelect";
import _ from "lodash";

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
    
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          {EditableTextUtils.renderTextArea(options)}
        </td>
      </tr>
    );
  }

  public static renderInputGroup(options: {
    title: string,
    placeholder?: string,
    value: string | undefined,
    isEditing: boolean,
    onChange: ((value: string) => void),
  }) {
    let optionsCopy = _.clone(options);
    optionsCopy.placeholder = options.placeholder || options.title;
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          { !options.isEditing
              ? <h4>{optionsCopy.value}</h4>
              : <Form.Control
                  defaultValue={options.value}
                  placeholder={optionsCopy.placeholder}
                  onChange={ (event:any) => options.onChange(event.target.value) }
                />
          }
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
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          {EditableTextUtils.renderHtmlSelect(options)}
        </td>
      </tr>
    );
  }

  // TODO: isediting
  public static renderFilterSelect(options: {
    title: string,
    type: "performers" | "studios" | "tags",
    initialId: string | undefined,
    onChange: ((id: string | undefined) => void),
  }) {
    return (
      <tr>
        <td>{options.title}</td>
        <td>
          <FilterSelect
            type={options.type}
            onSelectItem={(item) => options.onChange(item ? item.id : undefined)}
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
