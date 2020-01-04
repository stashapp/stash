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
    asURL?: boolean,
    urlPrefix?: string,
    onChange: ((value: string) => void),
  }) {
    function maybeRenderURL() {
      if (options.asURL) {
        let url = options.value;
        if (options.urlPrefix) {
          url = options.urlPrefix + url;
        }
        return <a href={url}>{options.value}</a>
      } else {
        return options.value;
      }
    }

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
        element = <Label>{maybeRenderURL()}</Label>;
      } else {
        element = <span>{maybeRenderURL()}</span>;
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