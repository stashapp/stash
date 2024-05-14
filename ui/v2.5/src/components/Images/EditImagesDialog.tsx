import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkImageUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "src/components/Shared/Select";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { MultiSelect, MultiString } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateGalleryIds,
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
  selected: GQL.SlimImageDataFragment[];
  onClose: (applied: boolean) => void;
  showAllFields?: boolean;
}

const imageFields = [
  "title",
  "scene_code",
  "details",
  "photographer",
  "date",
  "urls",
];

export const EditImagesDialog: React.FC<IListOperationProps> = (
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
  const selectedUrls = props.selected.map((image) => ({
    urls: image.urls.map((url) => ({ value: url }))
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

  const [galleryMode, setGalleryMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [galleryIds, setGalleryIds] = useState<string[]>();
  const [existingGalleryIds, setExistingGalleryIds] = useState<string[]>();

  const [organized, setOrganized] = useState<boolean | undefined>();

  const [updateInput, setUpdateInput] = useState<GQL.BulkImageUpdateInput>(
    {}
  );

  const [showAllFields, setShowAllFields] = useState(props.showAllFields ?? false);

  const [updateImages] = useBulkImageUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function setUpdateField(input: Partial<GQL.BulkImageUpdateInput>) {
    setUpdateInput({ ...updateInput, ...input });
  }

  function getImageInput(): GQL.BulkImageUpdateInput {
    // need to determine what we are actually setting on each image
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateGalleryIds = getAggregateGalleryIds(props.selected);
    const aggregateUrls = getAggregateUrls(selectedUrls);

    const imageInput: GQL.BulkImageUpdateInput = {
      ids: props.selected.map((image) => {
        return image.id;
      }),
      ...updateInput,
    };

    imageInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    imageInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);

    imageInput.urls = getAggregateInputStrings(
      urlsMode,
      urls,
      aggregateUrls
    );

    imageInput.performer_ids = getAggregateInputIDs(
      performerMode,
      performerIds,
      aggregatePerformerIds
    );
    imageInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);
    imageInput.gallery_ids = getAggregateInputIDs(
      galleryMode,
      galleryIds,
      aggregateGalleryIds
    );

    if (organized !== undefined) {
      imageInput.organized = organized;
    }

    return imageInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateImages({
        variables: {
          input: getImageInput(),
        },
      });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "images" }).toLocaleLowerCase() }
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
    let updateGalleryIds: string[] = [];
    let updateUrls: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      getAggregateStateObject(state, image, imageFields, first);
      const imageRating = image.rating100;
      const imageStudioID = image?.studio?.id;
      const imagePerformerIDs = (image.performers ?? [])
        .map((p) => p.id)
        .sort();
      const imageTagIDs = (image.tags ?? []).map((p) => p.id).sort();
      const imageGalleryIDs = (image.galleries ?? []).map((p) => p.id).sort();
      const imageUrls = (image.urls ?? []);

      if (first) {
        updateRating = imageRating ?? undefined;
        updateStudioID = imageStudioID;
        updatePerformerIds = imagePerformerIDs;
        updateTagIds = imageTagIDs;
        updateGalleryIds = imageGalleryIDs;
        updateUrls = imageUrls;
        updateOrganized = image.organized;
        first = false;
      } else {
        if (imageRating !== updateRating) {
          updateRating = undefined;
        }
        if (imageStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!isEqual(imagePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!isEqual(imageUrls, updateUrls)) {
          updateUrls = [];
        }
        if (!isEqual(imageTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (!isEqual(imageGalleryIDs, updateGalleryIds)) {
          updateGalleryIds = [];
        }
        if (image.organized !== updateOrganized) {
          updateOrganized = undefined;
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioID);
    setExistingPerformerIds(updatePerformerIds);
    setExistingTagIds(updateTagIds);
    setExistingGalleryIds(updateGalleryIds);
    setExistingUrls(updateUrls);
    setOrganized(updateOrganized);
  }, [props.selected]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  function cycleOrganized() {
    if (organized) {
      setOrganized(undefined);
    } else if (organized === undefined) {
      setOrganized(false);
    } else {
      setOrganized(true);
    }
  }

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

  function renderTextField(
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void,
    isDetails: Boolean = false
  ) {
    return (
      <Form.Group>
        <BulkUpdateTextInput
        as={isDetails ? 'textarea' : undefined}
        value={value === null ? "" : value ?? undefined}
        valueChanged={(newValue) => setter(newValue)}
        unsetDisabled={props.selected.length < 2}
        />
      </Form.Group>
    );
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
            singularEntity: intl.formatMessage({ id: "image" }),
            pluralEntity: intl.formatMessage({ id: "images" }),
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

          {showAllFields && 
          <Form.Group controlId="text-input" as={Row}>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "title" })})}
            <Col xs={9}>
              {renderTextField(updateInput.title, (v) => setUpdateField({ title: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "scene_code" })})}
            <Col xs={9}>
              {renderTextField(updateInput.code, (v) => setUpdateField({ code: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "urls" })})}
            <Col xs={9}>
              {renderURLMultiSelect(urls)}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "date" })})}
            <Col xs={9}>
              {renderTextField(updateInput.date, (v) => setUpdateField({ date: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "photographer" })})}
            <Col xs={9}>
              {renderTextField(updateInput.photographer, (v) => setUpdateField({ photographer: v }))}
            </Col>
          </Form.Group>
          }

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
            <MultiSelect
              type="performers"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setPerformerIds(itemIDs)}
              onSetMode={(newMode) => setPerformerMode(newMode)}
              existing={existingPerformerIds ?? []}
              ids={performerIds ?? []}
              mode={performerMode}
            />
          </Form.Group>

          <Form.Group controlId="galleries">
            <Form.Label>
              <FormattedMessage id="galleries" />
            </Form.Label>
            <MultiSelect
              type="galleries"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setGalleryIds(itemIDs)}
              onSetMode={(newMode) => setGalleryMode(newMode)}
              existing={existingGalleryIds ?? []}
              ids={galleryIds ?? []}
              mode={galleryMode}
            />
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            <MultiSelect
              type="tags"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setTagIds(itemIDs)}
              onSetMode={(newMode) => setTagMode(newMode)}
              existing={existingTagIds ?? []}
              ids={tagIds ?? []}
              mode={tagMode}
            />
          </Form.Group>

          {showAllFields && 
          <Form.Group controlId="details">
            <Form.Label>
              <FormattedMessage id="details" />
            </Form.Label>
            {renderTextField(updateInput.details, (v) => setUpdateField({ details: v }), true)}
          </Form.Group>
          }

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
