import React, { useEffect, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useBulkGroupUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregateStateObject,
  getAggregateTagIds,
  getAggregateStudioId,
  getAggregateIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { BulkUpdateFormGroup, BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { BulkUpdateDateInput } from "../Shared/DateInput";
import { IRelatedGroupEntry } from "./GroupDetails/RelatedGroupTable";
import { ContainingGroupsMultiSet } from "./ContainingGroupsMultiSet";
import { getDateError } from "src/utils/yup";

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

  const unsetDisabled = props.selected.length < 2;

  const [updateGroups] = useBulkGroupUpdate();

  const [dateError, setDateError] = useState<string | undefined>();

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

  useEffect(() => {
    setDateError(getDateError(updateInput.date ?? "", intl));
  }, [updateInput.date, intl]);

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
        disabled={isUpdating || !!dateError}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        isRunning={isUpdating}
      >
        <Form>
          <BulkUpdateFormGroup name="rating">
            <RatingSystem
              value={updateInput.rating100}
              onSetRating={(value) =>
                setUpdateField({ rating100: value ?? undefined })
              }
              disabled={isUpdating}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="date">
            <BulkUpdateDateInput
              value={updateInput.date}
              valueChanged={(newValue) => setUpdateField({ date: newValue })}
              unsetDisabled={unsetDisabled}
              error={dateError}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="director">
            <BulkUpdateTextInput
              value={updateInput.director}
              valueChanged={(newValue) =>
                setUpdateField({ director: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="studio">
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
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup
            name="containing-groups"
            messageId="containing_groups"
            inline={false}
          >
            <ContainingGroupsMultiSet
              disabled={isUpdating}
              onUpdate={(v) => setGroups(v)}
              onSetMode={(newMode) => setGroupMode(newMode)}
              existingValue={aggregateState.containingGroups ?? []}
              value={containingGroups ?? []}
              mode={containingGroupsMode}
              menuPortalTarget={document.body}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="tags" inline={false}>
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
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="synopsis" inline={false}>
            <BulkUpdateTextInput
              value={updateInput.synopsis}
              valueChanged={(newValue) =>
                setUpdateField({ synopsis: newValue })
              }
              unsetDisabled={unsetDisabled}
              as="textarea"
            />
          </BulkUpdateFormGroup>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
