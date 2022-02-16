import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import _ from "lodash";
import { useBulkPerformerUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";
import {
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregateRating,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";

interface IListOperationProps {
  selected: GQL.SlimPerformerDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditPerformersDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating, setRating] = useState<number>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [favorite, setFavorite] = useState<boolean | undefined>();

  const [updatePerformers] = useBulkPerformerUpdate(getPerformerInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function getPerformerInput(): GQL.BulkPerformerUpdateInput {
    // need to determine what we are actually setting on each performer
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateRating = getAggregateRating(props.selected);

    const performerInput: GQL.BulkPerformerUpdateInput = {
      ids: props.selected.map((performer) => {
        return performer.id;
      }),
    };

    performerInput.rating = getAggregateInputValue(rating, aggregateRating);

    performerInput.tag_ids = getAggregateInputIDs(
      tagMode,
      tagIds,
      aggregateTagIds
    );

    if (favorite !== undefined) {
      performerInput.favorite = favorite;
    }

    return performerInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updatePerformers();
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "performers" })
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

  useEffect(() => {
    const state = props.selected;
    let updateTagIds: string[] = [];
    let updateFavorite: boolean | undefined;
    let updateRating: number | undefined;
    let first = true;

    state.forEach((performer: GQL.SlimPerformerDataFragment) => {
      const performerTagIDs = (performer.tags ?? []).map((p) => p.id).sort();
      const performerRating = performer.rating;

      if (first) {
        updateTagIds = performerTagIDs;
        first = false;
        updateFavorite = performer.favorite;
        updateRating = performerRating ?? undefined;
      } else {
        if (!_.isEqual(performerTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (performer.favorite !== updateFavorite) {
          updateFavorite = undefined;
        }
        if (performerRating !== updateRating) {
          updateRating = undefined;
        }
      }
    });

    setExistingTagIds(updateTagIds);
    setFavorite(updateFavorite);
    setRating(updateRating);
  }, [props.selected, tagMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = favorite === undefined;
    }
  }, [favorite, checkboxRef]);

  function cycleFavorite() {
    if (favorite) {
      setFavorite(undefined);
    } else if (favorite === undefined) {
      setFavorite(false);
    } else {
      setFavorite(true);
    }
  }

  function render() {
    return (
      <Modal
        show
        icon="pencil-alt"
        header="Edit Performers"
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
        <Form>
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
            />
          </Form.Group>

          <Form.Group controlId="favorite">
            <Form.Check
              type="checkbox"
              label="Favorite"
              checked={favorite}
              ref={checkboxRef}
              onChange={() => cycleFavorite()}
            />
          </Form.Group>
        </Form>
      </Modal>
    );
  }

  return render();
};
