import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Dropdown, Form, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faHashtag } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import {
  NumberCriterion,
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { Option, SidebarListFilter, SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierSelectorButtons } from "../ModifierSelect";
import { NumberField } from "src/utils/form";

// ============================================================================
// LEGACY EXPORTS FOR BACKWARDS COMPATIBILITY
// ============================================================================

interface IDurationFilterProps {
  criterion: NumberCriterion;
  onValueChanged: (value: INumberValue) => void;
}

export function useNumberCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as NumberCriterion;

    const newCriterion = filter.makeCriterion(option.type) as NumberCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: NumberCriterion) => {
      const newFilter = cloneDeep(filter);

      if (false && !c.isValid()) {
        const newCriteria = filter.criteria.filter((cc) => {
          return cc.criterionOption.type !== c.criterionOption.type;
        });

        newFilter.criteria = newCriteria;
      } else {
        let found = false;

        const newCriteria = filter.criteria.map((cc) => {
          if (cc.criterionOption.type === c.criterionOption.type) {
            found = true;
            return c;
          }

          return cc;
        });

        if (!found) {
          newCriteria.push(c);
        }

        newFilter.criteria = newCriteria;
      }
      setFilter(newFilter);
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  function onValueChanged(value: INumberValue) {
    const newCriterion = cloneDeep(criterion);
    newCriterion.value = value;
    setCriterion(newCriterion);
  }

  return {
    criterion,
    setCriterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  };
}

interface INumberSelectedItemsProps {
  criterion: NumberCriterion | null;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
  onClear: () => void;
}

export const NumberSelectedItems: React.FC<INumberSelectedItemsProps> = ({
  criterion,
  defaultModifier,
  onChangedModifierSelect,
  onClear,
}) => {
  const intl = useIntl();

  if (criterion?.value.value === undefined) {
    return null;
  }

  const getValueLabel = () => {
    if (!criterion?.value) return null;

    const { value, value2 } = criterion.value;

    switch (criterion.modifier) {
      case CriterionModifier.Equals:
        return value?.toString();
      case CriterionModifier.NotEquals:
        return `≠ ${value}`;
      case CriterionModifier.GreaterThan:
        return `> ${value}`;
      case CriterionModifier.LessThan:
        return `< ${value}`;
      case CriterionModifier.Between:
        return `${value} - ${value2}`;
      case CriterionModifier.NotBetween:
        return `not ${value} - ${value2}`;
      default:
        return value?.toString();
    }
  };

  const valueLabel = getValueLabel();

  return (
    <ul className="selected-list">
      {criterion?.modifier != defaultModifier && criterion?.modifier ? (
        <SelectedItem
          className="modifier-object"
          label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
          onClick={() => onChangedModifierSelect(defaultModifier)}
        />
      ) : null}
      {valueLabel ? (
        <SelectedItem label={valueLabel} onClick={onClear} />
      ) : null}
    </ul>
  );
};

