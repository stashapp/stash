import React from "react";

export const DoubleRangeInput: React.FC<{
  className?: string;
  minInput: React.ReactNode;
  maxInput: React.ReactNode;
  min?: number;
  max: number;
  value: [number, number];
  onChange(value: [number, number]): void;
}> = ({
  className = "",
  minInput,
  maxInput,
  min = 0,
  max,
  value,
  onChange,
}) => {
  const minValue = value[0];
  const maxValue = value[1];

  return (
    <div className={`double-range-input ${className}`}>
      <div className="double-range-input-labels">
        {minInput}
        {maxInput}
      </div>
      <div className="double-range-sliders">
        <input
          type="range"
          min={min}
          max={max}
          step={1}
          value={minValue}
          onChange={(e) => {
            const rawValue = parseInt(e.target.value);
            if (rawValue < maxValue) {
              onChange([rawValue, maxValue]);
            }
          }}
          className="double-range-slider double-range-slider-min"
        />
        <input
          type="range"
          min={min}
          max={max}
          step={1}
          value={maxValue}
          onChange={(e) => {
            const rawValue = parseInt(e.target.value);
            if (rawValue > minValue) {
              onChange([minValue, rawValue]);
            }
          }}
          className="double-range-slider double-range-slider-max"
        />
      </div>
    </div>
  );
};
