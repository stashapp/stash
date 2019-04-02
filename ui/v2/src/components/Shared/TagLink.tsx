import {
  ITagProps,
  Tag,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { TagDataFragment } from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";

interface IProps extends ITagProps {
  tag: TagDataFragment;
}

export const TagLink: FunctionComponent<IProps> = (props: IProps) => {
  return (
    <Tag
      className="tag-item"
      interactive={true}
    >
      <Link to={NavigationUtils.makeTagScenesUrl(props.tag)}>{props.tag.name}</Link>
    </Tag>
  );
};
