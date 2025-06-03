import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { PerformerList } from "src/components/Performers/PerformerList";
import { FilteredSceneList } from "src/components/Scenes/SceneList";
import { View } from "src/components/List/views";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  TagsCriterion,
  PerformerTagsCriterionOption,
} from "src/models/list-filter/criteria/tags";

interface ITagPerformersPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
  showPerformerScenes?: boolean;
}

export const TagPerformersPanel: React.FC<ITagPerformersPanel> = ({
  active,
  tag,
  showSubTagContent,
  showPerformerScenes,
}) => {
  // Hook for filtering performers by tag
  const performerFilterHook = useTagFilterHook(tag, showSubTagContent);

  // Custom filter hook for finding scenes with performers that have this tag
  const performerScenesFilterHook = (baseFilter: ListFilterModel) => {
    // Clone to avoid mutating original
    const newFilter = baseFilter.clone();

    // Remove any existing performer_tags criteria to prevent duplicates
    newFilter.criteria = newFilter.criteria.filter(
      (c) => c.criterionOption.type !== "performer_tags"
    );

    const tagValue = { id: tag.id, label: tag.name };

    const performerTagsCriterion = new TagsCriterion(
      PerformerTagsCriterionOption
    );
    performerTagsCriterion.value = {
      items: [tagValue],
      excluded: [],
      depth: showSubTagContent ? -1 : 0,
    };
    performerTagsCriterion.modifier = GQL.CriterionModifier.IncludesAll;

    newFilter.criteria.push(performerTagsCriterion);

    return newFilter;
  };

  if (showPerformerScenes) {
    return (
      <FilteredSceneList
        filterHook={performerScenesFilterHook}
        alterQuery={false}
        view={View.TagScenes}
      />
    );
  }

  return (
    <PerformerList
      filterHook={performerFilterHook}
      alterQuery={active}
      view={View.TagPerformers}
    />
  );
};
