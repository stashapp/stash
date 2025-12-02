import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { CountryCriterion } from "src/models/list-filter/criteria/country";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { getCountries } from "src/utils/country";

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export const SidebarCountryFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();

  const options = useMemo(() => {
    return getCountries(intl.locale).map((country) => ({
      id: country.value,
      label: country.label,
    }));
  }, [intl.locale]);

  const criteria = filter.criteriaFor(option.type) as CountryCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.modifier === CriterionModifier.Equals) {
      return options.filter((o) => criterion.value === o.id);
    }

    return [];
  }, [options, criterion]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (criterion && criterion.modifier === CriterionModifier.Equals) {
      if (criterion.value === item.id) {
        // Remove selection if same item is selected
        setFilter(filter.removeCriterion(option.type));
        return;
      }
    }

    // Set new selection
    newCriterion.modifier = CriterionModifier.Equals;
    newCriterion.value = item.id;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (
      criterion &&
      criterion.modifier === CriterionModifier.Equals &&
      criterion.value === item.id
    ) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  // handle filtering of selected options
  const candidates = useMemo(() => {
    return options.filter(
      (p) => selected.find((s) => s.id === p.id) === undefined
    );
  }, [options, selected]);

  return (
    <>
      <SidebarListFilter
        title={title}
        candidates={candidates}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selected}
        singleValue={true}
      />
    </>
  );
};
