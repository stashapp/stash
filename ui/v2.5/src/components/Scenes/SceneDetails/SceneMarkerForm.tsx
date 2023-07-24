import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { useFormik } from "formik";
import * as yup from "yup";
import * as GQL from "src/core/generated-graphql";
import {
  useSceneMarkerCreate,
  useSceneMarkerUpdate,
  useSceneMarkerDestroy,
} from "src/core/StashService";
import { DurationInput } from "src/components/Shared/DurationInput";
import { TagSelect, MarkerTitleSuggest } from "src/components/Shared/Select";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import { useToast } from "src/hooks/Toast";
import isEqual from "lodash-es/isEqual";

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

  const schema = yup.object({
    title: yup.string().ensure(),
    seconds: yup.number().required().integer(),
    primary_tag_id: yup.string().required(),
    tag_ids: yup.array(yup.string().required()).defined(),
  });

  const initialValues = {
    title: editingMarker?.title ?? "",
    seconds: editingMarker?.seconds ?? Math.round(getPlayerPosition() ?? 0),
    primary_tag_id: editingMarker?.primary_tag.id ?? "",
    tag_ids: editingMarker?.tags.map((tag) => tag.id) ?? [],
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    validationSchema: schema,
    enableReinitialize: true,
    onSubmit: (values) => onSave(values),
  });

  async function onSave(input: InputValues) {
    try {
      if (!editingMarker) {
        await sceneMarkerCreate({
          variables: {
            scene_id: sceneID,
            ...input,
          },
        });
      } else {
        await sceneMarkerUpdate({
          variables: {
            id: editingMarker.id,
            scene_id: sceneID,
            ...input,
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
    if (!editingMarker) return;

    try {
      await sceneMarkerDestroy({ variables: { id: editingMarker.id } });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  const primaryTagId = formik.values.primary_tag_id;

  return (
    <Form noValidate onSubmit={formik.handleSubmit}>
      <div>
        <Form.Group className="row">
          <Form.Label className="col-sm-3 col-md-2 col-xl-12 col-form-label">
            Marker Title
          </Form.Label>
          <div className="col-sm-9 col-md-10 col-xl-12">
            <MarkerTitleSuggest
              initialMarkerTitle={formik.values.title}
              onChange={(query: string) => formik.setFieldValue("title", query)}
            />
          </div>
        </Form.Group>
        <Form.Group className="row">
          <Form.Label className="col-sm-3 col-md-2 col-xl-12 col-form-label">
            Primary Tag
          </Form.Label>
          <div className="col-sm-4 col-md-6 col-xl-12 mb-3 mb-sm-0 mb-xl-3">
            <TagSelect
              onSelect={(tags) =>
                formik.setFieldValue("primary_tag_id", tags[0]?.id)
              }
              ids={primaryTagId ? [primaryTagId] : []}
              noSelectionString="Select/create tag..."
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.primary_tag_id}
            </Form.Control.Feedback>
          </div>
          <div className="col-sm-5 col-md-4 col-xl-12">
            <div className="row">
              <Form.Label className="col-sm-4 col-md-4 col-xl-12 col-form-label text-sm-right text-xl-left">
                Time
              </Form.Label>
              <div className="col-sm-8 col-xl-12">
                <DurationInput
                  onValueChange={(s) => formik.setFieldValue("seconds", s)}
                  onReset={() =>
                    formik.setFieldValue(
                      "seconds",
                      Math.round(getPlayerPosition() ?? 0)
                    )
                  }
                  numericValue={formik.values.seconds}
                  mandatory
                />
              </div>
            </div>
          </div>
        </Form.Group>
        <Form.Group className="row">
          <Form.Label className="col-sm-3 col-md-2 col-xl-12 col-form-label">
            Tags
          </Form.Label>
          <div className="col-sm-9 col-md-10 col-xl-12">
            <TagSelect
              isMulti
              onSelect={(tags) =>
                formik.setFieldValue(
                  "tag_ids",
                  tags.map((tag) => tag.id)
                )
              }
              ids={formik.values.tag_ids}
              noSelectionString="Select/create tags..."
            />
          </div>
        </Form.Group>
      </div>
      <div className="buttons-container row">
        <div className="col d-flex">
          <Button
            variant="primary"
            disabled={
              (editingMarker && !formik.dirty) || !isEqual(formik.errors, {})
            }
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
    </Form>
  );
};
