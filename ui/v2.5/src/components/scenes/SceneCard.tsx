import React, { useState } from "react";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { VideoHoverHook } from "src/hooks";
import { Icon, TagLink, HoverPopover } from "src/components/Shared";
import { TextUtils } from "src/utils";

interface ISceneCardProps {
  scene: GQL.SlimSceneDataFragment;
  selected: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged: (selected: boolean, shiftKey: boolean) => void;
}

export const SceneCard: React.FC<ISceneCardProps> = (
  props: ISceneCardProps
) => {
  const [previewPath, setPreviewPath] = useState<string>();
  const videoHoverHook = VideoHoverHook.useVideoHover({
    resetOnMouseLeave: false
  });

  const config = StashService.useConfiguration();
  const showStudioAsText =
    config?.data?.configuration.interface.showStudioAsText ?? false;

  function maybeRenderRatingBanner() {
    if (!props.scene.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          props.scene.rating ? `rating-${props.scene.rating}` : ""
        }`}
      >
        RATING: {props.scene.rating}
      </div>
    );
  }

  function maybeRenderSceneSpecsOverlay() {
    return (
      <div className="scene-specs-overlay">
        {props.scene.file.height ? (
          <span className="overlay-resolution">
            {" "}
            {TextUtils.resolution(props.scene.file.height)}
          </span>
        ) : (
          ""
        )}
        {(props.scene.file.duration ?? 0) >= 1
          ? TextUtils.secondsToTimestamp(props.scene.file.duration ?? 0)
          : ""}
      </div>
    );
  }

  function maybeRenderSceneStudioOverlay() {
    if (!props.scene.studio) return;

    let style: React.CSSProperties = {
      backgroundImage: `url('${props.scene.studio.image_path}')`
    };

    let text = "";
    if (showStudioAsText) {
      style = {};
      text = props.scene.studio.name;
    }

    return (
      <div className="scene-studio-overlay">
        <Link to={`/studios/${props.scene.studio.id}`} style={style}>
          {text}
        </Link>
      </div>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.scene.tags.length <= 0) return;

    const popoverContent = props.scene.tags.map(tag => (
      <TagLink key={tag.id} tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="tag" />
          {props.scene.tags.length}
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.scene.performers.length <= 0) return;

    const popoverContent = props.scene.performers.map(performer => (
      <div className="performer-tag-container" key="performer">
        <Link
          to={`/performers/${performer.id}`}
          className="performer-tag previewable image"
          style={{ backgroundImage: `url(${performer.image_path})` }}
        />
        <TagLink key={performer.id} performer={performer} />
      </div>
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="user" />
          {props.scene.performers.length}
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderSceneMarkerPopoverButton() {
    if (props.scene.scene_markers.length <= 0) return;

    const popoverContent = props.scene.scene_markers.map(marker => {
      const markerPopover = { ...marker, scene: { id: props.scene.id } };
      return <TagLink key={marker.id} marker={markerPopover} />;
    });

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="map-marker-alt" />
          {props.scene.scene_markers.length}
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.scene.tags.length > 0 ||
      props.scene.performers.length > 0 ||
      props.scene.scene_markers.length > 0
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="scene-popovers">
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
    const { file } = props.scene;
    const width = file.width ? file.width : 0;
    const height = file.height ? file.height : 0;
    return height > width;
  }

  let shiftKey = false;

  return (
    <Card
      className={`zoom-${props.zoomIndex} scene-card`}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <Form.Control
        type="checkbox"
        className="card-select"
        checked={props.selected}
        onChange={() => props.onSelectedChanged(!props.selected, shiftKey)}
        onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
          // eslint-disable-next-line prefer-destructuring
          shiftKey = event.shiftKey;
          event.stopPropagation();
        }}
      />
      {maybeRenderSceneStudioOverlay()}
      <Link
        to={`/scenes/${props.scene.id}`}
        className={cx("image", "previewable", { portrait: isPortrait() })}
      >
        <div className="video-container">
          {maybeRenderRatingBanner()}
          {maybeRenderSceneSpecsOverlay()}
          <video
            loop
            className={cx("preview", { portrait: isPortrait() })}
            poster={props.scene.paths.screenshot || ""}
            ref={videoHoverHook.videoEl}
          >
            {previewPath ? <source src={previewPath} /> : ""}
          </video>
        </div>
      </Link>
      <div className="card-section">
        <h5 className="text-truncate">
          {props.scene.title
            ? props.scene.title
            : TextUtils.fileNameFromPath(props.scene.path)}
        </h5>
        <span>{props.scene.date}</span>
        <p>
          {TextUtils.truncate(
            props.scene.details ?? "",
            100,
            "... (continued)"
          )}
        </p>
      </div>

      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
