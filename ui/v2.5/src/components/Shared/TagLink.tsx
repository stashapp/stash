import { Badge } from 'react-bootstrap';
import React from "react";
import { Link } from "react-router-dom";
import { PerformerDataFragment, SceneMarkerDataFragment, TagDataFragment } from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";
import { TextUtils } from "../../utils/text";

interface IProps {
  tag?: Partial<TagDataFragment>;
  performer?: Partial<PerformerDataFragment>;
  marker?: Partial<SceneMarkerDataFragment>;
}

export const TagLink: React.FC<IProps> = (props: IProps) => {
  let link: string = "#";
  let title: string = "";
  if (!!props.tag) {
    link = NavigationUtils.makeTagScenesUrl(props.tag);
    title = props.tag.name || "";
  } else if (!!props.performer) {
    link = NavigationUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (!!props.marker) {
    link = NavigationUtils.makeSceneMarkerUrl(props.marker);
    title = `${props.marker.title} - ${TextUtils.secondsToTimestamp(props.marker.seconds || 0)}`;
  }
  return (
    <Badge 
      className="tag-item"
      variant="secondary"
    >
      <Link to={link}>{title}</Link>
    </Badge>
  );
};
