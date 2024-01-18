import React, {
  useState,
  useCallback,
  useMemo,
  useEffect,
  useContext,
} from "react";
import * as GQL from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useHistory, useLocation } from "react-router-dom";
import isEqual from "lodash-es/isEqual";
import { useConfigureUI } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";

export interface ICriterionOption {
  option: CriterionOption;
  showInSidebar: boolean;
}

export function useFilterConfig(mode: GQL.FilterMode) {
  const { configuration } = useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const ui = (configuration?.ui ?? {}) as IUIConfig;

  const savedOrder: string[] = useMemo(
    () => ui.criterionOrder?.[mode.toLowerCase()] ?? [],
    [mode, ui.criterionOrder]
  );

  const savedSidebar: string[] | undefined =
    ui.sidebarCriteria?.[mode.toLocaleLowerCase()];

  const defaultOptions = useMemo(() => {
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

  const [criterionOptions, setCriterionOptionsState] = useState(defaultOptions);

  useEffect(() => {
    const newOrder: ICriterionOption[] = [];
    savedOrder.forEach((o) => {
      const option = defaultOptions.find((d) => d.option.type === o);
      if (option) {
        newOrder.push({ ...option });
      }
    });

    // insert any missing options at the index they would be in the default order
    defaultOptions.forEach((o, i) => {
      if (!newOrder.some((n) => n.option.type === o.option.type)) {
        newOrder.splice(i, 0, { ...o });
      }
    });

    // override sidebar options
    if (savedSidebar) {
      newOrder.forEach((o) => {
        o.showInSidebar = savedSidebar.includes(o.option.type);
      });
    }

    setCriterionOptionsState(newOrder);
  }, [defaultOptions, savedOrder, savedSidebar]);

  function saveCriterionOptions(newOptions: ICriterionOption[]) {
    const criteriaOrder = newOptions.map((o) => o.option.type);
    const sidebarCriteria = newOptions
      .filter((o) => o.showInSidebar)
      .map((o) => o.option.type);

    saveUI({
      variables: {
        partial: {
          criterionOrder: {
            [mode.toLowerCase()]: criteriaOrder,
          },
          sidebarCriteria: {
            [mode.toLowerCase()]: sidebarCriteria,
          },
        },
      },
    });
  }

  function setCriterionOptions(newOptions: ICriterionOption[]) {
    setCriterionOptionsState(newOptions);
    saveCriterionOptions(newOptions);
  }

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
      let newFilter = prevFilter.empty();
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
