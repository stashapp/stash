import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Form as FormikForm, Formik } from "formik";
import * as yup from "yup";
import * as GQL from "src/core/generated-graphql";
import {
  useGalleryChapterCreate,
  useGalleryChapterUpdate,
  useGalleryChapterDestroy,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import isEqual from "lodash-es/isEqual";

interface IFormFields {
  title: string;
  imageIndex: number;
}

interface IGalleryChapterForm {
  galleryID: string;
  editingChapter?: GQL.GalleryChapterDataFragment;
  onClose: () => void;
}

export const GalleryChapterForm: React.FC<IGalleryChapterForm> = ({
  galleryID,
  editingChapter,
  onClose,
}) => {
  const intl = useIntl();

  const [galleryChapterCreate] = useGalleryChapterCreate();
  const [galleryChapterUpdate] = useGalleryChapterUpdate();
  const [galleryChapterDestroy] = useGalleryChapterDestroy();
  const Toast = useToast();

  const schema = yup.object({
    title: yup.string().ensure(),
    imageIndex: yup
      .number()
      .required()
      .label(intl.formatMessage({ id: "image_index" }))
      .moreThan(0),
  });

  const onSubmit = (values: IFormFields) => {
    const variables:
      | GQL.GalleryChapterUpdateInput
      | GQL.GalleryChapterCreateInput = {
      title: values.title,
      image_index: values.imageIndex,
      gallery_id: galleryID,
    };

    if (!editingChapter) {
      galleryChapterCreate({ variables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    } else {
      const updateVariables = variables as GQL.GalleryChapterUpdateInput;
      updateVariables.id = editingChapter!.id;
      galleryChapterUpdate({ variables: updateVariables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    }
  };

  const onDelete = () => {
    if (!editingChapter) return;

    galleryChapterDestroy({ variables: { id: editingChapter.id } })
      .then(onClose)
      .catch((err) => Toast.error(err));
  };

  const values: IFormFields = {
    title: editingChapter?.title ?? "",
    imageIndex: editingChapter?.image_index ?? 1,
  };

  return (
    <Formik
      initialValues={values}
      onSubmit={onSubmit}
      validationSchema={schema}
    >
      {(formik) => (
        <FormikForm>
          <div>
            <Form.Group>
              <Form.Label>
                <FormattedMessage id="title" />
              </Form.Label>

              <Form.Control
                className="text-input"
                placeholder={intl.formatMessage({ id: "title" })}
                {...formik.getFieldProps("title")}
                isInvalid={!!formik.getFieldMeta("title").error}
              />
              <Form.Control.Feedback type="invalid">
                {formik.getFieldMeta("title").error}
              </Form.Control.Feedback>
            </Form.Group>

            <Form.Group>
              <Form.Label>
                <FormattedMessage id="image_index" />
              </Form.Label>

              <Form.Control
                className="text-input"
                placeholder={intl.formatMessage({ id: "image_index" })}
                {...formik.getFieldProps("imageIndex")}
                isInvalid={!!formik.getFieldMeta("imageIndex").error}
              />
              <Form.Control.Feedback type="invalid">
                {formik.getFieldMeta("imageIndex").error}
              </Form.Control.Feedback>
            </Form.Group>
          </div>
          <div className="buttons-container row">
            <div className="col d-flex">
              <Button
                variant="primary"
                disabled={
                  (editingChapter && !formik.dirty) ||
                  !isEqual(formik.errors, {})
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
              {editingChapter && (
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
      )}
    </Formik>
  );
};
