// Gaps in video.js 8's bundled (JSDoc-generated) types.
import type Component from "video.js/dist/types/component";
import type ControlBar from "video.js/dist/types/control-bar/control-bar";
// videojs-mobile-ui ships its Player augmentation but doesn't reference it from package.json
import type {} from "videojs-mobile-ui/player";

declare module "video.js/dist/types/player" {
  export default interface Player {
    controlBar: ControlBar;
    // Getter overloads: the generated combined getter/setter signatures return `T | undefined`.
    currentTime(): number;
    duration(): number;
    volume(): number;
    muted(): boolean;
    isFullscreen(): boolean;
    poster(): string;
    playbackRate(): number;
    hasStarted(): boolean;
    // typed as returning Tech
    currentSource(): { src: string; type?: string };
  }
}

// Plugins get the evented mixin at runtime, but the class types omit it.
declare module "video.js/dist/types/plugin" {
  export default interface Plugin
    extends Pick<Component, "on" | "off" | "one" | "trigger"> {}
}
