/* eslint-disable @typescript-eslint/naming-convention */

declare module "videojs-vtt-thumbnails-freetube" {
  import videojs from "video.js";

  function vttThumbnails(options?: vttThumbnails.Options): void;

  class vttThumbnailsPlugin {
    src(source: string | undefined): void;
  }

  namespace vttThumbnails {
    const VERSION: typeof videojs.VERSION;

    interface Options {
      /**
       * Source URL to use for thumbnails.
       */
      src?: string;
      /**
       * Disables the timestamp that is shown on hover.
       * @default false
       */
      showTimestamp?: boolean;
    }
  }

  export = vttThumbnails;

  declare module "video.js" {
    interface VideoJsPlayer {
      vttThumbnails: vttThumbnailsPlugin;
    }
    interface VideoJsPlayerPluginOptions {
      vttThumbnails?: vttThumbnails.Options;
    }
  }
}
