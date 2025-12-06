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
// PRESET CONFIGURATIONS FOR NUMBER FILTERS
// ============================================================================

interface NumberPreset {
  id: string;
  label: string;
  modifier: CriterionModifier;
  value: number;
  value2?: number;
}

// Preset configurations for different filter types
const presetConfigs: Record<string, NumberPreset[]> = {
  // =========================================================================
  // COUNT PRESETS (for performers, tags, studios, etc.)
  // =========================================================================
  scene_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-10", label: "1-10", modifier: CriterionModifier.Between, value: 1, value2: 10 },
    { id: "11-50", label: "11-50", modifier: CriterionModifier.Between, value: 11, value2: 50 },
    { id: "51-100", label: "51-100", modifier: CriterionModifier.Between, value: 51, value2: 100 },
    { id: "101-500", label: "101-500", modifier: CriterionModifier.Between, value: 101, value2: 500 },
    { id: "500+", label: "500+", modifier: CriterionModifier.GreaterThan, value: 500 },
  ],
  image_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-100", label: "1-100", modifier: CriterionModifier.Between, value: 1, value2: 100 },
    { id: "101-1000", label: "101-1K", modifier: CriterionModifier.Between, value: 101, value2: 1000 },
    { id: "1001-10000", label: "1K-10K", modifier: CriterionModifier.Between, value: 1001, value2: 10000 },
    { id: "10001-50000", label: "10K-50K", modifier: CriterionModifier.Between, value: 10001, value2: 50000 },
    { id: "50000+", label: "50K+", modifier: CriterionModifier.GreaterThan, value: 50000 },
  ],
  gallery_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-10", label: "1-10", modifier: CriterionModifier.Between, value: 1, value2: 10 },
    { id: "11-50", label: "11-50", modifier: CriterionModifier.Between, value: 11, value2: 50 },
    { id: "51-100", label: "51-100", modifier: CriterionModifier.Between, value: 51, value2: 100 },
    { id: "101-500", label: "101-500", modifier: CriterionModifier.Between, value: 101, value2: 500 },
    { id: "500+", label: "500+", modifier: CriterionModifier.GreaterThan, value: 500 },
  ],
  performer_count: [
    { id: "0", label: "Solo", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1", label: "1", modifier: CriterionModifier.Equals, value: 1 },
    { id: "2", label: "2", modifier: CriterionModifier.Equals, value: 2 },
    { id: "3-5", label: "3-5", modifier: CriterionModifier.Between, value: 3, value2: 5 },
    { id: "6+", label: "6+", modifier: CriterionModifier.GreaterThan, value: 5 },
  ],
  tag_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-5", label: "1-5", modifier: CriterionModifier.Between, value: 1, value2: 5 },
    { id: "6-10", label: "6-10", modifier: CriterionModifier.Between, value: 6, value2: 10 },
    { id: "11-20", label: "11-20", modifier: CriterionModifier.Between, value: 11, value2: 20 },
    { id: "20+", label: "20+", modifier: CriterionModifier.GreaterThan, value: 20 },
  ],
  studio_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-5", label: "1-5", modifier: CriterionModifier.Between, value: 1, value2: 5 },
    { id: "6-10", label: "6-10", modifier: CriterionModifier.Between, value: 6, value2: 10 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 10 },
  ],
  group_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-5", label: "1-5", modifier: CriterionModifier.Between, value: 1, value2: 5 },
    { id: "6-10", label: "6-10", modifier: CriterionModifier.Between, value: 6, value2: 10 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 10 },
  ],
  marker_count: [
    { id: "0", label: "None", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-5", label: "1-5", modifier: CriterionModifier.Between, value: 1, value2: 5 },
    { id: "6-10", label: "6-10", modifier: CriterionModifier.Between, value: 6, value2: 10 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 10 },
  ],
  file_count: [
    { id: "1", label: "Single", modifier: CriterionModifier.Equals, value: 1 },
    { id: "2+", label: "Multiple", modifier: CriterionModifier.GreaterThan, value: 1 },
    { id: "5+", label: "5+", modifier: CriterionModifier.GreaterThan, value: 4 },
  ],
  // Group-specific counts
  containing_group_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1", label: "1", modifier: CriterionModifier.Equals, value: 1 },
    { id: "2+", label: "2+", modifier: CriterionModifier.GreaterThan, value: 1 },
    { id: "5+", label: "5+", modifier: CriterionModifier.GreaterThan, value: 4 },
  ],
  sub_group_count: [
    { id: "0", label: "0", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1-5", label: "1-5", modifier: CriterionModifier.Between, value: 1, value2: 5 },
    { id: "6-10", label: "6-10", modifier: CriterionModifier.Between, value: 6, value2: 10 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 10 },
  ],

  // =========================================================================
  // PLAY/ACTIVITY PRESETS
  // =========================================================================
  play_count: [
    { id: "unplayed", label: "Unplayed", modifier: CriterionModifier.Equals, value: 0 },
    { id: "watched", label: "Watched", modifier: CriterionModifier.GreaterThan, value: 0 },
    { id: "2+", label: "2+", modifier: CriterionModifier.GreaterThan, value: 1 },
    { id: "5+", label: "5+", modifier: CriterionModifier.GreaterThan, value: 4 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 9 },
  ],
  o_counter: [
    { id: "none", label: "None", modifier: CriterionModifier.Equals, value: 0 },
    { id: "1+", label: "1+", modifier: CriterionModifier.GreaterThan, value: 0 },
    { id: "5+", label: "5+", modifier: CriterionModifier.GreaterThan, value: 4 },
    { id: "10+", label: "10+", modifier: CriterionModifier.GreaterThan, value: 9 },
  ],

  // =========================================================================
  // RATING PRESETS (1-100 scale)
  // =========================================================================
  rating100: [
    { id: "1star", label: "★", modifier: CriterionModifier.Between, value: 1, value2: 29 },
    { id: "2star", label: "★★", modifier: CriterionModifier.Between, value: 30, value2: 49 },
    { id: "3star", label: "★★★", modifier: CriterionModifier.Between, value: 50, value2: 69 },
    { id: "4star", label: "★★★★", modifier: CriterionModifier.Between, value: 70, value2: 89 },
    { id: "5star", label: "★★★★★", modifier: CriterionModifier.Between, value: 90, value2: 100 },
  ],

  // =========================================================================
  // VIDEO TECHNICAL PRESETS
  // =========================================================================
  framerate: [
    { id: "24", label: "24 fps", modifier: CriterionModifier.Equals, value: 24 },
    { id: "25", label: "25 fps", modifier: CriterionModifier.Equals, value: 25 },
    { id: "30", label: "30 fps", modifier: CriterionModifier.Equals, value: 30 },
    { id: "50", label: "50 fps", modifier: CriterionModifier.Equals, value: 50 },
    { id: "60", label: "60 fps", modifier: CriterionModifier.Equals, value: 60 },
    { id: "60+", label: "60+ fps", modifier: CriterionModifier.GreaterThan, value: 60 },
  ],
  // Bitrate in bits per second
  bitrate: [
    { id: "low", label: "Low (<5 Mbps)", modifier: CriterionModifier.LessThan, value: 5000000 },
    { id: "medium", label: "Medium (5-15 Mbps)", modifier: CriterionModifier.Between, value: 5000000, value2: 15000000 },
    { id: "high", label: "High (15-30 Mbps)", modifier: CriterionModifier.Between, value: 15000000, value2: 30000000 },
    { id: "ultra", label: "Ultra (30+ Mbps)", modifier: CriterionModifier.GreaterThan, value: 30000000 },
  ],
  interactive_speed: [
    { id: "slow", label: "Slow (<50)", modifier: CriterionModifier.LessThan, value: 50 },
    { id: "medium", label: "Medium (50-150)", modifier: CriterionModifier.Between, value: 50, value2: 150 },
    { id: "fast", label: "Fast (150-250)", modifier: CriterionModifier.Between, value: 150, value2: 250 },
    { id: "veryfast", label: "Very Fast (250+)", modifier: CriterionModifier.GreaterThan, value: 250 },
  ],

  // =========================================================================
  // PERFORMER PHYSICAL PRESETS
  // =========================================================================
  // Height in cm
  height_cm: [
    { id: "petite", label: "Petite (<155 cm)", modifier: CriterionModifier.LessThan, value: 155 },
    { id: "short", label: "Short (155-165 cm)", modifier: CriterionModifier.Between, value: 155, value2: 165 },
    { id: "average", label: "Average (165-175 cm)", modifier: CriterionModifier.Between, value: 165, value2: 175 },
    { id: "tall", label: "Tall (175-185 cm)", modifier: CriterionModifier.Between, value: 175, value2: 185 },
    { id: "verytall", label: "Very Tall (185+ cm)", modifier: CriterionModifier.GreaterThan, value: 185 },
  ],
  // Weight in kg
  weight: [
    { id: "light", label: "<50 kg", modifier: CriterionModifier.LessThan, value: 50 },
    { id: "50-60", label: "50-60 kg", modifier: CriterionModifier.Between, value: 50, value2: 60 },
    { id: "60-70", label: "60-70 kg", modifier: CriterionModifier.Between, value: 60, value2: 70 },
    { id: "70-80", label: "70-80 kg", modifier: CriterionModifier.Between, value: 70, value2: 80 },
    { id: "80-90", label: "80-90 kg", modifier: CriterionModifier.Between, value: 80, value2: 90 },
    { id: "90+", label: "90+ kg", modifier: CriterionModifier.GreaterThan, value: 90 },
  ],
  penis_length: [
    { id: "small", label: "<12 cm", modifier: CriterionModifier.LessThan, value: 12 },
    { id: "average", label: "12-15 cm", modifier: CriterionModifier.Between, value: 12, value2: 15 },
    { id: "large", label: "15-18 cm", modifier: CriterionModifier.Between, value: 15, value2: 18 },
    { id: "xlarge", label: "18-21 cm", modifier: CriterionModifier.Between, value: 18, value2: 21 },
    { id: "xxlarge", label: "21+ cm", modifier: CriterionModifier.GreaterThan, value: 21 },
  ],

  // =========================================================================
  // AGE PRESETS
  // =========================================================================
  age: [
    { id: "18-19", label: "18-19", modifier: CriterionModifier.Between, value: 18, value2: 19 },
    { id: "20s", label: "20s", modifier: CriterionModifier.Between, value: 20, value2: 29 },
    { id: "30s", label: "30s", modifier: CriterionModifier.Between, value: 30, value2: 39 },
    { id: "40s", label: "40s", modifier: CriterionModifier.Between, value: 40, value2: 49 },
    { id: "50s", label: "50s", modifier: CriterionModifier.Between, value: 50, value2: 59 },
    { id: "60+", label: "60+", modifier: CriterionModifier.GreaterThan, value: 59 },
  ],
  performer_age: [
    { id: "18-19", label: "18-19", modifier: CriterionModifier.Between, value: 18, value2: 19 },
    { id: "20s", label: "20s", modifier: CriterionModifier.Between, value: 20, value2: 29 },
    { id: "30s", label: "30s", modifier: CriterionModifier.Between, value: 30, value2: 39 },
    { id: "40s", label: "40s", modifier: CriterionModifier.Between, value: 40, value2: 49 },
    { id: "50+", label: "50+", modifier: CriterionModifier.GreaterThan, value: 49 },
  ],
  birth_year: [
    { id: "2000s", label: "2000s", modifier: CriterionModifier.Between, value: 2000, value2: 2009 },
    { id: "1990s", label: "1990s", modifier: CriterionModifier.Between, value: 1990, value2: 1999 },
    { id: "1980s", label: "1980s", modifier: CriterionModifier.Between, value: 1980, value2: 1989 },
    { id: "1970s", label: "1970s", modifier: CriterionModifier.Between, value: 1970, value2: 1979 },
    { id: "1960s", label: "1960s", modifier: CriterionModifier.Between, value: 1960, value2: 1969 },
    { id: "pre1960", label: "Pre-1960", modifier: CriterionModifier.LessThan, value: 1960 },
  ],
  death_year: [
    { id: "2020s", label: "2020s", modifier: CriterionModifier.Between, value: 2020, value2: 2029 },
    { id: "2010s", label: "2010s", modifier: CriterionModifier.Between, value: 2010, value2: 2019 },
    { id: "2000s", label: "2000s", modifier: CriterionModifier.Between, value: 2000, value2: 2009 },
    { id: "1990s", label: "1990s", modifier: CriterionModifier.Between, value: 1990, value2: 1999 },
    { id: "pre1990", label: "Pre-1990", modifier: CriterionModifier.LessThan, value: 1990 },
  ],
};

