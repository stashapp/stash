import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form, Col, Row } from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
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
import {
  queryScrapeImage,
  queryScrapeImageURL,
  useListImageScrapers,
  mutateReloadScrapers,
} from "../../../core/StashService";
import { ImageScrapeDialog } from "./ImageScrapeDialog";
import { Studio, StudioSelect } from "src/components/Studios/StudioSelect";
import { galleryTitle } from "src/core/galleries";
import {
  Gallery,
  GallerySelect,
  excludeFileBasedGalleries,
} from "src/components/Galleries/GallerySelect";
import { useTagsEdit } from "src/hooks/tagsEdit";
import { ScraperMenu } from "src/components/Shared/ScraperMenu";

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

  const [galleries, setGalleries] = useState<Gallery[]>([]);
  const [performers, setPerformers] = useState<Performer[]>([]);
  const [studio, setStudio] = useState<Studio | null>(null);

  const isNew = image.id === undefined;

  useEffect(() => {
    setGalleries(
      image.galleries?.map((g) => ({
        id: g.id,
        title: galleryTitle(g),
        files: g.files,
        folder: g.folder,
      })) ?? []
    );
  }, [image.galleries]);

  const scrapers = useListImageScrapers();
  const [scrapedImage, setScrapedImage] = useState<GQL.ScrapedImage | null>();

  const schema = yup.object({
    title: yup.string().ensure(),
    code: yup.string().ensure(),
    urls: yupUniqueStringList(intl),
    date: yupDateString(intl),
    details: yup.string().ensure(),
    photographer: yup.string().ensure(),
    gallery_ids: yup.array(yup.string().required()).defined(),
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
    gallery_ids: (image.galleries ?? []).map((g) => g.id),
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

  const { tags, updateTagsStateFromScraper, tagsControl } = useTagsEdit(
    image.tags,
    (ids) => formik.setFieldValue("tag_ids", ids)
  );

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
    setPerformers(image.performers ?? []);
  }, [image.performers]);

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

  const fragmentScrapers = useMemo(() => {
    return (scrapers?.data?.listScrapers ?? []).filter((s) =>
      s.image?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );
  }, [scrapers]);

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

  async function onScrapeClicked(s: GQL.ScraperSourceInput) {
    if (!image || !image.id) return;

    setIsLoading(true);
    try {
      const result = await queryScrapeImage(s.scraper_id!, image.id);
      if (!result.data || !result.data.scrapeSingleImage?.length) {
        Toast.success("No images found");
        return;
      }
      setScrapedImage(result.data.scrapeSingleImage[0]);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (scrapers?.data?.listScrapers ?? []).some((s) =>
      (s?.image?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateImageFromScrapedGallery(
    imageData: GQL.ScrapedImageDataFragment
  ) {
    if (imageData.title) {
      formik.setFieldValue("title", imageData.title);
    }

    if (imageData.code) {
      formik.setFieldValue("code", imageData.code);
    }

    if (imageData.details) {
      formik.setFieldValue("details", imageData.details);
    }

    if (imageData.photographer) {
      formik.setFieldValue("photographer", imageData.photographer);
    }

    if (imageData.date) {
      formik.setFieldValue("date", imageData.date);
    }

    if (imageData.urls) {
      formik.setFieldValue("urls", imageData.urls);
    }

    if (imageData.studio?.stored_id) {
      onSetStudio({
        id: imageData.studio.stored_id,
        name: imageData.studio.name ?? "",
        aliases: [],
      });
    }

    if (imageData.performers?.length) {
      const idPerfs = imageData.performers.filter((p) => {
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

    updateTagsStateFromScraper(imageData.tags ?? undefined);
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

  async function onScrapeDialogClosed(data?: GQL.ScrapedImageDataFragment) {
    if (data) {
      updateImageFromScrapedGallery(data);
    }
    setScrapedImage(undefined);
  }

  async function onScrapeImageURL(url: string) {
    if (!url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeImageURL(url);
      if (!result || !result.data || !result.data.scrapeImageURL) {
        return;
      }
      setScrapedImage(result.data.scrapeImageURL);
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

  function renderGalleriesField() {
    const title = intl.formatMessage({ id: "galleries" });
    const control = (
      <GallerySelect
        values={galleries}
        onSelect={(items) => onSetGalleries(items)}
        isMulti
        extraCriteria={excludeFileBasedGalleries}
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

  function maybeRenderScrapeDialog() {
    if (!scrapedImage) {
      return;
    }

    const currentImage = {
      id: image.id!,
      ...formik.values,
    };

    return (
      <ImageScrapeDialog
        image={currentImage}
        imageStudio={studio}
        imageTags={tags}
        imagePerformers={performers}
        scraped={scrapedImage}
        onClose={(data) => {
          onScrapeDialogClosed(data);
        }}
      />
    );
  }

  return (
    <div id="image-edit-details">
      <Prompt
        when={formik.dirty}
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />

      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <Row className="form-container edit-buttons-container px-3 pt-3">
          <div className="edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={
                (!isNew && !formik.dirty) || !isEqual(formik.errors, {})
              }
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
              onScrapeImageURL,
              urlScrapable,
              "urls",
              urlProps
            )}

            {renderDateField("date")}
            {renderInputField("photographer")}

            {renderGalleriesField()}
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
