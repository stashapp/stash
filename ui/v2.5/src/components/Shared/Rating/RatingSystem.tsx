import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import { RatingNumber } from "./RatingNumber";
import { RatingStars } from "./RatingStars";

export interface IRatingSystemProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
}

export const RatingSystem: React.FC<IRatingSystemProps> = (
  props: IRatingSystemProps
) => {
  const { configuration: config } = React.useContext(ConfigurationContext);

  function getRatingStars() {
    return (
      <RatingStars
        value={props.value}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
        ratingSystem={
          config?.interface.ratingSystem ?? GQL.RatingSystem.FiveStar
        }
      />
    );
  }

  let toReturn;
  switch (config?.interface?.ratingSystem) {
    // case GQL.RatingSystem.TenStar:
    // case GQL.RatingSystem.TenPointFiveStar:
    // case GQL.RatingSystem.TenPointTwoFiveStar:
    case GQL.RatingSystem.FiveStar:
    case GQL.RatingSystem.FivePointFiveStar:
    case GQL.RatingSystem.FivePointTwoFiveStar:
      toReturn = getRatingStars();
      break;
    case GQL.RatingSystem.TenPointDecimal:
      toReturn = (
        <RatingNumber
          value={props.value}
          onSetRating={props.onSetRating}
          disabled={props.disabled}
        />
      );
      break;
    default:
      toReturn = getRatingStars();
      break;
  }

  return toReturn;
};