// Get presets for a filter type
function getPresets(filterType: string): NumberPreset[] {
  // Check for exact match first
  if (presetConfigs[filterType]) {
    return presetConfigs[filterType];
  }

  // Check for partial matches (e.g., "scene_count" in "performer_scene_count")
  const partialMatches: [string, string[]][] = [
    ["scene_count", ["scene_count"]],
    ["image_count", ["image_count"]],
    ["gallery_count", ["gallery_count"]],
    ["performer_count", ["performer_count"]],
    ["tag_count", ["tag_count"]],
    ["studio_count", ["studio_count"]],
    ["group_count", ["group_count", "sub_group", "containing_group"]],
    ["marker_count", ["marker_count"]],
    ["file_count", ["file_count"]],
    ["play_count", ["play_count"]],
    ["o_counter", ["o_counter", "o_count"]],
    ["rating100", ["rating"]],
    ["framerate", ["framerate", "frame_rate"]],
    ["bitrate", ["bitrate", "bit_rate"]],
    ["interactive_speed", ["interactive_speed"]],
    ["height_cm", ["height"]],
    ["weight", ["weight"]],
    ["penis_length", ["penis_length"]],
    ["age", ["age"]],
    ["performer_age", ["performer_age"]],
    ["birth_year", ["birth_year"]],
    ["death_year", ["death_year"]],
  ];

  for (const [presetKey, patterns] of partialMatches) {
    for (const pattern of patterns) {
      if (filterType.includes(pattern)) {
        return presetConfigs[presetKey];
      }
    }
  }

  return [];
}

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

  // Get presets for this filter type
  const presets = useMemo(() => getPresets(option.type), [option.type]);

  // Check if a preset matches the current filter state
  const matchesPreset = useCallback(
    (preset: NumberPreset): boolean => {
      if (modifier !== preset.modifier) return false;
      if (value.value !== preset.value) return false;
      if (preset.value2 !== undefined && value.value2 !== preset.value2) return false;
      return true;
    },
    [modifier, value]
  );

  // Check if current value matches any preset
  const activePreset = useMemo((): string | null => {
    if (!hasActiveValue) return null;
    
    for (const preset of presets) {
      if (matchesPreset(preset)) {
        return preset.id;
      }
    }
    
    // Has a value but doesn't match any preset = custom
    return "custom";
  }, [hasActiveValue, presets, matchesPreset]);

  // Build candidates list (modifier options + presets + custom)
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

    // Add presets as quick options (excluding the active one)
    if (presets.length > 0) {
      presets
        .filter((preset) => preset.id !== activePreset)
        .forEach((preset) => {
          items.push({
            id: `preset_${preset.id}`,
            label: preset.label,
            className: "preset-option",
            canExclude: false,
          });
        });
    }

    // Add custom option if not already in custom mode
    if (activePreset !== "custom") {
      items.push({
        id: "custom",
        label: intl.formatMessage({ id: "actions.custom", defaultMessage: "Custom..." }),
        className: "preset-option",
        canExclude: false,
      });
    }

    return items;
  }, [intl, selectedModifiers, hasActiveValue, supportsIsNull, supportsNotNull, presets, activePreset]);

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
      } else if (v.id === "custom") {
        // Show custom input
        setShowCustom(true);
      } else if (v.className === "preset-option" && v.id.startsWith("preset_")) {
        // Handle preset selection
        const presetId = v.id.replace("preset_", "");
        const preset = presets.find((p) => p.id === presetId);
        if (preset) {
          const newCriterion = cloneDeep(criterion);
          newCriterion.modifier = preset.modifier;
          newCriterion.value = { value: preset.value, value2: preset.value2 };
          setCriterion(newCriterion);
          setShowCustom(false);
        }
      }
    },
    [criterion, setCriterion, presets]
  );

  const onUnselect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.id === "any" || v.id === "none") {
        setCriterion(null);
      } else if (v.id === "value") {
        setCriterion(null);
        setInputValue("");
        setInputValue2("");
        setShowCustom(false);
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
    showCustom,
    setShowCustom,
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

  // Handle submit and close the custom input
  const handleCustomSubmit = useCallback(
    (value1: string, value2: string, modifier: CriterionModifier) => {
      state.onInputSubmit(value1, value2, modifier);
      state.setShowCustom(false);
    },
    [state]
  );

  // Custom input shown when "Custom..." is selected
  const customInput = state.showCustom ? (
    <div className="custom-number-input">
      <NumberInput
        inputValue={state.inputValue}
        setInputValue={state.setInputValue}
        inputValue2={state.inputValue2}
        setInputValue2={state.setInputValue2}
        onSubmit={handleCustomSubmit}
        modifiers={state.numberModifiers}
        defaultModifier={state.defaultModifier}
        placeholder={placeholder}
        disabled={inputDisabled}
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
