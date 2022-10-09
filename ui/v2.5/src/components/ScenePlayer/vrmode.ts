/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import "videojs-vr";
import * as GQL from "../../core/generated-graphql";
import { projectionStrings } from "../../utils/projection";
import { stereoModeStrings } from "../../utils/stereoMode";

function isVrDevice() {
  return navigator.userAgent.match(/oculusbrowser|\svr\s/i);
}

const vrmode = function (this: VideoJsPlayer, scene: GQL.SceneDataFragment) {
  const player = this;
  let projection;
  let stereo_mode;
  let lowercaseTags = scene.tags.map((tag) => tag.name.toLowerCase());
  let file = scene.files[0];

  if (!isVrDevice()) {
    if (
      file.projection &&
      file.projection !== "AUTO" &&
      projectionStrings.includes(file.projection)
    ) {
      ({ projection } = file);
    } else {
      if (
        lowercaseTags.filter((tag) =>
          /^vr$|virtual reality|^vr porn|^3d$/.exec(tag)
        ).length > 0
      ) {
        if (
          lowercaseTags.filter((tag) => /^dome|^180$|^180\D$/.exec(tag))
            .length > 0
        ) {
          projection = "DOME";
        } else if (
          lowercaseTags.filter((tag) =>
            /^sphere$|^equirectangular$|^360$|^360[\D\-_]/.exec(tag)
          ).length > 0
        ) {
          projection = "SPHERE";
        } else if (
          lowercaseTags.filter((tag) => /^mkx200$/.exec(tag)).length > 0
        ) {
          projection = "MKX200";
        } else if (
          lowercaseTags.filter((tag) => /^cube$/.exec(tag)).length > 0
        ) {
          projection = "CUBE";
        } else if (
          lowercaseTags.filter((tag) =>
            /^eac$|^eac[\-_]|^equi[\-_\s]?angular[\-_\s]?cube[\-_\s]?/.exec(tag)
          ).length > 0
        ) {
          projection = "EAC";
        } else if (
          lowercaseTags.filter((tag) => /^rf52$/.exec(tag)).length > 0
        ) {
          projection = "RF52";
        } else if (
          lowercaseTags.filter((tag) => /^fisheye$/.exec(tag)).length > 0
        ) {
          projection = "FISHEYE";
        } else {
          projection = "DOME";
        }
      }
    }

    if (
      file.stereo_mode &&
      file.stereo_mode !== "AUTO" &&
      stereoModeStrings.includes(file.stereo_mode)
    ) {
      ({ stereo_mode } = file);
    } else {
      if (
        lowercaseTags.filter((tag) =>
          /^vr$|virtual reality|^vr porn|^3d$/.exec(tag)
        ).length > 0
      ) {
        if (
          lowercaseTags.filter((tag) =>
            /^mono$|^monoscopic$|^[0-9]{1,3}[\D\-_]mono$/.exec(tag)
          ).length > 0
        ) {
          stereo_mode = "MONO";
        } else if (
          lowercaseTags.filter((tag) =>
            /^stereo$|^stereoscopic|^[0-9]{1,3}[\D\-_]lr$|^[0-9]{1,3}[\D\-_]left[\D\-_]right$/.exec(
              tag
            )
          ).length > 0
        ) {
          stereo_mode = "LEFT_RIGHT";
        } else if (
          lowercaseTags.filter((tag) =>
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
