import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Dropdown, Form, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faFolder } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import { CriterionModifier } from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import {
  ModifierCriterion,
  CriterionValue,
  CriterionOption,
  StringCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";

// ============================================================================
// LEGACY EXPORTS FOR BACKWARDS COMPATIBILITY
// ============================================================================

interface IInputFilterProps {
  criterion: ModifierCriterion<CriterionValue>;
  onValueChanged: (value: string) => void;
}

export const PathFilter: React.FC<IInputFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  // don't show folder select for regex
  const regex =
    criterion.modifier === CriterionModifier.MatchesRegex ||
    criterion.modifier === CriterionModifier.NotMatchesRegex;

  return (
    <Form.Group>
      {regex ? (
        <Form.Control
          className="btn-secondary"
          type={criterion.modifierCriterionOption().inputType}
          onChange={(v) => onValueChanged(v.target.value)}
          value={criterion.value ? criterion.value.toString() : ""}
        />
      ) : (
        <FolderSelect
          currentDirectory={criterion.value ? criterion.value.toString() : ""}
          onChangeDirectory={onValueChanged}
          collapsible
          quotePath
          hideError
          defaultDirectories={libraryPaths}
        />
      )}
    </Form.Group>
  );
};

// ============================================================================
// NEW IMPROVED SIDEBAR PATH FILTER
// ============================================================================

// Get modifier label for display
function getModifierLabel(intl: ReturnType<typeof useIntl>, modifier: CriterionModifier): string {
  const labels: Partial<Record<CriterionModifier, string>> = {
    [CriterionModifier.Equals]: intl.formatMessage({ id: "criterion_modifier.equals", defaultMessage: "is" }),
    [CriterionModifier.NotEquals]: intl.formatMessage({ id: "criterion_modifier.not_equals", defaultMessage: "is not" }),
    [CriterionModifier.Includes]: intl.formatMessage({ id: "criterion_modifier.includes", defaultMessage: "includes" }),
    [CriterionModifier.Excludes]: intl.formatMessage({ id: "criterion_modifier.excludes", defaultMessage: "excludes" }),
    [CriterionModifier.MatchesRegex]: intl.formatMessage({ id: "criterion_modifier.matches_regex", defaultMessage: "regex" }),
    [CriterionModifier.NotMatchesRegex]: intl.formatMessage({ id: "criterion_modifier.not_matches_regex", defaultMessage: "not regex" }),
  };
  return labels[modifier] || modifier;
}

// Create icon for path value
function createPathIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faFolder}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function usePathFilterState(props: {
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
    if (ret) return ret as StringCriterion;

    const newCriterion = filter.makeCriterion(option.type) as StringCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: StringCriterion | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier ?? CriterionModifier.Includes;
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

  // Get display label for the current value
  const getValueLabel = useCallback(() => {
    if (!hasActiveValue || !value) return null;

    // Shorten path for display
    const displayPath = value.length > 30 ? "..." + value.slice(-27) : value;

    switch (modifier) {
      case CriterionModifier.Equals:
        return displayPath;
      case CriterionModifier.NotEquals:
        return `≠ ${displayPath}`;
      case CriterionModifier.Includes:
        return displayPath;
      case CriterionModifier.Excludes:
        return `excl. ${displayPath}`;
      case CriterionModifier.MatchesRegex:
        return `/${value}/`;
      case CriterionModifier.NotMatchesRegex:
        return `!/${value}/`;
      default:
        return displayPath;
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
        icon: createPathIcon(),
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
      if (v.id === "any" || v.id === "none" || v.id === "value") {
        setCriterion(null);
        setInputValue("");
      }
    },
    [setCriterion]
  );

  const onPathChange = useCallback(
    (pathValue: string, mod: CriterionModifier) => {
      if (!pathValue.trim()) {
        setCriterion(null);
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = mod;
      newCriterion.value = pathValue;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  // Get available path modifiers (excluding null modifiers)
  const pathModifiers = useMemo(() => {
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
    onPathChange,
    pathModifiers,
    defaultModifier,
    selectedModifiers,
    hasActiveValue,
    currentValue: value ?? "",
    currentModifier: modifier,
  };
}

// Path input component with modifier dropdown and folder select
interface IPathInputProps {
  currentValue: string;
  onPathChange: (value: string, modifier: CriterionModifier) => void;
  modifiers: CriterionModifier[];
  defaultModifier: CriterionModifier;
  disabled?: boolean;
}

const PathInput: React.FC<IPathInputProps> = ({
  currentValue,
  onPathChange,
  modifiers,
  defaultModifier,
  disabled,
}) => {
  const intl = useIntl();
  const [selectedModifier, setSelectedModifier] = useState(defaultModifier);
  const { configuration } = React.useContext(ConfigurationContext);
  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  // Check if we should show regex input or folder select
  const isRegex =
    selectedModifier === CriterionModifier.MatchesRegex ||
    selectedModifier === CriterionModifier.NotMatchesRegex;

  const handlePathChange = (value: string) => {
    onPathChange(value, selectedModifier);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && e.currentTarget.value.trim()) {
      onPathChange(e.currentTarget.value, selectedModifier);
    }
  };

  return (
    <div className="path-input-container">
      <InputGroup className="path-input-group">
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
      </InputGroup>

      {isRegex ? (
        <Form.Control
          type="text"
          onChange={(e) => handlePathChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={intl.formatMessage({ id: "path", defaultMessage: "Path regex" })}
          disabled={disabled}
          className="path-regex-input"
        />
      ) : (
        <div className="folder-select-container">
          <FolderSelect
            currentDirectory={currentValue}
            onChangeDirectory={handlePathChange}
            collapsible
            quotePath
            hideError
            defaultDirectories={libraryPaths}
          />
        </div>
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

export const SidebarPathFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = usePathFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled = state.selectedModifiers.any || state.selectedModifiers.none;

  const pathInput = (
    <PathInput
      currentValue={state.currentValue}
      onPathChange={state.onPathChange}
      modifiers={state.pathModifiers}
      defaultModifier={state.defaultModifier}
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
      preCandidates={pathInput}
    />
  );
};
