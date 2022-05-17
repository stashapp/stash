import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export function galleryTitle(
  s: Partial<Pick<GQL.GalleryDataFragment, "title" | "files">>
) {
  if (s.title) {
    return s.title;
  }
  if (s.files && s.files.length > 0) {
    return TextUtils.fileNameFromPath(s.files[0].path);
  }
  return "";
}

export function galleryPath(s: Pick<GQL.GalleryDataFragment, "files">) {
  if (s.files && s.files.length > 0) {
    return s.files[0].path;
  }
  return "";
}
