import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { Form } from "react-bootstrap";
import ImageUtils from "src/utils/image";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import Mousetrap from "mousetrap";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import isEqual from "lodash-es/isEqual";
import { useToast } from "src/hooks/Toast";
import { handleUnsavedChanges } from "src/utils/navigation";
import { formikUtils } from "src/utils/form";
import { yupFormikValidate, yupUniqueAliases } from "src/utils/yup";
import { Tag, TagSelect } from "../TagSelect";

interface ITagEditPanel {
  tag: Partial<GQL.TagDataFragment>;
  onSubmit: (tag: GQL.TagCreateInput) => Promise<void>;
  onCancel: () => void;
  onDelete: () => void;
  setImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const TagEditPanel: React.FC<ITagEditPanel> = ({
  tag,
  onSubmit,
  onCancel,
  onDelete,
  setImage,
  setEncodingImage,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const isNew = tag.id === undefined;

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [childTags, setChildTags] = useState<Tag[]>([]);
  const [parentTags, setParentTags] = useState<Tag[]>([]);

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yupUniqueAliases(intl, "name"),
    description: yup.string().ensure(),
    parent_ids: yup.array(yup.string().required()).defined(),
    child_ids: yup.array(yup.string().required()).defined(),
    ignore_auto_tag: yup.boolean().defined(),
    image: yup.string().nullable().optional(),
  });

  const initialValues = {
    name: tag?.name ?? "",
    aliases: tag?.aliases ?? [],
    description: tag?.description ?? "",
    parent_ids: (tag?.parents ?? []).map((t) => t.id),
    child_ids: (tag?.children ?? []).map((t) => t.id),
    ignore_auto_tag: tag?.ignore_auto_tag ?? false,
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  function onSetParentTags(items: Tag[]) {
    setParentTags(items);
    formik.setFieldValue(
      "parent_ids",
      items.map((item) => item.id)
    );
  }

  function onSetChildTags(items: Tag[]) {
    setChildTags(items);
    formik.setFieldValue(
      "child_ids",
      items.map((item) => item.id)
    );
  }

  useEffect(() => {
    setParentTags(tag.parents ?? []);
  }, [tag.parents]);

  useEffect(() => {
    setChildTags(tag.children ?? []);
  }, [tag.children]);

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("s s", () => {
      if (formik.dirty) {
        formik.submitForm();
      }
    });

    return () => {
      Mousetrap.unbind("s s");
    };
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

  const encodingImage = ImageUtils.usePasteImage(onImageLoad);

  useEffect(() => {
    setImage(formik.values.image);
  }, [formik.values.image, setImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

  function onImageLoad(imageData: string | null) {
    formik.setFieldValue("image", imageData);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  const { renderField, renderInputField, renderStringListField } = formikUtils(
    intl,
    formik
  );

  function renderParentTagsField() {
    const title = intl.formatMessage({ id: "parent_tags" });
    const control = (
      <TagSelect
        isMulti
        onSelect={onSetParentTags}
        values={parentTags}
        excludeIds={[...(tag?.id ? [tag.id] : []), ...formik.values.child_ids]}
        creatable={false}
        hoverPlacement="right"
      />
    );

    return renderField("parent_ids", title, control);
  }

  function renderSubTagsField() {
    const title = intl.formatMessage({ id: "sub_tags" });
    const control = (
      <TagSelect
        isMulti
        onSelect={onSetChildTags}
        values={childTags}
        excludeIds={[...(tag?.id ? [tag.id] : []), ...formik.values.parent_ids]}
        creatable={false}
        hoverPlacement="right"
      />
    );

    return renderField("child_ids", title, control);
  }

  if (isLoading) return <LoadingIndicator />;

  // TODO: CSS class
  return (
    <div>
      {isNew && (
        <h2>
          <FormattedMessage
            id="actions.add_entity"
            values={{ entityType: intl.formatMessage({ id: "tag" }) }}
          />
        </h2>
      )}

      <Prompt
        when={formik.dirty}
        message={(location, action) => {
          // Check if it's a redirect after tag creation
          if (action === "PUSH" && location.pathname.startsWith("/tags/")) {
            return true;
          }

          return handleUnsavedChanges(intl, "tags", tag.id)(location);
        }}
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="tag-edit">
        {renderInputField("name")}
        {renderStringListField("aliases")}
        {renderInputField("description", "textarea")}
        {renderParentTagsField()}
        {renderSubTagsField()}
        <hr />
        {renderInputField("ignore_auto_tag", "checkbox")}
      </Form>

      <DetailsEditNavbar
        objectName={tag?.name ?? intl.formatMessage({ id: "tag" })}
        classNames="col-xl-9 mt-3"
        isNew={isNew}
        isEditing
        onToggleEdit={onCancel}
        onSave={formik.handleSubmit}
        saveDisabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
        onImageChange={onImageChange}
        onImageChangeURL={onImageLoad}
        onClearImage={() => onImageLoad(null)}
        onDelete={onDelete}
        acceptSVG
      />
    </div>
  );
};
