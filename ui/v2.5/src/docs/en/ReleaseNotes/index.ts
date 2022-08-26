import v0170 from "./v0170.md";
import r20220826 from "./20220826.md";

export type Module = typeof v0170;

interface IReleaseNotes {
  // handle should be in the form of YYYYMMDD
  date: number;
  content: Module;
}

export const releaseNotes: IReleaseNotes[] = [
  {
    date: 20220801,
    content: v0170,
  },
  {
    date: 20220826,
    content: r20220826,
  },
];
