import React from "react";
import { IUIConfig } from "src/core/config";
import { ConfigurationContext } from "src/hooks/Config";
import {
  defaultRatingStarPrecision,
  defaultRatingSystemOptions,
  RatingSystemType,
} from "src/utils/rating";
import { RatingNumber } from "./RatingNumber";
import { RatingStars } from "./RatingStars";

export interface IRatingSystemProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  valueRequired?: boolean;
}

export const RatingSystem: React.FC<IRatingSystemProps> = (
  props: IRatingSystemProps
) => {
  const { configuration: config } = React.useContext(ConfigurationContext);
  const ratingSystemOptions =
    (config?.ui as IUIConfig)?.ratingSystemOptions ??
    defaultRatingSystemOptions;

  function getRatingStars() {
    return (
      <RatingStars
        value={props.value}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
        precision={
          ratingSystemOptions.starPrecision ?? defaultRatingStarPrecision
        }
        valueRequired={props.valueRequired}
      />
    );
  }

  if (ratingSystemOptions.type === RatingSystemType.Stars) {
    return getRatingStars();
  } else {
    return (
      <RatingNumber
        value={props.value}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
      />
    );
  }
};
