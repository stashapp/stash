import React, { useCallback, useMemo, useState } from "react";
import { Badge, Button, Collapse, Form } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { faCheckCircle, faMinus, faPlus, faTimes, faTimesCircle } from "@fortawesome/free-solid-svg-icons";
import { ClearableInput } from "src/components/Shared/ClearableInput";
import { CriterionType, ILabeledId } from "src/models/list-filter/types";
import { cloneDeep, debounce } from "lodash-es";
import {
  IHierarchicalLabeledIdCriterion,
  ILabeledIdCriterion,
} from "src/models/list-filter/criteria/criterion";
import { useIntl } from "react-intl";
import { faEye, faEyeSlash } from "@fortawesome/free-regular-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";

interface ISelectableFilter {
  query: string;
  setQuery: (query: string) => void;
  queryResults: ILabeledId[];
  selected: ILabeledId[];
  excluded: ILabeledId[];
  onSelect: (value: ILabeledId, include: boolean) => void;
  onUnselect: (value: ILabeledId) => void;
}

const SelectableFilter: React.FC<ISelectableFilter> = ({
  query,
  setQuery,
  queryResults,
  selected,
  excluded,
  onSelect,
  onUnselect,
}) => {
  const [internalQuery, setInternalQuery] = useState(query);

  const onInputChange = useMemo(() => {
    return debounce((input: string) => {
      setQuery(input);
    }, 250);
  }, [setQuery]);

  function onInternalInputChange(input: string) {
    setInternalQuery(input);
    onInputChange(input);
  }

  const objects = useMemo(() => {
    return queryResults.filter(
      (p) => selected.find((s) => s.id === p.id) === undefined && excluded.find((s) => s.id === p.id) === undefined
    );
  }, [queryResults, selected]);

  return (
    <div className="selectable-filter">
      <ClearableInput
        value={internalQuery}
        setValue={(v) => onInternalInputChange(v)}
      />
      <ul>
        {selected.map((p) => (
          <li key={p.id} className="selected-object">
            <a onClick={() => onUnselect(p)}>
              <div>
                <Icon className="fa-fw include-button" icon={faCheckCircle} />
                <span className="selected-object-label">{p.label}</span>
              </div>
              <div></div>
            </a>
          </li>
        ))}
        {excluded.map((p) => (
          <li key={p.id} className="excluded-object">
            <a onClick={() => onUnselect(p)}>
              <div>
                <Icon className="fa-fw exclude-icon" icon={faTimesCircle} />
                <span className="excluded-object-label">{p.label}</span>
              </div>
              <div></div>
            </a>
          </li>
        ))}
        {objects.map((p) => (
          <li key={p.id} className="unselected-object">
            <a onClick={() => onSelect(p, true)}>
              <div>
                <Icon className="fa-fw include-button" icon={faPlus} />
                <span>{p.label}</span>
              </div>
              <div>
                <span className="object-count">{p.id}</span>
                <Button onClick={() => onSelect(p, false)} className="minimal exclude-button">
                  <Icon className="fa-fw exclude-icon" icon={faMinus} />
                </Button>
              </div>
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
};

interface IHeader {
  title: string;
  selected: number;
  include: boolean | undefined;
  alwaysShown?: JSX.Element;
}

export const Header: React.FC<React.PropsWithChildren<IHeader>> = (
  props: React.PropsWithChildren<IHeader>
) => {
  const [open, setOpen] = useState(false);

  const icon = props.include ? faEye : faEyeSlash;

  return (
    <div>
      <Button onClick={() => setOpen(!open)} className="filter-header">
        <span className="header-title">{props.title}</span>
        {!!props.selected && (
          <span>
            <Icon icon={icon} />
            <Badge>{props.selected}</Badge>
          </span>
        )}
      </Button>
      {props.alwaysShown}
      <Collapse in={open}>
        <div>{props.children}</div>
      </Collapse>
    </div>
  );
};

interface IObjectsFilter<T extends ILabeledIdCriterion> {
  type: CriterionType;
  criterion: T;
  setCriterion: (criterion: T) => void;
  queryHook: (query: string) => ILabeledId[];
}

export const ObjectsFilter = <T extends ILabeledIdCriterion>(
  props: IObjectsFilter<T>
) => {
  const { type, criterion, setCriterion, queryHook } = props;

  const [excluded, setExcluded] = useState<ILabeledId[]>([]);

  const intl = useIntl();

  const [query, setQuery] = useState("");

  const queryResults = queryHook(query);

  function onSelect(value: ILabeledId, newInclude: boolean) {
    let newCriterion: T = cloneDeep(criterion);

    if (newInclude) {
      newCriterion.value.push(value);
      setCriterion(newCriterion);
    } else {
      setExcluded([...excluded, value]);
    }
  }

  const onUnselect = useCallback(
    (value: ILabeledId) => {
      if (!criterion) return;

      let newCriterion: T = cloneDeep(criterion);

      newCriterion.value = criterion.value.filter((v) => v.id !== value.id);

      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const sortedSelected = useMemo(() => {
    const ret = criterion.value.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  return (
    <SelectableFilter
      query={query}
      setQuery={setQuery}
      selected={sortedSelected}
      queryResults={queryResults}
      onSelect={onSelect}
      onUnselect={onUnselect}
      excluded={excluded}
    />
  );
};

interface IHierarchicalObjectsFilter<
  T extends IHierarchicalLabeledIdCriterion
> {
  type: CriterionType;
  criterion: T;
  setCriterion: (filter: T) => void;
  queryHook: (query: string) => ILabeledId[];
}

export const HierarchicalObjectsFilter = <
  T extends IHierarchicalLabeledIdCriterion
>(
  props: IHierarchicalObjectsFilter<T>
) => {
  const { type, criterion, setCriterion, queryHook } = props;

  const intl = useIntl();

  const [query, setQuery] = useState("");

  const queryResults = queryHook(query);

  // const criterion = useMemo(() => {
  //   return filter.criteria.find((cc) => {
  //     return cc.criterionOption.type === type;
  //   }) as T | undefined;
  // }, [filter, type]);

  const include = useMemo(() => {
    if (!criterion) return;

    switch (criterion.modifier) {
      case CriterionModifier.IncludesAll:
        return true;
      case CriterionModifier.Excludes:
        return false;
    }
  }, [criterion]);

  function onSelect(value: ILabeledId, newInclude: boolean) {
    let newCriterion: T = cloneDeep(criterion);

    // if (newInclude) {
    //   newCriterion.modifier = CriterionModifier.IncludesAll;
    // } else {
    //   newCriterion.modifier = CriterionModifier.Excludes;
    // }

    newCriterion.value.items.push(value);

    // const newFilter = filter.setCriterion(type, newCriterion);

    setCriterion(newCriterion);
  }

  const onUnselect = useCallback(
    (id: string) => {
      if (!criterion) return;

      let newCriterion: T | undefined = cloneDeep(criterion);

      newCriterion.value.items = criterion.value.items.filter(
        (v) => v.id !== id
      );
      // if (newCriterion.value.items.length === 0) {
      //   newCriterion = undefined;
      // }

      // const newFilter = filter.setCriterion(type, newCriterion);

      setCriterion(newCriterion);
    },
    [setCriterion, criterion, type]
  );

  const sortedSelected = useMemo(() => {
    const ret = criterion.value.items.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  const selected = useMemo(() => {
    if (!sortedSelected.length) return;

    return (
      <ul className="selected-objects">
        {sortedSelected.map((s) => {
          return (
            <li key={s.id}>
              <Badge
                className="selected-object"
                variant="secondary"
                onClick={() => onUnselect(s.id)}
              >
                <span>{s.label}</span>
                <Button variant="secondary">
                  <Icon icon={faTimes} />
                </Button>
              </Badge>
            </li>
          );
        })}
      </ul>
    );
  }, [sortedSelected, onUnselect]);

  return (
    <Header
      title={intl.formatMessage({ id: type })}
      include={include}
      selected={criterion?.value.items.length ?? 0}
      alwaysShown={selected}
    >
      <SelectableFilter
        query={query}
        setQuery={setQuery}
        selected={criterion?.value.items ?? []}
        queryResults={queryResults}
        onSelect={onSelect}
      />
    </Header>
  );
};
