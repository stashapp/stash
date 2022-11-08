import React, { useContext } from "react";
import { FormattedMessage } from "react-intl";
import { convertToRatingFormat } from "src/utils/rating";
import { RatingSystem } from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) => {
  const { configuration: config } = useContext(ConfigurationContext);

  return rating ? (
    <div className={config?.interface.ratingSystem == RatingSystem.FiveStar ? `rating-banner rating-${convertToRatingFormat(rating,
      config?.interface.ratingSystem)}` : `rating-banner rating100-${Math.trunc(rating / 5)}`}>
      <FormattedMessage id="rating" />:{" "}
      {convertToRatingFormat(
        rating,
        config?.interface.ratingSystem ?? RatingSystem.FiveStar
      )}
    </div>
  ) : (
    <></>
  );
};
