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
import { MultiSet } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateGalleryIds,
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioId,
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
  const [studioId, setStudioId] = useState<string>();
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
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateGalleryIds = getAggregateGalleryIds(props.selected);

    const imageInput: GQL.BulkImageUpdateInput = {
      ids: props.selected.map((image) => {
        return image.id;
      }),
    };

    imageInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    imageInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);

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
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      const imageRating = image.rating100;
      const imageStudioID = image?.studio?.id;
      const imagePerformerIDs = (image.performers ?? [])
        .map((p) => p.id)
        .sort();
      const imageTagIDs = (image.tags ?? []).map((p) => p.id).sort();
      const imageGalleryIDs = (image.galleries ?? []).map((p) => p.id).sort();

      if (first) {
        updateRating = imageRating ?? undefined;
        updateStudioID = imageStudioID;
        updatePerformerIds = imagePerformerIDs;
        updateTagIds = imageTagIDs;
        updateGalleryIds = imageGalleryIDs;
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
                menuPortalTarget={document.body}
              />
            </Col>
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
