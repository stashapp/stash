import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Dropdown, Form, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faFont } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { CriterionModifier } from "../../../core/generated-graphql";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter, SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierSelectorButtons } from "../ModifierSelect";
import useFocus from "src/utils/focus";

// Legacy hook exports for backwards compatibility
export function useStringCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<string>;

    const newCriterion = filter.makeCriterion(
      option.type
    ) as ModifierCriterion<string>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<string>) => {
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

export function useModifierCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const { criterion, setCriterion } = useStringCriterion(
    option,
    filter,
    setFilter
  );
  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  const onValueChange = useCallback(
    (value: string) => {
      if (!value.trim()) {
        setFilter(filter.removeCriterion(option.type));
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = criterion?.modifier
        ? criterion.modifier
        : defaultModifier;
      newCriterion.value = value;
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
    },
    [criterion, setCriterion, filter, setFilter, option.type, defaultModifier]
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
    onValueChange,
    onChangedModifierSelect,
  };
}

// Legacy components for backwards compatibility
interface ISelectedItemsProps {
  criterion: ModifierCriterion<string> | null;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
  onValueChange: (value: string) => void;
}

export const SelectedItems: React.FC<ISelectedItemsProps> = ({
  criterion,
  defaultModifier,
  onChangedModifierSelect,
  onValueChange,
}) => {
  const intl = useIntl();

  if (!criterion?.value) {
    return null;
  }

  return (
    <ul className="selected-list">
      {criterion?.modifier != defaultModifier && criterion?.modifier ? (
        <SelectedItem
          className="modifier-object"
          label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
          onClick={() => onChangedModifierSelect(defaultModifier)}
        />
      ) : null}
      {criterion?.value ? (
        <SelectedItem
          label={criterion.value}
          onClick={() => onValueChange("")}
        />
      ) : null}
    </ul>
  );
};

interface IModifierControlsProps {
  modifierOptions: CriterionModifier[] | undefined;
  currentModifier: CriterionModifier;
  onChangedModifierSelect: (m: CriterionModifier) => void;
}

export const ModifierControls: React.FC<IModifierControlsProps> = ({
  modifierOptions,
  currentModifier,
  onChangedModifierSelect,
}) => (
  <ModifierSelectorButtons
    options={modifierOptions}
    value={currentModifier}
    onChanged={onChangedModifierSelect}
  />
);

interface IStringFilterProps {
  criterion: ModifierCriterion<string>;
  onValueChanged: (value: string) => void;
  placeholder?: string;
}

export const StringFilter: React.FC<IStringFilterProps> = ({
  criterion,
  onValueChanged,
  placeholder,
}) => {
  const intl = useIntl();

  function onValueChange(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged(event.target.value);
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={onValueChange}
          value={criterion.value}
          placeholder={placeholder || intl.formatMessage({ id: "search" })}
        />
      </Form.Group>
    </div>
  );
};

// ============================================================================
// NEW IMPROVED SIDEBAR STRING FILTER WITH MULTI-VALUE SUPPORT
// ============================================================================

// Get modifier label for display
function getModifierLabel(intl: ReturnType<typeof useIntl>, modifier: CriterionModifier): string {
  const labels: Partial<Record<CriterionModifier, string>> = {
    [CriterionModifier.Equals]: intl.formatMessage({ id: "criterion_modifier.equals", defaultMessage: "is" }),
    [CriterionModifier.NotEquals]: intl.formatMessage({ id: "criterion_modifier.not_equals", defaultMessage: "is not" }),
    [CriterionModifier.Includes]: intl.formatMessage({ id: "criterion_modifier.includes", defaultMessage: "includes" }),
    [CriterionModifier.Excludes]: intl.formatMessage({ id: "criterion_modifier.excludes", defaultMessage: "excludes" }),
    [CriterionModifier.MatchesRegex]: intl.formatMessage({ id: "criterion_modifier.matches_regex", defaultMessage: "matches regex" }),
    [CriterionModifier.NotMatchesRegex]: intl.formatMessage({ id: "criterion_modifier.not_matches_regex", defaultMessage: "not matches regex" }),
  };
  return labels[modifier] || modifier;
}

// Create icon for string value
function createStringIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faFont}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

