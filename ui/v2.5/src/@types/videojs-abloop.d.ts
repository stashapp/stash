/* XXbiome-ignore @typescript-eslint/naming-convention */

declare module "videojs-abloop" {
  import videojs from "video.js";
  import type VjsPlugin from "video.js/dist/types/plugin";

  declare function abLoopPlugin(
    window: Window & typeof globalThis,
    player: videojs
  ): abLoopPlugin.Plugin;

  declare namespace abLoopPlugin {
    interface Options {
      start: number | boolean;
      end: number | boolean;
      enabled: boolean;
      loopIfBeforeStart: boolean;
      loopIfAfterEnd: boolean;
      pauseBeforeLooping: boolean;
      pauseAfterLooping: boolean;
    }

    class Plugin extends VjsPlugin {
      getOptions(): Options;
      setOptions(o: Options): void;
    }
  }

  export = abLoopPlugin;

  declare module "video.js/dist/types/player" {
    export default interface Player {
      abLoopPlugin: abLoopPlugin.Plugin;
    }
  }
}
