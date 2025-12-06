import React, { ReactNode, useCallback, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faStar } from "@fortawesome/free-solid-svg-icons";
import { CriterionModifier } from "src/core/generated-graphql";
import { INumberValue } from "src/models/list-filter/types";
import {
  CriterionOption,
  ModifierCriterion,
} from "src/models/list-filter/criteria/criterion";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { RatingStars } from "src/components/Shared/Rating/RatingStars";
import {
  convertToRatingFormat,
  defaultRatingStarPrecision,
  defaultRatingSystemOptions,
  RatingSystemType,
} from "src/utils/rating";
import { ConfigurationContext } from "src/hooks/Config";
import { RatingCriterion } from "src/models/list-filter/criteria/rating";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

// ============================================================================
// LEGACY EXPORTS FOR BACKWARDS COMPATIBILITY
// ============================================================================

interface IRatingFilterProps {
  criterion: ModifierCriterion<INumberValue>;
  onValueChanged: (value: INumberValue) => void;
}

export const RatingFilter: React.FC<IRatingFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  function getRatingSystem(field: "value" | "value2") {
    const defaultValue = field === "value" ? 0 : undefined;

    return (
      <div>
        <RatingSystem
          value={criterion.value[field]}
          onSetRating={(value) => {
            onValueChanged({
              ...criterion.value,
              [field]: value ?? defaultValue,
            });
          }}
          valueRequired
        />
      </div>
    );
  }

  if (
    criterion.modifier === CriterionModifier.Equals ||
    criterion.modifier === CriterionModifier.NotEquals ||
    criterion.modifier === CriterionModifier.GreaterThan ||
    criterion.modifier === CriterionModifier.LessThan
  ) {
    return getRatingSystem("value");
  }

  if (
    criterion.modifier === CriterionModifier.Between ||
    criterion.modifier === CriterionModifier.NotBetween
  ) {
    return (
      <div className="rating-filter">
        {getRatingSystem("value")}
        <span className="and-divider">
          <FormattedMessage id="between_and" />
        </span>
        {getRatingSystem("value2")}
      </div>
    );
  }

  return <></>;
};

// ============================================================================
// NEW IMPROVED SIDEBAR RATING FILTER
// ============================================================================

// Format rating value for display
function formatRatingValue(value: number | undefined, precision: string): string {
  if (value === undefined) return "";
  
  const rating = convertToRatingFormat(value, {
    type: RatingSystemType.Stars,
    starPrecision: precision as any,
  });
  
  return rating?.toString() ?? value.toString();
}

// Create icon for rating value
function createRatingIcon(): React.ReactNode {
  return (
    <FontAwesomeIcon
      icon={faStar}
      style={{ marginRight: "0.5em", color: "#f5c518", opacity: 0.7 }}
      fixedWidth
    />
  );
}

function useRatingFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const { configuration: config } = React.useContext(ConfigurationContext);
  const ratingSystemOptions =
    config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions;
  const starPrecision = ratingSystemOptions.starPrecision ?? defaultRatingStarPrecision;

  // Track pending rating (selected but no modifier chosen yet)
  const [pendingRating, setPendingRating] = useState<number | null>(null);

  const criteria = filter.criteriaFor(option.type) as RatingCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const setCriterion = useCallback(
    (c: RatingCriterion | null) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );

      if (c && c.isValid()) newCriteria.push(c);

      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const modifier = criterion?.modifier;
  const value = criterion?.value;

  // Get modifier label for display
  const getModifierLabel = useCallback(
    (mod: CriterionModifier) => {
      switch (mod) {
        case CriterionModifier.Equals:
          return intl.formatMessage({
            id: "criterion_modifier.equals",
            defaultMessage: "is",
          });
        case CriterionModifier.NotEquals:
          return intl.formatMessage({
            id: "criterion_modifier.not_equals",
            defaultMessage: "is not",
          });
        case CriterionModifier.GreaterThan:
          return intl.formatMessage({
            id: "criterion_modifier.greater_than",
            defaultMessage: "greater than",
          });
        case CriterionModifier.LessThan:
          return intl.formatMessage({
            id: "criterion_modifier.less_than",
            defaultMessage: "less than",
          });
        case CriterionModifier.NotNull:
          return intl.formatMessage({
            id: "criterion_modifier_values.any",
            defaultMessage: "any",
          });
        case CriterionModifier.IsNull:
          return intl.formatMessage({
            id: "criterion_modifier_values.none",
            defaultMessage: "none",
          });
        default:
          return "";
      }
    },
    [intl]
  );

  // Build selected items list
  const selected = useMemo(() => {
    const selectedItems: Option[] = [];

    // Check for any/none modifiers first
    if (modifier === CriterionModifier.NotNull) {
      selectedItems.push({
        id: "any",
        label: `(${getModifierLabel(CriterionModifier.NotNull)})`,
        className: "modifier-object",
      });
      return selectedItems;
    }
    if (modifier === CriterionModifier.IsNull) {
      selectedItems.push({
        id: "none",
        label: `(${getModifierLabel(CriterionModifier.IsNull)})`,
        className: "modifier-object",
      });
      return selectedItems;
    }

    // If filter is active (has value in criterion), show modifier and value
    if (value?.value !== undefined) {
      // Add the modifier indicator
      selectedItems.push({
        id: "modifier",
        label: `(${getModifierLabel(modifier!)})`,
        className: "modifier-object",
      });
      // Add the rating value
      const ratingDisplay = formatRatingValue(value.value, starPrecision);
      selectedItems.push({
        id: "rating",
        label: `${ratingDisplay} ★`,
        icon: createRatingIcon(),
      });
    }
    // If there's a pending rating (selected but no modifier yet), show it
    else if (pendingRating !== null) {
      const ratingDisplay = formatRatingValue(pendingRating, starPrecision);
      selectedItems.push({
        id: "pending",
        label: `${ratingDisplay} ★`,
        icon: createRatingIcon(),
      });
    }

    return selectedItems;
  }, [value, modifier, getModifierLabel, pendingRating, starPrecision]);

  // Build candidates list
  const candidates = useMemo(() => {
    // If a rating is pending (selected but no modifier yet), show modifier options
    if (pendingRating !== null) {
      return [
        {
          id: "equals",
          label: `(${getModifierLabel(CriterionModifier.Equals)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "not_equals",
          label: `(${getModifierLabel(CriterionModifier.NotEquals)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "greater_than",
          label: `(${getModifierLabel(CriterionModifier.GreaterThan)})`,
          className: "modifier-object",
          canExclude: false,
        },
        {
          id: "less_than",
          label: `(${getModifierLabel(CriterionModifier.LessThan)})`,
          className: "modifier-object",
          canExclude: false,
        },
      ];
    }

    // If filter is already active, don't show any candidates
    if (value?.value !== undefined || modifier === CriterionModifier.NotNull || modifier === CriterionModifier.IsNull) {
      return [];
    }

    // Show any/none options when nothing is selected
    return [
      {
        id: "any",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.any",
        })})`,
        className: "modifier-object",
        canExclude: false,
      },
      {
        id: "none",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.none",
        })})`,
        className: "modifier-object",
        canExclude: false,
      },
    ];
  }, [value, modifier, getModifierLabel, pendingRating, intl]);

  const onSelect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object") {
        // Handle any/none selection
        if (v.id === "any") {
          const newCriterion = criterion
            ? (criterion.clone() as RatingCriterion)
            : (option.makeCriterion() as RatingCriterion);
          newCriterion.modifier = CriterionModifier.NotNull;
          newCriterion.value = { value: undefined, value2: undefined };
          setCriterion(newCriterion);
          setPendingRating(null);
          return;
        }
        if (v.id === "none") {
          const newCriterion = criterion
            ? (criterion.clone() as RatingCriterion)
            : (option.makeCriterion() as RatingCriterion);
          newCriterion.modifier = CriterionModifier.IsNull;
          newCriterion.value = { value: undefined, value2: undefined };
          setCriterion(newCriterion);
          setPendingRating(null);
          return;
        }

        // User selected a modifier after choosing a rating
        if (pendingRating !== null) {
          const newCriterion = criterion
            ? (criterion.clone() as RatingCriterion)
            : (option.makeCriterion() as RatingCriterion);
          newCriterion.value = { value: pendingRating, value2: undefined };

          let mod = CriterionModifier.Equals;
          switch (v.id) {
            case "equals":
              mod = CriterionModifier.Equals;
              break;
            case "not_equals":
              mod = CriterionModifier.NotEquals;
              break;
            case "greater_than":
              mod = CriterionModifier.GreaterThan;
              break;
            case "less_than":
              mod = CriterionModifier.LessThan;
              break;
          }
          newCriterion.modifier = mod;
          setCriterion(newCriterion);
          setPendingRating(null);
        }
      }
    },
    [criterion, option, setCriterion, pendingRating]
  );

  const onUnselect = useCallback(
    (v: Option, _exclude: boolean) => {
      if (v.className === "modifier-object" || v.id === "rating" || v.id === "pending") {
        // Clear the filter
        setCriterion(null);
        setPendingRating(null);
      }
    },
    [setCriterion]
  );

  const onRatingSelect = useCallback(
    (ratingValue: number | null) => {
      if (ratingValue === null) {
        setPendingRating(null);
        return;
      }
      // Store as pending - wait for modifier selection
      setPendingRating(ratingValue);
    },
    []
  );

  return {
    selected,
    candidates,
    onSelect,
    onUnselect,
    onRatingSelect,
    starPrecision,
    pendingRating,
    hasActiveFilter: value?.value !== undefined || modifier === CriterionModifier.NotNull || modifier === CriterionModifier.IsNull,
  };
}

interface ISidebarFilter {
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarRatingFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const state = useRatingFilterState({ option, filter, setFilter });

  // Show rating stars input when nothing is selected and no pending rating
  const showRatingStars = !state.hasActiveFilter && state.pendingRating === null;

  const ratingStarsInput = showRatingStars ? (
    <div className="rating-stars-input">
      <RatingStars
        value={null}
        onSetRating={state.onRatingSelect}
        precision={state.starPrecision as any}
      />
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
      preCandidates={ratingStarsInput}
    />
  );
};
