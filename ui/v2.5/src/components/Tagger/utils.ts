import * as GQL from "src/core/generated-graphql";
import { ParseMode } from "./constants";
import { queryFindStudio } from "src/core/StashService";
import { mergeStashIDs } from "src/utils/stashbox";

const months = [
  "jan",
  "feb",
  "mar",
  "apr",
  "may",
  "jun",
  "jul",
  "aug",
  "sep",
  "oct",
  "nov",
  "dec",
];

const ddmmyyRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
const yyyymmddRegex = /(\d{4})[-.](\d{2})[-.](\d{2})/;
const mmddyyRegex = /(\d{2})[-.](\d{2})[-.](\d{4})/;
const ddMMyyRegex = new RegExp(
  `(\\d{1,2}).(${months.join("|")})\\.?.(\\d{4})`,
  "i"
);
const MMddyyRegex = new RegExp(
  `(${months.join("|")})\\.?.(\\d{1,2}),?.(\\d{4})`,
  "i"
);
const javcodeRegex = /([a-zA-Z|tT28|tT38]+-\d+[zZeE]?)/;

const handleSpecialStrings = (input: string): string => {
  let output = input;
  const ddmmyy = output.match(ddmmyyRegex);
  if (ddmmyy) {
    output = output.replace(
      ddmmyy[0],
      ` 20${ddmmyy[1]}-${ddmmyy[2]}-${ddmmyy[3]} `
    );
  }
  const mmddyy = output.match(mmddyyRegex);
  if (mmddyy) {
    output = output.replace(
      mmddyy[0],
      ` ${mmddyy[1]}-${mmddyy[2]}-${mmddyy[3]} `
    );
  }
  const ddMMyy = output.match(ddMMyyRegex);
  if (ddMMyy) {
    const month = (months.indexOf(ddMMyy[2].toLowerCase()) + 1)
      .toString()
      .padStart(2, "0");
    output = output.replace(
      ddMMyy[0],
      ` ${ddMMyy[3]}-${month}-${ddMMyy[1].padStart(2, "0")} `
    );
  }
  const MMddyy = output.match(MMddyyRegex);
  if (MMddyy) {
    const month = (months.indexOf(MMddyy[1].toLowerCase()) + 1)
      .toString()
      .padStart(2, "0");
    output = output.replace(
      MMddyy[0],
      ` ${MMddyy[3]}-${month}-${MMddyy[2].padStart(2, "0")} `
    );
  }

  const yyyymmdd = output.search(yyyymmddRegex);
  // if we find a date, then replace hyphens with spaces outside of the date
  // replace dots with hyphens in the date
  if (yyyymmdd !== -1)
    return (
      output.slice(0, yyyymmdd).replace(/-/g, " ") +
      output.slice(yyyymmdd, yyyymmdd + 10).replace(/\./g, "-") +
      output.slice(yyyymmdd + 10).replace(/-/g, " ")
    );

  const javcodeIndex = output.search(javcodeRegex);
  // if we find a javcode, then replace hyphens with spaces outside of the javcode
  if (javcodeIndex !== -1) {
    const javcodeLength = output.match(javcodeRegex)![1].length;
    return (
      output.slice(0, javcodeIndex).replace(/-/g, " ") +
      output.slice(javcodeIndex, javcodeIndex + javcodeLength) +
      output.slice(javcodeIndex + javcodeLength).replace(/-/g, " ")
    );
  }
  // otherwise just replace hyphens with spaces
  return output.replace(/-/g, " ");
};

export function prepareQueryString(
  scene: Partial<GQL.SlimSceneDataFragment>,
  paths: string[],
  filename: string,
  mode: ParseMode,
  blacklist: string[]
) {
  const regexs = blacklist
    .map((b) => {
      try {
        return new RegExp(b, "gi");
      } catch {
        // ignore
        return null;
      }
    })
    .filter((r) => r !== null) as RegExp[];

  if ((mode === "auto" && scene.date && scene.studio) || mode === "metadata") {
    let str = [
      scene.date,
      scene.studio?.name ?? "",
      (scene?.performers ?? []).map((p) => p.name).join(" "),
      scene?.title ? scene.title.replace(/[^a-zA-Z0-9 ]+/g, "") : "",
    ]
      .filter((s) => s !== "")
      .join(" ");
    regexs.forEach((re) => {
      str = str.replace(re, " ");
    });
    return str;
  }
  let s = "";

  if (mode === "auto" || mode === "filename") {
    s = filename;
  } else if (mode === "path") {
    s = [...paths, filename].join(" ");
  } else if (mode === "dir" && paths.length) {
    s = paths[paths.length - 1];
  }

  regexs.forEach((re) => {
    s = s.replace(re, " ");
  });
  s = handleSpecialStrings(s);
  return s.replace(/\./g, " ").replace(/ +/g, " ");
}

export const parsePath = (filePath: string) => {
  if (!filePath) {
    return {
      paths: [],
      file: "",
      ext: "",
    };
  }

  const path = filePath.toLowerCase();
  // Absolute paths on Windows start with a drive letter (e.g. C:\)
  // Alternatively, they may start with a UNC path (e.g. \\server\share)
  // Remove the drive letter/UNC and replace backslashes with forward slashes
  const normalizedPath = path.replace(/^[a-z]:|\\\\/, "").replace(/\\/g, "/");
  const pathComponents = normalizedPath
    .split("/")
    .filter((component) => component.trim().length > 0);
  const fileName = pathComponents[pathComponents.length - 1];

  const ext = fileName.match(/\.[a-z0-9]*$/)?.[0] ?? "";
  const file = fileName.slice(0, ext.length * -1);

  // remove any .. or . paths
  const paths = (
    pathComponents.length >= 1
      ? pathComponents.slice(0, pathComponents.length - 1)
      : []
  ).filter((p) => p !== ".." && p !== ".");

  return { paths, file, ext };
};

export async function mergeStudioStashIDs(
  id: string,
  newStashIDs: GQL.StashIdInput[]
) {
  const existing = await queryFindStudio(id);
  if (existing?.data?.findStudio?.stash_ids) {
    return mergeStashIDs(existing.data.findStudio.stash_ids, newStashIDs);
  }

  return newStashIDs;
}
