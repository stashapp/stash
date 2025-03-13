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

function parseDate(input: string): string {
  return input
    .replace(/\./g, " ")
    .replace(/ {2,}/g, " ")
    .trim();
}

export function prepareQueryString(
  scene: Partial<GQL.SlimSceneDataFragment>,
  paths: string[],
  filename: string,
  mode: ParseMode,
  blacklist: string[]
) {
  let str = "";

  if ((mode === "auto" && scene.date && scene.studio) || mode === "metadata") {
    str = [
      scene.date,
      scene.studio?.name ?? "",
      (scene?.performers ?? []).map((p) => p.name).join(" "),
      scene?.title ? scene.title.replace(/[^a-zA-Z0-9 ]+/g, "") : "",
    ]
      .filter((s) => s !== "")
      .join(" ");
  } else {
    if (mode === "auto" || mode === "filename") {
      str = filename;
    } else if (mode === "path") {
      str = [...paths, filename].join(" ");
    } else if (mode === "dir" && paths.length) {
      str = paths[paths.length - 1];
    }
  }

  const regexReplacements = new Map<string, string>();
  const regexs = blacklist
    .map((entry) => {
      let [pattern, replacement] = entry.split("||"); // Extract regex and replacement

      try {
        const compiledRegex = new RegExp(pattern, "gi");

        if (replacement  !== undefined ){
          regexReplacements.set(compiledRegex.source, replacement); // Store replacement
        }

        return compiledRegex;
      } catch {
        return null; // Ignore invalid regex patterns
      }
    })
    .filter((r) => r !== null) as RegExp[];

  // Apply regex filtering and replacements
  regexs.forEach((regex) => {
    const replacement = regexReplacements.get(regex.source);
    if (replacement) {
      str = str.replace(regex, (match:string, ...groups: string[]) => {
        return replacement.replace(/\\(\d+)/g, (_, groupIndex) => groups[parseInt(groupIndex, 10) - 1] || "");
      });
    } else {
      str = str.replace(regex, "");
    }
  });
  str = handleSpecialStrings(str);
  str = parseDate(str);
  return str;
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
