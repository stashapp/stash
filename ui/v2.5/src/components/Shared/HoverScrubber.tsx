import React, { useMemo } from "react";
import cx from "classnames";

interface IHoverScrubber {
  totalSprites: number;
  activeIndex: number | undefined;
  setActiveIndex: (index: number | undefined) => void;
  onClick?: () => void;
}

export const HoverScrubber: React.FC<IHoverScrubber> = ({
  totalSprites,
  activeIndex,
  setActiveIndex,
  onClick,
}) => {
  function getActiveIndex(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    const { width } = e.currentTarget.getBoundingClientRect();
    const x = e.nativeEvent.offsetX;

    const i = Math.round((x / width) * (totalSprites - 1));

    // clamp to [0, totalSprites)
    if (i < 0) return 0;
    if (i >= totalSprites) return totalSprites - 1;
    return i;
  }

  function onMouseMove(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    const relatedTarget = e.currentTarget;

    if (relatedTarget !== e.target) return;

    setActiveIndex(getActiveIndex(e));
  }

  function onMouseLeave() {
    setActiveIndex(undefined);
  }

  function onScrubberClick(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    if (!onClick) return;

    const relatedTarget = e.currentTarget;

    if (relatedTarget !== e.target) return;

    e.preventDefault();
    onClick();
  }

  const indicatorStyle = useMemo(() => {
    if (activeIndex === undefined || !totalSprites) return {};

    const width = ((activeIndex + 1) / totalSprites) * 100;

    return {
      width: `${width}%`,
    };
  }, [activeIndex, totalSprites]);

  return (
    <div
      className={cx("hover-scrubber", {
        "hover-scrubber-inactive": !totalSprites,
      })}
    >
      <div
        className="hover-scrubber-area"
        onMouseMove={onMouseMove}
        onMouseLeave={onMouseLeave}
        onClick={onScrubberClick}
      />
      <div className="hover-scrubber-indicator">
        {activeIndex !== undefined && (
          <div
            className="hover-scrubber-indicator-marker"
            style={indicatorStyle}
          ></div>
        )}
      </div>
    </div>
  );
};
