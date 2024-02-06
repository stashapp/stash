import React, { useEffect, useState } from "react";
import { Button, Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { ConfigurationContext } from "src/hooks/Config";
import isEqual from "lodash-es/isEqual";
import {
  yupDateString,
  yupFormikValidate,
  yupUniqueStringList,
} from "src/utils/yup";
import {
  Performer,
  PerformerSelect,
} from "src/components/Performers/PerformerSelect";
import { formikUtils } from "src/utils/form";
import { Tag, TagSelect } from "src/components/Tags/TagSelect";
import { Studio, StudioSelect } from "src/components/Studios/StudioSelect";

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

  const [performers, setPerformers] = useState<Performer[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [studio, setStudio] = useState<Studio | null>(null);

  const schema = yup.object({
    title: yup.string().ensure(),
    code: yup.string().ensure(),
    urls: yupUniqueStringList(intl),
    date: yupDateString(intl),
    details: yup.string().ensure(),
    photographer: yup.string().ensure(),
    rating100: yup.number().integer().nullable().defined(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
  });

  const initialValues = {
    title: image.title ?? "",
    code: image.code ?? "",
    urls: image?.urls ?? [],
    date: image?.date ?? "",
    details: image.details ?? "",
    photographer: image.photographer ?? "",
    rating100: image.rating100 ?? null,
    studio_id: image.studio?.id ?? null,
    performer_ids: (image.performers ?? []).map((p) => p.id),
    tag_ids: (image.tags ?? []).map((t) => t.id),
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  function onSetPerformers(items: Performer[]) {
    setPerformers(items);
    formik.setFieldValue(
      "performer_ids",
      items.map((item) => item.id)
    );
  }

  function onSetTags(items: Tag[]) {
    setTags(items);
    formik.setFieldValue(
      "tag_ids",
      items.map((item) => item.id)
    );
  }

  function onSetStudio(item: Studio | null) {
    setStudio(item);
    formik.setFieldValue("studio_id", item ? item.id : null);
  }

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  useEffect(() => {
    setPerformers(image.performers ?? []);
  }, [image.performers]);

  useEffect(() => {
    setTags(image.tags ?? []);
  }, [image.tags]);

  useEffect(() => {
    setStudio(image.studio ?? null);
  }, [image.studio]);

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

  if (isLoading) return <LoadingIndicator />;

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
  const {
    renderField,
    renderInputField,
    renderDateField,
    renderRatingField,
    renderURLListField,
  } = formikUtils(intl, formik, splitProps);

  function renderStudioField() {
    const title = intl.formatMessage({ id: "studio" });
    const control = (
      <StudioSelect
        onSelect={(items) => onSetStudio(items.length > 0 ? items[0] : null)}
        values={studio ? [studio] : []}
      />
    );

    return renderField("studio_id", title, control);
  }

  function renderPerformersField() {
    const title = intl.formatMessage({ id: "performers" });
    const control = (
      <PerformerSelect isMulti onSelect={onSetPerformers} values={performers} />
    );

    return renderField("performer_ids", title, control, fullWidthProps);
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

  function renderDetailsField() {
    const props = {
      labelProps: {
        column: true,
        sm: 3,
        lg: 12,
      },
      fieldProps: {
        sm: 9,
        lg: 12,
      },
    };

    return renderInputField("details", "textarea", "details", props);
  }

  return (
    <div id="image-edit-details">
      <Prompt
        when={formik.dirty}
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />

      <Form noValidate onSubmit={formik.handleSubmit}>
        <Row className="form-container edit-buttons-container px-3 pt-3">
          <div className="edit-buttons mb-3 pl-0">
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
        </Row>
        <Row className="form-container px-3">
          <Col lg={7} xl={12}>
            {renderInputField("title")}
            {renderInputField("code", "text", "scene_code")}

            {renderURLListField("urls")}

            {renderDateField("date")}
            {renderInputField("photographer")}
            {renderRatingField("rating100", "rating")}

            {renderStudioField()}
            {renderPerformersField()}
            {renderTagsField()}
          </Col>
          <Col lg={5} xl={12}>
            {renderDetailsField()}
          </Col>
        </Row>
      </Form>
    </div>
  );
};
