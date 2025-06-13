import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkImageUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { MultiSet } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateGalleryIds,
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioIds,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";

interface IListOperationProps {
  selected: GQL.SlimImageDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditImagesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number>();
  const [studioIds, setStudioIds] = useState<string[]>();
  const [studioMode, setStudioMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [existingStudioIds, setExistingStudioIds] = useState<string[]>();
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

  const [updateImages] = useBulkImageUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function getImageInput(): GQL.BulkImageUpdateInput {
    // need to determine what we are actually setting on each image
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioIds = getAggregateStudioIds(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateGalleryIds = getAggregateGalleryIds(props.selected);

    const imageInput: GQL.BulkImageUpdateInput = {
      ids: props.selected.map((image) => {
        return image.id;
      }),
    };

    imageInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    imageInput.studio_ids = getAggregateInputIDs(
      studioMode,
      studioIds,
      aggregateStudioIds
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
    let updateStudioIDs: string[] = [];
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateGalleryIds: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      const imageRating = image.rating100;
      const imageStudioIDs = (image.studios ?? []).map((s) => s.id);
      const imagePerformerIDs = (image.performers ?? [])
        .map((p) => p.id)
        .sort();
      const imageTagIDs = (image.tags ?? []).map((p) => p.id).sort();
      const imageGalleryIDs = (image.galleries ?? []).map((p) => p.id).sort();

      if (first) {
        updateRating = imageRating ?? undefined;
        updateStudioIDs = imageStudioIDs;
        updatePerformerIds = imagePerformerIDs;
        updateTagIds = imageTagIDs;
        updateGalleryIds = imageGalleryIDs;
        updateOrganized = image.organized;
        first = false;
      } else {
        if (imageRating !== updateRating) {
          updateRating = undefined;
        }
        if (!isEqual(imageStudioIDs, updateStudioIDs)) {
          updateStudioIDs = [];
        }
        if (!isEqual(imagePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
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
    setStudioIds(updateStudioIDs);
    setExistingStudioIds(updateStudioIDs);
    setExistingPerformerIds(updatePerformerIds);
    setExistingTagIds(updateTagIds);
    setExistingGalleryIds(updateGalleryIds);
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
          <Form.Group controlId="studios">
            <Form.Label>
              <FormattedMessage id="studios" />
            </Form.Label>
            <MultiSet
              type="studios"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setStudioIds(itemIDs)}
              onSetMode={(newMode) => setStudioMode(newMode)}
              existingIds={existingStudioIds ?? []}
              ids={studioIds ?? []}
              mode={studioMode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          <Form.Group controlId="performers">
            <Form.Label>
              <FormattedMessage id="performers" />
            </Form.Label>
            <MultiSet
              type="performers"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setPerformerIds(itemIDs)}
              onSetMode={(newMode) => setPerformerMode(newMode)}
              existingIds={existingPerformerIds ?? []}
              ids={performerIds ?? []}
              mode={performerMode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            <MultiSet
              type="tags"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setTagIds(itemIDs)}
              onSetMode={(newMode) => setTagMode(newMode)}
              existingIds={existingTagIds ?? []}
              ids={tagIds ?? []}
              mode={tagMode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          <Form.Group controlId="galleries">
            <Form.Label>
              <FormattedMessage id="galleries" />
            </Form.Label>
            <MultiSet
              type="galleries"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setGalleryIds(itemIDs)}
              onSetMode={(newMode) => setGalleryMode(newMode)}
              existingIds={existingGalleryIds ?? []}
              ids={galleryIds ?? []}
              mode={galleryMode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

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
