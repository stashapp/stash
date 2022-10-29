import videojs, { VideoJsPlayer } from "video.js";

class OffsetPlugin extends videojs.getPlugin("plugin") {
  private offsetDuration?: number;
  private offsetStart?: number;

  constructor(player: VideoJsPlayer) {
    super(player);

    const plugin = this;

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
        return (plugin.offsetStart ?? 0) + _currentTime();
      }
      if (plugin.offsetDuration === undefined) {
        return _currentTime(seconds);
      }

      plugin.offsetStart = seconds;

      const srcUrl = new URL(this.src());
      srcUrl.searchParams.delete("start");
      srcUrl.searchParams.append("start", seconds.toString());
      const currentSrc = this.currentSource();
      const newSources = this.currentSources().map(
        (source: videojs.Tech.SourceObject) => {
          return {
            ...source,
            src: source.src === currentSrc.src ? srcUrl.toString() : source.src,
          };
        }
      );
      this.src(newSources);
      this.play();
    }
    player.currentTime = currentTime;

    const _getCache = player.getCache.bind(player);
    function getCache(this: VideoJsPlayer) {
      const cache = _getCache();
      if (plugin.offsetDuration !== undefined)
        return {
          ...cache,
          currentTime: player.currentTime(),
        };
      return cache;
    }
    player.getCache = getCache;
  }

  setOffsetDuration(duration: number) {
    this.offsetDuration = duration;
  }

  clearOffsetDuration() {
    this.offsetDuration = undefined;
    this.offsetStart = undefined;
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
