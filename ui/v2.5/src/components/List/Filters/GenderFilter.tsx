import React, { ReactNode, useCallback, useMemo } from "react";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faVenus,
  faTransgenderAlt,
  faMars,
} from "@fortawesome/free-solid-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { GenderCriterion } from "src/models/list-filter/criteria/gender";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

// Create a gender icon element
function createGenderIcon(genderId: string): React.ReactNode {
  let icon;
  let color;

  switch (genderId) {
    case "Male":
      icon = faMars;
      color = "#6fa8dc"; // blue
      break;
    case "Female":
      icon = faVenus;
      color = "#e06666"; // pink/red
      break;
    case "Transgender Male":
    case "Transgender Female":
    case "Intersex":
    case "Non-Binary":
      icon = faTransgenderAlt;
      color = "#9966cc"; // purple
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
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  // Gender options with icons
  const genderOptions = useMemo(() => {
    return [
      {
        id: "Male",
        label: intl.formatMessage({ id: "gender_types.MALE" }),
        icon: createGenderIcon("Male"),
      },
      {
        id: "Female",
        label: intl.formatMessage({ id: "gender_types.FEMALE" }),
        icon: createGenderIcon("Female"),
      },
      {
        id: "Transgender Male",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_MALE" }),
        icon: createGenderIcon("Transgender Male"),
      },
      {
        id: "Transgender Female",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_FEMALE" }),
        icon: createGenderIcon("Transgender Female"),
      },
      {
        id: "Intersex",
        label: intl.formatMessage({ id: "gender_types.INTERSEX" }),
        icon: createGenderIcon("Intersex"),
      },
      {
        id: "Non-Binary",
        label: intl.formatMessage({ id: "gender_types.NON_BINARY" }),
        icon: createGenderIcon("Non-Binary"),
      },
    ];
  }, [intl]);

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

    // If genders are selected with Includes modifier, add them
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
  const excluded = useMemo(() => {
    if (modifier === CriterionModifier.Excludes && value.length > 0) {
      return value
        .map((genderId) => {
          const gender = genderOptions.find((g) => g.id === genderId);
          if (gender) {
            return {
              id: genderId,
              label: gender.label,
              icon: gender.icon,
            };
          }
          return null;
        })
        .filter((g): g is Option => g !== null);
    }
    return [];
  }, [modifier, value, genderOptions]);

  // Build candidates list
  const candidates = useMemo(() => {
    const modifierCandidates: Option[] = [];

    // Show modifier options when no specific genders are selected
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

    // Don't show gender options if modifier is any/none
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return modifierCandidates;
    }

    // Filter genders to exclude already selected ones
    const filteredGenders = genderOptions.filter((gender) => {
      return !value.includes(gender.id);
    });

    return modifierCandidates.concat(
      filteredGenders.map((gender) => ({
        id: gender.id,
        label: gender.label,
        canExclude: true,
        icon: gender.icon,
      }))
    );
  }, [modifier, value, intl, genderOptions]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as GenderCriterion;

      if (v.className === "modifier-object") {
        // Handle modifier selection
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = [];
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = [];
        }
      } else {
        // Handle gender selection
        if (exclude) {
          // If currently including, switch to excluding
          if (newCriterion.modifier === CriterionModifier.Includes) {
            newCriterion.modifier = CriterionModifier.Excludes;
            newCriterion.value = [v.id];
          } else if (newCriterion.modifier === CriterionModifier.Excludes) {
            // Add to excluded list
            if (!newCriterion.value.includes(v.id)) {
              newCriterion.value = [...newCriterion.value, v.id];
            }
          } else {
            // Start new exclude
            newCriterion.modifier = CriterionModifier.Excludes;
            newCriterion.value = [v.id];
          }
        } else {
          // Include
          if (newCriterion.modifier === CriterionModifier.Excludes) {
            // Switch from exclude to include
            newCriterion.modifier = CriterionModifier.Includes;
            newCriterion.value = [v.id];
          } else if (newCriterion.modifier === CriterionModifier.Includes) {
            // Add to included list
            if (!newCriterion.value.includes(v.id)) {
              newCriterion.value = [...newCriterion.value, v.id];
            }
          } else {
            // Start new include
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
        // Clear modifier
        setCriterion(null);
        return;
      }

      const newCriterion = criterion.clone() as GenderCriterion;

      // Remove from current selection
      newCriterion.value = newCriterion.value.filter((id) => id !== v.id);

      if (newCriterion.value.length === 0) {
        // Reset to default
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
  const state = useGenderFilterState({ option, filter, setFilter });

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
