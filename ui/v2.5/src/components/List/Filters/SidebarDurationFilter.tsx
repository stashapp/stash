import React, { useEffect, useMemo, useState } from "react";
import { CriterionModifier } from "../../../core/generated-graphql";
import { CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { DurationCriterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import TextUtils from "src/utils/text";
import { DoubleRangeInput } from "src/components/Shared/DoubleRangeInput";
import { useDebounce } from "src/hooks/debounce";

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
const MAX_LABEL = "2+ hrs"; // Display label for maximum duration

// Custom step values: 0, 2min (120s), 5min (300s), then 5 minute intervals
const DURATION_STEPS = [
  0, 120, 300, 600, 900, 1200, 1500, 1800, 2100, 2400, 2700, 3000, 3300, 3600,
  3900, 4200, 4500, 4800, 5100, 5400, 5700, 6000, 6300, 6600, 6900, 7200,
];

// Snap a value to the nearest valid step
function snapToStep(value: number): number {
  if (value <= 0) return 0;
  if (value >= MAX_DURATION) return MAX_DURATION;

  // Find the closest step
  let closest = DURATION_STEPS[0];
  let minDiff = Math.abs(value - closest);

  for (const step of DURATION_STEPS) {
    const diff = Math.abs(value - step);
    if (diff < minDiff) {
      minDiff = diff;
      closest = step;
    }
  }

  return closest;
}

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
    currentMin === 0 ? "0m" : TextUtils.secondsAsTimeString(currentMin)
  );
  const [maxInput, setMaxInput] = useState(
    currentMax >= MAX_DURATION
      ? MAX_LABEL
      : TextUtils.secondsAsTimeString(currentMax)
  );

  // Reset slider when criterion is removed externally (via filter tag X)
  useEffect(() => {
    if (!criterion) {
      setSliderMin(0);
      setSliderMax(MAX_DURATION);
      setMinInput("0m");
      setMaxInput(MAX_LABEL);
    }
  }, [criterion]);

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
    setMinInput(
      preset.min === 0 ? "0m" : TextUtils.secondsAsTimeString(preset.min)
    );
    setMaxInput(
      preset.max === null
        ? MAX_LABEL
        : TextUtils.secondsAsTimeString(preset.max)
    );
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselectPreset() {
    setFilter(filter.removeCriterion(option.type));
    setSliderMin(0);
    setSliderMax(MAX_DURATION);
    setMinInput("0m");
    setMaxInput(MAX_LABEL);
  }

  // Parse time input (supports formats like "10", "1:30", "1:30:00", "2+ hrs")
  function parseTimeInput(input: string): number | null {
    const trimmed = input.trim().toLowerCase();

    if (trimmed === "max" || trimmed === MAX_LABEL.toLowerCase()) {
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

  // Debounced filter update
  function updateFilter(min: number, max: number) {
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

  const updateFilterDebounceMS = 300;
  const debounceUpdateFilter = useDebounce(
    updateFilter,
    updateFilterDebounceMS
  );

  function handleSliderChange(min: number, max: number) {
    if (min < 0 || max > MAX_DURATION || min >= max) {
      return;
    }

    setSliderMin(min);
    setSliderMax(max);
    setMinInput(min === 0 ? "0m" : TextUtils.secondsAsTimeString(min));
    setMaxInput(
      max >= MAX_DURATION ? MAX_LABEL : TextUtils.secondsAsTimeString(max)
    );

    debounceUpdateFilter(min, max);
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
      setMinInput(
        sliderMin === 0 ? "0m" : TextUtils.secondsAsTimeString(sliderMin)
      );
    }
  }

  function handleMaxInputBlur() {
    const parsed = parseTimeInput(maxInput);
    if (parsed !== null && parsed > sliderMin && parsed <= MAX_DURATION) {
      handleSliderChange(sliderMin, parsed);
    } else {
      // Reset to current value if invalid
      setMaxInput(
        sliderMax >= MAX_DURATION
          ? MAX_LABEL
          : TextUtils.secondsAsTimeString(sliderMax)
      );
    }
  }

  const customSlider = (
    <DoubleRangeInput
      className="duration-slider"
      minInput={
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
      }
      maxInput={
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
          placeholder={MAX_LABEL}
        />
      }
      min={0}
      max={MAX_DURATION}
      value={[sliderMin, sliderMax]}
      onChange={(vals) => {
        handleSliderChange(snapToStep(vals[0]), snapToStep(vals[1]));
      }}
    />
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
