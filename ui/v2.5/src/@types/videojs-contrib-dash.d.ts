/* eslint-disable @typescript-eslint/naming-convention */

declare module "videojs-contrib-dash" {
  class Html5DashJS {
    /**
     * Get a list of hooks for a specific lifecycle.
     *
     * @param type the lifecycle to get hooks from
     * @param hook optionally add a hook to the lifecycle
     * @return an array of hooks or empty if none
     */
    static hooks(type: string, hook: Function | Function[]): Function[];

    /**
     * Add a function hook to a specific dash lifecycle.
     *
     * @param type the lifecycle to hook the function to
     * @param hook the function or array of functions to attach
     */
    static hook(type: string, hook: Function | Function[]): void;

    /**
     * Remove a hook from a specific dash lifecycle.
     *
     * @param type the lifecycle that the function hooked to
     * @param hook the hooked function to remove
     * @return true if the function was removed, false if not found
     */
    static removeHook(type: string, hook: Function): boolean;
  }
}
