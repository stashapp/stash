/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import "videojs-vr";

export interface VRMenuOptions {
  /**
   * Whether to show the vr button.
   * @default false
   */
  showButton?: boolean;
}

enum VRType {
  Spherical = "360",
  Off = "Off",
}

const vrTypeProjection = {
  [VRType.Spherical]: "360",
  [VRType.Off]: "NONE",
};

function isVrDevice() {
  return navigator.userAgent.match(/oculusbrowser|\svr\s/i);
}

class VRMenuItem extends videojs.getComponent("MenuItem") {
  public type: VRType;
  public isSelected = false;

  constructor(parent: VRMenuButton, type: VRType) {
    const options = {} as videojs.MenuItemOptions;
    options.selectable = true;
    options.multiSelectable = false;
    options.label = type;

    super(parent.player(), options);

    this.type = type;

    this.addClass("vjs-source-menu-item");
  }

  selected(selected: boolean): void {
    super.selected(selected);
    this.isSelected = selected;
  }

  handleClick() {
    if (this.isSelected) return;

    this.trigger("selected");
  }
}

class VRMenuButton extends videojs.getComponent("MenuButton") {
  private items: VRMenuItem[] = [];
  private selectedType: VRType = VRType.Off;

  constructor(player: VideoJsPlayer) {
    super(player);
    this.setTypes();
  }

  private onSelected(item: VRMenuItem) {
    this.selectedType = item.type;

    this.items.forEach((i) => {
      i.selected(i.type === this.selectedType);
    });

    this.trigger("typeselected", item.type);
  }

  public setTypes() {
    this.items = Object.values(VRType).map((type) => {
      const item = new VRMenuItem(this, type);

      item.on("selected", () => {
        this.onSelected(item);
      });

      return item;
    });
    this.update();
  }

  createEl() {
    return videojs.dom.createEl("div", {
      className:
        "vjs-vr-selector vjs-menu-button vjs-menu-button-popup vjs-control vjs-button",
    });
  }

  createItems() {
    if (this.items === undefined) return [];

    for (const item of this.items) {
      item.selected(item.type === this.selectedType);
    }

    return this.items;
  }
}

class VRMenuPlugin extends videojs.getPlugin("plugin") {
  private menu: VRMenuButton;

  maybeToggleVrOffClass(projection: string) {
    const playerVideoContainer = document.getElementById("VideoJsPlayer");
    if (!playerVideoContainer) {
      return;
    }

    if (projection == vrTypeProjection[VRType.Off]) {
      playerVideoContainer.classList.add("vjs-vr-off");
    } else {
      playerVideoContainer.classList.remove("vjs-vr-off");
    }
  }

  constructor(player: VideoJsPlayer, options: VRMenuOptions) {
    super(player);

    this.menu = new VRMenuButton(player);

    if (isVrDevice() || !options.showButton) return;

    this.menu.on("typeselected", (_, type: VRType) => {
      const projection = vrTypeProjection[type];
      this.maybeToggleVrOffClass(projection);
      player.vr({ projection });
      player.load();
    });

    player.on("ready", () => {
      const { controlBar } = player;
      const fullscreenToggle = controlBar.getChild("fullscreenToggle")!.el();
      controlBar.addChild(this.menu);
      controlBar.el().insertBefore(this.menu.el(), fullscreenToggle);
    });
  }
}

// Register the plugin with video.js.
videojs.registerComponent("VRMenuButton", VRMenuButton);
videojs.registerPlugin("vrMenu", VRMenuPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    vrMenu: () => VRMenuPlugin;
    vr: (options: Object) => void;
  }
  interface VideoJsPlayerPluginOptions {
    vrMenu?: VRMenuOptions;
  }
}

export default VRMenuPlugin;
