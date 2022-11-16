import { SceneMarker, Tag } from "./generated-graphql";

type SceneMarkerFragment = Pick<SceneMarker, "id" | "title"> & {
  primary_tag: Pick<Tag, "id" | "name">;
};

export function markerTitle(s: SceneMarkerFragment) {
  if (s.title) {
    return s.title;
  }

  if (s.primary_tag?.name) {
    return s.primary_tag?.name;
  }

  return "";
}
