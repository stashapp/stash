import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MoviesCriterion } from "src/models/list-filter/criteria/movies";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneList } from "src/components/Scenes/SceneList";
import { View } from "src/components/List/views";

interface IGroupScenesPanel {
  active: boolean;
  group: GQL.MovieDataFragment;
}

export const GroupScenesPanel: React.FC<IGroupScenesPanel> = ({
  active,
  group,
}) => {
  function filterHook(filter: ListFilterModel) {
    const movieValue = { id: group.id, label: group.name };
    // if movie is already present, then we modify it, otherwise add
    let movieCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "movies";
    }) as MoviesCriterion | undefined;

    if (
      movieCriterion &&
      (movieCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        movieCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the movie if not present
      if (
        !movieCriterion.value.find((p) => {
          return p.id === group.id;
        })
      ) {
        movieCriterion.value.push(movieValue);
      }

      movieCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      movieCriterion = new MoviesCriterion();
      movieCriterion.value = [movieValue];
      filter.criteria.push(movieCriterion);
    }

    return filter;
  }

  if (group && group.id) {
    return (
      <SceneList
        filterHook={filterHook}
        defaultSort="movie_scene_number"
        alterQuery={active}
        view={View.GroupScenes}
      />
    );
  }
  return <></>;
};
