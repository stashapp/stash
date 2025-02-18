import videojs, { VideoJsPlayer } from "video.js";
import { WebVTT } from "videojs-vtt.js";

export interface IVTTThumbnailsOptions {
  /**
   * Source URL to use for thumbnails.
   */
  src?: string;
  /**
   * Whether to show the timestamp on hover.
   * @default false
   */
  showTimestamp?: boolean;
}

interface IVTTData {
  start: number;
  end: number;
  style: IVTTStyle | null;
}

interface IVTTStyle {
  background: string;
  width: string;
  height: string;
}

class VTTThumbnailsPlugin extends videojs.getPlugin("plugin") {
  private source: string | null;
  private showTimestamp: boolean;

  private progressBar?: HTMLElement;
  private thumbnailHolder?: HTMLDivElement;

  private showing = false;

  private vttData?: IVTTData[];
  private lastStyle?: IVTTStyle;

  constructor(player: VideoJsPlayer, options: IVTTThumbnailsOptions) {
    super(player, options);
    this.source = options.src ?? null;
    this.showTimestamp = options.showTimestamp ?? false;

    player.ready(() => {
      player.addClass("vjs-vtt-thumbnails");
      this.initializeThumbnails();
    });
  }

  src(source: string | null): void {
    this.resetPlugin();
    this.source = source;
    this.initializeThumbnails();
  }

  detach(): void {
    this.resetPlugin();
  }

  private resetPlugin() {
    this.showing = false;

    if (this.thumbnailHolder) {
      this.thumbnailHolder.remove();
      delete this.thumbnailHolder;
    }

    if (this.progressBar) {
      this.progressBar.removeEventListener(
        "pointerenter",
        this.onBarPointerEnter
      );
      this.progressBar.removeEventListener(
        "pointermove",
        this.onBarPointerMove
      );
      this.progressBar.removeEventListener(
        "pointerleave",
        this.onBarPointerLeave
      );

      delete this.progressBar;
    }

    delete this.vttData;
    delete this.lastStyle;
  }

  /**
   * Bootstrap the plugin.
   */
  private initializeThumbnails() {
    if (!this.source) {
      return;
    }

    const baseUrl = this.getBaseUrl();
    const url = this.getFullyQualifiedUrl(this.source, baseUrl);

    this.getVttFile(url).then((data) => {
      this.vttData = this.processVtt(data);
      this.setupThumbnailElement();
    });
  }

  /**
   * Builds a base URL should we require one.
   */
  private getBaseUrl() {
    return [
      window.location.protocol,
      "//",
      window.location.hostname,
      window.location.port ? ":" + window.location.port : "",
      window.location.pathname,
    ]
      .join("")
      .split(/([^\/]*)$/gi)[0];
  }

  /**
   * Grabs the contents of the VTT file.
   */
  private getVttFile(url: string): Promise<string> {
    return new Promise((resolve, reject) => {
      const req = new XMLHttpRequest();

      req.addEventListener("load", () => {
        resolve(req.responseText);
      });
      req.addEventListener("error", (e) => {
        reject(e);
      });
      req.open("GET", url);
      req.send();
    });
  }

  private setupThumbnailElement() {
    const progressBar = this.player.$(".vjs-progress-control") as HTMLElement;
    if (!progressBar) return;
    this.progressBar = progressBar;

    const thumbHolder = document.createElement("div");
    thumbHolder.setAttribute("class", "vjs-vtt-thumbnail-display");
    progressBar.appendChild(thumbHolder);
    this.thumbnailHolder = thumbHolder;

    if (!this.showTimestamp) {
      this.player.$(".vjs-mouse-display")?.classList.add("vjs-hidden");
    }

    progressBar.addEventListener("pointerover", this.onBarPointerEnter);
    progressBar.addEventListener("pointerout", this.onBarPointerLeave);
  }

  private onBarPointerEnter = () => {
    this.showThumbnailHolder();
    this.progressBar?.addEventListener("pointermove", this.onBarPointerMove);
  };

  private onBarPointerMove = (e: Event) => {
    const { progressBar } = this;
    if (!progressBar) return;

    this.showThumbnailHolder();
    this.updateThumbnailStyle(
      videojs.dom.getPointerPosition(progressBar, e).x,
      progressBar.offsetWidth
    );
  };

  private onBarPointerLeave = () => {
    this.hideThumbnailHolder();
    this.progressBar?.removeEventListener("pointermove", this.onBarPointerMove);
  };

  private getStyleForTime(time: number) {
    if (!this.vttData) return null;

    for (const element of this.vttData) {
      const item = element;

      if (time >= item.start && time < item.end) {
        return item.style;
      }
    }

    return null;
  }

  private showThumbnailHolder() {
    if (this.thumbnailHolder && !this.showing) {
      this.showing = true;
      this.thumbnailHolder.style.opacity = "1";
    }
  }

  private hideThumbnailHolder() {
    if (this.thumbnailHolder && this.showing) {
      this.showing = false;
      this.thumbnailHolder.style.opacity = "0";
    }
  }

