import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useFormik } from "formik";
import * as yup from "yup";
import * as GQL from "src/core/generated-graphql";
import {
  useGalleryChapterCreate,
  useGalleryChapterUpdate,
  useGalleryChapterDestroy,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import isEqual from "lodash-es/isEqual";

interface IGalleryChapterForm {
  galleryID: string;
  chapter?: GQL.GalleryChapterDataFragment;
  onClose: () => void;
}

export const GalleryChapterForm: React.FC<IGalleryChapterForm> = ({
  galleryID,
  chapter,
  onClose,
}) => {
  const intl = useIntl();

  const [galleryChapterCreate] = useGalleryChapterCreate();
  const [galleryChapterUpdate] = useGalleryChapterUpdate();
  const [galleryChapterDestroy] = useGalleryChapterDestroy();
  const Toast = useToast();

  const isNew = chapter === undefined;

  const schema = yup.object({
    title: yup.string().ensure(),
    image_index: yup
      .number()
      .integer()
      .required()
      .moreThan(0)
      .label(intl.formatMessage({ id: "image_index" })),
  });

  const initialValues = {
    title: chapter?.title ?? "",
    image_index: chapter?.image_index ?? 1,
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
      if (isNew) {
        await galleryChapterCreate({
          variables: {
            gallery_id: galleryID,
            ...input,
          },
        });
      } else {
        await galleryChapterUpdate({
          variables: {
            id: chapter.id,
            gallery_id: galleryID,
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
    if (isNew) return;

    try {
      await galleryChapterDestroy({ variables: { id: chapter.id } });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  return (
    <Form noValidate onSubmit={formik.handleSubmit}>
      <div>
        <Form.Group>
          <Form.Label>
            <FormattedMessage id="title" />
          </Form.Label>

          <Form.Control
            className="text-input"
            placeholder={intl.formatMessage({ id: "title" })}
            isInvalid={!!formik.errors.title}
            {...formik.getFieldProps("title")}
          />
          <Form.Control.Feedback type="invalid">
            {formik.errors.title}
          </Form.Control.Feedback>
        </Form.Group>

        <Form.Group>
          <Form.Label>
            <FormattedMessage id="image_index" />
          </Form.Label>

          <Form.Control
            className="text-input"
            placeholder={intl.formatMessage({ id: "image_index" })}
            isInvalid={!!formik.errors.image_index}
            {...formik.getFieldProps("image_index")}
          />
          <Form.Control.Feedback type="invalid">
            {formik.errors.image_index}
          </Form.Control.Feedback>
        </Form.Group>
      </div>
      <div className="buttons-container row">
        <div className="col d-flex">
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
