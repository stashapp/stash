/**
 * VideoJS Settings Menu Plugin
 *
 * Custom settings menu for the video player providing:
 * - Quality/source selection with auto-fallback on errors
 * - Playback speed control
 * - Subtitle/CC track selection
 *
 * This plugin auto-registers when imported.
 */
import videojs, { VideoJsPlayer } from "video.js";

// Icons as SVG data URIs
const ICONS = {
  quality: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18">
    <rect x="3" y="6" width="18" height="12" rx="2" stroke="currentColor" stroke-width="2" fill="none"/>
    <text x="12" y="15" text-anchor="middle" font-size="8" font-weight="bold" fill="currentColor">HQ</text>
  </svg>`,
  speed: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18">
    <circle cx="12" cy="12" r="9"/>
    <path d="M12 6v6l4 2"/>
  </svg>`,
  subtitles: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18">
    <rect x="2" y="5" width="20" height="14" rx="2" stroke="currentColor" stroke-width="2" fill="none"/>
    <text x="12" y="15" text-anchor="middle" font-size="7" font-weight="bold" fill="currentColor">CC</text>
  </svg>`,
  chevron: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
    <polyline points="9 18 15 12 9 6"/>
  </svg>`,
  back: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
    <polyline points="15 18 9 12 15 6"/>
  </svg>`,
};

export interface ISource extends videojs.Tech.SourceObject {
  label?: string;
  errored?: boolean;
}

type SettingsSection = "main" | "quality" | "speed" | "subtitles";

class SettingsMenuPlugin extends videojs.getPlugin("plugin") {
  private menuButton: HTMLElement;
  private menuContainer: HTMLElement;
  private sources: ISource[] = [];
  private selectedSourceIndex = 0;
  private selectedSpeed = 1;
  private currentSection: SettingsSection = "main";
  private isOpen = false;
  private cleanupTextTracks: HTMLTrackElement[] = [];
  private manualTextTracks: HTMLTrackElement[] = [];
  private manuallySelected = false;

  private readonly playbackRates = [0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2];

  constructor(player: VideoJsPlayer) {
    super(player);

    // Create settings button
    this.menuButton = this.createSettingsButton();
    this.menuContainer = this.createMenuContainer();

    player.on("ready", () => {
      const { controlBar } = player;
      const fullscreenToggle = controlBar.getChild("fullscreenToggle")?.el();

      // Insert settings button before fullscreen
      if (fullscreenToggle) {
        controlBar.el().insertBefore(this.menuButton, fullscreenToggle);
        controlBar.el().insertBefore(this.menuContainer, fullscreenToggle);
      } else {
        controlBar.el().appendChild(this.menuButton);
        controlBar.el().appendChild(this.menuContainer);
      }

      // Hide default playback rate and subs-caps buttons
      const playbackRateBtn = controlBar.getChild("playbackRateMenuButton");
      if (playbackRateBtn) {
        playbackRateBtn.hide();
      }
      const subsCapsBtn = controlBar.getChild("subsCapsButton");
      if (subsCapsBtn) {
        subsCapsBtn.hide();
      }
    });

    // Stop propagation for all clicks inside the menu container
    // This prevents the document click handler from closing the menu
    this.menuContainer.addEventListener("click", (e) => {
      e.stopPropagation();
    });

    // Close menu when clicking outside
    document.addEventListener("click", (e) => {
      if (
        !this.menuButton.contains(e.target as Node) &&
        !this.menuContainer.contains(e.target as Node)
      ) {
        this.closeMenu();
      }
    });

    // Close menu when video plays
    player.on("play", () => this.closeMenu());

    // Handle keyboard navigation
    this.menuContainer.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        if (this.currentSection !== "main") {
          this.showSection("main");
        } else {
          this.closeMenu();
        }
        e.preventDefault();
      }
    });

    player.on("loadedmetadata", () => {
      if (!player.videoWidth() && !player.videoHeight()) {
        if (player.error() !== null) return;
        const currentSrc = player.currentSrc();
        if (currentSrc === null) return;
        if (currentSrc.includes(".m3u8") || currentSrc.includes(".mpd")) {
          player.play();
        } else {
          player.error(MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED);
          return;
        }
      }
    });

    player.on("error", () => {
      const error = player.error();
      if (!error) return;

      if (
        error.code !== MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED &&
        error.code !== MediaError.MEDIA_ERR_DECODE
      )
        return;

      const currentSource = player.currentSource() as ISource;
      console.log(`Source '${currentSource.label}' is unsupported`);

      currentSource.errored = true;

      if (this.manuallySelected) {
        return;
      }

      if (
        this.selectedSourceIndex !== -1 &&
        this.selectedSourceIndex + 1 < this.sources.length
      ) {
        this.selectedSourceIndex += 1;
        const newSource = this.sources[this.selectedSourceIndex];
        console.log(`Trying next source in playlist: '${newSource.label}'`);

        const currentTime = player.currentTime();
        player.src(newSource);
        player.load();
        player.one("canplay", () => {
          player.currentTime(currentTime);
        });
        player.play();
        this.renderMenu();
      } else {
        console.log("No more sources in playlist");
      }
    });

    // Sync playback rate with player
    player.on("ratechange", () => {
      this.selectedSpeed = player.playbackRate();
      this.renderMenu();
    });
  }

  private createSettingsButton(): HTMLElement {
    const button = videojs.dom.createEl("button", {
      className: "vjs-settings-button vjs-control vjs-button",
      title: "Settings",
    }) as HTMLElement;

    button.innerHTML = `
      <span class="vjs-icon-placeholder" aria-hidden="true">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </span>
      <span class="vjs-control-text">Settings</span>
    `;

    button.addEventListener("click", (e) => {
      e.stopPropagation();
      this.toggleMenu();
    });

    return button;
  }

  private createMenuContainer(): HTMLElement {
    const container = videojs.dom.createEl("div", {
      className: "vjs-settings-menu-container",
    }) as HTMLElement;

    return container;
  }

  private toggleMenu() {
    if (this.isOpen) {
      this.closeMenu();
    } else {
      this.openMenu();
    }
  }

  private openMenu() {
    this.isOpen = true;
    this.currentSection = "main";
    this.menuContainer.classList.add("vjs-settings-menu-open");
    this.renderMenu();
  }

  private closeMenu() {
    this.isOpen = false;
    this.menuContainer.classList.remove("vjs-settings-menu-open");
  }

  private showSection(section: SettingsSection) {
    this.currentSection = section;
    this.renderMenu();
  }

  private renderMenu() {
    switch (this.currentSection) {
      case "main":
        this.renderMainMenu();
        break;
      case "quality":
        this.renderQualityMenu();
        break;
      case "speed":
        this.renderSpeedMenu();
        break;
      case "subtitles":
        this.renderSubtitlesMenu();
        break;
    }
  }

  private renderMainMenu() {
    const selectedSource = this.selectedSourceIndex >= 0 && this.selectedSourceIndex < this.sources.length
      ? this.sources[this.selectedSourceIndex]
      : null;
    const qualityLabel = selectedSource?.label || selectedSource?.type || "Auto";
    const speedLabel = this.selectedSpeed === 1 ? "1x" : `${this.selectedSpeed}x`;
    const subtitleLabel = this.getSelectedSubtitleLabel();

    this.menuContainer.innerHTML = `
      <div class="vjs-settings-menu">
        <div class="vjs-settings-menu-item" data-section="quality">
          <span class="vjs-settings-menu-icon">${ICONS.quality}</span>
          <span class="vjs-settings-menu-label">Quality</span>
          <span class="vjs-settings-menu-value">${qualityLabel}</span>
          <span class="vjs-settings-menu-chevron">${ICONS.chevron}</span>
        </div>
        <div class="vjs-settings-menu-item" data-section="speed">
          <span class="vjs-settings-menu-icon">${ICONS.speed}</span>
          <span class="vjs-settings-menu-label">Playback Speed</span>
          <span class="vjs-settings-menu-value">${speedLabel}</span>
          <span class="vjs-settings-menu-chevron">${ICONS.chevron}</span>
        </div>
        <div class="vjs-settings-menu-item" data-section="subtitles">
          <span class="vjs-settings-menu-icon">${ICONS.subtitles}</span>
          <span class="vjs-settings-menu-label">Subtitle/CC</span>
          <span class="vjs-settings-menu-value">${subtitleLabel}</span>
          <span class="vjs-settings-menu-chevron">${ICONS.chevron}</span>
        </div>
      </div>
    `;

    this.menuContainer.querySelectorAll(".vjs-settings-menu-item").forEach((item) => {
      item.addEventListener("click", (e) => {
        e.stopPropagation();
        const section = item.getAttribute("data-section") as SettingsSection;
        this.showSection(section);
      });
    });
  }

  private renderQualityMenu() {
    let items = "";
    
    if (this.sources.length === 0) {
      items = `
        <div class="vjs-settings-submenu-item vjs-disabled">
          <span class="vjs-settings-check"></span>
          <span class="vjs-settings-submenu-label">No sources available</span>
        </div>
      `;
    } else {
      items = this.sources
        .map((source, index) => {
          const isSelected = index === this.selectedSourceIndex;
          const isErrored = source.errored;
          const label = source.label || source.type || `Source ${index + 1}`;
          return `
            <div class="vjs-settings-submenu-item ${isSelected ? "vjs-selected" : ""} ${isErrored ? "vjs-errored" : ""}" data-index="${index}">
              <span class="vjs-settings-check">${isSelected ? "✓" : ""}</span>
              <span class="vjs-settings-submenu-label">${label}</span>
            </div>
          `;
        })
        .join("");
    }

    this.menuContainer.innerHTML = `
      <div class="vjs-settings-menu vjs-settings-submenu">
        <div class="vjs-settings-submenu-header" data-action="back">
          <span class="vjs-settings-back-icon">${ICONS.back}</span>
          <span class="vjs-settings-submenu-title">Quality</span>
        </div>
        <div class="vjs-settings-submenu-items">
          ${items}
        </div>
      </div>
    `;

    this.menuContainer.querySelector("[data-action='back']")?.addEventListener("click", (e) => {
      e.stopPropagation();
      this.showSection("main");
    });

    this.menuContainer.querySelectorAll(".vjs-settings-submenu-item:not(.vjs-disabled)").forEach((item) => {
      item.addEventListener("click", (e) => {
        e.stopPropagation();
        const index = parseInt(item.getAttribute("data-index") || "0", 10);
        this.selectSource(index);
      });
    });
  }

  private renderSpeedMenu() {
    const items = this.playbackRates
      .map((rate) => {
        const isSelected = rate === this.selectedSpeed;
        const label = rate === 1 ? "Normal" : `${rate}x`;
        return `
          <div class="vjs-settings-submenu-item ${isSelected ? "vjs-selected" : ""}" data-rate="${rate}">
            <span class="vjs-settings-check">${isSelected ? "✓" : ""}</span>
            <span class="vjs-settings-submenu-label">${label}</span>
          </div>
        `;
      })
      .join("");

    this.menuContainer.innerHTML = `
      <div class="vjs-settings-menu vjs-settings-submenu">
        <div class="vjs-settings-submenu-header" data-action="back">
          <span class="vjs-settings-back-icon">${ICONS.back}</span>
          <span class="vjs-settings-submenu-title">Playback Speed</span>
        </div>
        <div class="vjs-settings-submenu-items">
          ${items}
        </div>
      </div>
    `;

    this.menuContainer.querySelector("[data-action='back']")?.addEventListener("click", (e) => {
      e.stopPropagation();
      this.showSection("main");
    });

    this.menuContainer.querySelectorAll(".vjs-settings-submenu-item").forEach((item) => {
      item.addEventListener("click", (e) => {
        e.stopPropagation();
        const rate = parseFloat(item.getAttribute("data-rate") || "1");
        this.selectSpeed(rate);
      });
    });
  }

  private renderSubtitlesMenu() {
    const textTracks = this.player.textTracks();
    const items: string[] = [];

    // Add "Off" option
    let hasSelected = false;
    for (let i = 0; i < textTracks.length; i++) {
      const track = textTracks[i];
      if (track.kind === "captions" || track.kind === "subtitles") {
        if (track.mode === "showing") {
          hasSelected = true;
        }
      }
    }

    items.push(`
      <div class="vjs-settings-submenu-item ${!hasSelected ? "vjs-selected" : ""}" data-track-index="-1">
        <span class="vjs-settings-check">${!hasSelected ? "✓" : ""}</span>
        <span class="vjs-settings-submenu-label">subtitles off</span>
      </div>
    `);

    for (let i = 0; i < textTracks.length; i++) {
      const track = textTracks[i];
      if (track.kind === "captions" || track.kind === "subtitles") {
        const isSelected = track.mode === "showing";
        items.push(`
          <div class="vjs-settings-submenu-item ${isSelected ? "vjs-selected" : ""}" data-track-index="${i}">
            <span class="vjs-settings-check">${isSelected ? "✓" : ""}</span>
            <span class="vjs-settings-submenu-label">${track.label || track.language}</span>
          </div>
        `);
      }
    }

    this.menuContainer.innerHTML = `
      <div class="vjs-settings-menu vjs-settings-submenu">
        <div class="vjs-settings-submenu-header" data-action="back">
          <span class="vjs-settings-back-icon">${ICONS.back}</span>
          <span class="vjs-settings-submenu-title">Subtitle/CC</span>
        </div>
        <div class="vjs-settings-submenu-items">
          ${items.join("")}
        </div>
      </div>
    `;

    this.menuContainer.querySelector("[data-action='back']")?.addEventListener("click", (e) => {
      e.stopPropagation();
      this.showSection("main");
    });

    this.menuContainer.querySelectorAll(".vjs-settings-submenu-item").forEach((item) => {
      item.addEventListener("click", (e) => {
        e.stopPropagation();
        const index = parseInt(item.getAttribute("data-track-index") || "-1", 10);
        this.selectSubtitle(index);
      });
    });
  }

  private getSelectedSubtitleLabel(): string {
    const textTracks = this.player.textTracks();
    for (let i = 0; i < textTracks.length; i++) {
      const track = textTracks[i];
      if (
        (track.kind === "captions" || track.kind === "subtitles") &&
        track.mode === "showing"
      ) {
        return track.label || track.language || "On";
      }
    }
    return "subtitles off";
  }

  private selectSource(index: number) {
    if (index === this.selectedSourceIndex) {
      this.showSection("main");
      return;
    }

    this.manuallySelected = true;
    this.selectedSourceIndex = index;
    const source = this.sources[index];

    const currentTime = this.player.currentTime();
    const paused = this.player.paused();

    this.player.src(source);
    this.player.one("canplay", () => {
      if (paused) {
        this.player.pause();
      }
      this.player.currentTime(currentTime);
    });
    this.player.play();

    this.showSection("main");
  }

  private selectSpeed(rate: number) {
    this.selectedSpeed = rate;
    this.player.playbackRate(rate);
    this.showSection("main");
  }

  private selectSubtitle(trackIndex: number) {
    const textTracks = this.player.textTracks();

    // Disable all tracks first
    for (let i = 0; i < textTracks.length; i++) {
      const track = textTracks[i];
      if (track.kind === "captions" || track.kind === "subtitles") {
        track.mode = "disabled";
      }
    }

    // Enable selected track
    if (trackIndex >= 0 && trackIndex < textTracks.length) {
      textTracks[trackIndex].mode = "showing";
    }

    this.showSection("main");
  }

  // Public API
  setSources(sources: ISource[]) {
    const cleanupTracks = this.cleanupTextTracks.splice(0);
    for (const track of cleanupTracks) {
      this.player.removeRemoteTextTrack(track);
    }

    // Store the sources array
    this.sources = sources.slice(); // Make a copy to ensure we have our own reference
    
    if (this.sources.length !== 0) {
      this.selectedSourceIndex = 0;
      this.player.src(this.sources[0]);
    } else {
      this.selectedSourceIndex = -1;
    }

    this.manuallySelected = false;
  }

  get textTracks(): HTMLTrackElement[] {
    return [...this.cleanupTextTracks, ...this.manualTextTracks];
  }

  addTextTrack(options: videojs.TextTrackOptions, manualCleanup: boolean) {
    const track = this.player.addRemoteTextTrack(options, true);
    if (manualCleanup) {
      this.manualTextTracks.push(track);
    } else {
      this.cleanupTextTracks.push(track);
    }
    return track;
  }

  removeTextTrack(track: HTMLTrackElement) {
    this.player.removeRemoteTextTrack(track);
    let index = this.manualTextTracks.indexOf(track);
    if (index != -1) {
      this.manualTextTracks.splice(index, 1);
    }
    index = this.cleanupTextTracks.indexOf(track);
    if (index != -1) {
      this.cleanupTextTracks.splice(index, 1);
    }
  }
}

// Register the plugin with video.js.
videojs.registerPlugin("settingsMenu", SettingsMenuPlugin);

/* eslint-disable @typescript-eslint/naming-convention */
declare module "video.js" {
  interface VideoJsPlayer {
    settingsMenu: () => SettingsMenuPlugin;
  }
  interface VideoJsPlayerPluginOptions {
    settingsMenu?: {};
  }
}

export default SettingsMenuPlugin;


