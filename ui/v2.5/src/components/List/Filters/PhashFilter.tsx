import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { Dropdown, Form, InputGroup } from "react-bootstrap";
import { useIntl, IntlShape } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faFingerprint } from "@fortawesome/free-solid-svg-icons";
import { IPhashDistanceValue } from "../../../models/list-filter/types";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { PhashCriterion } from "../../../models/list-filter/criteria/phash";
import { CriterionModifier } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { Icon } from "src/components/Shared/Icon";
import { cloneDeep } from "lodash-es";

// ============================================================================
// LEGACY EXPORTS FOR BACKWARDS COMPATIBILITY
// ============================================================================

interface IPhashFilterProps {
  criterion: ModifierCriterion<IPhashDistanceValue>;
  onValueChanged: (value: IPhashDistanceValue) => void;
}

export const PhashFilter: React.FC<IPhashFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function valueChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      value: event.target.value,
      distance: criterion.value.distance,
    });
  }

  function distanceChanged(event: React.ChangeEvent<HTMLInputElement>) {
    let distance = parseInt(event.target.value);
    if (distance < 0 || isNaN(distance)) {
      distance = 0;
    }

    onValueChanged({
      distance,
      value: criterion.value.value,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={valueChanged}
          value={value ? value.value : ""}
          placeholder={intl.formatMessage({ id: "media_info.phash" })}
        />
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <Form.Control
              type="number"
              className="btn-secondary"
              onChange={distanceChanged}
              value={value ? value.distance : ""}
              placeholder={intl.formatMessage({ id: "distance" })}
            />
          </Form.Group>
        )}
    </div>
  );
};

// ============================================================================
// NEW IMPROVED SIDEBAR PHASH FILTER
// ============================================================================

// Distance presets with meaningful labels
interface DistancePreset {
  id: string;
  label: string;
  distance: number;
}

const distancePresets: DistancePreset[] = [
  { id: "exact", label: "Exact (0)", distance: 0 },
  { id: "very_similar", label: "Very Similar (≤3)", distance: 3 },
  { id: "similar", label: "Similar (≤8)", distance: 8 },
  { id: "loose", label: "Loose (≤16)", distance: 16 },
];

// Create icon for phash value
function createPhashIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faFingerprint}
      style={{ marginRight: "0.5em", opacity: 0.7 }}
      fixedWidth
    />
  );
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

function usePhashFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [inputValue, setInputValue] = useState("");
  const [inputDistance, setInputDistance] = useState("");
  const [showCustomDistance, setShowCustomDistance] = useState(false);

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as PhashCriterion;

    return filter.makeCriterion(option.type) as PhashCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: PhashCriterion | null) => {
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

  // Determine if there's an active phash value
  const hasActiveValue = useMemo(() => {
    return (
      value?.value &&
      modifier !== CriterionModifier.IsNull &&
      modifier !== CriterionModifier.NotNull
    );
  }, [value, modifier]);

  // Get display label for the current value
  const getValueLabel = useCallback(() => {
    if (!hasActiveValue || !value?.value) return null;

    const phashShort = value.value.length > 16
      ? value.value.slice(0, 8) + "..." + value.value.slice(-8)
      : value.value;

    if (value.distance && value.distance > 0) {
      return `${phashShort} (dist: ${value.distance})`;
    }
    return phashShort;
  }, [hasActiveValue, value]);

  // Get modifier label for display
  const getModifierDisplayLabel = useCallback(() => {
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
      const modifierLabel = getModifierDisplayLabel();
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
        icon: createPhashIcon(),
      });
    }

    return items;
  }, [intl, selectedModifiers, getValueLabel, getModifierDisplayLabel, modifier]);

  // Get active distance preset
  const activeDistancePreset = useMemo((): string | null => {
    if (!hasActiveValue || !value?.distance === undefined) return null;
    
    const preset = distancePresets.find(p => p.distance === value.distance);
    return preset?.id || "custom";
  }, [hasActiveValue, value]);

  // Build candidates list (modifier options + distance presets)
  const candidates = useMemo(() => {
    const items: Option[] = [];

    // Don't show candidates if any/none is selected
    if (selectedModifiers.any || selectedModifiers.none) {
      return items;
    }

    // Show modifier options when nothing is selected
    if (!hasActiveValue) {
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

    // Add distance presets when phash input has a value
    if (inputValue.trim()) {
      distancePresets.forEach((preset) => {
        items.push({
          id: `distance_${preset.id}`,
          label: preset.label,
          className: "preset-option",
          canExclude: false,
        });
      });

      // Add custom distance option
      if (!showCustomDistance) {
        items.push({
          id: "distance_custom",
          label: intl.formatMessage({ id: "actions.custom", defaultMessage: "Custom distance..." }),
          className: "preset-option",
          canExclude: false,
        });
      }
    }

    return items;
  }, [intl, selectedModifiers, hasActiveValue, inputValue, showCustomDistance]);

  const onSelect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Handle modifier selection
        const newCriterion = cloneDeep(criterion);
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = { value: "", distance: 0 };
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = { value: "", distance: 0 };
        }
        setCriterion(newCriterion);
        setShowCustomDistance(false);
      } else if (v.id === "distance_custom") {
        // Show custom distance input
        setShowCustomDistance(true);
      } else if (v.id.startsWith("distance_")) {
        // Handle distance preset selection
        const presetId = v.id.replace("distance_", "");
        const preset = distancePresets.find(p => p.id === presetId);
        if (preset && inputValue.trim()) {
          const newCriterion = cloneDeep(criterion);
          newCriterion.modifier = CriterionModifier.Equals;
          newCriterion.value = { value: inputValue.trim(), distance: preset.distance };
          setCriterion(newCriterion);
          setInputValue("");
          setShowCustomDistance(false);
        }
      }
    },
    [criterion, setCriterion, inputValue]
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
        setInputValue("");
        setInputDistance("");
        setShowCustomDistance(false);
      }
    },
    [setCriterion]
  );

  const onInputSubmit = useCallback(
    (phashValue: string, distance: number, notEquals: boolean) => {
      if (!phashValue.trim()) {
        setCriterion(null);
        return;
      }

      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = notEquals
        ? CriterionModifier.NotEquals
        : CriterionModifier.Equals;
      newCriterion.value = { value: phashValue.trim(), distance };
      setCriterion(newCriterion);
      setInputValue("");
      setInputDistance("");
    },
    [criterion, setCriterion]
  );

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    inputValue,
    setInputValue,
    inputDistance,
    setInputDistance,
    onInputSubmit,
    selectedModifiers,
    hasActiveValue,
    showCustomDistance,
    setShowCustomDistance,
  };
}

