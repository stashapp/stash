import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import _ from "lodash";
import { useBulkImageUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect, Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";

interface IListOperationProps {
  selected: GQL.SlimImageDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditImagesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [
    performerMode,
    setPerformerMode,
  ] = React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();

  const [updateImages] = useBulkImageUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function makeBulkUpdateIds(
    ids: string[],
    mode: GQL.BulkUpdateIdMode
  ): GQL.BulkUpdateIds {
    return {
      mode,
      ids,
    };
  }

  function getImageInput(): GQL.BulkImageUpdateInput {
    // need to determine what we are actually setting on each image
    const aggregateRating = getRating(props.selected);
    const aggregateStudioId = getStudioId(props.selected);
    const aggregatePerformerIds = getPerformerIds(props.selected);
    const aggregateTagIds = getTagIds(props.selected);

    const imageInput: GQL.BulkImageUpdateInput = {
      ids: props.selected.map((image) => {
        return image.id;
      }),
    };

    // if rating is undefined
    if (rating === undefined) {
      // and all images have the same rating, then we are unsetting the rating.
      if (aggregateRating) {
        // null rating to unset it
        imageInput.rating = null;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      imageInput.rating = rating;
    }

    // if studioId is undefined
    if (studioId === undefined) {
      // and all images have the same studioId,
      // then unset the studioId, otherwise ignoring studioId
      if (aggregateStudioId) {
        // null studio_id to unset it
        imageInput.studio_id = null;
      }
    } else {
      // if studioId is set, then we are setting it
      imageInput.studio_id = studioId;
    }

    // if performerIds are empty
    if (
      performerMode === GQL.BulkUpdateIdMode.Set &&
      (!performerIds || performerIds.length === 0)
    ) {
      // and all images have the same ids,
      if (aggregatePerformerIds.length > 0) {
        // then unset the performerIds, otherwise ignore
        imageInput.performer_ids = makeBulkUpdateIds(
          performerIds || [],
          performerMode
        );
      }
    } else {
      // if performerIds non-empty, then we are setting them
      imageInput.performer_ids = makeBulkUpdateIds(
        performerIds || [],
        performerMode
      );
    }

    // if tagIds non-empty, then we are setting them
    if (
      tagMode === GQL.BulkUpdateIdMode.Set &&
      (!tagIds || tagIds.length === 0)
    ) {
      // and all images have the same ids,
      if (aggregateTagIds.length > 0) {
        // then unset the tagIds, otherwise ignore
        imageInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
      }
    } else {
      // if tagIds non-empty, then we are setting them
      imageInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
    }

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
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "countables.images" }, { count: 100 })
              .toLocaleLowerCase(),
          }
        ),
      });
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  function getRating(state: GQL.SlimImageDataFragment[]) {
    let ret: number | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      if (first) {
        ret = image.rating ?? undefined;
        first = false;
      } else if (ret !== image.rating) {
        ret = undefined;
      }
    });

    return ret;
  }

  function getStudioId(state: GQL.SlimImageDataFragment[]) {
    let ret: string | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      if (first) {
        ret = image?.studio?.id;
        first = false;
      } else {
        const studio = image?.studio?.id;
        if (ret !== studio) {
          ret = undefined;
        }
      }
    });

    return ret;
  }

  function getPerformerIds(state: GQL.SlimImageDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      if (first) {
        ret = image.performers ? image.performers.map((p) => p.id).sort() : [];
        first = false;
      } else {
        const perfIds = image.performers
          ? image.performers.map((p) => p.id).sort()
          : [];

        if (!_.isEqual(ret, perfIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function getTagIds(state: GQL.SlimImageDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      if (first) {
        ret = image.tags ? image.tags.map((t) => t.id).sort() : [];
        first = false;
      } else {
        const tIds = image.tags ? image.tags.map((t) => t.id).sort() : [];

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
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((image: GQL.SlimImageDataFragment) => {
      const imageRating = image.rating;
      const imageStudioID = image?.studio?.id;
      const imagePerformerIDs = (image.performers ?? [])
        .map((p) => p.id)
        .sort();
      const imageTagIDs = (image.tags ?? []).map((p) => p.id).sort();

      if (first) {
        updateRating = imageRating ?? undefined;
        updateStudioID = imageStudioID;
        updatePerformerIds = imagePerformerIDs;
        updateTagIds = imageTagIDs;
        updateOrganized = image.organized;
        first = false;
      } else {
        if (imageRating !== updateRating) {
          updateRating = undefined;
        }
        if (imageStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!_.isEqual(imagePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!_.isEqual(imageTagIDs, updateTagIds)) {
          updateTagIds = [];
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
    setOrganized(updateOrganized);
  }, [props.selected, performerMode, tagMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

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
      <MultiSet
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
        existingIds={existingIds ?? []}
        ids={ids ?? []}
        mode={mode}
      />
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
      <Modal
        show
        icon="pencil-alt"
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage(
              { id: "countables.images" },
              { count: 1 }
            ),
            pluralEntity: intl.formatMessage(
              { id: "countables.images" },
              { count: props?.selected?.length ?? 1 }
            ),
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
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
                disabled={isUpdating}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage(
                { id: "countables.studios" },
                { count: 1 }
              ),
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
              <FormattedMessage
                id="countables.performers"
                values={{ count: 100 }}
              />
            </Form.Label>
            {renderMultiSelect("performers", performerIds)}
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="countables.tags" values={{ count: 100 }} />
            </Form.Label>
            {renderMultiSelect("tags", tagIds)}
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
      </Modal>
    );
  }

  return render();
};
