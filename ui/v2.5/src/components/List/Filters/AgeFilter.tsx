import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUser } from "@fortawesome/free-solid-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { NumberCriterion } from "src/models/list-filter/criteria/criterion";
import { INumberValue } from "src/models/list-filter/types";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";

// Age range preset type
interface AgeRange {
  id: string;
  labelId: string;
  defaultLabel: string;
  min?: number;
  max?: number;
}

// Predefined age ranges
const AGE_RANGES: AgeRange[] = [
  { id: "18-25", labelId: "age_ranges.18_25", defaultLabel: "18-25", min: 18, max: 25 },
  { id: "26-35", labelId: "age_ranges.26_35", defaultLabel: "26-35", min: 26, max: 35 },
  { id: "36-45", labelId: "age_ranges.36_45", defaultLabel: "36-45", min: 36, max: 45 },
  { id: "46-55", labelId: "age_ranges.46_55", defaultLabel: "46-55", min: 46, max: 55 },
  { id: "56-65", labelId: "age_ranges.56_65", defaultLabel: "56-65", min: 56, max: 65 },
  { id: "65+", labelId: "age_ranges.65_plus", defaultLabel: "65+", min: 65, max: undefined },
];

// Check if current criterion matches an age range
function matchesAgeRange(
  criterion: NumberCriterion,
  range: AgeRange
): boolean {
  const { modifier, value } = criterion;

  // Range with both min and max
  if (range.min !== undefined && range.max !== undefined) {
    if (modifier === CriterionModifier.Between) {
      return value.value === range.min && value.value2 === range.max;
    }
    return false;
  }

  // Range with only min (e.g., 65+)
  if (range.min !== undefined && range.max === undefined) {
    if (modifier === CriterionModifier.GreaterThan) {
      return value.value === range.min - 1; // GreaterThan 64 means 65+
    }
    return false;
  }

  return false;
}

// Create icon for age range
function createAgeIcon(rangeId: string): React.ReactNode {
  // Use different opacity/styling based on age group to give visual hints
  let opacity = 0.7;
  
  switch (rangeId) {
    case "18-25":
      opacity = 0.6;
      break;
    case "26-35":
      opacity = 0.7;
      break;
    case "36-45":
      opacity = 0.8;
      break;
    case "46-55":
      opacity = 0.85;
      break;
    case "56-65":
      opacity = 0.9;
      break;
    case "65+":
      opacity = 1.0;
      break;
  }

  return (
    <FontAwesomeIcon
      icon={faUser}
      style={{ marginRight: "0.5em", opacity }}
      fixedWidth
    />
  );
}

function useAgeFilterState(props: {
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

  // Check which age range is currently active
  const activeRange = useMemo((): string | null => {
    if (
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      return null;
    }

    for (const range of AGE_RANGES) {
      if (matchesAgeRange(criterion, range)) {
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

  // Age range options with localized labels
  const ageRangeOptions = useMemo(() => {
    return AGE_RANGES.map((range) => ({
      id: range.id,
      label: intl.formatMessage({ id: range.labelId, defaultMessage: range.defaultLabel }),
      icon: createAgeIcon(range.id),
      min: range.min,
      max: range.max,
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

    // Add active age range or custom value
    if (activeRange && activeRange !== "custom") {
      const range = ageRangeOptions.find((r) => r.id === activeRange);
      if (range) {
        items.push({
          id: activeRange,
          label: range.label,
          icon: range.icon,
        });
      }
    } else if (activeRange === "custom" || (value.value !== undefined && !selectedModifiers.any && !selectedModifiers.none)) {
      // Custom age value
      let label = "";
      if (modifier === CriterionModifier.Equals) {
        label = `${value.value}`;
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
          icon: createAgeIcon("custom"),
        });
      }
    }

    return items;
  }, [intl, selectedModifiers, activeRange, ageRangeOptions, value, modifier]);

  // Build candidates list
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Show modifier options when nothing is selected
    if (!selectedModifiers.any && !selectedModifiers.none && !activeRange) {
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

    // Don't show age ranges if modifier is any/none
    if (selectedModifiers.any || selectedModifiers.none) {
      return items;
    }

    // Add age range options (excluding active one)
    ageRangeOptions
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
        label: intl.formatMessage({ id: "age_ranges.custom", defaultMessage: "Custom age..." }),
        icon: createAgeIcon("custom"),
        canExclude: false,
      });
    }

    return items;
  }, [intl, selectedModifiers, activeRange, ageRangeOptions]);

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
        // Show custom age inputs
        setShowCustom(true);
        return;
      }

      // Handle age range selection
      const range = AGE_RANGES.find((r) => r.id === v.id);
      if (!range) return;

      const newCriterion = cloneDeep(criterion);

      if (range.min !== undefined && range.max !== undefined) {
        // Range with both bounds - use Between
        newCriterion.modifier = CriterionModifier.Between;
        newCriterion.value = { value: range.min, value2: range.max };
      } else if (range.min !== undefined && range.max === undefined) {
        // Open-ended range (e.g., 65+) - use GreaterThan
        newCriterion.modifier = CriterionModifier.GreaterThan;
        newCriterion.value = { value: range.min - 1, value2: undefined }; // GreaterThan 64 = 65+
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

// Custom age input section
interface ICustomAgeInputProps {
  criterion: NumberCriterion;
  onValueChanged: (value: INumberValue) => void;
  onModifierChanged: (m: CriterionModifier) => void;
}

const CustomAgeInput: React.FC<ICustomAgeInputProps> = ({
  criterion,
  onValueChanged,
  onModifierChanged,
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

  function onChanged(newValue: number | undefined, property: "value" | "value2") {
    const valueCopy = { ...value };
    valueCopy[property] = newValue;
    onValueChanged(valueCopy);
  }

  const showSecondInput =
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween;

  return (
    <div className="custom-age-input">
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
          type="number"
          min={0}
          max={120}
          value={value?.value ?? ""}
          onChange={(e) => onChanged(e.target.value ? parseInt(e.target.value, 10) : undefined, "value")}
          placeholder={
            showSecondInput
              ? intl.formatMessage({ id: "criterion.from", defaultMessage: "From" })
              : intl.formatMessage({ id: "age", defaultMessage: "Age" })
          }
        />
      </Form.Group>
      {showSecondInput && (
        <Form.Group>
          <Form.Control
            type="number"
            min={0}
            max={120}
            value={value?.value2 ?? ""}
            onChange={(e) => onChanged(e.target.value ? parseInt(e.target.value, 10) : undefined, "value2")}
            placeholder={intl.formatMessage({ id: "criterion.to", defaultMessage: "To" })}
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
}

export const SidebarAgeFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = useAgeFilterState({ option, filter, setFilter });

  const customInput = state.showCustom ? (
    <CustomAgeInput
      criterion={state.criterion}
      onValueChanged={state.onCustomValueChanged}
      onModifierChanged={state.onCustomModifierChanged}
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

