import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Dropdown, Form, InputGroup } from "react-bootstrap";
import { useIntl, IntlShape } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faLink } from "@fortawesome/free-solid-svg-icons";
import { IStashIDValue } from "src/models/list-filter/types";
import {
  ModifierCriterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { StashIDCriterion } from "src/models/list-filter/criteria/stash-ids";
import { CriterionModifier } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { Icon } from "src/components/Shared/Icon";
import { cloneDeep } from "lodash-es";

// ============================================================================
// LEGACY EXPORTS FOR BACKWARDS COMPATIBILITY
// ============================================================================

interface IStashIDFilterProps {
  criterion: ModifierCriterion<IStashIDValue>;
  onValueChanged: (value: IStashIDValue) => void;
}

// Legacy hook for backwards compatibility
export function useStashIDCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as StashIDCriterion;

    const newCriterion = filter.makeCriterion(option.type) as StashIDCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: StashIDCriterion) => {
      const newFilter = cloneDeep(filter);

      // replace or add the criterion
      const newCriteria = filter.criteria.filter((cc) => {
        return cc.criterionOption.type !== c.criterionOption.type;
      });
      newCriteria.push(c);
      newFilter.criteria = newCriteria;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  const onValueChanged = useCallback(
    (value: IStashIDValue) => {
      const newCriterion = criterion.clone();
      newCriterion.value = value;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onChangedModifierSelect = useCallback(
    (modifier: CriterionModifier) => {
      const newCriterion = criterion.clone();
      newCriterion.modifier = modifier;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  return {
    criterion,
    setCriterion,
    onValueChanged,
    onChangedModifierSelect,
    defaultModifier,
    modifierOptions,
  };
}

export const StashIDFilter: React.FC<IStashIDFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function onEndpointChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      endpoint: event.target.value,
      stashID: criterion.value.stashID,
    });
  }

  function onStashIDChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      stashID: event.target.value,
      endpoint: criterion.value.endpoint,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={onEndpointChanged}
          value={value ? value.endpoint : ""}
          placeholder={intl.formatMessage({ id: "stash_id_endpoint" })}
        />
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <Form.Control
              className="btn-secondary"
              onChange={onStashIDChanged}
              value={value ? value.stashID : ""}
              placeholder={intl.formatMessage({ id: "stash_id" })}
            />
          </Form.Group>
        )}
    </div>
  );
};

// ============================================================================
// NEW IMPROVED SIDEBAR STASH ID FILTER
// ============================================================================

// Known stash-box URL patterns
const STASHBOX_URL_PATTERNS = [
  // Standard stash-box URLs: https://stashdb.org/performers/uuid
  /^https?:\/\/([^/]+)\/(performers|scenes|studios|tags|galleries|groups|movies)\/([a-f0-9-]+)/i,
  // Direct UUID URLs: https://stashdb.org/uuid
  /^https?:\/\/([^/]+)\/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i,
];

// Parse a stash-box URL and extract endpoint and stash ID
function parseStashBoxURL(input: string): { endpoint: string; stashID: string } | null {
  const trimmed = input.trim();
  
  for (const pattern of STASHBOX_URL_PATTERNS) {
    const match = trimmed.match(pattern);
    if (match) {
      // First pattern: [full, endpoint, type, uuid]
      // Second pattern: [full, endpoint, uuid]
      const endpoint = match[1];
      const stashID = match[3] || match[2];
      
      // Validate it looks like a UUID
      if (/^[a-f0-9-]+$/i.test(stashID) && stashID.length >= 32) {
        return { endpoint, stashID };
      }
    }
  }
  
  return null;
}

// Check if input looks like a URL
function looksLikeURL(input: string): boolean {
  return /^https?:\/\//i.test(input.trim());
}

