import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import {
  ContainingGroupsCriterionOption,
  GroupsCriterion,
} from "src/models/list-filter/criteria/groups";
import { ListFilterModel } from "src/models/list-filter/filter";
import TextUtils from "src/utils/text";

export const scrapedGroupToCreateInput = (toCreate: GQL.ScrapedGroup) => {
  const input: GQL.GroupCreateInput = {
    name: toCreate.name ?? "",
    urls: toCreate.urls,
    aliases: toCreate.aliases,
    front_image: toCreate.front_image,
    back_image: toCreate.back_image,
    synopsis: toCreate.synopsis,
    date: toCreate.date,
    director: toCreate.director,
    // #788 - convert duration and rating to the correct type
    duration: TextUtils.timestampToSeconds(toCreate.duration),
    studio_id: toCreate.studio?.stored_id,
    rating100: parseInt(toCreate.rating ?? "0", 10) * 20,
  };

  if (!input.duration) {
    input.duration = undefined;
  }

  if (!input.rating100 || Number.isNaN(input.rating100)) {
    input.rating100 = undefined;
  }

  return input;
};

export const useContainingGroupFilterHook = (
  group: Pick<GQL.StudioDataFragment, "id" | "name">
) => {
  const { configuration } = React.useContext(ConfigurationContext);
  return (filter: ListFilterModel) => {
    const groupValue = { id: group.id, label: group.name };
    // if studio is already present, then we modify it, otherwise add
    let groupCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "containing_groups";
    }) as GroupsCriterion | undefined;

    if (groupCriterion) {
      // add the group if not present
      if (
        !groupCriterion.value.items.find((p) => {
          return p.id === group.id;
        })
      ) {
        groupCriterion.value.items.push(groupValue);
      }
    } else {
      groupCriterion = new GroupsCriterion(ContainingGroupsCriterionOption);
      groupCriterion.value = {
        items: [groupValue],
        excluded: [],
        depth: configuration?.ui.showChildStudioContent ? -1 : 0,
      };
      groupCriterion.modifier = GQL.CriterionModifier.Includes;
      filter.criteria.push(groupCriterion);
    }

    return filter;
  };
};
