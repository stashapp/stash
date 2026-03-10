import React, { useEffect, useMemo, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
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
import { BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { BulkUpdateDateInput } from "../Shared/DateInput";

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

  const [performerIds, setPerformerIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [groupIds, setGroupIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });

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
            singularEntity: intl.formatMessage({ id: "scene" }),
            pluralEntity: intl.formatMessage({ id: "scenes" }),
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

          {renderTextField("scene_code", updateInput.code, (newValue) =>
            setUpdateField({ code: newValue })
          )}
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

          <Form.Group controlId="performers">
            <Form.Label>
              <FormattedMessage id="performers" />
            </Form.Label>
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
          </Form.Group>

          <Form.Group controlId="groups">
            <Form.Label>
              <FormattedMessage id="groups" />
            </Form.Label>
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
            "details",
            updateInput.details,
            (v) => setUpdateField({ details: v }),
            true
          )}

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
