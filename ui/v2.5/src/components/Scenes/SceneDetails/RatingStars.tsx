import React, { useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";

export interface IRatingStarsProps {
  value?: number;
  onSetRating?: (value?: number) => void;
}

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const [hoverRating, setHoverRating] = useState<number | undefined>();
  const disabled = !props.onSetRating;

  function setRating(rating: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = rating;

    // unset if we're clicking on the current rating
    if (props.value === rating) {
      newRating = undefined;
    }

    // set the hover rating to undefined so that it doesn't immediately clear
    // the stars
    setHoverRating(undefined);

    props.onSetRating(newRating);
  }

  function getIconPrefix(rating: number) {
    if (hoverRating && hoverRating >= rating) {
      if (hoverRating === props.value) {
        return "far";
      }

      return "fas";
    }

    if (!hoverRating && props.value && props.value >= rating) {
      return "fas";
    }

    return "far";
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
      if (hoverRating === props.value) {
        return "unsetting";
      }

      return "setting";
    }

    if (props.value && props.value >= rating) {
      return "set";
    }

    return "unset";
  }

  function getTooltip(rating: number) {
    if (disabled && props.value) {
      // always return current rating for disabled control
      return props.value.toString();
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
    >
      <Icon
        icon={[getIconPrefix(rating), "star"]}
        className={getClassName(rating)}
      />
    </Button>
  );

  const maxRating = 5;

  return (
    <div className="rating-stars">
      {Array.from(Array(maxRating)).map((value, index) =>
        renderRatingButton(index + 1)
      )}
    </div>
  );
};
