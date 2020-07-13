import { Card, Button, ButtonGroup } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils } from "src/utils";
import { Icon } from "../Shared";

interface IProps {
  tag: GQL.TagDataFragment;
  zoomIndex: number;
}

export const TagCard: React.FC<IProps> = ({ tag, zoomIndex }) => {
  function maybeRenderScenesPopoverButton() {
    if (!tag.scene_count) return;

    return (
      <Link to={NavUtils.makeTagScenesUrl(tag)}>
        <Button className="minimal">
          <Icon icon="play-circle" />
          <span>{tag.scene_count}</span>
        </Button>
      </Link>
    );
  }

  function maybeRenderSceneMarkersPopoverButton() {
    if (!tag.scene_marker_count) return;

    return (
      <Link to={NavUtils.makeTagSceneMarkersUrl(tag)}>
        <Button className="minimal">
          <Icon icon="map-marker-alt" />
          <span>{tag.scene_marker_count}</span>
        </Button>
      </Link>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (tag) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderSceneMarkersPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <Card className={`tag-card zoom-${zoomIndex}`}>
      <Link to={`/tags/${tag.id}`} className="tag-card-header">
        <img
          className="tag-card-image"
          alt={tag.name}
          src={tag.image_path ?? ""}
        />
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">{tag.name}</h5>
      </div>
      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
