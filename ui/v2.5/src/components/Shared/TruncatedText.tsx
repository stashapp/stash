import React, { useRef, useState } from "react";
import { Overlay, Tooltip } from "react-bootstrap";
import { Placement } from "react-bootstrap/Overlay";
import debounce from "lodash-es/debounce";
import cx from "classnames";

const CLASSNAME = "TruncatedText";
const CLASSNAME_TOOLTIP = `${CLASSNAME}-tooltip`;

interface ITruncatedTextProps {
  text?: JSX.Element | string | null;
  lineCount?: number;
  placement?: Placement;
  delay?: number;
  className?: string;
}

const TruncatedText: React.FC<ITruncatedTextProps> = ({
  text,
  className,
  lineCount = 1,
  placement = "bottom",
  delay = 1000,
}) => {
  const [showTooltip, setShowTooltip] = useState(false);
  const target = useRef(null);

  if (!text) return <></>;

  const startShowingTooltip = debounce(() => setShowTooltip(true), delay);

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
};

export default TruncatedText;
