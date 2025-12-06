import React, { ReactNode, useCallback, useContext, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { languageData, getLanguageCode } from "src/utils/caption";
import {
  CaptionCriterion,
} from "src/models/list-filter/criteria/captions";
import { FacetCountsContext } from "src/hooks/useFacetCounts";

// Create language code badge element
function createCodeBadge(code: string): React.ReactNode {
  if (!code) return null;
  return (
    <span 
      style={{ 
        marginRight: "0.5em",
        fontSize: "0.7em",
        fontWeight: 600,
        padding: "0.15em 0.4em",
        borderRadius: "3px",
        backgroundColor: "rgba(255,255,255,0.15)",
        textTransform: "uppercase",
      }}
    >
      {code}
    </span>
  );
}

function useCaptionsFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  counts?: Map<string, number>;
}) {
  const intl = useIntl();
  const { option, filter, setFilter, counts } = props;

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

    // If a language is selected and it's included, add it
    if (criterion.value && modifier === CriterionModifier.Includes) {
      const code = getLanguageCode(criterion.value);
      modifierValues.push({
        id: criterion.value,
        label: criterion.value,
        icon: createCodeBadge(code),
      });
    }

    return modifierValues;
  }, [intl, selectedModifiers, criterion.value, modifier]);

  // Build excluded items list
  const excluded = useMemo(() => {
    // If a language is selected and it's excluded, add it
    // (the exclude icon is already shown by the list component)
    if (criterion.value && modifier === CriterionModifier.Excludes) {
      const code = getLanguageCode(criterion.value);
      return [{
        id: criterion.value,
        label: criterion.value,
        icon: createCodeBadge(code),
      }];
    }
    return [];
  }, [criterion.value, modifier]);

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
    // Also filter out zero-count languages when counts are available
    const hasLoadedCounts = counts && counts.size > 0;
    const filteredLanguages = languageData.filter((lang) => {
      if (criterion.value === lang.name) return false;
      if (!query) {
        // When no query, filter by count if available
        if (hasLoadedCounts) {
          const count = counts.get(lang.code.toLowerCase());
          return count !== undefined && count > 0;
        }
        return true;
      }
      const matchesQuery = lang.name.toLowerCase().includes(query.toLowerCase()) ||
             lang.code.toLowerCase().includes(query.toLowerCase());
      if (!matchesQuery) return false;
      // Even with query, filter by count if available
      if (hasLoadedCounts) {
        const count = counts.get(lang.code.toLowerCase());
        return count !== undefined && count > 0;
      }
      return true;
    });

    return modifierCandidates.concat(
      filteredLanguages.map((lang) => {
        const count = counts?.get(lang.code.toLowerCase());
        return {
          id: lang.name,
          label: lang.name,
          canExclude: true,
          icon: createCodeBadge(lang.code.toUpperCase()),
          count,
        };
      })
    );
  }, [modifier, criterion.value, query, intl, counts]);

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
  // Get facet counts from context - include loading state to avoid stale data filtering
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);

  const state = useCaptionsFilterState({
    filter,
    setFilter,
    option,
    // Pass empty map when loading to prevent filtering with stale data
    counts: facetsLoading ? new Map() : facetCounts.captions,
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