// Create icon for stash ID value
function createStashIDIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faLink}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function useStashIDFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [inputEndpoint, setInputEndpoint] = useState("");
  const [inputStashID, setInputStashID] = useState("");

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as StashIDCriterion;

    return filter.makeCriterion(option.type) as StashIDCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: StashIDCriterion | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier, value } = criterion;

  // Build selected modifiers (any/none)
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Determine if there's an active stash ID value
  const hasActiveValue = useMemo(() => {
    return (
      value?.stashID &&
      modifier !== CriterionModifier.IsNull &&
      modifier !== CriterionModifier.NotNull
    );
  }, [value, modifier]);

  // Get display label for the current value
  const getValueLabel = useCallback(() => {
    if (!hasActiveValue || !value?.stashID) return null;

    const stashIDShort =
      value.stashID.length > 24
        ? value.stashID.slice(0, 12) + "..." + value.stashID.slice(-8)
        : value.stashID;

    if (value.endpoint) {
      return `${value.endpoint}: ${stashIDShort}`;
    }
    return stashIDShort;
  }, [hasActiveValue, value]);

  // Get modifier label for display
  const getModifierLabel = useCallback(() => {
    if (modifier === CriterionModifier.Equals) {
      return intl.formatMessage({ id: "criterion_modifier.equals" });
    } else if (modifier === CriterionModifier.NotEquals) {
      return intl.formatMessage({ id: "criterion_modifier.not_equals" });
    }
    return null;
  }, [modifier, intl]);

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

    // Add active value with modifier
    const valueLabel = getValueLabel();
    if (valueLabel) {
      const modifierLabel = getModifierLabel();
      if (modifierLabel && modifier !== CriterionModifier.Equals) {
        items.push({
          id: "modifier",
          label: `(${modifierLabel})`,
          className: "modifier-object",
        });
      }
      items.push({
        id: "value",
        label: valueLabel,
        icon: createStashIDIcon(),
      });
    }

    return items;
  }, [intl, selectedModifiers, getValueLabel, getModifierLabel, modifier]);

  // Build candidates list (modifier options)
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Show modifier options when nothing is selected
    if (!selectedModifiers.any && !selectedModifiers.none && !hasActiveValue) {
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

    return items;
  }, [intl, selectedModifiers, hasActiveValue]);

  const onSelect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Handle modifier selection
        const newCriterion = cloneDeep(criterion);
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = { endpoint: "", stashID: "" };
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = { endpoint: "", stashID: "" };
        }
        setCriterion(newCriterion);
      }
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (
        v.id === "any" ||
        v.id === "none" ||
        v.id === "value" ||
        v.id === "modifier"
      ) {
        setCriterion(null);
        setInputEndpoint("");
        setInputStashID("");
      }
    },
    [setCriterion]
  );

  const onInputSubmit = useCallback(
    (endpoint: string, stashID: string, notEquals: boolean) => {
      if (!stashID.trim()) {
        setCriterion(null);
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = notEquals
        ? CriterionModifier.NotEquals
        : CriterionModifier.Equals;
      newCriterion.value = {
        endpoint: endpoint.trim(),
        stashID: stashID.trim(),
      };
      setCriterion(newCriterion);
      setInputEndpoint("");
      setInputStashID("");
    },
    [criterion, setCriterion]
  );

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    inputEndpoint,
    setInputEndpoint,
    inputStashID,
    setInputStashID,
    onInputSubmit,
    selectedModifiers,
    hasActiveValue,
  };
}

// Get localized label for modifier
function getModifierLabel(intl: IntlShape, modifier: CriterionModifier): string {
  const labels: Record<string, string> = {
    [CriterionModifier.Equals]: intl.formatMessage({
      id: "criterion_modifier.equals",
      defaultMessage: "is",
    }),
    [CriterionModifier.NotEquals]: intl.formatMessage({
      id: "criterion_modifier.not_equals",
      defaultMessage: "is not",
    }),
  };
  return labels[modifier] || modifier;
}

