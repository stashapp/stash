import React from "react";
import Icon from "./Icon";

const CLASSNAME = "RatingStars";
const CLASSNAME_FILLED = `${CLASSNAME}-filled`;
const CLASSNAME_UNFILLED = `${CLASSNAME}-unfilled`;

interface IProps {
  rating?: number | null;
}

export const RatingStars: React.FC<IProps> = ({ rating }) =>
  rating ? (
    <div className={CLASSNAME}>
      <Icon icon={["fas", "star"]} className={CLASSNAME_FILLED} />
      <Icon
        icon={[rating >= 2 ? "fas" : "far", "star"]}
        className={rating >= 2 ? CLASSNAME_FILLED : CLASSNAME_UNFILLED}
      />
      <Icon
        icon={[rating >= 3 ? "fas" : "far", "star"]}
        className={rating >= 3 ? CLASSNAME_FILLED : CLASSNAME_UNFILLED}
      />
      <Icon
        icon={[rating >= 4 ? "fas" : "far", "star"]}
        className={rating >= 4 ? CLASSNAME_FILLED : CLASSNAME_UNFILLED}
      />
      <Icon
        icon={[rating === 5 ? "fas" : "far", "star"]}
        className={rating === 5 ? CLASSNAME_FILLED : CLASSNAME_UNFILLED}
      />
    </div>
  ) : (
    <></>
  );
