import React, { useEffect, useState } from "react";
import { Button, Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
} from "src/components/Shared/Select";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { URLField } from "src/components/Shared/URLField";
import { useToast } from "src/hooks/Toast";
import FormUtils from "src/utils/form";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { ConfigurationContext } from "src/hooks/Config";
import isEqual from "lodash-es/isEqual";
import { DateInput } from "src/components/Shared/DateInput";

interface IProps {
  image: GQL.ImageDataFragment;
  isVisible: boolean;
  onSubmit: (input: GQL.ImageUpdateInput) => Promise<void>;
  onDelete: () => void;
}

export const ImageEditPanel: React.FC<IProps> = ({
  image,
  isVisible,
  onSubmit,
  onDelete,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { configuration } = React.useContext(ConfigurationContext);

  const schema = yup.object({
    title: yup.string().ensure(),
    url: yup.string().ensure(),
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
    rating100: yup.number().nullable().defined(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
  });

  const initialValues = {
    title: image.title ?? "",
    url: image?.url ?? "",
    date: image?.date ?? "",
    rating100: image.rating100 ?? null,
    studio_id: image.studio?.id ?? null,
    performer_ids: (image.performers ?? []).map((p) => p.id),
    tag_ids: (image.tags ?? []).map((t) => t.id),
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  useRatingKeybinds(
    true,
    configuration?.ui?.ratingSystemOptions?.type,
    setRating
  );

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        if (formik.dirty) {
          formik.submitForm();
        }
      });
      Mousetrap.bind("d d", () => {
        onDelete();
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");
      };
    }
  });

  async function onSave(input: InputValues) {
    setIsLoading(true);
    try {
      await onSubmit({
        id: image.id,
        ...input,
      });
      formik.resetForm();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  function renderTextField(field: string, title: string, placeholder?: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        {FormUtils.renderLabel({
          title,
        })}
        <Col xs={9}>
          <Form.Control
            className="text-input"
            placeholder={placeholder ?? title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
          <Form.Control.Feedback type="invalid">
            {formik.getFieldMeta(field).error}
          </Form.Control.Feedback>
        </Col>
      </Form.Group>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="image-edit-details">
      <Prompt
        when={formik.dirty}
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />

      <Form noValidate onSubmit={formik.handleSubmit}>
        <div className="form-container row px-3 pt-3">
          <div className="col edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={!formik.dirty || !isEqual(formik.errors, {})}
              onClick={() => formik.submitForm()}
            >
              <FormattedMessage id="actions.save" />
            </Button>
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => onDelete()}
            >
              <FormattedMessage id="actions.delete" />
            </Button>
          </div>
        </div>
        <div className="form-container row px-3">
          <div className="col-12 col-lg-6 col-xl-12">
            {renderTextField("title", intl.formatMessage({ id: "title" }))}
            <Form.Group controlId="url" as={Row}>
              <Col xs={3} className="pr-0 url-label">
                <Form.Label className="col-form-label">
                  <FormattedMessage id="url" />
                </Form.Label>
              </Col>
              <Col xs={9}>
                <URLField
                  {...formik.getFieldProps("url")}
                  onScrapeClick={() => {}}
                  urlScrapable={() => {
                    return false;
                  }}
                  isInvalid={!!formik.getFieldMeta("url").error}
                />
              </Col>
            </Form.Group>
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
            <Form.Group controlId="rating" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "rating" }),
              })}
              <Col xs={9}>
                <RatingSystem
                  value={formik.values.rating100 ?? undefined}
                  onSetRating={(value) =>
                    formik.setFieldValue("rating100", value ?? null)
                  }
                />
              </Col>
            </Form.Group>
            <Form.Group controlId="studio" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "studio" }),
              })}
              <Col xs={9}>
                <StudioSelect
                  onSelect={(items) =>
                    formik.setFieldValue(
                      "studio_id",
                      items.length > 0 ? items[0]?.id : null
                    )
                  }
                  ids={formik.values.studio_id ? [formik.values.studio_id] : []}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="performers" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "performers" }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <PerformerSelect
                  isMulti
                  onSelect={(items) =>
                    formik.setFieldValue(
                      "performer_ids",
                      items.map((item) => item.id)
                    )
                  }
                  ids={formik.values.performer_ids}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="tags" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "tags" }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <TagSelect
                  isMulti
                  onSelect={(items) =>
                    formik.setFieldValue(
                      "tag_ids",
                      items.map((item) => item.id)
                    )
                  }
                  ids={formik.values.tag_ids}
                  hoverPlacement="right"
                />
              </Col>
            </Form.Group>
          </div>
        </div>
      </Form>
    </div>
  );
};
