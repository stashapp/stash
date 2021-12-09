/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useContext, useEffect, useRef, useState } from "react";
import VideoJS, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-vtt-thumbnails-freetube";
import "videojs-seek-buttons";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { ConfigurationContext } from "src/hooks/Config";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment;
  sceneStreams: GQL.SceneStreamEndpoint[];
  timestamp: number;
  autoplay?: boolean;
  onReady?: () => void;
  onSeeked?: () => void;
  onTime?: () => void;
  onComplete?: () => void;
}

export const ScenePlayer: React.FC<IScenePlayerProps> = ({
  className,
  autoplay,
  scene,
  sceneStreams,
  timestamp,
}) => {
  const { configuration } = useContext(ConfigurationContext);
  const config = configuration?.interface;
  const videoRef = useRef(null);
  const playerRef = useRef<VideoJsPlayer | undefined>();
  const [time, setTime] = useState(0);

  useEffect(() => {
    if (playerRef.current) playerRef.current.currentTime(timestamp);
  }, [timestamp]);

  useEffect(() => {
    const videoElement = videoRef.current;
    if (!videoElement) return;

    const options: VideoJsPlayerOptions = {
      autoplay: autoplay || (config?.autostartVideo ?? false) || timestamp > 0,
      poster: scene.paths.screenshot ?? undefined,
      controls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
    };

    const player = VideoJS(videoElement, options);
    player.src(
      sceneStreams.map((stream) => ({
        src: stream.url,
        type: stream.mime_type ?? undefined,
      }))
    );

    (player as any).vttThumbnails({
      src: scene.paths.vtt,
      showTimestamp: true,
    });
    player.seekButtons({
      forward: 10,
      back: 10,
    });
    player.on("timeupdate", function (this: VideoJsPlayer) {
      setTime(this.currentTime());
    });
    player.on("loadeddata", function (this: VideoJsPlayer) {
      if (timestamp > 0) this.currentTime(timestamp);
    });

    playerRef.current = player;

    return () => {
      if (playerRef.current) {
        playerRef.current.dispose();
        playerRef.current = undefined;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [playerRef]);

  const onScrubberScrolled = () => {
    playerRef.current?.pause();
  };
  const onScrubberSeek = (seconds: number) => {
    playerRef.current?.currentTime(seconds);
  };

  const isPortrait =
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
      <ScenePlayerScrubber
        scene={scene}
        position={time}
        onSeek={onScrubberSeek}
        onScrolled={onScrubberScrolled}
      />
    </div>
  );
};

export const getPlayerPosition = () =>
  VideoJS.getPlayer(VIDEO_PLAYER_ID).currentTime();
