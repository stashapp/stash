import {
  ITagProps,
  Tag,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { DvdDataFragment,PerformerDataFragment, SceneMarkerDataFragment, TagDataFragment } from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";
import { TextUtils } from "../../utils/text";

interface IProps extends ITagProps {
  tag?: Partial<TagDataFragment>;
  performer?: Partial<PerformerDataFragment>;
  dvd?: Partial<DvdDataFragment>;
  marker?: Partial<SceneMarkerDataFragment>;
}

export const TagLink: FunctionComponent<IProps> = (props: IProps) => {
  let link: string = "#";
  let title: string = "";
  if (!!props.tag) {
    link = NavigationUtils.makeTagScenesUrl(props.tag);
    title = props.tag.name || "";
  } else if (!!props.performer) {
    link = NavigationUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (!!props.dvd) {
    link = NavigationUtils.makeDvdScenesUrl(props.dvd);
    title = props.dvd.name || "";
  } else if (!!props.marker) {
    link = NavigationUtils.makeSceneMarkerUrl(props.marker);
    title = `${props.marker.title} - ${TextUtils.secondsToTimestamp(props.marker.seconds || 0)}`;
  }
  return (
    <Tag
      className="tag-item"
      interactive={true}
    >
      <Link to={link}>{title}</Link>
    </Tag>
  );
};
