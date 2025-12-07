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
} from "src/models/list-filter/criteria/criterion";
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

type InputMode = "none" | "browse" | "regex";

function usePathFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;
  const { configuration } = React.useContext(ConfigurationContext);

  const [inputValue, setInputValue] = useState("");
  const [inputMode, setInputMode] = useState<InputMode>("none");

  // Get library paths from configuration
  const libraryPaths = useMemo(() => {
    return configuration?.general.stashes.map((s) => s.path) ?? [];
  }, [configuration]);

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

  // Check if current value matches a library path
  const activeLibraryPath = useMemo(() => {
    if (!hasActiveValue || !value) return null;
    return libraryPaths.find((p) => p === value) ?? null;
  }, [hasActiveValue, value, libraryPaths]);

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

  // Build candidates list (modifier options + library paths + browse/regex options)
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Don't show candidates if any/none is selected
    if (selectedModifiers.any || selectedModifiers.none) {
      return items;
    }

    // Show modifier options when nothing is selected
    if (!hasActiveValue) {
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

    // Add library paths as presets (excluding the currently active one)
    libraryPaths
      .filter((p) => p !== activeLibraryPath)
      .forEach((path) => {
        // Get just the folder name for display
        const folderName = path.split(/[/\\]/).filter(Boolean).pop() || path;
        items.push({
          id: `library_${path}`,
          label: folderName,
          icon: createPathIcon(),
          className: "preset-option",
          canExclude: true,
        });
      });

    // Add browse and regex options
    if (inputMode !== "browse") {
      items.push({
        id: "browse",
        label: intl.formatMessage({ id: "actions.browse", defaultMessage: "Browse..." }),
        icon: createPathIcon(),
        className: "preset-option",
        canExclude: false,
      });
    }
    if (inputMode !== "regex") {
      items.push({
        id: "regex",
        label: intl.formatMessage({ id: "actions.regex", defaultMessage: "Regex..." }),
        className: "preset-option",
        canExclude: false,
      });
    }

    return items;
  }, [intl, selectedModifiers, hasActiveValue, supportsIsNull, supportsNotNull, libraryPaths, activeLibraryPath, inputMode]);

  const onSelect = useCallback(
    (v: Option, exclude: boolean) => {
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
        setInputMode("none");
      } else if (v.id === "browse") {
        // Show folder browser
        setInputMode("browse");
      } else if (v.id === "regex") {
        // Show regex input
        setInputMode("regex");
      } else if (v.id.startsWith("library_")) {
        // Handle library path selection
        const path = v.id.replace("library_", "");
        const newCriterion = cloneDeep(criterion);
        newCriterion.modifier = exclude ? CriterionModifier.Excludes : CriterionModifier.Includes;
        newCriterion.value = path;
        setCriterion(newCriterion);
        setInputMode("none");
      }
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.id === "any" || v.id === "none" || v.id === "value") {
        setCriterion(null);
        setInputValue("");
        setInputMode("none");
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
      setInputMode("none");
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
    inputMode,
    setInputMode,
    libraryPaths,
  };
}

// Browse input component with folder select
interface IBrowseInputProps {
  currentValue: string;
  onPathChange: (value: string, modifier: CriterionModifier) => void;
  onCancel: () => void;
  libraryPaths: string[];
  disabled?: boolean;
}

const BrowseInput: React.FC<IBrowseInputProps> = ({
  currentValue,
  onPathChange,
  onCancel,
  libraryPaths,
  disabled,
}) => {
  const intl = useIntl();
  const [selectedModifier, setSelectedModifier] = useState(CriterionModifier.Includes);

  const modifiers = [
    CriterionModifier.Includes,
    CriterionModifier.Excludes,
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
  ];

  const handlePathChange = (value: string) => {
    if (value.trim()) {
      onPathChange(value, selectedModifier);
    }
  };

  return (
    <div className="path-browse-input">
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

// Regex input component
interface IRegexInputProps {
  onPathChange: (value: string, modifier: CriterionModifier) => void;
  onCancel: () => void;
  disabled?: boolean;
}

const RegexInput: React.FC<IRegexInputProps> = ({
  onPathChange,
  onCancel,
  disabled,
}) => {
  const intl = useIntl();
  const [inputValue, setInputValue] = useState("");
  const [selectedModifier, setSelectedModifier] = useState(CriterionModifier.MatchesRegex);

  const modifiers = [
    CriterionModifier.MatchesRegex,
    CriterionModifier.NotMatchesRegex,
  ];

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputValue.trim()) {
      onPathChange(inputValue, selectedModifier);
    }
    if (e.key === "Escape") {
      onCancel();
    }
  };

  return (
    <div className="path-regex-input-container">
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
        <Form.Control
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={intl.formatMessage({ id: "dialogs.path_filter.regex_placeholder", defaultMessage: "Path regex pattern..." })}
          disabled={disabled}
          autoFocus
        />
      </InputGroup>
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

  // Render appropriate input based on mode
  let inputSection = null;
  if (state.inputMode === "browse") {
    inputSection = (
      <BrowseInput
        currentValue={state.currentValue}
        onPathChange={state.onPathChange}
        onCancel={() => state.setInputMode("none")}
        libraryPaths={state.libraryPaths}
        disabled={inputDisabled}
      />
    );
  } else if (state.inputMode === "regex") {
    inputSection = (
      <RegexInput
        onPathChange={state.onPathChange}
        onCancel={() => state.setInputMode("none")}
        disabled={inputDisabled}
      />
    );
  }

  return (
    <SidebarListFilter
      title={title}
      candidates={state.candidates}
      onSelect={state.onSelect}
      onUnselect={state.onUnselect}
      selected={state.selected}
      canExclude={true}
      singleValue={true}
      sectionID={sectionID}
      preCandidates={inputSection}
    />
  );
};
