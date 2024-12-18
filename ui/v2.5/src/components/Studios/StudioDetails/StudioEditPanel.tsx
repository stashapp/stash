import React, { useEffect, useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { Form } from "react-bootstrap";
import ImageUtils from "src/utils/image";
import { getStashIDs } from "src/utils/stashIds";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import isEqual from "lodash-es/isEqual";
import { useToast } from "src/hooks/Toast";
import { handleUnsavedChanges } from "src/utils/navigation";
import { formikUtils } from "src/utils/form";
import { yupFormikValidate, yupUniqueAliases } from "src/utils/yup";
import { Studio, StudioSelect } from "../StudioSelect";
import { useTagsEdit } from "src/hooks/tagsEdit";

interface IStudioEditPanel {
  studio: Partial<GQL.StudioDataFragment>;
  onSubmit: (studio: GQL.StudioCreateInput) => Promise<void>;
  onCancel: () => void;
  onDelete: () => void;
  setImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const StudioEditPanel: React.FC<IStudioEditPanel> = ({
  studio,
  onSubmit,
  onCancel,
  onDelete,
  setImage,
  setEncodingImage,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const isNew = studio.id === undefined;

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [parentStudio, setParentStudio] = useState<Studio | null>(null);

  const schema = yup.object({
    name: yup.string().required(),
    url: yup.string().ensure(),
    details: yup.string().ensure(),
    parent_id: yup.string().required().nullable(),
    aliases: yupUniqueAliases(intl, "name"),
    tag_ids: yup.array(yup.string().required()).defined(),
    ignore_auto_tag: yup.boolean().defined(),
    stash_ids: yup.mixed<GQL.StashIdInput[]>().defined(),
    image: yup.string().nullable().optional(),
  });

  const initialValues = {
    id: studio.id,
    name: studio.name ?? "",
    url: studio.url ?? "",
    details: studio.details ?? "",
    parent_id: studio.parent_studio?.id ?? null,
    aliases: studio.aliases ?? [],
    tag_ids: (studio.tags ?? []).map((t) => t.id),
    ignore_auto_tag: studio.ignore_auto_tag ?? false,
    stash_ids: getStashIDs(studio.stash_ids),
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  const { tagsControl } = useTagsEdit(studio.tags, (ids) =>
    formik.setFieldValue("tag_ids", ids)
  );

  function onSetParentStudio(item: Studio | null) {
    setParentStudio(item);
    formik.setFieldValue("parent_id", item ? item.id : null);
  }

  const encodingImage = ImageUtils.usePasteImage((imageData) =>
    formik.setFieldValue("image", imageData)
  );

  useEffect(() => {
    setParentStudio(
      studio.parent_studio
        ? {
            id: studio.parent_studio.id,
            name: studio.parent_studio.name,
            aliases: [],
          }
        : null
    );
  }, [studio.parent_studio]);

  useEffect(() => {
    setImage(formik.values.image);
  }, [formik.values.image, setImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

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

  function onImageLoad(imageData: string | null) {
    formik.setFieldValue("image", imageData);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  const {
    renderField,
    renderInputField,
    renderStringListField,
    renderStashIDsField,
  } = formikUtils(intl, formik);

  function renderParentStudioField() {
    const title = intl.formatMessage({ id: "parent_studio" });
    const control = (
      <StudioSelect
        onSelect={(items) =>
          onSetParentStudio(items.length > 0 ? items[0] : null)
        }
        values={parentStudio ? [parentStudio] : []}
      />
    );

    return renderField("parent_id", title, control);
  }

  function renderTagsField() {
    const title = intl.formatMessage({ id: "tags" });
    return renderField("tag_ids", title, tagsControl());
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <>
      <Prompt
        when={formik.dirty}
        message={(location, action) => {
          // Check if it's a redirect after studio creation
          if (action === "PUSH" && location.pathname.startsWith("/studios/"))
            return true;

          return handleUnsavedChanges(intl, "studios", studio.id)(location);
        }}
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="studio-edit">
        {renderInputField("name")}
        {renderStringListField("aliases")}
        {renderInputField("url")}
        {renderInputField("details", "textarea")}
        {renderParentStudioField()}
        {renderTagsField()}
        {renderStashIDsField("stash_ids", "studios")}
        <hr />
        {renderInputField("ignore_auto_tag", "checkbox")}
      </Form>

      <DetailsEditNavbar
        objectName={studio?.name ?? intl.formatMessage({ id: "studio" })}
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
    </>
  );
};
