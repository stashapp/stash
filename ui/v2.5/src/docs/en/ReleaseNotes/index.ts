import v0170 from "./v0170.md";
import v0200 from "./v0200.md";
import v0240 from "./v0240.md";
import v0250 from "./v0250.md";
import v0260 from "./v0260.md";
import v0270 from "./v0270.md";
import v0290 from "./v0290.md";

export interface IReleaseNotes {
  // handle should be in the form of YYYYMMDD
  date: number;
  version: string;
  content: string;
}

export const releaseNotes: IReleaseNotes[] = [
  {
    date: 20251026,
    version: "v0.29.0",
    content: v0290,
  },
  {
    date: 20240826,
    version: "v0.27.0",
    content: v0270,
  },
  {
    date: 20240510,
    version: "v0.26.0",
    content: v0260,
  },
  {
    date: 20240228,
    version: "v0.25.0",
    content: v0250,
  },
  {
    date: 20231212,
    version: "v0.24.0",
    content: v0240,
  },
  {
    date: 20230301,
    version: "v0.20.0",
    content: v0200,
  },
  {
    date: 20220906,
    version: "v0.17.0",
    content: v0170,
  },
];
