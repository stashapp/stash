import { useEffect, useRef } from "react";
import { StashService } from "../core/StashService";

export interface IVideoHoverHookData {
  videoEl: React.RefObject<HTMLVideoElement>;
  isPlaying: React.MutableRefObject<boolean>;
  isHovering: React.MutableRefObject<boolean>;
  options: IVideoHoverHookOptions;
}

export interface IVideoHoverHookOptions {
  resetOnMouseLeave: boolean;
}

export class VideoHoverHook {
  public static useVideoHover(options: IVideoHoverHookOptions): IVideoHoverHookData {
    const videoEl = useRef<HTMLVideoElement>(null);
    const isPlaying = useRef<boolean>(false);
    const isHovering = useRef<boolean>(false);

    const config = StashService.useConfiguration();
    const soundEnabled = !!config.data && !!config.data.configuration ? config.data.configuration.interface.soundOnPreview : true;

    useEffect(() => {
      const videoTag = videoEl.current;
      if (!videoTag) { return; }
      videoTag.onplaying = () => {
        if (isHovering.current === true) {
          isPlaying.current = true;
        } else {
          videoTag.pause();
        }
      };
      videoTag.onpause = () => isPlaying.current = false;
    }, [videoEl]);

    useEffect(() => {
      const videoTag = videoEl.current;
      if (!videoTag) { return; }
      videoTag.volume = soundEnabled ? 0.05 : 0;
    }, [soundEnabled]);

    return {videoEl, isPlaying, isHovering, options};
  }

  public static onMouseEnter(data: IVideoHoverHookData) {
    data.isHovering.current = true;

    const videoTag = data.videoEl.current;
    if (!videoTag) { return; }
    if (videoTag.paused && !data.isPlaying.current) {
      videoTag.play().catch((error) => {
        console.log(error.message);
      });
    }
  }

  public static onMouseLeave(data: IVideoHoverHookData) {
    data.isHovering.current = false;

    const videoTag = data.videoEl.current;
    if (!videoTag) { return; }
    if (!videoTag.paused && data.isPlaying) {
      videoTag.pause();
      if (data.options.resetOnMouseLeave) {
        videoTag.removeAttribute("src");
        videoTag.load();
        data.isPlaying.current = false;
      }
    }
  }
}