// Stash ID input component with URL paste support
interface IStashIDInputProps {
  inputEndpoint: string;
  setInputEndpoint: (value: string) => void;
  inputStashID: string;
  setInputStashID: (value: string) => void;
  onSubmit: (endpoint: string, stashID: string, notEquals: boolean) => void;
  disabled?: boolean;
}

const StashIDInput: React.FC<IStashIDInputProps> = ({
  inputEndpoint,
  setInputEndpoint,
  inputStashID,
  setInputStashID,
  onSubmit,
  disabled,
}) => {
  const intl = useIntl();
  const [selectedModifier, setSelectedModifier] = useState(
    CriterionModifier.Equals
  );
  const [urlParsed, setUrlParsed] = useState(false);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputStashID.trim()) {
      onSubmit(
        inputEndpoint,
        inputStashID,
        selectedModifier === CriterionModifier.NotEquals
      );
    }
  };

  // Handle paste event to detect URLs
  const handlePaste = useCallback(
    (e: React.ClipboardEvent<HTMLInputElement>) => {
      const pastedText = e.clipboardData.getData("text");
      
      if (looksLikeURL(pastedText)) {
        const parsed = parseStashBoxURL(pastedText);
        if (parsed) {
          e.preventDefault();
          setInputEndpoint(parsed.endpoint);
          setInputStashID(parsed.stashID);
          setUrlParsed(true);
          // Auto-clear the "URL parsed" message after 3 seconds
          setTimeout(() => setUrlParsed(false), 3000);
        }
      }
    },
    [setInputEndpoint, setInputStashID]
  );

  // Handle change to detect URLs typed/pasted without using paste event
  const handleStashIDChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.target.value;
      
      // Check if it looks like a URL was pasted
      if (looksLikeURL(value)) {
        const parsed = parseStashBoxURL(value);
        if (parsed) {
          setInputEndpoint(parsed.endpoint);
          setInputStashID(parsed.stashID);
          setUrlParsed(true);
          setTimeout(() => setUrlParsed(false), 3000);
          return;
        }
      }
      
      setInputStashID(value);
      setUrlParsed(false);
    },
    [setInputEndpoint, setInputStashID]
  );

  const modifiers = [CriterionModifier.Equals, CriterionModifier.NotEquals];

  return (
    <div className="stash-id-input-container">
      <Form.Control
        type="text"
        value={inputEndpoint}
        onChange={(e) => setInputEndpoint(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={intl.formatMessage({
          id: "stash_id_endpoint",
          defaultMessage: "Endpoint (optional)",
        })}
        disabled={disabled}
        className="endpoint-input"
      />
      <InputGroup className="stash-id-input-group">
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
          value={inputStashID}
          onChange={handleStashIDChange}
          onPaste={handlePaste}
          onKeyDown={handleKeyDown}
          placeholder={intl.formatMessage({
            id: "dialogs.stash_id_filter.placeholder",
            defaultMessage: "Stash ID or paste URL",
          })}
          disabled={disabled}
        />
      </InputGroup>
      {urlParsed && (
        <div className="url-parsed-hint">
          ✓ {intl.formatMessage({ 
            id: "dialogs.stash_id_filter.url_parsed", 
            defaultMessage: "URL parsed automatically" 
          })}
        </div>
      )}
      <div className="stash-id-hint">
        {intl.formatMessage({ 
          id: "dialogs.stash_id_filter.hint", 
          defaultMessage: "Tip: Paste a stash-box URL to auto-fill" 
        })}
      </div>
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

export const SidebarStashIDFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = useStashIDFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled =
    state.selectedModifiers.any || state.selectedModifiers.none;

  const stashIDInput = (
    <StashIDInput
      inputEndpoint={state.inputEndpoint}
      setInputEndpoint={state.setInputEndpoint}
      inputStashID={state.inputStashID}
      setInputStashID={state.setInputStashID}
      onSubmit={state.onInputSubmit}
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
      preCandidates={stashIDInput}
    />
  );
};
