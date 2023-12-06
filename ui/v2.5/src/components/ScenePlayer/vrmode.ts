/* eslint-disable @typescript-eslint/naming-convention */
import videojs, { VideoJsPlayer } from "video.js";
import "videojs-vr";
// separate type import, otherwise typescript elides the above import
// and the plugin does not get initialized
import type { ProjectionType, Plugin as VideoJsVRPlugin } from "videojs-vr";

export interface VRMenuOptions {
  /**
   * Whether to show the vr button.
   * @default false
   */
  showButton?: boolean;
}

enum VRType {
  LR180 = "180 LR",
  TB360 = "360 TB",
  Mono360 = "360 Mono",
  Off = "Off",
}

const vrTypeProjection: Record<VRType, ProjectionType> = {
  [VRType.LR180]: "180_LR",
  [VRType.TB360]: "360_TB",
  [VRType.Mono360]: "360",
  [VRType.Off]: "NONE",
};

function isVrDevice() {
  return navigator.userAgent.match(/oculusbrowser|\svr\s/i);
}

class VRMenuItem extends videojs.getComponent("MenuItem") {
  public type: VRType;
  public isSelected = false;

  constructor(parent: VRMenuButton, type: VRType) {
    const options: videojs.MenuItemOptions = {};
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
  private showButton: boolean;
  private vr?: VideoJsVRPlugin;

  constructor(player: VideoJsPlayer, options: VRMenuOptions) {
    super(player);

    this.menu = new VRMenuButton(player);
    this.showButton = options.showButton ?? false;

    if (isVrDevice()) return;

    this.vr = this.player.vr();

    this.menu.on("typeselected", (_, type: VRType) => {
      this.loadVR(type);
    });

    player.on("ready", () => {
      if (this.showButton) {
        this.addButton();
      }
    });
  }

  private loadVR(type: VRType) {
    const projection = vrTypeProjection[type];
    this.vr?.setProjection(projection);
    this.vr?.init();
  }

  private addButton() {
    const { controlBar } = this.player;
    const fullscreenToggle = controlBar.getChild("fullscreenToggle")!.el();
    controlBar.addChild(this.menu);
    controlBar.el().insertBefore(this.menu.el(), fullscreenToggle);
  }

  private removeButton() {
    const { controlBar } = this.player;
    controlBar.removeChild(this.menu);
  }

  public setShowButton(showButton: boolean) {
    if (isVrDevice()) return;

    if (showButton === this.showButton) return;

    this.showButton = showButton;
    if (showButton) {
      this.addButton();
    } else {
      this.removeButton();
      this.loadVR(VRType.Off);
    }
  }
}

// Register the plugin with video.js.
videojs.registerComponent("VRMenuButton", VRMenuButton);
videojs.registerPlugin("vrMenu", VRMenuPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    vrMenu: () => VRMenuPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    vrMenu?: VRMenuOptions;
  }
}

export default VRMenuPlugin;
