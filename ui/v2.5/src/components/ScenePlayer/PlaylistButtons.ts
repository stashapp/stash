import videojs from "video.js";
import type VideoJsPlayer from "video.js/dist/types/player";
import type Button from "video.js/dist/types/button";
import type Plugin from "video.js/dist/types/plugin";

type ControlOptions = ConstructorParameters<typeof Button>[1] & {
  direction: "forward" | "back";
  parent: SkipButtonPlugin;
};

class SkipButtonPlugin extends (videojs.getPlugin("plugin") as typeof Plugin) {
  onNext?: () => void;
  onPrevious?: () => void;

  constructor(player: VideoJsPlayer) {
    super(player);
    player.ready(() => {
      this.ready();
    });
  }

  public setForwardHandler(handler?: () => void) {
    this.onNext = handler;
    if (handler !== undefined) this.player.addClass("vjs-skip-buttons-next");
    else this.player.removeClass("vjs-skip-buttons-next");
  }

  public setBackwardHandler(handler?: () => void) {
    this.onPrevious = handler;
    if (handler !== undefined) this.player.addClass("vjs-skip-buttons-prev");
    else this.player.removeClass("vjs-skip-buttons-prev");
  }

  handleForward() {
    this.onNext?.();
  }

  handleBackward() {
    this.onPrevious?.();
  }

  ready() {
    this.player.addClass("vjs-skip-buttons");

    this.player.controlBar.addChild(
      "skipButton",
      {
        direction: "forward",
        parent: this,
      },
      1
    );

    this.player.controlBar.addChild(
      "skipButton",
      {
        direction: "back",
        parent: this,
      },
      0
    );
  }
}

class SkipButton extends (videojs.getComponent("button") as typeof Button) {
  private parentPlugin: SkipButtonPlugin;
  private direction: "forward" | "back";

  constructor(player: VideoJsPlayer, options: ControlOptions) {
    super(player, options);
    this.parentPlugin = options.parent;
    this.direction = options.direction;
    if (options.direction === "forward") {
      this.controlText(this.localize("Skip to next video"));
      this.addClass(`vjs-icon-next-item`);
    } else if (options.direction === "back") {
      this.controlText(this.localize("Skip to previous video"));
      this.addClass(`vjs-icon-previous-item`);
    }
  }

  /**
   * Return button class names
   */
  buildCSSClass() {
    return `vjs-skip-button ${super.buildCSSClass()}`;
  }

  /**
   * Seek with the button's configured offset
   */
  handleClick() {
    if (this.direction === "forward") this.parentPlugin.handleForward();
    else this.parentPlugin.handleBackward();
  }
}

videojs.registerComponent("SkipButton", SkipButton);
videojs.registerPlugin("skipButtons", SkipButtonPlugin);

declare module "video.js/dist/types/player" {
  export default interface Player {
    skipButtons: () => SkipButtonPlugin;
  }
}

export default SkipButtonPlugin;
