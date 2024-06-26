import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkGroupUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "../Shared/Modal";
import { StudioSelect } from "../Shared/Select";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { isEqual } from "lodash-es";
import { MultiSet } from "../Shared/MultiSet";

interface IListOperationProps {
  selected: GQL.GroupDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditGroupsDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number | undefined>();
  const [studioId, setStudioId] = useState<string | undefined>();
  const [director, setDirector] = useState<string | undefined>();

  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();

  const [updateGroups] = useBulkGroupUpdate(getGroupInput());

  const [isUpdating, setIsUpdating] = useState(false);

  function getGroupInput(): GQL.BulkGroupUpdateInput {
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);

    const groupInput: GQL.BulkGroupUpdateInput = {
      ids: props.selected.map((group) => group.id),
      director,
    };

    // if rating is undefined
    groupInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    groupInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);
    groupInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);

    return groupInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateGroups();
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "groups" }).toLocaleLowerCase(),
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
    let updateStudioId: string | undefined;
    let updateTagIds: string[] = [];
    let updateDirector: string | undefined;
    let first = true;

    state.forEach((group: GQL.GroupDataFragment) => {
      const groupTagIDs = (group.tags ?? []).map((p) => p.id).sort();

      if (first) {
        first = false;
        updateRating = group.rating100 ?? undefined;
        updateStudioId = group.studio?.id ?? undefined;
        updateTagIds = groupTagIDs;
        updateDirector = group.director ?? undefined;
      } else {
        if (group.rating100 !== updateRating) {
          updateRating = undefined;
        }
        if (group.studio?.id !== updateStudioId) {
          updateStudioId = undefined;
        }
        if (group.director !== updateDirector) {
          updateDirector = undefined;
        }
        if (!isEqual(groupTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioId);
    setExistingTagIds(updateTagIds);
    setDirector(updateDirector);
  }, [props.selected]);

  function render() {
    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "actions.edit_entity" },
          { entityType: intl.formatMessage({ id: "groups" }) }
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
              />
            </Col>
          </Form.Group>
          <Form.Group controlId="director">
            <Form.Label>
              <FormattedMessage id="director" />
            </Form.Label>
            <Form.Control
              className="input-control"
              type="text"
              value={director}
              onChange={(event) => setDirector(event.currentTarget.value)}
              placeholder={intl.formatMessage({ id: "director" })}
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
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
