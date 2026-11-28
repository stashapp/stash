import videojs from "video.js";
import type VideoJsPlayer from "video.js/dist/types/player";
import type Button from "video.js/dist/types/button";
import type Plugin from "video.js/dist/types/plugin";

interface IAutostartButtonOptions {
  enabled?: boolean;
}

type AutostartButtonOptions = ConstructorParameters<typeof Button>[1] & {
  autostartEnabled: boolean;
};

class AutostartButton extends (videojs.getComponent(
  "Button"
) as typeof Button) {
  private autostartEnabled: boolean;

  constructor(player: VideoJsPlayer, options: AutostartButtonOptions) {
    super(player, options);
    this.autostartEnabled = options.autostartEnabled;
    this.updateIcon();
  }

  buildCSSClass() {
    return `vjs-autostart-button ${super.buildCSSClass()}`;
  }

  private updateIcon() {
    this.removeClass("vjs-icon-play-circle");
    this.removeClass("vjs-icon-cancel");

    if (this.autostartEnabled) {
      this.addClass("vjs-icon-play-circle");
      this.controlText(this.localize("Auto-start enabled (click to disable)"));
    } else {
      this.addClass("vjs-icon-cancel");
      this.controlText(this.localize("Auto-start disabled (click to enable)"));
    }
  }

  handleClick(event: Event) {
    // Prevent the click from bubbling up and affecting the video player
    event.stopPropagation();

    this.autostartEnabled = !this.autostartEnabled;
    this.updateIcon();
    this.trigger("autostartchanged", { enabled: this.autostartEnabled });
  }

  public setEnabled(enabled: boolean) {
    this.autostartEnabled = enabled;
    this.updateIcon();
  }
}

class AutostartButtonPlugin extends (videojs.getPlugin(
  "plugin"
) as typeof Plugin) {
  private button: AutostartButton;
  private autostartEnabled: boolean;
  updateAutoStart: (enabled: boolean) => Promise<void> = () => {
    return Promise.resolve();
  };

  constructor(player: VideoJsPlayer, options?: IAutostartButtonOptions) {
    super(player);

    this.autostartEnabled = options?.enabled ?? false;

    this.button = new AutostartButton(player, {
      autostartEnabled: this.autostartEnabled,
    });

    player.ready(() => {
      this.ready();
    });
  }

  private ready() {
    // Add button to control bar, before the fullscreen button
    const { controlBar } = this.player;
    const fullscreenToggle = controlBar.getChild("fullscreenToggle");
    if (fullscreenToggle) {
      controlBar.addChild(this.button);
      controlBar.el().insertBefore(this.button.el(), fullscreenToggle.el());
    } else {
      controlBar.addChild(this.button);
    }

    // Listen for changes
    this.button.on(
      "autostartchanged",
      (_: Event, data: { enabled: boolean }) => {
        this.autostartEnabled = data.enabled;
        this.updateAutoStart(this.autostartEnabled);
      }
    );
  }

  public isEnabled(): boolean {
    return this.autostartEnabled;
  }

  public getEnabled(): boolean {
    return this.autostartEnabled;
  }

  public setEnabled(enabled: boolean) {
    this.autostartEnabled = enabled;
    this.button.setEnabled(enabled);
  }

  public syncWithConfig(configEnabled: boolean) {
    // Sync button state with external config changes
    if (this.autostartEnabled !== configEnabled) {
      this.setEnabled(configEnabled);
    }
  }
}

// Register the plugin with video.js.
videojs.registerComponent("AutostartButton", AutostartButton);
videojs.registerPlugin("autostartButton", AutostartButtonPlugin);

declare module "video.js/dist/types/player" {
  export default interface Player {
    autostartButton: () => AutostartButtonPlugin;
  }
}

export default AutostartButtonPlugin;
