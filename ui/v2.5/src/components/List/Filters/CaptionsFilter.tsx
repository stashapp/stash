import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { languageMap } from "src/utils/caption";
import {
  CaptionCriterion,
} from "src/models/list-filter/criteria/captions";

// Get language options from the languageMap
const languageOptions = Array.from(languageMap.entries()).map(
  ([code, name]) => ({
    code,
    name,
  })
);

function useCaptionsFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [query, setQuery] = useState("");

  // Get or create criterion
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as CaptionCriterion;

    return filter.makeCriterion(option.type) as CaptionCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: CaptionCriterion) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier } = criterion;

  // Build selected modifiers
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Build selected items list (for included languages)
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

    // If a language is selected and it's included (not excluded), add it with "(includes)" prefix
    if (
      criterion.value &&
      modifier === CriterionModifier.Includes
    ) {
      // Add the includes modifier indicator
      modifierValues.push({
        id: "includes_modifier",
        label: `(${intl.formatMessage({
          id: "criterion_modifier.includes",
          defaultMessage: "includes",
        })})`,
        className: "modifier-object",
      });
      // Add the language value
      modifierValues.push({
        id: criterion.value,
        label: criterion.value,
      });
    }

    return modifierValues;
  }, [intl, selectedModifiers, criterion.value, modifier]);

  // Build excluded items list
  const excluded = useMemo(() => {
    const excludedItems: Option[] = [];
    
    if (
      criterion.value &&
      modifier === CriterionModifier.Excludes
    ) {
      // Add the excludes modifier indicator
      excludedItems.push({
        id: "excludes_modifier",
        label: `(${intl.formatMessage({
          id: "criterion_modifier.excludes",
          defaultMessage: "excludes",
        })})`,
        className: "modifier-object",
      });
      // Add the language value
      excludedItems.push({
        id: criterion.value,
        label: criterion.value,
      });
    }
    return excludedItems;
  }, [criterion.value, modifier, intl]);

  // Build candidates list (language options)
  const candidates = useMemo(() => {
    const modifierCandidates: Option[] = [];

    // Show modifier options when no language selected
    if (
      (modifier === CriterionModifier.Includes ||
        modifier === CriterionModifier.Excludes) &&
      !criterion.value
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

    // Don't show language options if modifier is any/none
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return modifierCandidates;
    }

    // Filter languages by query and exclude already selected
    const filteredLanguages = languageOptions.filter((lang) => {
      if (criterion.value === lang.name) return false;
      if (!query) return true;
      return lang.name.toLowerCase().includes(query.toLowerCase());
    });

    return modifierCandidates.concat(
      filteredLanguages.map((lang) => ({
        id: lang.name,
        label: lang.name,
        canExclude: true,
      }))
    );
  }, [modifier, criterion.value, query, intl]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as CaptionCriterion;

      if (v.className === "modifier-object") {
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = "";
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = "";
        }
        setCriterion(newCriterion);
        return;
      }

      // Set the language value with appropriate modifier
      newCriterion.value = v.id;
      newCriterion.modifier = exclude
        ? CriterionModifier.Excludes
        : CriterionModifier.Includes;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as CaptionCriterion;

      if (v.className === "modifier-object") {
        // Reset to default modifier and clear value
        newCriterion.modifier = CriterionModifier.Includes;
        newCriterion.value = "";
        setCriterion(newCriterion);
        return;
      }

      // Clear the language value
      newCriterion.value = "";
      newCriterion.modifier = CriterionModifier.Includes;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  return {
    candidates,
    onSelect,
    onUnselect,
    selected,
    excluded,
    canExclude: true,
    query,
    setQuery,
  };
}

export const SidebarCaptionsFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}> = ({ title, option, filter, setFilter, sectionID }) => {
  const state = useCaptionsFilterState({
    filter,
    setFilter,
    option,
  });

  return (
    <SidebarListFilter
      {...state}
      title={title}
      sectionID={sectionID}
      singleValue
    />
  );
};

export default SidebarCaptionsFilter;

