import React, { useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "../Icon";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";
import {
  convertFromRatingFormat,
  convertToRatingFormat,
  getRatingPrecision,
  RatingStarPrecision,
  RatingSystemType,
} from "src/utils/rating";
import { useIntl } from "react-intl";

export interface IRatingStarsProps {
  value: number | null;
  onSetRating?: (value: number | null) => void;
  disabled?: boolean;
  precision: RatingStarPrecision;
  valueRequired?: boolean;
}

type HoverState = {
  star: number;
  fraction: number;
};

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const intl = useIntl();
  const [hoverRating, setHoverRating] = useState<HoverState | undefined>();
  const [hoveredStar, setHoveredStar] = useState<number | undefined>();
  const disabled = props.disabled || !props.onSetRating;

  const rating = convertToRatingFormat(props.value, {
    type: RatingSystemType.Stars,
    starPrecision: props.precision,
  });
  const stars = rating ? Math.floor(rating) : 0;
  // the upscaling was necesary to fix rounding issue present with tenth place precision
  const fraction = rating ? ((rating * 10) % 10) / 10 : 0;

  const max = 5;
  const precision = getRatingPrecision(props.precision);

  function newToggleFraction() {
    if (precision !== 1) {
      if (fraction !== precision) {
        if (fraction == 0) {
          return 1 - precision;
        }

        return fraction - precision;
      }
    }
  }

  function setRating(thisStar: number, clickFraction?: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = thisStar;

    // if we have a click fraction (from half-star clicking), use it directly
    if (clickFraction !== undefined && precision !== 1) {
      const targetRating = thisStar - 1 + clickFraction;

      // check if clicking on the same rating to toggle/unset
      const currentTotal = stars + fraction;
      if (Math.abs(currentTotal - targetRating) < 0.01) {
        if (props.valueRequired) {
          newRating = targetRating;
        } else {
          newRating = undefined;
        }
      } else {
        newRating = targetRating;
      }
    } else {
      // original toggle logic for non-half-star clicks or full precision
      if (
        (stars === thisStar && !fraction) ||
        (stars + 1 === thisStar && fraction)
      ) {
        const f = newToggleFraction();
        if (!f) {
          if (props.valueRequired) {
            if (fraction) {
              newRating = stars + 1;
            } else {
              newRating = stars;
            }
          } else {
            newRating = undefined;
          }
        } else if (fraction) {
          // we're toggling from an existing fraction so use the stars value
          newRating = stars + f;
        } else {
          // we're toggling from a whole value, so decrement from current rating
          newRating = stars - 1 + f;
        }
      }
    }

    // set the hover rating to undefined so that it doesn't immediately clear
    // the stars
    setHoverRating(undefined);

    if (!newRating) {
      props.onSetRating(null);
      return;
    }

    props.onSetRating(
      convertFromRatingFormat(newRating, RatingSystemType.Stars)
    );
  }

  function getNextRatingForStar(
    thisStar: number,
    mouseX?: number
  ): { star: number; fraction: number } | null {
    // for half precision with mouse position, detect which half
    if (precision === 0.5 && mouseX !== undefined) {
      const isRightHalf = mouseX > 0.5; // mouseX should be normalized 0-1
      const hoverFraction = isRightHalf ? 1 : 0.5;
      return { star: thisStar, fraction: hoverFraction };
    }

    // for other precisions, use the original toggle logic
    let nextRating = thisStar;
    let nextFraction = 0;

    if (
      (stars === thisStar && !fraction) ||
      (stars + 1 === thisStar && fraction)
    ) {
      const f = newToggleFraction();
      if (!f) {
        if (props.valueRequired) {
          if (fraction) {
            nextRating = stars + 1;
            nextFraction = 0;
          } else {
            nextRating = stars;
            nextFraction = 0;
          }
        } else {
          // unset rating
          return null;
        }
      } else if (fraction) {
        nextRating = stars;
        nextFraction = f;
      } else {
        nextRating = stars - 1;
        nextFraction = f;
      }
    } else {
      nextRating = thisStar - 1;
      nextFraction = 1;
    }

    return { star: nextRating + 1, fraction: nextFraction };
  }

  function updateHoverRating(
    thisStar: number,
    event?: React.MouseEvent<HTMLButtonElement>
  ) {
    if (disabled) return;

    setHoveredStar(thisStar);

    let mouseX: number | undefined;
    if (event && precision === 0.5) {
      const rect = event.currentTarget.getBoundingClientRect();
      mouseX = (event.clientX - rect.left) / rect.width;
    }

    const nextRating = getNextRatingForStar(thisStar, mouseX);
    if (nextRating) {
      setHoverRating(nextRating);
    } else {
      setHoverRating({ star: 0, fraction: 0 }); // unset
    }
  }

  function handleClick(
    event: React.MouseEvent<HTMLButtonElement>,
    thisStar: number
  ) {
    // for half precision, detect which half of the button was clicked
    if (precision === 0.5) {
      const rect = event.currentTarget.getBoundingClientRect();
      const clickX = event.clientX - rect.left;
      const buttonWidth = rect.width;
      const isRightHalf = clickX > buttonWidth / 2;

      const clickFraction = isRightHalf ? 1 : 0.5;
      setRating(thisStar, clickFraction);
    } else {
      setRating(thisStar);
    }
  }

  function onMouseOver(
    event: React.MouseEvent<HTMLButtonElement>,
    thisStar: number
  ) {
    updateHoverRating(thisStar, event);
  }

  function onMouseMove(
    event: React.MouseEvent<HTMLButtonElement>,
    thisStar: number
  ) {
    if (precision === 0.5) {
      updateHoverRating(thisStar, event);
    }
  }

  function onMouseOut(thisStar: number) {
    if (!disabled && hoveredStar === thisStar) {
      setHoverRating(undefined);
      setHoveredStar(undefined);
    }
  }

  function getClassName(thisStar: number) {
    const hoverTotal = hoverRating
      ? hoverRating.star - 1 + hoverRating.fraction
      : undefined;
    const currentTotal = stars + fraction;

    if (hoverTotal !== undefined && hoverTotal >= thisStar) {
      if (Math.abs(hoverTotal - currentTotal) < 0.01) {
        return "unsetting";
      }
      return "setting";
    }

    if (stars && stars >= thisStar) {
      return "set";
    }

    return "unset";
  }

  function getTooltip(thisStar: number, current: RatingFraction | undefined) {
    if (disabled) {
      if (rating) {
        // always return current rating for disabled control
        return rating.toString();
      }

      return undefined;
    }

    // adjust tooltip to use fractions
    if (!current) {
      return intl.formatMessage({ id: "actions.unset" });
    }

    return (current.rating + current.fraction).toString();
  }

  type RatingFraction = {
    rating: number;
    fraction: number;
  };

  function getCurrentSelectedRating(): RatingFraction | undefined {
    if (hoverRating) {
      // star 0 means unset
      if (hoverRating.star === 0) {
        return undefined;
      }

      const hoverTotal = hoverRating.star - 1 + hoverRating.fraction;
      const currentTotal = stars + fraction;

      // check if hovering over current rating (for unsetting)
      if (Math.abs(hoverTotal - currentTotal) < 0.01 && !props.valueRequired) {
        return undefined;
      }

      return {
        rating: hoverRating.star - 1,
        fraction: hoverRating.fraction,
      };
    }

    return { rating: stars, fraction: fraction };
  }

  function getButtonClassName(
    thisStar: number,
    current: RatingFraction | undefined
  ) {
    if (!current || thisStar > current.rating + 1) {
      return "star-fill-0";
    }

    if (thisStar <= current.rating) {
      return "star-fill-100";
    }

    let w = current.fraction * 100;
    return `star-fill-${w}`;
  }

  const renderRatingButton = (thisStar: number) => {
    const ratingFraction = getCurrentSelectedRating();

    return (
      <Button
        disabled={disabled}
        className={`minimal ${getButtonClassName(thisStar, ratingFraction)}`}
        onClick={(event: React.MouseEvent<HTMLButtonElement>) =>
          handleClick(event, thisStar)
        }
        variant="secondary"
        onMouseEnter={(event: React.MouseEvent<HTMLButtonElement>) =>
          onMouseOver(event, thisStar)
        }
        onMouseMove={(event: React.MouseEvent<HTMLButtonElement>) =>
          onMouseMove(event, thisStar)
        }
        onMouseLeave={() => onMouseOut(thisStar)}
        onFocus={() => updateHoverRating(thisStar)}
        onBlur={() => onMouseOut(thisStar)}
        title={getTooltip(thisStar, ratingFraction)}
        key={`star-${thisStar}`}
      >
        <div className="filled-star">
          <Icon icon={fasStar} className="set" />
        </div>
        <div className="unfilled-star">
          <Icon icon={farStar} className={getClassName(thisStar)} />
        </div>
      </Button>
    );
  };

  const maybeRenderStarRatingNumber = () => {
    const ratingFraction = getCurrentSelectedRating();
    if (
      !ratingFraction ||
      (ratingFraction.rating == 0 && ratingFraction.fraction == 0)
    ) {
      return;
    }

    return (
      <span className="star-rating-number">
        {ratingFraction.rating + ratingFraction.fraction}
      </span>
    );
  };

  const precisionClassName = `rating-stars-precision-${props.precision}`;

  return (
    <div className={`rating-stars ${precisionClassName}`}>
      {Array.from(Array(max)).map((value, index) =>
        renderRatingButton(index + 1)
      )}
      {maybeRenderStarRatingNumber()}
    </div>
  );
};
