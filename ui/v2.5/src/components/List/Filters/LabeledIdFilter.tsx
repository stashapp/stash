import React, { useCallback, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { FilterSelect, SelectObject } from "src/components/Shared/Select";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import { ILoadResults, useCacheResults } from "src/hooks/data";
import {
  CriterionOption,
  ModifierCriterion,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  IHierarchicalLabelValue,
  ILabeledId,
  ILabeledValueListValue,
} from "src/models/list-filter/types";
import { Option } from "./SidebarListFilter";
import {
  CriterionModifier,
  FilterMode,
  InputMaybe,
  IntCriterionInput,
  SceneFilterType,
} from "src/core/generated-graphql";
import { useIntl } from "react-intl";

interface ILabeledIdFilterProps {
  criterion: ModifierCriterion<ILabeledId[]>;
  onValueChanged: (value: ILabeledId[]) => void;
}

export const LabeledIdFilter: React.FC<ILabeledIdFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const criterionOption = criterion.modifierCriterionOption();
  const { inputType } = criterionOption;

  if (
    inputType !== "performers" &&
    inputType !== "studios" &&
    inputType !== "scene_tags" &&
    inputType !== "performer_tags" &&
    inputType !== "tags" &&
    inputType !== "scenes" &&
    inputType !== "groups" &&
    inputType !== "galleries"
  ) {
    return null;
  }

  function getLabel(i: SelectObject) {
    switch (inputType) {
      case "galleries":
        return galleryTitle(i);
      case "scenes":
        return objectTitle(i);
    }

    return i.name ?? i.title ?? "";
  }

  function onSelectionChanged(items: SelectObject[]) {
    onValueChanged(
      items.map((i) => ({
        id: i.id,
        label: getLabel(i),
      }))
    );
  }

  return (
    <Form.Group>
      <FilterSelect
        type={inputType}
        isMulti
        onSelect={onSelectionChanged}
        ids={criterion.value.map((labeled) => labeled.id)}
        menuPortalTarget={document.body}
      />
    </Form.Group>
  );
};

type ModifierValue = "any" | "none" | "any_of" | "only" | "include_subs";

export function getModifierCandidates(props: {
  modifier: CriterionModifier;
  defaultModifier: CriterionModifier;
  hasSelected?: boolean;
  hasExcluded?: boolean;
  singleValue?: boolean;
  hierarchical?: boolean;
}) {
  const {
    modifier,
    defaultModifier,
    hasSelected,
    hasExcluded,
    singleValue,
    hierarchical,
  } = props;
  const ret: ModifierValue[] = [];

  if (modifier === defaultModifier && !hasSelected && !hasExcluded) {
    ret.push("any");
  }
  if (modifier === defaultModifier && !hasSelected && !hasExcluded) {
    ret.push("none");
  }
  if (!singleValue && modifier === defaultModifier && hasSelected) {
    ret.push("any_of");
  }
  if (
    hierarchical &&
    modifier === defaultModifier &&
    (hasSelected || hasExcluded)
  ) {
    ret.push("include_subs");
  }
  if (
    !singleValue &&
    modifier === defaultModifier &&
    hasSelected &&
    !hasExcluded
  ) {
    ret.push("only");
  }
  return ret;
}

export function modifierValueToModifier(key: ModifierValue): CriterionModifier {
  switch (key) {
    case "any":
      return CriterionModifier.NotNull;
    case "none":
      return CriterionModifier.IsNull;
    case "any_of":
      return CriterionModifier.Includes;
    case "only":
      return CriterionModifier.Equals;
  }

  throw new Error("Invalid modifier value");
}

function getDefaultModifier(singleValue: boolean) {
  if (singleValue) {
    return CriterionModifier.Includes;
  }
  return CriterionModifier.IncludesAll;
}

