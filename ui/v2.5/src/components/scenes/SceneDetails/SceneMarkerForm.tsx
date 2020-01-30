import React from "react";
import { Button, Form } from "react-bootstrap";
import { Field, FieldProps, Form as FormikForm, Formik } from "formik";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import {
  DurationInput,
  TagSelect,
  MarkerTitleSuggest
} from "src/components/Shared";
import { useToast } from "src/hooks";

interface IFormFields {
  title: string;
  seconds: string;
  primaryTagId: string;
  tagIds: string[];
}

interface ISceneMarkerForm {
  sceneID: string;
  editingMarker?: GQL.SceneMarkerDataFragment;
  playerPosition?: number;
  onClose: () => void;
}

export const SceneMarkerForm: React.FC<ISceneMarkerForm> = ({
  sceneID,
  editingMarker,
  playerPosition,
  onClose
}) => {
  const [sceneMarkerCreate] = StashService.useSceneMarkerCreate();
  const [sceneMarkerUpdate] = StashService.useSceneMarkerUpdate();
  const [sceneMarkerDestroy] = StashService.useSceneMarkerDestroy();
  const Toast = useToast();

  const onSubmit = (values: IFormFields) => {
    const variables: GQL.SceneMarkerUpdateInput | GQL.SceneMarkerCreateInput = {
      title: values.title,
      seconds: parseFloat(values.seconds),
      scene_id: sceneID,
      primary_tag_id: values.primaryTagId,
      tag_ids: values.tagIds
    };
    if (!editingMarker) {
      sceneMarkerCreate({ variables })
        .then(onClose)
        .catch(err => Toast.error(err));
    } else {
      const updateVariables = variables as GQL.SceneMarkerUpdateInput;
      updateVariables.id = editingMarker!.id;
      sceneMarkerUpdate({ variables: updateVariables })
        .then(onClose)
        .catch(err => Toast.error(err));
    }
  };

  const onDelete = () => {
    if (!editingMarker) return;

    sceneMarkerDestroy({ variables: { id: editingMarker.id } })
      .then(onClose)
      .catch(err => Toast.error(err));
  };
  const renderTitleField = (fieldProps: FieldProps<string>) => (
    <div className="col-10">
      <MarkerTitleSuggest
        initialMarkerTitle={fieldProps.field.value}
        onChange={(query: string) =>
          fieldProps.form.setFieldValue("title", query)
        }
      />
    </div>
  );

  const renderSecondsField = (fieldProps: FieldProps<string>) => (
    <div className="col-3">
      <DurationInput
        onValueChange={s => fieldProps.form.setFieldValue("seconds", s)}
        onReset={() =>
          fieldProps.form.setFieldValue(
            "seconds",
            Math.round(playerPosition ?? 0)
          )
        }
        numericValue={Number.parseInt(fieldProps.field.value ?? "0", 10)}
      />
    </div>
  );

  const renderPrimaryTagField = (fieldProps: FieldProps<string>) => (
    <TagSelect
      onSelect={tags =>
        fieldProps.form.setFieldValue("primaryTagId", tags[0]?.id)
      }
      ids={fieldProps.field.value ? [fieldProps.field.value] : []}
      noSelectionString="Select or create tag..."
    />
  );

  const renderTagsField = (fieldProps: FieldProps<string[]>) => (
    <TagSelect
      isMulti
      onSelect={tags =>
        fieldProps.form.setFieldValue(
          "tagIds",
          tags.map(tag => tag.id)
        )
      }
      ids={fieldProps.field.value}
      noSelectionString="Select or create tags..."
    />
  );

  const values: IFormFields = {
    title: editingMarker?.title ?? "",
    seconds: (
      editingMarker?.seconds ?? Math.round(playerPosition ?? 0)
    ).toString(),
    primaryTagId: editingMarker?.primary_tag.id ?? "",
    tagIds: editingMarker?.tags.map(tag => tag.id) ?? []
  };

  return (
    <Formik initialValues={values} onSubmit={onSubmit}>
      <FormikForm>
        <div>
          <Form.Group className="row">
            <Form.Label htmlFor="title" className="col-2">
              Scene Marker Title
            </Form.Label>
            <Field name="title">{renderTitleField}</Field>
          </Form.Group>
          <Form.Group className="row">
            <Form.Label htmlFor="primaryTagId" className="col-2">
              Primary Tag
            </Form.Label>
            <div className="col-6">
              <Field name="primaryTagId">{renderPrimaryTagField}</Field>
            </div>
            <Form.Label htmlFor="seconds" className="col-1">
              Time
            </Form.Label>
            <Field name="seconds">{renderSecondsField}</Field>
          </Form.Group>
          <Form.Group className="row">
            <Form.Label htmlFor="tagIds" className="col-2">
              Tags
            </Form.Label>
            <div className="col-10">
              <Field name="tagIds">{renderTagsField}</Field>
            </div>
          </Form.Group>
        </div>
        <div className="buttons-container row">
          <Button variant="primary" type="submit">
            Submit
          </Button>
          <Button variant="secondary" type="button" onClick={onClose} className="ml-2">
            Cancel
          </Button>
          {editingMarker && (
            <Button
              variant="danger"
              className="ml-auto"
              onClick={() => onDelete()}
            >
              Delete
            </Button>
          )}
        </div>
      </FormikForm>
    </Formik>
  );
};
