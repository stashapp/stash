import React, { useRef, useState } from "react";
import { Overlay, Tooltip } from "react-bootstrap";
import { Placement } from "react-bootstrap/Overlay";
import debounce from "lodash-es/debounce";
import cx from "classnames";
import { ConfigurationContext } from "src/hooks/Config";

const CLASSNAME = "TruncatedText";
const CLASSNAME_TOOLTIP = `${CLASSNAME}-tooltip`;

interface ITruncatedTextProps {
  text?: string | null;
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
  const { isTouch } = React.useContext(ConfigurationContext);

  if (!text) return <></>;

  const startDeferredShowingTooltip = debounce(
    () => setShowTooltip(true),
    delay
  );

  const handleFocus = (element: HTMLElement, defer: boolean = true) => {
    // Check if visible size is smaller than the content size
    if (
      element.offsetWidth < element.scrollWidth ||
      element.offsetHeight + 10 < element.scrollHeight
    ) {
      if (defer) {
        startDeferredShowingTooltip();
      } else {
        setShowTooltip(true);
      }
    }
  };

  const handleBlur = () => {
    startDeferredShowingTooltip.cancel();
    setShowTooltip(false);
  };

  const handleClick = (element: HTMLElement) => {
    showTooltip ? handleBlur() : handleFocus(element, false);
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
      onMouseEnter={isTouch ? undefined : (e) => handleFocus(e.currentTarget)}
      onFocus={isTouch ? undefined : (e) => handleFocus(e.currentTarget)}
      onMouseLeave={handleBlur}
      onBlur={handleBlur}
      onClick={isTouch ? (e) => handleClick(e.currentTarget) : undefined}
    >
      {text}
      {overlay}
    </div>
  );
};

export default TruncatedText;
