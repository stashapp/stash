import React, {
  KeyboardEvent,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import videojs, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-seek-buttons";
import "videojs-landscape-fullscreen";
import "./live";
import "./PlaylistButtons";
import "./source-selector";
import "./persist-volume";
import "./markers";
import "./vtt-thumbnails";
import "./big-buttons";
import cx from "classnames";
import {
  useSceneSaveActivity,
  useSceneIncrementPlayCount,
} from "src/core/StashService";

import * as GQL from "src/core/generated-graphql";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { ConfigurationContext } from "src/hooks/Config";
import {
  ConnectionState,
  InteractiveContext,
} from "src/hooks/Interactive/context";
import { SceneInteractiveStatus } from "src/hooks/Interactive/status";
import { languageMap } from "src/utils/caption";
import { VIDEO_PLAYER_ID } from "./util";
import { IUIConfig } from "src/core/config";

function handleHotkeys(player: VideoJsPlayer, event: videojs.KeyboardEvent) {
  function seekPercent(percent: number) {
    const duration = player.duration();
    const time = duration * percent;
    player.currentTime(time);
  }

  function seekPercentRelative(percent: number) {
    const duration = player.duration();
    const currentTime = player.currentTime();
    const time = currentTime + duration * percent;
    if (time > duration) return;
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
    case 221: // ]
      seekPercentRelative(0.1);
      break;
    case 219: // [
      seekPercentRelative(-0.1);
      break;
  }
}

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment | undefined | null;
  hideScrubberOverride: boolean;
  autoplay?: boolean;
  permitLoop?: boolean;
  initialTimestamp: number;
  sendSetTimestamp: (setTimestamp: (value: number) => void) => void;
  onComplete: () => void;
  onNext: () => void;
  onPrevious: () => void;
}

