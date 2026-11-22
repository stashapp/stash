import React, { useEffect, useMemo, useState } from "react";
import { Alert, Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregateGroupIds,
  getAggregatePerformerIds,
  getAggregateStateObject,
  getAggregateTagIds,
  getAggregateStudioId,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { IndeterminateCheckbox } from "../Shared/IndeterminateCheckbox";
import { BulkUpdateFormGroup, BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { BulkUpdateDateInput } from "../Shared/DateInput";
import { getDateError } from "src/utils/yup";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (applied: boolean) => void;
}

const sceneFields = [
  "code",
  "rating100",
  "details",
  "organized",
  "director",
  "date",
  "production_date",
];

// bulk input fields that map to lockable scene fields
const lockableInputFields: [
  keyof GQL.BulkSceneUpdateInput,
  GQL.SceneLockableField,
][] = [
  ["code", GQL.SceneLockableField.Code],
  ["details", GQL.SceneLockableField.Details],
  ["director", GQL.SceneLockableField.Director],
  ["date", GQL.SceneLockableField.Date],
  ["production_date", GQL.SceneLockableField.ProductionDate],
  ["studio_id", GQL.SceneLockableField.Studio],
];

export const EditScenesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const [updateInput, setUpdateInput] = useState<GQL.BulkSceneUpdateInput>({
    ids: props.selected.map((scene) => {
      return scene.id;
    }),
  });

  const [dateError, setDateError] = useState<string | undefined>();
  const [productionDateError, setProductionDateError] = useState<
    string | undefined
  >();

  const [performerIds, setPerformerIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [groupIds, setGroupIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });

  const unsetDisabled = props.selected.length < 2;

  const [updateScenes] = useBulkSceneUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const aggregateState = useMemo(() => {
    const updateState: Partial<GQL.BulkSceneUpdateInput> = {};
    const state = props.selected;
    updateState.studio_id = getAggregateStudioId(props.selected);
    const updateTagIds = getAggregateTagIds(props.selected);
    const updatePerformerIds = getAggregatePerformerIds(props.selected);
    const updateGroupIds = getAggregateGroupIds(props.selected);
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      getAggregateStateObject(updateState, scene, sceneFields, first);
      first = false;
    });

    return {
      state: updateState,
      tagIds: updateTagIds,
      performerIds: updatePerformerIds,
      groupIds: updateGroupIds,
    };
  }, [props.selected]);

  // update initial state from aggregate
  useEffect(() => {
    setUpdateInput((current) => ({ ...current, ...aggregateState.state }));
  }, [aggregateState]);

  useEffect(() => {
    setDateError(getDateError(updateInput.date ?? "", intl));
  }, [updateInput.date, intl]);

  useEffect(() => {
    setProductionDateError(
      getDateError(updateInput.production_date ?? "", intl)
    );
  }, [updateInput.production_date, intl]);

  // number of selected scenes where a changed field is locked, and so will be skipped
  const lockedSceneCount = useMemo(() => {
    const changed = lockableInputFields
      .filter(
        ([key]) =>
          updateInput[key] !== undefined &&
          updateInput[key] !== aggregateState.state[key]
      )
      .map(([, field]) => field);

    return props.selected.filter((s) =>
      s.locked_fields.some((f) => changed.includes(f))
    ).length;
  }, [updateInput, aggregateState, props.selected]);

  function setUpdateField(input: Partial<GQL.BulkSceneUpdateInput>) {
    setUpdateInput((current) => ({ ...current, ...input }));
  }

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    const sceneInput: GQL.BulkSceneUpdateInput = {
      ...updateInput,
      tag_ids: tagIds,
      performer_ids: performerIds,
      group_ids: groupIds,
    };

    // we don't have unset functionality for the rating star control
    // so need to determine if we are setting a rating or not
    sceneInput.rating100 = getAggregateInputValue(
      updateInput.rating100,
      aggregateState.state.rating100
    );

    return sceneInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateScenes({ variables: { input: getSceneInput() } });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "scenes" }).toLocaleLowerCase() }
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
            singularEntity: intl.formatMessage({ id: "scene" }),
            pluralEntity: intl.formatMessage({ id: "scenes" }),
          }
        )}
        accept={{
          onClick: onSave,
          text: intl.formatMessage({ id: "actions.apply" }),
        }}
        disabled={isUpdating || !!dateError || !!productionDateError}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        isRunning={isUpdating}
      >
        {lockedSceneCount > 0 && (
          <Alert variant="warning">
            {intl.formatMessage(
              { id: "dialogs.bulk_edit_locked_fields_warning" },
              { count: lockedSceneCount }
            )}
          </Alert>
        )}
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

          <BulkUpdateFormGroup name="scene_code">
            <BulkUpdateTextInput
              value={updateInput.code}
              valueChanged={(newValue) => setUpdateField({ code: newValue })}
              unsetDisabled={unsetDisabled}
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

          <BulkUpdateFormGroup name="production_date">
            <BulkUpdateDateInput
              value={updateInput.production_date}
              valueChanged={(newValue) =>
                setUpdateField({ production_date: newValue })
              }
              unsetDisabled={unsetDisabled}
              error={productionDateError}
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

          <BulkUpdateFormGroup name="performers" inline={false}>
            <MultiSet
              type={"performers"}
              disabled={isUpdating}
              onUpdate={(itemIDs) => {
                setPerformerIds((c) => ({ ...c, ids: itemIDs }));
              }}
              onSetMode={(newMode) => {
                setPerformerIds((c) => ({ ...c, mode: newMode }));
              }}
              ids={performerIds.ids ?? []}
              existingIds={aggregateState.performerIds}
              mode={performerIds.mode}
              menuPortalTarget={document.body}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="groups" inline={false}>
            <MultiSet
              type={"groups"}
              disabled={isUpdating}
              onUpdate={(itemIDs) => {
                setGroupIds((c) => ({ ...c, ids: itemIDs }));
              }}
              onSetMode={(newMode) => {
                setGroupIds((c) => ({ ...c, mode: newMode }));
              }}
              ids={groupIds.ids ?? []}
              existingIds={aggregateState.groupIds}
              mode={groupIds.mode}
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

          <BulkUpdateFormGroup name="details" inline={false}>
            <BulkUpdateTextInput
              value={updateInput.details}
              valueChanged={(newValue) => setUpdateField({ details: newValue })}
              unsetDisabled={unsetDisabled}
              as="textarea"
            />
          </BulkUpdateFormGroup>

          <Form.Group controlId="organized">
            <IndeterminateCheckbox
              label={intl.formatMessage({ id: "organized" })}
              setChecked={(checked) => setUpdateField({ organized: checked })}
              checked={updateInput.organized ?? undefined}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
