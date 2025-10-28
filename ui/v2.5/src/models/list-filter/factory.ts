import { FilterMode } from "src/core/generated-graphql";
import { ListFilterOptions } from "./filter-options";
import { GalleryListFilterOptions } from "./galleries";
import { ImageListFilterOptions } from "./images";
import { GroupListFilterOptions } from "./groups";
import { PerformerListFilterOptions } from "./performers";
import { SceneMarkerListFilterOptions } from "./scene-markers";
import { SceneListFilterOptions } from "./scenes";
import { StudioListFilterOptions } from "./studios";
import { TagListFilterOptions } from "./tags";

export function getFilterOptions(mode: FilterMode): ListFilterOptions {
  switch (mode) {
    case FilterMode.Scenes:
      return SceneListFilterOptions;
    case FilterMode.Performers:
      return PerformerListFilterOptions;
    case FilterMode.Studios:
      return StudioListFilterOptions;
    case FilterMode.Galleries:
      return GalleryListFilterOptions;
    case FilterMode.SceneMarkers:
      return SceneMarkerListFilterOptions;
    case FilterMode.Movies:
    case FilterMode.Groups:
      return GroupListFilterOptions;
    case FilterMode.Tags:
      return TagListFilterOptions;
    case FilterMode.Images:
      return ImageListFilterOptions;
  }
}
