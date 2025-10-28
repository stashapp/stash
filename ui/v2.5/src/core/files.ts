import TextUtils from "src/utils/text";
import * as GQL from "src/core/generated-graphql";

export interface IFile {
  path: string;
}

interface IObjectWithFiles {
  files?: GQL.Maybe<IFile[]>;
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

interface IObjectWithVisualFiles {
  visual_files?: IFile[];
}

export interface IObjectWithTitleVisualFiles extends IObjectWithVisualFiles {
  title?: GQL.Maybe<string>;
}

export function imageTitle(s: Partial<IObjectWithTitleVisualFiles>) {
  if (s.title) {
    return s.title;
  }
  if (s.visual_files && s.visual_files.length > 0) {
    return TextUtils.fileNameFromPath(s.visual_files[0].path);
  }
  return "";
}

export function imagePath(s: IObjectWithVisualFiles) {
  if (s.visual_files && s.visual_files.length > 0) {
    return s.visual_files[0].path;
  }
  return "";
}
