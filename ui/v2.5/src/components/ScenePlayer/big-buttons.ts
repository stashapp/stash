import videojs, { VideoJsPlayer } from "video.js";
import type { AbLoopOptions, AbLoopPluginApi } from "./util";

// prettier-ignore
const BigPlayButton = videojs.getComponent(
  "BigPlayButton"
) as unknown as typeof videojs.BigPlayButton;

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

function getAbLoopApi(player: VideoJsPlayer) {
  return player.abLoopPlugin as unknown as AbLoopPluginApi | undefined;
}

interface AbBoundButtonOptions extends videojs.ComponentOptions {
  bound: "start" | "end";
}

// Sets the AB-loop start/end point to the current playback position. Lives
// in the floating button group (mobile only, see styles.scss) as the
// touch-friendly counterpart to the Start/End buttons the videojs-abloop
// plugin itself renders in the control bar on desktop.
class AbBoundButton extends videojs.getComponent("Button") {
  private bound: "start" | "end";

  constructor(player: VideoJsPlayer, options: AbBoundButtonOptions) {
    super(player, options);
    this.bound = options.bound;
    this.controlText(
      this.localize(
        this.bound === "start" ? "Set AB loop start" : "Set AB loop end"
      )
    );

    // recolor once this bound has actually been set (as opposed to still
    // sitting at the vendor plugin's start:0/end:false defaults), so
    // there's some feedback beyond remembering whether it was tapped -
    // see styles.scss's vjs-ab-bound-active
    const refresh = () => {
      const opts = getAbLoopApi(player)?.getOptions();
      const isSet =
        this.bound === "start"
          ? typeof opts?.start === "number" && opts.start > 0
          : typeof opts?.end === "number";
      this.toggleClass("vjs-ab-bound-active", isSet);
    };
    player.on("abloopchange", refresh);
    refresh();
  }

  buildCSSClass() {
    // called synchronously from within super(), before the constructor body
    // below runs - so this.bound isn't set yet at this point. this.options_
    // is though (Component's own constructor populates it before calling
    // createEl(), which is what triggers this).
    const bound = (this.options_ as AbBoundButtonOptions).bound;
    return `vjs-ab-${bound}-button ${super.buildCSSClass()}`;
  }

  handleClick(event: videojs.EventTarget.Event) {
    // Prevent the click from bubbling up and affecting the video player
    event.stopPropagation();

    const player = this.player() as unknown as VideoJsPlayer;
    const api = getAbLoopApi(player);
    if (!api) return;

    const opts = api.getOptions();
    const changes: Partial<AbLoopOptions> =
      this.bound === "start"
        ? { start: player.currentTime() }
        : { end: player.currentTime() };
    api.setOptions({ ...opts, ...changes });
  }
}

interface BigButtonGroupOptions extends videojs.ComponentOptions {
  showAbLoop?: boolean;
}

class BigButtonGroup extends videojs.getComponent("Component") {
  constructor(player: VideoJsPlayer, options?: BigButtonGroupOptions) {
    super(player, options);

    if (options?.showAbLoop) {
      this.addChild("AbBoundButton", { bound: "start" });
    }

    this.addChild("seekButton", {
      direction: "back",
      seconds: 10,
    });

    this.addChild("BigPlayPauseButton");

    this.addChild("seekButton", {
      direction: "forward",
      seconds: 10,
    });

    if (options?.showAbLoop) {
      this.addChild("AbBoundButton", { bound: "end" });
    }
  }

  createEl() {
    return super.createEl("div", {
      className: "vjs-big-button-group",
    });
  }
}

interface BigButtonsOptions {
  showAbLoop?: boolean;
}

class BigButtonsPlugin extends videojs.getPlugin("plugin") {
  constructor(player: VideoJsPlayer, options?: BigButtonsOptions) {
    super(player, options);

    player.ready(() => {
      player.addChild("BigButtonGroup", {
        showAbLoop: options?.showAbLoop ?? false,
      });
    });
  }
}

// Register the plugin with video.js.
videojs.registerComponent("BigButtonGroup", BigButtonGroup);
videojs.registerComponent("BigPlayPauseButton", BigPlayPauseButton);
videojs.registerComponent("AbBoundButton", AbBoundButton);
videojs.registerPlugin("bigButtons", BigButtonsPlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    bigButtons: () => BigButtonsPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    bigButtons?: BigButtonsOptions;
  }
}

export default BigButtonsPlugin;
