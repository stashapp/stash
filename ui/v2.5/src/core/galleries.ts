import TextUtils from "src/utils/text";
import * as GQL from "src/core/generated-graphql";

interface IFile {
  path: string;
}

interface IGallery {
  files: GQL.Maybe<IFile[]>;
  folder?: GQL.Maybe<IFile>;
}

interface IGalleryWithTitle extends IGallery {
  title: GQL.Maybe<string>;
}

export function galleryTitle(s: Partial<IGalleryWithTitle>) {
  if (s.title) {
    return s.title;
  }
  if (s.files && s.files.length > 0) {
    return TextUtils.fileNameFromPath(s.files[0].path);
  }
  if (s.folder) {
    return TextUtils.fileNameFromPath(s.folder.path);
  }
  return "";
}

export function galleryPath(s: IGallery) {
  if (s.files && s.files.length > 0) {
    return s.files[0].path;
  }
  if (s.folder) {
    return s.folder.path;
  }
  return "";
}
