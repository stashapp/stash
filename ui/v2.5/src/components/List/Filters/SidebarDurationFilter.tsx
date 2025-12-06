import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faClock } from "@fortawesome/free-solid-svg-icons";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { NumberCriterion } from "src/models/list-filter/criteria/criterion";
import { CriterionModifier } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { INumberValue } from "src/models/list-filter/types";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import TextUtils from "src/utils/text";

// Duration range preset type
interface DurationRange {
  id: string;
  labelId: string;
  defaultLabel: string;
  minSeconds?: number;
  maxSeconds?: number;
}

// Predefined duration ranges (in seconds)
const DURATION_RANGES: DurationRange[] = [
  { id: "under-5", labelId: "duration_ranges.under_5_min", defaultLabel: "Under 5 min", minSeconds: undefined, maxSeconds: 5 * 60 },
  { id: "5-15", labelId: "duration_ranges.5_15_min", defaultLabel: "5-15 min", minSeconds: 5 * 60, maxSeconds: 15 * 60 },
  { id: "15-30", labelId: "duration_ranges.15_30_min", defaultLabel: "15-30 min", minSeconds: 15 * 60, maxSeconds: 30 * 60 },
  { id: "30-60", labelId: "duration_ranges.30_60_min", defaultLabel: "30-60 min", minSeconds: 30 * 60, maxSeconds: 60 * 60 },
  { id: "1-2h", labelId: "duration_ranges.1_2_hours", defaultLabel: "1-2 hours", minSeconds: 60 * 60, maxSeconds: 2 * 60 * 60 },
  { id: "over-2h", labelId: "duration_ranges.over_2_hours", defaultLabel: "Over 2 hours", minSeconds: 2 * 60 * 60, maxSeconds: undefined },
];

// Check if current criterion matches a duration range
function matchesDurationRange(
  criterion: NumberCriterion,
  range: DurationRange
): boolean {
  const { modifier, value } = criterion;

  // Range with both min and max (Between)
  if (range.minSeconds !== undefined && range.maxSeconds !== undefined) {
    if (modifier === CriterionModifier.Between) {
      return value.value === range.minSeconds && value.value2 === range.maxSeconds;
    }
    return false;
  }

  // Range with only max (LessThan) - e.g., "Under 5 min"
  if (range.minSeconds === undefined && range.maxSeconds !== undefined) {
    if (modifier === CriterionModifier.LessThan) {
      return value.value === range.maxSeconds;
    }
    return false;
  }

  // Range with only min (GreaterThan) - e.g., "Over 2 hours"
  if (range.minSeconds !== undefined && range.maxSeconds === undefined) {
    if (modifier === CriterionModifier.GreaterThan) {
      return value.value === range.minSeconds;
    }
    return false;
  }

  return false;
}

// Create icon for duration range
function createDurationIcon(rangeId: string): React.ReactNode {
  // Use different opacity based on duration to give visual hints
  let opacity = 0.7;
  
  switch (rangeId) {
    case "under-5":
      opacity = 0.5;
      break;
    case "5-15":
      opacity = 0.6;
      break;
    case "15-30":
      opacity = 0.7;
      break;
    case "30-60":
      opacity = 0.8;
      break;
    case "1-2h":
      opacity = 0.9;
      break;
    case "over-2h":
      opacity = 1.0;
      break;
  }

  return (
    <FontAwesomeIcon
      icon={faClock}
      style={{ marginRight: "0.5em", opacity }}
      fixedWidth
    />
  );
}

function useDurationFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [showCustom, setShowCustom] = useState(false);

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as NumberCriterion;

    const newCriterion = filter.makeCriterion(option.type) as NumberCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: NumberCriterion | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier, value } = criterion;

  // Check which duration range is currently active
  const activeRange = useMemo((): string | null => {
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return null;
    }

    for (const range of DURATION_RANGES) {
      if (matchesDurationRange(criterion, range)) {
        return range.id;
      }
    }

    // If there's a value but it doesn't match any preset, it's custom
    if (value.value !== undefined) {
      return "custom";
    }

    return null;
  }, [criterion, modifier, value]);

  // Build selected modifiers (any/none)
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Duration range options with localized labels
  const durationRangeOptions = useMemo(() => {
    return DURATION_RANGES.map((range) => ({
      id: range.id,
      label: intl.formatMessage({ id: range.labelId, defaultMessage: range.defaultLabel }),
      icon: createDurationIcon(range.id),
      minSeconds: range.minSeconds,
      maxSeconds: range.maxSeconds,
    }));
  }, [intl]);

  // Format duration value for display
  const formatDuration = (seconds: number | undefined): string => {
    if (seconds === undefined) return "";
    return TextUtils.secondsAsTimeString(seconds);
  };

  // Build selected items list
  const selected = useMemo(() => {
    const items: Option[] = [];

    // Add modifier if any/none
    if (selectedModifiers.any) {
      items.push({
        id: "any",
        label: `(${intl.formatMessage({ id: "criterion_modifier_values.any" })})`,
        className: "modifier-object",
      });
    }
    if (selectedModifiers.none) {
      items.push({
        id: "none",
        label: `(${intl.formatMessage({ id: "criterion_modifier_values.none" })})`,
        className: "modifier-object",
      });
    }

    // Add active duration range or custom value
    if (activeRange && activeRange !== "custom") {
      const range = durationRangeOptions.find((r) => r.id === activeRange);
      if (range) {
        items.push({
          id: activeRange,
          label: range.label,
          icon: range.icon,
        });
      }
    } else if (activeRange === "custom" || (value.value !== undefined && !selectedModifiers.any && !selectedModifiers.none)) {
      // Custom duration value
      let label = "";
      if (modifier === CriterionModifier.Equals) {
        label = formatDuration(value.value);
      } else if (modifier === CriterionModifier.NotEquals) {
        label = `≠ ${formatDuration(value.value)}`;
      } else if (modifier === CriterionModifier.GreaterThan) {
        label = `> ${formatDuration(value.value)}`;
      } else if (modifier === CriterionModifier.LessThan) {
        label = `< ${formatDuration(value.value)}`;
      } else if (modifier === CriterionModifier.Between) {
        label = `${formatDuration(value.value)} - ${formatDuration(value.value2)}`;
      } else if (modifier === CriterionModifier.NotBetween) {
        label = `not ${formatDuration(value.value)} - ${formatDuration(value.value2)}`;
      }

      if (label) {
        items.push({
          id: "custom",
          label,
          icon: createDurationIcon("custom"),
        });
      }
    }

    return items;
  }, [intl, selectedModifiers, activeRange, durationRangeOptions, value, modifier, formatDuration]);

  // Build candidates list
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Show modifier options when nothing is selected
    if (!selectedModifiers.any && !selectedModifiers.none && !activeRange) {
      items.push({
        id: "any",
        label: `(${intl.formatMessage({ id: "criterion_modifier_values.any" })})`,
        className: "modifier-object",
        canExclude: false,
      });
      items.push({
        id: "none",
        label: `(${intl.formatMessage({ id: "criterion_modifier_values.none" })})`,
        className: "modifier-object",
        canExclude: false,
      });
    }

    // Don't show duration ranges if modifier is any/none
    if (selectedModifiers.any || selectedModifiers.none) {
      return items;
    }

    // Add duration range options (excluding active one)
    durationRangeOptions
      .filter((r) => r.id !== activeRange)
      .forEach((range) => {
        items.push({
          id: range.id,
          label: range.label,
          icon: range.icon,
          canExclude: false,
        });
      });

    // Add custom option if not in custom mode
    if (activeRange !== "custom") {
      items.push({
        id: "custom",
        label: intl.formatMessage({ id: "duration_ranges.custom", defaultMessage: "Custom duration..." }),
        icon: createDurationIcon("custom"),
        canExclude: false,
      });
    }

    return items;
  }, [intl, selectedModifiers, activeRange, durationRangeOptions]);

  const onSelect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Handle modifier selection
        const newCriterion = cloneDeep(criterion);
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = { value: undefined, value2: undefined };
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = { value: undefined, value2: undefined };
        }
        setCriterion(newCriterion);
        setShowCustom(false);
        return;
      }

      if (v.id === "custom") {
        // Show custom duration inputs
        setShowCustom(true);
        return;
      }

      // Handle duration range selection
      const range = DURATION_RANGES.find((r) => r.id === v.id);
      if (!range) return;

      const newCriterion = cloneDeep(criterion);

      if (range.minSeconds !== undefined && range.maxSeconds !== undefined) {
        // Range with both bounds - use Between
        newCriterion.modifier = CriterionModifier.Between;
        newCriterion.value = { value: range.minSeconds, value2: range.maxSeconds };
      } else if (range.minSeconds === undefined && range.maxSeconds !== undefined) {
        // Max only (e.g., Under 5 min) - use LessThan
        newCriterion.modifier = CriterionModifier.LessThan;
        newCriterion.value = { value: range.maxSeconds, value2: undefined };
      } else if (range.minSeconds !== undefined && range.maxSeconds === undefined) {
        // Min only (e.g., Over 2 hours) - use GreaterThan
        newCriterion.modifier = CriterionModifier.GreaterThan;
        newCriterion.value = { value: range.minSeconds, value2: undefined };
      }

      setCriterion(newCriterion);
      setShowCustom(false);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (_v: Option, _exclude: boolean) => {
      setCriterion(null);
      setShowCustom(false);
    },
    [setCriterion]
  );

  const onCustomValueChanged = useCallback(
    (newValue: INumberValue) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.value = newValue;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onCustomModifierChanged = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    showCustom,
    setShowCustom,
    criterion,
    onCustomValueChanged,
    onCustomModifierChanged,
  };
}

