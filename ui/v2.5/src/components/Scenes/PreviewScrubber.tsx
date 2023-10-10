import React, { useRef, useMemo, useState, useLayoutEffect } from "react";
import { useSpriteInfo } from "src/hooks/sprite";
import { useThrottle } from "src/hooks/throttle";
import TextUtils from "src/utils/text";

interface IHoverScrubber {
  totalSprites: number;
  activeIndex: number | undefined;
  setActiveIndex: (index: number | undefined) => void;
  onClick?: () => void;
}

const HoverScrubber: React.FC<IHoverScrubber> = ({
  totalSprites,
  activeIndex,
  setActiveIndex,
  onClick,
}) => {
  function getActiveIndex(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    const { width } = e.currentTarget.getBoundingClientRect();
    const x = e.nativeEvent.offsetX;

    const i = Math.floor((x / width) * totalSprites);

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

    const width = (activeIndex / totalSprites) * 100;

    return {
      width: `${width}%`,
    };
  }, [activeIndex, totalSprites]);

  return (
    <div className="hover-scrubber">
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

interface IScenePreviewProps {
  vttPath: string | undefined;
  onClick?: (timestamp: number) => void;
}

function scaleToFit(dimensions: { w: number; h: number }, bounds: DOMRect) {
  const rw = bounds.width / dimensions.w;
  const rh = bounds.height / dimensions.h;

  // for consistency, use max by default and min for portrait
  if (dimensions.w > dimensions.h) {
    return Math.max(rw, rh);
  }

  return Math.min(rw, rh);
}

export const PreviewScrubber: React.FC<IScenePreviewProps> = ({
  vttPath,
  onClick,
}) => {
  const imageParentRef = useRef<HTMLDivElement>(null);
  const [style, setStyle] = useState({});

  const [activeIndex, setActiveIndex] = useState<number>();

  const debounceSetActiveIndex = useThrottle(setActiveIndex, 50);

  const spriteInfo = useSpriteInfo(vttPath);

  const sprite = useMemo(() => {
    if (!spriteInfo || activeIndex === undefined) {
      return undefined;
    }
    return spriteInfo[activeIndex];
  }, [activeIndex, spriteInfo]);

  useLayoutEffect(() => {
    const imageParent = imageParentRef.current;

    if (!sprite || !imageParent) {
      return setStyle({});
    }

    const clientRect = imageParent.getBoundingClientRect();
    const scale = scaleToFit(sprite, clientRect);

    setStyle({
      backgroundPosition: `${-sprite.x}px ${-sprite.y}px`,
      backgroundImage: `url(${sprite.url})`,
      width: `${sprite.w}px`,
      height: `${sprite.h}px`,
      transform: `scale(${scale})`,
    });
  }, [sprite]);

  const currentTime = useMemo(() => {
    if (!sprite) return undefined;

    const start = TextUtils.secondsToTimestamp(sprite.start);

    return start;
  }, [sprite]);

  function onScrubberClick() {
    if (!sprite || !onClick) {
      return;
    }

    onClick(sprite.start);
  }

  if (!spriteInfo) return null;

  return (
    <div className="preview-scrubber">
      {sprite && (
        <div className="scene-card-preview-image" ref={imageParentRef}>
          <div className="scrubber-image" style={style}></div>
          {currentTime !== undefined && (
            <div className="scrubber-timestamp">{currentTime}</div>
          )}
        </div>
      )}
      <HoverScrubber
        totalSprites={spriteInfo.length}
        activeIndex={activeIndex}
        setActiveIndex={(i) => debounceSetActiveIndex(i)}
        onClick={onScrubberClick}
      />
    </div>
  );
};
