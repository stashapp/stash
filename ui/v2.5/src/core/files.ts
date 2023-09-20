import TextUtils from "src/utils/text";
import * as GQL from "src/core/generated-graphql";

export interface IFile {
  path: string;
}

interface IObjectWithFiles {
  files?: IFile[];
}

export interface IObjectWithTitleFiles extends IObjectWithFiles {
  title?: GQL.Maybe<string>;
}

export function objectTitle(s: Partial<IObjectWithTitleFiles>) {
  if (s.title) {
    return s.title;
  }
  if (s.files && s.files.length > 0) {
    return TextUtils.fileNameFromPath(s.files[0].path);
  }
  return "";
}

export function objectPath(s: IObjectWithFiles) {
  if (s.files && s.files.length > 0) {
    return s.files[0].path;
  }
  return "";
}
