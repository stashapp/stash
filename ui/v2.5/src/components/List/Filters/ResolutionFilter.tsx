import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { resolutionStrings } from "src/utils/resolution";
import { ResolutionCriterion } from "src/models/list-filter/criteria/resolution";

// Resolution options from the resolution utility
const resolutionOptions = resolutionStrings.map((res) => ({
  id: res,
  label: res,
}));

function useResolutionFilterState(props: {
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

  // Get modifier label for display
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

    // If filter is active (has value in criterion), show modifier and value
    if (criterion.value) {
      // Add the modifier indicator
      selectedItems.push({
        id: "modifier",
        label: `(${getModifierLabel(modifier)})`,
        className: "modifier-object",
      });
      // Add the resolution value
      selectedItems.push({
        id: criterion.value,
        label: criterion.value,
      });
    }
    // If there's a pending resolution (selected but no modifier yet), show it
    else if (pendingResolution) {
      selectedItems.push({
        id: pendingResolution,
        label: pendingResolution,
      });
    }

    return selectedItems;
  }, [criterion.value, modifier, getModifierLabel, pendingResolution]);

  // Build candidates list (resolution options or modifier options)
  const candidates = useMemo(() => {
    // If a resolution is pending (selected but no modifier yet), show modifier options
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

    // If filter is already active, don't show any candidates
    if (criterion.value) {
      return [];
    }

    // Show resolution options (no modifier options until resolution is selected)
    const filteredResolutions = resolutionOptions.filter((res) => {
      if (!query) return true;
      return res.label.toLowerCase().includes(query.toLowerCase());
    });

    return filteredResolutions.map((res) => ({
      id: res.id,
      label: res.label,
      canExclude: false,
    }));
  }, [criterion.value, query, getModifierLabel, pendingResolution]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      if (v.className === "modifier-object") {
        // User selected a modifier after choosing a resolution
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

      // User selected a resolution - store it as pending and wait for modifier
      setPendingResolution(v.id);
    },
    [criterion, setCriterion, pendingResolution]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Clicking on modifier clears the filter
        const newCriterion = criterion.clone() as ResolutionCriterion;
        newCriterion.modifier = CriterionModifier.Equals;
        newCriterion.value = "";
        setCriterion(newCriterion);
        setPendingResolution("");
        return;
      }

      // Clear the resolution value or pending resolution
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
  const state = useResolutionFilterState({
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

export default SidebarResolutionFilter;

