import {
  Button,
  ButtonGroup,
  Card,
  Checkbox,
  Divider,
  Elevation,
  H4,
  Popover,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { VideoHoverHook } from "../../hooks/VideoHover";
import { ColorUtils } from "../../utils/color";
import { TextUtils } from "../../utils/text";
import { TagLink } from "../Shared/TagLink";
import { ZoomUtils } from "../../utils/zoom";
import { StashService } from "../../core/StashService";

interface ISceneCardProps {
  scene: GQL.SlimSceneDataFragment;
  selected: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged: (selected : boolean, shiftKey : boolean) => void;
}

export const SceneCard: FunctionComponent<ISceneCardProps> = (props: ISceneCardProps) => {
  const [previewPath, setPreviewPath] = useState<string | undefined>(undefined);
  const videoHoverHook = VideoHoverHook.useVideoHover({resetOnMouseLeave: false});
  
  const config = StashService.useConfiguration();
  const showStudioAsText = !!config.data && !!config.data.configuration ? config.data.configuration.interface.showStudioAsText : false;

  function maybeRenderRatingBanner() {
    if (!props.scene.rating) { return; }
    return (
      <div className={`rating-banner ${ColorUtils.classForRating(props.scene.rating)}`}>
        RATING: {props.scene.rating}
      </div>
    );
  }

  function maybeRenderSceneSpecsOverlay() {
    return (
      <div className={`scene-specs-overlay`}>
        {!!props.scene.file.height ? <span className={`overlay-resolution`}> {TextUtils.resolution(props.scene.file.height)}</span> : undefined}
        {props.scene.file.duration !== undefined && props.scene.file.duration >= 1 ? TextUtils.secondsToTimestamp(props.scene.file.duration) : ""}
      </div>
    );
  }

  function maybeRenderSceneStudioOverlay() {
    if (!props.scene.studio) {
      return;
    }

    let style: React.CSSProperties = {
      backgroundImage: `url('${props.scene.studio.image_path}')`,
    };

    let text = "";

    if (showStudioAsText) {
      style = {};
      text = props.scene.studio.name;
    }

    return (
      <div className={`scene-studio-overlay`}>
        <Link
          to={`/studios/${props.scene.studio.id}`}
          style={style}
        >
          {text}
        </Link>
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

    const performers = props.scene.performers.map((performer) => {
      return (
        <>
        <div className="performer-tag-container">
          <Link
            to={`/performers/${performer.id}`}
            className="performer-tag previewable image"
            style={{backgroundImage: `url(${performer.image_path})`}}
          ></Link>
          <TagLink key={performer.id} performer={performer} />
        </div>
        </>
      );
    });
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

  function isPortrait() {
    let file = props.scene.file;
    let width = file.width ? file.width : 0;
    let height = file.height ? file.height : 0;
    return height > width;
  }

  function getLinkClassName() {
    let ret = "image previewable";
    
    if (isPortrait()) {
      ret += " portrait";
    }

    return ret;
  }

  function getVideoClassName() {
    let ret = "preview";
    
    if (isPortrait()) {
      ret += " portrait";
    }

    return ret;
  }

  var shiftKey = false;

  return (
    <Card
      className={"grid-item scene-card " + ZoomUtils.classForZoom(props.zoomIndex)}
      elevation={Elevation.ONE}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <Checkbox
        className="card-select"
        checked={props.selected}
        onChange={() => props.onSelectedChanged(!props.selected, shiftKey)}
        onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => { shiftKey = event.shiftKey; event.stopPropagation(); } }
      />
      <Link to={`/scenes/${props.scene.id}`} className={getLinkClassName()}>
        <div className="video-container">
          {maybeRenderRatingBanner()}
          {maybeRenderSceneSpecsOverlay()}
          {maybeRenderSceneStudioOverlay()}
          <video className={getVideoClassName()} loop={true} poster={props.scene.paths.screenshot || ""} ref={videoHoverHook.videoEl}>
            {!!previewPath ? <source src={previewPath} /> : ""}
          </video>
        </div>
      </Link>
      <div className="card-section">
        <H4 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
          {!!props.scene.title ? props.scene.title : TextUtils.fileNameFromPath(props.scene.path)}
        </H4>
        <span className="bp3-text-small bp3-text-muted">{props.scene.date}</span>
        <p>{TextUtils.truncate(props.scene.details, 100, "... (continued)")}</p>
      </div>

      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
