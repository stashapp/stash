declare module "intersection-observer";

declare module "*.md" {
  const src: string;
  export default src;
}

// XXbiome-ignore @typescript-eslint/naming-convention: intentional
interface ImportMetaEnv {
  readonly VITE_APP_GITHASH?: string;
  readonly VITE_APP_STASH_VERSION?: string;
  readonly VITE_APP_DATE?: string;
  readonly VITE_APP_PLATFORM_URL?: string;
}
