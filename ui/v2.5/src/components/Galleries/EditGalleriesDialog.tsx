import React, { useEffect, useMemo, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkGalleryUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
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
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { BulkUpdateDateInput } from "../Shared/DateInput";

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
          { id: "dialogs.edit_entity_title" },
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

          {renderTextField("photographer", updateInput.photographer, (v) =>
            setUpdateField({ photographer: v })
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

          <Form.Group controlId="scenes">
            <Form.Label>
              <FormattedMessage id="scenes" />
            </Form.Label>
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
