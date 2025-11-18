import React, {
  useRef,
  useState,
  useEffect,
  useCallback,
  MouseEvent,
  useMemo,
} from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import NavUtils from "src/utils/navigation";
import cx from "classnames";
import { SceneQueue } from "src/models/sceneQueue";
import { useConfigurationContext } from "src/hooks/Config";
import { markerTitle } from "src/core/markers";
import { objectTitle } from "src/core/files";

export type WallItemType = keyof WallItemData;

export type WallItemData = {
  scene: GQL.SlimSceneDataFragment;
  sceneMarker: GQL.SceneMarkerDataFragment;
  image: GQL.SlimImageDataFragment;
};

interface IWallItemProps<T extends WallItemType> {
  type: T;
  index?: number;
  data: WallItemData[T];
  sceneQueue?: SceneQueue;
  clickHandler?: (e: MouseEvent, item: WallItemData[T]) => void;
  className: string;
}

interface IPreviews {
  video?: string;
  animation?: string;
  image?: string;
}

const Preview: React.FC<{
  previews: IPreviews;
  config?: GQL.ConfigDataFragment;
  active: boolean;
}> = ({ previews, config, active }) => {
  const videoEl = useRef<HTMLVideoElement>(null);
  const [isMissing, setIsMissing] = useState(false);

  const previewType = config?.interface?.wallPlayback;
  const soundOnPreview = config?.interface?.soundOnPreview ?? false;

  useEffect(() => {
    const video = videoEl.current;
    if (!video) return;

    video.muted = !(soundOnPreview && active);
    if (previewType !== "video") {
      if (active) {
        video.play();
      } else {
        video.pause();
      }
    }
  }, [previewType, soundOnPreview, active]);

  const image = (
    <img
      loading="lazy"
      alt=""
      className="wall-item-media"
      src={
        (previewType === "animation" && previews.animation) || previews.image
      }
    />
  );
  const video = (
    <video
      disableRemotePlayback
      playsInline
      src={previews.video}
      poster={previews.image}
      autoPlay={previewType === "video"}
      loop
      muted
      className={cx("wall-item-media", {
        "wall-item-preview": previewType !== "video",
      })}
      onError={(error: React.SyntheticEvent<HTMLVideoElement>) => {
        // Error code 4 indicates media not found or unsupported
        setIsMissing(error.currentTarget.error?.code === 4);
      }}
      ref={videoEl}
    />
  );

  if (isMissing) {
    // show the image if the video preview is unavailable
    if (previews.image) {
      return image;
    }

    return (
      <div className="wall-item-media wall-item-missing">
        Pending preview generation
      </div>
    );
  }

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

export const WallItem = <T extends WallItemType>({
  type,
  index,
  data,
  sceneQueue,
  clickHandler,
  className,
}: IWallItemProps<T>) => {
  const [active, setActive] = useState(false);
  const itemEl = useRef<HTMLDivElement>(null);
  const { configuration: config } = useConfigurationContext();

  const showTextContainer = config?.interface.wallShowTitle ?? true;

  const previews = useMemo(() => {
    switch (type) {
      case "scene":
        const scene = data as GQL.SlimSceneDataFragment;
        return {
          video: scene.paths.preview ?? undefined,
          animation: scene.paths.webp ?? undefined,
          image: scene.paths.screenshot ?? undefined,
        };
      case "sceneMarker":
        const sceneMarker = data as GQL.SceneMarkerDataFragment;
        return {
          video: sceneMarker.stream,
          animation: sceneMarker.preview,
          image: sceneMarker.screenshot,
        };
      case "image":
        const image = data as GQL.SlimImageDataFragment;
        return {
          image: image.paths.thumbnail ?? undefined,
        };
      default:
        // this is unreachable, inference fails for some reason
        return type as never;
    }
  }, [type, data]);
  const linkSrc = useMemo(() => {
    switch (type) {
      case "scene":
        const scene = data as GQL.SlimSceneDataFragment;
        return sceneQueue
          ? sceneQueue.makeLink(scene.id, { sceneIndex: index })
          : `/scenes/${scene.id}`;
      case "sceneMarker":
        const sceneMarker = data as GQL.SceneMarkerDataFragment;
        return NavUtils.makeSceneMarkerUrl(sceneMarker);
      case "image":
        const image = data as GQL.SlimImageDataFragment;
        return `/images/${image.id}`;
      default:
        return type;
    }
  }, [type, data, sceneQueue, index]);
  const title = useMemo(() => {
    switch (type) {
      case "scene":
        const scene = data as GQL.SlimSceneDataFragment;
        return objectTitle(scene);
      case "sceneMarker":
        const sceneMarker = data as GQL.SceneMarkerDataFragment;
        const newTitle = markerTitle(sceneMarker);
        const seconds = TextUtils.formatTimestampRange(
          sceneMarker.seconds,
          sceneMarker.end_seconds ?? undefined
        );
        if (newTitle) {
          return `${newTitle} - ${seconds}`;
        } else {
          return seconds;
        }
      case "image":
        return "";
      default:
        return type;
    }
  }, [type, data]);
  const tags = useMemo(() => {
    if (type === "sceneMarker") {
      const sceneMarker = data as GQL.SceneMarkerDataFragment;
      return [sceneMarker.primary_tag, ...sceneMarker.tags];
    }
  }, [type, data]);

  const setInactive = () => setActive(false);
  const toggleActive = useCallback((e: TransitionEvent) => {
    if (e.propertyName === "transform" && e.elapsedTime === 0) {
      // Get the current scale of the wall-item. If it's smaller than 1.1 the item is being scaled up, otherwise down.
      const matrixScale = getComputedStyle(itemEl.current!).transform.match(
        /-?\d+\.?\d+|\d+/g
      )?.[0];
      const scale = Number.parseFloat(matrixScale ?? "2") || 2;
      setActive((value) => scale <= 1.1 && !value);
    }
  }, []);

  useEffect(() => {
    const item = itemEl.current!;
    item.addEventListener("transitioncancel", setInactive);
    item.addEventListener("transitionstart", toggleActive);
    return () => {
      item.removeEventListener("transitioncancel", setInactive);
      item.removeEventListener("transitionstart", toggleActive);
    };
  }, [toggleActive]);

  const onClick = (e: MouseEvent) => {
    clickHandler?.(e, data);
  };

  const renderText = () => {
    if (!showTextContainer) return;

    return (
      <div className="wall-item-text">
        <div>{title}</div>
        {tags?.map((tag) => (
          <span key={tag.id} className="wall-tag">
            {tag.name}
          </span>
        ))}
      </div>
    );
  };

  return (
    <div className="wall-item">
      <div className={`wall-item-container ${className}`} ref={itemEl}>
        <Link onClick={onClick} to={linkSrc} className="wall-item-anchor">
          <Preview previews={previews} config={config} active={active} />
          {renderText()}
        </Link>
      </div>
    </div>
  );
};
