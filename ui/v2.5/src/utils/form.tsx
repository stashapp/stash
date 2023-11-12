import { Form, FormLabelProps } from "react-bootstrap";

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

const renderLabel = (options: {
  title: string;
  labelProps?: FormLabelProps;
}) => (
  <Form.Label column {...getLabelProps(options.labelProps)}>
    {options.title}
  </Form.Label>
);

const FormUtils = {
  renderLabel,
};

export default FormUtils;
