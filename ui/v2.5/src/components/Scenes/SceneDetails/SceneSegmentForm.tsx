import React, { useMemo } from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useFormik } from "formik";
import * as yup from "yup";
import * as GQL from "src/core/generated-graphql";
import {
  useSceneSegmentCreate,
  useSceneSegmentUpdate,
  useSceneSegmentDestroy,
} from "src/core/StashService";
import { DurationInput } from "src/components/Shared/DurationInput";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import { useToast } from "src/hooks/Toast";
import isEqual from "lodash-es/isEqual";
import { formikUtils } from "src/utils/form";
import { yupFormikValidate } from "src/utils/yup";

interface ISceneSegmentForm {
  sceneID: string;
  segment?: GQL.SceneSegmentDataFragment;
  onClose: () => void;
}

export const SceneSegmentForm: React.FC<ISceneSegmentForm> = ({
  sceneID,
  segment,
  onClose,
}) => {
  const intl = useIntl();

  const [sceneSegmentCreate] = useSceneSegmentCreate();
  const [sceneSegmentUpdate] = useSceneSegmentUpdate();
  const [sceneSegmentDestroy] = useSceneSegmentDestroy();
  const Toast = useToast();

  const isNew = segment === undefined;

  const schema = yup.object({
    title: yup.string().required("Title is required"),
    start_seconds: yup.number().min(0, "Start time must be >= 0").required(),
    end_seconds: yup
      .number()
      .min(0, "End time must be >= 0")
      .required()
      .test(
        "is-greater-than-start",
        "End time must be greater than start time",
        function (value) {
          return value > this.parent.start_seconds;
        }
      ),
  });

  const initialValues = useMemo(
    () => ({
      title: segment?.title ?? "",
      start_seconds:
        segment?.start_seconds ?? Math.round(getPlayerPosition() ?? 0),
      end_seconds:
        segment?.end_seconds ?? Math.round(getPlayerPosition() ?? 0) + 60,
    }),
    [segment]
  );

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  async function onSave(input: InputValues) {
    try {
      if (isNew) {
        await sceneSegmentCreate({
          variables: {
            input: {
              scene_id: sceneID,
              ...input,
            },
          },
        });
      } else {
        await sceneSegmentUpdate({
          variables: {
            input: {
              id: segment.id,
              ...input,
            },
          },
        });
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  async function onDelete() {
    if (isNew) return;

    try {
      await sceneSegmentDestroy({ variables: { id: segment.id } });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  const splitProps = {
    labelProps: {
      column: true,
      sm: 3,
    },
    fieldProps: {
      sm: 9,
    },
  };

  const { renderField } = formikUtils(intl, formik, splitProps);

  function renderTitleField() {
    const title = intl.formatMessage({ id: "title" });
    const control = (
      <Form.Control
        className="text-input"
        type="text"
        placeholder="Segment title"
        {...formik.getFieldProps("title")}
      />
    );

    return renderField("title", title, control);
  }

  function renderStartTimeField() {
    const { error } = formik.getFieldMeta("start_seconds");

    const title = intl.formatMessage({
      id: "time_start",
      defaultMessage: "Start Time",
    });
    const control = (
      <DurationInput
        value={formik.values.start_seconds}
        setValue={(v) => formik.setFieldValue("start_seconds", v)}
        onReset={() =>
          formik.setFieldValue("start_seconds", getPlayerPosition() ?? 0)
        }
        error={error}
      />
    );

    return renderField("start_seconds", title, control);
  }

  function renderEndTimeField() {
    const { error } = formik.getFieldMeta("end_seconds");

    const title = intl.formatMessage({ id: "time_end" });
    const control = (
      <>
        <DurationInput
          value={formik.values.end_seconds}
          setValue={(v) => formik.setFieldValue("end_seconds", v)}
          onReset={() =>
            formik.setFieldValue("end_seconds", getPlayerPosition() ?? 0)
          }
          error={error}
        />
        {formik.touched.end_seconds && formik.errors.end_seconds && (
          <Form.Control.Feedback type="invalid">
            {formik.errors.end_seconds}
          </Form.Control.Feedback>
        )}
      </>
    );

    return renderField("end_seconds", title, control);
  }

  return (
    <Form noValidate onSubmit={formik.handleSubmit}>
      <div className="form-container px-3">
        <h4>
          {isNew ? (
            <FormattedMessage
              id="actions.create_segment"
              defaultMessage="Create Segment"
            />
          ) : (
            <FormattedMessage
              id="actions.edit_segment"
              defaultMessage="Edit Segment"
            />
          )}
        </h4>
        {renderTitleField()}
        {renderStartTimeField()}
        {renderEndTimeField()}
      </div>
      <div className="buttons-container px-3">
        <div className="d-flex">
          <Button
            variant="primary"
            disabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
            onClick={() => formik.submitForm()}
          >
            <FormattedMessage id="actions.save" />
          </Button>
          <Button
            variant="secondary"
            type="button"
            onClick={onClose}
            className="ml-2"
          >
            <FormattedMessage id="actions.cancel" />
          </Button>
          {!isNew && (
            <Button
              variant="danger"
              className="ml-auto"
              onClick={() => onDelete()}
            >
              <FormattedMessage id="actions.delete" />
            </Button>
          )}
        </div>
      </div>
    </Form>
  );
};
