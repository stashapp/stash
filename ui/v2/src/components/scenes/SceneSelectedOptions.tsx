import _ from "lodash";
import {
  AnchorButton,
  Button,
  ButtonGroup,
  ControlGroup,
  FormGroup,
  HTMLSelect,
  InputGroup,
  Menu,
  MenuItem,
  Popover,
  Spinner,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent, SyntheticEvent, useEffect, useRef, useState } from "react";
import { FilterSelect } from "../select/FilterSelect";
import { FilterMultiSelect } from "../select/FilterMultiSelect";
import { StashService } from "../../core/StashService";
import * as GQL from "../../core/generated-graphql";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[],
  onScenesUpdated: () => void;
}

export const SceneSelectedOptions: FunctionComponent<IListOperationProps> = (props: IListOperationProps) => {
  const [rating, setRating] = useState<string>("");
  const [studioId, setStudioId] = useState<string | undefined>(undefined);
  const [performerIds, setPerformerIds] = useState<string[] | undefined>(undefined);
  const [tagIds, setTagIds] = useState<string[] | undefined>(undefined);

  const updateScenes = StashService.useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  function getSceneInput() : GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    var aggregateRating = getRating(props.selected);
    var aggregateStudioId = getStudioId(props.selected);
    var aggregatePerformerIds = getPerformerIds(props.selected);
    var aggregateTagIds = getTagIds(props.selected);

    var sceneInput : GQL.BulkSceneUpdateInput = {
      ids: props.selected.map((scene) => {
        return scene.id;
      })
    };

    // if rating is undefined 
    if (rating === "") {
      // and all scenes have the same rating, then we are unsetting the rating.
      if(aggregateRating) {
        // an undefined rating is ignored in the server, so set it to 0 instead
        sceneInput.rating = 0;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      sceneInput.rating = Number.parseInt(rating);
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
      const result = await updateScenes();
      ToastUtils.success("Updated scenes");
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    props.onScenesUpdated();
  }

  function getRating(state: GQL.SlimSceneDataFragment[]) {
    var ret : number | undefined;
    var first = true;

    state.forEach((scene : GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.rating;
        first = false;
      } else {
        if (ret !== scene.rating) {
          ret = undefined;
        }
      }
    });

    return ret;
  }

  function getStudioId(state: GQL.SlimSceneDataFragment[]) {
    var ret : string | undefined;
    var first = true;

    state.forEach((scene : GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.studio ? scene.studio.id : undefined;
        first = false;
      } else {
        if (scene.studio && ret != scene.studio.id) {
          ret = undefined;
        }
      }
    });

    return ret;
  }

  function toId(object : any) {
    return object.id;
  }

  function getPerformerIds(state: GQL.SlimSceneDataFragment[]) {
    var ret : string[] = [];
    var first = true;

    state.forEach((scene : GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = !!scene.performers ? scene.performers.map(toId).sort() : [];
        first = false;
      } else {
        const perfIds = !!scene.performers ? scene.performers.map(toId).sort() : [];
        
        if (!_.isEqual(performerIds, perfIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function getTagIds(state: GQL.SlimSceneDataFragment[]) {
    var ret : string[] = [];
    var first = true;

    state.forEach((scene : GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = !!scene.tags ? scene.tags.map(toId).sort() : [];
        first = false;
      } else {
        const tIds = !!scene.tags ? scene.tags.map(toId).sort() : [];
        
        if (!_.isEqual(ret, tIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function updateScenesEditState(state: GQL.SlimSceneDataFragment[]) {
    function toId(object : any) {
      return object.id;
    }

    var rating : string = "";
    var studioId : string | undefined;
    var performerIds : string[] = [];
    var tagIds : string[] = [];
    var first = true;

    state.forEach((scene : GQL.SlimSceneDataFragment) => {
      var thisRating = scene.rating ? scene.rating.toString() : "";
      var thisStudio = scene.studio ? scene.studio.id : undefined;

      if (first) {
        rating = thisRating;
        studioId = thisStudio;
        performerIds = !!scene.performers ? scene.performers.map(toId).sort() : [];
        tagIds = !!scene.tags ? scene.tags.map(toId).sort() : [];
        first = false;
      } else {
        if (rating !== thisRating) {
          rating = "";
        }
        if (studioId != thisStudio) {
          studioId = undefined;
        }
        const perfIds = !!scene.performers ? scene.performers.map(toId).sort() : [];
        const tIds = !!scene.tags ? scene.tags.map(toId).sort() : [];
        
        if (!_.isEqual(performerIds, perfIds)) {
          performerIds = [];
        }

        if (!_.isEqual(tagIds, tIds)) {
          tagIds = [];
        }
      }
    });
    
    setRating(rating);
    setStudioId(studioId);
    setPerformerIds(performerIds);
    setTagIds(tagIds);
  }

  useEffect(() => {
    updateScenesEditState(props.selected);
  }, [props.selected]);

  function renderMultiSelect(type: "performers" | "tags", initialIds: string[] | undefined) {
    return (
      <FilterMultiSelect
        type={type}
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          switch (type) {
            case "performers": setPerformerIds(ids); break;
            case "tags": setTagIds(ids); break;
          }
        }}
        initialIds={initialIds}
      />
    );
  }
  
  function render() {
    return (
      <>
        {isLoading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
        <div className="operation-container">
          <FormGroup className="operation-item" label="Rating">
            <HTMLSelect
              options={["", 1, 2, 3, 4, 5]}
              onChange={(event) => setRating(event.target.value)}
              value={rating}
            />
          </FormGroup>
          
          <FormGroup className="operation-item" label="Studio">
            <FilterSelect
              type="studios"
              onSelectItem={(item : any) => setStudioId(item ? item.id : undefined)}
              initialId={studioId}
            />
          </FormGroup>

          <FormGroup className="operation-item" label="Performers">
            {renderMultiSelect("performers", performerIds)}
          </FormGroup>

          <FormGroup className="operation-item" label="Tags">
            {renderMultiSelect("tags", tagIds)}
          </FormGroup>
          
          <ButtonGroup className="operation-item">
            <Button 
              intent="primary"
              onClick={() => onSave()}>
                Apply
            </Button>
          </ButtonGroup>
        </div>
      </>
    );
  }

  return render();
};
