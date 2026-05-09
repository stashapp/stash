import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Prompt } from "react-router-dom";
import { Button, Dropdown, Form, Col, Row, SplitButton } from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  queryScrapeGallery,
  queryScrapeGalleryURL,
  useListGalleryScrapers,
  mutateReloadScrapers,
} from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { useFormik } from "formik";
import { GalleryScrapeDialog } from "./GalleryScrapeDialog";
import isEqual from "lodash-es/isEqual";
import { handleUnsavedChanges } from "src/utils/navigation";
import {
  Performer,
  PerformerSelect,
} from "src/components/Performers/PerformerSelect";
import {
  yupDateString,
  yupFormikValidate,
  yupUniqueStringList,
} from "src/utils/yup";
import { formikUtils } from "src/utils/form";
import { Studio, StudioSelect } from "src/components/Studios/StudioSelect";
import { Scene, SceneSelect } from "src/components/Scenes/SceneSelect";
import { useTagsEdit } from "src/hooks/tagsEdit";
import { ScraperMenu } from "src/components/Shared/ScraperMenu";
import {
  CustomFieldsInput,
  formatCustomFieldInput,
} from "src/components/Shared/CustomFields";
import cloneDeep from "lodash-es/cloneDeep";
import {
  applyStoryMetadata,
  getStoryMetadata,
  isStoryGallery,
  withoutStoryCustomFields,
} from "./storyMetadata";

interface IProps {
  gallery: Partial<GQL.GalleryDataFragment>;
  isVisible: boolean;
  onSubmit: (input: GQL.GalleryCreateInput, andNew?: boolean) => Promise<void>;
  onDelete: () => void;
}

