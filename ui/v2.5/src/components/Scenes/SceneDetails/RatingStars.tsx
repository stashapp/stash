import React, { useState } from "react";
import { Button } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";

export interface IRatingStarsProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  precision: number;
  max: number;
}

const maxRating = 100;

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const [hoverRating, setHoverRating] = useState<number | undefined>();
  const disabled = props.disabled || !props.onSetRating;

  const stars = props.value ? convertFrom100(props.value) : 0;

  function convertFrom100(rating100: number) {
    return (
      Math.round(
        (1 / props.precision) * (rating100 / maxRating) * props.max
      ) /
      (1 / props.precision)
    );
  }

  function convertTo100(rating: number | undefined) {
    return rating == undefined
      ? undefined
      : Math.round((rating / props.max) * maxRating);
  }

  function setRating(rating: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = rating;

    // unset if we're clicking on the current rating
    if (stars === rating) {
      newRating = undefined;
    }

    // set the hover rating to undefined so that it doesn't immediately clear
    // the stars
    setHoverRating(undefined);

    props.onSetRating(convertTo100(newRating));
  }

  function getIcon(rating: number) {
    if (hoverRating && hoverRating >= rating) {
      if (hoverRating === stars) {
        return farStar;
      }

      return fasStar;
    }

    if (!hoverRating && stars && stars >= rating) {
      return fasStar;
    }

    return farStar;
  }

  function onMouseOver(rating: number) {
    if (!disabled) {
      setHoverRating(rating);
    }
  }

  function onMouseOut(rating: number) {
    if (!disabled && hoverRating === rating) {
      setHoverRating(undefined);
    }
  }

  function getClassName(rating: number) {
    if (hoverRating && hoverRating >= rating) {
      if (hoverRating === stars) {
        return "unsetting";
      }

      return "setting";
    }

    if (stars && stars >= rating) {
      return "set";
    }

    return "unset";
  }

  function getTooltip(rating: number) {
    if (disabled && stars) {
      // always return current rating for disabled control
      return stars.toString();
    }

    if (!disabled) {
      return rating.toString();
    }
  }

  const renderRatingButton = (rating: number) => (
    <Button
      disabled={disabled}
      className="minimal"
      onClick={() => setRating(rating)}
      variant="secondary"
      onMouseOver={() => onMouseOver(rating)}
      onMouseOut={() => onMouseOut(rating)}
      onFocus={() => onMouseOver(rating)}
      onBlur={() => onMouseOut(rating)}
      title={getTooltip(rating)}
      key={`star-${rating}`}
    >
      <Icon icon={getIcon(rating)} className={getClassName(rating)} />
    </Button>
  );

  return (
    <div className="rating-stars align-middle">
      {Array.from(Array(props.max)).map((value, index) =>
        renderRatingButton(index + 1)
      )}
    </div>
  );
};