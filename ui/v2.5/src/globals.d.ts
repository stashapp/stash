// eslint-disable-next-line no-var
declare var STASH_BASE_URL: string;
declare module "intersection-observer";

declare module "*.md" {
  const src: string;
  export default src;
}

/* eslint-disable-next-line  @typescript-eslint/naming-convention */
interface ImportMetaEnv {
  readonly VITE_APP_GITHASH?: string;
  readonly VITE_APP_STASH_VERSION?: string;
  readonly VITE_APP_DATE?: string;
  readonly VITE_APP_PLATFORM_PORT?: string;
  readonly VITE_APP_HTTPS?: string;
}
