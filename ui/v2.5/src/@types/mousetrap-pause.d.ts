/* eslint-disable @typescript-eslint/naming-convention */

declare module "mousetrap-pause" {
  import { MousetrapStatic } from "mousetrap";

  function MousetrapPause(mousetrap: MousetrapStatic): MousetrapStatic;

  export default MousetrapPause;

  module "mousetrap" {
    interface MousetrapStatic {
      pause(): void;
      unpause(): void;
      pauseCombo(combo: string): void;
      unpauseCombo(combo: string): void;
    }
    interface MousetrapInstance {
      pause(): void;
      unpause(): void;
      pauseCombo(combo: string): void;
      unpauseCombo(combo: string): void;
    }
  }
}
