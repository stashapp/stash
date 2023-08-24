/* eslint-disable @typescript-eslint/naming-convention */

declare module "videojs-vr" {
  import videojs from "video.js";
  // we don't want to depend on THREE.js directly, these are just typedefs for videojs-vr
  // eslint-disable-next-line import/no-extraneous-dependencies
  import * as THREE from "three";

  declare function videojsVR(options?: videojsVR.Options): videojsVR.Plugin;

  declare namespace videojsVR {
    const VERSION: typeof videojs.VERSION;

    type ProjectionType =
      // The video is half sphere and the user should not be able to look behind themselves
      | "180"
      // Used for side-by-side 180 videos The video is half sphere and the user should not be able to look behind themselves
      | "180_LR"
      // Used for monoscopic 180 videos The video is half sphere and the user should not be able to look behind themselves
      | "180_MONO"
      // The video is a sphere
      | "360"
      | "Sphere"
      | "equirectangular"
      // The video is a cube
      | "360_CUBE"
      | "Cube"
      // This video is not a 360 video
      | "NONE"
      // Check player.mediainfo.projection to see if the current video is a 360 video.
      | "AUTO"
      // Used for side-by-side 360 videos
      | "360_LR"
      // Used for top-to-bottom 360 videos
      | "360_TB"
      // Used for Equi-Angular Cubemap videos
      | "EAC"
      // Used for side-by-side Equi-Angular Cubemap videos
      | "EAC_LR";

    interface Options {
      /**
       * Force the cardboard button to display on all devices even if we don't think they support it.
       *
       * @default false
       */
      forceCardboard?: boolean;

      /**
       * Whether motion/gyro controls should be enabled.
       *
       * @default true on iOS and Android
       */
      motionControls?: boolean;

      /**
       * Defines the projection type.
       *
       * @default "AUTO"
       */
      projection?: ProjectionType;

      /**
       * This alters the number of segments in the spherical mesh onto which equirectangular videos are projected.
       * The default is 32 but in some circumstances you may notice artifacts and need to increase this number.
       *
       * @default 32
       */
      sphereDetail?: number;

      /**
       * Enable debug logging for this plugin
       *
       * @default false
       */
      debug?: boolean;

      /**
       * Use this property to pass the Omnitone library object to the plugin. Please be aware of, the Omnitone library is not included in the build files.
       */
      omnitone?: object;

      /**
       * Default options for the Omnitone library. Please check available options on https://github.com/GoogleChrome/omnitone
       */
      omnitoneOptions?: object;

      /**
       * Feature to disable the togglePlay manually. This functionality is useful in live events so that users cannot stop the live, but still have a controlBar available.
       *
       * @default false
       */
      disableTogglePlay?: boolean;
    }

    interface PlayerMediaInfo {
      /**
       * This should be set on a source-by-source basis to turn 360 videos on an off depending upon the video.
       * Note that AUTO is the same as NONE for player.mediainfo.projection.
       */
      projection?: ProjectionType;
    }

    class Plugin extends videojs.Plugin {
      setProjection(projection: ProjectionType): void;
      init(): void;
      reset(): void;

      cameraVector: THREE.Vector3;

      camera: THREE.Camera;
      scene: THREE.Scene;
      renderer: THREE.Renderer;
    }
  }

  export = videojsVR;

  declare module "video.js" {
    interface VideoJsPlayer {
      vr: typeof videojsVR;
      mediainfo?: videojsVR.PlayerMediaInfo;
    }
  }
}
