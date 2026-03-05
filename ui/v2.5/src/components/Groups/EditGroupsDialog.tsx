import React, { useEffect, useMemo, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkGroupUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregateStateObject,
  getAggregateTagIds,
  getAggregateStudioId,
  getAggregateIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { BulkUpdateDateInput } from "../Shared/DateInput";
import { IRelatedGroupEntry } from "./GroupDetails/RelatedGroupTable";
import { ContainingGroupsMultiSet } from "./ContainingGroupsMultiSet";

interface IListOperationProps {
  selected: GQL.ListGroupDataFragment[];
  onClose: (applied: boolean) => void;
}

export function getAggregateContainingGroups(
  state: Pick<GQL.ListGroupDataFragment, "containing_groups">[]
) {
  const sortedLists: IRelatedGroupEntry[][] = state.map((o) =>
    o.containing_groups
      .map((oo) => ({
        group: oo.group,
        description: oo.description,
      }))
      .sort((a, b) => a.group.id.localeCompare(b.group.id))
  );

  return getAggregateIds(sortedLists);
}

function getAggregateContainingGroupInput(
  mode: GQL.BulkUpdateIdMode,
  input: IRelatedGroupEntry[] | undefined,
  aggregateValues: IRelatedGroupEntry[]
): GQL.BulkUpdateGroupDescriptionsInput | undefined {
  if (mode === GQL.BulkUpdateIdMode.Set && (!input || input.length === 0)) {
    // and all scenes have the same ids,
    if (aggregateValues.length > 0) {
      // then unset, otherwise ignore
      return { mode, groups: [] };
    }
  } else {
    // if input non-empty, then we are setting them
    return {
      mode,
      groups:
        input?.map((e) => {
          return { group_id: e.group.id, description: e.description };
        }) || [],
    };
  }

  return undefined;
}

const groupFields = ["rating100", "synopsis", "director", "date"];

export const EditGroupsDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const [updateInput, setUpdateInput] = useState<GQL.BulkGroupUpdateInput>({
    ids: props.selected.map((group) => {
      return group.id;
    }),
  });

  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [containingGroupsMode, setGroupMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [containingGroups, setGroups] = useState<IRelatedGroupEntry[]>();

  const [updateGroups] = useBulkGroupUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const aggregateState = useMemo(() => {
    const updateState: Partial<GQL.BulkGroupUpdateInput> = {};
    const state = props.selected;
    updateState.studio_id = getAggregateStudioId(props.selected);
    const updateTagIds = getAggregateTagIds(props.selected);
    const aggregateGroups = getAggregateContainingGroups(props.selected);
    let first = true;

    state.forEach((group: GQL.ListGroupDataFragment) => {
      getAggregateStateObject(updateState, group, groupFields, first);
      first = false;
    });

    return {
      state: updateState,
      tagIds: updateTagIds,
      containingGroups: aggregateGroups,
    };
  }, [props.selected]);

  // update initial state from aggregate
  useEffect(() => {
    setUpdateInput((current) => ({ ...current, ...aggregateState.state }));
  }, [aggregateState]);

  function setUpdateField(input: Partial<GQL.BulkGroupUpdateInput>) {
    setUpdateInput((current) => ({ ...current, ...input }));
  }

  function getGroupInput(): GQL.BulkGroupUpdateInput {
    const groupInput: GQL.BulkGroupUpdateInput = {
      ...updateInput,
      tag_ids: tagIds,
    };

    // we don't have unset functionality for the rating star control
    // so need to determine if we are setting a rating or not
    groupInput.rating100 = getAggregateInputValue(
      updateInput.rating100,
      aggregateState.state.rating100
    );

    groupInput.containing_groups = getAggregateContainingGroupInput(
      containingGroupsMode,
      containingGroups,
      aggregateState.containingGroups
    );

    return groupInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateGroups({ variables: { input: getGroupInput() } });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "groups" }).toLocaleLowerCase() }
        )
      );
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  function renderTextField(
    name: string,
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void,
    area: boolean = false
  ) {
    const control = (
      <BulkUpdateTextInput
        value={value === null ? "" : value ?? undefined}
        valueChanged={(newValue) => setter(newValue)}
        unsetDisabled={props.selected.length < 2}
        as={area ? "textarea" : undefined}
      />
    );

    if (area) {
      return (
        <Form.Group
          controlId={name}
          data-field={name}
          as={area ? undefined : Row}
        >
          <Form.Label>
            <FormattedMessage id={name} />
          </Form.Label>
          {control}
        </Form.Group>
      );
    }

    return (
      <Form.Group controlId={name} data-field={name} as={Row}>
        {FormUtils.renderLabel({
          title: intl.formatMessage({ id: name }),
        })}
        <Col xs={9}>{control}</Col>
      </Form.Group>
    );
  }

  function render() {
    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_count_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "group" }),
            pluralEntity: intl.formatMessage({ id: "groups" }),
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
                value={updateInput.rating100}
                onSetRating={(value) =>
                  setUpdateField({ rating100: value ?? undefined })
                }
                disabled={isUpdating}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="date" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "date" }),
            })}
            <Col xs={9}>
              <BulkUpdateDateInput
                value={updateInput.date ?? undefined}
                valueChanged={(newValue) => setUpdateField({ date: newValue })}
                unsetDisabled={props.selected.length < 2}
              />
            </Col>
          </Form.Group>

          {renderTextField("director", updateInput.director, (v) =>
            setUpdateField({ director: v })
          )}
          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "studio" }),
            })}
            <Col xs={9}>
              <StudioSelect
                onSelect={(items) =>
                  setUpdateField({
                    studio_id: items.length > 0 ? items[0]?.id : undefined,
                  })
                }
                ids={updateInput.studio_id ? [updateInput.studio_id] : []}
                isDisabled={isUpdating}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
          <Form.Group controlId="containing-groups">
            <Form.Label>
              <FormattedMessage id="containing_groups" />
            </Form.Label>
            <ContainingGroupsMultiSet
              disabled={isUpdating}
              onUpdate={(v) => setGroups(v)}
              onSetMode={(newMode) => setGroupMode(newMode)}
              existingValue={aggregateState.containingGroups ?? []}
              value={containingGroups ?? []}
              mode={containingGroupsMode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            <MultiSet
              type={"tags"}
              disabled={isUpdating}
              onUpdate={(itemIDs) => {
                setTagIds((c) => ({ ...c, ids: itemIDs }));
              }}
              onSetMode={(newMode) => {
                setTagIds((c) => ({ ...c, mode: newMode }));
              }}
              ids={tagIds.ids ?? []}
              existingIds={aggregateState.tagIds}
              mode={tagIds.mode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          {renderTextField(
            "synopsis",
            updateInput.synopsis,
            (v) => setUpdateField({ synopsis: v }),
            true
          )}
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
