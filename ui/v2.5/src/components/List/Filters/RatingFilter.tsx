import React, { useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { RatingStars } from "src/components/Shared/Rating/RatingStars";
import {
  defaultRatingStarPrecision,
  defaultRatingSystemOptions,
} from "src/utils/rating";
import { useConfigurationContext } from "src/hooks/Config";
import { RatingCriterion } from "src/models/list-filter/criteria/rating";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

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

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

const any = "any";
const none = "none";

export const SidebarRatingFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();

  const anyLabel = `(${intl.formatMessage({
    id: "criterion_modifier_values.any",
  })})`;
  const noneLabel = `(${intl.formatMessage({
    id: "criterion_modifier_values.none",
  })})`;

  const anyOption = useMemo(
    () => ({
      id: "any",
      label: anyLabel,
      className: "modifier-object",
    }),
    [anyLabel]
  );

  const noneOption = useMemo(
    () => ({
      id: "none",
      label: noneLabel,
      className: "modifier-object",
    }),
    [noneLabel]
  );

  const { configuration: config } = useConfigurationContext();
  const ratingSystemOptions =
    config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions;

  const options: Option[] = useMemo(() => {
    return [anyOption, noneOption];
  }, [anyOption, noneOption]);

  const criteria = filter.criteriaFor(option.type) as RatingCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.modifier === CriterionModifier.NotNull) {
      return [anyOption];
    } else if (criterion.modifier === CriterionModifier.IsNull) {
      return [noneOption];
    }

    return [];
  }, [anyOption, noneOption, criterion]);

  const ratingValue = useMemo(() => {
    if (!criterion || criterion.modifier !== CriterionModifier.GreaterThan) {
      return null;
    }

    return criterion.value.value ?? null;
  }, [criterion]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (item.id === any) {
      newCriterion.modifier = CriterionModifier.NotNull;
      // newCriterion.value
    } else if (item.id === none) {
      newCriterion.modifier = CriterionModifier.IsNull;
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect() {
    setFilter(filter.removeCriterion(option.type));
  }

  function onRatingValueChange(value: number | null) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();
    if (value === null) {
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    newCriterion.modifier = CriterionModifier.GreaterThan;
    newCriterion.value.value = value - 1;

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  const ratingStars = (
    <div className="no-icon-margin">
      <RatingStars
        value={ratingValue}
        onSetRating={onRatingValueChange}
        precision={
          ratingSystemOptions.starPrecision ?? defaultRatingStarPrecision
        }
        orMore
      />
    </div>
  );
  return (
    <>
      <SidebarListFilter
        title={title}
        candidates={options}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selected}
        singleValue
        preCandidates={ratingValue === null ? ratingStars : undefined}
        preSelected={ratingValue !== null ? ratingStars : undefined}
        sectionID={sectionID}
      />
      <div></div>
    </>
  );
};
