import React, { useEffect, useState, useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Dropdown, DropdownButton, Form, Col, Row, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import FormUtils from "src/utils/form";
import { DateInput } from "src/components/Shared/DateInput";
import { useToast } from "src/hooks/Toast";

import { useFormik } from "formik";

interface ISceneHistoryProps {
  scene: GQL.SceneDataFragment; //Partial<GQL.SceneDataFragment>;
  onSubmit: (input: GQL.SceneCreateInput) => Promise<void>;
}

const SceneHistoryPanel: React.FC<ISceneHistoryProps> = ({
  scene,
  onSubmit,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const schema = yup.object({
    date: yup
      .string()
      .ensure()
      .test({
        name: "date",
        test: (value) => {
          if (!value) return true;
          if (!value.match(/^\d{4}-\d{2}-\d{2}$/)) return false;
          if (Number.isNaN(Date.parse(value))) return false;
          return true;
        },
        message: intl.formatMessage({ id: "validation.date_invalid_form" }),
      }),
    cover_image: yup.string().nullable().optional(),
  });

  const initialValues = useMemo(
    () => ({
      date: scene.date ?? "",
    }),
    [scene]
  );

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
  });

  async function onSave(input: InputValues) {
    setIsLoading(true);
    try {
      await onSubmit(input);
      formik.resetForm();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  type InputValues = yup.InferType<typeof schema>;

  return (
    <Form.Group controlId="date" as={Row}>
      {FormUtils.renderLabel({
        title: intl.formatMessage({ id: "date" }),
      })}

      <Col xs={9}>
        <DateInput
          value={formik.values.date}
          onValueChange={(value) => formik.setFieldValue("date", value)}
          error={formik.errors.date}
        />
      </Col>
    </Form.Group>
  );
};

export default SceneHistoryPanel;
