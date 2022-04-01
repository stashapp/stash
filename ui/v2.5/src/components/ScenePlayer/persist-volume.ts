import videojs, { VideoJsPlayer } from "video.js";
import localForage from "localforage";

const persistVolume = function (this: VideoJsPlayer) {
  const player = this;
  const levelKey = "volume-level";
  const mutedKey = "volume-muted";

  player.on("volumechange", function () {
    localForage.setItem(levelKey, player.volume());
    localForage.setItem(mutedKey, player.muted());
  });

  localForage.getItem(levelKey).then((value) => {
    if (value !== null) {
      player.volume(value as number);
    }
  });

  localForage.getItem(mutedKey).then((value) => {
    if (value !== null) {
      player.muted(value as boolean);
    }
  });
};

// Register the plugin with video.js.
videojs.registerPlugin("persistVolume", persistVolume);

export default persistVolume;
