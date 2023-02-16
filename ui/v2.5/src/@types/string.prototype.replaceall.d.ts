declare module "string.prototype.replaceall" {
  function replaceAll(
    searchValue: string | RegExp,
    replaceValue: string
  ): string;
  function replaceAll(
    searchValue: string | RegExp,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    replacer: (substring: string, ...args: any[]) => string
  ): string;

  namespace replaceAll {
    function getPolyfill(): typeof replaceAll;
    function implementation(): typeof replaceAll;
    function shim(): void;
  }

  export default replaceAll;
}
