import { useState, useCallback, useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";

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
