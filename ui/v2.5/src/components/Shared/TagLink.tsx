import { Badge } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import {
  PerformerDataFragment,
  SceneMarkerDataFragment,
  TagDataFragment,
  MovieDataFragment,
  SceneDataFragment,
} from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";

interface IProps {
  tag?: Partial<TagDataFragment>;
  performer?: Partial<PerformerDataFragment>;
  marker?: Partial<SceneMarkerDataFragment>;
  movie?: Partial<MovieDataFragment>;
  scene?: Partial<SceneDataFragment>;
  className?: string;
}

export const TagLink: React.FC<IProps> = (props: IProps) => {
  let link: string = "#";
  let title: string = "";
  if (props.tag) {
    link = NavUtils.makeTagScenesUrl(props.tag);
    title = props.tag.name || "";
  } else if (props.performer) {
    link = NavUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (props.movie) {
    link = NavUtils.makeMovieScenesUrl(props.movie);
    title = props.movie.name || "";
  } else if (props.marker) {
    link = NavUtils.makeSceneMarkerUrl(props.marker);
    title = `${props.marker.title} - ${TextUtils.secondsToTimestamp(
      props.marker.seconds || 0
    )}`;
  } else if (props.scene) {
    link = `/scenes/${props.scene.id}`;
    title = props.scene.title
      ? props.scene.title
      : TextUtils.fileNameFromPath(props.scene.path ?? "");
  }
  return (
    <Badge className={`tag-item ${props.className}`} variant="secondary">
      <Link to={link}>{title}</Link>
    </Badge>
  );
};
