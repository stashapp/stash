import React from "react";
import { Form, Col, Row, ColProps, FormLabelProps } from "react-bootstrap";
import EditableTextUtils from "./editabletext";

function getLabelProps(labelProps?: FormLabelProps) {
  let ret = labelProps;
  if (!ret) {
    ret = {
      column: true,
      xs: 3,
    };
  }

  return ret;
}

function getInputProps(inputProps?: ColProps) {
  let ret = inputProps;
  if (!ret) {
    ret = {
      xs: 9,
    };
  }

  return ret;
}

const renderLabel = (options: {
  title: string;
  labelProps?: FormLabelProps;
}) => (
  <Form.Label column {...getLabelProps(options.labelProps)}>
    {options.title}
  </Form.Label>
);

const renderEditableText = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderEditableText(options)}
    </Col>
  </Form.Group>
);

const renderTextArea = (options: {
  title: string;
  value: string | undefined;
  isEditing: boolean;
  onChange: (value: string) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderTextArea(options)}
    </Col>
  </Form.Group>
);

const renderInputGroup = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  url?: string;
  onChange: (value: string) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderInputGroup(options)}
    </Col>
  </Form.Group>
);

const renderDurationInput = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  asString?: boolean;
  onChange: (value: string | undefined) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => {
  return (
    <Form.Group controlId={options.title} as={Row}>
      {renderLabel(options)}
      <Col {...getInputProps(options.inputProps)}>
        {EditableTextUtils.renderDurationInput(options)}
      </Col>
    </Form.Group>
  );
};

const renderHtmlSelect = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
  selectOptions: Array<string | number>;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderHtmlSelect(options)}
    </Col>
  </Form.Group>
);

// TODO: isediting
const renderFilterSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialId: string | undefined;
  onChange: (id: string | undefined) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderFilterSelect(options)}
    </Col>
  </Form.Group>
);

// TODO: isediting
const renderMultiSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialIds: string[] | undefined;
  onChange: (ids: string[]) => void;
  labelProps?: FormLabelProps;
  inputProps?: ColProps;
}) => (
  <Form.Group controlId={options.title} as={Row}>
    {renderLabel(options)}
    <Col {...getInputProps(options.inputProps)}>
      {EditableTextUtils.renderMultiSelect(options)}
    </Col>
  </Form.Group>
);

const FormUtils = {
  renderLabel,
  renderEditableText,
  renderTextArea,
  renderInputGroup,
  renderDurationInput,
  renderHtmlSelect,
  renderFilterSelect,
  renderMultiSelect,
};
export default FormUtils;
