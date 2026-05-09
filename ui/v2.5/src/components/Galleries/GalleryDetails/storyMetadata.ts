import * as GQL from "src/core/generated-graphql";
import { CustomFieldMap } from "src/components/Shared/CustomFields";

export const storyTextExtensions = new Set([
  "txt",
  "text",
  "md",
  "markdown",
  "nfo",
  "rst",
  "log",
]);

export const storyCustomFieldKeys = {
  author: "story_author",
  language: "story_language",
  sourceWebsite: "story_source_website",
  tagLine: "story_tag_line",
  audioUrl: "story_audio_url",
  backCoverUrl: "story_back_cover_url",
} as const;

export interface StoryMetadata {
  author: string;
  language: string;
  sourceWebsite: string;
  sourceUrl: string;
  tagLine: string;
  audioUrl: string;
  backCoverUrl: string;
}

function fileExtension(path?: string) {
  if (!path) {
    return "";
  }

  const match = path.toLowerCase().match(/\.([^.]+)$/);
  return match?.[1] ?? "";
}

function getStringValue(value: unknown) {
  if (typeof value === "string") {
    return value;
  }

  if (typeof value === "number" || typeof value === "boolean") {
    return `${value}`;
  }

  return "";
}

function setOrDelete(
  values: CustomFieldMap,
  key: string,
  value: string
): CustomFieldMap {
  const nextValues = { ...values };
  const trimmedValue = value.trim();

  if (trimmedValue) {
    nextValues[key] = trimmedValue;
  } else {
    delete nextValues[key];
  }

  return nextValues;
}

export function isStoryGallery(
  gallery: Pick<Partial<GQL.GalleryDataFragment>, "files" | "paths">
) {
  if (!gallery.paths?.text || gallery.files?.length !== 1) {
    return false;
  }

  return storyTextExtensions.has(fileExtension(gallery.files[0]?.path));
}

export function getStoryMetadata(
  gallery: Pick<Partial<GQL.GalleryDataFragment>, "custom_fields" | "urls">
): StoryMetadata {
  const customFields = gallery.custom_fields ?? {};

  return {
    author: getStringValue(customFields[storyCustomFieldKeys.author]),
    language: getStringValue(customFields[storyCustomFieldKeys.language]),
    sourceWebsite: getStringValue(
      customFields[storyCustomFieldKeys.sourceWebsite]
    ),
    sourceUrl: gallery.urls?.[0] ?? "",
    tagLine: getStringValue(customFields[storyCustomFieldKeys.tagLine]),
    audioUrl: getStringValue(customFields[storyCustomFieldKeys.audioUrl]),
    backCoverUrl: getStringValue(
      customFields[storyCustomFieldKeys.backCoverUrl]
    ),
  };
}

export function withoutStoryCustomFields(values: CustomFieldMap = {}) {
  const nextValues = { ...values };

  Object.values(storyCustomFieldKeys).forEach((key) => {
    delete nextValues[key];
  });

  return nextValues;
}

export function applyStoryMetadata(
  values: CustomFieldMap = {},
  metadata: Omit<StoryMetadata, "sourceUrl">
) {
  let nextValues = { ...values };

  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.author,
    metadata.author
  );
  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.language,
    metadata.language
  );
  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.sourceWebsite,
    metadata.sourceWebsite
  );
  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.tagLine,
    metadata.tagLine
  );
  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.audioUrl,
    metadata.audioUrl
  );
  nextValues = setOrDelete(
    nextValues,
    storyCustomFieldKeys.backCoverUrl,
    metadata.backCoverUrl
  );

  return nextValues;
}