export const NumberFilter: React.FC<IDurationFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const { value } = criterion;

  function onChanged(
    event: React.ChangeEvent<HTMLInputElement>,
    property: "value" | "value2"
  ) {
    const numericValue = parseInt(event.target.value, 10);
    const valueCopy = { ...value };

    valueCopy[property] = !Number.isNaN(numericValue) ? numericValue : 0;
    onValueChanged(valueCopy);
  }

  let equalsControl: JSX.Element | null = null;
  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals
  ) {
    equalsControl = (
      <Form.Group>
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
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
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(e, "value")
          }
          value={value?.value ?? ""}
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
        <NumberField
          className="btn-secondary"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            onChanged(
              e,
              criterion.modifier === CriterionModifier.LessThan
                ? "value"
                : "value2"
            )
          }
          value={
            (criterion.modifier === CriterionModifier.LessThan
              ? value?.value
              : value?.value2) ?? ""
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

// ============================================================================
// NEW IMPROVED SIDEBAR NUMBER FILTER
// ============================================================================

// Get modifier label for display
function getModifierLabel(intl: ReturnType<typeof useIntl>, modifier: CriterionModifier): string {
  const labels: Partial<Record<CriterionModifier, string>> = {
    [CriterionModifier.Equals]: intl.formatMessage({ id: "criterion_modifier.equals", defaultMessage: "is" }),
    [CriterionModifier.NotEquals]: intl.formatMessage({ id: "criterion_modifier.not_equals", defaultMessage: "is not" }),
    [CriterionModifier.GreaterThan]: intl.formatMessage({ id: "criterion_modifier.greater_than", defaultMessage: ">" }),
    [CriterionModifier.LessThan]: intl.formatMessage({ id: "criterion_modifier.less_than", defaultMessage: "<" }),
    [CriterionModifier.Between]: intl.formatMessage({ id: "criterion_modifier.between", defaultMessage: "between" }),
    [CriterionModifier.NotBetween]: intl.formatMessage({ id: "criterion_modifier.not_between", defaultMessage: "not between" }),
  };
  return labels[modifier] || modifier;
}

// Create icon for number value
function createNumberIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faHashtag}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function useNumberFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [inputValue, setInputValue] = useState("");
  const [inputValue2, setInputValue2] = useState("");

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

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier ?? CriterionModifier.Equals;
  const modifierOptions = modifierCriterionOption?.modifierOptions ?? [];

  const { modifier, value } = criterion;

  // Check if any/none modifiers are supported
  const supportsIsNull = modifierOptions.includes(CriterionModifier.IsNull);
  const supportsNotNull = modifierOptions.includes(CriterionModifier.NotNull);

  // Build selected modifiers (any/none)
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Determine if there's an active value filter
  const hasActiveValue = useMemo(() => {
    return (
      value.value !== undefined &&
      modifier !== CriterionModifier.IsNull &&
      modifier !== CriterionModifier.NotNull
    );
  }, [value, modifier]);

  // Get display label for the current value
  const getValueLabel = useCallback(() => {
    if (!hasActiveValue) return null;

    const { value: v, value2: v2 } = value;

    switch (modifier) {
      case CriterionModifier.Equals:
        return `${v}`;
      case CriterionModifier.NotEquals:
        return `≠ ${v}`;
      case CriterionModifier.GreaterThan:
        return `> ${v}`;
      case CriterionModifier.LessThan:
        return `< ${v}`;
      case CriterionModifier.Between:
        return `${v} - ${v2}`;
      case CriterionModifier.NotBetween:
        return `not ${v} - ${v2}`;
      default:
        return `${v}`;
    }
  }, [hasActiveValue, value, modifier]);

  // Build selected items list
  const selected = useMemo(() => {
    const items: Option[] = [];

    // Add modifier if any/none
    if (selectedModifiers.any) {
      items.push({
        id: "any",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.any",
        })})`,
        className: "modifier-object",
      });
    }
    if (selectedModifiers.none) {
      items.push({
        id: "none",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.none",
        })})`,
        className: "modifier-object",
      });
    }

    // Add active value
    const valueLabel = getValueLabel();
    if (valueLabel) {
      items.push({
        id: "value",
        label: valueLabel,
        icon: createNumberIcon(),
      });
    }

    return items;
  }, [intl, selectedModifiers, getValueLabel]);

  // Build candidates list (modifier options)
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Show modifier options when nothing is selected
    if (!selectedModifiers.any && !selectedModifiers.none && !hasActiveValue) {
      if (supportsNotNull) {
        items.push({
          id: "any",
          label: `(${intl.formatMessage({
            id: "criterion_modifier_values.any",
          })})`,
          className: "modifier-object",
          canExclude: false,
        });
      }
      if (supportsIsNull) {
        items.push({
          id: "none",
          label: `(${intl.formatMessage({
            id: "criterion_modifier_values.none",
          })})`,
          className: "modifier-object",
          canExclude: false,
        });
      }
    }

    return items;
  }, [intl, selectedModifiers, hasActiveValue, supportsIsNull, supportsNotNull]);

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
      }
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.id === "any" || v.id === "none") {
        setCriterion(null);
      } else if (v.id === "value") {
        setCriterion(null);
        setInputValue("");
        setInputValue2("");
      }
    },
    [setCriterion]
  );

  const onInputSubmit = useCallback(
    (val1: string, val2: string, selectedModifier: CriterionModifier) => {
      const num1 = val1.trim() ? parseInt(val1, 10) : undefined;
      const num2 = val2.trim() ? parseInt(val2, 10) : undefined;

      if (num1 === undefined || Number.isNaN(num1)) {
        setCriterion(null);
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = selectedModifier;
      newCriterion.value = { value: num1, value2: num2 };
      setCriterion(newCriterion);
      setInputValue("");
      setInputValue2("");
    },
    [criterion, setCriterion]
  );

  // Get available number modifiers (excluding null modifiers)
  const numberModifiers = useMemo(() => {
    return modifierOptions.filter(
      (m) => m !== CriterionModifier.IsNull && m !== CriterionModifier.NotNull
    );
  }, [modifierOptions]);

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    inputValue,
    setInputValue,
    inputValue2,
    setInputValue2,
    onInputSubmit,
    numberModifiers,
    defaultModifier,
    selectedModifiers,
    hasActiveValue,
  };
}

