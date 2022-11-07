import React, { useState } from "react";
import { Button } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";
import {
  convertFromRatingFormat,
  convertToRatingFormat,
  getMaxStars,
  getRatingPrecision,
} from "src/utils/rating";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";

export interface IRatingStarsProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  ratingSystem: GQL.RatingSystem;
}

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const intl = useIntl();
  const [hoverRating, setHoverRating] = useState<number | undefined>();
  const disabled = props.disabled || !props.onSetRating;

  const rating = convertToRatingFormat(props.value, props.ratingSystem);
  const stars = rating ? Math.floor(rating) : 0;
  const fraction = rating ? rating % 1 : 0;

  const max = getMaxStars(props.ratingSystem);
  const precision = getRatingPrecision(props.ratingSystem);

  function newToggleFraction() {
    if (precision > 0) {
      if (fraction !== precision) {
        if (fraction == 0) {
          return 1 - precision;
        }

        return fraction - precision;
      }
    }
  }

  function setRating(thisStar: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = thisStar;

    // toggle rating fraction if we're clicking on the current rating
    if (
      (stars === thisStar && !fraction) ||
      (stars + 1 === thisStar && fraction)
    ) {
      const f = newToggleFraction();
      if (!f) {
        newRating = undefined;
      } else if (fraction) {
        // we're toggling from an existing fraction so use the stars value
        newRating = stars + f;
      } else {
        // we're toggling from a whole value, so decrement from current rating
        newRating = stars - 1 + f;
      }
    }

    // set the hover rating to undefined so that it doesn't immediately clear
    // the stars
    setHoverRating(undefined);

    if (!newRating) {
      props.onSetRating(undefined);
      return;
    }

    props.onSetRating(convertFromRatingFormat(newRating, props.ratingSystem));
  }

  function onMouseOver(thisStar: number) {
    if (!disabled) {
      setHoverRating(thisStar);
    }
  }

  function onMouseOut(thisStar: number) {
    if (!disabled && hoverRating === thisStar) {
      setHoverRating(undefined);
    }
  }

  function getClassName(thisStar: number) {
    if (hoverRating && hoverRating >= thisStar) {
      if (hoverRating === stars) {
        return "unsetting";
      }

      return "setting";
    }

    if (stars && stars >= thisStar) {
      return "set";
    }

    return "unset";
  }

  function getTooltip(thisStar: number) {
    if (disabled && rating) {
      // always return current rating for disabled control
      return rating.toString();
    }

    if (!disabled) {
      // adjust tooltip to use fractions
      if (thisStar === stars && !fraction) {
        const f = newToggleFraction();
        if (!f) {
          return intl.formatMessage({ id: "actions.unset" });
        }
        return (thisStar - 1 + (f ?? 0)).toString();
      } else if (thisStar === stars + 1 && fraction) {
        const f = newToggleFraction();
        if (!f) {
          return intl.formatMessage({ id: "actions.unset" });
        }
        return (thisStar + (f ?? 0)).toString();
      }

      return thisStar.toString();
    }
  }

  function getStyle(thisStar: number) {
    let r: number = hoverRating ? hoverRating : stars;
    let f: number | undefined = fraction;

    if (hoverRating) {
      if (hoverRating === stars && !precision) {
        // unsetting
        return { width: 0 };
      }
      if (hoverRating === stars + 1 && fraction && fraction === precision) {
        // unsetting
        return { width: 0 };
      }

      if (hoverRating === thisStar) {
        if (f && hoverRating === stars + 1) {
          f = newToggleFraction();
          r--;
        } else if (!f && hoverRating === stars) {
          f = newToggleFraction();
          r--;
        }
      } else {
        f = 0;
      }
    }

    return { width: `${getStarWidth(thisStar, r, f ?? 0)}%` };
  }

  function getStarWidth(
    thisRating: number,
    currentStars: number,
    currentFraction: number
  ) {
    if (thisRating > currentStars + 1) {
      return 0;
    }

    if (thisRating <= currentStars) {
      return 100;
    }

    let w = currentFraction * 100;
    // adjust width for 1/4 and 3/4
    if (w == 25) w = 35;
    if (w == 75) w = 65;

    return w;
  }

  const renderRatingButton = (thisStar: number) => (
    <Button
      disabled={disabled}
      className="minimal"
      onClick={() => setRating(thisStar)}
      variant="secondary"
      onMouseEnter={() => onMouseOver(thisStar)}
      onMouseLeave={() => onMouseOut(thisStar)}
      onFocus={() => onMouseOver(thisStar)}
      onBlur={() => onMouseOut(thisStar)}
      title={getTooltip(thisStar)}
      key={`star-${thisStar}`}
    >
      <div className="filled-star" style={getStyle(thisStar)}>
        <Icon icon={fasStar} className="set" />
      </div>
      <div>
        <Icon icon={farStar} className={getClassName(thisStar)} />
      </div>
    </Button>
  );

  return (
    <div className="rating-stars align-middle">
      {Array.from(Array(max)).map((value, index) =>
        renderRatingButton(index + 1)
      )}
    </div>
  );
};
