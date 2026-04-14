import { Star, X } from "lucide-react";
import { useState, type MouseEvent } from "react";
import { Field } from "./EditModal";
import type { RatingStarPrecision, RatingSystemOptions } from "../api/types";
import { useAppConfig } from "../state/AppConfigContext";

export const defaultRatingSystemOptions: RatingSystemOptions = {
  type: "stars",
  starPrecision: "full",
};

export function normalizeRatingSystemType(value: string | null | undefined): RatingSystemOptions["type"] {
  return value?.toLowerCase() === "decimal" ? "decimal" : "stars";
}

export function normalizeRatingStarPrecision(value: string | null | undefined): RatingStarPrecision {
  switch (value?.toLowerCase()) {
    case "half":
      return "half";
    case "quarter":
      return "quarter";
    case "tenth":
      return "tenth";
    case "full":
    default:
      return "full";
  }
}

export function normalizeRatingOptions(options?: Partial<RatingSystemOptions> | null): RatingSystemOptions {
  return {
    type: normalizeRatingSystemType(options?.type),
    starPrecision: normalizeRatingStarPrecision(options?.starPrecision),
  };
}

function roundToStep(value: number, step: number) {
  if (!step) return value;
  return Number((Math.round(value / step) * step).toFixed(2));
}

function clampDisplayRating(value: number, options?: RatingSystemOptions) {
  const maxValue = getRatingMax(options);
  const step = getRatingStep(options);
  return roundToStep(Math.min(maxValue, Math.max(0, value)), step);
}

function trimTrailingZeros(value: number) {
  return value.toFixed(2).replace(/\.00$/, "").replace(/(\.\d)0$/, "$1");
}

export function getRatingPrecision(precision: RatingStarPrecision) {
  switch (precision) {
    case "half":
      return 0.5;
    case "quarter":
      return 0.25;
    case "tenth":
      return 0.1;
    case "full":
    default:
      return 1;
  }
}

export function getRatingStep(options?: RatingSystemOptions) {
  const normalized = normalizeRatingOptions(options);
  return normalized.type === "decimal" ? 0.1 : getRatingPrecision(normalized.starPrecision);
}

export function getRatingMax(options?: RatingSystemOptions) {
  return normalizeRatingOptions(options).type === "decimal" ? 10 : 5;
}

export function convertToRatingFormat(rating: number | null | undefined, options?: RatingSystemOptions) {
  if (!rating) {
    return null;
  }

  const normalized = normalizeRatingOptions(options);
  const maxValue = getRatingMax(normalized);
  const step = getRatingStep(normalized);
  return roundToStep(rating / (100 / maxValue), step);
}

export function convertFromRatingFormat(rating: number, options?: RatingSystemOptions) {
  const normalized = normalizeRatingOptions(options);
  const maxValue = getRatingMax(normalized);
  return Math.round(rating * (100 / maxValue));
}

export function formatDisplayRating(rating: number | null | undefined, options?: RatingSystemOptions) {
  const converted = convertToRatingFormat(rating, options);
  if (converted === null) {
    return null;
  }

  const normalized = normalizeRatingOptions(options);
  return normalized.type === "decimal" ? converted.toFixed(1) : trimTrailingZeros(converted);
}

export function getRatingInputLabel(options?: RatingSystemOptions) {
  const normalized = normalizeRatingOptions(options);
  const maxValue = getRatingMax(normalized);
  const step = getRatingStep(normalized);

  if (normalized.type === "decimal") {
    return "Rating (0-10.0)";
  }

  return step === 1 ? `Rating (0-${maxValue} stars)` : `Rating (0-${maxValue} stars, step ${step})`;
}

function useRatingOptions() {
  const { config } = useAppConfig();
  return normalizeRatingOptions(config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions);
}

