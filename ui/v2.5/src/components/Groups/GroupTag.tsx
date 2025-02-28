import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { GroupLink } from "../Shared/TagLink";

export const GroupTag: React.FC<{
  group: Pick<GQL.GroupDataFragment, "id" | "name" | "front_image_path">;
  linkType?: "scene" | "sub_group" | "details";
  description?: string;
}> = ({ group, linkType, description }) => {
  return (
    <div className="group-tag-container">
      <Link to={`/groups/${group.id}`} className="group-tag col m-auto zoom-2">
        <img
          className="image-thumbnail"
          alt={group.name ?? ""}
          src={group.front_image_path ?? ""}
        />
      </Link>
      <GroupLink
        group={group}
        description={description}
        linkType={linkType}
        className="d-block"
      />
    </div>
  );
};
