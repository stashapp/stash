import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import _ from "lodash";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect, Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
import { RatingStars } from "./SceneDetails/RatingStars";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditScenesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [
    performerMode,
    setPerformerMode,
  ] = React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [movieMode, setMovieMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [movieIds, setMovieIds] = useState<string[]>();
  const [existingMovieIds, setExistingMovieIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();

  const [updateScenes] = useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function makeBulkUpdateIds(
    ids: string[],
    mode: GQL.BulkUpdateIdMode
  ): GQL.BulkUpdateIds {
    return {
      mode,
      ids,
    };
  }

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    const aggregateRating = getRating(props.selected);
    const aggregateStudioId = getStudioId(props.selected);
    const aggregatePerformerIds = getPerformerIds(props.selected);
    const aggregateTagIds = getTagIds(props.selected);
    const aggregateMovieIds = getMovieIds(props.selected);

    const sceneInput: GQL.BulkSceneUpdateInput = {
      ids: props.selected.map((scene) => {
        return scene.id;
      }),
    };

    // if rating is undefined
    if (rating === undefined) {
      // and all scenes have the same rating, then we are unsetting the rating.
      if (aggregateRating) {
        // null rating unsets it
        sceneInput.rating = null;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      sceneInput.rating = rating;
    }

    // if studioId is undefined
    if (studioId === undefined) {
      // and all scenes have the same studioId,
      // then unset the studioId, otherwise ignoring studioId
      if (aggregateStudioId) {
        // null studio_id unsets it
        sceneInput.studio_id = null;
      }
    } else {
      // if studioId is set, then we are setting it
      sceneInput.studio_id = studioId;
    }

    // if performerIds are empty
    if (
      performerMode === GQL.BulkUpdateIdMode.Set &&
      (!performerIds || performerIds.length === 0)
    ) {
      // and all scenes have the same ids,
      if (aggregatePerformerIds.length > 0) {
        // then unset the performerIds, otherwise ignore
        sceneInput.performer_ids = makeBulkUpdateIds(
          performerIds || [],
          performerMode
        );
      }
    } else {
      // if performerIds non-empty, then we are setting them
      sceneInput.performer_ids = makeBulkUpdateIds(
        performerIds || [],
        performerMode
      );
    }

    // if tagIds non-empty, then we are setting them
    if (
      tagMode === GQL.BulkUpdateIdMode.Set &&
      (!tagIds || tagIds.length === 0)
    ) {
      // and all scenes have the same ids,
      if (aggregateTagIds.length > 0) {
        // then unset the tagIds, otherwise ignore
        sceneInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
      }
    } else {
      // if tagIds non-empty, then we are setting them
      sceneInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
    }

    // if movieIds non-empty, then we are setting them
    if (
      movieMode === GQL.BulkUpdateIdMode.Set &&
      (!movieIds || movieIds.length === 0)
    ) {
      // and all scenes have the same ids,
      if (aggregateMovieIds.length > 0) {
        // then unset the movieIds, otherwise ignore
        sceneInput.movie_ids = makeBulkUpdateIds(movieIds || [], movieMode);
      }
    } else {
      // if movieIds non-empty, then we are setting them
      sceneInput.movie_ids = makeBulkUpdateIds(movieIds || [], movieMode);
    }

    if (organized !== undefined) {
      sceneInput.organized = organized;
    }

    return sceneInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateScenes();
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "scenes" }).toLocaleLowerCase() }
        ),
      });
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
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
        ret = scene.performers ? scene.performers.map((p) => p.id).sort() : [];
        first = false;
      } else {
        const perfIds = scene.performers
          ? scene.performers.map((p) => p.id).sort()
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
        ret = scene.tags ? scene.tags.map((t) => t.id).sort() : [];
        first = false;
      } else {
        const tIds = scene.tags ? scene.tags.map((t) => t.id).sort() : [];

        if (!_.isEqual(ret, tIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  function getMovieIds(state: GQL.SlimSceneDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      if (first) {
        ret = scene.movies ? scene.movies.map((m) => m.movie.id).sort() : [];
        first = false;
      } else {
        const mIds = scene.movies
          ? scene.movies.map((m) => m.movie.id).sort()
          : [];

        if (!_.isEqual(ret, mIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateMovieIds: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      const sceneRating = scene.rating;
      const sceneStudioID = scene?.studio?.id;
      const scenePerformerIDs = (scene.performers ?? [])
        .map((p) => p.id)
        .sort();
      const sceneTagIDs = (scene.tags ?? []).map((p) => p.id).sort();
      const sceneMovieIDs = (scene.movies ?? []).map((m) => m.movie.id).sort();

      if (first) {
        updateRating = sceneRating ?? undefined;
        updateStudioID = sceneStudioID;
        updatePerformerIds = scenePerformerIDs;
        updateTagIds = sceneTagIDs;
        updateMovieIds = sceneMovieIDs;
        first = false;
        updateOrganized = scene.organized;
      } else {
        if (sceneRating !== updateRating) {
          updateRating = undefined;
        }
        if (sceneStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!_.isEqual(scenePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!_.isEqual(sceneTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (!_.isEqual(sceneMovieIDs, updateMovieIds)) {
          updateMovieIds = [];
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
    setExistingMovieIds(updateMovieIds);
    setOrganized(updateOrganized);
  }, [props.selected, performerMode, tagMode, movieMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  function renderMultiSelect(
    type: "performers" | "tags" | "movies",
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
      case "movies":
        mode = movieMode;
        existingIds = existingMovieIds;
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
            case "movies":
              setMovieIds(itemIDs);
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
            case "movies":
              setMovieMode(newMode);
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
      <Modal
        show
        icon="pencil-alt"
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
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
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

          <Form.Group controlId="movies">
            <Form.Label>
              <FormattedMessage id="movies" />
            </Form.Label>
            {renderMultiSelect("movies", movieIds)}
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
      </Modal>
    );
  }

  return render();
};
