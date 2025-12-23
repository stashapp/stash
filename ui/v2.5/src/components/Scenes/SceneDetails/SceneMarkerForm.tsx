import React, { useEffect, useMemo, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useFormik } from "formik";
import * as yup from "yup";
import * as GQL from "src/core/generated-graphql";
import {
  useSceneMarkerCreate,
  useSceneMarkerUpdate,
  useSceneMarkerDestroy,
} from "src/core/StashService";
import { DurationInput } from "src/components/Shared/DurationInput";
import { MarkerTitleSuggest } from "src/components/Shared/Select";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import { useToast } from "src/hooks/Toast";
import isEqual from "lodash-es/isEqual";
import { formikUtils } from "src/utils/form";
import { yupFormikValidate } from "src/utils/yup";
import { Tag, TagSelect } from "src/components/Tags/TagSelect";

interface ISceneMarkerForm {
  sceneID: string;
  marker?: GQL.SceneMarkerDataFragment;
  onClose: () => void;
}

export const SceneMarkerForm: React.FC<ISceneMarkerForm> = ({
  sceneID,
  marker,
  onClose,
}) => {
  const intl = useIntl();

  const [sceneMarkerCreate] = useSceneMarkerCreate();
  const [sceneMarkerUpdate] = useSceneMarkerUpdate();
  const [sceneMarkerDestroy] = useSceneMarkerDestroy();
  const Toast = useToast();

  const [primaryTag, setPrimaryTag] = useState<Tag>();
  const [tags, setTags] = useState<Tag[]>([]);

  const isNew = marker === undefined;

  const schema = yup.object({
    title: yup.string().ensure(),
    seconds: yup.number().min(0).required(),
    end_seconds: yup
      .number()
      .min(0)
      .nullable()
      .defined()
      .test(
        "is-greater-than-seconds",
        intl.formatMessage({ id: "validation.end_time_before_start_time" }),
        function (value) {
          return value === null || value >= this.parent.seconds;
        }
      ),
    primary_tag_id: yup.string().required(),
    tag_ids: yup.array(yup.string().required()).defined(),
  });

  // useMemo to only run getPlayerPosition when the input marker actually changes
  const initialValues = useMemo(
    () => ({
      title: marker?.title ?? "",
      seconds: marker?.seconds ?? Math.round(getPlayerPosition() ?? 0),
      end_seconds: marker?.end_seconds ?? null,
      primary_tag_id: marker?.primary_tag.id ?? "",
      tag_ids: marker?.tags.map((tag) => tag.id) ?? [],
    }),
    [marker]
  );

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  function onSetPrimaryTag(item: Tag) {
    setPrimaryTag(item);
    formik.setFieldValue("primary_tag_id", item.id);
  }

  function onSetTags(items: Tag[]) {
    setTags(items);
    formik.setFieldValue(
      "tag_ids",
      items.map((item) => item.id)
    );
  }

  useEffect(() => {
    setPrimaryTag(
      marker?.primary_tag
        ? { ...marker.primary_tag, aliases: [], stash_ids: [] }
        : undefined
    );
  }, [marker?.primary_tag]);

  useEffect(() => {
    setTags(
      marker?.tags.map((t) => ({
        ...t,
        aliases: [],
        stash_ids: [],
      })) ?? []
    );
  }, [marker?.tags]);

  async function onSave(input: InputValues) {
    try {
      if (isNew) {
        await sceneMarkerCreate({
          variables: {
            scene_id: sceneID,
            ...input,
            // undefined means setting to null, not omitting the field
            end_seconds: input.end_seconds ?? null,
          },
        });
      } else {
        await sceneMarkerUpdate({
          variables: {
            id: marker.id,
            scene_id: sceneID,
            ...input,
            // undefined means setting to null, not omitting the field
            end_seconds: input.end_seconds ?? null,
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
      await sceneMarkerDestroy({ variables: { id: marker.id } });
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
  const fullWidthProps = {
    labelProps: {
      column: true,
      sm: 3,
      xl: 12,
    },
    fieldProps: {
      sm: 9,
      xl: 12,
    },
  };
  const { renderField } = formikUtils(intl, formik, splitProps);

  function renderTitleField() {
    const title = intl.formatMessage({ id: "title" });
    const control = (
      <MarkerTitleSuggest
        initialMarkerTitle={formik.values.title}
        onChange={(v) => formik.setFieldValue("title", v)}
      />
    );

    return renderField("title", title, control);
  }

  function renderPrimaryTagField() {
    const title = intl.formatMessage({ id: "primary_tag" });
    const control = (
      <>
        <TagSelect
          onSelect={(t) => onSetPrimaryTag(t[0])}
          values={primaryTag ? [primaryTag] : []}
          hoverPlacement="right"
        />
        {formik.touched.primary_tag_id && (
          <Form.Control.Feedback type="invalid">
            {formik.errors.primary_tag_id}
          </Form.Control.Feedback>
        )}
      </>
    );

    return renderField("primary_tag_id", title, control);
  }

  function renderTimeField() {
    const { error } = formik.getFieldMeta("seconds");

    const title = intl.formatMessage({ id: "time" });
    const control = (
      <DurationInput
        value={formik.values.seconds}
        setValue={(v) => formik.setFieldValue("seconds", v)}
        onReset={() =>
          formik.setFieldValue("seconds", getPlayerPosition() ?? 0)
        }
        error={error}
      />
    );

    return renderField("seconds", title, control);
  }

  function renderEndTimeField() {
    const { error } = formik.getFieldMeta("end_seconds");

    const title = intl.formatMessage({ id: "time_end" });
    const control = (
      <>
        <DurationInput
          value={formik.values.end_seconds}
          setValue={(v) => formik.setFieldValue("end_seconds", v ?? null)}
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

  function renderTagsField() {
    const title = intl.formatMessage({ id: "tags" });
    const control = (
      <TagSelect
        isMulti
        onSelect={onSetTags}
        values={tags}
        hoverPlacement="right"
      />
    );

    return renderField("tag_ids", title, control, fullWidthProps);
  }

  return (
    <Form noValidate onSubmit={formik.handleSubmit}>
      <div className="form-container px-3">
        {renderTitleField()}
        {renderPrimaryTagField()}
        {renderTimeField()}
        {renderEndTimeField()}
        {renderTagsField()}
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
