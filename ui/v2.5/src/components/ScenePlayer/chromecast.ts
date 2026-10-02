import videojs, { VideoJsPlayer } from "video.js";
import {
  type ICastFile,
  type ICastSource,
  type ICastStream,
  isLocalHost,
  pickCastSource,
  rewriteCastUrl,
} from "src/utils/castMedia";
import { CastSession, isCastSenderSupported } from "./cast-session";

export interface IChromecastScene {
  title: string;
  streams: ICastStream[];
  file?: ICastFile;
}

export interface IChromecastOptions {
  enabled: boolean;
  // Address the device should fetch from, from systemStatus.localIPs.
  lanIp: string;
  scene: IChromecastScene | null;
  onError: (message: string) => void;
}

const DISABLED: IChromecastOptions = {
  enabled: false,
  lanIp: "",
  scene: null,
  onError: () => undefined,
};

class ChromecastButton extends videojs.getComponent("Button") {
  public onCastClick: () => void = () => undefined;

  buildCSSClass() {
    return `vjs-chromecast-button ${super.buildCSSClass()}`;
  }

  handleClick() {
    this.onCastClick();
  }
}

/**
 * Casts the current scene to a Chromecast.
 *
 * Google's Cast Application Framework is driven directly rather than through a
 * Video.js tech: a tech hands the device whatever the browser is playing,
 * which is a stream chosen for a browser, at an address only this machine can
 * resolve.
 */
class ChromecastPlugin extends videojs.getPlugin("plugin") {
  private button: ChromecastButton;
  private session = new CastSession();
  private options = DISABLED;
  private busy = false;

  constructor(player: VideoJsPlayer) {
    super(player);

    this.button = new ChromecastButton(player);
    this.button.onCastClick = () => {
      void this.toggleCast();
    };

    // Set before ready() so a device that connects while the player is still
    // setting up is not missed.
    this.session.onConnectionChange = () => this.updateButton();

    player.ready(() => {
      const { controlBar } = player;
      controlBar.addChild(this.button);

      const fullscreenToggle = controlBar.getChild("fullscreenToggle");
      if (fullscreenToggle) {
        controlBar.el().insertBefore(this.button.el(), fullscreenToggle.el());
      }

      this.applyOptions();
    });

    player.on("dispose", () => this.session.destroy());
  }

  public configure(options: IChromecastOptions) {
    this.options = options;
    this.applyOptions();
  }

  private applyOptions() {
    if (this.options.enabled) {
      this.button.show();
      void this.session.watch();
    } else {
      this.button.hide();
    }

    this.updateButton();
  }

  private updateButton() {
    const casting = this.session.isConnected();

    this.button.toggleClass("vjs-chromecast-casting", casting);
    this.button.controlText(
      this.player.localize(casting ? "Stop casting" : "Cast")
    );
  }

  /**
   * Why casting cannot start, or "" when it can. Checked before opening the
   * picker, so the user gets a reason instead of a device that connects and
   * then sits there.
   */
  private blockedReason(source: ICastSource | null): string {
    if (!isCastSenderSupported()) {
      return this.player.localize("Casting needs Chrome, Edge, or Opera");
    }

    // The sender SDK only initialises on HTTPS or localhost.
    if (!window.isSecureContext) {
      return this.player.localize(
        "Casting needs Stash to be open over HTTPS or as http://localhost"
      );
    }

    if (!source) {
      return this.player.localize(
        "This scene has no stream a Chromecast can play"
      );
    }

    // Parsed once already by pickCastSource, so this cannot throw.
    const { hostname } = new URL(source.url, window.location.href);
    if (isLocalHost(hostname) && !this.options.lanIp) {
      return this.player.localize(
        "A Chromecast cannot fetch localhost, and Stash found no LAN address for this machine"
      );
    }

    return "";
  }

  private async toggleCast() {
    if (this.busy) return;

    if (this.session.isConnected()) {
      this.session.disconnect();
      return;
    }

    const { scene } = this.options;
    const source = scene ? pickCastSource(scene.streams, scene.file) : null;

    const reason = this.blockedReason(source);
    if (reason || !scene || !source) {
      this.options.onError(reason);
      return;
    }

    this.busy = true;
    try {
      await this.session.connect();
      await this.session.load({
        url: rewriteCastUrl(source.url, this.options.lanIp),
        contentType: source.contentType,
        title: scene.title,
        duration: scene.file?.duration,
        startTime: this.player.currentTime() || 0,
        autoplay: !this.player.paused(),
      });
    } catch (err) {
      this.reportError(err);
    } finally {
      this.busy = false;
      this.updateButton();
    }
  }

  private reportError(err: unknown) {
    const message =
      (err as { message?: string })?.message ??
      (err as { description?: string })?.description ??
      String(err);

    // Dismissing the device picker rejects too; that is not an error.
    if (/cancel/i.test(message)) return;

    this.options.onError(message);
  }
}

videojs.registerComponent("ChromecastButton", ChromecastButton);
videojs.registerPlugin("chromecast", ChromecastPlugin);

declare module "video.js" {
  interface VideoJsPlayer {
    chromecast: () => ChromecastPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    chromecast?: Record<string, never>;
  }
}

export default ChromecastPlugin;
