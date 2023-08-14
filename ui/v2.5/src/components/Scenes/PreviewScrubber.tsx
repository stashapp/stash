import React, { useMemo } from "react";
import { useDebounce } from "src/hooks/debounce";
import { useSpriteInfo } from "src/hooks/sprite";

interface IHoverScrubber {
  totalSprites: number;
  activeIndex: number | undefined;
  setActiveIndex: (index: number | undefined) => void;
}

const HoverScrubber: React.FC<IHoverScrubber> = ({
  totalSprites,
  activeIndex,
  setActiveIndex,
}) => {
  function onMouseMove(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    const relatedTarget = e.currentTarget;

    if (relatedTarget !== e.target) return;

    const { width } = relatedTarget.getBoundingClientRect();
    const x = e.nativeEvent.offsetX;

    const index = Math.floor((x / width) * totalSprites);
    setActiveIndex(index);
  }

  function onMouseLeave() {
    setActiveIndex(undefined);
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
      ></div>
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
}

export const PreviewScrubber: React.FC<IScenePreviewProps> = ({ vttPath }) => {
  const [activeIndex, setActiveIndex] = React.useState<number | undefined>();

  const debounceSetActiveIndex = useDebounce(
    setActiveIndex,
    [setActiveIndex],
    10
  );

  const spriteInfo = useSpriteInfo(vttPath);

  const style = useMemo(() => {
    if (!spriteInfo || activeIndex === undefined) {
      return {};
    }

    const sprite = spriteInfo[activeIndex];
    const totalWidth = spriteInfo.reduce(
      (acc, cur) => Math.max(acc, cur.x + cur.w),
      0
    );
    const totalHeight = spriteInfo.reduce(
      (acc, cur) => Math.max(acc, cur.y + cur.h),
      0
    );

    const spriteX = sprite.x / totalWidth;
    const spriteY = sprite.y / totalHeight;

    const spritesX = Math.floor(totalWidth / sprite.w);
    const spritesY = Math.floor(totalHeight / sprite.h);

    return {
      "background-size": `calc(100% * ${spritesX}) calc(100% * ${spritesY})`,
      backgroundPosition: `calc(${-spriteX} * 100% * ${spritesX}) calc(${-spriteY} * 100% * ${spritesY})`,
      backgroundImage: `url(${sprite.url})`,
    };
  }, [spriteInfo, activeIndex]);

  if (!spriteInfo) return null;

  return (
    <div className="preview-scrubber">
      {activeIndex !== undefined && spriteInfo && (
        <div className="scene-card-preview-image">
          <div className="scrubber-image" style={style}></div>
        </div>
      )}
      <HoverScrubber
        totalSprites={81}
        activeIndex={activeIndex}
        setActiveIndex={(i) => debounceSetActiveIndex(i)}
      />
    </div>
  );
};
