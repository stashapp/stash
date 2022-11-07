import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import * as GQL from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";

export const usePerformerFilterHook = (
  performer: GQL.PerformerDataFragment
) => {
  return (filter: ListFilterModel) => {
    const performerValue = {
      id: performer.id,
      label: performer.name ?? `Performer ${performer.id}`,
    };
    // if performers is already present, then we modify it, otherwise add
    let performerCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "performers";
    }) as PerformersCriterion;

    if (
      performerCriterion &&
      (performerCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        performerCriterion.modifier === GQL.CriterionModifier.Includes)
    ) {
      // add the performer if not present
      if (
        !performerCriterion.value.find((p) => {
          return p.id === performer.id;
        })
      ) {
        performerCriterion.value.push(performerValue);
      }

      performerCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      // overwrite
      performerCriterion = new PerformersCriterion();
      performerCriterion.value = [performerValue];
      filter.criteria.push(performerCriterion);
    }

    return filter;
  };
};

interface IPerformerFragment {
  name?: GQL.Maybe<string>;
  gender?: GQL.Maybe<GQL.GenderEnum>;
}

export function sortPerformers<T extends IPerformerFragment>(performers: T[]) {
  const ret = performers.slice();
  ret.sort((a, b) => {
    if (a.gender === b.gender) {
      // sort by name
      return (a.name ?? "").localeCompare(b.name ?? "");
    }

    // TODO - may want to customise gender order
    const genderOrder = [
      GQL.GenderEnum.Female,
      GQL.GenderEnum.TransgenderFemale,
      GQL.GenderEnum.Male,
      GQL.GenderEnum.TransgenderMale,
      GQL.GenderEnum.Intersex,
      GQL.GenderEnum.NonBinary,
    ];

    const aIndex = a.gender
      ? genderOrder.indexOf(a.gender)
      : genderOrder.length;
    const bIndex = b.gender
      ? genderOrder.indexOf(b.gender)
      : genderOrder.length;
    return aIndex - bIndex;
  });

  return ret;
}
