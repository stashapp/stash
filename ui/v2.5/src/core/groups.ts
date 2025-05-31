import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import {
  GroupsCriterion,
  GroupsCriterionOption,
} from "src/models/list-filter/criteria/groups";
import { ListFilterModel } from "src/models/list-filter/filter";

export const useGroupFilterHook = (
  group: GQL.GroupDataFragment,
  showChildGroupContent?: boolean
) => {
  return (filter: ListFilterModel) => {
    const groupValue = { id: group.id, label: group.name };
    // if group is already present, then we modify it, otherwise add
    let groupCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "groups";
    }) as GroupsCriterion | undefined;

    if (groupCriterion) {
      // we should be showing group only. Remove other values
      groupCriterion.value.items = [groupValue];
      groupCriterion.modifier = GQL.CriterionModifier.Includes;
    } else {
      groupCriterion = new GroupsCriterion(GroupsCriterionOption);
      groupCriterion.value = {
        items: [groupValue],
        excluded: [],
        depth: showChildGroupContent ? -1 : 0,
      };
      groupCriterion.modifier = GQL.CriterionModifier.Includes;
      filter.criteria.push(groupCriterion);
    }

    return filter;
  };
};

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
