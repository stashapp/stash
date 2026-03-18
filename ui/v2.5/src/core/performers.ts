import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import * as GQL from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { stringToGender } from "src/utils/gender";
import { filterData } from "src/utils/data";
import { stringToCircumcised } from "src/utils/circumcised";

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
    }) as PerformersCriterion | undefined;

    if (performerCriterion) {
      if (
        performerCriterion.modifier === GQL.CriterionModifier.IncludesAll ||
        performerCriterion.modifier === GQL.CriterionModifier.Includes
      ) {
        // add the performer if not present
        if (
          !performerCriterion.value.items.find((p) => {
            return p.id === performer.id;
          })
        ) {
          performerCriterion.value.items.push(performerValue);
        }
      } else {
        // overwrite
        performerCriterion.value.items = [performerValue];
      }

      performerCriterion.modifier = GQL.CriterionModifier.IncludesAll;
    } else {
      performerCriterion = new PerformersCriterion();
      performerCriterion.value.items = [performerValue];
      performerCriterion.modifier = GQL.CriterionModifier.IncludesAll;
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

export function normalizeAliases(
  aliases: GQL.PerformerAliasInput[]
): GQL.PerformerAliasInput[] {
  const aliasesMap = new Map<string, boolean>();
  aliases.forEach((a) => {
    const existing = aliasesMap.get(a.alias);
    // If duplicates exist, and their ignore_auto_tag differs, we default to true (ignore auto tag)
    const safeIgnore =
      existing === undefined
        ? a.ignore_auto_tag
        : existing || a.ignore_auto_tag;
    aliasesMap.set(a.alias, safeIgnore);
  });

  return [...aliasesMap].map(([alias, ignore_auto_tag]) => ({
    alias,
    ignore_auto_tag,
  }));
}

export function mergeScrapedAliases(
  scrapedAliases: string | null | undefined,
  existingAliases: GQL.PerformerAliasInput[] = []
): GQL.PerformerAliasInput[] {
  if (!scrapedAliases) {
    return existingAliases;
  }

  const existingMap = new Map<string, boolean>();
  existingAliases.forEach((a) => {
    existingMap.set(a.alias, a.ignore_auto_tag);
  });

  return scrapedAliases
    .split(",")
    .map((a) => {
      const trimmed = a.trim();
      const existing = existingMap.get(trimmed);
      return {
        alias: trimmed,
        ignore_auto_tag: existing !== undefined ? existing : true,
      };
    })
    .filter((a) => a.alias.length > 0);
}

export const scrapedPerformerToCreateInput = (
  toCreate: GQL.ScrapedPerformer,
  endpoint?: string
) => {
  const input: GQL.PerformerCreateInput = {
    name: toCreate.name ?? "",
    gender: stringToGender(toCreate.gender),
    birthdate: toCreate.birthdate,
    disambiguation: toCreate.disambiguation,
    ethnicity: toCreate.ethnicity,
    country: toCreate.country,
    eye_color: toCreate.eye_color,
    height_cm: toCreate.height ? Number(toCreate.height) : undefined,
    measurements: toCreate.measurements,
    fake_tits: toCreate.fake_tits,
    career_start: toCreate.career_start,
    career_end: toCreate.career_end,
    tattoos: toCreate.tattoos,
    piercings: toCreate.piercings,
    aliases: mergeScrapedAliases(toCreate.aliases),
    urls: toCreate.urls,
    tag_ids: filterData((toCreate.tags ?? []).map((t) => t.stored_id)),
    image:
      (toCreate.images ?? []).length > 0
        ? (toCreate.images ?? [])[0]
        : undefined,
    details: toCreate.details,
    death_date: toCreate.death_date,
    hair_color: toCreate.hair_color,
    weight: toCreate.weight ? Number(toCreate.weight) : undefined,
    penis_length: toCreate.penis_length
      ? Number(toCreate.penis_length)
      : undefined,
    circumcised: stringToCircumcised(toCreate.circumcised),
  };

  if (endpoint && toCreate.remote_site_id) {
    input.stash_ids = [
      {
        endpoint,
        stash_id: toCreate.remote_site_id,
      },
    ];
  }

  return input;
};
