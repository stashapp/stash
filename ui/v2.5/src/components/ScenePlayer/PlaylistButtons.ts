/* eslint-disable @typescript-eslint/naming-convention */
import VideoJS, { VideoJsPlayer } from "video.js";

const Button = VideoJS.getComponent("Button");

interface ControlOptions extends VideoJS.ComponentOptions {
  direction: "forward" | "back";
  parent: SkipButtonPlugin;
}

/**
 * A video.js plugin.
 *
 * In the plugin function, the value of `this` is a video.js `Player`
 * instance. You cannot rely on the player being in a "ready" state here,
 * depending on how the plugin is invoked. This may or may not be important
 * to you; if not, remove the wait for "ready"!
 *
 * @function skipButtons
 * @param    {Object} [options={}]
 *           An object of options left to the plugin author to define.
 */
class SkipButtonPlugin extends VideoJS.getPlugin("plugin") {
  onNext?: () => void | undefined;
  onPrevious?: () => void | undefined;

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

class SkipButton extends Button {
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

VideoJS.registerComponent("SkipButton", SkipButton);
VideoJS.registerPlugin("skipButtons", SkipButtonPlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    skipButtons: () => void | SkipButtonPlugin;
  }
}

export default SkipButtonPlugin;
