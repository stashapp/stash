import * as GQL from "src/core/generated-graphql";
import { SceneListFilterOptions } from "../scenes";
import { MovieListFilterOptions } from "../movies";
import { GalleryListFilterOptions } from "../galleries";
import { PerformerListFilterOptions } from "../performers";
import { ImageListFilterOptions } from "../images";
import { SceneMarkerListFilterOptions } from "../scene-markers";
import { StudioListFilterOptions } from "../studios";
import { TagListFilterOptions } from "../tags";
import { CriterionType } from "../types";

const filterModeOptions = {
  [GQL.FilterMode.Galleries]: GalleryListFilterOptions.criterionOptions,
  [GQL.FilterMode.Images]: ImageListFilterOptions.criterionOptions,
  [GQL.FilterMode.Movies]: MovieListFilterOptions.criterionOptions,
  [GQL.FilterMode.Performers]: PerformerListFilterOptions.criterionOptions,
  [GQL.FilterMode.SceneMarkers]: SceneMarkerListFilterOptions.criterionOptions,
  [GQL.FilterMode.Scenes]: SceneListFilterOptions.criterionOptions,
  [GQL.FilterMode.Studios]: StudioListFilterOptions.criterionOptions,
  [GQL.FilterMode.Tags]: TagListFilterOptions.criterionOptions,
};

export function makeCriteria(
  mode: GQL.FilterMode,
  type: CriterionType = "none"
) {
  const criterionOptions = filterModeOptions[mode];

  const option = criterionOptions.find((o) => o.type === type);

  if (!option) {
    throw new Error(`Unknown criterion parameter name: ${type}`);
  }

  return option?.makeCriterion();
}
