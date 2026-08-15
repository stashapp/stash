import videojs, { VideoJsPlayer } from "video.js";
import type { AbLoopPluginApi } from "./util";

function getAbLoopApi(player: VideoJsPlayer) {
  return player.abLoopPlugin as unknown as AbLoopPluginApi | undefined;
}

interface AbLoopToggleOptions {
  /**
   * Whether to show the button at all (mirrors showAbLoopControls).
   * @default false
   */
  enabled?: boolean;
}

// Icon-only counterpart to the vendor videojs-abloop plugin's own
// "LOOP ON"/"Loop off" text button, shown on mobile only (see styles.scss)
// where that text button is hidden. Uses the normal Button structure
// (icon-placeholder + control-text, same as every built-in button) rather
// than reusing the vendor's own button element, since that element's text
// gets overwritten wholesale (element.textContent = ...) on every change,
// which would wipe out any icon markup injected into it.
class AbLoopToggleButton extends videojs.getComponent("Button") {
  constructor(player: VideoJsPlayer, options?: videojs.ComponentOptions) {
    super(player, options);
    this.controlText(this.localize("Toggle AB loop"));

    const refresh = () => {
      const enabled = getAbLoopApi(player)?.getOptions().enabled;
      if (enabled) this.addClass("vjs-ab-loop-active");
      else this.removeClass("vjs-ab-loop-active");
    };
    player.on("abloopchange", refresh);
    refresh();
  }

  buildCSSClass() {
    return `vjs-ab-loop-toggle ${super.buildCSSClass()}`;
  }

  handleClick(event: videojs.EventTarget.Event) {
    // Prevent the click from bubbling up and affecting the video player
    event.stopPropagation();

    const player = this.player() as unknown as VideoJsPlayer;
    const api = getAbLoopApi(player);
    if (!api) return;

    const opts = api.getOptions();
    api.setOptions({ ...opts, enabled: !opts.enabled });
  }
}

class AbLoopTogglePlugin extends videojs.getPlugin("plugin") {
  private button: AbLoopToggleButton;
  private added = false;

  constructor(player: VideoJsPlayer, options?: AbLoopToggleOptions) {
    super(player, options);

    this.button = new AbLoopToggleButton(player);

    player.on("ready", () => {
      if (options?.enabled) this.addButton();
    });
  }

  private addButton() {
    if (this.added) return;
    const { controlBar } = this.player;
    const fullscreenToggle = controlBar.getChild("fullscreenToggle");
    controlBar.addChild(this.button);
    if (fullscreenToggle) {
      controlBar.el().insertBefore(this.button.el(), fullscreenToggle.el());
    }
    this.added = true;
  }
}

// Register the plugin with video.js.
videojs.registerComponent("AbLoopToggleButton", AbLoopToggleButton);
videojs.registerPlugin("abLoopToggle", AbLoopTogglePlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    abLoopToggle: () => AbLoopTogglePlugin;
  }
  interface VideoJsPlayerPluginOptions {
    abLoopToggle?: AbLoopToggleOptions;
  }
}

export default AbLoopTogglePlugin;