export function useSelectionState(props: {
  criterion: ModifierCriterion<ILabeledValueListValue>;
  setCriterion: (c: ModifierCriterion<ILabeledValueListValue>) => void;
  singleValue?: boolean;
  hierarchical?: boolean;
  includeSubMessageID?: string;
}) {
  const intl = useIntl();

  const {
    criterion,
    setCriterion,
    singleValue = false,
    hierarchical = false,
    includeSubMessageID,
  } = props;
  const { modifier } = criterion;

  const defaultModifier = getDefaultModifier(singleValue);

  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
      any_of: !singleValue && modifier === CriterionModifier.Includes,
      only: !singleValue && modifier === CriterionModifier.Equals,
      include_subs:
        hierarchical &&
        modifier === defaultModifier &&
        (criterion.value as IHierarchicalLabelValue).depth === -1,
    };
  }, [modifier, singleValue, criterion.value, defaultModifier, hierarchical]);

  const selected = useMemo(() => {
    const modifierValues: Option[] = Object.entries(selectedModifiers)
      .filter((v) => v[1])
      .map((v) => {
        const messageID =
          v[0] === "include_subs"
            ? includeSubMessageID
            : `criterion_modifier_values.${v[0]}`;

        return {
          id: v[0],
          label: `(${intl.formatMessage({
            id: messageID,
          })})`,
          className: "modifier-object",
        };
      });

    return modifierValues.concat(
      criterion.value.items.map((s) => ({
        id: s.id,
        label: s.label,
      }))
    );
  }, [intl, selectedModifiers, criterion.value.items, includeSubMessageID]);

  const excluded = useMemo(() => {
    return criterion.value.excluded.map((s) => ({
      id: s.id,
      label: s.label,
    }));
  }, [criterion.value.excluded]);

  const includingOnly = modifier == CriterionModifier.Equals;
  const excludingOnly =
    modifier == CriterionModifier.Excludes ||
    modifier == CriterionModifier.NotEquals;

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion: ModifierCriterion<ILabeledValueListValue> =
        criterion.clone();

      if (v.className === "modifier-object") {
        if (v.id === "include_subs") {
          (newCriterion.value as IHierarchicalLabelValue).depth = -1;
          setCriterion(newCriterion);
          return;
        }

        newCriterion.modifier = modifierValueToModifier(v.id as ModifierValue);
        setCriterion(newCriterion);
        return;
      }

      // if only exclude is allowed, then add to excluded
      if (excludingOnly) {
        exclude = true;
      }

      const items = !exclude ? criterion.value.items : criterion.value.excluded;
      const newItems = [...items, v];

      if (!exclude) {
        newCriterion.value.items = newItems;
      } else {
        newCriterion.value.excluded = newItems;
      }
      setCriterion(newCriterion);
    },
    [excludingOnly, criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, exclude: boolean) => {
      const newCriterion = criterion.clone();

      if (v.className === "modifier-object") {
        if (v.id === "include_subs") {
          newCriterion.value.depth = 0;
          setCriterion(newCriterion);
          return;
        }
        newCriterion.modifier = defaultModifier;
        setCriterion(newCriterion);
        return;
      }

      const items = !exclude ? criterion.value.items : criterion.value.excluded;
      const newItems = items.filter((i) => i.id !== v.id);

      if (!exclude) {
        newCriterion.value.items = newItems;
      } else {
        newCriterion.value.excluded = newItems;
      }
      setCriterion(newCriterion);
    },
    [criterion, setCriterion, defaultModifier]
  );

  return { selected, excluded, onSelect, onUnselect, includingOnly };
}

export function useCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<ILabeledValueListValue>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<ILabeledValueListValue>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<ILabeledValueListValue>) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  return { criterion, setCriterion };
}

export interface IUseQueryHookProps {
  q: string;
  filter?: ListFilterModel;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  skip: boolean;
}

export function useQueryState(
  useQuery: (props: IUseQueryHookProps) => ILoadResults<ILabeledId[]>,
  filter: ListFilterModel,
  skip: boolean,
  options?: {
    filterHook?: (filter: ListFilterModel) => ListFilterModel;
  }
) {
  const [query, setQuery] = useState("");
  const { results: queryResults } = useCacheResults(
    useQuery({ q: query, filter, filterHook: options?.filterHook, skip })
  );

  return { query, setQuery, queryResults };
}

