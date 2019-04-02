import {
  Button,
  ButtonGroup,
  Card,
  Divider,
  Elevation,
  H4,
  Popover,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent, RefObject, useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { VideoHoverHook } from "../../hooks/VideoHover";
import { ColorUtils } from "../../utils/color";
import { TextUtils } from "../../utils/text";
import { TagLink } from "../Shared/TagLink";
import { SceneHelpers } from "./helpers";

interface ISceneCardProps {
  scene: GQL.SlimSceneDataFragment;
}

export const SceneCard: FunctionComponent<ISceneCardProps> = (props: ISceneCardProps) => {
  const [previewPath, setPreviewPath] = useState<string | undefined>(undefined);
  const videoHoverHook = VideoHoverHook.useVideoHover();

  function maybeRenderRatingBanner() {
    if (!props.scene.rating) { return; }
    return (
      <div className={`rating-banner ${ColorUtils.classForRating(props.scene.rating)}`}>
        RATING: {props.scene.rating}
      </div>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.scene.tags.length <= 0) { return; }

    const tags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <Popover interactionKind={"hover"} position="bottom">
        <Button
          icon="tag"
          text={props.scene.tags.length}
        />
        <>{tags}</>
      </Popover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.scene.performers.length <= 0) { return; }

    const performers = props.scene.performers.map((performer) => (
      <TagLink key={performer.id} performer={performer} />
    ));
    return (
      <Popover interactionKind={"hover"} position="bottom">
        <Button
          icon="person"
          text={props.scene.performers.length}
        />
        <>{performers}</>
      </Popover>
    );
  }

  function maybeRenderSceneMarkerPopoverButton() {
    if (props.scene.scene_markers.length <= 0) { return; }

    const sceneMarkers = props.scene.scene_markers.map((marker) => {
      (marker as any).scene = {};
      (marker as any).scene.id = props.scene.id;
      return <TagLink key={marker.id} marker={marker} />;
    });
    return (
      <Popover interactionKind={"hover"} position="bottom">
        <Button
          icon="map-marker"
          text={props.scene.scene_markers.length}
        />
        <>{sceneMarkers}</>
      </Popover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (props.scene.tags.length > 0 ||
        props.scene.performers.length > 0 ||
        props.scene.scene_markers.length > 0) {
      return (
        <>
          <Divider />
          <ButtonGroup minimal={true} className="card-section centered">
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderSceneMarkerPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  function onMouseEnter() {
    if (!previewPath || previewPath === "") {
      setPreviewPath(props.scene.paths.preview || "");
    }
    VideoHoverHook.onMouseEnter(videoHoverHook);
  }
  function onMouseLeave() {
    VideoHoverHook.onMouseLeave(videoHoverHook);
    setPreviewPath("");
  }

  return (
    <Card
      className="grid-item"
      elevation={Elevation.ONE}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <Link to={`/scenes/${props.scene.id}`} className="image previewable">
        {maybeRenderRatingBanner()}
        <video className="preview" loop={true} poster={props.scene.paths.screenshot || ""} ref={videoHoverHook.videoEl}>
          {!!previewPath ? <source src={previewPath} /> : ""}
        </video>
      </Link>
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {!!props.scene.title ? props.scene.title : TextUtils.fileNameFromPath(props.scene.path)}
        </H4>
        <span className="bp3-text-small bp3-text-muted">{props.scene.date}</span>
        <p>{TextUtils.truncate(props.scene.details, 100, "... (continued)")}</p>
      </div>

      {maybeRenderPopoverButtonGroup()}

      <Divider />
      <span className="card-section centered">
        {props.scene.file.size !== undefined ? TextUtils.fileSize(parseInt(props.scene.file.size, 10)) : ""}
        &nbsp;|&nbsp;
        {props.scene.file.duration !== undefined ? TextUtils.secondsToTimestamp(props.scene.file.duration) : ""}
        &nbsp;|&nbsp;
        {props.scene.file.width} x {props.scene.file.height}
      </span>
      {SceneHelpers.maybeRenderStudio(props.scene, 50, true)}
    </Card>
  );
};
