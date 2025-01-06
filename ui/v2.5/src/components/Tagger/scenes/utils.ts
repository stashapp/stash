import { SlimSceneDataFragment } from "src/core/generated-graphql";
import { IScrapedScene } from "../context";
import { distance } from "src/utils/hamming";

export function minDistance(hash: string, stashScene: SlimSceneDataFragment) {
  let ret = 9999;
  stashScene.files.forEach((cv) => {
    if (ret === 0) return;

    const stashHash = cv.fingerprints.find((fp) => fp.type === "phash");
    if (!stashHash) {
      return;
    }

    const d = distance(hash, stashHash.value);
    if (d < ret) {
      ret = d;
    }
  });

  return ret;
}

export function calculatePhashComparisonScore(
  stashScene: SlimSceneDataFragment,
  scrapedScene: IScrapedScene
) {
  const phashFingerprints =
    scrapedScene.fingerprints?.filter((f) => f.algorithm === "PHASH") ?? [];
  const filteredFingerprints = phashFingerprints.filter(
    (f) => minDistance(f.hash, stashScene) <= 8
  );

  if (phashFingerprints.length == 0) return [0, 0];

  return [
    filteredFingerprints.length,
    filteredFingerprints.length / phashFingerprints.length,
  ];
}

export function minDurationDiff(
  stashScene: SlimSceneDataFragment,
  duration: number
) {
  let ret = 9999;
  stashScene.files.forEach((cv) => {
    if (ret === 0) return;

    const d = Math.abs(duration - cv.duration);
    if (d < ret) {
      ret = d;
    }
  });

  return ret;
}

export function calculateDurationComparisonScore(
  stashScene: SlimSceneDataFragment,
  scrapedScene: IScrapedScene
) {
  if (scrapedScene.fingerprints && scrapedScene.fingerprints.length > 0) {
    const durations = scrapedScene.fingerprints.map((f) => f.duration);
    const diffs = durations.map((d) => minDurationDiff(stashScene, d));
    const filteredDurations = diffs.filter((duration) => duration <= 5);

    const minDiff = Math.min(...diffs);

    return [
      filteredDurations.length,
      filteredDurations.length / durations.length,
      minDiff,
    ];
  }
  return [0, 0, 0];
}

export function compareScenesForSort(
  stashScene: SlimSceneDataFragment,
  sceneA: IScrapedScene,
  sceneB: IScrapedScene
) {
  // Compare sceneA and sceneB to each other for sorting based on similarity to stashScene
  // Order of priority is: nb. phash match > nb. duration match > ratio duration match > ratio phash match

  // scenes without any fingerprints should be sorted to the end
  if (!sceneA.fingerprints?.length && sceneB.fingerprints?.length) {
    return 1;
  }
  if (!sceneB.fingerprints?.length && sceneA.fingerprints?.length) {
    return -1;
  }

  const [nbPhashMatchSceneA, ratioPhashMatchSceneA] =
    calculatePhashComparisonScore(stashScene, sceneA);
  const [nbPhashMatchSceneB, ratioPhashMatchSceneB] =
    calculatePhashComparisonScore(stashScene, sceneB);

  // If only one scene has matching phash, prefer that scene
  if (
    (nbPhashMatchSceneA != nbPhashMatchSceneB && nbPhashMatchSceneA === 0) ||
    nbPhashMatchSceneB === 0
  ) {
    return nbPhashMatchSceneB - nbPhashMatchSceneA;
  }

  // Prefer scene with highest ratio of phash matches
  if (ratioPhashMatchSceneA !== ratioPhashMatchSceneB) {
    return ratioPhashMatchSceneB - ratioPhashMatchSceneA;
  }

  // Same ratio of phash matches, check duration
  const [
    nbDurationMatchSceneA,
    ratioDurationMatchSceneA,
    minDurationDiffSceneA,
  ] = calculateDurationComparisonScore(stashScene, sceneA);
  const [
    nbDurationMatchSceneB,
    ratioDurationMatchSceneB,
    minDurationDiffSceneB,
  ] = calculateDurationComparisonScore(stashScene, sceneB);

  if (nbDurationMatchSceneA != nbDurationMatchSceneB) {
    return nbDurationMatchSceneB - nbDurationMatchSceneA;
  }

  // Same number of phash & duration, check duration ratio
  if (ratioDurationMatchSceneA != ratioDurationMatchSceneB) {
    return ratioDurationMatchSceneB - ratioDurationMatchSceneA;
  }

  // fall back to duration difference - less is better
  return minDurationDiffSceneA - minDurationDiffSceneB;
}
