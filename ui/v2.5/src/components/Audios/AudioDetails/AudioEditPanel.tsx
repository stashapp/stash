import React, { useEffect, useState, useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form, Col, Row } from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ImageInput } from "src/components/Shared/ImageInput";
import { useToast } from "src/hooks/Toast";
import ImageUtils from "src/utils/image";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { IGroupEntry, AudioGroupTable } from "./AudioGroupTable";
import { galleryTitle } from "src/core/galleries";
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
import { Studio, StudioSelect } from "src/components/Studios/StudioSelect";
import { Gallery, GallerySelect } from "src/components/Galleries/GallerySelect";
import { Group } from "src/components/Groups/GroupSelect";
import { useTagsEdit } from "src/hooks/tagsEdit";
import {
  CustomFieldsInput,
  formatCustomFieldInput,
} from "src/components/Shared/CustomFields";
import cloneDeep from "lodash-es/cloneDeep";

interface IProps {
  audio: Partial<GQL.AudioDataFragment>;
  initialCoverImage?: string;
  isNew?: boolean;
  isVisible: boolean;
  onSubmit: (input: GQL.AudioCreateInput) => Promise<void>;
  onDelete?: () => void;
}

// NOTE - audio has no scrapers, no stash-ids and no cover generation, so the
// scrape menus and generate-thumbnail actions from SceneEditPanel are omitted.
export const AudioEditPanel: React.FC<IProps> = ({
  audio,
  initialCoverImage,
  isNew = false,
  isVisible,
  onSubmit,
  onDelete,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const [galleries, setGalleries] = useState<Gallery[]>([]);
  const [performers, setPerformers] = useState<Performer[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [studio, setStudio] = useState<Studio | null>(null);

  useEffect(() => {
    setGalleries(
      audio.galleries?.map((g) => ({
        id: g.id,
        title: galleryTitle(g),
        files: g.files,
        folder: g.folder,
      })) ?? []
    );
  }, [audio.galleries]);

  useEffect(() => {
    setPerformers(audio.performers ?? []);
  }, [audio.performers]);

  useEffect(() => {
    setGroups(audio.groups?.map((m) => m.group) ?? []);
  }, [audio.groups]);

  useEffect(() => {
    setStudio(audio.studio ?? null);
  }, [audio.studio]);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const schema = yup.object({
    title: yup.string().ensure(),
    code: yup.string().ensure(),
    urls: yupUniqueStringList(intl),
    date: yupDateString(intl),
    gallery_ids: yup.array(yup.string().required()).defined(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    groups: yup
      .array(
        yup.object({
          group_id: yup.string().required(),
          audio_index: yup.number().integer().nullable().defined(),
        })
      )
      .defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
    details: yup.string().ensure(),
    cover_image: yup.string().nullable().optional(),
    custom_fields: yup.object().required().defined(),
  });

  const initialValues = useMemo(
    () => ({
      title: audio.title ?? "",
      code: audio.code ?? "",
      urls: audio.urls ?? [],
      date: audio.date ?? "",
      gallery_ids: (audio.galleries ?? []).map((g) => g.id),
      studio_id: audio.studio?.id ?? null,
      performer_ids: (audio.performers ?? []).map((p) => p.id),
      groups: (audio.groups ?? []).map((m) => {
        return { group_id: m.group.id, audio_index: m.audio_index ?? null };
      }),
      tag_ids: (audio.tags ?? []).map((t) => t.id),
      details: audio.details ?? "",
      cover_image: initialCoverImage,
      custom_fields: cloneDeep(audio.custom_fields ?? {}),
    }),
    [audio, initialCoverImage]
  );

  type InputValues = yup.InferType<typeof schema>;

  const [customFieldsError, setCustomFieldsError] = useState<string>();

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

  function submit(values: InputValues) {
    const input = {
      ...schema.cast(values),
      custom_fields: formatCustomFieldInput(isNew, values.custom_fields),
    };
    onSave(input);
  }

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: submit,
  });

  const { tagsControl } = useTagsEdit(audio.tags, (ids) =>
    formik.setFieldValue("tag_ids", ids)
  );

  const coverImagePreview = useMemo(() => {
    const audioImage = audio.paths?.screenshot;
    const formImage = formik.values.cover_image;
    if (formImage === null && audioImage) {
      const audioImageURL = new URL(audioImage);
      audioImageURL.searchParams.set("default", "true");
      return audioImageURL.toString();
    } else if (formImage) {
      return formImage;
    }
    return audioImage;
  }, [formik.values.cover_image, audio.paths?.screenshot]);

  const groupEntries = useMemo(() => {
    return formik.values.groups
      .map((m) => {
        return {
          group: groups.find((mm) => mm.id === m.group_id),
          audio_index: m.audio_index,
        };
      })
      .filter((m) => m.group !== undefined) as IGroupEntry[];
  }, [formik.values.groups, groups]);

  function onSetGalleries(items: Gallery[]) {
    setGalleries(items);
    formik.setFieldValue(
      "gallery_ids",
      items.map((i) => i.id)
    );
  }

  function onSetPerformers(items: Performer[]) {
    setPerformers(items);
    formik.setFieldValue(
      "performer_ids",
      items.map((item) => item.id)
    );
  }

  function onSetStudio(item: Studio | null) {
    setStudio(item);
    formik.setFieldValue("studio_id", item ? item.id : null);
  }

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        if (formik.dirty) {
          formik.submitForm();
        }
      });
      Mousetrap.bind("d d", () => {
        if (onDelete) {
          onDelete();
        }
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");
      };
    }
  });

  function onCoverImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onImageLoad(imageData: string | null) {
    formik.setFieldValue("cover_image", imageData);
  }

  function onResetCover() {
    formik.setFieldValue("cover_image", null);
  }

  function onSetGroupEntries(input: IGroupEntry[]) {
    setGroups(input.map((m) => m.group));

    const newGroups = input.map((m) => ({
      group_id: m.group.id,
      audio_index: m.audio_index ?? null,
    }));

    formik.setFieldValue("groups", newGroups);
  }

  const image = useMemo(() => {
    if (coverImagePreview) {
      return (
        <img
          className="audio-cover"
          src={coverImagePreview}
          alt={intl.formatMessage({ id: "cover_image" })}
        />
      );
    }
  }, [coverImagePreview, intl]);

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
  const urlProps = isNew
    ? splitProps
    : {
        labelProps: {
          column: true,
          md: 3,
          lg: 12,
        },
        fieldProps: {
          md: 9,
          lg: 12,
        },
      };
  const { renderField, renderInputField, renderDateField, renderURLListField } =
    formikUtils(intl, formik, splitProps);

  function renderGalleriesField() {
    const title = intl.formatMessage({ id: "galleries" });
    const control = (
      <GallerySelect
        values={galleries}
        onSelect={(items) => onSetGalleries(items)}
        isMulti
      />
    );

    return renderField("gallery_ids", title, control);
  }

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
    const date = (() => {
      try {
        return schema.validateSyncAt("date", formik.values);
      } catch (_e) {
        return undefined;
      }
    })();

    const title = intl.formatMessage({ id: "performers" });
    const control = (
      <PerformerSelect
        isMulti
        onSelect={onSetPerformers}
        values={performers}
        ageFromDate={date}
      />
    );

    return renderField("performer_ids", title, control, fullWidthProps);
  }

  function renderGroupsField() {
    const title = intl.formatMessage({ id: "groups" });
    const control = (
      <AudioGroupTable value={groupEntries} onUpdate={onSetGroupEntries} />
    );

    return renderField("groups", title, control, fullWidthProps);
  }

  function renderTagsField() {
    const title = intl.formatMessage({ id: "tags" });
    return renderField("tag_ids", title, tagsControl(), fullWidthProps);
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
    <div id="audio-edit-details">
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
              disabled={
                (!isNew && !formik.dirty) ||
                !isEqual(formik.errors, {}) ||
                customFieldsError !== undefined
              }
              onClick={() => formik.submitForm()}
            >
              <FormattedMessage id="actions.save" />
            </Button>
            {onDelete && (
              <Button
                className="edit-button"
                variant="danger"
                onClick={() => onDelete()}
              >
                <FormattedMessage id="actions.delete" />
              </Button>
            )}
          </div>
        </Row>
        <Row className="form-container px-3">
          <Col lg={7} xl={12}>
            {renderInputField("title")}
            {renderInputField("code", "text", "audio_code")}

            {renderURLListField("urls", undefined, undefined, "urls", urlProps)}

            {renderDateField("date")}

            {renderGalleriesField()}
            {renderStudioField()}
            {renderPerformersField()}
            {renderGroupsField()}
            {renderTagsField()}
          </Col>
          <Col lg={5} xl={12}>
            {renderDetailsField()}
            <Form.Group controlId="cover_image">
              <Form.Label>
                <FormattedMessage id="cover_image" />
              </Form.Label>
              {image}
              <ImageInput
                isEditing
                onImageChange={onCoverImageChange}
                onImageURL={onImageLoad}
                onReset={
                  formik.values.cover_image ||
                  (formik.values.cover_image !== null &&
                    audio.paths?.screenshot)
                    ? () => onResetCover()
                    : undefined
                }
              />
            </Form.Group>

            <CustomFieldsInput
              values={formik.values.custom_fields}
              onChange={(v) => formik.setFieldValue("custom_fields", v)}
              error={customFieldsError}
              setError={(e) => setCustomFieldsError(e)}
            />
          </Col>
        </Row>
      </Form>
    </div>
  );
};

export default AudioEditPanel;
