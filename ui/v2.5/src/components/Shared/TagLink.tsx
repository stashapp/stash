import { Badge } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import {
  PerformerDataFragment,
  TagDataFragment,
  MovieDataFragment,
  SceneDataFragment,
} from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import * as GQL from "src/core/generated-graphql";
import { TagPopover } from "../Tags/TagPopover";
import { markerTitle } from "src/core/markers";
import { contrastingTextColor } from "src/utils/display";
import styles from "src/styles/globalStyles.module.scss";

interface IFile {
  path: string;
}
interface IGallery {
  id: string;
  files: IFile[];
  folder?: GQL.Maybe<IFile>;
  title: GQL.Maybe<string>;
}

type SceneMarkerFragment = Pick<GQL.SceneMarker, "id" | "title" | "seconds"> & {
  scene: Pick<GQL.Scene, "id">;
  primary_tag: Pick<GQL.Tag, "id" | "name">;
};

interface IProps {
  tag?: Partial<TagDataFragment>;
  tagType?: "performer" | "scene" | "gallery" | "image";
  performer?: Partial<PerformerDataFragment>;
  marker?: SceneMarkerFragment;
  movie?: Partial<MovieDataFragment>;
  scene?: Partial<Pick<SceneDataFragment, "id" | "title" | "files">>;
  gallery?: Partial<IGallery>;
  className?: string;
}

export const TagLink: React.FC<IProps> = (props: IProps) => {
  let id: string = "";
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
    id = props.tag.id || "";
    title = props.tag.name || "";
  } else if (props.performer) {
    link = NavUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (props.movie) {
    link = NavUtils.makeMovieScenesUrl(props.movie);
    title = props.movie.name || "";
  } else if (props.marker) {
    link = NavUtils.makeSceneMarkerUrl(props.marker);
    title = `${markerTitle(props.marker)} - ${TextUtils.secondsToTimestamp(
      props.marker.seconds || 0
    )}`;
  } else if (props.gallery) {
    link = `/galleries/${props.gallery.id}`;
    title = galleryTitle(props.gallery);
  } else if (props.scene) {
    link = `/scenes/${props.scene.id}`;
    title = objectTitle(props.scene);
  }

  return (
    <Badge className={cx("tag-item", props.className)} variant="secondary" 
      style={{
        ["--tag-bg-color" as  string]: props.tag?.color ?? styles.textMuted,
        ["--tag-text-color" as  string]: contrastingTextColor(props.tag?.color) ?? styles.darkText
      }}
    >
      <TagPopover id={id}>
        <Link to={link}>{title}</Link>
      </TagPopover>
    </Badge>
  );
};
