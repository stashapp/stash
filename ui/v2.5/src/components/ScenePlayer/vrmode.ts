/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import "videojs-vr";
import * as GQL from "../../core/generated-graphql";
import { projectionStrings } from "../../utils/projection";
import { stereoModeStrings } from "../../utils/stereoMode";

const vrmode = function (this: VideoJsPlayer, scene: GQL.SceneDataFragment) {
  const player = this;
  let projection;
  let stereo_mode;
  let lowercaseSceneTags = scene.tags.map((tag) => tag.name.toLowerCase());

  if (
    scene.projection &&
    scene.projection !== "AUTO" &&
    projectionStrings.includes(scene.projection)
  ) {
    ({ projection } = scene);
  } else {
    if (
      lowercaseSceneTags.filter((tag) =>
        /^vr$|virtual reality|^vr porn|^3d$/.exec(tag)
      ).length > 0
    ) {
      if (
        lowercaseSceneTags.filter((tag) => /^dome|^180$|^180\D$/.exec(tag))
          .length > 0
      ) {
        projection = "DOME";
      } else if (
        lowercaseSceneTags.filter((tag) =>
          /^sphere$|^equirectangular$|^360$|^360[\D\-_]/.exec(tag)
        ).length > 0
      ) {
        projection = "SPHERE";
      } else if (
        lowercaseSceneTags.filter((tag) => /^mkx200$/.exec(tag)).length > 0
      ) {
        projection = "MKX200";
      } else if (
        lowercaseSceneTags.filter((tag) => /^cube$/.exec(tag)).length > 0
      ) {
        projection = "CUBE";
      } else if (
        lowercaseSceneTags.filter((tag) =>
          /^eac$|^eac[\-_]|^equi[\-_\s]?angular[\-_\s]?cube[\-_\s]?/.exec(tag)
        ).length > 0
      ) {
        projection = "EAC";
      } else if (
        lowercaseSceneTags.filter((tag) => /^rf52$/.exec(tag)).length > 0
      ) {
        projection = "RF52";
      } else if (
        lowercaseSceneTags.filter((tag) => /^fisheye$/.exec(tag)).length > 0
      ) {
        projection = "FISHEYE";
      }else {
        projection = "DOME";
      }
    }
  }

  if (
    scene.stereo_mode &&
    scene.stereo_mode !== "AUTO" &&
    stereoModeStrings.includes(scene.stereo_mode)
  ) {
    ({ stereo_mode } = scene);
  } else {
    if (
      lowercaseSceneTags.filter((tag) =>
        /^vr$|virtual reality|^vr porn|^3d$/.exec(tag)
      ).length > 0
    ) {
      if (
        lowercaseSceneTags.filter((tag) =>
          /^mono$|^monoscopic$|^[0-9]{1,3}[\D\-_]mono$/.exec(tag)
        ).length > 0
      ) {
        stereo_mode = "MONO";
      } else if (
        lowercaseSceneTags.filter((tag) =>
          /^stereo$|^stereoscopic|^[0-9]{1,3}[\D\-_]lr$|^[0-9]{1,3}[\D\-_]left[\D\-_]right$/.exec(
            tag
          )
        ).length > 0
      ) {
        stereo_mode = "LEFT_RIGHT";
      } else if (
        lowercaseSceneTags.filter((tag) =>
          /^top[\D\-_]bottom$|^tb$|^[0-9]{1,3}[\D\-_]tb$|^[0-9]{1,3}[\D\-_]top[\D\-_]bottom$/.exec(
            tag
          )
        ).length > 0
      ) {
        stereo_mode = "TOP_BOTTOM";
      } else {
        stereo_mode = "LEFT_RIGHT";
      }
    }
  }

  switch (projection) {
    case "SPHERE":
      switch (stereo_mode) {
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
    case "FISHEYE":
      switch (stereo_mode) {
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
