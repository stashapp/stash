import * as GQL from "src/core/generated-graphql";

export const parsePath = (filePath: string) => {
  const path = filePath.toLowerCase();
  const isWin = /^([a-z]:|\\\\)/.test(path);
  const normalizedPath = isWin
    ? path.replace(/^[a-z]:/, "").replace(/\\/g, "/")
    : path;
  const pathComponents = normalizedPath
    .split("/")
    .filter((component) => component.trim().length > 0);
  const fileName = pathComponents[pathComponents.length - 1];

  const ext = fileName.match(/\.[a-z0-9]*$/)?.[0] ?? "";
  const file = fileName.slice(0, ext.length * -1);
  const paths =
    pathComponents.length > 2
      ? pathComponents.slice(0, pathComponents.length - 2)
      : [];

  return { paths, file, ext };
};

export interface IStashBoxFingerprint {
  hash: string;
  algorithm: string;
  duration: number;
}

export interface IStashBoxPerformer {
  id: string;
  name: string;
  gender?: GQL.GenderEnum;
  url?: string;
  twitter?: string;
  instagram?: string;
  birthdate?: string;
  ethnicity?: string;
  country?: string;
  eye_color?: string;
  height?: string;
  measurements?: string;
  fake_tits?: string;
  career_length?: string;
  tattoos?: string;
  piercings?: string;
  aliases?: string;
  images: string[];
}

export interface IStashBoxTag {
  id: string;
  name: string;
}

export interface IStashBoxStudio {
  id: string;
  name: string;
  url?: string;
  image?: string;
}

export interface IStashBoxScene {
  id: string;
  title: string;
  date: string;
  duration: number;
  details?: string;
  url?: string;

  studio: IStashBoxStudio;
  images: string[];
  tags: IStashBoxTag[];
  performers: IStashBoxPerformer[];
  fingerprints: IStashBoxFingerprint[];
}

const selectStudio = (
  studio?: GQL.ScrapedSceneStudio | null
): IStashBoxStudio => ({
  id: studio?.stored_id ?? "",
  name: studio?.name ?? "",
  url: studio?.url ?? undefined,
});

const selectFingerprints = (scene: GQL.ScrapedScene|null): IStashBoxFingerprint[] => (
  scene?.fingerprints ?? []
);

const selectTags = (tags: GQL.ScrapedSceneTag[]): IStashBoxTag[] =>
  tags.map((t) => ({
    id: t.stored_id ?? "",
    name: t.name ?? "",
  }));

const selectPerformers = (
  performers: GQL.ScrapedScenePerformer[]
): IStashBoxPerformer[] =>
  performers.map((p) => ({
    id: p.site_id ?? "",
    name: p.name ?? "",
    gender: (p.gender ?? GQL.GenderEnum.Female) as GQL.GenderEnum,
    url: p.url ?? undefined,
    twitter: p.twitter ?? undefined,
    instagram: p.instagram ?? undefined,
    birthdate: p.birthdate ?? undefined,
    ethnicity: p.ethnicity ?? undefined,
    country: p.country ?? undefined,
    eye_color: p.eye_color ?? undefined,
    height: p.height ?? undefined,
    measurements: p.measurements ?? undefined,
    fake_tits: p.fake_tits ?? undefined,
    career_length: p.career_length ?? undefined,
    tattoos: p.tattoos ?? undefined,
    piercings: p.piercings ?? undefined,
    aliases: p.aliases ?? undefined,
    images: p.images ?? [],
  }));

export const selectScenes = (
  scenes?: (GQL.ScrapedScene | null)[]
): IStashBoxScene[] => {
  const result = (scenes ?? [])
    .filter((s) => s !== null)
    .map(
      (s) =>
        ({
          id: s?.site_id ?? "",
          title: s?.title ?? "",
          date: s?.date ?? "",
          duration: s?.duration ?? 0,
          details: s?.details,
          url: s?.url,
          images: s?.image ? [s.image] : [],
          studio: selectStudio(s?.studio),
          fingerprints: selectFingerprints(s),
          performers: selectPerformers(s?.performers ?? []),
          tags: selectTags(s?.tags ?? []),
        } as IStashBoxScene)
    );

  return result;
};
