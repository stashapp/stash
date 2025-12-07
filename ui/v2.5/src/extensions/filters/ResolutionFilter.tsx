import React, { ReactNode, useCallback, useContext, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTv } from "@fortawesome/free-solid-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { resolutionStrings, stringToResolution } from "src/utils/resolution";
import { ResolutionCriterion } from "src/models/list-filter/criteria/resolution";
import { FacetCountsContext } from "src/extensions/hooks/useFacetCounts";

function createResolutionIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faTv}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function useResolutionFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  counts?: Map<string, number>;
}) {
  const intl = useIntl();
  const { option, filter, setFilter, counts } = props;

  const [query, setQuery] = useState("");

  // Resolution options with counts from facets
  const resolutionOptions = useMemo(() => {
    return resolutionStrings.map((res) => {
      const resEnum = stringToResolution(res);
      const count = resEnum ? counts?.get(resEnum) : undefined;
      return {
        id: res,
        label: res,
        count,
      };
    });
  }, [counts]);

  // Get or create criterion
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ResolutionCriterion;
    return filter.makeCriterion(option.type) as ResolutionCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ResolutionCriterion) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );
      if (c.isValid()) newCriteria.push(c);
      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier } = criterion;

  const getModifierLabel = useCallback(
    (mod: CriterionModifier) => {
      switch (mod) {
        case CriterionModifier.Equals:
          return intl.formatMessage({
            id: "criterion_modifier.equals",
            defaultMessage: "is",
          });
        case CriterionModifier.NotEquals:
          return intl.formatMessage({
            id: "criterion_modifier.not_equals",
            defaultMessage: "is not",
          });
        case CriterionModifier.GreaterThan:
          return intl.formatMessage({
            id: "criterion_modifier.greater_than",
            defaultMessage: "greater than",
          });
        case CriterionModifier.LessThan:
          return intl.formatMessage({
            id: "criterion_modifier.less_than",
            defaultMessage: "less than",
          });
        default:
          return "";
      }
    },
    [intl]
  );

  // Track pending resolution (selected but no modifier chosen yet)
  const [pendingResolution, setPendingResolution] = useState<string>("");

  // Build selected items list
  const selected = useMemo(() => {
    const selectedItems: Option[] = [];

    if (criterion.value) {
      selectedItems.push({
        id: "modifier",
        label: `(${getModifierLabel(modifier)})`,
        className: "modifier-object",
      });
      selectedItems.push({
        id: criterion.value,
        label: criterion.value,
        icon: createResolutionIcon(),
      });
    } else if (pendingResolution) {
      selectedItems.push({
        id: pendingResolution,
        label: pendingResolution,
        icon: createResolutionIcon(),
      });
    }

    return selectedItems;
  }, [criterion.value, modifier, getModifierLabel, pendingResolution]);

  // Build candidates list (resolution options or modifier options)
  const candidates = useMemo(() => {
    if (pendingResolution) {
      return [
        {
          id: "equals",
          label: `(${getModifierLabel(CriterionModifier.Equals)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "not_equals",
          label: `(${getModifierLabel(CriterionModifier.NotEquals)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "greater_than",
          label: `(${getModifierLabel(CriterionModifier.GreaterThan)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "less_than",
          label: `(${getModifierLabel(CriterionModifier.LessThan)})`,
          className: "modifier-object",
          canExclude: false,
        },
      ];
    }

    if (criterion.value) {
      return [];
    }

    // Filter out zero-count options
    const hasLoadedCounts = counts && counts.size > 0;
    const filteredResolutions = resolutionOptions.filter((res) => {
      const matchesQuery = !query || res.label.toLowerCase().includes(query.toLowerCase());
      if (!matchesQuery) return false;
      if (!hasLoadedCounts) return true;
      return res.count !== undefined && res.count > 0;
    });

    return filteredResolutions.map((res) => ({
      id: res.id,
      label: res.label,
      icon: createResolutionIcon(),
      canExclude: false,
      count: res.count,
    }));
  }, [criterion.value, query, getModifierLabel, pendingResolution, resolutionOptions]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      if (v.className === "modifier-object") {
        if (pendingResolution) {
          const newCriterion = criterion.clone() as ResolutionCriterion;
          newCriterion.value = pendingResolution;
          
          let mod = CriterionModifier.Equals;
          switch (v.id) {
            case "equals":
              mod = CriterionModifier.Equals;
              break;
            case "not_equals":
              mod = CriterionModifier.NotEquals;
              break;
            case "greater_than":
              mod = CriterionModifier.GreaterThan;
              break;
            case "less_than":
              mod = CriterionModifier.LessThan;
              break;
          }
          newCriterion.modifier = mod;
          setCriterion(newCriterion);
          setPendingResolution("");
        }
        return;
      }

      setPendingResolution(v.id);
    },
    [criterion, setCriterion, pendingResolution]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      if (v.className === "modifier-object") {
        const newCriterion = criterion.clone() as ResolutionCriterion;
        newCriterion.modifier = CriterionModifier.Equals;
        newCriterion.value = "";
        setCriterion(newCriterion);
        setPendingResolution("");
        return;
      }

      if (criterion.value) {
        const newCriterion = criterion.clone() as ResolutionCriterion;
        newCriterion.value = "";
        newCriterion.modifier = CriterionModifier.Equals;
        setCriterion(newCriterion);
      }
      setPendingResolution("");
    },
    [criterion, setCriterion]
  );

  return {
    candidates,
    onSelect,
    onUnselect,
    selected,
    excluded: [] as Option[],
    canExclude: false,
    query,
    setQuery,
  };
}

export const SidebarResolutionFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}> = ({ title, option, filter, setFilter, sectionID }) => {
  // Get facet counts from context - include loading state to avoid stale data filtering
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);
  
  const state = useResolutionFilterState({
    filter,
    setFilter,
    option,
    // Pass empty map when loading to prevent filtering with stale data
    counts: facetsLoading ? new Map() : facetCounts.resolutions,
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

export default SidebarResolutionFilter;
