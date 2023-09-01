import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  faCheckCircle,
  faMinus,
  faPlus,
  faTimesCircle,
} from "@fortawesome/free-solid-svg-icons";
import { faTimesCircle as faTimesCircleRegular } from "@fortawesome/free-regular-svg-icons";
import { ClearableInput } from "src/components/Shared/ClearableInput";
import {
  IHierarchicalLabelValue,
  ILabeledId,
  ILabeledValueListValue,
} from "src/models/list-filter/types";
import { cloneDeep } from "lodash-es";
import {
  Criterion,
  IHierarchicalLabeledIdCriterion,
} from "src/models/list-filter/criteria/criterion";
import { defineMessages, MessageDescriptor, useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { keyboardClickHandler } from "src/utils/keyboard";
import { useDebouncedSetState } from "src/hooks/debounce";
import useFocus from "src/utils/focus";

interface ISelectedItem {
  item: ILabeledId;
  excluded?: boolean;
  onClick: () => void;
}

const SelectedItem: React.FC<ISelectedItem> = ({
  item,
  excluded = false,
  onClick,
}) => {
  const iconClassName = excluded ? "exclude-icon" : "include-button";
  const spanClassName = excluded
    ? "excluded-object-label"
    : "selected-object-label";
  const [hovered, setHovered] = useState(false);

  const icon = useMemo(() => {
    if (!hovered) {
      return excluded ? faTimesCircle : faCheckCircle;
    }

    return faTimesCircleRegular;
  }, [hovered, excluded]);

  function onMouseOver() {
    setHovered(true);
  }

  function onMouseOut() {
    setHovered(false);
  }

  return (
    <a
      onClick={() => onClick()}
      onKeyDown={keyboardClickHandler(onClick)}
      onMouseEnter={() => onMouseOver()}
      onMouseLeave={() => onMouseOut()}
      onFocus={() => onMouseOver()}
      onBlur={() => onMouseOut()}
      tabIndex={0}
    >
      <div>
        <Icon className={`fa-fw ${iconClassName}`} icon={icon} />
        <span className={spanClassName}>{item.label}</span>
      </div>
      <div></div>
    </a>
  );
};

interface ISelectableFilter {
  query: string;
  onQueryChange: (query: string) => void;
  modifier: CriterionModifier;
  inputFocus: ReturnType<typeof useFocus>;
  canExclude: boolean;
  queryResults: ILabeledId[];
  selected: ILabeledId[];
  excluded: ILabeledId[];
  onSelect: (value: ILabeledId, exclude: boolean) => void;
  onUnselect: (value: ILabeledId) => void;
}

const SelectableFilter: React.FC<ISelectableFilter> = ({
  query,
  onQueryChange,
  modifier,
  inputFocus,
  canExclude,
  queryResults,
  selected,
  excluded,
  onSelect,
  onUnselect,
}) => {
  const objects = useMemo(() => {
    return queryResults.filter(
      (p) =>
        selected.find((s) => s.id === p.id) === undefined &&
        excluded.find((s) => s.id === p.id) === undefined
    );
  }, [queryResults, selected, excluded]);

  const includingOnly = modifier == CriterionModifier.Equals;
  const excludingOnly =
    modifier == CriterionModifier.Excludes ||
    modifier == CriterionModifier.NotEquals;

  const includeIcon = <Icon className="fa-fw include-button" icon={faPlus} />;
  const excludeIcon = <Icon className="fa-fw exclude-icon" icon={faMinus} />;

  return (
    <div className="selectable-filter">
      <ClearableInput
        focus={inputFocus}
        value={query}
        setValue={(v) => onQueryChange(v)}
      />
      <ul>
        {selected.map((p) => (
          <li key={p.id} className="selected-object">
            <SelectedItem
              item={p}
              excluded={excludingOnly}
              onClick={() => onUnselect(p)}
            />
          </li>
        ))}
        {excluded.map((p) => (
          <li key={p.id} className="excluded-object">
            <SelectedItem item={p} excluded onClick={() => onUnselect(p)} />
          </li>
        ))}
        {objects.map((p) => (
          <li key={p.id} className="unselected-object">
            <a
              onClick={() => onSelect(p, false)}
              onKeyDown={keyboardClickHandler(() => onSelect(p, false))}
              tabIndex={0}
            >
              <div>
                {!excludingOnly ? includeIcon : excludeIcon}
                <span>{p.label}</span>
              </div>
              <div>
                {/* TODO item count */}
                {/* <span className="object-count">{p.id}</span> */}
                {canExclude && !includingOnly && !excludingOnly && (
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      onSelect(p, true);
                    }}
                    onKeyDown={(e) => e.stopPropagation()}
                    className="minimal exclude-button"
                  >
                    <span className="exclude-button-text">exclude</span>
                    {excludeIcon}
                  </Button>
                )}
              </div>
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
};

interface IObjectsFilter<T extends Criterion<ILabeledValueListValue>> {
  criterion: T;
  setCriterion: (criterion: T) => void;
  useResults: (query: string) => { results: ILabeledId[]; loading: boolean };
}

export const ObjectsFilter = <
  T extends Criterion<ILabeledValueListValue | IHierarchicalLabelValue>
>({
  criterion,
  setCriterion,
  useResults,
}: IObjectsFilter<T>) => {
  const [query, setQuery] = useState("");
  const [displayQuery, setDisplayQuery] = useState(query);

  const debouncedSetQuery = useDebouncedSetState(setQuery, 250);
  const onQueryChange = useCallback(
    (input: string) => {
      setDisplayQuery(input);
      debouncedSetQuery(input);
    },
    [debouncedSetQuery, setDisplayQuery]
  );

  const [queryResults, setQueryResults] = useState<ILabeledId[]>([]);
  const { results, loading: resultsLoading } = useResults(query);
  useEffect(() => {
    if (!resultsLoading) {
      setQueryResults(results);
    }
  }, [results, resultsLoading]);

  const inputFocus = useFocus();
  const [, setInputFocus] = inputFocus;

  function onSelect(value: ILabeledId, newExclude: boolean) {
    let newCriterion: T = cloneDeep(criterion);

    if (newExclude) {
      if (newCriterion.value.excluded) {
        newCriterion.value.excluded.push(value);
      } else {
        newCriterion.value.excluded = [value];
      }
    } else {
      newCriterion.value.items.push(value);
    }

    setCriterion(newCriterion);

    // reset filter query after selecting
    debouncedSetQuery.cancel();
    setQuery("");
    setDisplayQuery("");

    // focus the input box
    setInputFocus();
  }

  const onUnselect = useCallback(
    (value: ILabeledId) => {
      if (!criterion) return;

      let newCriterion: T = cloneDeep(criterion);

      newCriterion.value.items = criterion.value.items.filter(
        (v) => v.id !== value.id
      );
      newCriterion.value.excluded = criterion.value.excluded.filter(
        (v) => v.id !== value.id
      );

      setCriterion(newCriterion);

      // focus the input box
      setInputFocus();
    },
    [criterion, setCriterion, setInputFocus]
  );

  const sortedSelected = useMemo(() => {
    const ret = criterion.value.items.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  const sortedExcluded = useMemo(() => {
    if (!criterion.value.excluded) return [];
    const ret = criterion.value.excluded.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  // if excludes is not a valid modifierOption then we can use `value.excluded`
  const canExclude =
    criterion.criterionOption.modifierOptions.find(
      (m) => m === CriterionModifier.Excludes
    ) === undefined;

  return (
    <SelectableFilter
      query={displayQuery}
      onQueryChange={onQueryChange}
      modifier={criterion.modifier}
      inputFocus={inputFocus}
      canExclude={canExclude}
      selected={sortedSelected}
      queryResults={queryResults}
      onSelect={onSelect}
      onUnselect={onUnselect}
      excluded={sortedExcluded}
    />
  );
};

interface IHierarchicalObjectsFilter<T extends IHierarchicalLabeledIdCriterion>
  extends IObjectsFilter<T> {}

export const HierarchicalObjectsFilter = <
  T extends IHierarchicalLabeledIdCriterion
>(
  props: IHierarchicalObjectsFilter<T>
) => {
  const intl = useIntl();
  const { criterion, setCriterion } = props;

  const messages = defineMessages({
    studio_depth: {
      id: "studio_depth",
      defaultMessage: "Levels (empty for all)",
    },
  });

  function onDepthChanged(depth: number) {
    let newCriterion: T = cloneDeep(criterion);
    newCriterion.value.depth = depth;
    setCriterion(newCriterion);
  }

  function criterionOptionTypeToIncludeID(): string {
    if (criterion.criterionOption.type === "studios") {
      return "include-sub-studios";
    }
    if (criterion.criterionOption.type === "children") {
      return "include-parent-tags";
    }
    return "include-sub-tags";
  }

  function criterionOptionTypeToIncludeUIString(): MessageDescriptor {
    const optionType =
      criterion.criterionOption.type === "studios"
        ? "include_sub_studios"
        : criterion.criterionOption.type === "children"
        ? "include_parent_tags"
        : "include_sub_tags";
    return {
      id: optionType,
    };
  }

  return (
    <Form>
      {criterion.modifier !== CriterionModifier.Equals && (
        <Form.Group>
          <Form.Check
            id={criterionOptionTypeToIncludeID()}
            checked={criterion.value.depth !== 0}
            label={intl.formatMessage(criterionOptionTypeToIncludeUIString())}
            onChange={() =>
              onDepthChanged(criterion.value.depth !== 0 ? 0 : -1)
            }
          />
        </Form.Group>
      )}

      {criterion.value.depth !== 0 && (
        <Form.Group>
          <Form.Control
            className="btn-secondary"
            type="number"
            placeholder={intl.formatMessage(messages.studio_depth)}
            onChange={(e) =>
              onDepthChanged(e.target.value ? parseInt(e.target.value, 10) : -1)
            }
            defaultValue={
              criterion.value && criterion.value.depth !== -1
                ? criterion.value.depth
                : ""
            }
            min="1"
          />
        </Form.Group>
      )}
      <ObjectsFilter {...props} />
    </Form>
  );
};
