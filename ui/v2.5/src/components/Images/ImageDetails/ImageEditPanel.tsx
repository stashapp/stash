import React, { useEffect, useState } from "react";
import {
  Button,
  Form,
  Col,
  Row,
  Dropdown,
  DropdownButton,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { TagSelect, StudioSelect } from "src/components/Shared/Select";
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
import { Icon } from "../../Shared/Icon";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import {
  mutateReloadScrapers,
  queryScrapeImage,
  queryScrapeImageURL,
  useListImageScrapers,
} from "../../../core/StashService";
import { ImageScrapeDialog } from "./ImageScrapeDialog";
import { ScrapedImageDataFragment } from "src/core/generated-graphql";

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

  const Scrapers = useListImageScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);
  const [scrapedImage, setScrapedImage] = useState<GQL.ScrapedImage | null>();

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

  useRatingKeybinds(
    true,
    configuration?.ui?.ratingSystemOptions?.type,
    setRating
  );

  useEffect(() => {
    setPerformers(image.performers ?? []);
  }, [image.performers]);

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

  useEffect(() => {
    const newQueryableScrapers = (Scrapers?.data?.listScrapers ?? []).filter(
      (s) => s.image?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

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

  async function onScrapeClicked(scraper: GQL.Scraper) {
    if (!image || !image.id) return;

    setIsLoading(true);
    try {
      const result = await queryScrapeImage(scraper.id, image.id);
      if (!result.data || !result.data.scrapeSingleImage?.length) {
        Toast.success("No galleries found");
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
    return (Scrapers?.data?.listScrapers ?? []).some((s) =>
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
      formik.setFieldValue("studio_id", imageData.studio.stored_id);
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

    if (imageData?.tags?.length) {
      const idTags = imageData.tags.filter((t) => {
        return t.stored_id !== undefined && t.stored_id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map((t) => t.stored_id);
        formik.setFieldValue("tag_ids", newIds as string[]);
      }
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

  async function onScrapeDialogClosed(data?: ScrapedImageDataFragment) {
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
        onSelect={(items) =>
          formik.setFieldValue(
            "studio_id",
            items.length > 0 ? items[0]?.id : null
          )
        }
        ids={formik.values.studio_id ? [formik.values.studio_id] : []}
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
        onSelect={(items) =>
          formik.setFieldValue(
            "tag_ids",
            items.map((item) => item.id)
          )
        }
        ids={formik.values.tag_ids}
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
        imagePerformer={performers}
        scraped={scrapedImage}
        onClose={(data) => {
          onScrapeDialogClosed(data);
        }}
      />
    );
  }

  function renderScraperMenu() {
    /*
    if (isNew) {
      return;
    }
     */

    return (
      <DropdownButton
        className="d-inline-block"
        id="gallery-scrape"
        title={intl.formatMessage({ id: "actions.scrape_with" })}
      >
        {queryableScrapers.map((s) => (
          <Dropdown.Item key={s.name} onClick={() => onScrapeClicked(s)}>
            {s.name}
          </Dropdown.Item>
        ))}
        <Dropdown.Item onClick={() => onReloadScrapers()}>
          <span className="fa-icon">
            <Icon icon={faSyncAlt} />
          </span>
          <span>
            <FormattedMessage id="actions.reload_scrapers" />
          </span>
        </Dropdown.Item>
      </DropdownButton>
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
          <div className="ml-auto text-right d-flex">{renderScraperMenu()}</div>
        </Row>
        <Row className="form-container px-3">
          <Col lg={7} xl={12}>
            {renderInputField("title")}
            {renderInputField("code", "text", "scene_code")}

            {renderURLListField("urls", onScrapeImageURL, urlScrapable)}

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