// Number input component with modifier dropdown
interface INumberInputProps {
  inputValue: string;
  setInputValue: (value: string) => void;
  inputValue2: string;
  setInputValue2: (value: string) => void;
  onSubmit: (value1: string, value2: string, modifier: CriterionModifier) => void;
  modifiers: CriterionModifier[];
  defaultModifier: CriterionModifier;
  placeholder?: string;
  disabled?: boolean;
}

const NumberInput: React.FC<INumberInputProps> = ({
  inputValue,
  setInputValue,
  inputValue2,
  setInputValue2,
  onSubmit,
  modifiers,
  defaultModifier,
  placeholder,
  disabled,
}) => {
  const intl = useIntl();
  const [selectedModifier, setSelectedModifier] = useState(defaultModifier);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputValue.trim()) {
      onSubmit(inputValue, inputValue2, selectedModifier);
    }
  };

  const showModifierDropdown = modifiers.length > 1;
  const showSecondInput =
    selectedModifier === CriterionModifier.Between ||
    selectedModifier === CriterionModifier.NotBetween;

  return (
    <div className="number-input-container">
      <InputGroup className="number-input-group">
        {showModifierDropdown && (
          <InputGroup.Prepend>
            <Dropdown>
              <Dropdown.Toggle
                variant="secondary"
                disabled={disabled}
                className="modifier-dropdown-toggle"
              >
                {getModifierLabel(intl, selectedModifier)}
                <Icon icon={faChevronDown} className="dropdown-icon" />
              </Dropdown.Toggle>
              <Dropdown.Menu className="bg-secondary text-white">
                {modifiers.map((m) => (
                  <Dropdown.Item
                    key={m}
                    className="bg-secondary text-white"
                    active={m === selectedModifier}
                    onClick={() => setSelectedModifier(m)}
                  >
                    {getModifierLabel(intl, m)}
                  </Dropdown.Item>
                ))}
              </Dropdown.Menu>
            </Dropdown>
          </InputGroup.Prepend>
        )}
        <Form.Control
          type="number"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder || intl.formatMessage({ id: "criterion.value", defaultMessage: "Value" })}
          disabled={disabled}
        />
      </InputGroup>
      {showSecondInput && (
        <InputGroup className="number-input-group number-input-second">
          <InputGroup.Prepend>
            <InputGroup.Text>-</InputGroup.Text>
          </InputGroup.Prepend>
          <Form.Control
            type="number"
            value={inputValue2}
            onChange={(e) => setInputValue2(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={intl.formatMessage({ id: "criterion.to", defaultMessage: "To" })}
            disabled={disabled}
          />
        </InputGroup>
      )}
    </div>
  );
};

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  placeholder?: string;
  sectionID?: string;
}

export const SidebarNumberFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  placeholder,
  sectionID,
}) => {
  const state = useNumberFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled = state.selectedModifiers.any || state.selectedModifiers.none;

  const numberInput = (
    <NumberInput
      inputValue={state.inputValue}
      setInputValue={state.setInputValue}
      inputValue2={state.inputValue2}
      setInputValue2={state.setInputValue2}
      onSubmit={state.onInputSubmit}
      modifiers={state.numberModifiers}
      defaultModifier={state.defaultModifier}
      placeholder={placeholder}
      disabled={inputDisabled}
    />
  );

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
      preCandidates={numberInput}
    />
  );
};