// Parse duration string (HH:MM:SS or MM:SS or seconds) to seconds
function parseDuration(input: string): number | undefined {
  if (!input.trim()) return undefined;
  
  // If it's just a number, treat as seconds
  if (/^\d+$/.test(input.trim())) {
    return parseInt(input.trim(), 10);
  }
  
  // Parse HH:MM:SS or MM:SS format
  const parts = input.split(":").map((p) => parseInt(p, 10));
  if (parts.some(isNaN)) return undefined;
  
  if (parts.length === 3) {
    return parts[0] * 3600 + parts[1] * 60 + parts[2];
  } else if (parts.length === 2) {
    return parts[0] * 60 + parts[1];
  }
  
  return undefined;
}

// Custom duration input section
interface ICustomDurationInputProps {
  criterion: NumberCriterion;
  onValueChanged: (value: INumberValue) => void;
  onModifierChanged: (m: CriterionModifier) => void;
  onCancel: () => void;
}

const CustomDurationInput: React.FC<ICustomDurationInputProps> = ({
  criterion,
  onValueChanged,
  onModifierChanged,
  onCancel,
}) => {
  const intl = useIntl();
  const { value, modifier } = criterion;
  
  const [inputValue, setInputValue] = useState(
    value?.value !== undefined ? TextUtils.secondsAsTimeString(value.value) : ""
  );
  const [inputValue2, setInputValue2] = useState(
    value?.value2 !== undefined ? TextUtils.secondsAsTimeString(value.value2) : ""
  );

  const modifierOptions = [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.GreaterThan,
    CriterionModifier.LessThan,
    CriterionModifier.Between,
    CriterionModifier.NotBetween,
  ];

  const modifierLabels: Record<CriterionModifier, string> = {
    [CriterionModifier.Equals]: intl.formatMessage({ id: "criterion_modifier.equals", defaultMessage: "is" }),
    [CriterionModifier.NotEquals]: intl.formatMessage({ id: "criterion_modifier.not_equals", defaultMessage: "is not" }),
    [CriterionModifier.GreaterThan]: intl.formatMessage({ id: "criterion_modifier.greater_than", defaultMessage: "is greater than" }),
    [CriterionModifier.LessThan]: intl.formatMessage({ id: "criterion_modifier.less_than", defaultMessage: "is less than" }),
    [CriterionModifier.Between]: intl.formatMessage({ id: "criterion_modifier.between", defaultMessage: "between" }),
    [CriterionModifier.NotBetween]: intl.formatMessage({ id: "criterion_modifier.not_between", defaultMessage: "not between" }),
    [CriterionModifier.Includes]: "",
    [CriterionModifier.IncludesAll]: "",
    [CriterionModifier.Excludes]: "",
    [CriterionModifier.IsNull]: "",
    [CriterionModifier.NotNull]: "",
    [CriterionModifier.MatchesRegex]: "",
    [CriterionModifier.NotMatchesRegex]: "",
  };

  function onInputChange(newValue: string, property: "value" | "value2") {
    if (property === "value") {
      setInputValue(newValue);
    } else {
      setInputValue2(newValue);
    }
    
    const parsed = parseDuration(newValue);
    const valueCopy = { ...value };
    valueCopy[property] = parsed;
    onValueChanged(valueCopy);
  }

  const showSecondInput =
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween;

  return (
    <div className="custom-duration-input">
      <Form.Group className="modifier-select">
        <Form.Control
          as="select"
          value={modifier}
          onChange={(e) => onModifierChanged(e.target.value as CriterionModifier)}
        >
          {modifierOptions.map((m) => (
            <option key={m} value={m}>
              {modifierLabels[m]}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <Form.Control
          type="text"
          value={inputValue}
          onChange={(e) => onInputChange(e.target.value, "value")}
          placeholder={
            showSecondInput
              ? intl.formatMessage({ id: "criterion.from", defaultMessage: "From (HH:MM:SS)" })
              : "HH:MM:SS"
          }
        />
      </Form.Group>
      {showSecondInput && (
        <Form.Group>
          <Form.Control
            type="text"
            value={inputValue2}
            onChange={(e) => onInputChange(e.target.value, "value2")}
            placeholder={intl.formatMessage({ id: "criterion.to", defaultMessage: "To (HH:MM:SS)" })}
          />
        </Form.Group>
      )}
      <button
        type="button"
        className="custom-cancel-button"
        onClick={onCancel}
      >
        ✕
      </button>
    </div>
  );
};

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarDurationFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = useDurationFilterState({ option, filter, setFilter });

  const customInput = state.showCustom ? (
    <CustomDurationInput
      criterion={state.criterion}
      onValueChanged={state.onCustomValueChanged}
      onModifierChanged={state.onCustomModifierChanged}
      onCancel={() => state.setShowCustom(false)}
    />
  ) : null;

  return (
    <SidebarListFilter
      title={title}
      candidates={state.candidates}
      onSelect={state.onSelect}
      onUnselect={state.onUnselect}
      selected={state.selected}
      canExclude={false}
      singleValue={true}
      sectionID={sectionID}
      preCandidates={customInput}
    />
  );
};
