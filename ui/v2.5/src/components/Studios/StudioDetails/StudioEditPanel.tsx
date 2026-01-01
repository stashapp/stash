import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ImageInput } from "src/components/Shared/ImageInput";
import cx from "classnames";
import { Button, Dropdown, Form } from "react-bootstrap";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import ImageUtils from "src/utils/image";
import { addUpdateStashID, getStashIDs } from "src/utils/stashIds";
import { stashboxDisplayName } from "src/utils/stashbox";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import isEqual from "lodash-es/isEqual";
import { useToast } from "src/hooks/Toast";
import { useConfigurationContext } from "src/hooks/Config";
import { handleUnsavedChanges } from "src/utils/navigation";
import { formikUtils } from "src/utils/form";
import { yupFormikValidate, yupUniqueAliases } from "src/utils/yup";
import { Studio, StudioSelect } from "../StudioSelect";
import { useTagsEdit } from "src/hooks/tagsEdit";
import { Icon } from "src/components/Shared/Icon";
import StashBoxIDSearchModal from "src/components/Shared/StashBoxIDSearchModal";
import StudioStashBoxModal, { IStashBox } from "./StudioStashBoxModal";
import { StudioScrapeDialog } from "./StudioScrapeDialog";

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
  const { configuration: stashConfig } = useConfigurationContext();

  const isNew = studio.id === undefined;

  // Editing/scraper state
  const [scraper, setScraper] = useState<IStashBox>();
  const [isScraperModalOpen, setIsScraperModalOpen] = useState(false);
  const [scrapedStudio, setScrapedStudio] = useState<GQL.ScrapedStudio>();
  const [isStashIDSearchOpen, setIsStashIDSearchOpen] = useState(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [parentStudio, setParentStudio] = useState<Studio | null>(null);

  const schema = yup.object({
    name: yup.string().required(),
    urls: yup.array(yup.string().required()).defined(),
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
    urls: studio.urls ?? [],
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

  const { tags, updateTagsStateFromScraper, tagsControl } = useTagsEdit(
    studio.tags,
    (ids) => formik.setFieldValue("tag_ids", ids)
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

  function onStashIDSelected(item?: GQL.StashIdInput) {
    if (!item) return;
    formik.setFieldValue(
      "stash_ids",
      addUpdateStashID(formik.values.stash_ids, item)
    );
  }

  function updateStashIDs(remoteSiteID: string | null | undefined) {
    if (remoteSiteID && scraper?.endpoint) {
      const newIDs =
        formik.values.stash_ids?.filter(
          (s) => s.endpoint !== scraper.endpoint
        ) ?? [];
      newIDs.push({
        endpoint: scraper.endpoint,
        stash_id: remoteSiteID,
        updated_at: new Date().toISOString(),
      });
      formik.setFieldValue("stash_ids", newIDs);
    }
  }

  function updateStudioEditStateFromScraper(
    state: Partial<GQL.ScrapedStudioDataFragment>
  ) {
    if (state.name) {
      formik.setFieldValue("name", state.name);
    }
    if (state.urls) {
      formik.setFieldValue("urls", state.urls);
    }
    if (state.details) {
      formik.setFieldValue("details", state.details);
    }
    if (state.aliases) {
      formik.setFieldValue(
        "aliases",
        state.aliases.split(",").map((a) => a.trim())
      );
    }
    updateTagsStateFromScraper(state.tags ?? undefined);

    // image is a base64 string
    // overwrite if not new since it came from a dialog
    // overwrite if image is unset
    if ((!isNew || !formik.values.image) && state.image) {
      formik.setFieldValue("image", state.image);
    }

    updateStashIDs(state.remote_site_id);
  }

  function onScrapeStashBox(studioResult: GQL.ScrapedStudio) {
    setIsScraperModalOpen(false);

    const result: GQL.ScrapedStudioDataFragment = {
      ...studioResult,
      __typename: "ScrapedStudio",
    };

    // if this is a new studio, just dump the data
    if (isNew) {
      updateStudioEditStateFromScraper(result);
      setScraper(undefined);
    } else {
      setScrapedStudio(result);
    }
  }

  function onScraperSelected(s: IStashBox) {
    setScraper(s);
    setIsScraperModalOpen(true);
  }

  function renderScraperMenu() {
    if (!studio) {
      return;
    }
    const stashBoxes = stashConfig?.general.stashBoxes ?? [];

    if (stashBoxes.length === 0) {
      return;
    }

    const popover = (
      <Dropdown.Menu id="studio-scraper-popover">
        {stashBoxes.map((s, index) => (
          <Dropdown.Item
            as={Button}
            key={s.endpoint}
            className="minimal"
            onClick={() => onScraperSelected({ ...s, index })}
          >
            {stashboxDisplayName(s.name, index)}
          </Dropdown.Item>
        ))}
      </Dropdown.Menu>
    );

    return (
      <Dropdown className="d-inline-block">
        <Dropdown.Toggle variant="secondary" className="mr-2">
          <FormattedMessage id="actions.scrape_with" />
        </Dropdown.Toggle>
        {popover}
      </Dropdown>
    );
  }

  function renderButtons(classNames: string) {
    return (
      <div className={cx("details-edit", "col-xl-9", classNames)}>
        {!isNew && (
          <Button className="mr-2" variant="primary" onClick={onCancel}>
            <FormattedMessage id="actions.cancel" />
          </Button>
        )}
        {renderScraperMenu()}
        <ImageInput
          isEditing
          onImageChange={onImageChange}
          onImageURL={onImageLoad}
          acceptSVG
        />
        <div>
          <Button
            className="mr-2"
            variant="danger"
            onClick={() => onImageLoad(null)}
          >
            <FormattedMessage id="actions.clear_image" />
          </Button>
        </div>
        <Button
          variant="success"
          disabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
          onClick={() => formik.submitForm()}
        >
          <FormattedMessage id="actions.save" />
        </Button>
      </div>
    );
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedStudio || !scraper) {
      return;
    }

    const currentStudio = {
      ...formik.values,
      image: formik.values.image ?? studio.image_path,
    };

    return (
      <StudioScrapeDialog
        studio={currentStudio}
        studioTags={tags}
        scraped={scrapedStudio}
        scraper={scraper}
        onClose={(s) => {
          onScrapeDialogClosed(s);
        }}
      />
    );
  }

  function onScrapeDialogClosed(s?: GQL.ScrapedStudioDataFragment) {
    if (s) {
      updateStudioEditStateFromScraper(s);
    }
    setScrapedStudio(undefined);
    setScraper(undefined);
  }

  function renderScrapeModal() {
    if (!isScraperModalOpen || !scraper) {
      return;
    }

    return (
      <StudioStashBoxModal
        instance={scraper}
        onHide={() => setScraper(undefined)}
        onSelectStudio={onScrapeStashBox}
        name={formik.values.name || ""}
      />
    );
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
      {renderScrapeModal()}
      {maybeRenderScrapeDialog()}
      {isStashIDSearchOpen && (
        <StashBoxIDSearchModal
          entityType="studio"
          stashBoxes={stashConfig?.general.stashBoxes ?? []}
          excludedStashBoxEndpoints={formik.values.stash_ids.map(
            (s) => s.endpoint
          )}
          onSelectItem={(item) => {
            onStashIDSelected(item);
            setIsStashIDSearchOpen(false);
          }}
          initialQuery={studio.name ?? ""}
        />
      )}

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
        {renderStringListField("urls")}
        {renderInputField("details", "textarea")}
        {renderParentStudioField()}
        {renderTagsField()}
        {renderStashIDsField(
          "stash_ids",
          "studios",
          "stash_ids",
          undefined,
          <Button
            variant="success"
            className="mr-2 py-0"
            onClick={() => setIsStashIDSearchOpen(true)}
            disabled={!stashConfig?.general.stashBoxes?.length}
            title={intl.formatMessage({ id: "actions.add_stash_id" })}
          >
            <Icon icon={faPlus} />
          </Button>
        )}
        <hr />
        {renderInputField("ignore_auto_tag", "checkbox")}
      </Form>

      {renderButtons("mt-3")}
    </>
  );
};