function StaticStars({ value, sizeClass }: { value: number; sizeClass: string }) {
  return (
    <span className="flex items-center gap-0.5 text-plex-accent">
      {Array.from({ length: 5 }, (_, index) => {
        const fill = Math.max(0, Math.min(1, value - index));
        return (
          <span key={index} className="relative inline-flex">
            <Star className={`${sizeClass} text-plex-text-muted`} />
            <span className="absolute inset-y-0 left-0 overflow-hidden" style={{ width: `${fill * 100}%` }}>
              <Star className={`${sizeClass} fill-current text-plex-accent`} />
            </span>
          </span>
        );
      })}
    </span>
  );
}

function RatingNumberInput({
  value,
  onChange,
  options,
  inputClassName,
}: {
  value?: number;
  onChange: (value: number | undefined) => void;
  options: RatingSystemOptions;
  inputClassName?: string;
}) {
  const displayValue = convertToRatingFormat(value, options);

  return (
    <input
      type="number"
      value={displayValue ?? ""}
      min={0}
      max={getRatingMax(options)}
      step={getRatingStep(options)}
      onChange={(event) => {
        const nextValue = event.target.value;
        if (!nextValue) {
          onChange(undefined);
          return;
        }

        const parsed = Number(nextValue);
        if (!Number.isFinite(parsed)) {
          return;
        }

        onChange(convertFromRatingFormat(clampDisplayRating(parsed, options), options));
      }}
      onBlur={(event) => {
        if (!event.target.value) {
          return;
        }

        const parsed = Number(event.target.value);
        if (!Number.isFinite(parsed)) {
          onChange(undefined);
          return;
        }

        onChange(convertFromRatingFormat(clampDisplayRating(parsed, options), options));
      }}
      className={inputClassName ?? "w-full rounded border border-gray-700 bg-gray-800 px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"}
    />
  );
}

export function RatingBadge({ rating }: { rating?: number }) {
  const options = useRatingOptions();
  const displayValue = convertToRatingFormat(rating, options);
  const label = formatDisplayRating(rating, options);

  if (displayValue === null || label === null) {
    return null;
  }

  return (
    <span className="flex items-center gap-1 text-xs font-medium text-plex-accent">
      {options.type === "stars" ? <StaticStars value={displayValue} sizeClass="h-3 w-3" /> : <Star className="h-3 w-3 fill-current" />}
      <span>{label}</span>
    </span>
  );
}

/** Diagonal rating banner for grid cards (like Stash's RatingBanner) */
export function RatingBanner({ rating }: { rating?: number }) {
  const options = useRatingOptions();
  const displayValue = convertToRatingFormat(rating, options);
  if (displayValue === null) return null;
  const normalized = normalizeRatingOptions(options);
  const label = formatDisplayRating(rating, options);
  const showStars = normalized.type === "stars" && normalized.starPrecision === "full";
  const displayText = showStars ? "★".repeat(Math.round(displayValue)) : (label ?? "");
  if (!displayText) return null;

  // Color gradient from gray (low) through gold to red (high), matching original Stash
  const bannerColor = getRatingBannerColor(rating, options);

  return (
    <div className="rating-banner-container absolute top-0 left-0 overflow-hidden pointer-events-none" style={{ width: 80, height: 80 }}>
      <div
        className="absolute text-white font-bold text-center tracking-wide shadow-md"
        style={{
          left: -48,
          top: 14,
          transform: "rotate(-36deg)",
          padding: "6px 45px",
          fontSize: "1rem",
          lineHeight: "1.6rem",
          letterSpacing: 1,
          background: bannerColor,
        }}
      >
        {displayText}
      </div>
    </div>
  );
}

export function getRatingBannerColor(rating: number | undefined, options?: RatingSystemOptions): string {
  if (rating == null) return "#939393";
  // Normalize to 0-100 scale
  const normalized = normalizeRatingOptions(options);
  let pct: number;
  if (normalized.type === "stars") {
    const display = convertToRatingFormat(rating, options) ?? 0;
    pct = (display / 5) * 100;
  } else {
    pct = rating;
  }
  // Color gradient matching original Stash: gray → gold → orange → red
  const colors: [number, string][] = [
    [0, "#939393"], [10, "#9b8c7d"], [20, "#9e8974"], [30, "#a7805b"],
    [40, "#af7944"], [50, "#b47435"], [60, "#c39f2b"], [70, "#d2ca20"],
    [80, "#e7a811"], [85, "#ff8000"], [90, "#ff6a07"], [95, "#ff4812"],
    [100, "#ff0000"],
  ];
  for (let i = colors.length - 1; i >= 0; i--) {
    if (pct >= colors[i][0]) return colors[i][1];
  }
  return "#939393";
}

