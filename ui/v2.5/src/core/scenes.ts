import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export function sceneTitle(
  s: Partial<Pick<GQL.SceneDataFragment, "title" | "files">>
) {
  if (s.title) {
    return s.title;
  }
  if (s.files && s.files.length > 0) {
    return TextUtils.fileNameFromPath(s.files[0].path);
  }
  return "";
}

export function scenePath(s: Pick<GQL.SceneDataFragment, "files">) {
  if (s.files && s.files.length > 0) {
    return s.files[0].path;
  }
  return "";
}
