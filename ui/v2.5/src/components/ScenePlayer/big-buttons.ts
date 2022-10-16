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

interface BigButtonGroupEvent {
  direction: "forward" | "backward";
}

class BigButtonGroup extends videojs.getComponent("Component") {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    const backButton = this.addChild("seekButton", { direction: "back" });

    this.addChild("BigPlayPauseButton", {});

    const forwardButton = this.addChild("seekButton", { direction: "forward" });

    this.on("display", (e: BigButtonGroupEvent) => {
      let el = e.direction === "forward" ? forwardButton : backButton;
      el.addClass("visible");
      setTimeout(() => {
        el.removeClass("visible");
      }, 100);
    });
  }

  createEl() {
    return super.createEl("div", {
      className: "vjs-big-button-group",
    });
  }
}

class touchControls extends videojs.getPlugin("plugin") {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  timer: NodeJS.Timeout | null = null;
  bigButtonGroup: BigButtonGroup | undefined;

  constructor(player: VideoJsPlayer, options: any) {
    super(player, options);

    this.bigButtonGroup = player.getChild("BigButtonGroup");
    if (!videojs.browser.TOUCH_ENABLED || !this.bigButtonGroup) {
      return;
    }
    player.on("touchstart", this.onTouchStart);
  }

  onTouchStart = (event: TouchEvent) => {
    this.handleTouchStart(event);
  };

  handleTouchStart(event: TouchEvent) {
    if (!this.player.hasStarted()) {
      return;
    }
    const playerWidth = this.player.el().getBoundingClientRect().width;
    const touchLocation = event.changedTouches[0]?.clientX;
    const touchPercent = touchLocation / playerWidth;

    if (this.timer == null) {
      this.timer = setTimeout(() => {
        this.timer = null;
      }, 250);
      return;
    }
    clearTimeout(this.timer);
    this.timer = null;

    if (touchPercent > 0.65) {
      this.bigButtonGroup?.trigger({type: "display", direction: "forward" });
      this.player.currentTime(this.player.currentTime() + 10);
    } else if (touchPercent < 0.35) {
      this.bigButtonGroup?.trigger({type: "display", direction: "backward" });
      this.player.currentTime(Math.max(0, this.player.currentTime() - 10));
    }
  }
}

const bigButtons = function (this: VideoJsPlayer) {
  this.addChild("BigButtonGroup");
};

// Register the plugin with video.js.
videojs.registerComponent("BigButtonGroup", BigButtonGroup);
videojs.registerComponent("BigPlayPauseButton", BigPlayPauseButton);
videojs.registerPlugin("bigButtons", bigButtons);
videojs.registerPlugin("touchControls", touchControls);

declare module "video.js" {
  export interface VideoJsPlayer {
    bigButtons: () => typeof bigButtons;
    touchControls: () => touchControls;
  }
}

export default bigButtons;
