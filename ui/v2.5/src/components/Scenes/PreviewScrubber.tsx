import React, {
  useRef,
  useMemo,
  useState,
  useLayoutEffect,
  useEffect,
} from "react";
import { useSpriteInfo } from "src/hooks/sprite";
import { useThrottle } from "src/hooks/throttle";
import TextUtils from "src/utils/text";
import { HoverScrubber } from "../Shared/HoverScrubber";

interface IScenePreviewProps {
  vttPath: string | undefined;
  onClick?: (timestamp: number) => void;
  disabled?: boolean;
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

const defaultSprites = 81; // 9x9 grid by default

export const PreviewScrubber: React.FC<IScenePreviewProps> = ({
  vttPath,
  onClick,
  disabled,
}) => {
  const imageParentRef = useRef<HTMLDivElement>(null);
  const [style, setStyle] = useState({});

  const [activeIndex, setActiveIndex] = useState<number>();

  const debounceSetActiveIndex = useThrottle(setActiveIndex, 50);

  // hold off on loading vtt until first mouse over
  const [hasLoaded, setHasLoaded] = useState(false);
  const spriteInfo = useSpriteInfo(hasLoaded ? vttPath : undefined);

  const spriteSheetSize = useMemo(() => {
    if (!spriteInfo) {
      return { x: 0, y: 0 };
    }

    // calculate total width/height of scrubber image so we can scale it
    const maxX = Math.max(...spriteInfo.map((sprite) => sprite.x + sprite.w));
    const maxY = Math.max(...spriteInfo.map((sprite) => sprite.y + sprite.h));

    return { x: maxX, y: maxY };
  }, [spriteInfo]);

  const sprite = useMemo(() => {
    if (!spriteInfo || activeIndex === undefined) {
      return undefined;
    }
    return spriteInfo[activeIndex];
  }, [activeIndex, spriteInfo]);

  // mark as loaded on the first hover
  useEffect(() => {
    if (activeIndex !== undefined) {
      setHasLoaded(true);
    }
  }, [activeIndex]);

  useLayoutEffect(() => {
    const imageParent = imageParentRef.current;

    if (!sprite || !imageParent) {
      return setStyle({});
    }

    const clientRect = imageParent.getBoundingClientRect();
    const scale = scaleToFit(sprite, clientRect);

    setStyle({
      backgroundPosition: `${-sprite.x * scale}px ${-sprite.y * scale}px`,
      backgroundImage: `url(${sprite.url})`,
      backgroundSize: `${spriteSheetSize.x * scale}px ${
        spriteSheetSize.y * scale
      }px`,
      width: `${sprite.w * scale}px`,
      height: `${sprite.h * scale}px`,
    });
  }, [sprite, spriteSheetSize]);

  const currentTime = useMemo(() => {
    if (!sprite) return undefined;

    const start = TextUtils.secondsToTimestamp(sprite.start);

    return start;
  }, [sprite]);

  function onScrubberClick(index: number) {
    if (!onClick || !spriteInfo) {
      return;
    }

    const s = spriteInfo[index];
    onClick(s.start);
  }

  if (spriteInfo === null || !vttPath) return null;

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
        totalSprites={spriteInfo?.length ?? defaultSprites}
        activeIndex={activeIndex}
        setActiveIndex={(i) => debounceSetActiveIndex(i)}
        onClick={onScrubberClick}
        disabled={disabled}
      />
    </div>
  );
};
