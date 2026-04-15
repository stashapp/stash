import { gql } from "@apollo/client";
import * as GQL from "src/core/generated-graphql";
import { getClient } from "src/core/StashService";
import {
  TagsCriterion,
  TagsCriterionOption,
} from "src/models/list-filter/criteria/tags";
import { ListFilterModel } from "src/models/list-filter/filter";

export const useTagFilterHook = (
  tag: GQL.TagDataFragment,
  showSubTagContent?: boolean
) => {
  return (filter: ListFilterModel) => {
    const tagValue = { id: tag.id, label: tag.name };
    // if tag is already present, then we modify it, otherwise add
    let tagCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "tags";
    }) as TagsCriterion | undefined;

    if (tagCriterion) {
      if (
        tagCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        tagCriterion.modifier === GQL.CriterionModifier.Includes
      ) {
        // add the tag if not present
        if (
          !tagCriterion.value.items.find((p) => {
            return p.id === tag.id;
          })
        ) {
          tagCriterion.value.items.push(tagValue);
        }
      } else {
        // overwrite
        tagCriterion.value.items = [tagValue];
      }

      tagCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      tagCriterion = new TagsCriterion(TagsCriterionOption);
      tagCriterion.value = {
        items: [tagValue],
        excluded: [],
        depth: showSubTagContent ? -1 : 0,
      };
      tagCriterion.modifier = GQL.CriterionModifier.IncludesAll;
      filter.criteria.push(tagCriterion);
    }

    return filter;
  };
};

interface ITagRelationTuple {
  parents: GQL.SlimTagDataFragment[];
  children: GQL.SlimTagDataFragment[];
}

export const tagRelationHook = (
  tag: GQL.SlimTagDataFragment | GQL.TagDataFragment | GQL.TagListDataFragment,
  old: ITagRelationTuple,
  updated: ITagRelationTuple
) => {
  const { cache } = getClient();

  const tagRef = cache.writeFragment({
    data: tag,
    fragment: gql`
      fragment Tag on Tag {
        id
      }
    `,
  });

  function updater(
    property: "parents" | "children",
    oldTags: GQL.SlimTagDataFragment[],
    updatedTags: GQL.SlimTagDataFragment[]
  ) {
    oldTags.forEach((o) => {
      if (!updatedTags.some((u) => u.id === o.id)) {
        cache.modify({
          id: cache.identify(o),
          fields: {
            [property](value, { readField }) {
              return (value as GQL.SlimTagDataFragment[]).filter(
                (t) => readField("id", t) !== tag.id
              );
            },
          },
        });
      }
    });

    updatedTags.forEach((u) => {
      if (!oldTags.some((o) => o.id === u.id)) {
        cache.modify({
          id: cache.identify(u),
          fields: {
            [property](value) {
              return [...(value as unknown[]), tagRef];
            },
          },
        });
      }
    });
  }

  updater("children", old.parents, updated.parents);
  updater("parents", old.children, updated.children);
};
