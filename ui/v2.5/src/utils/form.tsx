import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { FormikValues, useFormik } from "formik";
import React, { InputHTMLAttributes, useEffect, useRef } from "react";
import {
  Button,
  Col,
  ColProps,
  Form,
  FormControlProps,
  FormLabelProps,
  Row,
} from "react-bootstrap";
import { IntlShape } from "react-intl";
import { DateInput } from "src/components/Shared/DateInput";
import { DurationInput } from "src/components/Shared/DurationInput";
import { Icon } from "src/components/Shared/Icon";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { LinkType, StashIDPill } from "src/components/Shared/StashID";
import { StringListInput } from "src/components/Shared/StringListInput";
import { URLListInput } from "src/components/Shared/URLField";
import * as GQL from "src/core/generated-graphql";

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

export function renderLabel(options: {
  title: string;
  labelProps?: FormLabelProps;
}) {
  return (
    <Form.Label column {...getLabelProps(options.labelProps)}>
      {options.title}
    </Form.Label>
  );
}

// useStopWheelScroll is a hook to provide a workaround for a bug in React/Chrome.
// If a number field is focused and the mouse pointer is over the field, then scrolling
// the mouse wheel will change the field value _and_ scroll the window.
// This hook prevents the propagation that causes the window to scroll.
export function useStopWheelScroll(ref: React.RefObject<HTMLElement>) {
  // removed the dependency array because the underlying ref value may change
  useEffect(() => {
    const { current } = ref;

    function stopWheelScroll(e: WheelEvent) {
      if (current) {
        e.stopPropagation();
      }
    }

    if (current) {
      current.addEventListener("wheel", stopWheelScroll);
    }

    return () => {
      if (current) {
        current.removeEventListener("wheel", stopWheelScroll);
      }
    };
  });
}

// NumberField is a wrapper around Form.Control that prevents wheel events from scrolling the window.
export const NumberField: React.FC<
  InputHTMLAttributes<HTMLInputElement> & FormControlProps
> = (props) => {
  const inputRef = useRef<HTMLInputElement>(null);

  useStopWheelScroll(inputRef);

  return <Form.Control {...props} type="number" ref={inputRef} />;
};

type Formik<V extends FormikValues> = ReturnType<typeof useFormik<V>>;

interface IProps {
  labelProps?: FormLabelProps;
  fieldProps?: ColProps;
}

