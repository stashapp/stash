import { Badge } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import {
  PerformerDataFragment,
  SceneMarkerDataFragment,
  TagDataFragment,
  MovieDataFragment,
  SceneDataFragment,
} from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import * as GQL from "src/core/generated-graphql";

interface IFile {
  path: string;
}
interface IGallery {
  id: string;
  files: IFile[];
  folder?: GQL.Maybe<IFile>;
  title: GQL.Maybe<string>;
}

interface IProps {
  tag?: Partial<TagDataFragment>;
  tagType?: "performer" | "scene" | "gallery" | "image";
  performer?: Partial<PerformerDataFragment>;
  marker?: Partial<SceneMarkerDataFragment>;
  movie?: Partial<MovieDataFragment>;
  scene?: Partial<Pick<SceneDataFragment, "id" | "title" | "files">>;
  gallery?: Partial<IGallery>;
  className?: string;
}

export const TagLink: React.FC<IProps> = (props: IProps) => {
  let link: string = "#";
  let title: string = "";
  if (props.tag) {
    switch (props.tagType) {
      case "scene":
      case undefined:
        link = NavUtils.makeTagScenesUrl(props.tag);
        break;
      case "performer":
        link = NavUtils.makeTagPerformersUrl(props.tag);
        break;
      case "gallery":
        link = NavUtils.makeTagGalleriesUrl(props.tag);
        break;
      case "image":
        link = NavUtils.makeTagImagesUrl(props.tag);
        break;
    }
    title = props.tag.name || "";
  } else if (props.performer) {
    link = NavUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (props.movie) {
    link = NavUtils.makeMovieScenesUrl(props.movie);
    title = props.movie.name || "";
  } else if (props.marker) {
    link = NavUtils.makeSceneMarkerUrl(props.marker);
    title = `${
      props.marker.title || props.marker.primary_tag?.name || ""
    } - ${TextUtils.secondsToTimestamp(props.marker.seconds || 0)}`;
  } else if (props.gallery) {
    link = `/galleries/${props.gallery.id}`;
    title = galleryTitle(props.gallery);
  } else if (props.scene) {
    link = `/scenes/${props.scene.id}`;
    title = objectTitle(props.scene);
  }
  return (
    <Badge className={cx("tag-item", props.className)} variant="secondary">
      <Link to={link}>{title}</Link>
    </Badge>
  );
};
