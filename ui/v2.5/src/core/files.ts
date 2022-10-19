import { TextUtils } from "src/utils";
import * as GQL from "src/core/generated-graphql";

interface IFile {
  path: string;
}

interface IObjectWithFiles {
  files: IFile[];
}

interface IObjectWithTitleFiles extends IObjectWithFiles {
  title: GQL.Maybe<string>;
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
