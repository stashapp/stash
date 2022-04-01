import videojs, { VideoJsPlayer } from "video.js";

const BigPlayButton = videojs.getComponent("BigPlayButton");

class BigPlayPauseButton extends BigPlayButton {
  handleClick(event: videojs.EventTarget.Event) {
    if (this.player().paused()) {
      // @ts-ignore for some reason handleClick isn't defined in BigPlayButton type. Not sure why
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
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

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

const bigButtons = function (this: VideoJsPlayer) {
  this.addChild("BigButtonGroup");
};

// Register the plugin with video.js.
videojs.registerComponent("BigButtonGroup", BigButtonGroup);
videojs.registerComponent("BigPlayPauseButton", BigPlayPauseButton);
videojs.registerPlugin("bigButtons", bigButtons);

export default bigButtons;
