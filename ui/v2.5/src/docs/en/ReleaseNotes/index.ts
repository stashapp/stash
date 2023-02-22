import v0170 from "./v0170.md";

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
];
