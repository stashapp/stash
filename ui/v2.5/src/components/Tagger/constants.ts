import { ScraperSourceInput } from "src/core/generated-graphql";

export const STASH_BOX_PREFIX = "stashbox:";
export const SCRAPER_PREFIX = "scraper:";

export interface ITaggerSource {
  id: string;
  stashboxEndpoint?: string;
  sourceInput: ScraperSourceInput;
  displayName: string;
  supportSceneQuery?: boolean;
  supportSceneFragment?: boolean;
}

export const LOCAL_FORAGE_KEY = "tagger";
export const DEFAULT_BLACKLIST = [
  "\\sXXX\\s",
  "1080p",
  "720p",
  "2160p",
  "KTR",
  "RARBG",
  "\\scom\\s",
  "\\[",
  "\\]",
];
export const DEFAULT_EXCLUDED_PERFORMER_FIELDS = ["name"];
export const DEFAULT_EXCLUDED_STUDIO_FIELDS = ["name"];

export const initialConfig: ITaggerConfig = {
  blacklist: DEFAULT_BLACKLIST,
  showMales: false,
  mode: "auto",
  setCoverImage: true,
  setTags: false,
  tagOperation: "merge",
  fingerprintQueue: {},
  excludedPerformerFields: DEFAULT_EXCLUDED_PERFORMER_FIELDS,
  excludedStudioFields: DEFAULT_EXCLUDED_STUDIO_FIELDS,
  createParentStudios: true,
};

export type ParseMode = "auto" | "filename" | "dir" | "path" | "metadata";
export type TagOperation = "merge" | "overwrite";
export interface ITaggerConfig {
  blacklist: string[];
  showMales: boolean;
  mode: ParseMode;
  setCoverImage: boolean;
  setTags: boolean;
  tagOperation: TagOperation;
  selectedEndpoint?: string;
  fingerprintQueue: Record<string, string[]>;
  excludedPerformerFields?: string[];
  excludedStudioFields?: string[];
  createParentStudios: boolean;
}

export const PERFORMER_FIELDS = [
  "name",
  "aliases",
  "image",
  "gender",
  "birthdate",
  "ethnicity",
  "country",
  "eye_color",
  "hair_color",
  "height",
  "measurements",
  "fake_tits",
  "career_length",
  "tattoos",
  "piercings",
  "url",
  "twitter",
  "instagram",
  "details",
  "death_date",
  "weight",
];

export const STUDIO_FIELDS = ["name", "image", "url", "parent"];
