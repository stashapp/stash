import React, { useState, useRef } from "react";
import * as GQL from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import Box from "@mui/material/Box";
import StarIcon from "@mui/icons-material/Star";
import StarBorderIcon from "@mui/icons-material/StarOutline";
import { FormattedMessage, useIntl } from "react-intl";
import { Col, Form, Row } from "react-bootstrap";
import { FormUtils } from "src/utils";

export interface IRatingSystemProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  excludeLabel?: boolean;
}

export interface IRatingStarsProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
  precision: number;
  maxRating: number;
}

export interface IRatingNumberProps {
  value?: number;
  onSetRating?: (value?: number) => void;
  disabled?: boolean;
}

export const RatingNumber: React.FC<IRatingNumberProps> = (
  props: IRatingNumberProps
) => {
  const [input, setInput] = useState<string | "0.0">();
  const [previous, setPrevious] = useState<string | "0.0">();
  const [useValidation, setValidation] = useState<boolean | true>();
  function stepChange() {
    setValidation(false);
  }

  function nonStepChange() {
    setValidation(true);
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    let val = e.target.value;
    if (!useValidation && props.onSetRating != null) {
      e.target.value = Number(val).toFixed(1);
      setInput(Number(val).toFixed(1));
      setPrevious(Number(val).toFixed(1));
      let tempVal = Number(val) * 10;
      props.onSetRating(tempVal != 0 ? tempVal : undefined);
      setValidation(true);
      return;
    }
    const match = /(\d?)(\d?)(.?)((\d)?)/g.exec(val);
    const matchOld = /(\d?)(\d?)(.?)((\d{0,2})?)/g.exec(previous ?? "");

    if (match == null || props.onSetRating == null) {
      return;
    }

    if (match[2] && !(match[2] == "0" && match[1] == "1")) {
      match[2] = "";
    }
    if (match[4] == null || match[4] == "") {
      match[4] = "0";
    }
    let value = match[1] + match[2] + "." + match[4];
    e.target.value = value;
    if (val.length > 0) {
      if (Number(value) > 10) {
        value = "10.0";
      }
      e.target.value = Number(value).toFixed(1);
      setInput(Number(value).toFixed(1));
      setPrevious(Number(value).toFixed(1));
      let tempVal = Number(value) * 10;
      props.onSetRating(tempVal != 0 ? tempVal : undefined);

      // This is a workaround to a missing feature where you can't set cursor position in input numbers.
      // See https://stackoverflow.com/questions/33406169/failed-to-execute-setselectionrange-on-htmlinputelement-the-input-elements
      e.target.type = "text";
      let cursorPosition = 0;
      if (match[2] && !match[4]) {
        cursorPosition = 3;
      } else if (matchOld != null && match[1] !== matchOld[1]) {
        cursorPosition = 2;
      } else if (
        matchOld != null &&
        match[1] === matchOld[1] &&
        match[2] === matchOld[2] &&
        match[4] === matchOld[4]
      ) {
        cursorPosition = 2;
      }
      e.target.setSelectionRange(cursorPosition, cursorPosition);
      e.target.type = "number";
    }
  }

  if (props.disabled) {
    return <text>{Number((props.value ?? 0) / 10).toFixed(1)}</text>;
  } else {
    return (
      <div>
        <input
          className="text-input"
          type="number"
          onMouseDown={stepChange}
          onKeyDown={nonStepChange}
          onChange={handleChange}
          value={input}
          defaultValue={
            props.value == null || props.value == undefined
              ? "0.0"
              : Number(props.value / 10).toFixed(1)
          }
          min="0.0"
          step="0.1"
          max="10"
          style={{ fontSize: "22px", padding: "4px" }}
          placeholder="0.0"
        />
      </div>
    );
  }
};

