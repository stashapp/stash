import videojs, { VideoJsPlayer } from "video.js";

export interface ISource extends videojs.Tech.SourceObject {
  offset?: boolean;
  duration?: number;
}

interface ICue extends TextTrackCue {
  _startTime?: number;
  _endTime?: number;
}

class OffsetPlugin extends videojs.getPlugin("plugin") {
  private offsetDuration?: number;
  private offsetStart = 0;
  private currentSrc: string | null = null;

  constructor(player: VideoJsPlayer) {
    super(player);

    const plugin = this;

    const _src = player.src.bind(player);
    function src(this: VideoJsPlayer): string;
    function src(
      this: VideoJsPlayer,
      source: string | videojs.Tech.SourceObject | videojs.Tech.SourceObject[]
    ): void;
    function src(
      this: VideoJsPlayer,
      source?: string | videojs.Tech.SourceObject | videojs.Tech.SourceObject[]
    ) {
      if (source === undefined) {
        return _src();
      }

      plugin.currentSrc = null;
      plugin.offsetDuration = undefined;
      plugin.updateOffsetStart(0);

      return _src(source);
    }
    player.src = src;

    const _duration = player.duration.bind(player);
    function duration(this: VideoJsPlayer): number;
    function duration(this: VideoJsPlayer, seconds: number): void;
    function duration(this: VideoJsPlayer, seconds?: number) {
      if (seconds !== undefined) {
        return _duration(seconds);
      }
      if (plugin.offsetDuration !== undefined) {
        return plugin.offsetDuration;
      }
      return _duration();
    }
    player.duration = duration;

    const _currentTime = player.currentTime.bind(player);
    function currentTime(this: VideoJsPlayer): number;
    function currentTime(this: VideoJsPlayer, seconds: number): void;
    function currentTime(this: VideoJsPlayer, seconds?: number) {
      if (seconds === undefined) {
        return plugin.offsetStart + _currentTime();
      }
      if (plugin.offsetDuration === undefined) {
        return _currentTime(seconds);
      }

      if (seconds === _currentTime()) return;

      plugin.updateOffsetStart(seconds);

      const currentSource = this.currentSource();
      const newSources = this.currentSources();
      for (const source of newSources) {
        if (source.src === currentSource.src) {
          const srcUrl = new URL(currentSource.src);
          srcUrl.searchParams.set("start", seconds.toString());
          source.src = srcUrl.toString();
          plugin.currentSrc = source.src;
        }
      }
      const playbackRate = this.playbackRate();
      _src(newSources);
      this.play()?.then(() => {
        this.playbackRate(playbackRate);
      });
    }
    player.currentTime = currentTime;

    const _addRemoteTextTrack = player.addRemoteTextTrack.bind(player);
    function addRemoteTextTrack(
      this: VideoJsPlayer,
      options: videojs.TextTrackOptions,
      manualCleanup: boolean
    ) {
      const textTrack = _addRemoteTextTrack(options, manualCleanup);
      textTrack.addEventListener("load", () => {
        const { cues } = textTrack.track;
        if (cues) {
          for (let j = 0; j < cues.length; j++) {
            const cue = cues[j] as ICue;
            cue._startTime = cue.startTime;
            cue.startTime = cue._startTime - plugin.offsetStart;
            cue._endTime = cue.endTime;
            cue.endTime = cue._endTime - plugin.offsetStart;
          }
        }
      });

      return textTrack;
    }
    player.addRemoteTextTrack = addRemoteTextTrack;

    player.on("loadstart", () => {
      const source = player.currentSource() as ISource;
      if (source.src === this.currentSrc) return;

      this.currentSrc = source.src;
      if (source.offset && source.duration) {
        plugin.offsetDuration = source.duration;
      } else {
        plugin.offsetDuration = undefined;
      }
      this.updateOffsetStart(0);
    });
  }

  private updateOffsetStart(offset: number) {
    this.offsetStart = offset;

    const tracks = this.player.remoteTextTracks();
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
}

// Register the plugin with video.js.
videojs.registerPlugin("offset", OffsetPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    offset: () => OffsetPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    offset?: {};
  }
}

export default OffsetPlugin;
