import * as GQL from "../core/generated-graphql";
import { PerformersCriterion } from "../models/list-filter/criteria/performers";
import { StudiosCriterion } from "../models/list-filter/criteria/studios";
import { TagsCriterion } from "../models/list-filter/criteria/tags";
import { ListFilterModel } from "../models/list-filter/filter";
import { FilterMode } from "../models/list-filter/types";

const makePerformerScenesUrl = (
  performer: Partial<GQL.PerformerDataFragment>
) => {
  if (!performer.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new PerformersCriterion();
  criterion.value = [
    { id: performer.id, label: performer.name || `Performer ${performer.id}` }
  ];
  filter.criteria.push(criterion);
  return `/scenes?${filter.makeQueryParameters()}`;
};

const makeStudioScenesUrl = (studio: Partial<GQL.StudioDataFragment>) => {
  if (!studio.id) return "#";
  const filter = new ListFilterModel(FilterMode.Scenes);
  const criterion = new StudiosCriterion();
  criterion.value = [
    { id: studio.id, label: studio.name || `Studio ${studio.id}` }
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

const Nav = {
  makePerformerScenesUrl,
  makeStudioScenesUrl,
  makeTagSceneMarkersUrl,
  makeTagScenesUrl,
  makeSceneMarkerUrl
};
export default Nav;
