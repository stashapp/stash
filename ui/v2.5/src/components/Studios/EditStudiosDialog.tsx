import React, { useEffect, useMemo, useState } from "react";
import { Col, Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkStudioUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "../Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { MultiSet } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregateState,
  getAggregateStateObject,
} from "src/utils/bulkUpdate";
import { IndeterminateCheckbox } from "../Shared/IndeterminateCheckbox";
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import * as FormUtils from "src/utils/form";
import { StudioSelect } from "../Shared/Select";

interface IListOperationProps {
  selected: GQL.SlimStudioDataFragment[];
  onClose: (applied: boolean) => void;
}

const studioFields = ["favorite", "rating100", "details", "ignore_auto_tag"];

export const EditStudiosDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const [updateInput, setUpdateInput] = useState<GQL.BulkStudioUpdateInput>({
    ids: props.selected.map((studio) => {
      return studio.id;
    }),
  });

  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });

  const [updateStudios] = useBulkStudioUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const aggregateState = useMemo(() => {
    const updateState: Partial<GQL.BulkStudioUpdateInput> = {};
    const state = props.selected;
    let updateTagIds: string[] = [];
    let first = true;

    state.forEach((studio: GQL.SlimStudioDataFragment) => {
      getAggregateStateObject(updateState, studio, studioFields, first);

      // studio data fragment doesn't have parent_id, so handle separately
      updateState.parent_id = getAggregateState(
        updateState.parent_id,
        studio.parent_studio?.id,
        first
      );

      const studioTagIDs = (studio.tags ?? []).map((p) => p.id).sort();

      updateTagIds = getAggregateState(updateTagIds, studioTagIDs, first) ?? [];

      first = false;
    });

    return { state: updateState, tagIds: updateTagIds };
  }, [props.selected]);

  // update initial state from aggregate
  useEffect(() => {
    setUpdateInput((current) => ({ ...current, ...aggregateState.state }));
  }, [aggregateState]);

  function setUpdateField(input: Partial<GQL.BulkStudioUpdateInput>) {
    setUpdateInput((current) => ({ ...current, ...input }));
  }

  function getStudioInput(): GQL.BulkStudioUpdateInput {
    const studioInput: GQL.BulkStudioUpdateInput = {
      ...updateInput,
      tag_ids: tagIds,
    };

    // we don't have unset functionality for the rating star control
    // so need to determine if we are setting a rating or not
    studioInput.rating100 = getAggregateInputValue(
      updateInput.rating100,
      aggregateState.state.rating100
    );

    return studioInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateStudios({
        variables: {
          input: getStudioInput(),
        },
      });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "studios" }).toLocaleLowerCase(),
          }
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
    return (
      <Form.Group controlId={name}>
        <Form.Label>
          <FormattedMessage id={name} />
        </Form.Label>
        <BulkUpdateTextInput
          value={value === null ? "" : value ?? undefined}
          valueChanged={(newValue) => setter(newValue)}
          unsetDisabled={props.selected.length < 2}
          as={area ? "textarea" : undefined}
        />
      </Form.Group>
    );
  }

  function render() {
    return (
      <ModalComponent
        dialogClassName="edit-studios-dialog"
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "actions.edit_entity" },
          { entityType: intl.formatMessage({ id: "studios" }) }
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
        <Form.Group controlId="parent-studio" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "parent_studio" }),
          })}
          <Col xs={9}>
            <StudioSelect
              onSelect={(items) =>
                setUpdateField({
                  parent_id: items.length > 0 ? items[0]?.id : undefined,
                })
              }
              ids={updateInput.parent_id ? [updateInput.parent_id] : []}
              isDisabled={isUpdating}
              menuPortalTarget={document.body}
            />
          </Col>
        </Form.Group>
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
        <Form>
          <Form.Group controlId="favorite">
            <IndeterminateCheckbox
              setChecked={(checked) => setUpdateField({ favorite: checked })}
              checked={updateInput.favorite ?? undefined}
              label={intl.formatMessage({ id: "favourite" })}
            />
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            <MultiSet
              type="tags"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setTagIds((v) => ({ ...v, ids: itemIDs }))}
              onSetMode={(newMode) =>
                setTagIds((v) => ({ ...v, mode: newMode }))
              }
              existingIds={aggregateState.tagIds ?? []}
              ids={tagIds.ids ?? []}
              mode={tagIds.mode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          {renderTextField(
            "details",
            updateInput.details,
            (newValue) => setUpdateField({ details: newValue }),
            true
          )}

          <Form.Group controlId="ignore-auto-tags">
            <IndeterminateCheckbox
              label={intl.formatMessage({ id: "ignore_auto_tag" })}
              setChecked={(checked) =>
                setUpdateField({ ignore_auto_tag: checked })
              }
              checked={updateInput.ignore_auto_tag ?? undefined}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