export const RatingStars: React.FC<IRatingStarsProps> = (
  props: IRatingStarsProps
) => {
  const disabled = props.disabled || !props.onSetRating;
  const EmptyIcon = StarBorderIcon;
  const FilledIcon = StarIcon;
  const ratingContainerRef = useRef<HTMLInputElement>(null);
  const [isHovered, setIsHovered] = useState(false);
  const realMaxRating = 100;
  const [activeStar, setActiveStar] = useState(
    props.value != null && props.value != undefined
      ? convertFromBigScale(props.value)
      : -1
  );
  const [hoverActiveStar, setHoverActiveStar] = useState(-1);

  function setRating(rating: number) {
    if (!props.onSetRating) {
      return;
    }

    let newRating: number | undefined = rating;

    // unset if we're clicking on the current rating
    if (activeStar === rating) {
      newRating = undefined;
      setActiveStar(-1);
    }
    let temp = convertToBigScale(newRating);
    props.onSetRating(temp);
  }

  function convertFromBigScale(bigRating: number) {
    return (
      Math.round(
        (1 / props.precision) * (bigRating / realMaxRating) * props.maxRating
      ) /
      (1 / props.precision)
    );
  }

  function convertToBigScale(smallRating: number | undefined) {
    return smallRating == undefined
      ? undefined
      : Math.round((smallRating / props.maxRating) * realMaxRating);
  }

  function calculateRating(e: React.MouseEvent<Element, MouseEvent>) {
    if (ratingContainerRef != null && ratingContainerRef.current != null) {
      const {
        width,
        left,
      } = ratingContainerRef.current.getBoundingClientRect();
      let percent = (e.clientX - left) / width;
      const numberInStars = percent * props.maxRating;
      const nearestNumber =
        Math.round((numberInStars + props.precision / 2) / props.precision) *
        props.precision;

      return Number(
        nearestNumber.toFixed(
          props.precision.toString().split(".")[1]?.length || 0
        )
      );
    } else {
      return -1;
    }
  }

  function handleClick(e: React.MouseEvent<Element, MouseEvent>) {
    if (disabled) return;
    setIsHovered(false);
    let calculatedRating = calculateRating(e);
    setActiveStar(calculatedRating);
    setRating(calculatedRating);
  }

  function handleMouseMove(e: React.MouseEvent<Element, MouseEvent>) {
    if (disabled) return;
    setIsHovered(true);
    setHoverActiveStar(calculateRating(e));
  }

  function handleMouseLeave() {
    if (disabled) return;
    setHoverActiveStar(-1);
    setIsHovered(false);
  }

  return (
    <Box
      sx={{
        display: "inline-flex",
        position: "relative",
        cursor: "pointer",
        textAlign: "left",
      }}
      onClick={(e: React.MouseEvent<HTMLElement>) => handleClick(e)}
      onMouseMove={(e: React.MouseEvent<HTMLElement>) => handleMouseMove(e)}
      onMouseLeave={() => handleMouseLeave()}
      ref={ratingContainerRef}
    >
      {[...new Array(props.maxRating)].map((arr, index) => {
        const activeState = isHovered ? hoverActiveStar : activeStar;

        const showEmptyIcon = activeState === -1 || activeState < index + 1;

        const isActiveRating = activeState !== 1;
        const isRatingWithPrecision = activeState % 1 !== 0;
        const isRatingEqualToIndex = Math.ceil(activeState) === index + 1;
        const showRatingWithPrecision =
          isActiveRating && isRatingWithPrecision && isRatingEqualToIndex;

        return (
          <Box
            position={"relative"}
            sx={{
              cursor: "pointer",
              transform: disabled ? "" : "translateY(25%)",
            }}
            key={index}
          >
            <Box
              sx={{
                width: showRatingWithPrecision
                  ? `${
                      (activeState % 1 == 0.25
                        ? 0.35
                        : activeState % 1 == 0.75
                        ? 0.6
                        : activeState % 1) * 100
                    }%`
                  : "0%",
                overflow: "hidden",
                position: "fixed",
              }}
            >
              <FilledIcon sx={{ color: "gold" }} />
            </Box>
            <Box
              sx={{
                color: showEmptyIcon ? "gold" : "inherit",
              }}
            >
              {showEmptyIcon ? (
                <EmptyIcon className="unsetting" />
              ) : (
                <FilledIcon className="set" sx={{ color: "gold" }} />
              )}
            </Box>
          </Box>
        );
      })}
    </Box>
  );
};

export const RatingSystem: React.FC<IRatingSystemProps> = (
  props: IRatingSystemProps
) => {
  function getRatingStars(maxRating: number, precision: number) {
    return (
      <RatingStars
        value={props.value}
        onSetRating={props.onSetRating}
        disabled={props.disabled}
        precision={precision}
        maxRating={maxRating}
      />
    );
  }

  const { configuration: config } = React.useContext(ConfigurationContext);
  let toReturn;
  const intl = useIntl();
  switch (config?.interface?.ratingSystem) {
    case GQL.RatingSystem.TenStar:
      toReturn = getRatingStars(10, 1);
      break;
    case GQL.RatingSystem.TenPointFiveStar:
      toReturn = getRatingStars(10, 0.5);
      break;
    case GQL.RatingSystem.TenPointTwoFiveStar:
      toReturn = getRatingStars(10, 0.25);
      break;
    case GQL.RatingSystem.FiveStar:
      toReturn = getRatingStars(5, 1);
      break;
    case GQL.RatingSystem.FivePointFiveStar:
      toReturn = getRatingStars(5, 0.5);
      break;
    case GQL.RatingSystem.FivePointTwoFiveStar:
      toReturn = getRatingStars(5, 0.25);
      break;
    case GQL.RatingSystem.TenPointDecimal:
      toReturn = (
        <RatingNumber
          value={props.value}
          onSetRating={props.onSetRating}
          disabled={props.disabled}
        />
      );
      break;
    case GQL.RatingSystem.None:
      return <></>;
      break;
    default:
      toReturn = getRatingStars(5, 0.5);
      break;
  }

  if (props.excludeLabel) {
    return toReturn;
  }

  if (props.disabled) {
    return (
      <h6>
        <FormattedMessage id="rating" />: {toReturn}
      </h6>
    );
  } else {
    return (
      <>
        <Form.Group controlId="rating" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "rating" }),
          })}
          <Col xs={9}>{toReturn}</Col>
        </Form.Group>
      </>
    );
  }
};
