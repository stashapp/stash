import React, { useMemo, useState } from "react";
import { CriterionModifier } from "../../../core/generated-graphql";
import { CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { DurationCriterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import TextUtils from "src/utils/text";

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

// Duration presets in seconds
const DURATION_PRESETS = [
  { id: "0-5", label: "0-5 min", min: 0, max: 300 },
  { id: "5-10", label: "5-10 min", min: 300, max: 600 },
  { id: "10-20", label: "10-20 min", min: 600, max: 1200 },
  { id: "20-40", label: "20-40 min", min: 1200, max: 2400 },
  { id: "40+", label: "40+ min", min: 2400, max: null },
];

const MAX_DURATION = 7200; // 2 hours in seconds for the slider

export const SidebarDurationFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const criteria = filter.criteriaFor(option.type) as DurationCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  // Get current values from criterion
  const currentMin = criterion?.value?.value ?? 0;
  const currentMax = criterion?.value?.value2 ?? MAX_DURATION;

  const [sliderMin, setSliderMin] = useState(currentMin);
  const [sliderMax, setSliderMax] = useState(currentMax);
  const [minInput, setMinInput] = useState(
    TextUtils.secondsAsTimeString(currentMin)
  );
  const [maxInput, setMaxInput] = useState(
    currentMax >= MAX_DURATION ? "max" : TextUtils.secondsAsTimeString(currentMax)
  );

  // Determine which preset is selected
  const selectedPreset = useMemo(() => {
    if (!criterion) return null;

    // Check if current values match any preset
    for (const preset of DURATION_PRESETS) {
      if (preset.max === null) {
        // For "40+ min" preset
        if (
          criterion.modifier === CriterionModifier.GreaterThan &&
          criterion.value.value === preset.min
        ) {
          return preset.id;
        }
      } else {
        // For range presets
        if (
          criterion.modifier === CriterionModifier.Between &&
          criterion.value.value === preset.min &&
          criterion.value.value2 === preset.max
        ) {
          return preset.id;
        }
      }
    }

    // Check if it's a custom range or custom GreaterThan
    if (
      criterion.modifier === CriterionModifier.Between ||
      criterion.modifier === CriterionModifier.GreaterThan
    ) {
      return "custom";
    }

    return null;
  }, [criterion]);

  const options: Option[] = useMemo(() => {
    return DURATION_PRESETS.map((preset) => ({
      id: preset.id,
      label: preset.label,
      className: "duration-preset",
    }));
  }, []);

  const selected: Option[] = useMemo(() => {
    if (!selectedPreset) return [];
    if (selectedPreset === "custom") return [];

    const preset = DURATION_PRESETS.find((p) => p.id === selectedPreset);
    if (preset) {
      return [
        {
          id: preset.id,
          label: preset.label,
          className: "duration-preset",
        },
      ];
    }
    return [];
  }, [selectedPreset]);

  function onSelectPreset(item: Option) {
    const preset = DURATION_PRESETS.find((p) => p.id === item.id);
    if (!preset) return;

    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (preset.max === null) {
      // "40+ min" - use GreaterThan
      newCriterion.modifier = CriterionModifier.GreaterThan;
      newCriterion.value.value = preset.min;
      newCriterion.value.value2 = undefined;
    } else {
      // Range preset - use Between
      newCriterion.modifier = CriterionModifier.Between;
      newCriterion.value.value = preset.min;
      newCriterion.value.value2 = preset.max;
    }

    setSliderMin(preset.min);
    setSliderMax(preset.max ?? MAX_DURATION);
    setMinInput(TextUtils.secondsAsTimeString(preset.min));
    setMaxInput(
      preset.max === null ? "max" : TextUtils.secondsAsTimeString(preset.max)
    );
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselectPreset() {
    setFilter(filter.removeCriterion(option.type));
    setSliderMin(0);
    setSliderMax(MAX_DURATION);
    setMinInput(TextUtils.secondsAsTimeString(0));
    setMaxInput("max");
  }

  // Parse time input (supports formats like "10", "1:30", "1:30:00")
  function parseTimeInput(input: string): number | null {
    const trimmed = input.trim().toLowerCase();

    if (trimmed === "max") {
      return MAX_DURATION;
    }

    // Try to parse as pure number (minutes)
    const minutesOnly = parseFloat(trimmed);
    if (!isNaN(minutesOnly) && trimmed.indexOf(":") === -1) {
      return Math.round(minutesOnly * 60);
    }

    // Parse HH:MM:SS or MM:SS format
    const parts = trimmed.split(":").map((p) => parseInt(p));
    if (parts.some(isNaN)) {
      return null;
    }

    if (parts.length === 2) {
      // MM:SS
      return parts[0] * 60 + parts[1];
    } else if (parts.length === 3) {
      // HH:MM:SS
      return parts[0] * 3600 + parts[1] * 60 + parts[2];
    }

    return null;
  }

  function handleSliderChange(min: number, max: number) {
    setSliderMin(min);
    setSliderMax(max);
    setMinInput(TextUtils.secondsAsTimeString(min));
    setMaxInput(max >= MAX_DURATION ? "max" : TextUtils.secondsAsTimeString(max));

    // If slider is at full range (0 to max), remove the filter entirely
    if (min === 0 && max >= MAX_DURATION) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    // If max is at MAX_DURATION (but min > 0), use GreaterThan
    if (max >= MAX_DURATION) {
      newCriterion.modifier = CriterionModifier.GreaterThan;
      newCriterion.value.value = min;
      newCriterion.value.value2 = undefined;
    } else {
      newCriterion.modifier = CriterionModifier.Between;
      newCriterion.value.value = min;
      newCriterion.value.value2 = max;
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function handleMinInputChange(value: string) {
    setMinInput(value);
  }

  function handleMaxInputChange(value: string) {
    setMaxInput(value);
  }

  function handleMinInputBlur() {
    const parsed = parseTimeInput(minInput);
    if (parsed !== null && parsed >= 0 && parsed < sliderMax) {
      handleSliderChange(parsed, sliderMax);
    } else {
      // Reset to current value if invalid
      setMinInput(TextUtils.secondsAsTimeString(sliderMin));
    }
  }

  function handleMaxInputBlur() {
    const parsed = parseTimeInput(maxInput);
    if (parsed !== null && parsed > sliderMin && parsed <= MAX_DURATION) {
      handleSliderChange(sliderMin, parsed);
    } else {
      // Reset to current value if invalid
      setMaxInput(
        sliderMax >= MAX_DURATION ? "max" : TextUtils.secondsAsTimeString(sliderMax)
      );
    }
  }

  const customSlider = (
    <div className="duration-slider-container">
      <div className="duration-slider-labels">
        <input
          type="text"
          className="duration-label-input"
          value={minInput}
          onChange={(e) => handleMinInputChange(e.target.value)}
          onBlur={handleMinInputBlur}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.currentTarget.blur();
            }
          }}
          placeholder="0:00"
        />
        <input
          type="text"
          className="duration-label-input"
          value={maxInput}
          onChange={(e) => handleMaxInputChange(e.target.value)}
          onBlur={handleMaxInputBlur}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.currentTarget.blur();
            }
          }}
          placeholder="max"
        />
      </div>
      <div className="duration-slider-inputs">
        <input
          type="range"
          min={0}
          max={MAX_DURATION}
          step={60}
          value={sliderMin}
          onChange={(e) => {
            const newMin = parseInt(e.target.value);
            if (newMin < sliderMax) {
              handleSliderChange(newMin, sliderMax);
            }
          }}
          className="duration-slider duration-slider-min"
        />
        <input
          type="range"
          min={0}
          max={MAX_DURATION}
          step={60}
          value={sliderMax}
          onChange={(e) => {
            const newMax = parseInt(e.target.value);
            if (newMax > sliderMin) {
              handleSliderChange(sliderMin, newMax);
            }
          }}
          className="duration-slider duration-slider-max"
        />
      </div>
    </div>
  );

  return (
    <SidebarListFilter
      title={title}
      candidates={options}
      onSelect={onSelectPreset}
      onUnselect={onUnselectPreset}
      selected={selected}
      singleValue
      preCandidates={selectedPreset === null ? customSlider : undefined}
      preSelected={
        selectedPreset === "custom" || selectedPreset ? customSlider : undefined
      }
      sectionID={sectionID}
    />
  );
};
