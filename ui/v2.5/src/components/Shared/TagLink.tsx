import { Badge } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import {
  PerformerDataFragment,
  SceneMarkerDataFragment,
  TagDataFragment
} from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";

interface IProps {
  tag?: Partial<TagDataFragment>;
  performer?: Partial<PerformerDataFragment>;
  marker?: Partial<SceneMarkerDataFragment>;
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
  } else if (props.marker) {
    link = NavUtils.makeSceneMarkerUrl(props.marker);
    title = `${props.marker.title} - ${TextUtils.secondsToTimestamp(
      props.marker.seconds || 0
    )}`;
  }
  return (
    <Badge className={`tag-item ${props.className}`} variant="secondary">
      <Link to={link}>{title}</Link>
    </Badge>
  );
};