  private updateThumbnailStyle(percent: number, width: number) {
    if (!this.thumbnailHolder) return;

    const duration = this.player.duration();
    const time = percent * duration;
    const currentStyle = this.getStyleForTime(time);

    if (!currentStyle) {
      this.hideThumbnailHolder();
      return;
    }

    const xPos = percent * width;
    const thumbnailWidth = parseInt(currentStyle.width, 10);
    const halfThumbnailWidth = thumbnailWidth >> 1;
    const marginRight = width - (xPos + halfThumbnailWidth);
    const marginLeft = xPos - halfThumbnailWidth;

    if (marginLeft > 0 && marginRight > 0) {
      this.thumbnailHolder.style.transform =
        "translateX(" + (xPos - halfThumbnailWidth) + "px)";
    } else if (marginLeft <= 0) {
      this.thumbnailHolder.style.transform = "translateX(" + 0 + "px)";
    } else if (marginRight <= 0) {
      this.thumbnailHolder.style.transform =
        "translateX(" + (width - thumbnailWidth) + "px)";
    }

    if (this.lastStyle && this.lastStyle === currentStyle) {
      return;
    }

    this.lastStyle = currentStyle;

    Object.assign(this.thumbnailHolder.style, currentStyle);
  }

  private processVtt(data: string) {
    const processedVtts: IVTTData[] = [];

    const parser = new WebVTT.Parser(window, WebVTT.StringDecoder());
    parser.oncue = (cue: VTTCue) => {
      processedVtts.push({
        start: cue.startTime,
        end: cue.endTime,
        style: this.getVttStyle(cue.text),
      });
    };
    parser.parse(data);
    parser.flush();

    return processedVtts;
  }

  private getFullyQualifiedUrl(path: string, base: string) {
    if (path.indexOf("//") >= 0) {
      // We have a fully qualified path.
      return path;
    }

    if (base.indexOf("//") === 0) {
      // We don't have a fully qualified path, but need to
      // be careful with trimming.
      return [base.replace(/\/$/gi, ""), this.trim(path, "/")].join("/");
    }

    if (base.indexOf("//") > 0) {
      // We don't have a fully qualified path, and should
      // trim both sides of base and path.
      return [this.trim(base, "/"), this.trim(path, "/")].join("/");
    }

    // If all else fails.
    return path;
  }

  private getPropsFromDef(def: string) {
    const match = def.match(/^([^#]*)#xywh=(\d+),(\d+),(\d+),(\d+)$/i);
    if (!match) return null;

    return {
      image: match[1],
      x: match[2],
      y: match[3],
      w: match[4],
      h: match[5],
    };
  }

  private getVttStyle(vttImageDef: string) {
    // If there isn't a protocol, use the VTT source URL.
    let baseSplit: string;

    if (this.source === null) {
      baseSplit = this.getBaseUrl();
    } else if (this.source.indexOf("//") >= 0) {
      baseSplit = this.source.split(/([^\/]*)$/gi)[0];
    } else {
      baseSplit = this.getBaseUrl() + this.source.split(/([^\/]*)$/gi)[0];
    }

    vttImageDef = this.getFullyQualifiedUrl(vttImageDef, baseSplit);

    const imageProps = this.getPropsFromDef(vttImageDef);
    if (!imageProps) return null;

    return {
      background:
        'url("' +
        imageProps.image +
        '") no-repeat -' +
        imageProps.x +
        "px -" +
        imageProps.y +
        "px",
      width: imageProps.w + "px",
      height: imageProps.h + "px",
    };
  }

  /**
   * trim
   *
   * @param  str      source string
   * @param  charlist characters to trim from text
   * @return          trimmed string
   */
  private trim(str: string, charlist: string) {
    let whitespace = [
      " ",
      "\n",
      "\r",
      "\t",
      "\f",
      "\x0b",
      "\xa0",
      "\u2000",
      "\u2001",
      "\u2002",
      "\u2003",
      "\u2004",
      "\u2005",
      "\u2006",
      "\u2007",
      "\u2008",
      "\u2009",
      "\u200a",
      "\u200b",
      "\u2028",
      "\u2029",
      "\u3000",
    ].join("");
    let l = 0;

    str += "";
    if (charlist) {
      whitespace = (charlist + "").replace(/([[\]().?/*{}+$^:])/g, "$1");
    }

    l = str.length;
    for (let i = 0; i < l; i++) {
      if (whitespace.indexOf(str.charAt(i)) === -1) {
        str = str.substring(i);
        break;
      }
    }

    l = str.length;
    for (let i = l - 1; i >= 0; i--) {
      if (whitespace.indexOf(str.charAt(i)) === -1) {
        str = str.substring(0, i + 1);
        break;
      }
    }
    return whitespace.indexOf(str.charAt(0)) === -1 ? str : "";
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("vttThumbnails", VTTThumbnailsPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    vttThumbnails: () => VTTThumbnailsPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    vttThumbnails?: IVTTThumbnailsOptions;
  }
}

export default VTTThumbnailsPlugin;
