import { HTMLSelect, InputGroup, IOptionProps, TextArea, Label } from "@blueprintjs/core";
import React from "react";

export class EditableTextUtils {
  public static renderTextArea(options: {
    value: string | undefined,
    isEditing: boolean,
    onChange: ((value: string) => void)
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
    return element;
  }

  public static renderInputGroup(options: {
    value: string | undefined,
    isEditing: boolean,
    placeholder?: string,
    asLabel?: boolean,
    onChange: ((value: string) => void),
  }) {
    let element: JSX.Element;
    if (options.isEditing) {
      element = (
        <InputGroup
          onChange={(newValue: any) => options.onChange(newValue.target.value)}
          value={options.value}
          placeholder={options.placeholder}
        />
      );
    } else {
      if (options.asLabel) {
        element = <Label>{options.value}</Label>;
      } else {
        element = <span>{options.value}</span>;
      }
    }
    return element;
  }

  public static renderHtmlSelect(options: {
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
    return element;
  }
}