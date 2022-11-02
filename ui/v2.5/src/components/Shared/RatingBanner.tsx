import React, { useContext } from "react";
import { FormattedMessage } from "react-intl";
import { convertToRatingFormat } from "src/components/Scenes/SceneDetails/RatingSystem";
import { RatingSystem } from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) => {
  const { configuration: config } = useContext(ConfigurationContext);

  return rating ? (
    <div className={`rating-banner rating-${Math.trunc(rating / 5)}`}>
      <FormattedMessage id="rating100" />:{" "}
      {convertToRatingFormat(
        rating,
        config?.interface.ratingSystem ?? RatingSystem.FiveStar
      )}
    </div>
  ) : (
    <></>
  );
};
