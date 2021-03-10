import { Button, ButtonGroup } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils } from "src/utils";
import { Icon, TruncatedText } from "../Shared";
import { BasicCard } from "../Shared/BasicCard";

interface IProps {
  tag: GQL.TagDataFragment;
  zoomIndex: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const TagCard: React.FC<IProps> = ({
  tag,
  zoomIndex,
  selecting,
  selected,
  onSelectedChanged,
}) => {
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

  function maybeRenderPerformersPopoverButton() {
    if (!tag.performer_count) return;

    return (
      <Link to={NavUtils.makeTagPerformersUrl(tag)}>
        <Button className="minimal">
          <Icon icon="user" />
          <span>{tag.performer_count}</span>
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
            {maybeRenderPerformersPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <BasicCard
      className={`tag-card zoom-${zoomIndex}`}
      url={`/tags/${tag.id}`}
      linkClassName="tag-card-header"
      image={
        <img
          className="tag-card-image"
          alt={tag.name}
          src={tag.image_path ?? ""}
        />
      }
      details={
        <h5>
          <TruncatedText text={tag.name} />
        </h5>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
