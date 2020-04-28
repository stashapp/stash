import * as GQL from "src/core/generated-graphql";

import {
  SearchScene_searchScene_performers_performer_urls as URL,
  SearchScene_searchScene_performers_performer_images as Image,
  SearchScene_searchScene_performers_performer_measurements as Measurements,
  SearchScene_searchScene_performers_performer_tattoos as BodyModification,
} from "src/definitions-box/SearchScene";
import {
  BreastTypeEnum,
  GenderEnum as StashGenderEnum,
} from "src/definitions-box/globalTypes";
import { GenderEnum } from "src/core/generated-graphql";

export const sortImageURLs = (
  images: Image[],
  orientation: "portrait" | "landscape"
) =>
  images
    .map((i) => ({
      id: i.id,
      url: i.url,
      width: i.width ?? 1,
      height: i.height ?? 1,
      aspect:
        orientation === "portrait"
          ? (i.height ?? 1) / (i.width ?? 1) > 1
          : (i.width ?? 1) / (i.height ?? 1) > 1,
    }))
    .sort((a, b) => {
      if (a.aspect > b.aspect) return -1;
      if (a.aspect < b.aspect) return 1;
      if (orientation === "portrait" && a.height > b.height) return -1;
      if (orientation === "portrait" && a.height < b.height) return 1;
      if (orientation === "landscape" && a.width > b.width) return -1;
      if (orientation === "landscape" && a.width < b.width) return 1;
      return 0;
    });

export const getImage = (
  images: Image[],
  orientation: "portrait" | "landscape"
) => sortImageURLs(images, orientation)?.[0]?.url ?? "";

export const getUrlByType = (urls: (URL | null)[], type: string) =>
  (urls && (urls.find((url) => url?.type === type) || {}).url) || "";

export const formatMeasurements = (measurements: Measurements) =>
  measurements.cup_size && measurements.waist && measurements.hip
    ? `${measurements.band_size}${measurements.cup_size}-${measurements.waist}-${measurements.hip}`
    : "";

export const formatBreastType = (type: BreastTypeEnum | null) =>
  type === BreastTypeEnum.FAKE
    ? "Yes"
    : type === BreastTypeEnum.NATURAL
    ? "No"
    : "";

export const formatGender = (type: StashGenderEnum | null) =>
  type === StashGenderEnum.FEMALE
    ? GenderEnum.Female
    : type === StashGenderEnum.MALE
    ? GenderEnum.Male
    : type === StashGenderEnum.TRANSGENDER_FEMALE
    ? GenderEnum.TransgenderFemale
    : type === StashGenderEnum.TRANSGENDER_MALE
    ? GenderEnum.TransgenderMale
    : type === StashGenderEnum.INTERSEX
    ? GenderEnum.Intersex
    : null;

export const formatCareerLength = (
  start?: number | null,
  end?: number | null
) => (start && end ? `${start} - ${end}` : start ? `${start} - ` : null);

export const formatBodyModification = (mods: BodyModification[] | null) =>
  (mods || [])
    .map((m) =>
      m.location && m.description
        ? `${m.location}, ${m.description}`
        : m.description || m.location
    )
    .join("; ");

export const formatURL = (urls: URL[], type: string) =>
  urls.find((u) => u.type === type)?.url ?? null;

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

const selectStudio = (studio?: GQL.ScrapedSceneStudio|null): IStashBoxStudio => ({
  id: studio?.stored_id ?? '',
  name: studio?.name ?? '',
  url: studio?.url ?? undefined
});

const selectFingerprints = (): IStashBoxFingerprint[] => [];

const selectTags = (tags: GQL.ScrapedSceneTag[]): IStashBoxTag[] => (
  tags.map(t => ({
    id: t.stored_id ?? '',
    name: t.name ?? '',
  }))
);

const selectPerformers = (performers: GQL.ScrapedScenePerformer[]): IStashBoxPerformer[] => (
  performers.map(p => ({
    id: p.stored_id ?? '',
    name: p.name ?? '',
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
    images: []
  }))
);

export const selectScenes = (scenes: (GQL.ScrapedScene | null)[]): IStashBoxScene[] => {
  const result = scenes.filter(s => s !== null).map(s => ({
    id: "ID",
    title: s?.title ?? "",
    date: s?.date ?? "",
    duration: 0,
    details: s?.details,
    url: s?.url,
    images: s?.image ? [s.image]  : [],
    studio: selectStudio(s?.studio),
    fingerprints: selectFingerprints(),
    performers: selectPerformers(s?.performers ?? []),
    tags: selectTags(s?.tags ?? []),
  } as IStashBoxScene));

  return result;
}
