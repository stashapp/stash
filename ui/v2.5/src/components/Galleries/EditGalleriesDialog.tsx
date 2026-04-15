import React, { useEffect, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useBulkGalleryUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregatePerformerIds,
  getAggregateStateObject,
  getAggregateTagIds,
  getAggregateStudioId,
  getAggregateSceneIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { IndeterminateCheckbox } from "../Shared/IndeterminateCheckbox";
import { BulkUpdateFormGroup, BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { BulkUpdateDateInput } from "../Shared/DateInput";
import { getDateError } from "src/utils/yup";

interface IListOperationProps {
  selected: GQL.SlimGalleryDataFragment[];
  onClose: (applied: boolean) => void;
}

const galleryFields = [
  "code",
  "rating100",
  "details",
  "organized",
  "photographer",
  "date",
];

export const EditGalleriesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const [updateInput, setUpdateInput] = useState<GQL.BulkGalleryUpdateInput>({
    ids: props.selected.map((gallery) => {
      return gallery.id;
    }),
  });

  const [performerIds, setPerformerIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [sceneIds, setSceneIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });

  const unsetDisabled = props.selected.length < 2;

  const [dateError, setDateError] = useState<string | undefined>();

  const [updateGalleries] = useBulkGalleryUpdate();

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const aggregateState = useMemo(() => {
    const updateState: Partial<GQL.BulkGalleryUpdateInput> = {};
    const state = props.selected;
    updateState.studio_id = getAggregateStudioId(props.selected);
    const updateTagIds = getAggregateTagIds(props.selected);
    const updatePerformerIds = getAggregatePerformerIds(props.selected);
    const updateSceneIds = getAggregateSceneIds(props.selected);
    let first = true;

    state.forEach((gallery: GQL.SlimGalleryDataFragment) => {
      getAggregateStateObject(updateState, gallery, galleryFields, first);
      first = false;
    });

    return {
      state: updateState,
      tagIds: updateTagIds,
      performerIds: updatePerformerIds,
      sceneIds: updateSceneIds,
    };
  }, [props.selected]);

  // update initial state from aggregate
  useEffect(() => {
    setUpdateInput((current) => ({ ...current, ...aggregateState.state }));
  }, [aggregateState]);

  useEffect(() => {
    setDateError(getDateError(updateInput.date ?? "", intl));
  }, [updateInput.date, intl]);

  function setUpdateField(input: Partial<GQL.BulkGalleryUpdateInput>) {
    setUpdateInput((current) => ({ ...current, ...input }));
  }

  function getGalleryInput(): GQL.BulkGalleryUpdateInput {
    const galleryInput: GQL.BulkGalleryUpdateInput = {
      ...updateInput,
      tag_ids: tagIds,
      performer_ids: performerIds,
      scene_ids: sceneIds,
    };

    // we don't have unset functionality for the rating star control
    // so need to determine if we are setting a rating or not
    galleryInput.rating100 = getAggregateInputValue(
      updateInput.rating100,
      aggregateState.state.rating100
    );

    return galleryInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateGalleries({ variables: { input: getGalleryInput() } });
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "galleries" }).toLocaleLowerCase(),
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
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_count_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "gallery" }),
            pluralEntity: intl.formatMessage({ id: "galleries" }),
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

          <BulkUpdateFormGroup name="photographer">
            <BulkUpdateTextInput
              value={updateInput.photographer}
              valueChanged={(newValue) =>
                setUpdateField({ photographer: newValue })
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

          <BulkUpdateFormGroup name="scenes" inline={false}>
            <MultiSet
              type={"scenes"}
              disabled={isUpdating}
              onUpdate={(itemIDs) => {
                setSceneIds((c) => ({ ...c, ids: itemIDs }));
              }}
              onSetMode={(newMode) => {
                setSceneIds((c) => ({ ...c, mode: newMode }));
              }}
              ids={sceneIds.ids ?? []}
              existingIds={aggregateState.sceneIds}
              mode={sceneIds.mode}
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
