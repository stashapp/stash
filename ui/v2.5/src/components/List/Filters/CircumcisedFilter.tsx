import React, { useCallback, useContext, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle } from "@fortawesome/free-solid-svg-icons";
import { faCircle as faCircleRegular } from "@fortawesome/free-regular-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { CircumcisedCriterion } from "src/models/list-filter/criteria/circumcised";
import { circumcisedStrings, stringToCircumcised } from "src/utils/circumcised";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { FacetCountsContext } from "src/hooks/useFacetCounts";

function createCircumcisedIcon(value: string): React.ReactNode {
  // Use filled circle for "cut" and outline circle for "uncut"
  const icon = value.toLowerCase() === "cut" ? faCircle : faCircleRegular;
  return (
    <FontAwesomeIcon
      icon={icon}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarCircumcisedFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();
  
  // Get facet counts from context
  const { counts: facetCounts } = useContext(FacetCountsContext);

  // Build options from circumcised strings with counts
  const options: Option[] = useMemo(() => {
    return circumcisedStrings.map((value) => {
      const circumcisedEnum = stringToCircumcised(value);
      const count = circumcisedEnum ? facetCounts.circumcised.get(circumcisedEnum) : undefined;
      return {
        id: value,
        label: intl.formatMessage({ id: `circumcised_types.${value.toUpperCase()}` }),
        icon: createCircumcisedIcon(value),
        count,
      };
    });
  }, [intl, facetCounts.circumcised]);

  const criteria = filter.criteriaFor(option.type) as CircumcisedCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  // Build selected list from criterion
  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.modifier === CriterionModifier.IsNull) {
      return [{ id: "null", label: intl.formatMessage({ id: "criterion.none" }), className: "modifier-object" }];
    }
    if (criterion.modifier === CriterionModifier.NotNull) {
      return [{ id: "not-null", label: intl.formatMessage({ id: "criterion.any" }), className: "modifier-object" }];
    }

    if (
      criterion.modifier === CriterionModifier.Includes ||
      criterion.modifier === CriterionModifier.Excludes
    ) {
      return criterion.value.map((v) => {
        const opt = options.find((o) => o.id === v);
        return opt || { id: v, label: v };
      });
    }

    return [];
  }, [criterion, options, intl]);

  // Build excluded list
  const excluded: Option[] = useMemo(() => {
    if (!criterion || criterion.modifier !== CriterionModifier.Excludes) return [];
    
    return criterion.value.map((v) => {
      const opt = options.find((o) => o.id === v);
      return opt || { id: v, label: v };
    });
  }, [criterion, options]);

  // Build candidates (available options not yet selected)
  const candidates: Option[] = useMemo(() => {
    const modifierOptions: Option[] = [
      { id: "any", label: `(${intl.formatMessage({ id: "criterion.any" })})`, className: "modifier-object" },
      { id: "none", label: `(${intl.formatMessage({ id: "criterion.none" })})`, className: "modifier-object" },
    ];

    // If modifier is set, don't show value options
    if (criterion?.modifier === CriterionModifier.IsNull || 
        criterion?.modifier === CriterionModifier.NotNull) {
      return [];
    }

    // Filter out already selected options
    const selectedIds = new Set(selected.map((s) => s.id));
    const excludedIds = new Set(excluded.map((e) => e.id));
    
    const valueOptions = options
      .filter((o) => !selectedIds.has(o.id) && !excludedIds.has(o.id))
      .map((o) => ({ ...o, canExclude: true }));

    // Only show modifier options if nothing is selected yet
    if (selected.length === 0 && excluded.length === 0) {
      return [...modifierOptions, ...valueOptions];
    }

    return valueOptions;
  }, [criterion, options, selected, excluded, intl]);

  function onSelect(item: Option, exclude?: boolean) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    // Handle modifier options
    if (item.id === "any") {
      newCriterion.modifier = CriterionModifier.NotNull;
      newCriterion.value = [];
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
      return;
    }
    if (item.id === "none") {
      newCriterion.modifier = CriterionModifier.IsNull;
      newCriterion.value = [];
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
      return;
    }

    // Handle value selection
    if (exclude) {
      newCriterion.modifier = CriterionModifier.Excludes;
      if (!newCriterion.value.includes(item.id)) {
        newCriterion.value = [...newCriterion.value, item.id];
      }
    } else {
      newCriterion.modifier = CriterionModifier.Includes;
      if (!newCriterion.value.includes(item.id)) {
        newCriterion.value = [...newCriterion.value, item.id];
      }
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect(item: Option) {
    if (!criterion) return;

    // Handle modifier unselect
    if (item.id === "null" || item.id === "not-null" || item.id === "any" || item.id === "none") {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    // Handle value unselect
    const newValues = criterion.value.filter((v) => v !== item.id);
    if (newValues.length === 0) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    const newCriterion = criterion.clone();
    newCriterion.value = newValues;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  return (
    <SidebarListFilter
      title={title}
      candidates={candidates}
      selected={selected}
      excluded={excluded}
      onSelect={onSelect}
      onUnselect={onUnselect}
      singleValue={false}
      sectionID={sectionID}
    />
  );
};

