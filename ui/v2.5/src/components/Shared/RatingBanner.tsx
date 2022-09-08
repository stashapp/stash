import React from "react";
import { FormattedMessage } from "react-intl";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) =>
  rating ? (
    <div className={`rating-banner rating-${Math.trunc(rating / 5)}`}>
      <FormattedMessage id="rating" />: {rating}
    </div>
  ) : (
    <></>
  );
