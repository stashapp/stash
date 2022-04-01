import videojs, { VideoJsPlayer } from "video.js";

const offset = function (this: VideoJsPlayer) {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  const Player = this.constructor as any;

  if (!Player.__super__ || !Player.__super__.__offsetInit) {
    Player.__super__ = {
      __offsetInit: true,
      duration: Player.prototype.duration,
      currentTime: Player.prototype.currentTime,
      remainingTime: Player.prototype.remainingTime,
      getCache: Player.prototype.getCache,
    };

    Player.prototype.clearOffsetDuration = function () {
      this._offsetDuration = undefined;
      this._offsetStart = undefined;
    };

    Player.prototype.setOffsetDuration = function (duration: number) {
      this._offsetDuration = duration;
    };

    Player.prototype.duration = function () {
      if (this._offsetDuration !== undefined) {
        return this._offsetDuration;
      }
      return Player.__super__.duration.apply(this, arguments);
    };

    Player.prototype.currentTime = function (seconds: number) {
      if (seconds !== undefined && this._offsetDuration !== undefined) {
        this._offsetStart = seconds;

        const srcUrl = new URL(this.src());
        srcUrl.searchParams.delete("start");
        srcUrl.searchParams.append("start", seconds.toString());
        const currentSrc = this.currentSource();
        const newSources = this.currentSources().map(
          (source: videojs.Tech.SourceObject) => {
            return {
              ...source,
              src:
                source.src === currentSrc.src ? srcUrl.toString() : source.src,
            };
          }
        );
        this.src(newSources);
        this.play();

        return seconds;
      }
      return (
        (this._offsetStart ?? 0) +
        Player.__super__.currentTime.apply(this, arguments)
      );
    };

    Player.prototype.getCache = function () {
      const cache = Player.__super__.getCache.apply(this);
      if (this._offsetDuration !== undefined)
        return {
          ...cache,
          currentTime:
            (this._offsetStart ?? 0) + Player.__super__.currentTime.apply(this),
        };
      return cache;
    };

    Player.prototype.remainingTime = function () {
      if (this._offsetDuration !== undefined) {
        return this._offsetDuration - this.currentTime();
      }
      return this.duration() - this.currentTime();
    };
  }
};

// Register the plugin with video.js.
videojs.registerPlugin("offset", offset);

export default offset;
