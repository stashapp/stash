import * as GQL from "src/core/generated-graphql";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import {
  StudiosCriterion,
  ParentStudiosCriterion,
} from "src/models/list-filter/criteria/studios";
import { TagsCriterion } from "src/models/list-filter/criteria/tags";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterMode } from "src/models/list-filter/types";
import { MoviesCriterion } from "src/models/list-filter/criteria/movies";

const makePerformerScenesUrl = (
  performer: Partial<GQL.PerformerDataFragment>
) => {
  if (!performer.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new PerformersCriterion();
  criterion.value = [
    { id: performer.id, label: performer.name || `Performer ${performer.id}` },
  ];
  filter.criteria.push(criterion);
  return `/scenes?${filter.makeQueryParameters()}`;
};

const makeStudioScenesUrl = (studio: Partial<GQL.StudioDataFragment>) => {
  if (!studio.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new StudiosCriterion();
  criterion.value = [
    { id: studio.id, label: studio.name || `Studio ${studio.id}` },
  ];
  filter.criteria.push(criterion);
  return `/scenes?${filter.makeQueryParameters()}`;
};

const makeChildStudiosUrl = (studio: Partial<GQL.StudioDataFragment>) => {
  if (!studio.id) return "#";
  const filter = new ListFilterModel(FilterMode.Studios);
  const criterion = new ParentStudiosCriterion();
  criterion.value = [
    { id: studio.id, label: studio.name || `Studio ${studio.id}` },
  ];
  filter.criteria.push(criterion);
  return `/studios?${filter.makeQueryParameters()}`;
};

const makeMovieScenesUrl = (movie: Partial<GQL.MovieDataFragment>) => {
  if (!movie.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new MoviesCriterion();
  criterion.value = [
    { id: movie.id, label: movie.name || `Movie ${movie.id}` },
  ];
  filter.criteria.push(criterion);
  return `/scenes?${filter.makeQueryParameters()}`;
};

const makeTagScenesUrl = (tag: Partial<GQL.TagDataFragment>) => {
  if (!tag.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new TagsCriterion("tags");
  criterion.value = [{ id: tag.id, label: tag.name || `Tag ${tag.id}` }];
  filter.criteria.push(criterion);
  return `/scenes?${filter.makeQueryParameters()}`;
};

const makeTagSceneMarkersUrl = (tag: Partial<GQL.TagDataFragment>) => {
  if (!tag.id) return "#";
  const filter = new ListFilterModel(FilterMode.SceneMarkers);
  const criterion = new TagsCriterion("tags");
  criterion.value = [{ id: tag.id, label: tag.name || `Tag ${tag.id}` }];
  filter.criteria.push(criterion);
  return `/scenes/markers?${filter.makeQueryParameters()}`;
};

const makeSceneMarkerUrl = (
  sceneMarker: Partial<GQL.SceneMarkerDataFragment>
) => {
  if (!sceneMarker.id || !sceneMarker.scene) return "#";
  return `/scenes/${sceneMarker.scene.id}?t=${sceneMarker.seconds}`;
};

export default {
  makePerformerScenesUrl,
  makeStudioScenesUrl,
  makeTagSceneMarkersUrl,
  makeTagScenesUrl,
  makeSceneMarkerUrl,
  makeMovieScenesUrl,
  makeChildStudiosUrl,
};
