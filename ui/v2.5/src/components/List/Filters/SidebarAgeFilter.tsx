import React, { useEffect, useMemo, useState } from "react";
import { CriterionModifier } from "../../../core/generated-graphql";
import { CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { NumberCriterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { DoubleRangeInput } from "src/components/Shared/DoubleRangeInput";
import { useDebounce } from "src/hooks/debounce";

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

// Age presets
const AGE_PRESETS = [
  { id: "18-25", label: "18-25", min: 18, max: 25 },
  { id: "25-35", label: "25-35", min: 25, max: 35 },
  { id: "35-45", label: "35-45", min: 35, max: 45 },
  { id: "45-60", label: "45-60", min: 45, max: 60 },
  { id: "60+", label: "60+", min: 60, max: null },
];

const MAX_AGE = 60; // Maximum age for the slider
const MAX_LABEL = "60+"; // Display label for maximum age

export const SidebarAgeFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const criteria = filter.criteriaFor(option.type) as NumberCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  // Get current values from criterion
  const currentMin = criterion?.value?.value ?? 18;
  const currentMax = criterion?.value?.value2 ?? MAX_AGE;

  const [sliderMin, setSliderMin] = useState(currentMin);
  const [sliderMax, setSliderMax] = useState(currentMax);
  const [minInput, setMinInput] = useState(currentMin.toString());
  const [maxInput, setMaxInput] = useState(
    currentMax >= MAX_AGE ? MAX_LABEL : currentMax.toString()
  );

  // Reset slider when criterion is removed externally (via filter tag X)
  useEffect(() => {
    if (!criterion) {
      setSliderMin(18);
      setSliderMax(MAX_AGE);
      setMinInput("18");
      setMaxInput(MAX_LABEL);
    }
  }, [criterion]);

  // Determine which preset is selected
  const selectedPreset = useMemo(() => {
    if (!criterion) return null;

    // Check if current values match any preset
    for (const preset of AGE_PRESETS) {
      if (preset.max === null) {
        // For "60+" preset
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
    return AGE_PRESETS.map((preset) => ({
      id: preset.id,
      label: preset.label,
      className: "age-preset",
    }));
  }, []);

  const selected: Option[] = useMemo(() => {
    if (!selectedPreset) return [];
    if (selectedPreset === "custom") return [];

    const preset = AGE_PRESETS.find((p) => p.id === selectedPreset);
    if (preset) {
      return [
        {
          id: preset.id,
          label: preset.label,
          className: "age-preset",
        },
      ];
    }
    return [];
  }, [selectedPreset]);

  function onSelectPreset(item: Option) {
    const preset = AGE_PRESETS.find((p) => p.id === item.id);
    if (!preset) return;

    setSliderMin(preset.min);
    setSliderMax(preset.max ?? MAX_AGE);
    setMinInput(preset.min.toString());
    setMaxInput(preset.max === null ? MAX_LABEL : preset.max.toString());

    const currentCriteria = filter.criteriaFor(
      option.type
    ) as NumberCriterion[];
    const currentCriterion =
      currentCriteria.length > 0 ? currentCriteria[0] : null;
    const newCriterion = currentCriterion
      ? currentCriterion.clone()
      : option.makeCriterion();

    if (preset.max === null) {
      // "60+" - use GreaterThan
      newCriterion.modifier = CriterionModifier.GreaterThan;
      newCriterion.value.value = preset.min;
      newCriterion.value.value2 = undefined;
    } else {
      // Range preset - use Between
      newCriterion.modifier = CriterionModifier.Between;
      newCriterion.value.value = preset.min;
      newCriterion.value.value2 = preset.max;
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselectPreset() {
    setSliderMin(18);
    setSliderMax(MAX_AGE);
    setMinInput("18");
    setMaxInput(MAX_LABEL);
    setFilter(filter.removeCriterion(option.type));
  }

  // Parse age input (supports formats like "25", "100+")
  function parseAgeInput(input: string): number | null {
    const trimmed = input.trim().toLowerCase();

    if (trimmed === "max" || trimmed === MAX_LABEL.toLowerCase()) {
      return MAX_AGE;
    }

    const age = parseInt(trimmed);
    if (isNaN(age) || age < 18 || age > MAX_AGE) {
      return null;
    }

    return age;
  }

  // Filter update
  function updateFilter(min: number, max: number) {
    // If slider is at full range (18 to max), remove the filter entirely
    if (min === 18 && max >= MAX_AGE) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    const currentCriteria = filter.criteriaFor(
      option.type
    ) as NumberCriterion[];
    const currentCriterion =
      currentCriteria.length > 0 ? currentCriteria[0] : null;
    const newCriterion = currentCriterion
      ? currentCriterion.clone()
      : option.makeCriterion();

    // If max is at MAX_AGE (but min > 18), use GreaterThan
    if (max >= MAX_AGE) {
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
    setSliderMin(min);
    setSliderMax(max);
    setMinInput(min.toString());
    setMaxInput(max >= MAX_AGE ? MAX_LABEL : max.toString());

    debounceUpdateFilter(min, max);
  }

  function handleMinInputChange(value: string) {
    setMinInput(value);
  }

  function handleMaxInputChange(value: string) {
    setMaxInput(value);
  }

  function handleMinInputBlur() {
    const parsed = parseAgeInput(minInput);
    if (parsed !== null && parsed >= 18 && parsed < sliderMax) {
      handleSliderChange(parsed, sliderMax);
    } else {
      // Reset to current value if invalid
      setMinInput(sliderMin.toString());
    }
  }

  function handleMaxInputBlur() {
    const parsed = parseAgeInput(maxInput);
    if (parsed !== null && parsed > sliderMin && parsed <= MAX_AGE) {
      handleSliderChange(sliderMin, parsed);
    } else {
      // Reset to current value if invalid
      setMaxInput(sliderMax >= MAX_AGE ? MAX_LABEL : sliderMax.toString());
    }
  }

  const customSlider = (
    <div className="age-slider-container">
      <DoubleRangeInput
        min={18}
        max={MAX_AGE}
        value={[sliderMin, sliderMax]}
        onChange={([min, max]) => handleSliderChange(min, max)}
        minInput={
          <input
            type="text"
            className="age-label-input"
            value={minInput}
            onChange={(e) => handleMinInputChange(e.target.value)}
            onBlur={handleMinInputBlur}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.currentTarget.blur();
              }
            }}
            placeholder="18"
          />
        }
        maxInput={
          <input
            type="text"
            className="age-label-input"
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
      />
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
