import React from "react";
import { FormattedMessage } from "react-intl";

interface IProps {
  rating?: number | null;
}

export const RatingBanner: React.FC<IProps> = ({ rating }) =>
  rating ? (
    <div className={`rating-banner rating-${rating}`}>
      <FormattedMessage id="rating" />: {rating}
    </div>
  ) : (
    <></>
  );
