import _ from "lodash";
import React, { useRef, useState, useEffect } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { useConfiguration } from "src/core/StashService";
import { useVideoHover } from "src/hooks";
import { TextUtils, NavUtils } from "src/utils";

interface IWallItemProps {
  scene?: GQL.SlimSceneDataFragment;
  sceneMarker?: GQL.SceneMarkerDataFragment;
  origin?: string;
  onOverlay: (show: boolean) => void;
  clickHandler?: (
    item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment
  ) => void;
}

export const WallItem: React.FC<IWallItemProps> = (props: IWallItemProps) => {
  const [videoPath, setVideoPath] = useState<string>();
  const [previewPath, setPreviewPath] = useState<string>("");
  const [screenshotPath, setScreenshotPath] = useState<string>("");
  const [title, setTitle] = useState<string>("");
  const [tags, setTags] = useState<JSX.Element[]>([]);
  const config = useConfiguration();
  const hoverHandler = useVideoHover({
    resetOnMouseLeave: true,
  });
  const showTextContainer =
    config.data?.configuration.interface.wallShowTitle ?? true;

  function onMouseEnter() {
    hoverHandler.onMouseEnter();
    if (!videoPath || videoPath === "") {
      if (props.sceneMarker) {
        setVideoPath(props.sceneMarker.stream || "");
      } else if (props.scene) {
        setVideoPath(props.scene.paths.preview || "");
      }
    }
    props.onOverlay(true);
  }
  const debouncedOnMouseEnter = useRef(_.debounce(onMouseEnter, 500));

  function onMouseLeave() {
    hoverHandler.onMouseLeave();
    setVideoPath("");
    debouncedOnMouseEnter.current.cancel();
    props.onOverlay(false);
  }

  function onClick() {
    if (props.clickHandler === undefined) {
      return;
    }
    if (props.scene !== undefined) {
      props.clickHandler(props.scene);
    } else if (props.sceneMarker !== undefined) {
      props.clickHandler(props.sceneMarker);
    }
  }

  let linkSrc: string = "#";
  if (!props.clickHandler) {
    if (props.scene) {
      linkSrc = `/scenes/${props.scene.id}`;
    } else if (props.sceneMarker) {
      linkSrc = NavUtils.makeSceneMarkerUrl(props.sceneMarker);
    }
  }

  function onTransitionEnd(event: React.TransitionEvent<HTMLDivElement>) {
    const target = event.currentTarget;
    if (target.classList.contains("double-scale") && target.parentElement) {
      target.parentElement.style.zIndex = "10";
    } else if (target.parentElement) {
      target.parentElement.style.zIndex = "";
    }
  }

  useEffect(() => {
    if (props.sceneMarker) {
      setPreviewPath(props.sceneMarker.preview);
      setTitle(
        `${props.sceneMarker!.title} - ${TextUtils.secondsToTimestamp(
          props.sceneMarker.seconds
        )}`
      );
      const thisTags = props.sceneMarker.tags.map((tag) => (
        <span key={tag.id} className="wall-tag">
          {tag.name}
        </span>
      ));
      thisTags.unshift(
        <span key={props.sceneMarker.primary_tag.id} className="wall-tag">
          {props.sceneMarker.primary_tag.name}
        </span>
      );
      setTags(thisTags);
    } else if (props.scene) {
      setPreviewPath(props.scene.paths.webp || "");
      setScreenshotPath(props.scene.paths.screenshot || "");
      setTitle(props.scene.title || "");
    }
  }, [props.sceneMarker, props.scene]);

  function previewNotFound() {
    if (previewPath !== screenshotPath) {
      setPreviewPath(screenshotPath);
    }
  }

  const className = ["scene-wall-item-container"];
  if (hoverHandler.isHovering.current) {
    className.push("double-scale");
  }
  const style: React.CSSProperties = {};
  if (props.origin) {
    style.transformOrigin = props.origin;
  }
  return (
    <div className="wall-item">
      <div
        className={className.join(" ")}
        style={style}
        onTransitionEnd={onTransitionEnd}
        onMouseEnter={() => debouncedOnMouseEnter.current()}
        onMouseMove={() => debouncedOnMouseEnter.current()}
        onMouseLeave={onMouseLeave}
      >
        <Link onClick={onClick} to={linkSrc}>
          <video
            src={videoPath}
            poster={screenshotPath}
            className="scene-wall-video"
            style={hoverHandler.isHovering.current ? {} : { display: "none" }}
            autoPlay
            loop
            ref={hoverHandler.videoEl}
          />
          <img
            alt={title}
            className="scene-wall-image"
            src={previewPath || screenshotPath}
            onError={() => previewNotFound()}
          />
          {showTextContainer ? (
            <div className="scene-wall-item-text-container">
              <div>{title}</div>
              {tags}
            </div>
          ) : (
            ""
          )}
        </Link>
      </div>
    </div>
  );
};
