/* eslint-disable @typescript-eslint/no-explicit-any */
import React, {
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
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
import {
  ConnectionState,
  InteractiveContext,
} from "src/hooks/Interactive/context";
import { SceneInteractiveStatus } from "src/hooks/Interactive/status";
import { languageMap } from "src/utils/caption";

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

  const {
    interactive: interactiveClient,
    uploadScript,
    currentScript,
    initialised: interactiveInitialised,
    state: interactiveState,
  } = React.useContext(InteractiveContext);

  const [initialTimestamp] = useState(timestamp);
  const [ready, setReady] = useState(false);
  const started = useRef(false);
  const interactiveReady = useRef(false);

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
      playbackRates: [0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
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

    const settings = (player as any).textTrackSettings;
    settings.setValues({
      backgroundColor: "#000",
      backgroundOpacity: "0.5",
    });
    settings.updateDisplay();

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
    if (scene?.interactive && interactiveInitialised) {
      interactiveReady.current = false;
      uploadScript(scene.paths.funscript || "").then(() => {
        interactiveReady.current = true;
      });
    }
  }, [
    uploadScript,
    interactiveInitialised,
    scene?.interactive,
    scene?.paths.funscript,
  ]);

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

  const start = useCallback(() => {
    const player = playerRef.current;
    if (player && scene) {
      started.current = true;

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
  }, [scene, initialTimestamp]);

  useEffect(() => {
    let prevCaptionOffset = 0;

    function addCaptionOffset(player: VideoJsPlayer, offset: number) {
      const tracks = player.remoteTextTracks();
      for (let i = 0; i < tracks.length; i++) {
        const track = tracks[i];
        const { cues } = track;
        if (cues) {
          for (let j = 0; j < cues.length; j++) {
            const cue = cues[j];
            cue.startTime = cue.startTime + offset;
            cue.endTime = cue.endTime + offset;
          }
        }
      }
    }

    function removeCaptionOffset(player: VideoJsPlayer, offset: number) {
      const tracks = player.remoteTextTracks();
      for (let i = 0; i < tracks.length; i++) {
        const track = tracks[i];
        const { cues } = track;
        if (cues) {
          for (let j = 0; j < cues.length; j++) {
            const cue = cues[j];
            cue.startTime = cue.startTime + prevCaptionOffset - offset;
            cue.endTime = cue.endTime + prevCaptionOffset - offset;
          }
        }
      }
    }

    function handleOffset(player: VideoJsPlayer) {
      if (!scene) return;

      const currentSrc = player.currentSrc();

      const isDirect =
        currentSrc.endsWith("/stream") || currentSrc.endsWith("/stream.m3u8");

      const curTime = player.currentTime();
      if (!isDirect) {
        (player as any).setOffsetDuration(scene.file.duration);
      } else {
        (player as any).clearOffsetDuration();
      }

      if (curTime != prevCaptionOffset) {
        if (!isDirect) {
          removeCaptionOffset(player, curTime);
          prevCaptionOffset = curTime;
        } else {
          if (prevCaptionOffset != 0) {
            addCaptionOffset(player, prevCaptionOffset);
            prevCaptionOffset = 0;
          }
        }
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

    function getDefaultLanguageCode() {
      var languageCode = window.navigator.language;

      if (languageCode.indexOf("-") !== -1) {
        languageCode = languageCode.split("-")[0];
      }

      if (languageCode.indexOf("_") !== -1) {
        languageCode = languageCode.split("_")[0];
      }

      return languageCode;
    }

    function loadCaptions(player: VideoJsPlayer) {
      if (!scene) return;

      if (scene.captions) {
        var languageCode = getDefaultLanguageCode();
        var hasDefault = false;

        for (let caption of scene.captions) {
          var lang = caption.language_code;
          var label = lang;
          if (languageMap.has(lang)) {
            label = languageMap.get(lang)!;
          }

          label = label + " (" + caption.caption_type + ")";
          var setAsDefault = !hasDefault && languageCode == lang;
          if (!hasDefault && setAsDefault) {
            hasDefault = true;
          }
          player.addRemoteTextTrack(
            {
              src:
                scene.paths.caption +
                "?lang=" +
                lang +
                "&type=" +
                caption.caption_type,
              kind: "captions",
              srclang: lang,
              label: label,
              default: setAsDefault,
            },
            true
          );
        }
      }
    }

    // always stop the interactive client on initialisation
    interactiveClient.pause();
    interactiveReady.current = false;

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
    for (let i = 0; i < tracks.length; i++) {
      player.removeRemoteTextTrack(tracks[i] as any);
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

    if (scene.captions?.length! > 0) {
      loadCaptions(player);
    }

    player.currentTime(0);

    const looping =
      !!scene.file.duration &&
      maxLoopDuration !== 0 &&
      scene.file.duration < maxLoopDuration;
    player.loop(looping);
    interactiveClient.setLooping(looping);

    function loadstart(this: VideoJsPlayer) {
      // handle offset after loading so that we get the correct current source
      handleOffset(this);
    }

    player.on("loadstart", loadstart);

    function onPlay(this: VideoJsPlayer) {
      this.poster("");
      if (scene?.interactive && interactiveReady.current) {
        interactiveClient.play(this.currentTime());
      }
    }
    player.on("play", onPlay);

    function pause() {
      interactiveClient.pause();
    }
    player.on("pause", pause);

    function timeupdate(this: VideoJsPlayer) {
      if (scene?.interactive && interactiveReady.current) {
        interactiveClient.ensurePlaying(this.currentTime());
      }
      setTime(this.currentTime());
    }
    player.on("timeupdate", timeupdate);

    function seeking(this: VideoJsPlayer) {
      this.play();
    }
    player.on("seeking", seeking);

    function error() {
      handleError(true);
    }
    player.on("error", error);

    // changing source (eg when seeking) resets the playback rate
    // so set the default in addition to the current rate
    function ratechange(this: VideoJsPlayer) {
      this.defaultPlaybackRate(this.playbackRate());
    }
    player.on("ratechange", ratechange);

    function loadedmetadata(this: VideoJsPlayer) {
      if (!this.videoWidth() && !this.videoHeight()) {
        // Occurs during preload when videos with supported audio/unsupported video are preloaded.
        // Treat this as a decoding error and try the next source without playing.
        // However on Safari we get an media event when m3u8 is loaded which needs to be ignored.
        const currentFile = this.currentSrc();
        if (currentFile != null && !currentFile.includes("m3u8")) {
          // const play = !player.paused();
          // handleError(play);
          this.error(MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED);
        }
      }
    }
    player.on("loadedmetadata", loadedmetadata);

    player.load();

    if ((player as any).vttThumbnails?.src)
      (player as any).vttThumbnails?.src(scene?.paths.vtt);
    else
      (player as any).vttThumbnails({
        src: scene?.paths.vtt,
        showTimestamp: true,
      });

    setReady(true);
    started.current = false;

    return () => {
      setReady(false);

      // stop the interactive client
      interactiveClient.pause();

      player.off("loadstart", loadstart);
      player.off("play", onPlay);
      player.off("pause", pause);
      player.off("timeupdate", timeupdate);
      player.off("seeking", seeking);
      player.off("error", error);
      player.off("ratechange", ratechange);
      player.off("loadedmetadata", loadedmetadata);
    };
  }, [
    scene,
    config?.autostartVideo,
    maxLoopDuration,
    initialTimestamp,
    autoplay,
    interactiveClient,
    start,
  ]);

  useEffect(() => {
    if (!ready || started.current) {
      return;
    }

    const auto =
      autoplay || (config?.autostartVideo ?? false) || initialTimestamp > 0;

    // check if we're waiting for the interactive client
    const interactiveWaiting =
      scene?.interactive &&
      interactiveClient.handyKey &&
      currentScript !== scene.paths.funscript;

    if (scene && auto && !interactiveWaiting) {
      start();
    }
  }, [
    config?.autostartVideo,
    initialTimestamp,
    scene,
    ready,
    interactiveClient,
    currentScript,
    autoplay,
    start,
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
    const player = playerRef.current;
    if (player) {
      player.play()?.then(() => {
        player.currentTime(seconds);
      });
    }
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
      {scene?.interactive &&
        (interactiveState !== ConnectionState.Ready ||
          playerRef.current?.paused()) && <SceneInteractiveStatus />}
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
