import React, { useEffect, useState } from "react";
import { useThrottle } from "src/hooks/throttle";
import { HoverScrubber } from "../Shared/HoverScrubber";
import cx from "classnames";

export const GalleryPreviewScrubber: React.FC<{
  className?: string;
  previewPath: string;
  defaultPath: string;
  imageCount: number;
  onClick?: (imageIndex: number) => void;
  onPathChanged: React.Dispatch<React.SetStateAction<string | undefined>>;
}> = ({
  className,
  previewPath,
  defaultPath,
  imageCount,
  onClick,
  onPathChanged,
}) => {
  const [activeIndex, setActiveIndex] = useState<number>();
  const debounceSetActiveIndex = useThrottle(setActiveIndex, 50);

  function onScrubberClick(index: number) {
    if (!onClick) {
      return;
    }

    onClick(index);
  }

  useEffect(() => {
    function getPath() {
      if (activeIndex === undefined) {
        return defaultPath;
      }

      return `${previewPath}/${activeIndex}`;
    }

    onPathChanged(getPath());
  }, [activeIndex, defaultPath, previewPath, onPathChanged]);

  return (
    <div className={cx("preview-scrubber", className)}>
      <HoverScrubber
        totalSprites={imageCount}
        activeIndex={activeIndex}
        setActiveIndex={(i) => debounceSetActiveIndex(i)}
        onClick={onScrubberClick}
      />
    </div>
  );
};
