/* eslint-disable @typescript-eslint/naming-convention */

declare module "@blaineam/videojs-vr" {
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
      | "EAC_LR"
      // flat screen side-by-side
      | "SBS_MONO";

    interface mediaItem {
      title: string;
      thumbnail: string;
      url: string;
      duration?: number;
    }
    type mediaItems = mediaItem[];

    type orientationOffset = {
      x: number;
      y: number;
      z: number;
    };

    // options are taken verbaitum from the README
    interface Options {
      // Projection mode
      projection?: ProjectionType = "AUTO"; // see ProjectionType
      sphereDetails?: number = 32; // Sphere mesh detail (higher = smoother)

      // VR HUD options
      enableVRHud?: boolean = true; // Enable in-VR controls
      enableVRGallery?: boolean = true; // Enable in-VR video gallery
      showHUDOnStart?: boolean = true; // Show HUD when entering VR
      hudAutoHideDelay?: number = 5000; // Auto-hide HUD after ms (0 to disable)
      hudDistance?: number = 1.5; // Distance of HUD from viewer
      hudHeight?: number = 1.5; // Height of HUD
      hudScale?: number = 0.015; // Scale of HUD elements

      // Behavior options
      forceCardboard?: boolean = false; // Force cardboard button on all devices
      motionControls?: boolean = true; // Enable gyroscope/device orientation
      disableTogglePlay?: boolean = false; // Disable click-to-play

      // Spatial audio (requires Omnitone library)
      omnitone?: Object = null; // Pass Omnitone library object
      omnitoneOptions?: Record<string, unknown> = {}; // Omnitone configuration

      // Media gallery items
      mediaItems?: mediaItems = []; // Array of media items for gallery
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

      // VR HUD
      showHUD(): void; // Show the VR HUD
      hideHUD(): void; // Hide the VR HUD
      toggleHUD(): void; // Toggle HUD visibility

      // VR Gallery
      showGallery(): void; // Show the gallery panel
      hideGallery(): void; // Hide the gallery panel
      toggleGallery(): void; // Toggle gallery visibility
      setGalleryItems(mediaItems): void; // Update gallery media items

      // Favorite state
      setFavoriteState(boolean): void; // Set favorite button state
      getFavoriteState(): boolean; // Get current favorite state

      // Orientation
      setOrientationOffset(orientationOffset): void; // Tilt view
      resetOrientationOffset(): void; // Reset to default orientation
      recenter(): void; // Recenter VR view

      // Status
      isPresenting(): boolean; // Check if currently in VR mode

      camera: THREE.Camera;
      scene: THREE.Scene;
      renderer: THREE.Renderer;
      cameraVector: THREE.Vector3;
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
