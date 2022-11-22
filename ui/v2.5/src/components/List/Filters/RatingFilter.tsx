import React from "react";
import { FormattedMessage } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import { INumberValue } from "../../../models/list-filter/types";
import { Criterion } from "../../../models/list-filter/criteria/criterion";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";

interface IRatingFilterProps {
  criterion: Criterion<INumberValue>;
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
