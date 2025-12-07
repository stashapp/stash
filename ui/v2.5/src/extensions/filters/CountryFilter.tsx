import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { CountryCriterion } from "src/models/list-filter/criteria/country";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { getCountries } from "src/utils/country";

// Create a flag icon element using flag-icons CSS library
function createFlagIcon(countryCode: string): React.ReactNode {
  if (!countryCode || countryCode.length !== 2) return null;
  
  const code = countryCode.toLowerCase();
  return (
    <span
      className={`fi fi-${code}`}
      style={{ marginRight: "0.5em", fontSize: "1em" }}
    />
  );
}

function useCountryFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [query, setQuery] = useState("");

  // Get country options based on locale with flag icons
  const countryOptions = useMemo(() => {
    return getCountries(intl.locale).map((country) => ({
      id: country.value,
      label: country.label,
      icon: createFlagIcon(country.value),
    }));
  }, [intl.locale]);

  // Get or create criterion
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as CountryCriterion;

    return filter.makeCriterion(option.type) as CountryCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: CountryCriterion) => {
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

  // Build selected items list (for included countries)
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

    // If a country is selected and it's included (not excluded), add it
    if (
      criterion.value &&
      modifier === CriterionModifier.Equals
    ) {
      const country = countryOptions.find((c) => c.id === criterion.value);
      if (country) {
        modifierValues.push({
          id: criterion.value,
          label: country.label,
          icon: country.icon,
        });
      }
    }

    return modifierValues;
  }, [intl, selectedModifiers, criterion.value, modifier, countryOptions]);

  // Build excluded items list
  const excluded = useMemo(() => {
    if (
      criterion.value &&
      modifier === CriterionModifier.NotEquals
    ) {
      const country = countryOptions.find((c) => c.id === criterion.value);
      if (country) {
        return [{
          id: criterion.value,
          label: country.label,
          icon: country.icon,
        }];
      }
    }
    return [];
  }, [criterion.value, modifier, countryOptions]);

  // Build candidates list (country options)
  const candidates = useMemo(() => {
    const modifierCandidates: Option[] = [];

    // Show modifier options when no country selected
    if (
      (modifier === CriterionModifier.Equals ||
        modifier === CriterionModifier.NotEquals) &&
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

    // Don't show country options if modifier is any/none
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return modifierCandidates;
    }

    // Filter countries by query and exclude already selected
    const filteredCountries = countryOptions.filter((country) => {
      if (criterion.value === country.id) return false;
      if (!query) return true;
      return country.label.toLowerCase().includes(query.toLowerCase());
    });

    return modifierCandidates.concat(
      filteredCountries.map((country) => ({
        id: country.id,
        label: country.label,
        canExclude: true,
        icon: country.icon,
      }))
    );
  }, [modifier, criterion.value, query, intl, countryOptions]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as CountryCriterion;

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

      // Set the country value with appropriate modifier
      newCriterion.value = v.id;
      newCriterion.modifier = exclude
        ? CriterionModifier.NotEquals
        : CriterionModifier.Equals;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone() as CountryCriterion;

      if (v.className === "modifier-object") {
        // Reset to default modifier
        newCriterion.modifier = CriterionModifier.Equals;
        newCriterion.value = "";
        setCriterion(newCriterion);
        return;
      }

      // Clear the country value
      newCriterion.value = "";
      newCriterion.modifier = CriterionModifier.Equals;
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

export const SidebarCountryFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}> = ({ title, option, filter, setFilter, sectionID }) => {
  const state = useCountryFilterState({
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

export default SidebarCountryFilter;
