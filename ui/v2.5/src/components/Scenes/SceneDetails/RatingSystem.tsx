import React, { useState, MouseEventHandler, useRef } from "react";
import { Button } from "react-bootstrap";
import Icon from "src/components/Shared/Icon";
import { faStar as fasStar } from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";
import { faStarHalf as faStarHalf } from "@fortawesome/free-regular-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";

import Box from "@mui/material/Box";
import StarIcon from "@mui/icons-material/Star";
import StarBorderIcon from "@mui/icons-material/StarBorder";

export interface IRatingSystemProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
}

export interface IRatingStarsProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  precision: number;
  maxRating: number
}

export const RatingSystem: React.FC<IRatingSystemProps> = (
  props: IRatingSystemProps
) => {
  function getRatingStars(maxRating: number, precision: number) {
    return (
      <RatingStars value={props.value} onSetRating={props.onSetRating} disabled={props.disabled} precision={precision} maxRating={maxRating} />
    );
  }

  const { configuration: config } = React.useContext(ConfigurationContext);
  let toReturn;
  switch (config?.interface?.ratingSystem) {
    case GQL.RatingSystem.TenStar:
      toReturn = getRatingStars(10, 1);
      break;
    case GQL.RatingSystem.FiveStar:
      toReturn = getRatingStars(5, 1);
      break;
    default:
      toReturn = getRatingStars(5, 0.5);
      break
  }
  return toReturn;
};

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const [hoverRating, setHoverRating] = useState<number | undefined>();
  const disabled = props.disabled || !props.onSetRating;

  const EmptyIcon = StarBorderIcon;
  const FilledIcon = StarIcon;
  const ratingContainerRef = useRef<HTMLInputElement>(null);
  const [isHovered, setIsHovered] = useState(false);
  const realMaxRating = 100;
  let maxRating = props.maxRating;
  const [activeStar, setActiveStar] = useState(props.value != null && props.value != undefined ? convertFromBigScale(props.value) : -1);
  const [hoverActiveStar, setHoverActiveStar] = useState(-1);



  function setRating(rating: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = activeStar;

    // unset if we're clicking on the current rating
    if (props.value === activeStar) {
      newRating = undefined;
    }
    props.onSetRating(convertToBigScale(newRating));
  }

  function convertFromBigScale(bigRating: number) {
    return Math.round((bigRating / realMaxRating) * maxRating);
  }

  function convertToBigScale(smallRating: number | undefined) {
    return smallRating == undefined ? undefined : Math.round(smallRating / maxRating * realMaxRating );
  }

  function getIcon(rating: number) {
    if (hoverRating && hoverRating >= rating) {
      if (hoverRating === props.value) {
        return farStar;
      }

      return fasStar;
    }

    if (!hoverRating && props.value && props.value >= rating) {
      return fasStar;
    }

    return farStar;
  }

  function onMouseOver(rating: number) {
    if (!disabled) {
      setHoverRating(rating);
    }
  }

  function calculateRating(e: React.MouseEvent<Element, MouseEvent>) {
    if (ratingContainerRef != null && ratingContainerRef.current != null) {
      const { width, left } = ratingContainerRef.current.getBoundingClientRect();
      let percent = (e.clientX - left) / width;
      const numberInStars = percent * maxRating;
      const nearestNumber = Math.round((numberInStars + props.precision / 2) / props.precision) * props.precision;

      return Number(nearestNumber.toFixed(props.precision.toString().split('.')[1]?.length || 0));
    }
    else {
      return -1;
    }
  }


  function handleClick(e: React.MouseEvent<Element, MouseEvent>) {
    setIsHovered(false);
    setActiveStar(calculateRating(e));
    setRating(activeStar);
  };

  function handleMouseMove(e: React.MouseEvent<Element, MouseEvent>) {
    setIsHovered(true);
    setHoverActiveStar(calculateRating(e));
  };

  function handleMouseLeave(e: React.MouseEvent<Element, MouseEvent>) {
    setHoverActiveStar(-1); // Reset to default state
    setIsHovered(false);
  };

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
      key={`star-${rating}`}
    >
      <Icon icon={getIcon(rating)} className={getClassName(rating)} />
    </Button>
  );

  //return (
  //  <div className="rating-stars align-middle">
  //    {Array.from(Array(maxRating)).map((value, index) =>
  //      renderRatingButton(index + 1)
  //    )}
  //  </div>
  //);

  return (
    <Box
      sx={{
        display: 'inline-flex',
        position: 'relative',
        cursor: 'pointer',
        textAlign: 'left'
      }}
      onClick={(e: React.MouseEvent<HTMLElement>) => handleClick(e)}
      onMouseMove={(e: React.MouseEvent<HTMLElement>) => handleMouseMove(e)}
      onMouseLeave={(e: React.MouseEvent<HTMLElement>) => handleMouseLeave(e)}
      ref={ratingContainerRef}
    >
      {[...new Array(maxRating)].map((arr, index) => {
        const activeState = isHovered ? hoverActiveStar : activeStar;

        const showEmptyIcon = activeState === -1 || activeState < index + 1;

        const isActiveRating = activeState !== 1;
        const isRatingWithPrecision = activeState % 1 !== 0;
        const isRatingEqualToIndex = Math.ceil(activeState) === index + 1;
        const showRatingWithPrecision =
          isActiveRating && isRatingWithPrecision && isRatingEqualToIndex;

        return (
          <Box
            position={'relative'}
            sx={{
              cursor: 'pointer'
            }}
            key={index}
          >
            <Box
              sx={{
                width: showRatingWithPrecision ? `${(activeState % 1) * 100}%` : '0%',
                overflow: 'hidden',
                position: 'absolute'
              }}
            >
              <FilledIcon />
            </Box>
            {/*Note here */}
            <Box
              sx={{
                color: showEmptyIcon ? 'gray' : 'inherit'
              }}
            >
              {showEmptyIcon ? <EmptyIcon /> : <FilledIcon />}
            </Box>
          </Box>
        );
      })}
    </Box>
  );
};
