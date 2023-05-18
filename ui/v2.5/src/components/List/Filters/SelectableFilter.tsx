import React, { useCallback, useMemo, useState } from "react";
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
import { cloneDeep, debounce } from "lodash-es";
import {
  Criterion,
  IHierarchicalLabeledIdCriterion,
} from "src/models/list-filter/criteria/criterion";
import { defineMessages, MessageDescriptor, useIntl } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";

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
      onMouseEnter={() => onMouseOver()}
      onMouseLeave={() => onMouseOut()}
      onFocus={() => onMouseOver()}
      onBlur={() => onMouseOut()}
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
  setQuery: (query: string) => void;
  single: boolean;
  includeOnly: boolean;
  queryResults: ILabeledId[];
  selected: ILabeledId[];
  excluded: ILabeledId[];
  onSelect: (value: ILabeledId, include: boolean) => void;
  onUnselect: (value: ILabeledId) => void;
}

const SelectableFilter: React.FC<ISelectableFilter> = ({
  query,
  setQuery,
  single,
  queryResults,
  selected,
  excluded,
  includeOnly,
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
      (p) =>
        selected.find((s) => s.id === p.id) === undefined &&
        excluded.find((s) => s.id === p.id) === undefined
    );
  }, [queryResults, selected, excluded]);

  const includingOnly = includeOnly || (selected.length > 0 && single);
  const excludingOnly = excluded.length > 0 && single;

  const includeIcon = <Icon className="fa-fw include-button" icon={faPlus} />;
  const excludeIcon = <Icon className="fa-fw exclude-icon" icon={faMinus} />;

  return (
    <div className="selectable-filter">
      <ClearableInput
        value={internalQuery}
        setValue={(v) => onInternalInputChange(v)}
      />
      <ul>
        {selected.map((p) => (
          <li key={p.id} className="selected-object">
            <SelectedItem item={p} onClick={() => onUnselect(p)} />
          </li>
        ))}
        {excluded.map((p) => (
          <li key={p.id} className="excluded-object">
            <SelectedItem item={p} excluded onClick={() => onUnselect(p)} />
          </li>
        ))}
        {objects.map((p) => (
          <li key={p.id} className="unselected-object">
            {/* if excluding only, clicking on an item also excludes it */}
            <a onClick={() => onSelect(p, !excludingOnly)}>
              <div>
                {!excludingOnly ? includeIcon : excludeIcon}
                <span>{p.label}</span>
              </div>
              <div>
                {/* TODO item count */}
                {/* <span className="object-count">{p.id}</span> */}
                {!includingOnly && !excludingOnly && (
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      onSelect(p, false);
                    }}
                    className="minimal exclude-button"
                  >
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
  single?: boolean;
  setCriterion: (criterion: T) => void;
  queryHook: (query: string) => ILabeledId[];
}

export const ObjectsFilter = <
  T extends Criterion<ILabeledValueListValue | IHierarchicalLabelValue>
>(
  props: IObjectsFilter<T>
) => {
  const { criterion, setCriterion, queryHook, single = false } = props;

  const [query, setQuery] = useState("");

  const queryResults = queryHook(query);

  function onSelect(value: ILabeledId, newInclude: boolean) {
    let newCriterion: T = cloneDeep(criterion);

    if (newInclude) {
      newCriterion.value.items.push(value);
    } else {
      newCriterion.value.excluded.push(value);
    }

    setCriterion(newCriterion);
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
    },
    [criterion, setCriterion]
  );

  const sortedSelected = useMemo(() => {
    const ret = criterion.value.items.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  const sortedExcluded = useMemo(() => {
    const ret = criterion.value.excluded.slice();
    ret.sort((a, b) => a.label.localeCompare(b.label));
    return ret;
  }, [criterion]);

  return (
    <SelectableFilter
      single={single}
      includeOnly={criterion.modifier === CriterionModifier.Equals}
      query={query}
      setQuery={setQuery}
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
    if (criterion.criterionOption.type === "childTags") {
      return "include-parent-tags";
    }
    return "include-sub-tags";
  }

  function criterionOptionTypeToIncludeUIString(): MessageDescriptor {
    const optionType =
      criterion.criterionOption.type === "studios"
        ? "include_sub_studios"
        : criterion.criterionOption.type === "childTags"
        ? "include_parent_tags"
        : "include_sub_tags";
    return {
      id: optionType,
    };
  }

  return (
    <Form>
      <Form.Group>
        <Form.Check
          id={criterionOptionTypeToIncludeID()}
          checked={criterion.value.depth !== 0}
          label={intl.formatMessage(criterionOptionTypeToIncludeUIString())}
          onChange={() => onDepthChanged(criterion.value.depth !== 0 ? 0 : -1)}
        />
      </Form.Group>

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