export const GalleryEditPanel: React.FC<IProps> = ({
  gallery,
  isVisible,
  onSubmit,
  onDelete,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [scenes, setScenes] = useState<Scene[]>([]);

  const [performers, setPerformers] = useState<Performer[]>([]);
  const [studio, setStudio] = useState<Studio | null>(null);

  const isNew = gallery.id === undefined;

  const scrapers = useListGalleryScrapers();

  const [scrapedGallery, setScrapedGallery] =
    useState<GQL.ScrapedGallery | null>();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const titleRequired =
    isNew || (gallery?.files?.length === 0 && !gallery?.folder);
  const storyMetadata = getStoryMetadata(gallery);
  const shouldRenderStoryMetadata = isStoryGallery(gallery);

  const schema = yup.object({
    title: titleRequired ? yup.string().required() : yup.string().ensure(),
    code: yup.string().ensure(),
    urls: yupUniqueStringList(intl),
    date: yupDateString(intl),
    photographer: yup.string().ensure(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
    scene_ids: yup.array(yup.string().required()).defined(),
    details: yup.string().ensure(),
    custom_fields: yup.object().required().defined(),
    story_author: yup.string().ensure(),
    story_language: yup.string().ensure(),
    story_source_website: yup.string().ensure(),
    story_tag_line: yup.string().ensure(),
    story_audio_url: yup.string().ensure(),
    story_back_cover_url: yup.string().ensure(),
  });

  const initialValues = {
    title: gallery?.title ?? "",
    code: gallery?.code ?? "",
    urls: gallery?.urls ?? [],
    date: gallery?.date ?? "",
    photographer: gallery?.photographer ?? "",
    studio_id: gallery?.studio?.id ?? null,
    performer_ids: (gallery?.performers ?? []).map((p) => p.id),
    tag_ids: (gallery?.tags ?? []).map((t) => t.id),
    scene_ids: (gallery?.scenes ?? []).map((s) => s.id),
    details: gallery?.details ?? "",
    custom_fields: withoutStoryCustomFields(
      cloneDeep(gallery?.custom_fields ?? {})
    ),
    story_author: storyMetadata.author,
    story_language: storyMetadata.language,
    story_source_website: storyMetadata.sourceWebsite,
    story_tag_line: storyMetadata.tagLine,
    story_audio_url: storyMetadata.audioUrl,
    story_back_cover_url: storyMetadata.backCoverUrl,
  };

  type InputValues = yup.InferType<typeof schema>;

  const [customFieldsError, setCustomFieldsError] = useState<string>();

  function submit(values: InputValues) {
    const customFields = applyStoryMetadata(values.custom_fields, {
      author: values.story_author,
      language: values.story_language,
      sourceWebsite: values.story_source_website,
      tagLine: values.story_tag_line,
      audioUrl: values.story_audio_url,
      backCoverUrl: values.story_back_cover_url,
    });
    const inputValues = schema.cast(values);

    const input = {
      title: inputValues.title,
      code: inputValues.code,
      urls: inputValues.urls,
      date: inputValues.date,
      photographer: inputValues.photographer,
      studio_id: inputValues.studio_id,
      performer_ids: inputValues.performer_ids,
      tag_ids: inputValues.tag_ids,
      scene_ids: inputValues.scene_ids,
      details: inputValues.details,
      custom_fields: formatCustomFieldInput(isNew, customFields),
    };
    onSave(input);
  }

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: submit,
  });

  const { tags, updateTagsStateFromScraper, tagsControl } = useTagsEdit(
    gallery.tags,
    (ids) => formik.setFieldValue("tag_ids", ids)
  );

  function onSetScenes(items: Scene[]) {
    setScenes(items);
    formik.setFieldValue(
      "scene_ids",
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
    setPerformers(gallery.performers ?? []);
  }, [gallery.performers]);

  useEffect(() => {
    setStudio(gallery.studio ?? null);
  }, [gallery.studio]);

  useEffect(() => {
    setScenes(gallery.scenes ?? []);
  }, [gallery.scenes]);

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

  const fragmentScrapers = useMemo(() => {
    return (scrapers?.data?.listScrapers ?? []).filter((s) =>
      s.gallery?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );
  }, [scrapers]);

  const cover = useMemo(() => {
    if (gallery?.paths?.cover) {
      return (
        <div className="gallery-cover">
          <img
            src={gallery.paths.cover}
            alt={intl.formatMessage({ id: "cover_image" })}
          />
        </div>
      );
    }

    return <div></div>;
  }, [gallery?.paths?.cover, intl]);

  async function onSave(input: GQL.GalleryCreateInput, andNew?: boolean) {
    setIsLoading(true);
    try {
      await onSubmit(input, andNew);
      formik.resetForm();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onSaveAndNewClick() {
    const customFields = applyStoryMetadata(formik.values.custom_fields, {
      author: formik.values.story_author,
      language: formik.values.story_language,
      sourceWebsite: formik.values.story_source_website,
      tagLine: formik.values.story_tag_line,
      audioUrl: formik.values.story_audio_url,
      backCoverUrl: formik.values.story_back_cover_url,
    });
    const inputValues = schema.cast(formik.values);

    const input = {
      title: inputValues.title,
      code: inputValues.code,
      urls: inputValues.urls,
      date: inputValues.date,
      photographer: inputValues.photographer,
      studio_id: inputValues.studio_id,
      performer_ids: inputValues.performer_ids,
      tag_ids: inputValues.tag_ids,
      scene_ids: inputValues.scene_ids,
      details: inputValues.details,
      custom_fields: formatCustomFieldInput(isNew, customFields),
    };
    onSave(input, true);
  }

  async function onScrapeClicked(s: GQL.ScraperSourceInput) {
    if (!gallery || !gallery.id) return;

    setIsLoading(true);
    try {
      const result = await queryScrapeGallery(s.scraper_id!, gallery.id);
      if (!result.data || !result.data.scrapeSingleGallery?.length) {
        Toast.success("No galleries found");
        return;
      }
      setScrapedGallery(result.data.scrapeSingleGallery[0]);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onReloadScrapers() {
    setIsLoading(true);
    try {
      await mutateReloadScrapers();
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function onScrapeDialogClosed(data?: GQL.ScrapedGalleryDataFragment) {
    if (data) {
      updateGalleryFromScrapedGallery(data);
    }
    setScrapedGallery(undefined);
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedGallery) {
      return;
    }

    const currentGallery = {
      id: gallery.id!,
      ...formik.values,
    };

    return (
      <GalleryScrapeDialog
        gallery={currentGallery}
        galleryStudio={studio}
        galleryTags={tags}
        galleryPerformers={performers}
        scraped={scrapedGallery}
        onClose={(data) => {
          onScrapeDialogClosed(data);
        }}
      />
    );
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (scrapers?.data?.listScrapers ?? []).some((s) =>
      (s?.gallery?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateGalleryFromScrapedGallery(
    galleryData: GQL.ScrapedGalleryDataFragment
  ) {
    if (galleryData.title) {
      formik.setFieldValue("title", galleryData.title);
    }

    if (galleryData.code) {
      formik.setFieldValue("code", galleryData.code);
    }

    if (galleryData.details) {
      formik.setFieldValue("details", galleryData.details);
    }

    if (galleryData.photographer) {
      formik.setFieldValue("photographer", galleryData.photographer);
    }

    if (galleryData.date) {
      formik.setFieldValue("date", galleryData.date);
    }

    if (galleryData.urls) {
      formik.setFieldValue("urls", galleryData.urls);
    }

    if (galleryData.studio?.stored_id) {
      onSetStudio({
        id: galleryData.studio.stored_id,
        name: galleryData.studio.name ?? "",
        aliases: [],
      });
    }

    if (galleryData.performers?.length) {
      const idPerfs = galleryData.performers.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idPerfs.length > 0) {
        onSetPerformers(
          idPerfs.map((p) => {
            return {
              id: p.stored_id!,
              name: p.name ?? "",
              alias_list: [],
            };
          })
        );
      }
    }

    updateTagsStateFromScraper(galleryData.tags ?? undefined);
  }

  async function onScrapeGalleryURL(url: string) {
    if (!url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeGalleryURL(url);
      if (!result || !result.data || !result.data.scrapeGalleryURL) {
        return;
      }
      setScrapedGallery(result.data.scrapeGalleryURL);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
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

  function renderScenesField() {
    const title = intl.formatMessage({ id: "scenes" });
    const control = (
      <SceneSelect
        values={scenes}
        onSelect={(items) => onSetScenes(items)}
        isMulti
      />
    );

    return renderField("scene_ids", title, control);
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
      } catch (e) {
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

  function renderStoryMetadataFields() {
    if (!shouldRenderStoryMetadata) {
      return;
    }

    return (
      <>
        <h6 className="story-metadata-heading">
          <FormattedMessage id="story_metadata" />
        </h6>
        {renderInputField("story_author", "text", "author", fullWidthProps)}
        {renderInputField("story_language", "text", "language", fullWidthProps)}
        {renderInputField(
          "story_source_website",
          "text",
          "source_website",
          fullWidthProps
        )}
        {renderInputField(
          "story_tag_line",
          "text",
          "tag_line",
          fullWidthProps
        )}
        {renderInputField(
          "story_audio_url",
          "text",
          "audio_link",
          fullWidthProps
        )}
        {renderInputField(
          "story_back_cover_url",
          "text",
          "back_cover",
          fullWidthProps
        )}
      </>
    );
  }

  return (
    <div id="gallery-edit-details">
      <Prompt
        when={formik.dirty}
        message={handleUnsavedChanges(intl, "galleries", gallery?.id)}
      />

      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <Row className="form-container edit-buttons-container px-3 pt-3">
          <div className="edit-buttons mb-3 pl-0">
            {isNew ? (
              <SplitButton
                id="gallery-save-split-button"
                className="edit-button"
                variant="primary"
                disabled={
                  !isEqual(formik.errors, {}) || customFieldsError !== undefined
                }
                title={intl.formatMessage({ id: "actions.save" })}
                onClick={() => formik.submitForm()}
              >
                <Dropdown.Item onClick={() => onSaveAndNewClick()}>
                  <FormattedMessage id="actions.save_and_new" />
                </Dropdown.Item>
              </SplitButton>
            ) : (
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
            )}
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => onDelete()}
            >
              <FormattedMessage id="actions.delete" />
            </Button>
          </div>
          <div className="ml-auto text-right d-flex">
            {!isNew && (
              <ScraperMenu
                toggle={intl.formatMessage({ id: "actions.scrape_with" })}
                scrapers={fragmentScrapers}
                onScraperClicked={onScrapeClicked}
                onReloadScrapers={onReloadScrapers}
              />
            )}
          </div>
        </Row>
        <Row className="form-container px-3">
          <Col lg={7} xl={12}>
            {renderInputField("title")}
            {renderInputField("code", "text", "scene_code")}

            {renderURLListField(
              "urls",
              onScrapeGalleryURL,
              urlScrapable,
              "urls",
              urlProps
            )}

            {renderDateField("date")}
            {renderInputField("photographer")}

            {renderScenesField()}
            {renderStudioField()}
            {renderPerformersField()}
            {renderTagsField()}
          </Col>
          <Col lg={5} xl={12}>
            {renderDetailsField()}
            {renderStoryMetadataFields()}
            <Form.Group controlId="cover_image">
              <Form.Label>
                <FormattedMessage id="cover_image" />
              </Form.Label>
              {cover}
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
