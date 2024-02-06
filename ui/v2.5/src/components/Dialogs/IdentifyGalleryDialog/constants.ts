import * as GQL from "src/core/generated-graphql";

export interface IScraperSource {
  id: string;
  displayName: string;
  scraper_id?: string;
  options?: GQL.IdentifyGalleryMetadataOptionsInput;
}

export const galleryFields = [
  "title",
  "code",
  "date",
  "photographer",
  "details",
  "url",
  "studio",
  "performers",
  "tags",
] as const;
export type GalleryField = (typeof galleryFields)[number];

export const multiValueGalleryFields: GalleryField[] = [
  "studio",
  "performers",
  "tags",
];

export function galleryFieldMessageID(field: GalleryField) {
  if (field === "code") {
    return "scene_code";
  } else if (field === "studio") {
    return "studio_and_parent";
  }

  return field;
}