// Escape special regex characters in a string
function escapeRegex(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

// Convert multiple values to a regex pattern
function valuesToRegex(values: string[]): string {
  if (values.length === 0) return "";
  if (values.length === 1) return escapeRegex(values[0]);
  // Case-insensitive OR pattern
  return `(?i)(${values.map(escapeRegex).join("|")})`;
}

// Parse a regex pattern back to values (best effort)
function regexToValues(regex: string): string[] {
  if (!regex) return [];
  
  // Try to parse (?i)(val1|val2|val3) pattern
  const multiMatch = regex.match(/^\(\?i\)\((.+)\)$/);
  if (multiMatch) {
    // Split by | but handle escaped pipes
    const inner = multiMatch[1];
    const values = inner.split("|").map((v) => 
      v.replace(/\\(.)/g, "$1") // Unescape
    );
    return values;
  }
  
  // Single value - unescape
  return [regex.replace(/\\(.)/g, "$1")];
}

// Check if a regex pattern looks like our multi-value pattern
function isMultiValueRegex(regex: string): boolean {
  return /^\(\?i\)\(.+\)$/.test(regex) || !regex.includes("(");
}

function useMultiStringFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [inputValue, setInputValue] = useState("");

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ModifierCriterion<string>;

    const newCriterion = filter.makeCriterion(option.type) as ModifierCriterion<string>;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ModifierCriterion<string> | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
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

  // Parse current values from the criterion
  const currentValues = useMemo(() => {
    if (!value || selectedModifiers.any || selectedModifiers.none) return [];
    
    // If using regex modifier, try to parse values
    if (modifier === CriterionModifier.MatchesRegex || modifier === CriterionModifier.NotMatchesRegex) {
      if (isMultiValueRegex(value)) {
        return regexToValues(value);
      }
      // Complex regex - show as single value
      return [value];
    }
    
    // Single value for other modifiers
    return value ? [value] : [];
  }, [value, modifier, selectedModifiers]);

  // Determine if values are excluded (using NotMatchesRegex or Excludes)
  const isExcluded = useMemo(() => {
    return modifier === CriterionModifier.NotMatchesRegex || modifier === CriterionModifier.Excludes;
  }, [modifier]);

  // Determine if there's an active value filter
  const hasActiveValue = useMemo(() => {
    return currentValues.length > 0;
  }, [currentValues]);

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

    // Add each value as a separate item
    currentValues.forEach((val, index) => {
      items.push({
        id: `value_${index}`,
        label: val,
        icon: createStringIcon(),
        className: isExcluded ? "excluded-item" : undefined,
      });
    });

    return items;
  }, [intl, selectedModifiers, currentValues, isExcluded]);

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
          newCriterion.value = "";
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = "";
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
      } else if (v.id.startsWith("value_")) {
        // Remove this specific value
        const index = parseInt(v.id.replace("value_", ""), 10);
        const newValues = currentValues.filter((_, i) => i !== index);
        
        if (newValues.length === 0) {
          setCriterion(null);
        } else {
          const newCriterion = cloneDeep(criterion);
          newCriterion.value = valuesToRegex(newValues);
          // Keep the same modifier type (include/exclude)
          if (!isExcluded) {
            newCriterion.modifier = CriterionModifier.MatchesRegex;
          } else {
            newCriterion.modifier = CriterionModifier.NotMatchesRegex;
          }
          setCriterion(newCriterion);
        }
      }
    },
    [criterion, setCriterion, currentValues, isExcluded]
  );

  // Add a new value
  const addValue = useCallback(
    (newValue: string, exclude: boolean) => {
      const trimmed = newValue.trim();
      if (!trimmed) return;
      
      // Check if value already exists
      if (currentValues.includes(trimmed)) {
        setInputValue("");
        return;
      }

      const newValues = [...currentValues, trimmed];
      const newCriterion = cloneDeep(criterion);
      newCriterion.value = valuesToRegex(newValues);
      
      // Set modifier based on exclude flag
      // If we already have values with a modifier, keep that modifier unless changing modes
      if (currentValues.length === 0) {
        // First value - set based on exclude flag
        newCriterion.modifier = exclude ? CriterionModifier.NotMatchesRegex : CriterionModifier.MatchesRegex;
      } else {
        // Keep existing modifier
        newCriterion.modifier = isExcluded ? CriterionModifier.NotMatchesRegex : CriterionModifier.MatchesRegex;
      }
      
      setCriterion(newCriterion);
      setInputValue("");
    },
    [criterion, setCriterion, currentValues, isExcluded]
  );

  // Toggle between include/exclude mode
  const toggleExclude = useCallback(() => {
    if (currentValues.length === 0) return;
    
    const newCriterion = cloneDeep(criterion);
    newCriterion.modifier = isExcluded ? CriterionModifier.MatchesRegex : CriterionModifier.NotMatchesRegex;
    setCriterion(newCriterion);
  }, [criterion, setCriterion, currentValues, isExcluded]);

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    inputValue,
    setInputValue,
    addValue,
    toggleExclude,
    selectedModifiers,
    hasActiveValue,
    isExcluded,
    currentValues,
  };
}

