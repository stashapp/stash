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

export const initialConfig: ITaggerConfig = {
  blacklist: DEFAULT_BLACKLIST,
  showMales: false,
  mode: "auto",
  setCoverImage: true,
  setTags: false,
  tagOperation: "merge",
  fingerprintQueue: {},
};

export type ParseMode = "auto" | "filename" | "dir" | "path" | "metadata";
export const ModeDesc = {
  auto: "Uses metadata if present, or filename",
  metadata: "Only uses metadata",
  filename: "Only uses filename",
  dir: "Only uses parent directory of video file",
  path: "Uses entire file path",
};

export interface ITaggerConfig {
  blacklist: string[];
  showMales: boolean;
  mode: ParseMode;
  setCoverImage: boolean;
  setTags: boolean;
  tagOperation: string;
  selectedEndpoint?: string;
  fingerprintQueue: Record<string, string[]>;
}
