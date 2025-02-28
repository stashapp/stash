import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import {
  queryScrapeGroupURL,
  useListGroupScrapers,
} from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { useToast } from "src/hooks/Toast";
import { Modal as BSModal, Form, Button } from "react-bootstrap";
import TextUtils from "src/utils/text";
import ImageUtils from "src/utils/image";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { GroupScrapeDialog } from "./GroupScrapeDialog";
import isEqual from "lodash-es/isEqual";
import { handleUnsavedChanges } from "src/utils/navigation";
import { formikUtils } from "src/utils/form";
import {
  yupDateString,
  yupFormikValidate,
  yupUniqueStringList,
} from "src/utils/yup";
import { Studio, StudioSelect } from "src/components/Studios/StudioSelect";
import { useTagsEdit } from "src/hooks/tagsEdit";
import { Group } from "src/components/Groups/GroupSelect";
import { RelatedGroupTable, IRelatedGroupEntry } from "./RelatedGroupTable";

interface IGroupEditPanel {
  group: Partial<GQL.GroupDataFragment>;
  onSubmit: (group: GQL.GroupCreateInput) => Promise<void>;
  onCancel: () => void;
  onDelete: () => void;
  setFrontImage: (image?: string | null) => void;
  setBackImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const GroupEditPanel: React.FC<IGroupEditPanel> = ({
  group,
  onSubmit,
  onCancel,
  onDelete,
  setFrontImage,
  setBackImage,
  setEncodingImage,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const isNew = group.id === undefined;

  const [isLoading, setIsLoading] = useState(false);
  const [isImageAlertOpen, setIsImageAlertOpen] = useState<boolean>(false);

  const [imageClipboard, setImageClipboard] = useState<string>();

  const Scrapers = useListGroupScrapers();
  const [scrapedGroup, setScrapedGroup] = useState<GQL.ScrapedGroup>();

  const [studio, setStudio] = useState<Studio | null>(null);
  const [containingGroups, setContainingGroups] = useState<Group[]>([]);

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().ensure(),
    duration: yup.number().integer().min(0).nullable().defined(),
    date: yupDateString(intl),
    studio_id: yup.string().required().nullable(),
    tag_ids: yup.array(yup.string().required()).defined(),
    containing_groups: yup
      .array(
        yup.object({
          group_id: yup.string().required(),
          description: yup.string().nullable().ensure(),
        })
      )
      .defined(),
    director: yup.string().ensure(),
    urls: yupUniqueStringList(intl),
    synopsis: yup.string().ensure(),
    front_image: yup.string().nullable().optional(),
    back_image: yup.string().nullable().optional(),
  });

