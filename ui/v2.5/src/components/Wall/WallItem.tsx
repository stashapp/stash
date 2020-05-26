import React, { useRef, useState, useEffect } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { useConfiguration } from "src/core/StashService";
import { TextUtils, NavUtils } from "src/utils";
import cx from "classnames";

interface IWallItemProps {
  scene?: GQL.SlimSceneDataFragment;
  sceneMarker?: GQL.SceneMarkerDataFragment;
  clickHandler?: (
    item: GQL.SlimSceneDataFragment | GQL.SceneMarkerDataFragment
  ) => void;
  className: string;
}

interface IPreviews {
  video?: string;
  animation?: string;
  image?: string;
}

const Preview: React.FC<{
  previews?: IPreviews;
  config?: GQL.ConfigDataFragment;
  active: boolean;
}> = ({ previews, config, active }) => {
  const videoElement = useRef() as React.MutableRefObject<HTMLVideoElement>;
  const [isMissing, setIsMissing] = useState(false);

  const previewType = config?.interface?.wallPlayback;
  const soundOnPreview = config?.interface?.soundOnPreview ?? false;

  useEffect(() => {
    if (!videoElement.current) return;
    videoElement.current.muted = !(soundOnPreview && active);
    if (previewType !== "video") {
      if (active) videoElement.current.play();
      else videoElement.current.pause();
    }
  }, [videoElement, previewType, soundOnPreview, active]);

  if (!previews) return <div />;

  if (isMissing) {
    return (
      <div className="wall-item-media wall-item-missing">
        Pending preview generation
      </div>
    );
  }

  const image = (
    <img
      alt=""
      className="wall-item-media"
      src={
        (previewType === "animation" && previews.animation) || previews.image
      }
    />
  );
  const video = (
    <video
      src={previews.video}
      poster={previews.image}
      autoPlay={previewType === "video"}
      loop
      muted
      className={cx("wall-item-media", {
        "wall-item-preview": previewType !== "video",
      })}
      onError={() => setIsMissing(true)}
      ref={videoElement}
    />
  );

  if (previewType === "video") {
    return video;
  }
  return (
    <>
      {image}
      {video}
    </>
  );
};

export const WallItem: React.FC<IWallItemProps> = (props: IWallItemProps) => {
  const [active, setActive] = useState(false);
  const wallItem = useRef() as React.MutableRefObject<HTMLDivElement>;
  const config = useConfiguration();

  const showTextContainer =
    config.data?.configuration.interface.wallShowTitle ?? true;

  const previews = props.sceneMarker
    ? {
        video: props.sceneMarker.stream,
        animation: props.sceneMarker.preview,
      }
    : {
        video: props.scene?.paths.preview ?? undefined,
        animation: props.scene?.paths.webp ?? undefined,
        image: props.scene?.paths.screenshot ?? undefined,
      };

  const setInactive = () => setActive(false);
  const toggleActive = (e: TransitionEvent) => {
    if (e.propertyName === "transform" && e.elapsedTime === 0) {
      // Get the current scale of the wall-item. If it's smaller than 1.1 the item is being scaled up, otherwise down.
      const matrixScale = getComputedStyle(wallItem.current).transform.match(
        /-?\d+\.?\d+|\d+/g
      )?.[0];
      const scale = Number.parseFloat(matrixScale ?? "2") || 2;
      setActive(scale <= 1.1 && !active);
    }
  };

  useEffect(() => {
    const { current } = wallItem;
    current?.addEventListener("transitioncancel", setInactive);
    current?.addEventListener("transitionstart", toggleActive);
    return () => {
      current?.removeEventListener("transitioncancel", setInactive);
      current?.removeEventListener("transitionstart", toggleActive);
    };
  });

  const clickHandler = () => {
    if (props.scene) {
      props?.clickHandler?.(props.scene);
    }
    if (props.sceneMarker) {
      props?.clickHandler?.(props.sceneMarker);
    }
  };

  let linkSrc: string = "#";
  if (!props.clickHandler) {
    if (props.scene) {
      linkSrc = `/scenes/${props.scene.id}`;
    } else if (props.sceneMarker) {
      linkSrc = NavUtils.makeSceneMarkerUrl(props.sceneMarker);
    }
  }

  const renderText = () => {
    if (!showTextContainer) return;

    const title = props.sceneMarker
      ? `${props.sceneMarker!.title} - ${TextUtils.secondsToTimestamp(
          props.sceneMarker.seconds
        )}`
      : props.scene?.title ?? "";
    const tags = props.sceneMarker
      ? [props.sceneMarker.primary_tag, ...props.sceneMarker.tags]
      : [];

    return (
      <div className="wall-item-text">
        <div>{title}</div>
        {tags.map((tag) => (
          <span key={tag.id} className="wall-tag">
            {tag.name}
          </span>
        ))}
      </div>
    );
  };

  return (
    <div className="wall-item">
      <div className={`wall-item-container ${props.className}`} ref={wallItem}>
        <Link onClick={clickHandler} to={linkSrc} className="wall-item-anchor">
          <Preview
            previews={previews}
            config={config.data?.configuration}
            active={active}
          />
          {renderText()}
        </Link>
      </div>
    </div>
  );
};
