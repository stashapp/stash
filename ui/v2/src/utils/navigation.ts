import * as GQL from "../core/generated-graphql";
import { PerformersCriterion } from "../models/list-filter/criteria/performers";
import { StudiosCriterion } from "../models/list-filter/criteria/studios";
import { MoviesCriterion } from "../models/list-filter/criteria/movies";
import { TagsCriterion } from "../models/list-filter/criteria/tags";
import { ListFilterModel } from "../models/list-filter/filter";
import { FilterMode } from "../models/list-filter/types";

export class NavigationUtils {
  public static makePerformerScenesUrl(performer: Partial<GQL.PerformerDataFragment>): string {
    if (performer.id === undefined) { return "#"; }
    const filter = new ListFilterModel(FilterMode.Scenes);
    const criterion = new PerformersCriterion();
    criterion.value = [{ id: performer.id, label: performer.name || `Performer ${performer.id}` }];
    filter.criteria.push(criterion);
    return `/scenes?${filter.makeQueryParameters()}`;
  }

  public static makeStudioScenesUrl(studio: Partial<GQL.StudioDataFragment>): string {
    if (studio.id === undefined) { return "#"; }
    const filter = new ListFilterModel(FilterMode.Scenes);
    const criterion = new StudiosCriterion();
    criterion.value = [{ id: studio.id, label: studio.name || `Studio ${studio.id}` }];
    filter.criteria.push(criterion);
    return `/scenes?${filter.makeQueryParameters()}`;
  }

  public static makeMovieScenesUrl(movie: Partial<GQL.MovieDataFragment>): string {
    if (movie.id === undefined) { return "#"; }
    const filter = new ListFilterModel(FilterMode.Scenes);
    const criterion = new MoviesCriterion();
    criterion.value = [{ id: movie.id, label: movie.name || `Movie ${movie.id}` }];
    filter.criteria.push(criterion);
    return `/scenes?${filter.makeQueryParameters()}`;
  }

  public static makeTagScenesUrl(tag: Partial<GQL.TagDataFragment>): string {
    if (tag.id === undefined) { return "#"; }
    const filter = new ListFilterModel(FilterMode.Scenes);
    const criterion = new TagsCriterion("tags");
    criterion.value = [{ id: tag.id, label: tag.name || `Tag ${tag.id}` }];
    filter.criteria.push(criterion);
    return `/scenes?${filter.makeQueryParameters()}`;
  }

  public static makeTagSceneMarkersUrl(tag: Partial<GQL.TagDataFragment>): string {
    if (tag.id === undefined) { return "#"; }
    const filter = new ListFilterModel(FilterMode.SceneMarkers);
    const criterion = new TagsCriterion("tags");
    criterion.value = [{ id: tag.id, label: tag.name || `Tag ${tag.id}` }];
    filter.criteria.push(criterion);
    return `/scenes/markers?${filter.makeQueryParameters()}`;
  }

  public static makeSceneMarkerUrl(sceneMarker: Partial<GQL.SceneMarkerDataFragment>): string {
    if (sceneMarker.id === undefined || sceneMarker.scene === undefined) { return "#"; }
    return `/scenes/${sceneMarker.scene.id}?t=${sceneMarker.seconds}`;
  }
}
