/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useContext, useEffect, useRef, useState } from "react";
import VideoJS, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-vtt-thumbnails-freetube";
import "videojs-seek-buttons";
import "videojs-landscape-fullscreen";
import "./live";
import "./PlaylistButtons";
import "./source-selector";
import "./persist-volume";
import "./markers";
import "./big-buttons";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { ConfigurationContext } from "src/hooks/Config";
import { Interactive } from "src/utils/interactive";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

function handleHotkeys(player: VideoJsPlayer, event: VideoJS.KeyboardEvent) {
  function seekPercent(percent: number) {
    const duration = player.duration();
    const time = duration * percent;
    player.currentTime(time);
  }

  if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey) {
    return;
  }

  switch (event.which) {
    case 32: // space
    case 13: // enter
      if (player.paused()) player.play();
      else player.pause();
      break;
    case 77: // m
      player.muted(!player.muted());
      break;
    case 70: // f
      if (player.isFullscreen()) player.exitFullscreen();
      else player.requestFullscreen();
      break;
    case 39: // right arrow
      player.currentTime(Math.min(player.duration(), player.currentTime() + 5));
      break;
    case 37: // left arrow
      player.currentTime(Math.max(0, player.currentTime() - 5));
      break;
    case 38: // up arrow
      player.volume(player.volume() + 0.1);
      break;
    case 40: // down arrow
      player.volume(player.volume() - 0.1);
      break;
    case 48: // 0
      player.currentTime(0);
      break;
    case 49: // 1
      seekPercent(0.1);
      break;
    case 50: // 2
      seekPercent(0.2);
      break;
    case 51: // 3
      seekPercent(0.3);
      break;
    case 52: // 4
      seekPercent(0.4);
      break;
    case 53: // 5
      seekPercent(0.5);
      break;
    case 54: // 6
      seekPercent(0.6);
      break;
    case 55: // 7
      seekPercent(0.7);
      break;
    case 56: // 8
      seekPercent(0.8);
      break;
    case 57: // 9
      seekPercent(0.9);
      break;
  }
}

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment | undefined | null;
  timestamp: number;
  autoplay?: boolean;
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
  const sceneId = useRef<string | undefined>();
  const skipButtonsRef = useRef<any>();

  const [time, setTime] = useState(0);

  const [interactiveClient] = useState(
    new Interactive(config?.handyKey || "", config?.funscriptOffset || 0)
  );

  const [initialTimestamp] = useState(timestamp);

  const maxLoopDuration = config?.maximumLoopDuration ?? 0;

  useEffect(() => {
    if (playerRef.current && timestamp >= 0) {
      const player = playerRef.current;
      player.play()?.then(() => {
        player.currentTime(timestamp);
      });
    }
  }, [timestamp]);

  useEffect(() => {
    const videoElement = videoRef.current;
    if (!videoElement) return;

    const options: VideoJsPlayerOptions = {
      controls: true,
      controlBar: {
        pictureInPictureToggle: false,
        volumePanel: {
          inline: false,
        },
        chaptersButton: false,
      },
      nativeControlsForTouch: false,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      inactivityTimeout: 2000,
      preload: "none",
      userActions: {
        hotkeys: function (event) {
          const player = this as VideoJsPlayer;
          handleHotkeys(player, event);
        },
      },
    };

    const player = VideoJS(videoElement, options);

    (player as any).landscapeFullscreen({
      fullscreen: {
        enterOnRotate: true,
        exitOnRotate: true,
        alwaysInLandscapeMode: true,
        iOS: false,
      },
    });

    (player as any).markers();
    (player as any).offset();
    (player as any).sourceSelector();
    (player as any).persistVolume();
    (player as any).bigButtons();

    player.focus();
    playerRef.current = player;
  }, []);

  useEffect(() => {
    if (scene?.interactive) {
      interactiveClient.uploadScript(scene.paths.funscript || "");
    }
  }, [interactiveClient, scene?.interactive, scene?.paths.funscript]);

  useEffect(() => {
    if (skipButtonsRef.current) {
      skipButtonsRef.current.setForwardHandler(onNext);
      skipButtonsRef.current.setBackwardHandler(onPrevious);
    }
  }, [onNext, onPrevious]);

  useEffect(() => {
    const player = playerRef.current;
    if (player) {
      player.seekButtons({
        forward: 10,
        back: 10,
      });

      skipButtonsRef.current = player.skipButtons() ?? undefined;

      player.focus();
    }

    // Video player destructor
    return () => {
      if (playerRef.current) {
        playerRef.current.dispose();
        playerRef.current = undefined;
      }
    };
  }, []);

  useEffect(() => {
    function handleOffset(player: VideoJsPlayer) {
      if (!scene) return;

      const currentSrc = player.currentSrc();

      const isDirect =
        currentSrc.endsWith("/stream") || currentSrc.endsWith("/stream.m3u8");
      if (!isDirect) {
        (player as any).setOffsetDuration(scene.file.duration);
      } else {
        (player as any).clearOffsetDuration();
      }
    }

    function handleError(play: boolean) {
      const player = playerRef.current;
      if (!player) return;

      const currentFile = player.currentSource();
      if (currentFile) {
        // eslint-disable-next-line no-console
        console.log(`Source failed: ${currentFile.src}`);
        player.focus();
      }

      if (tryNextStream()) {
        // eslint-disable-next-line no-console
        console.log(`Trying next source in playlist: ${player.currentSrc()}`);
        player.load();
        if (play) {
          player.play();
        }
      } else {
        // eslint-disable-next-line no-console
        console.log("No more sources in playlist.");
      }
    }

    function tryNextStream() {
      const player = playerRef.current;
      if (!player) return;

      const sources = player.currentSources();

      if (sources.length > 1) {
        sources.shift();
        player.src(sources);
        return true;
      }

      return false;
    }

    if (!scene || scene.id === sceneId.current) return;
    sceneId.current = scene.id;

    const player = playerRef.current;
    if (!player) return;

    const auto =
      autoplay || (config?.autostartVideo ?? false) || initialTimestamp > 0;
    if (!auto && scene.paths?.screenshot) player.poster(scene.paths.screenshot);
    else player.poster("");

    // clear the offset before loading anything new.
    // otherwise, the offset will be applied to the next file when
    // currentTime is called.
    (player as any).clearOffsetDuration();

    const tracks = player.remoteTextTracks();
    if (tracks.length > 0) {
      player.removeRemoteTextTrack(tracks[0] as any);
    }

    player.src(
      scene.sceneStreams.map((stream) => ({
        src: stream.url,
        type: stream.mime_type ?? undefined,
        label: stream.label ?? undefined,
      }))
    );

    if (scene.paths.chapters_vtt) {
      player.addRemoteTextTrack(
        {
          src: scene.paths.chapters_vtt,
          kind: "chapters",
          default: true,
        },
        true
      );
    }

    player.currentTime(0);

    player.loop(
      !!scene.file.duration &&
        maxLoopDuration !== 0 &&
        scene.file.duration < maxLoopDuration
    );

    player.on("loadstart", function (this: VideoJsPlayer) {
      // handle offset after loading so that we get the correct current source
      handleOffset(this);
    });

    player.on("play", function (this: VideoJsPlayer) {
      player.poster("");
      if (scene.interactive) {
        interactiveClient.play(this.currentTime());
      }
    });

    player.on("pause", () => {
      if (scene.interactive) {
        interactiveClient.pause();
      }
    });

    player.on("timeupdate", function (this: VideoJsPlayer) {
      if (scene.interactive) {
        interactiveClient.ensurePlaying(this.currentTime());
      }

      setTime(this.currentTime());
    });

    player.on("seeking", function (this: VideoJsPlayer) {
      // backwards compatibility - may want to remove this in future
      this.play();
    });

    player.on("error", () => {
      handleError(true);
    });

    player.on("loadedmetadata", () => {
      if (!player.videoWidth() && !player.videoHeight()) {
        // Occurs during preload when videos with supported audio/unsupported video are preloaded.
        // Treat this as a decoding error and try the next source without playing.
        // However on Safari we get an media event when m3u8 is loaded which needs to be ignored.
        const currentFile = player.currentSrc();
        if (currentFile != null && !currentFile.includes("m3u8")) {
          // const play = !player.paused();
          // handleError(play);
          player.error(MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED);
        }
      }
    });

    player.load();

    if (auto) {
      player
        .play()
        ?.then(() => {
          if (initialTimestamp > 0) {
            player.currentTime(initialTimestamp);
          }
        })
        .catch(() => {
          if (scene.paths.screenshot) player.poster(scene.paths.screenshot);
        });
    }

    if ((player as any).vttThumbnails?.src)
      (player as any).vttThumbnails?.src(scene?.paths.vtt);
    else
      (player as any).vttThumbnails({
        src: scene?.paths.vtt,
        showTimestamp: true,
      });
  }, [
    scene,
    config?.autostartVideo,
    maxLoopDuration,
    initialTimestamp,
    autoplay,
    interactiveClient,
  ]);

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
          playsInline
          ref={videoRef}
          id={VIDEO_PLAYER_ID}
          className="video-js vjs-big-play-centered"
        />
      </div>
      {scene && (
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
