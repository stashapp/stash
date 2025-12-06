import React, { ReactNode, useCallback, useContext, useMemo } from "react";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faVenus,
  faTransgenderAlt,
  faMars,
} from "@fortawesome/free-solid-svg-icons";
import {
  CriterionModifier,
  GenderEnum,
} from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { GenderCriterion } from "src/models/list-filter/criteria/gender";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { FacetCountsContext } from "src/hooks/useFacetCounts";

// Map string IDs to GenderEnum values
const genderIdToEnum: Record<string, GenderEnum> = {
  "Male": GenderEnum.Male,
  "Female": GenderEnum.Female,
  "Transgender Male": GenderEnum.TransgenderMale,
  "Transgender Female": GenderEnum.TransgenderFemale,
  "Intersex": GenderEnum.Intersex,
  "Non-Binary": GenderEnum.NonBinary,
};

// Create a gender icon element
function createGenderIcon(genderId: string): React.ReactNode {
  let icon;
  let color;

  switch (genderId) {
    case "Male":
      icon = faMars;
      color = "#6fa8dc";
      break;
    case "Female":
      icon = faVenus;
      color = "#e06666";
      break;
    case "Transgender Male":
    case "Transgender Female":
    case "Intersex":
    case "Non-Binary":
      icon = faTransgenderAlt;
      color = "#9966cc";
      break;
    default:
      return null;
  }

  return (
    <FontAwesomeIcon
      icon={icon}
      style={{ marginRight: "0.5em", color }}
      fixedWidth
    />
  );
}

function useGenderFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  counts?: Map<GenderEnum, number>;
}) {
  const intl = useIntl();
  const { option, filter, setFilter, counts } = props;

  // Gender options with icons and counts
  const genderOptions = useMemo(() => {
    return [
      {
        id: "Male",
        label: intl.formatMessage({ id: "gender_types.MALE" }),
        icon: createGenderIcon("Male"),
        count: counts?.get(GenderEnum.Male),
      },
      {
        id: "Female",
        label: intl.formatMessage({ id: "gender_types.FEMALE" }),
        icon: createGenderIcon("Female"),
        count: counts?.get(GenderEnum.Female),
      },
      {
        id: "Transgender Male",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_MALE" }),
        icon: createGenderIcon("Transgender Male"),
        count: counts?.get(GenderEnum.TransgenderMale),
      },
      {
        id: "Transgender Female",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_FEMALE" }),
        icon: createGenderIcon("Transgender Female"),
        count: counts?.get(GenderEnum.TransgenderFemale),
      },
      {
        id: "Intersex",
        label: intl.formatMessage({ id: "gender_types.INTERSEX" }),
        icon: createGenderIcon("Intersex"),
        count: counts?.get(GenderEnum.Intersex),
      },
      {
        id: "Non-Binary",
        label: intl.formatMessage({ id: "gender_types.NON_BINARY" }),
        icon: createGenderIcon("Non-Binary"),
        count: counts?.get(GenderEnum.NonBinary),
      },
    ];
  }, [intl, counts]);

  const criteria = filter.criteriaFor(option.type) as GenderCriterion[];
  const criterion = useMemo(() => {
    return criteria.length > 0 ? criteria[0] : (option.makeCriterion() as GenderCriterion);
  }, [criteria, option]);

  const setCriterion = useCallback(
    (c: GenderCriterion | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier, value } = criterion;

  // Build selected modifiers (any/none)
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Build selected items list (for included genders)
  const selected = useMemo(() => {
    const modifierValues: Option[] = Object.entries(selectedModifiers)
      .filter((v) => v[1])
      .map((v) => ({
        id: v[0],
        label: `(${intl.formatMessage({
          id: `criterion_modifier_values.${v[0]}`,
        })})`,
        className: "modifier-object",
      }));

    if (modifier === CriterionModifier.Includes && value.length > 0) {
      value.forEach((genderId) => {
        const gender = genderOptions.find((g) => g.id === genderId);
        if (gender) {
          modifierValues.push({
            id: genderId,
            label: gender.label,
            icon: gender.icon,
          });
        }
      });
    }

    return modifierValues;
  }, [intl, selectedModifiers, modifier, value, genderOptions]);

  // Build excluded items list
  const excluded: Option[] = useMemo(() => {
    if (modifier === CriterionModifier.Excludes && value.length > 0) {
      const result: Option[] = [];
      value.forEach((genderId) => {
        const gender = genderOptions.find((g) => g.id === genderId);
        if (gender) {
          result.push({
            id: genderId,
            label: gender.label,
            icon: gender.icon,
          });
        }
      });
      return result;
    }
    return [];
  }, [modifier, value, genderOptions]);

  // Build candidates list
  const candidates = useMemo(() => {
    const modifierCandidates: Option[] = [];

    if (
      (modifier === CriterionModifier.Includes ||
        modifier === CriterionModifier.Excludes) &&
      value.length === 0
    ) {
      modifierCandidates.push({
        id: "any",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.any",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
      modifierCandidates.push({
        id: "none",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.none",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
    }

    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return modifierCandidates;
    }

    // Filter genders to exclude already selected ones and zero-count options
    const hasLoadedCounts = counts && counts.size > 0;
    const filteredGenders = genderOptions.filter((gender) => {
      if (value.includes(gender.id)) return false;
      if (!hasLoadedCounts) return true; // No counts loaded yet, show all
      return gender.count !== undefined && gender.count > 0;
    });

    return modifierCandidates.concat(
      filteredGenders.map((gender) => ({
        id: gender.id,
        label: gender.label,
        canExclude: true,
        icon: gender.icon,
        count: gender.count,
      }))
    );
  }, [modifier, value, intl, genderOptions]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as GenderCriterion;

      if (v.className === "modifier-object") {
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = [];
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = [];
        }
      } else {
        if (exclude) {
          if (newCriterion.modifier === CriterionModifier.Includes) {
            newCriterion.modifier = CriterionModifier.Excludes;
            newCriterion.value = [v.id];
          } else if (newCriterion.modifier === CriterionModifier.Excludes) {
            if (!newCriterion.value.includes(v.id)) {
              newCriterion.value = [...newCriterion.value, v.id];
            }
          } else {
            newCriterion.modifier = CriterionModifier.Excludes;
            newCriterion.value = [v.id];
          }
        } else {
          if (newCriterion.modifier === CriterionModifier.Excludes) {
            newCriterion.modifier = CriterionModifier.Includes;
            newCriterion.value = [v.id];
          } else if (newCriterion.modifier === CriterionModifier.Includes) {
            if (!newCriterion.value.includes(v.id)) {
              newCriterion.value = [...newCriterion.value, v.id];
            }
          } else {
            newCriterion.modifier = CriterionModifier.Includes;
            newCriterion.value = [v.id];
          }
        }
      }

      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      if (v.className === "modifier-object") {
        setCriterion(null);
        return;
      }

      const newCriterion = criterion.clone() as GenderCriterion;
      newCriterion.value = newCriterion.value.filter((id) => id !== v.id);

      if (newCriterion.value.length === 0) {
        setCriterion(null);
      } else {
        setCriterion(newCriterion);
      }
    },
    [criterion, setCriterion]
  );

  return {
    selected,
    excluded,
    candidates,
    onSelect,
    onUnselect,
    canExclude: true,
  };
}

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarGenderFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  // Get facet counts from context - include loading state to avoid stale data filtering
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);
  
  const state = useGenderFilterState({ 
    option, 
    filter, 
    setFilter, 
    // Pass undefined counts when loading to prevent filtering with stale data
    counts: facetsLoading ? new Map() : facetCounts.genders 
  });

  return (
    <SidebarListFilter
      title={title}
      candidates={state.candidates}
      onSelect={state.onSelect}
      onUnselect={state.onUnselect}
      selected={state.selected}
      excluded={state.excluded}
      canExclude={state.canExclude}
      singleValue={false}
      sectionID={sectionID}
    />
  );
};
