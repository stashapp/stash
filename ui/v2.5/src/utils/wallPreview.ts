export type PreviewMediaType = "image" | "video";

export interface PreviewSource {
  src?: string | null;
  mediaType: PreviewMediaType;
}

export interface SelectedPreviewSource {
  src: string;
  mediaType: PreviewMediaType;
}

export function getFirstValidPreviewSource(
  srcSet: readonly PreviewSource[],
  invalidSrcSet: string[]
): SelectedPreviewSource {
  const validSrcSet = srcSet.filter((s) => s.src);

  if (!validSrcSet.length) {
    return { src: "", mediaType: "image" };
  }

  const selected =
    validSrcSet.find(({ src }) => !invalidSrcSet.includes(src!)) ??
    ([...validSrcSet].pop() as PreviewSource);

  return {
    src: selected.src!,
    mediaType: selected.mediaType,
  };
}
