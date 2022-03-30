import videojs, { VideoJsPlayer } from "video.js";

interface ISource extends videojs.Tech.SourceObject {
  label?: string;
  selected?: boolean;
  sortIndex?: number;
}

const MenuButton = videojs.getComponent("MenuButton");
const MenuItem = videojs.getComponent("MenuItem");

class SourceMenuItem extends MenuItem {
  constructor(player: VideoJsPlayer, options: videojs.MenuItemOptions) {
    options.selectable = true;
    options.multiSelectable = false;

    super(player, options);
  }
}

class SourceMenuButton extends MenuButton {
  buildCSSClass() {
    return MenuButton.prototype.buildCSSClass.call(this) + " vjs-icon-cog";
  }

  createEl() {
    return videojs.dom.createEl("div", {
      className:
        "vjs-source-selector vjs-menu-button vjs-menu-button-popup vjs-control vjs-button",
    });
  }

  createItems() {
    const player = this.player();
    // slice so that we don't alter the order of the original array
    const sources = player.currentSources().slice() as ISource[];

    sources.sort((a, b) => (a.sortIndex ?? 0) - (b.sortIndex ?? 0));

    const hasSelected = sources.some((source) => source.selected);

    return sources.map((source, index) => {
      const label = source.label || source.type;
      const item = new SourceMenuItem(this.player(), {
        label: label,
        selected: source.selected || (!hasSelected && index === 0),
      });

      item.on("click", function () {
        // populate source sortIndex first if not present
        const currentSources = (player.currentSources() as ISource[]).map(
          (src, i) => {
            return {
              ...src,
              sortIndex: src.sortIndex ?? i,
              selected: false,
            };
          }
        );

        // put the selected source at the top of the list
        const selectedIndex = currentSources.findIndex(
          (src) => src.sortIndex === index
        );
        const selectedSrc = currentSources.splice(selectedIndex, 1)[0];
        selectedSrc.selected = true;
        currentSources.unshift(selectedSrc);

        const currentTime = player.currentTime();

        /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
        (player as any).clearOffsetDuration();
        player.src(currentSources);
        player.currentTime(currentTime);
        player.play();
      });

      item.addClass("vjs-source-menu-item");

      return item;
    });
  }
}

const sourceSelector = function (this: VideoJsPlayer) {
  const player = this;

  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  const PlayerConstructor = this.constructor as any;
  if (!PlayerConstructor.__sourceSelector) {
    PlayerConstructor.__sourceSelector = {
      selectSource: PlayerConstructor.prototype.selectSource,
    };
  }

  videojs.registerComponent("SourceMenuButton", SourceMenuButton);

  player.on("loadedmetadata", function () {
    const { controlBar } = player;
    const fullscreenToggle = controlBar.getChild("fullscreenToggle")!.el();

    const existingMenuButton = controlBar.getChild("SourceMenuButton");
    if (existingMenuButton) controlBar.removeChild(existingMenuButton);

    const menuButton = controlBar.addChild("SourceMenuButton");

    controlBar.el().insertBefore(menuButton.el(), fullscreenToggle);
  });
};

// Register the plugin with video.js.
videojs.registerPlugin("sourceSelector", sourceSelector);

export default sourceSelector;
