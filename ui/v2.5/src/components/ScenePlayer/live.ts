import { debounce } from "lodash-es";
import videojs, { VideoJsPlayer } from "video.js";

export interface ISource extends videojs.Tech.SourceObject {
  offset?: boolean;
  duration?: number;
}

interface ICue extends TextTrackCue {
  _startTime?: number;
  _endTime?: number;
}

// delay before loading new source after setting currentTime
const loadDelay = 200;

function offsetMiddleware(player: VideoJsPlayer) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- allow access to private tech methods
  let tech: any;
  let source: ISource;
  let offsetStart: number | undefined;
  let seeking = 0;

  function initCues(cues: TextTrackCueList) {
    const offset = offsetStart ?? 0;
    for (let j = 0; j < cues.length; j++) {
      const cue = cues[j] as ICue;
      cue._startTime = cue.startTime;
      cue.startTime = cue._startTime - offset;
      cue._endTime = cue.endTime;
      cue.endTime = cue._endTime - offset;
    }
  }

  function updateOffsetStart(offset: number | undefined) {
    offsetStart = offset;

    if (!tech) return;
    offset = offset ?? 0;

    const tracks = tech.remoteTextTracks();
    for (let i = 0; i < tracks.length; i++) {
      const { cues } = tracks[i];
      if (cues) {
        for (let j = 0; j < cues.length; j++) {
          const cue = cues[j] as ICue;
          if (cue._startTime === undefined || cue._endTime === undefined) {
            continue;
          }
          cue.startTime = cue._startTime - offset;
          cue.endTime = cue._endTime - offset;
        }
      }
    }
  }

  const loadSource = debounce(
    (seconds: number) => {
      const srcUrl = new URL(source.src);
      srcUrl.searchParams.set("start", seconds.toString());
      source.src = srcUrl.toString();

      const poster = player.poster();
      const playbackRate = tech.playbackRate();
      seeking = tech.paused() ? 1 : 2;
      player.poster("");
      tech.setSource(source);
      tech.setPlaybackRate(playbackRate);
      tech.one("canplay", () => {
        player.poster(poster);
        if (seeking === 1 || tech.scrubbing()) {
          tech.pause();
        }
        seeking = 0;
      });
      tech.trigger("timeupdate");
      tech.trigger("pause");
      tech.trigger("seeking");
      tech.play();
    },
    loadDelay,
    { leading: true }
  );

  return {
    setTech(newTech: videojs.Tech) {
      tech = newTech;

      const _addRemoteTextTrack = tech.addRemoteTextTrack.bind(tech);
      function addRemoteTextTrack(
        this: VideoJsPlayer,
        options: videojs.TextTrackOptions,
        manualCleanup: boolean
      ) {
        const textTrack = _addRemoteTextTrack(options, manualCleanup);
        textTrack.addEventListener("load", () => {
          const { cues } = textTrack.track;
          if (cues) {
            initCues(cues);
          }
        });

        return textTrack;
      }
      tech.addRemoteTextTrack = addRemoteTextTrack;

      const trackEls: HTMLTrackElement[] = tech.remoteTextTrackEls();
      for (let i = 0; i < trackEls.length; i++) {
        const trackEl = trackEls[i];
        const { track } = trackEl;
        if (track.cues) {
          initCues(track.cues);
        } else {
          trackEl.addEventListener("load", () => {
            if (track.cues) {
              initCues(track.cues);
            }
          });
        }
      }
    },
    setSource(
      srcObj: ISource,
      next: (err: unknown, src: videojs.Tech.SourceObject) => void
    ) {
      if (srcObj.offset && srcObj.duration) {
        updateOffsetStart(0);
      } else {
        updateOffsetStart(undefined);
      }
      source = srcObj;
      next(null, srcObj);
    },
    duration(seconds: number) {
      if (source.duration) {
        return source.duration;
      } else {
        return seconds;
      }
    },
    buffered(buffers: TimeRanges) {
      if (offsetStart === undefined) {
        return buffers;
      }

      const timeRanges: number[][] = [];
      for (let i = 0; i < buffers.length; i++) {
        const start = buffers.start(i) + offsetStart;
        const end = buffers.end(i) + offsetStart;

        timeRanges.push([start, end]);
      }

      // types for createTimeRanges are incorrect, should be number[][] not TimeRange[]
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      return videojs.createTimeRanges(timeRanges as any);
    },
    currentTime(seconds: number) {
      return (offsetStart ?? 0) + seconds;
    },
    setCurrentTime(seconds: number) {
      if (offsetStart === undefined) {
        return seconds;
      }

      const offsetSeconds = seconds - offsetStart;
      const buffers = tech.buffered() as TimeRanges;
      for (let i = 0; i < buffers.length; i++) {
        const start = buffers.start(i);
        const end = buffers.end(i);
        // seek point is in buffer, just seek normally
        if (start <= offsetSeconds && offsetSeconds <= end) {
          return offsetSeconds;
        }
      }

      updateOffsetStart(seconds);

      loadSource(seconds);

      return 0;
    },
    callPlay() {
      if (seeking) {
        seeking = 2;
        return videojs.middleware.TERMINATOR;
      }
    },
  };
}

videojs.use("*", offsetMiddleware);
