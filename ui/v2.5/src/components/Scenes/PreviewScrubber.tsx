import React, { useMemo } from "react";
import { useDebounce } from "src/hooks/debounce";
import { useSpriteInfo } from "src/hooks/sprite";
import TextUtils from "src/utils/text";

interface IHoverScrubber {
  totalSprites: number;
  activeIndex: number | undefined;
  setActiveIndex: (index: number | undefined) => void;
  onClick?: (index: number) => void;
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

    return Math.floor((x / width) * (totalSprites - 1));
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
    onClick(getActiveIndex(e));
  }

  const indicatorStyle = useMemo(() => {
    if (activeIndex === undefined) return {};

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
  const imageParentRef = React.useRef<HTMLDivElement>(null);

  const [activeIndex, setActiveIndex] = React.useState<number | undefined>();

  const debounceSetActiveIndex = useDebounce(
    setActiveIndex,
    [setActiveIndex],
    1
  );

  const spriteInfo = useSpriteInfo(vttPath);

  const style = useMemo(() => {
    if (!spriteInfo || activeIndex === undefined || !imageParentRef.current) {
      return {};
    }

    const sprite = spriteInfo[activeIndex];

    const clientRect = imageParentRef.current?.getBoundingClientRect();
    const scale = clientRect ? scaleToFit(sprite, clientRect) : 1;

    return {
      backgroundPosition: `${-sprite.x}px ${-sprite.y}px`,
      backgroundImage: `url(${sprite.url})`,
      width: `${sprite.w}px`,
      height: `${sprite.h}px`,
      transform: `scale(${scale})`,
    };
  }, [spriteInfo, activeIndex, imageParentRef]);

  const currentTime = useMemo(() => {
    if (!spriteInfo || activeIndex === undefined) {
      return undefined;
    }

    const sprite = spriteInfo[activeIndex];

    const start = TextUtils.secondsToTimestamp(sprite.start);

    return start;
  }, [activeIndex, spriteInfo]);

  function onScrubberClick(index: number) {
    if (!spriteInfo || !onClick) {
      return;
    }

    const sprite = spriteInfo[index];

    onClick(sprite.start);
  }

  if (!spriteInfo) return null;

  return (
    <div className="preview-scrubber">
      {activeIndex !== undefined && spriteInfo && (
        <div className="scene-card-preview-image" ref={imageParentRef}>
          <div className="scrubber-image" style={style}></div>
          {currentTime !== undefined && (
            <div className="scrubber-timestamp">{currentTime}</div>
          )}
        </div>
      )}
      <HoverScrubber
        totalSprites={81}
        activeIndex={activeIndex}
        setActiveIndex={(i) => debounceSetActiveIndex(i)}
        onClick={onScrubberClick}
      />
    </div>
  );
};
