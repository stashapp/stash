import React, { useCallback, useContext, useMemo } from "react";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMobileAlt, faDesktop, faSquare } from "@fortawesome/free-solid-svg-icons";
import {
  CriterionModifier,
} from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { OrientationCriterion } from "src/models/list-filter/criteria/orientation";
import { orientationStrings, stringToOrientation } from "src/utils/orientation";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { FacetCountsContext } from "src/hooks/useFacetCounts";

function createOrientationIcon(orientation: string): React.ReactNode {
  let icon;
  switch (orientation.toLowerCase()) {
    case "portrait":
      icon = faMobileAlt;
      break;
    case "landscape":
      icon = faDesktop;
      break;
    case "square":
      icon = faSquare;
      break;
    default:
      icon = faDesktop;
  }
  return (
    <FontAwesomeIcon
      icon={icon}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

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
  sectionID?: string;
}

export const SidebarOrientationFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();
  
  // Get facet counts from context - include loading state to avoid stale data filtering
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);

  const options = useMemo(() => {
    return orientationStrings.map((orientation) => {
      const orientationEnum = stringToOrientation(orientation);
      // Don't show counts when loading to prevent stale data filtering
      const count = (orientationEnum && !facetsLoading)
        ? facetCounts.orientations.get(orientationEnum)
        : undefined;
      return {
        id: orientation,
        label: orientation,
        icon: createOrientationIcon(orientation),
        count,
      };
    });
  }, [facetCounts, facetsLoading]);

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
        newCriterion.value = [...currentValues, item.id];
      }
    } else {
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

  // Filter out selected options and zero-count options
  const candidates = useMemo(() => {
    // Only filter by counts if facets have loaded AND are not currently loading
    const hasValidCounts = facetCounts.orientations.size > 0 && !facetsLoading;
    
    return options.filter((p) => {
      if (selected.find((s) => s.id === p.id)) return false;
      if (!hasValidCounts) return true;
      return p.count !== undefined && p.count > 0;
    });
  }, [options, selected, facetCounts.orientations.size, facetsLoading]);

  return (
    <SidebarListFilter
      title={title}
      candidates={candidates}
      onSelect={onSelect}
      onUnselect={onUnselect}
      selected={selected}
      singleValue={false}
      sectionID={sectionID}
    />
  );
};
