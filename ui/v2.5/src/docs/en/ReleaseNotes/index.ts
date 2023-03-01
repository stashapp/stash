import v0170 from "./v0170.md";
import v0200 from "./v0200.md";

export interface IReleaseNotes {
  // handle should be in the form of YYYYMMDD
  date: number;
  version: string;
  content: string;
}

export const releaseNotes: IReleaseNotes[] = [
  {
    date: 20230224,
    version: "v0.20.0",
    content: v0200,
  },
  {
    date: 20220906,
    version: "v0.17.0",
    content: v0170,
  },
];
