/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useContext, useEffect, useRef, useState } from "react";
import VideoJS, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-vtt-thumbnails-freetube";
import "videojs-seek-buttons";
import "./landscapeFullscreen";
import "./live";
import "./PlaylistButtons";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { ConfigurationContext } from "src/hooks/Config";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment | undefined | null;
  timestamp: number;
  autoplay?: boolean;
  onReady?: () => void;
  onSeeked?: () => void;
  onTime?: () => void;
  onComplete?: () => void;
  onNext?: () => void;
  onPrevious?: () => void;
}

export const ScenePlayer: React.FC<IScenePlayerProps> = ({
  className,
  autoplay,
  scene,
  timestamp,
  onComplete,
  onNext,
  onPrevious,
}) => {
  const { configuration } = useContext(ConfigurationContext);
  const config = configuration?.interface;
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<VideoJsPlayer | undefined>();
  const [time, setTime] = useState(0);

  const maxLoopDuration = config?.maximumLoopDuration ?? 0;

  useEffect(() => {
    if (playerRef.current) playerRef.current.currentTime(timestamp);
  }, [timestamp]);

  useEffect(() => {
    const videoElement = videoRef.current;
    if (!videoElement) return;

    const options: VideoJsPlayerOptions = {
      autoplay: false,
      controls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      inactivityTimeout: 2000,
    };

    const player = VideoJS(videoElement, options);
    player.seekButtons({
      forward: 10,
      back: 10,
    });
    const skipButtons = player.skipButtons();
    if (skipButtons) {
      skipButtons.setForwardHandler(onNext);
      skipButtons.setBackwardHandler(onPrevious);
    }
    player.on("error", () => {
      player.error(null);
    });
    (player as any).offset();
    (player as any).landscapeFullscreen({
      fullscreen: {
        enterOnRotate: true,
        exitOnRotate: true,
        alwaysInLandscapeMode: true,
        iOS: true
      }
    });

    playerRef.current = player;
  }, [playerRef]);

  useEffect(() => {
    const skipButtons = playerRef?.current?.skipButtons?.();
    if (skipButtons) {
      skipButtons.setForwardHandler(onNext);
      skipButtons.setBackwardHandler(onPrevious);
    }
  }, [onNext, onPrevious]);

  useEffect(() => {
    // Video player destructor
    return () => {
      if (playerRef.current) {
        playerRef.current.dispose();
        playerRef.current = undefined;
      }
    };
  }, []);

  useEffect(() => {
    if (!scene) return;

    const player = playerRef.current;
    if (!player) return;

    const auto = autoplay || (config?.autostartVideo ?? false) || timestamp > 0;
    if (!auto && scene.paths?.screenshot)
      player.poster(scene.paths.screenshot);
    else
      player.poster('');
    player.src(
      scene.sceneStreams.map((stream) => ({
        src: stream.url,
        type: stream.mime_type ?? undefined,
      }))
    );
    player.currentTime(0);
    if (auto)
      player.play();

    player.loop(!!scene.file.duration && maxLoopDuration !== 0 && scene.file.duration < maxLoopDuration);

    const isDirect = new URL(scene.sceneStreams[0].url).pathname.endsWith("/stream");
    if (!isDirect) {
      (player as any).setOffsetDuration(scene.file.duration);
    } else {
      (player as any).clearOffsetDuration();
    };

    player.on("timeupdate", function (this: VideoJsPlayer) {
      setTime(this.currentTime());
    });
    player.on("seeking", function (this: VideoJsPlayer) {
      if (!isDirect) {
        // remove the start parameter
        const srcUrl = new URL(player.src());
        srcUrl.searchParams.delete("start");

        /* eslint-disable no-param-reassign */
        srcUrl.searchParams.append("start", player.currentTime().toString());
        player.src(srcUrl.toString());
        /* eslint-enable no-param-reassign */

        player.play();
      }
    });
    player.play()?.catch(() => {
      if (scene.paths.screenshot)
        player.poster(scene.paths.screenshot);
    });

    if ((player as any).vttThumbnails?.src)
      (player as any).vttThumbnails?.src(scene?.paths.vtt);
    else
      (player as any).vttThumbnails({
        src: scene?.paths.vtt,
        showTimestamp: true,
      });
  }, [scene]);

  useEffect(() => {
    // Attach handler for onComplete event
    const player = playerRef.current;
    if (!player) return;

    player.on("ended", () => {
      onComplete?.();
    });

    return () => player.off("ended");
  }, [onComplete]);


  const onScrubberScrolled = () => {
    playerRef.current?.pause();
  };
  const onScrubberSeek = (seconds: number) => {
    playerRef.current?.currentTime(seconds);
  };

  const isPortrait =
    scene &&
    scene.file.height &&
    scene.file.width &&
    scene.file.height > scene.file.width;

  return (
    <div className={cx("VideoPlayer", { portrait: isPortrait })}>
      <div data-vjs-player className={cx("video-wrapper", className)}>
        <video
          ref={videoRef}
          id={VIDEO_PLAYER_ID}
          className="video-js vjs-big-play-centered"
        />
      </div>
      { scene && (
        <ScenePlayerScrubber
          scene={scene}
          position={time}
          onSeek={onScrubberSeek}
          onScrolled={onScrubberScrolled}
        />
      )}
    </div>
  );
};

export const getPlayerPosition = () =>
  VideoJS.getPlayer(VIDEO_PLAYER_ID).currentTime();
