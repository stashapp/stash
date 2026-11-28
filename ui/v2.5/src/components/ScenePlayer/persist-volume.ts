import videojs from "video.js";
import type VideoJsPlayer from "video.js/dist/types/player";
import type Plugin from "video.js/dist/types/plugin";
import localForage from "localforage";

const levelKey = "volume-level";
const mutedKey = "volume-muted";

interface IPersistVolumeOptions {
  enabled?: boolean;
}

class PersistVolumePlugin extends (videojs.getPlugin(
  "plugin"
) as typeof Plugin) {
  enabled: boolean;

  constructor(player: VideoJsPlayer, options?: IPersistVolumeOptions) {
    super(player);

    this.enabled = options?.enabled ?? true;

    player.on("volumechange", () => {
      if (this.enabled) {
        localForage.setItem(levelKey, player.volume());
        localForage.setItem(mutedKey, player.muted());
      }
    });

    player.ready(() => {
      this.ready();
    });
  }

  private ready() {
    localForage.getItem<number>(levelKey).then((value) => {
      if (value !== null) {
        this.player.volume(value);
      }
    });

    localForage.getItem<boolean>(mutedKey).then((value) => {
      if (value !== null) {
        this.player.muted(value);
      }
    });
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("persistVolume", PersistVolumePlugin);

declare module "video.js/dist/types/player" {
  export default interface Player {
    persistVolume: () => PersistVolumePlugin;
  }
}

export default PersistVolumePlugin;
