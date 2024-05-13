import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkGalleryUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { MultiSelect, MultiString } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputStrings,
  getAggregateInputValue,
  getAggregateStateObject,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
  getAggregateUrls,
} from "src/utils/bulkUpdate";
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";

interface IListOperationProps {
  selected: GQL.SlimGalleryDataFragment[];
  onClose: (applied: boolean) => void;
  showAllFields?: boolean;
}

const galleryFields = [
  "scene_code",
  "details",
  "photographer",
  "date",
  "urls",
];

export const EditGalleriesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [urlsMode, setUrlsMode] =
  React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [urls, setUrls] = useState<string[]>();
  const [existingUrls, setExistingUrls] = useState<string[]>();
  const selectedUrls = props.selected.map((gallery) => ({
    urls: gallery.urls.map((url) => ({ value: url }))
  }));
  const [performerMode, setPerformerMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [sceneMode, setSceneMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [sceneIds, setSceneIds] = useState<string[]>();
  const [existingSceneIds, setExistingSceneIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();
  const [updateInput, setUpdateInput] = useState<GQL.BulkGalleryUpdateInput>(
    {}
  );

  const [showAllFields, setShowAllFields] = useState(props.showAllFields ?? false);

  const [updateGalleries] = useBulkGalleryUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function setUpdateField(input: Partial<GQL.BulkGalleryUpdateInput>) {
    setUpdateInput({ ...updateInput, ...input });
  }

  function getGalleryInput(): GQL.BulkGalleryUpdateInput {
    // need to determine what we are actually setting on each gallery
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    //const aggregateSceneIds = getAggregateSceneIds(props.selected);
    const aggregateUrls = getAggregateUrls(selectedUrls);

    const galleryInput: GQL.BulkGalleryUpdateInput = {
      ids: props.selected.map((gallery) => {
        return gallery.id;
      }),
      ...updateInput,
    };

    galleryInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    galleryInput.studio_id = getAggregateInputValue(
      studioId,
      aggregateStudioId
    );
    
    // galleryInput.scene_ids = getAggregateInputIDs(
    //   sceneMode,
    //   sceneIds,
    //   //aggregateSceneIds
    // );

    galleryInput.urls = getAggregateInputStrings(
      urlsMode,
      urls,
      aggregateUrls
    );

    galleryInput.performer_ids = getAggregateInputIDs(
      performerMode,
      performerIds,
      aggregatePerformerIds
    );
    galleryInput.tag_ids = getAggregateInputIDs(
      tagMode,
      tagIds,
      aggregateTagIds
    );

    if (organized !== undefined) {
      galleryInput.organized = organized;
    }

    return galleryInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateGalleries({
        variables: {
          input: getGalleryInput(),
        },
      });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "galleries" }).toLocaleLowerCase(),
          }
        )
      );
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateSceneIds: string[] = [];
    let updateUrls: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((gallery: GQL.SlimGalleryDataFragment) => {
      getAggregateStateObject(state, gallery, galleryFields, first);
      const galleryRating = gallery.rating100;
      const GalleriestudioID = gallery?.studio?.id;
      const galleryPerformerIDs = (gallery.performers ?? [])
        .map((p) => p.id)
        .sort();
      const galleryTagIDs = (gallery.tags ?? []).map((p) => p.id).sort();
      const gallerySceneIDs = (gallery.scenes ?? []).map((s) => s.id).sort();
      const galleryUrls = (gallery.urls ?? []);

      if (first) {
        updateRating = galleryRating ?? undefined;
        updateStudioID = GalleriestudioID;
        updatePerformerIds = galleryPerformerIDs;
        updateTagIds = galleryTagIDs;
        updateSceneIds = gallerySceneIDs;
        updateUrls = galleryUrls;
        updateOrganized = gallery.organized;
        first = false;
      } else {
        if (galleryRating !== updateRating) {
          updateRating = undefined;
        }
        if (GalleriestudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!isEqual(galleryPerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!isEqual(gallerySceneIDs, updateSceneIds)) {
          updateSceneIds = [];
        }
        if (!isEqual(galleryUrls, updateUrls)) {
          updateUrls = [];
        }
        if (!isEqual(galleryTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (gallery.organized !== updateOrganized) {
          updateOrganized = undefined;
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioID);
    setExistingPerformerIds(updatePerformerIds);
    setExistingSceneIds(updateSceneIds);
    setExistingUrls(updateUrls);
    setExistingTagIds(updateTagIds);

    setOrganized(updateOrganized);
  }, [props.selected]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  function renderURLMultiSelect(
    urls: string[] | undefined
  ) {
    return (
      <MultiString
        disabled={isUpdating}
        onUpdate={(itemIDs) => {setUrls(itemIDs)}}
        onSetMode={(newMode) => {setUrlsMode(newMode)}}
        strings={urls ?? []}
        existing={existingUrls ?? []}
        mode={urlsMode}
      />
    );
  }

  function renderMultiSelect(
    type: "performers" | "tags",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    let existingIds: string[] | undefined = [];
    switch (type) {
      case "performers":
        mode = performerMode;
        existingIds = existingPerformerIds;
        break;
      case "tags":
        mode = tagMode;
        existingIds = existingTagIds;
        break;
    }

    return (
      <MultiSelect
        type={type}
        disabled={isUpdating}
        onUpdate={(itemIDs) => {
          switch (type) {
            case "performers":
              setPerformerIds(itemIDs);
              break;
            case "tags":
              setTagIds(itemIDs);
              break;
          }
        }}
        onSetMode={(newMode) => {
          switch (type) {
            case "performers":
              setPerformerMode(newMode);
              break;
            case "tags":
              setTagMode(newMode);
              break;
          }
        }}
        existing={existingIds ?? []}
        ids={ids ?? []}
        mode={mode}
      />
    );
  }

  function renderTextField(
    name: string,
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void,
    isDetails: Boolean = false
  ) {
    return (
      <Form.Group controlId={name}>
        <Form.Label>
          <FormattedMessage id={name} />
        </Form.Label>
        <BulkUpdateTextInput
          as={isDetails ? 'textarea' : undefined}
          value={value === null ? "" : value ?? undefined}
          valueChanged={(newValue) => setter(newValue)}
          unsetDisabled={props.selected.length < 2}
        />
      </Form.Group>
    );
  }

  function cycleOrganized() {
    if (organized) {
      setOrganized(undefined);
    } else if (organized === undefined) {
      setOrganized(false);
    } else {
      setOrganized(true);
    }
  }

  function render() {
    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "gallery" }),
            pluralEntity: intl.formatMessage({ id: "galleries" }),
          }
        )}
        accept={{
          onClick: onSave,
          text: intl.formatMessage({ id: "actions.apply" }),
        }}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        leftFooterButtons={
          <Form.Group controlId="toggle-all">
            <Form.Switch
              label={intl.formatMessage({ id: "actions.all_fields" })}
              checked={showAllFields}
              onChange={() => setShowAllFields(!showAllFields)}
            />
          </Form.Group>
        }
        isRunning={isUpdating}
      >
        <Form>
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "rating" }),
            })}
            <Col xs={9}>
              <RatingSystem
                value={rating100}
                onSetRating={(value) => setRating(value ?? undefined)}
                disabled={isUpdating}
              />
            </Col>
          </Form.Group>

          {showAllFields && renderTextField("scene_code", updateInput.code, (v) =>
            setUpdateField({ code: v })
          )}
          {showAllFields && 
          <Form.Group controlId="urls">
            <Form.Label>
              <FormattedMessage id="urls" />
            </Form.Label>
            {renderURLMultiSelect(urls)}
          </Form.Group>}
          {showAllFields && renderTextField("photographer", updateInput.photographer, (v) =>
            setUpdateField({ photographer: v })
          )}
          {showAllFields && renderTextField("date", updateInput.date, (v) =>
            setUpdateField({ date: v })
          )}

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "studio" }),
            })}
            <Col xs={9}>
              <StudioSelect
                onSelect={(items) =>
                  setStudioId(items.length > 0 ? items[0]?.id : undefined)
                }
                ids={studioId ? [studioId] : []}
                isDisabled={isUpdating}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="performers">
            <Form.Label>
              <FormattedMessage id="performers" />
            </Form.Label>
            {renderMultiSelect("performers", performerIds)}
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            {renderMultiSelect("tags", tagIds)}
          </Form.Group>

          {showAllFields && renderTextField("details", updateInput.details, (v) =>
            setUpdateField({ details: v }), true
          )}

          <Form.Group controlId="organized">
            <Form.Check
              type="checkbox"
              label={intl.formatMessage({ id: "organized" })}
              checked={organized}
              ref={checkboxRef}
              onChange={() => cycleOrganized()}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
