import localForage from "localforage";
import * as GQL from "src/core/generated-graphql";

export const DEFAULT_SCENE_FILE_GAIN = 100;
export const MIN_SCENE_FILE_GAIN = 0;
export const MAX_SCENE_FILE_GAIN = 600;

const MD5_FINGERPRINT = "md5";
const OSHASH_FINGERPRINT = "oshash";
const STORAGE_KEY_PREFIX = "scene-file-gain:";

export type SceneFileGainMap = Record<string, number>;

type SceneFileGainFile = Pick<GQL.VideoFileDataFragment, "id" | "path"> & {
  fingerprints: Array<Pick<GQL.Fingerprint, "type" | "value">>;
};

export function clampSceneFileGain(value: number) {
  if (!Number.isFinite(value)) {
    return DEFAULT_SCENE_FILE_GAIN;
  }

  return Math.min(
    MAX_SCENE_FILE_GAIN,
    Math.max(MIN_SCENE_FILE_GAIN, Math.round(value))
  );
}

function getFingerprintValue(file: SceneFileGainFile, type: string) {
  return file.fingerprints.find((fingerprint) => fingerprint.type === type)
    ?.value;
}

export function getSceneFileGainStorageId(file: SceneFileGainFile) {
  const md5 = getFingerprintValue(file, MD5_FINGERPRINT);
  if (md5) {
    return `md5:${md5}`;
  }

  const oshash = getFingerprintValue(file, OSHASH_FINGERPRINT);
  if (oshash) {
    return `oshash:${oshash}`;
  }

  return `path:${file.path}`;
}

function getSceneFileGainStorageKey(file: SceneFileGainFile) {
  return `${STORAGE_KEY_PREFIX}${getSceneFileGainStorageId(file)}`;
}

export function normalizeSceneFileGains(
  files: SceneFileGainFile[],
  gains: SceneFileGainMap = {}
) {
  return Object.fromEntries(
    files.map((file) => [
      file.id,
      clampSceneFileGain(gains[file.id] ?? DEFAULT_SCENE_FILE_GAIN),
    ])
  );
}

export async function loadSceneFileGains(files: SceneFileGainFile[]) {
  const entries = await Promise.all(
    files.map(async (file) => {
      const storedGain = await localForage.getItem<number>(
        getSceneFileGainStorageKey(file)
      );

      return [
        file.id,
        clampSceneFileGain(storedGain ?? DEFAULT_SCENE_FILE_GAIN),
      ] as const;
    })
  );

  return Object.fromEntries(entries);
}

export async function saveSceneFileGains(
  files: SceneFileGainFile[],
  gains: SceneFileGainMap
) {
  const normalizedGains = normalizeSceneFileGains(files, gains);

  await Promise.all(
    files.map(async (file) => {
      const gain = normalizedGains[file.id];
      const storageKey = getSceneFileGainStorageKey(file);

      if (gain === DEFAULT_SCENE_FILE_GAIN) {
        await localForage.removeItem(storageKey);
        return;
      }

      await localForage.setItem(storageKey, gain);
    })
  );

  return normalizedGains;
}
