import v0170 from "./v0170.md";
import v0200 from "./v0200.md";

interface IReleaseNotes {
  // handle should be in the form of YYYYMMDD
  date: number;
  content: string;
}

export const releaseNotes: IReleaseNotes[] = [
  {
    date: 20220906,
    content: v0170,
  },
  {
    date: 20230224,
    content: v0200,
  }
];
