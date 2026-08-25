import React, {
  KeyboardEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import videojs, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-mobile-ui";
import "videojs-seek-buttons";
import cx from "classnames";

// shared video.js plugins - registered globally on import
import "../ScenePlayer/persist-volume";
import "../ScenePlayer/big-buttons";
import "../ScenePlayer/PlaylistButtons";
import "../ScenePlayer/track-activity";
import "../ScenePlayer/media-session";
import "../ScenePlayer/wake-sentinel";

import {
  useAudioSaveActivity,
  useAudioIncrementPlayCount,
} from "src/core/StashService";

import * as GQL from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { PatchComponent } from "src/patch";

export const AUDIO_PLAYER_ID = "AudioJsPlayer";

// Audio is served as a direct stream only - there is no transcoding, no
// HLS/DASH variants and no generated previews, so the player is a plain
// video.js instance over the single direct-stream source.
function handleHotkeys(player: VideoJsPlayer, event: videojs.KeyboardEvent) {
  function seekStep(step: number) {
    const time = player.currentTime() + step;
    const duration = player.duration();
    if (time < 0) {
      player.currentTime(0);
    } else if (time < duration) {
      player.currentTime(time);
    } else {
      player.currentTime(duration);
    }
  }

  function seekPercent(percent: number) {
    player.currentTime(player.duration() * percent);
  }

  let seekFactor = 10;
  if (event.shiftKey) {
    seekFactor = 5;
  } else if (event.ctrlKey || event.altKey) {
    seekFactor = 60;
  }
  switch (event.which) {
    case 39: // right arrow
      seekStep(seekFactor);
      break;
    case 37: // left arrow
      seekStep(-seekFactor);
      break;
  }

  if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey) {
    return;
  }

  const skipButtons = player.skipButtons();
  if (skipButtons) {
    // handle multimedia keys
    switch (event.key) {
      case "MediaTrackNext":
        if (!skipButtons.onNext) return;
        skipButtons.onNext();
        break;
      case "MediaTrackPrevious":
        if (!skipButtons.onPrevious) return;
        skipButtons.onPrevious();
        break;
      // MediaPlayPause handled by videojs
    }
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

interface IAudioPlayerProps {
  audio: GQL.AudioDataFragment;
  autoplay?: boolean;
  permitLoop?: boolean;
  initialTimestamp?: number;
  onComplete?: () => void;
  onNext?: () => void;
  onPrevious?: () => void;
}

export const AudioPlayer: React.FC<IAudioPlayerProps> = PatchComponent(
  "AudioPlayer",
  ({
    audio,
    autoplay,
    permitLoop = true,
    initialTimestamp = 0,
    onComplete,
    onNext,
    onPrevious,
  }) => {
    const { configuration } = useConfigurationContext();
    const interfaceConfig = configuration?.interface;
    const uiConfig = configuration?.ui;
    const playerRef = useRef<HTMLDivElement>(null);
    const [_player, setPlayer] = useState<VideoJsPlayer>();
    const loadedAudioId = useRef<string>();
    const [audioSaveActivity] = useAudioSaveActivity();
    const [audioIncrementPlayCount] = useAudioIncrementPlayCount();

    const [ready, setReady] = useState(false);
    const started = useRef(false);
    const auto = useRef(false);

    const minimumPlayPercent = uiConfig?.minimumPlayPercent ?? 0;
    const trackActivity = uiConfig?.trackActivity ?? true;

    const file = useMemo(
      () => (audio.files.length > 0 ? audio.files[0] : undefined),
      [audio]
    );

    // the direct stream is the only source
    const source = useMemo(() => audio.audioStreams[0], [audio.audioStreams]);

    const maxLoopDuration = interfaceConfig?.maximumLoopDuration ?? 0;
    const looping = useMemo(
      () =>
        !!file?.duration &&
        permitLoop &&
        maxLoopDuration !== 0 &&
        file.duration < maxLoopDuration,
      [file, permitLoop, maxLoopDuration]
    );

    const getPlayer = useCallback(() => {
      if (!_player) return null;
      if (_player.isDisposed()) return null;
      return _player;
    }, [_player]);

    // Initialize VideoJS player
    useEffect(() => {
      const options: VideoJsPlayerOptions = {
        id: AUDIO_PLAYER_ID,
        controls: true,
        controlBar: {
          pictureInPictureToggle: false,
          fullscreenToggle: false,
          volumePanel: {
            inline: false,
          },
          chaptersButton: false,
        },
        nativeControlsForTouch: false,
        playbackRates: [0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
        inactivityTimeout: 0,
        preload: "none",
        playsinline: true,
        userActions: {
          hotkeys: function (this: VideoJsPlayer, event) {
            handleHotkeys(this, event);
          },
        },
        plugins: {
          persistVolume: {},
          bigButtons: {},
          seekButtons: {
            forward: 10,
            back: 10,
          },
          trackActivity: {},
          mediaSession: {},
          wakeSentinel: {},
          skipButtons: {},
        },
      };

      const audioEl = document.createElement("video-js");
      audioEl.setAttribute("data-vjs-player", "true");
      audioEl.setAttribute("crossorigin", "anonymous");
      audioEl.classList.add("vjs-big-play-centered");
      playerRef.current!.appendChild(audioEl);

      const vjs = videojs(audioEl, options);

      vjs.focus();
      setPlayer(vjs);

      return () => {
        vjs.dispose();
        audioEl.remove();
        setPlayer(undefined);

        // reset so that the source is reloaded if the player is recreated
        loadedAudioId.current = undefined;
      };
      // empty deps - only init once
    }, []);

    // Player event handlers
    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      function playing(this: VideoJsPlayer) {
        if (!started.current && !this.paused()) {
          started.current = true;
        }
      }

      function loadstart() {
        setReady(true);
      }

      player.on("playing", playing);
      player.on("loadstart", loadstart);

      return () => {
        player.off("playing", playing);
        player.off("loadstart", loadstart);
      };
    }, [getPlayer]);

    // load the source
    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      if (loadedAudioId.current === audio.id) return;
      loadedAudioId.current = audio.id;

      if (!source) {
        player.reset();
        return;
      }

      player.src({
        src: source.url,
        type: source.mime_type ?? undefined,
      });

      const resumeTime = audio.resume_time ?? 0;
      const alwaysStartFromBeginning =
        uiConfig?.alwaysStartFromBeginning ?? false;

      let startPosition = initialTimestamp;
      if (
        !startPosition &&
        !alwaysStartFromBeginning &&
        (file?.duration ?? 0) > resumeTime
      ) {
        startPosition = resumeTime;
      }

      player.load();
      player.focus();

      auto.current =
        !!autoplay ||
        (interfaceConfig?.autostartVideo ?? false) ||
        initialTimestamp > 0;

      player.ready(() => {
        if (startPosition) {
          player.currentTime(startPosition);
        }
      });

      started.current = false;
    }, [
      getPlayer,
      audio.id,
      audio.resume_time,
      source,
      file,
      autoplay,
      initialTimestamp,
      interfaceConfig?.autostartVideo,
      uiConfig?.alwaysStartFromBeginning,
    ]);

    useEffect(() => {
      const player = getPlayer();
      if (!player) return;
      const skipButtons = player.skipButtons();
      skipButtons.setForwardHandler(onNext);
      skipButtons.setBackwardHandler(onPrevious);
    }, [getPlayer, onNext, onPrevious]);

    // cover art is used as the poster
    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      player.poster(audio.paths.screenshot ?? "");
    }, [getPlayer, audio.paths.screenshot]);

    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      async function saveActivity(resumeTime: number, playDuration: number) {
        if (!audio.id) return;

        await audioSaveActivity({
          variables: {
            id: audio.id,
            playDuration,
            resume_time: resumeTime,
          },
        });
      }

      async function incrementPlayCount() {
        if (!audio.id) return;

        await audioIncrementPlayCount({
          variables: {
            id: audio.id,
          },
        });
      }

      const activity = player.trackActivity();
      activity.saveActivity = saveActivity;
      activity.incrementPlayCount = incrementPlayCount;
      activity.minimumPlayPercent = minimumPlayPercent;
      activity.setEnabled(trackActivity);
    }, [
      getPlayer,
      audio.id,
      trackActivity,
      minimumPlayPercent,
      audioIncrementPlayCount,
      audioSaveActivity,
    ]);

    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      player.loop(looping);
    }, [getPlayer, looping]);

    useEffect(() => {
      const player = getPlayer();
      if (!player || !ready || !auto.current) {
        return;
      }

      player.play();
      auto.current = false;
    }, [getPlayer, ready]);

    useEffect(() => {
      const player = getPlayer();
      if (!player || !onComplete) return;

      player.on("ended", onComplete);

      return () => player.off("ended");
    }, [getPlayer, onComplete]);

    useEffect(() => {
      const player = getPlayer();
      if (!player) return;

      const performers = audio.performers.map((p) => p.name).join(", ");
      player
        .mediaSession()
        .setMetadata(
          audio.title ?? "Stash",
          audio.studio?.name ?? performers ?? "Stash",
          audio.paths.screenshot || ""
        );
    }, [getPlayer, audio]);

    // Override spacebar to always pause/play
    function onKeyDown(this: HTMLDivElement, event: KeyboardEvent) {
      const player = getPlayer();
      if (!player) return;

      if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey) {
        return;
      }
      if (event.key === " ") {
        event.preventDefault();
        event.stopPropagation();
        if (player.paused()) {
          player.play();
        } else {
          player.pause();
        }
      }
    }

    return (
      <div
        className={cx("AudioPlayer", {
          "no-file": !file,
        })}
        onKeyDownCapture={onKeyDown}
      >
        <div className="audio-wrapper" ref={playerRef} />
      </div>
    );
  }
);

export default AudioPlayer;