// Phash input component - just the phash value input
interface IPhashInputProps {
  inputValue: string;
  setInputValue: (value: string) => void;
  disabled?: boolean;
}

const PhashInput: React.FC<IPhashInputProps> = ({
  inputValue,
  setInputValue,
  disabled,
}) => {
  const intl = useIntl();

  return (
    <div className="phash-input-container">
      <Form.Control
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        placeholder={intl.formatMessage({
          id: "media_info.phash",
          defaultMessage: "Enter PHash value...",
        })}
        disabled={disabled}
        className="phash-text-input"
      />
      {inputValue.trim() && (
        <div className="phash-hint">
          {intl.formatMessage({ 
            id: "dialogs.phash_filter.hint", 
            defaultMessage: "Select a distance threshold below" 
          })}
        </div>
      )}
    </div>
  );
};

// Custom distance input component
interface ICustomDistanceInputProps {
  inputValue: string;
  setInputValue: (value: string) => void;
  inputDistance: string;
  setInputDistance: (value: string) => void;
  onSubmit: (phashValue: string, distance: number, notEquals: boolean) => void;
  onCancel: () => void;
  disabled?: boolean;
}

const CustomDistanceInput: React.FC<ICustomDistanceInputProps> = ({
  inputValue,
  setInputValue,
  inputDistance,
  setInputDistance,
  onSubmit,
  onCancel,
  disabled,
}) => {
  const intl = useIntl();

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && inputValue.trim()) {
      const distance = parseInt(inputDistance, 10) || 0;
      onSubmit(inputValue, distance, false);
    }
    if (e.key === "Escape") {
      onCancel();
    }
  };

  return (
    <div className="custom-distance-input">
      <Form.Control
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        placeholder={intl.formatMessage({
          id: "media_info.phash",
          defaultMessage: "PHash",
        })}
        disabled={disabled}
        className="phash-text-input"
      />
      <Form.Control
        type="number"
        min={0}
        value={inputDistance}
        onChange={(e) => setInputDistance(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={intl.formatMessage({
          id: "distance",
          defaultMessage: "Distance (0-32)",
        })}
        disabled={disabled}
        autoFocus
      />
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

export const SidebarPhashFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = usePhashFilterState({ option, filter, setFilter });

  // Disable input when any/none modifier is selected
  const inputDisabled = state.selectedModifiers.any || state.selectedModifiers.none;

  // Show either the simple phash input or the custom distance input
  const inputSection = state.showCustomDistance ? (
    <CustomDistanceInput
      inputValue={state.inputValue}
      setInputValue={state.setInputValue}
      inputDistance={state.inputDistance}
      setInputDistance={state.setInputDistance}
      onSubmit={state.onInputSubmit}
      onCancel={() => state.setShowCustomDistance(false)}
      disabled={inputDisabled}
    />
  ) : (
    <PhashInput
      inputValue={state.inputValue}
      setInputValue={state.setInputValue}
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
      preCandidates={inputSection}
    />
  );
};