export function useCandidates(props: {
  criterion: ModifierCriterion<ILabeledValueListValue>;
  queryResults: ILabeledId[] | undefined;
  selected: Option[];
  excluded: Option[];
  hierarchical?: boolean;
  singleValue?: boolean;
  includeSubMessageID?: string;
}) {
  const intl = useIntl();

  const {
    criterion,
    queryResults,
    selected,
    excluded,
    hierarchical = false,
    singleValue = false,
    includeSubMessageID,
  } = props;
  const { modifier } = criterion;

  const results = useMemo(() => {
    if (
      !queryResults ||
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return [];
    }

    return queryResults.filter(
      (p) =>
        selected.find((s) => s.id === p.id) === undefined &&
        excluded.find((s) => s.id === p.id) === undefined
    );
  }, [queryResults, modifier, selected, excluded]);

  const defaultModifier = getDefaultModifier(singleValue);

  const candidates = useMemo(() => {
    const hierarchicalCandidate =
      hierarchical && (criterion.value as IHierarchicalLabelValue).depth !== -1;

    const modifierCandidates: Option[] = getModifierCandidates({
      modifier,
      defaultModifier,
      hasSelected: selected.length > 0,
      hasExcluded: excluded.length > 0,
      singleValue,
      hierarchical: hierarchicalCandidate,
    }).map((v) => {
      const messageID =
        v === "include_subs"
          ? includeSubMessageID
          : `criterion_modifier_values.${v}`;

      return {
        id: v,
        label: `(${intl.formatMessage({
          id: messageID,
        })})`,
        className: "modifier-object",
        canExclude: false,
      };
    });

    return modifierCandidates.concat(
      (results ?? []).map((r) => ({
        id: r.id,
        label: r.label,
      }))
    );
  }, [
    defaultModifier,
    intl,
    modifier,
    singleValue,
    results,
    selected,
    excluded,
    criterion.value,
    hierarchical,
    includeSubMessageID,
  ]);

  return candidates;
}

export function useLabeledIdFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  useQuery: (props: IUseQueryHookProps) => ILoadResults<ILabeledId[]>;
  singleValue?: boolean;
  hierarchical?: boolean;
  includeSubMessageID?: string;
}) {
  const {
    option,
    filter,
    setFilter,
    filterHook,
    useQuery,
    singleValue = false,
    hierarchical = false,
    includeSubMessageID,
  } = props;

  // defer querying until the user opens the filter
  const [skip, setSkip] = useState(true);

  const { query, setQuery, queryResults } = useQueryState(
    useQuery,
    filter,
    skip,
    { filterHook }
  );

  const { criterion, setCriterion } = useCriterion(option, filter, setFilter);

  const { selected, excluded, onSelect, onUnselect, includingOnly } =
    useSelectionState({
      criterion,
      setCriterion,
      singleValue,
      hierarchical,
      includeSubMessageID,
    });

  const candidates = useCandidates({
    criterion,
    queryResults,
    selected,
    excluded,
    hierarchical,
    singleValue,
    includeSubMessageID,
  });

  const onOpen = useCallback(() => {
    setSkip(false);
  }, []);

  return {
    candidates,
    onSelect,
    onUnselect,
    selected,
    excluded,
    canExclude: !includingOnly,
    query,
    setQuery,
    onOpen,
  };
}

export function makeQueryVariables(query: string, extraProps: {}) {
  return {
    filter: {
      q: query,
      per_page: 200,
    },
    ...extraProps,
  };
}

interface IFilterType {
  scenes_filter?: InputMaybe<SceneFilterType>;
  scene_count?: InputMaybe<IntCriterionInput>;
}

export function setObjectFilter(
  out: IFilterType,
  mode: FilterMode,
  relatedFilterOutput: SceneFilterType
) {
  const empty = Object.keys(relatedFilterOutput).length === 0;

  switch (mode) {
    case FilterMode.Scenes:
      // if empty, only get objects with scenes
      if (empty) {
        out.scene_count = {
          modifier: CriterionModifier.GreaterThan,
          value: 0,
        };
      }
      out.scenes_filter = relatedFilterOutput;
      break;
  }
}
