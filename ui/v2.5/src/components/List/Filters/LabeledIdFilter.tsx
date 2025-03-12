import React, { ReactNode, useCallback, useMemo, useState } from "react";
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
  ILabeledId,
  ILabeledValueListValue,
} from "src/models/list-filter/types";
import { SidebarListFilter } from "./SidebarListFilter";

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

export const LabeledIdQuickFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  useQuery: (q: string) => ILoadResults<ILabeledId[]>;
  singleValue?: boolean;
}> = ({ title, option, filter, setFilter, useQuery, singleValue }) => {
  const [query, setQuery] = useState("");

  const { results } = useCacheResults(useQuery(query));

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

  const onSelect = useCallback(
    (v: ILabeledId, exclude: boolean) => {
      const items = !exclude ? criterion.value.items : criterion.value.excluded;
      const newItems = [...items, v];
      const newCriterion = criterion.clone();
      if (!exclude) {
        newCriterion.value.items = newItems;
      } else {
        newCriterion.value.excluded = newItems;
      }
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: ILabeledId, exclude: boolean) => {
      const items = !exclude ? criterion.value.items : criterion.value.excluded;
      const newItems = items.filter((i) => i.id !== v.id);
      const newCriterion = criterion.clone();
      if (!exclude) {
        newCriterion.value.items = newItems;
      } else {
        newCriterion.value.excluded = newItems;
      }
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const selected = useMemo(() => {
    return criterion.value.items.map((s) => ({
      id: s.id,
      label: s.label,
    }));
  }, [criterion.value.items]);

  const excluded = useMemo(() => {
    return criterion.value.excluded.map((s) => ({
      id: s.id,
      label: s.label,
    }));
  }, [criterion.value.excluded]);

  const candidates = useMemo(() => {
    return (results ?? [])
      .filter((r) => {
        return (
          !selected.some((s) => s.id === r.id) &&
          !excluded.some((s) => s.id === r.id)
        );
      })
      .map((r) => ({
        id: r.id,
        label: r.label,
      }));
  }, [results, selected, excluded]);

  // FIXME - implement modifier optoins
  return (
    <SidebarListFilter
      title={title}
      candidates={candidates}
      onSelect={onSelect}
      onUnselect={onUnselect}
      selected={selected}
      excluded={excluded}
      canExclude
      query={query}
      setQuery={setQuery}
      singleValue={singleValue}
    />
  );
};
