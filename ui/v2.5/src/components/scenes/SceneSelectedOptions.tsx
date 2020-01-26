import React, { useEffect, useState } from "react";
import { Button, ButtonGroup, Form, Spinner } from "react-bootstrap";
import _ from "lodash";
import { StashService } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { FilterSelect, StudioSelect } from "src/components/Shared";
import { useToast } from "src/hooks";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[];
  onScenesUpdated: () => void;
}

export const SceneSelectedOptions: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const Toast = useToast();
  const [rating, setRating] = useState<string>("");
  const [studioId, setStudioId] = useState<string>();
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [tagIds, setTagIds] = useState<string[]>();

  const [updateScenes] = StashService.useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isLoading, setIsLoading] = useState(true);

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    const aggregateRating = getRating(props.selected);
    const aggregateStudioId = getStudioId(props.selected);
    const aggregatePerformerIds = getPerformerIds(props.selected);
    const aggregateTagIds = getTagIds(props.selected);

    const sceneInput: GQL.BulkSceneUpdateInput = {
      ids: props.selected.map(scene => {
        return scene.id;
      })
    };

    // if rating is undefined
    if (rating === "") {
      // and all scenes have the same rating, then we are unsetting the rating.
      if (aggregateRating) {
        // an undefined rating is ignored in the server, so set it to 0 instead
        sceneInput.rating = 0;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      sceneInput.rating = Number.parseInt(rating, 10);
    }

    // if studioId is undefined
    if (studioId === undefined) {
      // and all scenes have the same studioId,
      // then unset the studioId, otherwise ignoring studioId
      if (aggregateStudioId) {
        // an undefined studio_id is ignored in the server, so set it to empty string instead
        sceneInput.studio_id = "";
      }
    } else {
      // if studioId is set, then we are setting it
      sceneInput.studio_id = studioId;
    }

    // if performerIds are empty
    if (!performerIds || performerIds.length === 0) {
      // and all scenes have the same ids,
      if (aggregatePerformerIds.length > 0) {
        // then unset the performerIds, otherwise ignore
        sceneInput.performer_ids = performerIds;
      }
    } else {
      // if performerIds non-empty, then we are setting them
      sceneInput.performer_ids = performerIds;
    }

    // if tagIds non-empty, then we are setting them
    if (!tagIds || tagIds.length === 0) {
      // and all scenes have the same ids,
      if (aggregateTagIds.length > 0) {
        // then unset the tagIds, otherwise ignore
        sceneInput.tag_ids = tagIds;
      }
    } else {
      // if tagIds non-empty, then we are setting them
      sceneInput.tag_ids = tagIds;
    }

    return sceneInput;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      await updateScenes();
      Toast.success({ content: "Updated scenes" });
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
    props.onScenesUpdated();
  }

  function getRating(state: GQL.SlimSceneDataFragment[]) {
    let ret: number | undefined;
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.rating ?? undefined;
        first = false;
      } else if (ret !== scene.rating) {
        ret = undefined;
      }
    });

    return ret;
  }

  function getStudioId(state: GQL.SlimSceneDataFragment[]) {
    let ret: string | undefined;
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene?.studio?.id;
        first = false;
      } else {
        const studio = scene?.studio?.id;
        if (ret !== studio) {
          ret = undefined;
        }
      }
    });

    return ret;
  }

  function getPerformerIds(state: GQL.SlimSceneDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.performers ? scene.performers.map(p => p.id).sort() : [];
        first = false;
      } else {
        const perfIds = scene.performers
          ? scene.performers.map(p => p.id).sort()
          : [];

        if (!_.isEqual(ret, perfIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function getTagIds(state: GQL.SlimSceneDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.tags ? scene.tags.map(t => t.id).sort() : [];
        first = false;
      } else {
        const tIds = scene.tags ? scene.tags.map(t => t.id).sort() : [];

        if (!_.isEqual(ret, tIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function updateScenesEditState(state: GQL.SlimSceneDataFragment[]) {
    let updateRating = "";
    let updateStudioId: string|undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      const thisRating = scene.rating?.toString() ?? "";
      const thisStudio = scene?.studio?.id;

      if (first) {
        updateRating = thisRating;
        updateStudioId = thisStudio;
        updatePerformerIds = scene.performers
          ? scene.performers.map(p => p.id).sort()
          : [];
        updateTagIds = scene.tags ? scene.tags.map(p => p.id).sort() : [];
        first = false;
      } else {
        if (rating !== thisRating) {
          updateRating = "";
        }
        if (studioId !== thisStudio) {
          updateStudioId = undefined;
        }
        const perfIds = scene.performers
          ? scene.performers.map(p => p.id).sort()
          : [];
        const tIds = scene.tags ? scene.tags.map(t => t.id).sort() : [];

        if (!_.isEqual(performerIds, perfIds)) {
          updatePerformerIds = [];
        }

        if (!_.isEqual(tagIds, tIds)) {
          updateTagIds = [];
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioId);
    setPerformerIds(updatePerformerIds);
    setTagIds(updateTagIds);
  }

  useEffect(() => {
    updateScenesEditState(props.selected);
    setIsLoading(false);
  }, [props.selected]);

  function renderMultiSelect(
    type: "performers" | "tags",
    ids: string[] | undefined
  ) {
    return (
      <FilterSelect
        type={type}
        isMulti
        isClearable={false}
        onSelect={items => {
          const itemIDs = items.map(i => i.id);
          switch (type) {
            case "performers":
              setPerformerIds(itemIDS);
              break;
            case "tags":
              setTagIds(itemIDs);
              break;
          }
        }}
        ids={ids ?? []}
      />
    );
  }

  if(isLoading)
    return <Spinner animation="border" variant="light" />;

  function render() {
    return (
      <div className="operation-container">
        <Form.Group controlId="rating" className="operation-item rating-operation">
          <Form.Label>Rating</Form.Label>
          <Form.Control
            as="select"
            onChange={(event: any) => setRating(event.target.value)}
          >
            {["", '1', '2', '3', '4', '5'].map(opt => (
              <option selected={opt === rating} value={opt}>
                {opt}
              </option>
            ))}
          </Form.Control>
        </Form.Group>

        <Form.Group controlId="studio" className="operation-item">
          <Form.Label>Studio</Form.Label>
          <StudioSelect
            onSelect={items => setStudioId(items[0]?.id)}
            ids={studioId ? [studioId] : []}
          />
        </Form.Group>

        <Form.Group className="operation-item" controlId="performers">
          <Form.Label>Performers</Form.Label>
          {renderMultiSelect("performers", performerIds)}
        </Form.Group>

        <Form.Group className="operation-item" controlId="performers">
          <Form.Label>Tags</Form.Label>
          {renderMultiSelect("tags", tagIds)}
        </Form.Group>

        <Button variant="primary" onClick={onSave} className="apply-operation">
          Apply
        </Button>
      </div>
    );
  }

  return render();
};
