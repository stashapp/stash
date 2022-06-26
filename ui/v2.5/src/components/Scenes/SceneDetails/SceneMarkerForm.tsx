import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Field, FieldProps, Form as FormikForm, Formik } from "formik";
import * as GQL from "src/core/generated-graphql";
import {
  useSceneMarkerCreate,
  useSceneMarkerUpdate,
  useSceneMarkerDestroy,
} from "src/core/StashService";
import {
  DurationInput,
  TagSelect,
  MarkerTitleSuggest,
} from "src/components/Shared";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import useToast from "src/hooks/Toast";

interface IFormFields {
  title: string;
  seconds: string;
  primaryTagId: string;
  tagIds: string[];
}

interface ISceneMarkerForm {
  sceneID: string;
  editingMarker?: GQL.SceneMarkerDataFragment;
  onClose: () => void;
}

export const SceneMarkerForm: React.FC<ISceneMarkerForm> = ({
  sceneID,
  editingMarker,
  onClose,
}) => {
  const [sceneMarkerCreate] = useSceneMarkerCreate();
  const [sceneMarkerUpdate] = useSceneMarkerUpdate();
  const [sceneMarkerDestroy] = useSceneMarkerDestroy();
  const Toast = useToast();

  const onSubmit = (values: IFormFields) => {
    const variables: GQL.SceneMarkerUpdateInput | GQL.SceneMarkerCreateInput = {
      title: values.title,
      seconds: parseFloat(values.seconds),
      scene_id: sceneID,
      primary_tag_id: values.primaryTagId,
      tag_ids: values.tagIds,
    };
    if (!editingMarker) {
      sceneMarkerCreate({ variables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    } else {
      const updateVariables = variables as GQL.SceneMarkerUpdateInput;
      updateVariables.id = editingMarker!.id;
      sceneMarkerUpdate({ variables: updateVariables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    }
  };

  const onDelete = () => {
    if (!editingMarker) return;

    sceneMarkerDestroy({ variables: { id: editingMarker.id } })
      .then(onClose)
      .catch((err) => Toast.error(err));
  };
  const renderTitleField = (fieldProps: FieldProps<string>) => (
    <MarkerTitleSuggest
      initialMarkerTitle={fieldProps.field.value}
      onChange={(query: string) =>
        fieldProps.form.setFieldValue("title", query)
      }
    />
  );

  const renderSecondsField = (fieldProps: FieldProps<string>) => (
    <DurationInput
      onValueChange={(s) => fieldProps.form.setFieldValue("seconds", s)}
      onReset={() =>
        fieldProps.form.setFieldValue(
          "seconds",
          Math.round(getPlayerPosition() ?? 0)
        )
      }
      numericValue={Number.parseInt(fieldProps.field.value ?? "0", 10)}
      mandatory
    />
  );

  const renderPrimaryTagField = (fieldProps: FieldProps<string>) => (
    <TagSelect
      onSelect={(tags) =>
        fieldProps.form.setFieldValue("primaryTagId", tags[0]?.id)
      }
      ids={fieldProps.field.value ? [fieldProps.field.value] : []}
      noSelectionString="Select/create tag..."
    />
  );

  const renderTagsField = (fieldProps: FieldProps<string[]>) => (
    <TagSelect
      isMulti
      onSelect={(tags) =>
        fieldProps.form.setFieldValue(
          "tagIds",
          tags.map((tag) => tag.id)
        )
      }
      ids={fieldProps.field.value}
      noSelectionString="Select/create tags..."
    />
  );

  const values: IFormFields = {
    title: editingMarker?.title ?? "",
    seconds: (
      editingMarker?.seconds ?? Math.round(getPlayerPosition() ?? 0)
    ).toString(),
    primaryTagId: editingMarker?.primary_tag.id ?? "",
    tagIds: editingMarker?.tags.map((tag) => tag.id) ?? [],
  };

  return (
    <Formik initialValues={values} onSubmit={onSubmit}>
      <FormikForm>
        <div>
          <Form.Group className="row">
            <Form.Label
              htmlFor="title"
              className="col-sm-3 col-md-2 col-xl-12 col-form-label"
            >
              Marker Title
            </Form.Label>
            <div className="col-sm-9 col-md-10 col-xl-12">
              <Field name="title">{renderTitleField}</Field>
            </div>
          </Form.Group>
          <Form.Group className="row">
            <Form.Label
              htmlFor="primaryTagId"
              className="col-sm-3 col-md-2 col-xl-12 col-form-label"
            >
              Primary Tag
            </Form.Label>
            <div className="col-sm-4 col-md-6 col-xl-12 mb-3 mb-sm-0 mb-xl-3">
              <Field name="primaryTagId">{renderPrimaryTagField}</Field>
            </div>
            <div className="col-sm-5 col-md-4 col-xl-12">
              <div className="row">
                <Form.Label
                  htmlFor="seconds"
                  className="col-sm-4 col-md-4 col-xl-12 col-form-label text-sm-right text-xl-left"
                >
                  Time
                </Form.Label>
                <div className="col-sm-8 col-xl-12">
                  <Field name="seconds">{renderSecondsField}</Field>
                </div>
              </div>
            </div>
          </Form.Group>
          <Form.Group className="row">
            <Form.Label
              htmlFor="tagIds"
              className="col-sm-3 col-md-2 col-xl-12 col-form-label"
            >
              Tags
            </Form.Label>
            <div className="col-sm-9 col-md-10 col-xl-12">
              <Field name="tagIds">{renderTagsField}</Field>
            </div>
          </Form.Group>
        </div>
        <div className="buttons-container row">
          <div className="col d-flex">
            <Button variant="primary" type="submit">
              Submit
            </Button>
            <Button
              variant="secondary"
              type="button"
              onClick={onClose}
              className="ml-2"
            >
              <FormattedMessage id="actions.cancel" />
            </Button>
            {editingMarker && (
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
      </FormikForm>
    </Formik>
  );
};
