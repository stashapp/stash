import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneList } from "src/components/Scenes/SceneList";
import { TagsCriterion } from "src/models/list-filter/criteria/tags";

interface ITagScenesPanel {
  tag: Partial<GQL.TagDataFragment>;
}

export const TagScenesPanel: React.FC<ITagScenesPanel> = ({ tag }) => {
  function filterHook(filter: ListFilterModel) {
    const tagValue = { id: tag.id!, label: tag.name! };
    // if tag is already present, then we modify it, otherwise add
    let tagCriterion = filter.criteria.find((c) => {
      return c.type === "tags";
    }) as TagsCriterion;

    if (
      tagCriterion &&
      (tagCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        tagCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the tag if not present
      if (
        !tagCriterion.value.find((p) => {
          return p.id === tag.id;
        })
      ) {
        tagCriterion.value.push(tagValue);
      }

      tagCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      tagCriterion = new TagsCriterion("tags");
      tagCriterion.value = [tagValue];
      filter.criteria.push(tagCriterion);
    }

    return filter;
  }

  return <SceneList subComponent filterHook={filterHook} />;
};
