import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { OrientationCriterion } from "src/models/list-filter/criteria/orientation";
import { orientationStrings } from "src/utils/orientation";
import { Option, SidebarListFilter } from "./SidebarListFilter";

interface IOrientationFilterProps {
  criterion: OrientationCriterion;
  onValueChanged: (value: string[]) => void;
}

export const OrientationFilter: React.FC<IOrientationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  // This is for the main filter dialog - not implemented yet
  return null;
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export const SidebarOrientationFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();

  const options = useMemo(() => {
    return orientationStrings.map((orientation) => ({
      id: orientation,
      label: orientation,
    }));
  }, [intl]);

  const criteria = filter.criteriaFor(option.type) as OrientationCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (
      criterion.modifier === CriterionModifier.Includes ||
      criterion.modifier === CriterionModifier.Excludes
    ) {
      return options.filter((option) => criterion.value.includes(option.id));
    }

    return [];
  }, [options, criterion]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (criterion && criterion.modifier === CriterionModifier.Includes) {
      const currentValues = criterion.value;
      if (!currentValues.includes(item.id)) {
        // Add to selection
        newCriterion.value = [...currentValues, item.id];
      }
    } else {
      // Start new selection
      newCriterion.modifier = CriterionModifier.Includes;
      newCriterion.value = [item.id];
    }
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (
      criterion &&
      criterion.modifier === CriterionModifier.Includes &&
      criterion.value.includes(item.id)
    ) {
      const newValues = criterion.value.filter((v) => v !== item.id);
      if (newValues.length === 0) {
        setFilter(filter.removeCriterion(option.type));
        return;
      }
      newCriterion.value = newValues;
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
        singleValue={false}
      />
    </>
  );
};
