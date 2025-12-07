import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCalendarDay,
  faCalendarWeek,
  faCalendar,
  faCalendarAlt,
  faEdit,
} from "@fortawesome/free-solid-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";
import { IDateValue } from "src/models/list-filter/types";
import {
  ModifierCriterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { DateInput } from "src/components/Shared/DateInput";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import TextUtils from "src/utils/text";

// Date preset types
type DatePreset =
  | "today"
  | "yesterday"
  | "this_week"
  | "last_7_days"
  | "this_month"
  | "last_30_days"
  | "this_year"
  | "last_year"
  | "custom";

// Get date string in YYYY-MM-DD format
function formatDate(date: Date): string {
  return TextUtils.dateToString(date);
}

// Calculate date range for a preset
function getPresetDateRange(preset: DatePreset): { start: string; end: string } {
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

  switch (preset) {
    case "today":
      return { start: formatDate(today), end: formatDate(today) };
    case "yesterday": {
      const yesterday = new Date(today);
      yesterday.setDate(yesterday.getDate() - 1);
      return { start: formatDate(yesterday), end: formatDate(yesterday) };
    }
    case "this_week": {
      const startOfWeek = new Date(today);
      startOfWeek.setDate(today.getDate() - today.getDay());
      return { start: formatDate(startOfWeek), end: formatDate(today) };
    }
    case "last_7_days": {
      const sevenDaysAgo = new Date(today);
      sevenDaysAgo.setDate(today.getDate() - 6);
      return { start: formatDate(sevenDaysAgo), end: formatDate(today) };
    }
    case "this_month": {
      const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
      return { start: formatDate(startOfMonth), end: formatDate(today) };
    }
    case "last_30_days": {
      const thirtyDaysAgo = new Date(today);
      thirtyDaysAgo.setDate(today.getDate() - 29);
      return { start: formatDate(thirtyDaysAgo), end: formatDate(today) };
    }
    case "this_year": {
      const startOfYear = new Date(today.getFullYear(), 0, 1);
      return { start: formatDate(startOfYear), end: formatDate(today) };
    }
    case "last_year": {
      const lastYearStart = new Date(today.getFullYear() - 1, 0, 1);
      const lastYearEnd = new Date(today.getFullYear() - 1, 11, 31);
      return { start: formatDate(lastYearStart), end: formatDate(lastYearEnd) };
    }
    default:
      return { start: "", end: "" };
  }
}

// Check if current criterion matches a preset
function matchesPreset(
  criterion: ModifierCriterion<IDateValue>,
  preset: DatePreset
): boolean {
  if (preset === "custom") return false;

  const range = getPresetDateRange(preset);

  if (
    criterion.modifier === CriterionModifier.Equals &&
    range.start === range.end
  ) {
    return criterion.value.value === range.start;
  }

  if (criterion.modifier === CriterionModifier.Between) {
    return (
      criterion.value.value === range.start &&
      criterion.value.value2 === range.end
    );
  }

  return false;
}

// Get preset icon
function getPresetIcon(preset: DatePreset): React.ReactNode {
  const iconMap: Record<DatePreset, typeof faCalendarDay> = {
    today: faCalendarDay,
    yesterday: faCalendarDay,
    this_week: faCalendarWeek,
    last_7_days: faCalendarWeek,
    this_month: faCalendar,
    last_30_days: faCalendar,
    this_year: faCalendarAlt,
    last_year: faCalendarAlt,
    custom: faEdit,
  };

  return (
    <FontAwesomeIcon
      icon={iconMap[preset]}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function useDateFilterState(props: {
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
    if (ret) return ret as ModifierCriterion<IDateValue>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<IDateValue>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<IDateValue> | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier, value } = criterion;

  // Check which preset is currently active
  const activePreset = useMemo((): DatePreset | null => {
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return null;
    }

    const presets: DatePreset[] = [
      "today",
      "yesterday",
      "this_week",
      "last_7_days",
      "this_month",
      "last_30_days",
      "this_year",
      "last_year",
    ];

    for (const preset of presets) {
      if (matchesPreset(criterion, preset)) {
        return preset;
      }
    }

    // If there's a value but it doesn't match any preset, it's custom
    if (value.value || value.value2) {
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

  // Preset options
  const presetOptions = useMemo(() => {
    const presets: { id: DatePreset; labelId: string }[] = [
      { id: "today", labelId: "date_presets.today" },
      { id: "yesterday", labelId: "date_presets.yesterday" },
      { id: "this_week", labelId: "date_presets.this_week" },
      { id: "last_7_days", labelId: "date_presets.last_7_days" },
      { id: "this_month", labelId: "date_presets.this_month" },
      { id: "last_30_days", labelId: "date_presets.last_30_days" },
      { id: "this_year", labelId: "date_presets.this_year" },
      { id: "last_year", labelId: "date_presets.last_year" },
    ];

    return presets.map((p) => ({
      id: p.id,
      label: intl.formatMessage({ id: p.labelId, defaultMessage: p.id.replace(/_/g, " ") }),
      icon: getPresetIcon(p.id),
    }));
  }, [intl]);

  // Build selected items list
  const selected = useMemo(() => {
    const items: Option[] = [];

    // Add modifier if any/none
    Object.entries(selectedModifiers)
      .filter((v) => v[1])
      .forEach((v) => {
        items.push({
          id: v[0],
          label: `(${intl.formatMessage({
            id: `criterion_modifier_values.${v[0]}`,
          })})`,
          className: "modifier-object",
        });
      });

    // Add active preset or custom date
    if (activePreset && activePreset !== "custom") {
      const preset = presetOptions.find((p) => p.id === activePreset);
      if (preset) {
        items.push({
          id: activePreset,
          label: preset.label,
          icon: preset.icon,
        });
      }
    } else if (activePreset === "custom" || (value.value && !selectedModifiers.any && !selectedModifiers.none)) {
      // Custom date range
      let label = "";
      if (modifier === CriterionModifier.Equals) {
        label = value.value;
      } else if (modifier === CriterionModifier.NotEquals) {
        label = `≠ ${value.value}`;
      } else if (modifier === CriterionModifier.GreaterThan) {
        label = `> ${value.value}`;
      } else if (modifier === CriterionModifier.LessThan) {
        label = `< ${value.value}`;
      } else if (modifier === CriterionModifier.Between) {
        label = `${value.value} - ${value.value2}`;
      } else if (modifier === CriterionModifier.NotBetween) {
        label = `not ${value.value} - ${value.value2}`;
      }

      if (label) {
        items.push({
          id: "custom",
          label,
          icon: getPresetIcon("custom"),
        });
      }
    }

    return items;
  }, [intl, selectedModifiers, activePreset, presetOptions, value, modifier]);

  // Build candidates list
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Show modifier options when nothing is selected
    if (!selectedModifiers.any && !selectedModifiers.none && !activePreset) {
      items.push({
        id: "any",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.any",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
      items.push({
        id: "none",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.none",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
    }

    // Don't show presets if modifier is any/none
    if (selectedModifiers.any || selectedModifiers.none) {
      return items;
    }

    // Add preset options (excluding active one)
    presetOptions
      .filter((p) => p.id !== activePreset)
      .forEach((preset) => {
        items.push({
          id: preset.id,
          label: preset.label,
          icon: preset.icon,
          canExclude: false,
        });
      });

    // Add custom option if not in custom mode
    if (activePreset !== "custom") {
      items.push({
        id: "custom",
        label: intl.formatMessage({ id: "date_presets.custom", defaultMessage: "Custom date..." }),
        icon: getPresetIcon("custom"),
        canExclude: false,
      });
    }

    return items;
  }, [intl, selectedModifiers, activePreset, presetOptions]);

  const onSelect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Handle modifier selection
        const newCriterion = cloneDeep(criterion);
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = { value: "", value2: "" };
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = { value: "", value2: "" };
        }
        setCriterion(newCriterion);
        setShowCustom(false);
        return;
      }

      if (v.id === "custom") {
        // Show custom date inputs
        setShowCustom(true);
        return;
      }

      // Handle preset selection
      const preset = v.id as DatePreset;
      const range = getPresetDateRange(preset);
      const newCriterion = cloneDeep(criterion);

      if (range.start === range.end) {
        // Single day - use Equals
        newCriterion.modifier = CriterionModifier.Equals;
        newCriterion.value = { value: range.start, value2: "" };
      } else {
        // Date range - use Between
        newCriterion.modifier = CriterionModifier.Between;
        newCriterion.value = { value: range.start, value2: range.end };
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
    (newValue: IDateValue) => {
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

// Custom date input section
interface ICustomDateInputProps {
  criterion: ModifierCriterion<IDateValue>;
  onValueChanged: (value: IDateValue) => void;
  onModifierChanged: (m: CriterionModifier) => void;
  isTime?: boolean;
}

const CustomDateInput: React.FC<ICustomDateInputProps> = ({
  criterion,
  onValueChanged,
  onModifierChanged,
  isTime = false,
}) => {
  const intl = useIntl();
  const { value, modifier } = criterion;

  const modifierOptions = [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.GreaterThan,
    CriterionModifier.LessThan,
    CriterionModifier.Between,
    CriterionModifier.NotBetween,
  ];

  function onChanged(newValue: string, property: "value" | "value2") {
    const valueCopy = { ...value };
    valueCopy[property] = newValue;
    onValueChanged(valueCopy);
  }

  const showSecondInput =
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween;

  return (
    <div className="custom-date-input">
      <Form.Group className="modifier-select">
        <Form.Control
          as="select"
          value={modifier}
          onChange={(e) => onModifierChanged(e.target.value as CriterionModifier)}
        >
          {modifierOptions.map((m) => (
            <option key={m} value={m}>
              {ModifierCriterion.getModifierLabel(intl, m)}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={
            showSecondInput
              ? intl.formatMessage({ id: "criterion.from" })
              : intl.formatMessage({ id: "criterion.value" })
          }
          isTime={isTime}
        />
      </Form.Group>
      {showSecondInput && (
        <Form.Group>
          <DateInput
            value={value?.value2 ?? ""}
            onValueChange={(v) => onChanged(v, "value2")}
            placeholder={intl.formatMessage({ id: "criterion.to" })}
            isTime={isTime}
          />
        </Form.Group>
      )}
    </div>
  );
};

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
  isTime?: boolean;
}

export const SidebarDateFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
  isTime = false,
}) => {
  const state = useDateFilterState({ option, filter, setFilter });

  const customInput = state.showCustom ? (
    <div className="custom-date-input-wrapper">
      <CustomDateInput
        criterion={state.criterion}
        onValueChanged={state.onCustomValueChanged}
        onModifierChanged={state.onCustomModifierChanged}
        isTime={isTime}
      />
      <button
        type="button"
        className="custom-cancel-button"
        onClick={() => state.setShowCustom(false)}
      >
        ✕
      </button>
    </div>
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

// Keep old exports for backwards compatibility
export interface IDateFilterProps {
  criterion: ModifierCriterion<IDateValue>;
  onValueChanged: (value: IDateValue) => void;
}

export const DateFilter: React.FC<IDateFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const { value } = criterion;

  function onChanged(newValue: string, property: "value" | "value2") {
    const valueCopy = { ...value };

    valueCopy[property] = newValue;
    onValueChanged(valueCopy);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.value" })}
        />
      </Form.Group>
    );
  }

  let lowerControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.GreaterThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    lowerControl = (
      <Form.Group>
        <DateInput
          value={value?.value ?? ""}
          onValueChange={(v) => onChanged(v, "value")}
          placeholder={intl.formatMessage({ id: "criterion.greater_than" })}
        />
      </Form.Group>
    );
  }

  let upperControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.LessThan ||
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    upperControl = (
      <Form.Group>
        <DateInput
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
          }
          onValueChange={(v) =>
            onChanged(
              v,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          placeholder={intl.formatMessage({ id: "criterion.less_than" })}
        />
      </Form.Group>
    );
  }

  return (
    <>
      {equalsControl}
      {lowerControl}
      {upperControl}
    </>
  );
};

// Legacy hook export
export function useDateCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<IDateValue>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<IDateValue>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<IDateValue>) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onValueChanged = useCallback(
    (value: IDateValue) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.value = value;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  return {
    criterion,
    setCriterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  };
}
