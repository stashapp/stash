/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import localForage from "localforage";

const AUTOSTART_KEY = "video-autostart-enabled";

interface IAutostartButtonOptions {
  enabled?: boolean;
}

interface AutostartButtonOptions extends videojs.ComponentOptions {
  autostartEnabled: boolean;
}

class AutostartButton extends videojs.getComponent("Button") {
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

  handleClick() {
    this.autostartEnabled = !this.autostartEnabled;
    this.updateIcon();
    this.trigger("autostartchanged", { enabled: this.autostartEnabled });
  }

  public setEnabled(enabled: boolean) {
    this.autostartEnabled = enabled;
    this.updateIcon();
  }
}

class AutostartButtonPlugin extends videojs.getPlugin("plugin") {
  private button: AutostartButton;
  private autostartEnabled: boolean;
  private loaded: boolean = false;

  constructor(player: VideoJsPlayer, options?: IAutostartButtonOptions) {
    super(player, options);

    this.autostartEnabled = options?.enabled ?? false;
    
    // Load the saved preference immediately
    this.loadPreference();

    this.button = new AutostartButton(player, {
      autostartEnabled: this.autostartEnabled,
    });

    player.ready(() => {
      this.ready();
    });
  }

  private async loadPreference() {
    const value = await localForage.getItem<boolean>(AUTOSTART_KEY);
    if (value !== null) {
      this.autostartEnabled = value;
      if (this.button) {
        this.button.setEnabled(value);
      }
    }
    this.loaded = true;
  }

  private ready() {
    // Add button to control bar, before the fullscreen button
    const controlBar = this.player.controlBar;
    const fullscreenToggle = controlBar.getChild("fullscreenToggle");
    if (fullscreenToggle) {
      controlBar.addChild(this.button);
      controlBar.el().insertBefore(this.button.el(), fullscreenToggle.el());
    } else {
      controlBar.addChild(this.button);
    }

    // Listen for changes
    this.button.on("autostartchanged", (_, data: { enabled: boolean }) => {
      this.autostartEnabled = data.enabled;
      localForage.setItem(AUTOSTART_KEY, data.enabled);
    });
  }

  public isEnabled(): boolean {
    return this.autostartEnabled;
  }

  public async getEnabled(): Promise<boolean> {
    // Wait for the preference to be loaded if it hasn't been yet
    if (!this.loaded) {
      await this.loadPreference();
    }
    return this.autostartEnabled;
  }
}

// Register the plugin with video.js.
videojs.registerComponent("AutostartButton", AutostartButton);
videojs.registerPlugin("autostartButton", AutostartButtonPlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    autostartButton: () => AutostartButtonPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    autostartButton?: IAutostartButtonOptions;
  }
}

export default AutostartButtonPlugin;

