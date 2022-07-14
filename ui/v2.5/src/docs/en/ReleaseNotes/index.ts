import v0170 from "./v0170.md";

export type Module = typeof v0170;

interface IReleaseNotes {
  // handle should be in the form of YYYYMMDD
  date: number;
  content: Module;
}

export const releaseNotes: IReleaseNotes[] = [
  {
    date: 20220715,
    content: v0170,
  },
];
