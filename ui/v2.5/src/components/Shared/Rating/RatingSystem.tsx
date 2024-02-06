import React from "react";
import { ConfigurationContext } from "src/hooks/Config";
import {
  defaultRatingStarPrecision,
  defaultRatingSystemOptions,
  RatingSystemType,
} from "src/utils/rating";
import { RatingNumber } from "./RatingNumber";
import { RatingStars } from "./RatingStars";

export interface IRatingSystemProps {
  value: number | null | undefined;
  onSetRating?: (value: number | null) => void;
  disabled?: boolean;
  valueRequired?: boolean;
}

export const RatingSystem: React.FC<IRatingSystemProps> = (
  props: IRatingSystemProps
) => {
  const { configuration: config } = React.useContext(ConfigurationContext);
  const ratingSystemOptions =
    config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions;

  if (ratingSystemOptions.type === RatingSystemType.Stars) {
    return (
      <RatingStars
        value={props.value ?? null}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
        precision={
          ratingSystemOptions.starPrecision ?? defaultRatingStarPrecision
        }
        valueRequired={props.valueRequired}
      />
    );
  } else {
    return (
      <RatingNumber
        value={props.value ?? null}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
      />
    );
  }
};
