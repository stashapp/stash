import { useEffect, useRef } from "react";
import { useConfiguration } from "../core/StashService";

export interface IVideoHoverHookData {
  videoEl: React.RefObject<HTMLVideoElement>;
  isPlaying: React.MutableRefObject<boolean>;
  isHovering: React.MutableRefObject<boolean>;
  options: IVideoHoverHookOptions;
}

export interface IVideoHoverHookOptions {
  resetOnMouseLeave: boolean;
}

export const useVideoHover = (options: IVideoHoverHookOptions) => {
  const videoEl = useRef<HTMLVideoElement>(null);
  const isPlaying = useRef<boolean>(false);
  const isHovering = useRef<boolean>(false);
  const config = useConfiguration();

  const onMouseEnter = () => {
    isHovering.current = true;

    const videoTag = videoEl.current;
    if (!videoTag) {
      return;
    }
    if (videoTag.paused && !isPlaying.current) {
      videoTag.play().catch(() => {});
    }
  };

  const onMouseLeave = () => {
    isHovering.current = false;

    const videoTag = videoEl.current;
    if (!videoTag) {
      return;
    }
    if (!videoTag.paused && isPlaying) {
      videoTag.pause();
      if (options.resetOnMouseLeave) {
        videoTag.removeAttribute("src");
        videoTag.load();
        isPlaying.current = false;
      }
    }
  };

  const soundEnabled =
    config?.data?.configuration?.interface?.soundOnPreview ?? true;

  useEffect(() => {
    const videoTag = videoEl.current;
    if (!videoTag) {
      return;
    }
    videoTag.onplaying = () => {
      if (isHovering.current === true) {
        isPlaying.current = true;
      } else {
        videoTag.pause();
      }
    };
    videoTag.onpause = () => {
      isPlaying.current = false;
    };
  }, [videoEl]);

  useEffect(() => {
    const videoTag = videoEl.current;
    if (!videoTag) {
      return;
    }
    videoTag.volume = soundEnabled ? 0.05 : 0;
  }, [soundEnabled]);

  return {
    videoEl,
    isPlaying,
    isHovering,
    options,
    onMouseEnter,
    onMouseLeave,
  };
};