  const initialValues = {
    name: group?.name ?? "",
    aliases: group?.aliases ?? "",
    duration: group?.duration ?? null,
    date: group?.date ?? "",
    studio_id: group?.studio?.id ?? null,
    tag_ids: (group?.tags ?? []).map((t) => t.id),
    containing_groups: (group?.containing_groups ?? []).map((m) => {
      return { group_id: m.group.id, description: m.description ?? "" };
    }),
    director: group?.director ?? "",
    urls: group?.urls ?? [],
    synopsis: group?.synopsis ?? "",
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validate: yupFormikValidate(schema),
    onSubmit: (values) => onSave(schema.cast(values)),
  });

  const { tags, updateTagsStateFromScraper, tagsControl } = useTagsEdit(
    group.tags,
    (ids) => formik.setFieldValue("tag_ids", ids)
  );

  const containingGroupEntries = useMemo(() => {
    return formik.values.containing_groups
      .map((m) => {
        return {
          group: containingGroups.find((mm) => mm.id === m.group_id),
          description: m.description,
        };
      })
      .filter((m) => m.group !== undefined) as IRelatedGroupEntry[];
  }, [formik.values.containing_groups, containingGroups]);

  function onSetStudio(item: Studio | null) {
    setStudio(item);
    formik.setFieldValue("studio_id", item ? item.id : null);
  }

  useEffect(() => {
    setStudio(group.studio ?? null);
  }, [group.studio]);

  useEffect(() => {
    setContainingGroups(group.containing_groups?.map((m) => m.group) ?? []);
  }, [group.containing_groups]);

  // set up hotkeys
  useEffect(() => {
    // Mousetrap.bind("u", (e) => {
    //   setStudioFocus()
    //   e.preventDefault();
    // });
    Mousetrap.bind("s s", () => {
      if (formik.dirty) {
        formik.submitForm();
      }
    });

    return () => {
      // Mousetrap.unbind("u");
      Mousetrap.unbind("s s");
    };
  });

  function updateGroupEditStateFromScraper(
    state: Partial<GQL.ScrapedGroupDataFragment>
  ) {
    if (state.name) {
      formik.setFieldValue("name", state.name);
    }

    if (state.aliases) {
      formik.setFieldValue("aliases", state.aliases);
    }

    if (state.duration) {
      const seconds = TextUtils.timestampToSeconds(state.duration);
      if (seconds) {
        formik.setFieldValue("duration", seconds);
      }
    }

    if (state.date) {
      formik.setFieldValue("date", state.date);
    }

    if (state.studio && state.studio.stored_id) {
      onSetStudio({
        id: state.studio.stored_id,
        name: state.studio.name ?? "",
        aliases: [],
      });
    }

    if (state.director) {
      formik.setFieldValue("director", state.director);
    }
    if (state.synopsis) {
      formik.setFieldValue("synopsis", state.synopsis);
    }
    if (state.urls) {
      formik.setFieldValue("urls", state.urls);
    }
    updateTagsStateFromScraper(state.tags ?? undefined);

    if (state.front_image) {
      // image is a base64 string
      formik.setFieldValue("front_image", state.front_image);
    }
    if (state.back_image) {
      // image is a base64 string
      formik.setFieldValue("back_image", state.back_image);
    }
  }

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

  async function onScrapeGroupURL(url: string) {
    if (!url) return;
    setIsLoading(true);

    try {
      const result = await queryScrapeGroupURL(url);
      if (!result.data || !result.data.scrapeGroupURL) {
        return;
      }

      // if this is a new group, just dump the data
      if (isNew) {
        updateGroupEditStateFromScraper(result.data.scrapeGroupURL);
      } else {
        setScrapedGroup(result.data.scrapeGroupURL);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function urlScrapable(scrapedUrl: string) {
    return (
      !!scrapedUrl &&
      (Scrapers?.data?.listScrapers ?? []).some((s) =>
        (s?.group?.urls ?? []).some((u) => scrapedUrl.includes(u))
      )
    );
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedGroup) {
      return;
    }

    const currentGroup = {
      id: group.id!,
      ...formik.values,
    };

    // Get image paths for scrape gui
    currentGroup.front_image = group?.front_image_path;
    currentGroup.back_image = group?.back_image_path;

    return (
      <GroupScrapeDialog
        group={currentGroup}
        groupStudio={studio}
        groupTags={tags}
        scraped={scrapedGroup}
        onClose={(m) => {
          onScrapeDialogClosed(m);
        }}
      />
    );
  }

  function onScrapeDialogClosed(p?: GQL.ScrapedGroupDataFragment) {
    if (p) {
      updateGroupEditStateFromScraper(p);
    }
    setScrapedGroup(undefined);
  }

  const encodingImage = ImageUtils.usePasteImage(showImageAlert);

  useEffect(() => {
    setFrontImage(formik.values.front_image);
  }, [formik.values.front_image, setFrontImage]);

  useEffect(() => {
    setBackImage(formik.values.back_image);
  }, [formik.values.back_image, setBackImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

  function onFrontImageLoad(imageData: string | null) {
    formik.setFieldValue("front_image", imageData);
  }

  function onFrontImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onFrontImageLoad);
  }

  function onBackImageLoad(imageData: string | null) {
    formik.setFieldValue("back_image", imageData);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
  }

  function showImageAlert(imageData: string) {
    setImageClipboard(imageData);
    setIsImageAlertOpen(true);
  }

  function setImageFromClipboard(isFrontImage: boolean) {
    if (isFrontImage) {
      formik.setFieldValue("front_image", imageClipboard);
    } else {
      formik.setFieldValue("back_image", imageClipboard);
    }

    setImageClipboard(undefined);
    setIsImageAlertOpen(false);
  }

  function renderImageAlert() {
    return (
      <BSModal
        show={isImageAlertOpen}
        onHide={() => setIsImageAlertOpen(false)}
      >
        <BSModal.Body>
          <p>Select image to set</p>
        </BSModal.Body>
        <BSModal.Footer>
          <div>
            <Button
              className="mr-2"
              variant="secondary"
              onClick={() => setIsImageAlertOpen(false)}
            >
              <FormattedMessage id="actions.cancel" />
            </Button>

            <Button
              className="mr-2"
              onClick={() => setImageFromClipboard(false)}
            >
              Back Image
            </Button>
            <Button
              className="mr-2"
              onClick={() => setImageFromClipboard(true)}
            >
              Front Image
            </Button>
          </div>
        </BSModal.Footer>
      </BSModal>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  const {
    renderField,
    renderInputField,
    renderDateField,
    renderDurationField,
    renderURLListField,
  } = formikUtils(intl, formik);

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

  function renderTagsField() {
    const title = intl.formatMessage({ id: "tags" });
    return renderField("tag_ids", title, tagsControl());
  }

  function onSetContainingGroupEntries(input: IRelatedGroupEntry[]) {
    setContainingGroups(input.map((m) => m.group));

    const newGroups = input.map((m) => ({
      group_id: m.group.id,
      description: m.description,
    }));

    formik.setFieldValue("containing_groups", newGroups);
  }

  function renderContainingGroupsField() {
    const title = intl.formatMessage({ id: "containing_groups" });
    const control = (
      <RelatedGroupTable
        value={containingGroupEntries}
        onUpdate={onSetContainingGroupEntries}
        excludeIDs={group.id ? [group.id] : undefined}
      />
    );

    return renderField("containing_groups", title, control);
  }

  // TODO: CSS class
  return (
    <div>
      {isNew && (
        <h2>
          {intl.formatMessage(
            { id: "actions.add_entity" },
            { entityType: intl.formatMessage({ id: "group" }) }
          )}
        </h2>
      )}

      <Prompt
        when={formik.dirty}
        message={(location, action) => {
          // Check if it's a redirect after group creation
          if (action === "PUSH" && location.pathname.startsWith("/groups/"))
            return true;

          return handleUnsavedChanges(intl, "groups", group.id)(location);
        }}
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="group-edit">
        {renderInputField("name")}
        {renderInputField("aliases")}
        {renderDurationField("duration")}
        {renderDateField("date")}
        {renderContainingGroupsField()}
        {renderStudioField()}
        {renderInputField("director")}
        {renderURLListField("urls", onScrapeGroupURL, urlScrapable)}
        {renderInputField("synopsis", "textarea")}
        {renderTagsField()}
      </Form>

      <DetailsEditNavbar
        objectName={group?.name ?? intl.formatMessage({ id: "group" })}
        isNew={isNew}
        classNames="col-xl-9 mt-3"
        isEditing
        onToggleEdit={onCancel}
        onSave={formik.handleSubmit}
        saveDisabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
        onImageChange={onFrontImageChange}
        onImageChangeURL={onFrontImageLoad}
        onClearImage={() => onFrontImageLoad(null)}
        onBackImageChange={onBackImageChange}
        onBackImageChangeURL={onBackImageLoad}
        onClearBackImage={() => onBackImageLoad(null)}
        onDelete={onDelete}
      />

      {maybeRenderScrapeDialog()}
      {renderImageAlert()}
    </div>
  );
};
