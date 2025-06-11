import React, { useRef, useState } from "react";
import { Overlay, Tooltip } from "react-bootstrap";
import { Placement } from "react-bootstrap/Overlay";
import cx from "classnames";
import { useDebounce } from "src/hooks/debounce";
import { PatchComponent } from "src/patch";

const CLASSNAME = "TruncatedText";
const CLASSNAME_TOOLTIP = `${CLASSNAME}-tooltip`;

interface ITruncatedTextProps {
  text?: JSX.Element | string | null;
  lineCount?: number;
  placement?: Placement;
  delay?: number;
  className?: string;
}

export const TruncatedText: React.FC<ITruncatedTextProps> = PatchComponent(
  "TruncatedText",
  ({ text, className, lineCount = 1, placement = "bottom", delay = 1000 }) => {
    const [showTooltip, setShowTooltip] = useState(false);
    const target = useRef(null);

    const startShowingTooltip = useDebounce(() => setShowTooltip(true), delay);

    if (!text) return <></>;

    const handleFocus = (element: HTMLElement) => {
      // Check if visible size is smaller than the content size
      if (
        element.offsetWidth < element.scrollWidth ||
        element.offsetHeight + 10 < element.scrollHeight
      )
        startShowingTooltip();
    };

    const handleBlur = () => {
      startShowingTooltip.cancel();
      setShowTooltip(false);
    };

    const overlay = (
      <Overlay target={target.current} show={showTooltip} placement={placement}>
        <Tooltip id={CLASSNAME} className={CLASSNAME_TOOLTIP}>
          {text}
        </Tooltip>
      </Overlay>
    );

    return (
      <div
        className={cx(CLASSNAME, className)}
        style={{ WebkitLineClamp: lineCount }}
        ref={target}
        onMouseEnter={(e) => handleFocus(e.currentTarget)}
        onFocus={(e) => handleFocus(e.currentTarget)}
        onMouseLeave={handleBlur}
        onBlur={handleBlur}
      >
        {text}
        {overlay}
      </div>
    );
  }
);

export const TruncatedInlineText: React.FC<ITruncatedTextProps> = ({
  text,
  className,
  placement = "bottom",
  delay = 1000,
}) => {
  const [showTooltip, setShowTooltip] = useState(false);
  const target = useRef(null);

  const startShowingTooltip = useDebounce(() => setShowTooltip(true), delay);

  if (!text) return <></>;

  const handleFocus = (element: HTMLElement) => {
    // Check if visible size is smaller than the content size
    if (
      element.offsetWidth < element.scrollWidth ||
      element.offsetHeight + 10 < element.scrollHeight
    )
      startShowingTooltip();
  };

  const handleBlur = () => {
    startShowingTooltip.cancel();
    setShowTooltip(false);
  };

  const overlay = (
    <Overlay target={target.current} show={showTooltip} placement={placement}>
      <Tooltip id={CLASSNAME} className={CLASSNAME_TOOLTIP}>
        {text}
      </Tooltip>
    </Overlay>
  );

  return (
    <span
      className={cx(CLASSNAME, "inline", className)}
      ref={target}
      onMouseEnter={(e) => handleFocus(e.currentTarget)}
      onFocus={(e) => handleFocus(e.currentTarget)}
      onMouseLeave={handleBlur}
      onBlur={handleBlur}
    >
      {text}
      {overlay}
    </span>
  );
};