// Multi-value string input component
interface IMultiStringInputProps {
  inputValue: string;
  setInputValue: (value: string) => void;
  onAdd: (value: string, exclude: boolean) => void;
  placeholder?: string;
  disabled?: boolean;
  isExcluded: boolean;
  hasValues: boolean;
  onToggleExclude: () => void;
}

const MultiStringInput: React.FC<IMultiStringInputProps> = ({
  inputValue,
  setInputValue,
  onAdd,
  placeholder,
  disabled,
  isExcluded,
  hasValues,
  onToggleExclude,
}) => {
  const intl = useIntl();
  const [inputRef, setFocus] = useFocus();

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputValue.trim()) {
      e.preventDefault();
      // Shift+Enter to exclude
      onAdd(inputValue, e.shiftKey);
    }
  };

  return (
    <div className="multi-string-input">
      <InputGroup className="string-input-group">
        {hasValues && (
          <InputGroup.Prepend>
            <button
              type="button"
              className={`mode-toggle-button ${isExcluded ? "exclude-mode" : "include-mode"}`}
              onClick={onToggleExclude}
              title={isExcluded 
                ? intl.formatMessage({ id: "criterion_modifier.excludes", defaultMessage: "Excluding" })
                : intl.formatMessage({ id: "criterion_modifier.includes", defaultMessage: "Including" })
              }
            >
              {isExcluded ? "−" : "+"}
            </button>
          </InputGroup.Prepend>
        )}
        <Form.Control
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder || intl.formatMessage({ id: "actions.add_value", defaultMessage: "Add value..." })}
          disabled={disabled}
        />
      </InputGroup>
      <div className="input-hint">
        {intl.formatMessage({ id: "dialogs.string_filter.hint", defaultMessage: "Enter to add, Shift+Enter to exclude" })}
      </div>
    </div>
  );
};

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  placeholder?: string;
  modifier?: CriterionModifier;
  sectionID?: string;
}

export const SidebarStringFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  placeholder,
  sectionID,
}) => {
  const state = useMultiStringFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled = state.selectedModifiers.any || state.selectedModifiers.none;

  const stringInput = (
    <MultiStringInput
      inputValue={state.inputValue}
      setInputValue={state.setInputValue}
      onAdd={state.addValue}
      placeholder={placeholder}
      disabled={inputDisabled}
      isExcluded={state.isExcluded}
      hasValues={state.hasActiveValue}
      onToggleExclude={state.toggleExclude}
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
      singleValue={false}
      sectionID={sectionID}
      preCandidates={stringInput}
    />
  );
};

// Convenience exports for specific string filters
export const SidebarTattoosFilter: React.FC<ISidebarFilter> = (props) => (
  <SidebarStringFilter
    {...props}
    placeholder={props.placeholder || "tattoos"}
  />
);