export function formikUtils<V extends FormikValues>(
  intl: IntlShape,
  formik: Formik<V>,
  {
    labelProps = {
      column: true,
      sm: 3,
      xl: 2,
    },
    fieldProps = {
      sm: 9,
      xl: 7,
    },
  }: IProps = {}
) {
  type Field = keyof V & string;
  type ErrorMessage = string | undefined;

  function renderFormControl(field: Field, type: string, placeholder: string) {
    const formikProps = formik.getFieldProps({ name: field, type: type });
    const error = formik.errors[field] as ErrorMessage;

    let { value } = formikProps;
    if (value === null) {
      value = "";
    }

    let control: React.ReactNode;
    if (type === "checkbox") {
      control = (
        <Form.Check
          placeholder={placeholder}
          {...formikProps}
          value={value}
          isInvalid={!!error}
        />
      );
    } else if (type === "textarea") {
      control = (
        <Form.Control
          as="textarea"
          className="text-input"
          placeholder={placeholder}
          {...formikProps}
          value={value}
          isInvalid={!!error}
        />
      );
    } else if (type === "number") {
      control = (
        <NumberField
          type={type}
          className="text-input"
          placeholder={placeholder}
          {...formikProps}
          value={value}
          isInvalid={!!error}
        />
      );
    } else {
      control = (
        <Form.Control
          type={type}
          className="text-input"
          placeholder={placeholder}
          {...formikProps}
          value={value}
          isInvalid={!!error}
        />
      );
    }

    return (
      <>
        {control}
        <Form.Control.Feedback type="invalid">{error}</Form.Control.Feedback>
      </>
    );
  }

  function renderField(
    field: Field,
    title: string,
    control: React.ReactNode,
    props?: IProps
  ) {
    return (
      <Form.Group controlId={field} as={Row} data-field={field}>
        <Form.Label {...(props?.labelProps ?? labelProps)}>{title}</Form.Label>
        <Col {...(props?.fieldProps ?? fieldProps)}>{control}</Col>
      </Form.Group>
    );
  }

  function renderInputField(
    field: Field,
    type: string = "text",
    messageID: string = field,
    props?: IProps
  ) {
    const title = intl.formatMessage({ id: messageID });
    const control = renderFormControl(field, type, title);

    return renderField(field, title, control, props);
  }

  function renderSelectField(
    field: Field,
    entries: Map<string, string>,
    messageID: string = field,
    props?: IProps
  ) {
    const formikProps = formik.getFieldProps(field);

    let { value } = formikProps;
    if (value === null) {
      value = "";
    }

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <Form.Control
        as="select"
        className="input-control"
        {...formikProps}
        value={value}
      >
        <option value="" key=""></option>
        {Array.from(entries).map(([k, v]) => (
          <option value={v} key={v}>
            {k}
          </option>
        ))}
      </Form.Control>
    );

    return renderField(field, title, control, props);
  }

  function renderDateField(
    field: Field,
    messageID: string = field,
    props?: IProps
  ) {
    const value = formik.values[field] as string;
    const error = formik.errors[field] as ErrorMessage;

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <DateInput
        value={value}
        onValueChange={(v) => formik.setFieldValue(field, v)}
        error={error}
      />
    );

    return renderField(field, title, control, props);
  }

  function renderDurationField(
    field: Field,
    messageID: string = field,
    props?: IProps
  ) {
    const value = formik.values[field] as number | null;
    const error = formik.errors[field] as ErrorMessage;

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <DurationInput
        value={value}
        setValue={(v) => formik.setFieldValue(field, v)}
        error={error}
      />
    );

    return renderField(field, title, control, props);
  }

  function renderRatingField(
    field: Field,
    messageID: string = field,
    props?: IProps
  ) {
    const value = formik.values[field] as number | null;

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <RatingSystem
        value={value}
        onSetRating={(v) => formik.setFieldValue(field, v)}
      />
    );

    return renderField(field, title, control, props);
  }

  // flattens a potential list of errors into a [errorMsg, errorIdx] tuple
  // error messages are joined with newlines, and duplicate messages are skipped
  function flattenError(
    error: ErrorMessage[] | ErrorMessage
  ): [string | undefined, number[] | undefined] {
    if (Array.isArray(error)) {
      let errors: string[] = [];
      const errorIdx = [];
      for (let i = 0; i < error.length; i++) {
        const err = error[i];
        if (err) {
          if (!errors.includes(err)) {
            errors.push(err);
          }
          errorIdx.push(i);
        }
      }
      return [errors.join("\n"), errorIdx];
    } else {
      return [error, undefined];
    }
  }

  interface IStringListProps extends IProps {
    // defaults to true if not provided
    orderable?: boolean;
  }

  function renderStringListField(
    field: Field,
    messageID: string = field,
    props?: IStringListProps
  ) {
    const value = formik.values[field] as string[];
    const error = formik.errors[field] as ErrorMessage[] | ErrorMessage;

    const [errorMsg, errorIdx] = flattenError(error);

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <StringListInput
        value={value}
        setValue={(v) => formik.setFieldValue(field, v)}
        errors={errorMsg}
        errorIdx={errorIdx}
        orderable={props?.orderable}
      />
    );

    return renderField(field, title, control, props);
  }

  function renderURLListField(
    field: Field,
    onScrapeClick?: (url: string) => void,
    urlScrapable?: (url: string) => boolean,
    messageID: string = field,
    props?: IProps
  ) {
    const value = formik.values[field] as string[];
    const error = formik.errors[field] as ErrorMessage[] | ErrorMessage;

    const [errorMsg, errorIdx] = flattenError(error);

    const title = intl.formatMessage({ id: messageID });
    const control = (
      <URLListInput
        value={value}
        setValue={(v) => formik.setFieldValue(field, v)}
        errors={errorMsg}
        errorIdx={errorIdx}
        onScrapeClick={onScrapeClick}
        urlScrapable={urlScrapable}
      />
    );

    return renderField(field, title, control, props);
  }

  function renderStashIDsField(
    field: Field,
    linkType: LinkType,
    messageID: string = field,
    props?: IProps,
    addButton?: React.ReactNode
  ) {
    const values = formik.values[field] as GQL.StashIdInput[];

    const title = intl.formatMessage({ id: messageID });

    const removeStashID = (stashID: GQL.StashIdInput) => {
      const v = values.filter((s) => s !== stashID);
      formik.setFieldValue(field, v);
    };

    const control = (
      <>
        {values.length > 0 && (
          <ul className="pl-0 mb-2">
            {values.map((stashID) => {
              return (
                <Row as="li" key={stashID.stash_id} noGutters>
                  <Button
                    variant="danger"
                    className="mr-2 py-0"
                    title={intl.formatMessage(
                      { id: "actions.delete_entity" },
                      { entityType: intl.formatMessage({ id: "stash_id" }) }
                    )}
                    onClick={() => removeStashID(stashID)}
                  >
                    <Icon icon={faTrashAlt} />
                  </Button>
                  <StashIDPill stashID={stashID} linkType={linkType} />
                </Row>
              );
            })}
          </ul>
        )}
        {addButton}
      </>
    );

    return renderField(field, title, control, props);
  }

  return {
    renderFormControl,
    renderField,
    renderInputField,
    renderSelectField,
    renderDateField,
    renderDurationField,
    renderRatingField,
    renderStringListField,
    renderURLListField,
    renderStashIDsField,
  };
}
