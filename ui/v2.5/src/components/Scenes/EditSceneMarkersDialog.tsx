import React, { useEffect, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useBulkSceneMarkerUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "../Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { MultiSet } from "../Shared/MultiSet";
import {
  getAggregateState,
  getAggregateStateObject,
} from "src/utils/bulkUpdate";
import { BulkUpdateFormGroup, BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { TagSelect } from "../Shared/Select";

interface IListOperationProps {
  selected: GQL.SceneMarkerDataFragment[];
  onClose: (applied: boolean) => void;
}

const scenemarkerFields = ["title"];

export const EditSceneMarkersDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const [updateInput, setUpdateInput] =
    useState<GQL.BulkSceneMarkerUpdateInput>({
      ids: props.selected.map((scenemarker) => {
        return scenemarker.id;
      }),
    });

  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });

  const unsetDisabled = props.selected.length < 2;

  const [updateSceneMarkers] = useBulkSceneMarkerUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const aggregateState = useMemo(() => {
    const updateState: Partial<GQL.BulkSceneMarkerUpdateInput> = {};
    const state = props.selected;
    let updateTagIds: string[] = [];
    let first = true;

    state.forEach((scenemarker: GQL.SceneMarkerDataFragment) => {
      getAggregateStateObject(
        updateState,
        scenemarker,
        scenemarkerFields,
        first
      );

      // sceneMarker data fragment doesn't have primary_tag_id, so handle separately
      updateState.primary_tag_id = getAggregateState(
        updateState.primary_tag_id,
        scenemarker.primary_tag.id,
        first
      );

      const thisTagIDs = (scenemarker.tags ?? []).map((p) => p.id).sort();

      updateTagIds = getAggregateState(updateTagIds, thisTagIDs, first) ?? [];

      first = false;
    });

    return { state: updateState, tagIds: updateTagIds };
  }, [props.selected]);

  // update initial state from aggregate
  useEffect(() => {
    setUpdateInput((current) => ({ ...current, ...aggregateState.state }));
  }, [aggregateState]);

  function setUpdateField(input: Partial<GQL.BulkSceneMarkerUpdateInput>) {
    setUpdateInput((current) => ({ ...current, ...input }));
  }

  function getSceneMarkerInput(): GQL.BulkSceneMarkerUpdateInput {
    const sceneMarkerInput: GQL.BulkSceneMarkerUpdateInput = {
      ...updateInput,
      tag_ids: tagIds,
    };

    return sceneMarkerInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateSceneMarkers({
        variables: {
          input: getSceneMarkerInput(),
        },
      });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "markers" }).toLocaleLowerCase(),
          }
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
        dialogClassName="edit-scenemarkers-dialog"
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_count_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "marker" }),
            pluralEntity: intl.formatMessage({ id: "markers" }),
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
          <BulkUpdateFormGroup name="title">
            <BulkUpdateTextInput
              value={updateInput.title}
              valueChanged={(newValue) => setUpdateField({ title: newValue })}
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="primary-tag" messageId="primary_tag">
            <TagSelect
              onSelect={(t) => setUpdateField({ primary_tag_id: t[0]?.id })}
              ids={
                updateInput.primary_tag_id ? [updateInput.primary_tag_id] : []
              }
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
              existingIds={aggregateState.tagIds ?? []}
              mode={tagIds.mode}
              menuPortalTarget={document.body}
            />
          </BulkUpdateFormGroup>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
