import videojs from "video.js";
import type VideoJsPlayer from "video.js/dist/types/player";
import type Plugin from "video.js/dist/types/plugin";
import type BigPlayButtonType from "video.js/dist/types/big-play-button";

// prettier-ignore
const BigPlayButton = videojs.getComponent(
  "BigPlayButton"
) as typeof BigPlayButtonType;

class BigPlayPauseButton extends BigPlayButton {
  handleClick(event: KeyboardEvent | MouseEvent | TouchEvent) {
    if (this.player().paused()) {
      super.handleClick(event);
    } else {
      this.player().pause();
    }
  }

  buildCSSClass() {
    return "vjs-control vjs-button vjs-big-play-pause-button";
  }
}

class BigButtonGroup extends videojs.getComponent("Component") {
  // options carries playerOptions, which the skip buttons read their seconds from
  constructor(player: VideoJsPlayer, options?: object) {
    super(player, options);

    this.addChild("SkipBackward");
    this.addChild("BigPlayPauseButton");
    this.addChild("SkipForward");
  }

  createEl() {
    return super.createEl("div", {
      className: "vjs-big-button-group",
    });
  }
}

class BigButtonsPlugin extends (videojs.getPlugin("plugin") as typeof Plugin) {
  constructor(player: VideoJsPlayer) {
    super(player);

    player.ready(() => {
      player.addChild("BigButtonGroup");
    });
  }
}

// Register the plugin with video.js.
videojs.registerComponent("BigButtonGroup", BigButtonGroup);
videojs.registerComponent("BigPlayPauseButton", BigPlayPauseButton);
videojs.registerPlugin("bigButtons", BigButtonsPlugin);

declare module "video.js/dist/types/player" {
  export default interface Player {
    bigButtons: () => BigButtonsPlugin;
  }
}

export default BigButtonsPlugin;
