import React, { useContext } from "react";
import { FormattedMessage } from "react-intl";
import {
  convertToRatingFormat,
  defaultRatingSystemOptions,
  RatingStarPrecision,
  RatingSystemType,
} from "src/utils/rating";
import { ConfigurationContext } from "src/hooks/Config";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) => {
  const { configuration: config } = useContext(ConfigurationContext);
  const ratingSystemOptions =
    config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions;
  const isLegacy =
    ratingSystemOptions.type === RatingSystemType.Stars &&
    ratingSystemOptions.starPrecision === RatingStarPrecision.Full;

  const convertedRating = convertToRatingFormat(
    rating ?? undefined,
    ratingSystemOptions
  );

  return rating ? (
    <div
      className={
        isLegacy
          ? `rating-banner rating-${convertedRating}`
          : `rating-banner rating-100-${Math.trunc(rating / 5)}`
      }
    >
      <FormattedMessage id="rating" />: {convertedRating}
    </div>
  ) : (
    <></>
  );
};
