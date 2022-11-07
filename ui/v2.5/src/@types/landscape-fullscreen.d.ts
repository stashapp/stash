/* eslint-disable @typescript-eslint/naming-convention */

declare module "videojs-landscape-fullscreen" {
  import videojs from "video.js";

  function landscapeFullscreen(options?: {
    fullscreen: landscapeFullscreen.Options;
  }): void;

  namespace landscapeFullscreen {
    const VERSION: typeof videojs.VERSION;

    interface Options {
      /**
       * Enter fullscreen mode on rotating the device to landscape.
       * @default true
       */
      enterOnRotate?: boolean;
      /**
       * Exit fullscreen mode on rotating the device to portrait.
       * @default true
       */
      exitOnRotate?: boolean;
      /**
       * Always enter fullscreen in landscape mode even when device is in portrait mode (works on Chromium, Firefox, and IE >= 11).
       * @default true
       */
      alwaysInLandscapeMode?: boolean;
      /**
       * Whether to use fake fullscreen on iOS (needed for displaying player controls instead of system controls).
       * @default true
       */
      iOS?: boolean;
    }
  }

  export = landscapeFullscreen;

  declare module "video.js" {
    interface VideoJsPlayer {
      landscapeFullscreen: typeof landscapeFullscreen;
    }
    interface VideoJsPlayerPluginOptions {
      landscapeFullscreen?: { fullscreen: landscapeFullscreen.Options };
    }
  }
}
