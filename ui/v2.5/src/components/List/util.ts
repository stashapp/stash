import React, { useState, useCallback, useMemo, useEffect } from "react";
import * as GQL from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useHistory, useLocation } from "react-router-dom";
import isEqual from "lodash-es/isEqual";

export interface ICriterionOption {
  option: CriterionOption;
  showInSidebar: boolean;
}

export function useFilterConfig(mode: GQL.FilterMode) {
  // TODO - save this in the UI config

  const getCriterionOptions = useCallback(() => {
    const options = getFilterOptions(mode);

    return options.criterionOptions.map((o) => {
      return {
        option: o,
        showInSidebar: !options.defaultHiddenOptions.some(
          (c) => c.type === o.type
        ),
      } as ICriterionOption;
    });
  }, [mode]);

  const [criterionOptions, setCriterionOptions] = useState(
    getCriterionOptions()
  );

  const sidebarOptions = useMemo(
    () => criterionOptions.filter((o) => o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );
  const hiddenOptions = useMemo(
    () => criterionOptions.filter((o) => !o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );

  return {
    criterionOptions,
    sidebarOptions,
    hiddenOptions,
    setCriterionOptions,
  };
}

export function useFilterURL(
  filter: ListFilterModel,
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>,
  defaultFilter: ListFilterModel
) {
  const history = useHistory();
  const location = useLocation();

  // this hook causes the initial render to update the URL, losing
  // the existing URL params.
  // useEffect(() => {
  //   const newParams = filter.makeQueryParameters();
  //   history.replace({ ...history.location, search: newParams });
  // }, [filter, history]);

  // This hook runs on every page location change (ie navigation),
  // and updates the filter accordingly.
  useEffect(() => {
    // re-init to load default filter on empty new query params
    if (!location.search) {
      setFilter(defaultFilter.clone());
      return;
    }

    // the query has changed, update filter if necessary
    setFilter((prevFilter) => {
      let newFilter = prevFilter.clone();
      newFilter.configureFromQueryString(location.search);
      if (!isEqual(newFilter, prevFilter)) {
        return newFilter;
      } else {
        return prevFilter;
      }
    });
  }, [location.search, defaultFilter, setFilter]);

  // when the filter changes, update the URL
  const updateFilter = useCallback(
    (newFilter: ListFilterModel) => {
      const newParams = newFilter.makeQueryParameters();
      history.replace({ ...history.location, search: newParams });
    },
    [history]
  );

  return { setFilter: updateFilter };
}

// returns true if the filter has changed in a way that impacts the total count
function totalCountImpacted(
  oldFilter: ListFilterModel,
  newFilter: ListFilterModel
) {
  return (
    oldFilter.criteria.length !== newFilter.criteria.length ||
    oldFilter.criteria.some((c) => {
      const newCriterion = newFilter.criteria.find(
        (nc) => nc.getId() === c.getId()
      );
      return !newCriterion || !isEqual(c, newCriterion);
    })
  );
}

// this hook caches the total count of results, and only updates it when the filter changes
export function useResultCount(
  filter: ListFilterModel,
  loading: boolean,
  count: number
) {
  const [resultCount, setResultCount] = useState(count);
  const [lastFilter, setLastFilter] = useState(filter);

  // if we are only changing the page or sort, don't update the result count
  useEffect(() => {
    if (!loading) {
      setResultCount(count);
    } else {
      if (totalCountImpacted(lastFilter, filter)) {
        setResultCount(count);
      }
    }

    setLastFilter(filter);
  }, [loading, filter, count, lastFilter]);

  return resultCount;
}
