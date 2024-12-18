import { faUser } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { sortPerformers } from "src/core/performers";
import { HoverPopover } from "./HoverPopover";
import { Icon } from "./Icon";
import { PerformerLink, PerformerLinkType } from "./TagLink";

interface IProps {
  performers: Pick<
    GQL.Performer,
    "id" | "name" | "image_path" | "disambiguation" | "gender"
  >[];
  linkType?: PerformerLinkType;
}

export const PerformerPopoverButton: React.FC<IProps> = ({
  performers,
  linkType,
}) => {
  const sorted = sortPerformers(performers);
  const popoverContent = sorted.map((performer) => (
    <div className="performer-tag-container row" key={performer.id}>
      <Link
        to={`/performers/${performer.id}`}
        className="performer-tag col m-auto zoom-2"
      >
        <img
          className="image-thumbnail"
          alt={performer.name ?? ""}
          src={performer.image_path ?? ""}
        />
      </Link>
      <PerformerLink
        key={performer.id}
        performer={performer}
        className="d-block"
        linkType={linkType}
      />
    </div>
  ));

  return (
    <HoverPopover
      className="performer-count"
      placement="bottom"
      content={popoverContent}
    >
      <Button className="minimal">
        <Icon icon={faUser} />
        <span>{performers.length}</span>
      </Button>
    </HoverPopover>
  );
};
