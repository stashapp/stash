import React from "react";
import { FormattedMessage } from "react-intl";
import { ConvertToRatingFormat } from "src/components/Scenes/SceneDetails/RatingSystem";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) =>
  rating ? (
    <div className={`rating-banner rating-${Math.trunc(rating / 5)}`}>
      <FormattedMessage id="rating" />: {ConvertToRatingFormat(rating)}
    </div>
  ) : (
    <></>
  );