export const ScenePlayer: React.FC<IScenePlayerProps> = ({
  className,
  scene,
  hideScrubberOverride,
  autoplay,
  permitLoop = true,
  initialTimestamp: _initialTimestamp,
  sendSetTimestamp,
  onComplete,
  onNext,
  onPrevious,
}) => {
  const { configuration } = useContext(ConfigurationContext);
  const interfaceConfig = configuration?.interface;
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<VideoJsPlayer>();
  const sceneId = useRef<string>();
  const [sceneSaveActivity] = useSceneSaveActivity();
  const [sceneIncrementPlayCount] = useSceneIncrementPlayCount();

  const [time, setTime] = useState(0);
  const [ready, setReady] = useState(false);
  const [sessionInitialised, setSessionInitialised] = useState(false); // tracks play session. This is reset whenever ScenePlayer page is exited

  const {
    interactive: interactiveClient,
    uploadScript,
    currentScript,
    initialised: interactiveInitialised,
    state: interactiveState,
  } = React.useContext(InteractiveContext);

  const [fullscreen, setFullscreen] = useState(false);
  const [showScrubber, setShowScrubber] = useState(false);

  const initialTimestamp = useRef(-1);
  const started = useRef(false);
  const auto = useRef(false);
  const interactiveReady = useRef(false);

  const playDurationRef = useRef(0);
  const trackTime = useRef(false);
  const recordedActivity = useRef(false);
  const [updatePlayDuration, setUpdatePlayDuration] = useState(false);

  const ignoreInterval = uiConfig?.ignoreInterval ?? 0;
  const trackActivity = uiConfig?.trackActivity ?? false;

  const file = useMemo(
    () => ((scene?.files.length ?? 0) > 0 ? scene?.files[0] : undefined),
    [scene]
  );

  const maxLoopDuration = interfaceConfig?.maximumLoopDuration ?? 0;
  const looping = useMemo(
    () =>
      !!file?.duration &&
      permitLoop &&
      maxLoopDuration !== 0 &&
      file.duration < maxLoopDuration,
    [file, permitLoop, maxLoopDuration]
  );

  useEffect(() => {
    if (hideScrubberOverride || fullscreen) {
      setShowScrubber(false);
      return;
    }

    const onResize = () => {
      const show = window.innerHeight >= 450 && window.innerWidth >= 576;
      setShowScrubber(show);
    };
    onResize();

    window.addEventListener("resize", onResize);

    return () => window.removeEventListener("resize", onResize);
  }, [hideScrubberOverride, fullscreen]);

  useEffect(() => {
    sendSetTimestamp((value: number) => {
      const player = playerRef.current;
      if (player && value >= 0) {
        player.play()?.then(() => {
          player.currentTime(value);
        });
      }
    });
  }, [sendSetTimestamp]);

  useEffect(() => {
    const id = sceneId.current;
    if (trackActivity && updatePlayDuration && id && playerRef.current) {
      const playDuration = playDurationRef.current;
      let resume_time = playerRef.current.currentTime()!;
      const videoDuration = playerRef.current.duration();
      const percentPlayed = (100 / videoDuration) * playDuration;
      const percentCompleted = (100 / videoDuration) * resume_time;
      if (!recordedActivity.current && percentPlayed >= ignoreInterval) {
        sceneIncrementPlayCount({
          variables: {
            id,
          },
        });
        recordedActivity.current = true;
      }
      // if video is 98% or more complete then reset resume_time
      if (percentCompleted >= 98) {
        resume_time = 0;
      }
      sceneSaveActivity({
        variables: {
          id,
          resume_time,
          playDuration,
        },
      });
    }
  }, [
    updatePlayDuration,
    ignoreInterval,
    sceneIncrementPlayCount,
    sceneSaveActivity,
    trackActivity,
  ]);

  // Initialize VideoJS player
  useEffect(() => {
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
        hotkeys: function (this: VideoJsPlayer, event) {
          handleHotkeys(this, event);
        },
      },
      plugins: {
        vttThumbnails: {
          showTimestamp: true,
        },
        markers: {},
        sourceSelector: {},
        persistVolume: {},
        bigButtons: {},
        seekButtons: {
          forward: 10,
          back: 10,
        },
        skipButtons: {},
      },
    };

    const player = videojs(videoRef.current!, options);
    var playDurationHandler = window.setInterval(() => {
      if (trackTime.current) {
        playDurationRef.current++;
        if (playDurationRef.current % 10 == 0) {
          setUpdatePlayDuration(true);
        } else {
          setUpdatePlayDuration(false);
        }
      }
    }, 1000); // when scene is playing incrememt playDuration every second

    /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
    const settings = (player as any).textTrackSettings;
    settings.setValues({
      backgroundColor: "#000",
      backgroundOpacity: "0.5",
    });
    settings.updateDisplay();

    player.focus();
    playerRef.current = player;

    // Video player destructor
    return () => {
      if (playerRef.current) {
        clearInterval(playDurationHandler);
        if (trackActivity) {
          const id = sceneId.current;
          if (id) {
            let resume_time = playerRef.current.currentTime()!;
            const playDuration = playDurationRef.current;
            if (playDuration > 0) {
              const videoDuration = playerRef.current.duration();
              const percentCompleted = (100 / videoDuration) * resume_time;
              // if video is 98% or more complete then reset resume_time
              if (percentCompleted >= 98) {
                resume_time = 0;
              }
              sceneSaveActivity({
                variables: {
                  id,
                  resume_time,
                  playDuration,
                },
              });
            }
          }
        }
      }
      playerRef.current = undefined;
      player.dispose();
    };
  }, [sceneSaveActivity, ignoreInterval, trackActivity]);

  useEffect(() => {
    const player = playerRef.current;
    if (!player) return;
    const skipButtons = player.skipButtons();
    skipButtons.setForwardHandler(onNext);
    skipButtons.setBackwardHandler(onPrevious);
  }, [onNext, onPrevious]);

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

  // Player event handlers
  useEffect(() => {
    function canplay(this: VideoJsPlayer) {
      if (initialTimestamp.current !== -1) {
        this.currentTime(initialTimestamp.current);
        initialTimestamp.current = -1;
      }
    }

    function playing(this: VideoJsPlayer) {
      // This still runs even if autoplay failed on Safari,
      // only set flag if actually playing
      if (!started.current && !this.paused()) {
        started.current = true;
      }
    }

    function loadstart(this: VideoJsPlayer) {
      setReady(true);
    }

    function fullscreenchange(this: VideoJsPlayer) {
      setFullscreen(this.isFullscreen());
    }

    const player = playerRef.current;
    if (!player) return;

    player.on("canplay", canplay);
    player.on("playing", playing);
    player.on("loadstart", loadstart);
    player.on("fullscreenchange", fullscreenchange);

    return () => {
      player.off("canplay", canplay);
      player.off("playing", playing);
      player.off("loadstart", loadstart);
      player.off("fullscreenchange", fullscreenchange);
    };
  }, []);
  useEffect(() => {
    function onplay(this: VideoJsPlayer) {
      trackTime.current = true;
      this.persistVolume().enabled = true;
      if (scene?.interactive && interactiveReady.current) {
        interactiveClient.play(this.currentTime());
      }
    }

    function pause(this: VideoJsPlayer) {
      trackTime.current = false;
      interactiveClient.pause();

      const id = sceneId.current;
      if (id) {
        const playDuration = playDurationRef.current;
        let resume_time = this.currentTime()!;
        const videoDuration = this.duration();
        const percentCompleted = (100 / videoDuration) * resume_time;
        // if video is 98% or more complete then reset resume_time
        if (percentCompleted >= 98) {
          resume_time = 0;
        }
        sceneSaveActivity({
          variables: {
            id,
            resume_time,
            playDuration,
          },
        });
      }
    }

    function seeking(this: VideoJsPlayer) {
      if (this.paused()) return;
      if (scene?.interactive && interactiveReady.current) {
        interactiveClient.play(this.currentTime());
      }
    }

    function timeupdate(this: VideoJsPlayer) {
      if (this.paused()) return;
      if (scene?.interactive && interactiveReady.current) {
        interactiveClient.ensurePlaying(this.currentTime());
      }
      setTime(this.currentTime());
    }

    const player = playerRef.current;
    if (!player) return;

    player.on("play", onplay);
    player.on("pause", pause);
    player.on("seeking", seeking);
    player.on("timeupdate", timeupdate);

    return () => {
      player.off("play", onplay);
      player.off("pause", pause);
      player.off("seeking", seeking);
      player.off("timeupdate", timeupdate);
    };
  }, [interactiveClient, scene, sceneSaveActivity]);

  useEffect(() => {
    const player = playerRef.current;
    if (!player) return;

    // don't re-initialise the player unless the scene has changed
    if (!scene || !file || scene.id === sceneId.current) return;

    // if new scene was picked from playlist
    if (playerRef.current && sceneId.current) {
      if (trackActivity) {
        playDurationRef.current = 0;
        recordedActivity.current = false;
      }
    }

    sceneId.current = scene.id;

    setReady(false);

    // always stop the interactive client on initialisation
    interactiveClient.pause();
    interactiveReady.current = false;

    const alwaysStartFromBeginning =
      uiConfig?.alwaysStartFromBeginning ?? false;
    const isLandscape = file.height && file.width && file.width > file.height;

    if (isLandscape) {
      player.landscapeFullscreen({
        fullscreen: {
          enterOnRotate: true,
          exitOnRotate: true,
          alwaysInLandscapeMode: true,
          iOS: false,
        },
      });
    }

    const { duration } = file;
    const sourceSelector = player.sourceSelector();
    sourceSelector.setSources(
      scene.sceneStreams.map((stream) => {
        const src = new URL(stream.url);
        const isDirect =
          src.pathname.endsWith("/stream") ||
          src.pathname.endsWith("/stream.m3u8");

        return {
          src: stream.url,
          type: stream.mime_type ?? undefined,
          label: stream.label ?? undefined,
          offset: !isDirect,
          duration,
        };
      })
    );

    const markers = player.markers();
    markers.clearMarkers();
    for (const marker of scene.scene_markers) {
      markers.addMarker({
        title: marker.title,
        time: marker.seconds,
      });
    }

    function getDefaultLanguageCode() {
      let languageCode = window.navigator.language;

      if (languageCode.indexOf("-") !== -1) {
        languageCode = languageCode.split("-")[0];
      }

      if (languageCode.indexOf("_") !== -1) {
        languageCode = languageCode.split("_")[0];
      }

      return languageCode;
    }

    if (scene.captions && scene.captions.length > 0) {
      const languageCode = getDefaultLanguageCode();
      let hasDefault = false;

      for (let caption of scene.captions) {
        const lang = caption.language_code;
        let label = lang;
        if (languageMap.has(lang)) {
          label = languageMap.get(lang)!;
        }

        label = label + " (" + caption.caption_type + ")";
        const setAsDefault = !hasDefault && languageCode == lang;
        if (setAsDefault) {
          hasDefault = true;
        }
        sourceSelector.addTextTrack(
          {
            src: `${scene.paths.caption}?lang=${lang}&type=${caption.caption_type}`,
            kind: "captions",
            srclang: lang,
            label: label,
            default: setAsDefault,
          },
          false
        );
      }
    }

    if (scene.paths.screenshot) {
      player.poster(scene.paths.screenshot);
    } else {
      player.poster("");
    }

    auto.current =
      autoplay ||
      (interfaceConfig?.autostartVideo ?? false) ||
      _initialTimestamp > 0;

    var startPositition = _initialTimestamp;
    if (
      !(alwaysStartFromBeginning || sessionInitialised) &&
      file.duration > scene.resume_time!
    ) {
      startPositition = scene.resume_time!;
    }

    initialTimestamp.current = startPositition;
    setTime(startPositition);
    setSessionInitialised(true);

    player.load();
    player.focus();

    player.ready(() => {
      player.vttThumbnails().src(scene.paths.vtt ?? null);
    });

    started.current = false;

    return () => {
      // stop the interactive client
      interactiveClient.pause();
    };
  }, [
    file,
    scene,
    interactiveClient,
    sessionInitialised,
    trackActivity,
    autoplay,
    interfaceConfig?.autostartVideo,
    uiConfig?.alwaysStartFromBeginning,
    _initialTimestamp,
  ]);

  useEffect(() => {
    const player = playerRef.current;
    if (!player) return;

    player.loop(looping);
    interactiveClient.setLooping(looping);
  }, [interactiveClient, looping]);

  useEffect(() => {
    if (!scene || !ready || !auto.current) {
      return;
    }

    // check if we're waiting for the interactive client
    if (
      scene.interactive &&
      interactiveClient.handyKey &&
      currentScript !== scene.paths.funscript
    ) {
      return;
    }

    const player = playerRef.current;
    if (!player) return;

    player.play()?.catch(() => {
      // Browser probably blocking non-muted autoplay, so mute and try again
      player.persistVolume().enabled = false;
      player.muted(true);

      player.play();
    });
    auto.current = false;
  }, [scene, ready, interactiveClient, currentScript]);

  useEffect(() => {
    // Attach handler for onComplete event
    const player = playerRef.current;
    if (!player) return;

    player.on("ended", onComplete);

    return () => player.off("ended");
  }, [onComplete]);

  const onScrubberScroll = () => {
    if (started.current) {
      playerRef.current?.pause();
    }
  };
  const onScrubberSeek = (seconds: number) => {
    if (started.current) {
      playerRef.current?.currentTime(seconds);
    } else {
      initialTimestamp.current = seconds;
      setTime(seconds);
    }
  };

  // Override spacebar to always pause/play
  function onKeyDown(this: HTMLDivElement, event: KeyboardEvent) {
    const player = playerRef.current;
    if (!player) return;

    if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey) {
      return;
    }
    if (event.key == " ") {
      event.preventDefault();
      event.stopPropagation();
      if (player.paused()) {
        player.play();
      } else {
        player.pause();
      }
    }
  }

  const isPortrait =
    scene && file && file.height && file.width && file.height > file.width;

  return (
    <div
      className={cx("VideoPlayer", { portrait: isPortrait })}
      onKeyDownCapture={onKeyDown}
    >
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
      {scene && file && showScrubber && (
        <ScenePlayerScrubber
          file={file}
          scene={scene}
          time={time}
          onSeek={onScrubberSeek}
          onScroll={onScrubberScroll}
        />
      )}
    </div>
  );
};

export default ScenePlayer;
