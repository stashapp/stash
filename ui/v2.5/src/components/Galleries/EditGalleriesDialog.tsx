import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import _ from "lodash";
import { useBulkGalleryUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect, Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";

interface IListOperationProps {
  selected: GQL.GalleryDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditGalleriesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const Toast = useToast();
  const [rating, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [performerMode, setPerformerMode] = React.useState<
    GQL.BulkUpdateIdMode
  >(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();

  const [updateGalleries] = useBulkGalleryUpdate(getGalleryInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  function makeBulkUpdateIds(
    ids: string[],
    mode: GQL.BulkUpdateIdMode
  ): GQL.BulkUpdateIds {
    return {
      mode,
      ids,
    };
  }

  function getGalleryInput(): GQL.BulkGalleryUpdateInput {
    // need to determine what we are actually setting on each gallery
    const aggregateRating = getRating(props.selected);
    const aggregateStudioId = getStudioId(props.selected);
    const aggregatePerformerIds = getPerformerIds(props.selected);
    const aggregateTagIds = getTagIds(props.selected);

    const galleryInput: GQL.BulkGalleryUpdateInput = {
      ids: props.selected.map((gallery) => {
        return gallery.id;
      }),
    };

    // if rating is undefined
    if (rating === undefined) {
      // and all galleries have the same rating, then we are unsetting the rating.
      if (aggregateRating) {
        // an undefined rating is ignored in the server, so set it to 0 instead
        galleryInput.rating = 0;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      galleryInput.rating = rating;
    }

    // if studioId is undefined
    if (studioId === undefined) {
      // and all galleries have the same studioId,
      // then unset the studioId, otherwise ignoring studioId
      if (aggregateStudioId) {
        // an undefined studio_id is ignored in the server, so set it to empty string instead
        galleryInput.studio_id = "";
      }
    } else {
      // if studioId is set, then we are setting it
      galleryInput.studio_id = studioId;
    }

    // if performerIds are empty
    if (
      performerMode === GQL.BulkUpdateIdMode.Set &&
      (!performerIds || performerIds.length === 0)
    ) {
      // and all galleries have the same ids,
      if (aggregatePerformerIds.length > 0) {
        // then unset the performerIds, otherwise ignore
        galleryInput.performer_ids = makeBulkUpdateIds(
          performerIds || [],
          performerMode
        );
      }
    } else {
      // if performerIds non-empty, then we are setting them
      galleryInput.performer_ids = makeBulkUpdateIds(
        performerIds || [],
        performerMode
      );
    }

    // if tagIds non-empty, then we are setting them
    if (
      tagMode === GQL.BulkUpdateIdMode.Set &&
      (!tagIds || tagIds.length === 0)
    ) {
      // and all galleries have the same ids,
      if (aggregateTagIds.length > 0) {
        // then unset the tagIds, otherwise ignore
        galleryInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
      }
    } else {
      // if tagIds non-empty, then we are setting them
      galleryInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
    }

    return galleryInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateGalleries();
      Toast.success({ content: "Updated galleries" });
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  function getRating(state: GQL.GalleryDataFragment[]) {
    let ret: number | undefined;
    let first = true;

    state.forEach((gallery: GQL.GalleryDataFragment) => {
      if (first) {
        ret = gallery.rating ?? undefined;
        first = false;
      } else if (ret !== gallery.rating) {
        ret = undefined;
      }
    });

    return ret;
  }

  function getStudioId(state: GQL.GalleryDataFragment[]) {
    let ret: string | undefined;
    let first = true;

    state.forEach((gallery: GQL.GalleryDataFragment) => {
      if (first) {
        ret = gallery?.studio?.id;
        first = false;
      } else {
        const studio = gallery?.studio?.id;
        if (ret !== studio) {
          ret = undefined;
        }
      }
    });

    return ret;
  }

  function getPerformerIds(state: GQL.GalleryDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((gallery: GQL.GalleryDataFragment) => {
      if (first) {
        ret = gallery.performers
          ? gallery.performers.map((p) => p.id).sort()
          : [];
        first = false;
      } else {
        const perfIds = gallery.performers
          ? gallery.performers.map((p) => p.id).sort()
          : [];

        if (!_.isEqual(ret, perfIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function getTagIds(state: GQL.GalleryDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((gallery: GQL.GalleryDataFragment) => {
      if (first) {
        ret = gallery.tags ? gallery.tags.map((t) => t.id).sort() : [];
        first = false;
      } else {
        const tIds = gallery.tags ? gallery.tags.map((t) => t.id).sort() : [];

        if (!_.isEqual(ret, tIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let first = true;

    state.forEach((gallery: GQL.GalleryDataFragment) => {
      const galleryRating = gallery.rating;
      const GalleriestudioID = gallery?.studio?.id;
      const galleryPerformerIDs = (gallery.performers ?? [])
        .map((p) => p.id)
        .sort();
      const galleryTagIDs = (gallery.tags ?? []).map((p) => p.id).sort();

      if (first) {
        updateRating = galleryRating ?? undefined;
        updateStudioID = GalleriestudioID;
        updatePerformerIds = galleryPerformerIDs;
        updateTagIds = galleryTagIDs;
        first = false;
      } else {
        if (galleryRating !== updateRating) {
          updateRating = undefined;
        }
        if (GalleriestudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!_.isEqual(galleryPerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!_.isEqual(galleryTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioID);
    if (performerMode === GQL.BulkUpdateIdMode.Set) {
      setPerformerIds(updatePerformerIds);
    }

    if (tagMode === GQL.BulkUpdateIdMode.Set) {
      setTagIds(updateTagIds);
    }
  }, [props.selected, performerMode, tagMode]);

  function renderMultiSelect(
    type: "performers" | "tags",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    switch (type) {
      case "performers":
        mode = performerMode;
        break;
      case "tags":
        mode = tagMode;
        break;
    }

    return (
      <MultiSet
        type={type}
        disabled={isUpdating}
        onUpdate={(items) => {
          const itemIDs = items.map((i) => i.id);
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
        ids={ids ?? []}
        mode={mode}
      />
    );
  }

  function render() {
    return (
      <Modal
        show
        icon="pencil-alt"
        header="Edit Galleries"
        accept={{ onClick: onSave, text: "Apply" }}
        cancel={{
          onClick: () => props.onClose(false),
          text: "Cancel",
          variant: "secondary",
        }}
        isRunning={isUpdating}
      >
        <Form>
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: "Rating",
            })}
            <Col xs={9}>
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
                disabled={isUpdating}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: "Studio",
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
            <Form.Label>Performers</Form.Label>
            {renderMultiSelect("performers", performerIds)}
          </Form.Group>

          <Form.Group controlId="performers">
            <Form.Label>Tags</Form.Label>
            {renderMultiSelect("tags", tagIds)}
          </Form.Group>
        </Form>
      </Modal>
    );
  }

  return render();
};
