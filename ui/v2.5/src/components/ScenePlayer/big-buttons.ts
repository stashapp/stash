import videojs, { VideoJsPlayer } from "video.js";

// prettier-ignore
const BigPlayButton = videojs.getComponent("BigPlayButton") as unknown as typeof videojs.BigPlayButton;

class BigPlayPauseButton extends BigPlayButton {
  handleClick(event: videojs.EventTarget.Event) {
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
  constructor(player: VideoJsPlayer) {
    super(player);

    this.addChild("seekButton", {
      direction: "back",
      seconds: 10,
    });

    this.addChild("BigPlayPauseButton");

    this.addChild("seekButton", {
      direction: "forward",
      seconds: 10,
    });
  }

  createEl() {
    return super.createEl("div", {
      className: "vjs-big-button-group",
    });
  }
}

function bigButtons(this: VideoJsPlayer) {
  this.addChild("BigButtonGroup");
}

// Register the plugin with video.js.
videojs.registerComponent("BigButtonGroup", BigButtonGroup);
videojs.registerComponent("BigPlayPauseButton", BigPlayPauseButton);
videojs.registerPlugin("bigButtons", bigButtons);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    bigButtons: typeof bigButtons;
  }
  interface VideoJsPlayerPluginOptions {
    bigButtons?: {};
  }
}

export default bigButtons;
