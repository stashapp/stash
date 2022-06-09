/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import "videojs-vr";
import * as GQL from "../../core/generated-graphql";

const vrmode = function (this: VideoJsPlayer, scene: GQL.SceneDataFragment) {
  const player = this;

  switch (scene.projection) {
    case "SPHERE":
      switch (scene.stereo_mode) {
        case "TOP_BOTTOM":
          player.vr({ projection: "360_TB" });
          break;
        case "LEFT_RIGHT":
          player.vr({ projection: "360_LR" });
          break;
        default:
          player.vr({ projection: "360" });
      }
      break;
    case "DOME":
    case "MKX200":
    case "RF52":
      switch (scene.stereo_mode) {
        case "LEFT_RIGHT":
          player.vr({ projection: "180_LR" });
          break;
        default:
          player.vr({ projection: "180_MONO" });
      }
      break;
  }
};

// Register the plugin with video.js.
videojs.registerPlugin("vrmode", vrmode);

declare module "video.js" {
  interface VideoJsPlayer {
    vrmode: (scene: GQL.SceneDataFragment) => void;
    vr: (options: Object) => void;
  }
}

export default vrmode;