export function RatingField({ value, onChange }: { value?: number; onChange: (value: number | undefined) => void }) {
  const options = useRatingOptions();
  const displayValue = convertToRatingFormat(value, options);

  return (
    <Field label={getRatingInputLabel(options)}>
      <div className="space-y-2">
        <RatingNumberInput value={value} onChange={onChange} options={options} />
        {options.type === "stars" && displayValue !== null && (
          <div className="flex items-center gap-2 text-sm text-plex-text-secondary">
            <StaticStars value={displayValue} sizeClass="h-4 w-4" />
            <span>{formatDisplayRating(value, options)}</span>
          </div>
        )}
      </div>
    </Field>
  );
}

export function InteractiveRating({ value, onChange }: { value?: number; onChange: (value: number) => void }) {
  const options = useRatingOptions();
  const displayValue = convertToRatingFormat(value, options) ?? 0;
  const label = formatDisplayRating(value, options);
  const [hoverValue, setHoverValue] = useState<number | null>(null);

  if (options.type === "stars") {
    const step = getRatingStep(options);
    const activeValue = hoverValue ?? displayValue;

    const getValueFromPointer = (event: MouseEvent<HTMLButtonElement>, star: number) => {
      const rect = event.currentTarget.getBoundingClientRect();
      const ratio = rect.width > 0 ? (event.clientX - rect.left) / rect.width : 1;
      const clampedRatio = Math.min(1, Math.max(0, ratio));
      const segments = Math.max(1, Math.ceil(clampedRatio / step));
      const fractionalValue = Math.min(1, Math.max(step, Number((segments * step).toFixed(2))));
      return clampDisplayRating(star - 1 + fractionalValue, options);
    };

    return (
      <div className="flex items-center gap-2">
        <div className="flex items-center gap-0.5" onMouseLeave={() => setHoverValue(null)}>
          {[1, 2, 3, 4, 5].map((star) => (
            <button
              key={star}
              type="button"
              onMouseMove={(event) => setHoverValue(getValueFromPointer(event, star))}
              onFocus={() => setHoverValue(star)}
              onBlur={() => setHoverValue(null)}
              onClick={(event) => {
                const nextDisplayValue = getValueFromPointer(event, star);
                onChange(nextDisplayValue === displayValue ? 0 : convertFromRatingFormat(nextDisplayValue, options));
              }}
              className="relative text-plex-accent transition-transform hover:scale-110"
              title="Set rating"
            >
              <Star className="h-5 w-5 text-plex-text-muted" />
              <span
                className="absolute inset-y-0 left-0 overflow-hidden"
                style={{ width: `${Math.max(0, Math.min(1, activeValue - (star - 1))) * 100}%` }}
              >
                <Star className="h-5 w-5 fill-current text-plex-accent" />
              </span>
            </button>
          ))}
        </div>
        {(hoverValue ?? label) && <span className="text-sm text-plex-text-secondary">{hoverValue != null ? trimTrailingZeros(hoverValue) : label}</span>}
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-center gap-3">
      <div className="w-24">
        <RatingNumberInput
          value={value}
          onChange={(nextValue) => onChange(nextValue ?? 0)}
          options={options}
          inputClassName="w-full rounded border border-plex-border bg-plex-card px-3 py-2 text-sm text-plex-text focus:outline-none focus:border-plex-accent"
        />
      </div>
      {label && <span className="text-sm text-plex-text-secondary">{label}</span>}
      {!!value && (
        <button
          onClick={() => onChange(0)}
          className="inline-flex items-center gap-1 rounded border border-plex-border px-2 py-1 text-xs text-plex-text-secondary hover:text-plex-text"
        >
          <X className="h-3 w-3" /> Clear
        </button>
      )}
    </div>
  );
}