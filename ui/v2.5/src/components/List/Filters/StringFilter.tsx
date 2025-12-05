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
// NEW IMPROVED SIDEBAR STRING FILTER
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

function useStringFilterState(props: {
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
      value &&
      modifier !== CriterionModifier.IsNull &&
      modifier !== CriterionModifier.NotNull
    );
  }, [value, modifier]);

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

    // Add active value with its modifier
    if (hasActiveValue) {
      // Show modifier if not the default
      if (modifier !== defaultModifier) {
        items.push({
          id: "modifier",
          label: `(${getModifierLabel(intl, modifier)})`,
          className: "modifier-object",
        });
      }

      items.push({
        id: "value",
        label: value,
        icon: createStringIcon(),
      });
    }

    return items;
  }, [intl, selectedModifiers, hasActiveValue, value, modifier, defaultModifier]);

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
      if (v.id === "any" || v.id === "none" || v.id === "modifier") {
        // Reset modifier to default
        const newCriterion = cloneDeep(criterion);
        newCriterion.modifier = defaultModifier;
        if (v.id === "any" || v.id === "none") {
          newCriterion.value = "";
          setCriterion(null);
        } else {
          setCriterion(newCriterion);
        }
      } else if (v.id === "value") {
        // Clear value
        setCriterion(null);
        setInputValue("");
      }
    },
    [criterion, setCriterion, defaultModifier]
  );

  const onInputSubmit = useCallback(
    (inputVal: string, selectedModifier: CriterionModifier) => {
      if (!inputVal.trim()) {
        setCriterion(null);
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = selectedModifier;
      newCriterion.value = inputVal.trim();
      setCriterion(newCriterion);
      setInputValue("");
    },
    [criterion, setCriterion]
  );

  // Get available text modifiers (excluding null modifiers)
  const textModifiers = useMemo(() => {
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
    onInputSubmit,
    textModifiers,
    defaultModifier,
    selectedModifiers,
    hasActiveValue,
  };
}

// String input component with modifier dropdown
interface IStringInputProps {
  inputValue: string;
  setInputValue: (value: string) => void;
  onSubmit: (value: string, modifier: CriterionModifier) => void;
  modifiers: CriterionModifier[];
  defaultModifier: CriterionModifier;
  placeholder?: string;
  disabled?: boolean;
}

const StringInput: React.FC<IStringInputProps> = ({
  inputValue,
  setInputValue,
  onSubmit,
  modifiers,
  defaultModifier,
  placeholder,
  disabled,
}) => {
  const intl = useIntl();
  const [selectedModifier, setSelectedModifier] = useState(defaultModifier);
  const [focus] = useFocus();

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputValue.trim()) {
      onSubmit(inputValue, selectedModifier);
    }
  };

  const showModifierDropdown = modifiers.length > 1;

  return (
    <InputGroup className="string-input-group">
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
        ref={focus}
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder || intl.formatMessage({ id: "search" })}
        disabled={disabled}
      />
    </InputGroup>
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
  const state = useStringFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled = state.selectedModifiers.any || state.selectedModifiers.none;

  const stringInput = (
    <StringInput
      inputValue={state.inputValue}
      setInputValue={state.setInputValue}
      onSubmit={state.onInputSubmit}
      modifiers={state.textModifiers}
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
