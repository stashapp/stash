import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregateGroupIds,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditScenesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [performerMode, setPerformerMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [groupMode, setGroupMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [groupIds, setGroupIds] = useState<string[]>();
  const [existingGroupIds, setExistingGroupIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();

  const [updateScenes] = useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateGroupIds = getAggregateGroupIds(props.selected);

    const sceneInput: GQL.BulkSceneUpdateInput = {
      ids: props.selected.map((scene) => {
        return scene.id;
      }),
    };

    sceneInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    sceneInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);

    sceneInput.performer_ids = getAggregateInputIDs(
      performerMode,
      performerIds,
      aggregatePerformerIds
    );
    sceneInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);
    sceneInput.movie_ids = getAggregateInputIDs(
      groupMode,
      groupIds,
      aggregateGroupIds
    );

    if (organized !== undefined) {
      sceneInput.organized = organized;
    }

    return sceneInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateScenes();
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

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateGroupIds: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      const sceneRating = scene.rating100;
      const sceneStudioID = scene?.studio?.id;
      const scenePerformerIDs = (scene.performers ?? [])
        .map((p) => p.id)
        .sort();
      const sceneTagIDs = (scene.tags ?? []).map((p) => p.id).sort();
      const sceneGroupIDs = (scene.movies ?? []).map((m) => m.movie.id).sort();

      if (first) {
        updateRating = sceneRating ?? undefined;
        updateStudioID = sceneStudioID;
        updatePerformerIds = scenePerformerIDs;
        updateTagIds = sceneTagIDs;
        updateGroupIds = sceneGroupIDs;
        first = false;
        updateOrganized = scene.organized;
      } else {
        if (sceneRating !== updateRating) {
          updateRating = undefined;
        }
        if (sceneStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!isEqual(scenePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!isEqual(sceneTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (!isEqual(sceneGroupIDs, updateGroupIds)) {
          updateGroupIds = [];
        }
        if (scene.organized !== updateOrganized) {
          updateOrganized = undefined;
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioID);
    setExistingPerformerIds(updatePerformerIds);
    setExistingTagIds(updateTagIds);
    setExistingGroupIds(updateGroupIds);
    setOrganized(updateOrganized);
  }, [props.selected]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  function renderMultiSelect(
    type: "performers" | "tags" | "groups",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    let existingIds: string[] | undefined = [];
    switch (type) {
      case "performers":
        mode = performerMode;
        existingIds = existingPerformerIds;
        break;
      case "tags":
        mode = tagMode;
        existingIds = existingTagIds;
        break;
      case "groups":
        mode = groupMode;
        existingIds = existingGroupIds;
        break;
    }

    return (
      <MultiSet
        type={type}
        disabled={isUpdating}
        onUpdate={(itemIDs) => {
          switch (type) {
            case "performers":
              setPerformerIds(itemIDs);
              break;
            case "tags":
              setTagIds(itemIDs);
              break;
            case "groups":
              setGroupIds(itemIDs);
              break;
          }
        }}
        onSetMode={(newMode) => {
          switch (type) {
            case "performers":
              setPerformerMode(newMode);
              break;
            case "tags":
              setTagMode(newMode);
              break;
            case "groups":
              setGroupMode(newMode);
              break;
          }
        }}
        ids={ids ?? []}
        existingIds={existingIds ?? []}
        mode={mode}
      />
    );
  }

  function cycleOrganized() {
    if (organized) {
      setOrganized(undefined);
    } else if (organized === undefined) {
      setOrganized(false);
    } else {
      setOrganized(true);
    }
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

          <Form.Group controlId="performers">
            <Form.Label>
              <FormattedMessage id="performers" />
            </Form.Label>
            {renderMultiSelect("performers", performerIds)}
          </Form.Group>

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            {renderMultiSelect("tags", tagIds)}
          </Form.Group>

          <Form.Group controlId="groups">
            <Form.Label>
              <FormattedMessage id="groups" />
            </Form.Label>
            {renderMultiSelect("groups", groupIds)}
          </Form.Group>

          <Form.Group controlId="organized">
            <Form.Check
              type="checkbox"
              label={intl.formatMessage({ id: "organized" })}
              checked={organized}
              ref={checkboxRef}
              onChange={() => cycleOrganized()}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
