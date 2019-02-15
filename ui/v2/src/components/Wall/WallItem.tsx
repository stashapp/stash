import _ from "lodash";
import React, { FunctionComponent, useRef, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { VideoHoverHook } from "../../hooks/VideoHover";
import { TextUtils } from "../../utils/text";

interface IWallItemProps {
  scene?: GQL.SlimSceneDataFragment;
  sceneMarker?: GQL.SceneMarkerDataFragment;
  origin?: string;
  onOverlay: (show: boolean) => void;
  clickHandler?: (item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment) => void;
}

export const WallItem: FunctionComponent<IWallItemProps> = (props: IWallItemProps) => {
  const [videoPath, setVideoPath] = useState<string | undefined>(undefined);
  const videoHoverHook = VideoHoverHook.useVideoHover({resetOnMouseLeave: true});

  function onMouseEnter() {
    VideoHoverHook.onMouseEnter(videoHoverHook);
    if (!videoPath || videoPath === "") {
      if (!!props.sceneMarker) {
        setVideoPath(props.sceneMarker.stream || "");
      } else if (!!props.scene) {
        setVideoPath(props.scene.paths.preview || "");
      }
    }
    props.onOverlay(true);
  }
  const debouncedOnMouseEnter = useRef(_.debounce(onMouseEnter, 500));

  function onMouseLeave() {
    VideoHoverHook.onMouseLeave(videoHoverHook);
    setVideoPath("");
    debouncedOnMouseEnter.current.cancel();
    props.onOverlay(false);
  }

  function onClick() {
    if (props.clickHandler === undefined) { return; }
    if (props.scene !== undefined) {
      props.clickHandler(props.scene);
    } else if (props.sceneMarker !== undefined) {
      props.clickHandler(props.sceneMarker);
    }
  }

  let linkSrc: string = "#";
  if (props.clickHandler === undefined) {
    if (props.scene !== undefined) {
      linkSrc = `/scenes/${props.scene.id}`;
    } else if (props.sceneMarker !== undefined) {
      linkSrc = `/scenes/${props.sceneMarker.scene.id}?t=${props.sceneMarker.seconds}`;
    }
  }

  function onTransitionEnd(event: React.TransitionEvent<HTMLDivElement>) {
    const target = (event.target as any);
    if (target.classList.contains("double-scale")) {
      target.parentElement.style.zIndex = 10;
    } else {
      target.parentElement.style.zIndex = null;
    }
  }

  let previewSrc: string = "";
  let title: string = "";
  let tags: JSX.Element[] = [];
  if (!!props.sceneMarker) {
    previewSrc = props.sceneMarker.preview;
    title = `${props.sceneMarker!.title} - ${TextUtils.secondsToTimestamp(props.sceneMarker.seconds)}`;
    tags = props.sceneMarker.tags.map((tag) => (<span key={tag.id}>{tag.name}</span>));
    tags.unshift(<span key={props.sceneMarker.primary_tag.id}>{props.sceneMarker.primary_tag.name}</span>);
  } else if (!!props.scene) {
    previewSrc = props.scene.paths.webp || "";
    title = props.scene.title || "";
    // tags = props.scene.tags.map((tag) => (<span key={tag.id}>{tag.name}</span>));
  }

  const className = ["scene-wall-item-container"];
  if (videoHoverHook.isHovering.current) { className.push("double-scale"); }
  const style: React.CSSProperties = {};
  if (!!props.origin) { style.transformOrigin = props.origin; }
  return (
    <div className="wall grid-item">
      <div
        className={className.join(" ")}
        style={style}
        onTransitionEnd={onTransitionEnd}
        onMouseEnter={() => debouncedOnMouseEnter.current()}
        onMouseMove={() => debouncedOnMouseEnter.current()}
        onMouseLeave={onMouseLeave}
      >
        <Link onClick={() => onClick()} to={linkSrc}>
          <video
            src={videoPath}
            style={videoHoverHook.isHovering.current ? {} : {display: "none"}}
            autoPlay={true}
            loop={true}
            ref={videoHoverHook.videoEl}
          />
          <img src={previewSrc} />
          <div className="scene-wall-item-text-container">
            <div style={{lineHeight: 1}}>
              {title}
            </div>
            {tags}
          </div>
        </Link>
      </div>
    </div>
  );
};
