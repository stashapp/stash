import React, { useCallback, useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import Gallery, {
  GalleryI,
  PhotoProps,
  RenderImageProps,
} from "react-photo-gallery";
import { useConfigurationContext } from "src/hooks/Config";
import { objectTitle } from "src/core/files";
import { Link, useHistory } from "react-router-dom";
import { TruncatedText } from "../Shared/TruncatedText";
import TextUtils from "src/utils/text";
import cx from "classnames";
import NavUtils from "src/utils/navigation";
import { markerTitle } from "src/core/markers";

function wallItemTitle(sceneMarker: GQL.SceneMarkerDataFragment) {
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
}

interface IMarkerPhoto {
  marker: GQL.SceneMarkerDataFragment;
  link: string;
  onError?: (photo: PhotoProps<IMarkerPhoto>) => void;
}

interface IExtraProps {
  maxHeight: number;
}

export const MarkerWallItem: React.FC<
  RenderImageProps<IMarkerPhoto> & IExtraProps
> = (props: RenderImageProps<IMarkerPhoto> & IExtraProps) => {
  const { configuration } = useConfigurationContext();
  const playSound = configuration?.interface.soundOnPreview ?? false;
  const showTitle = configuration?.interface.wallShowTitle ?? false;

  const [active, setActive] = useState(false);

  const height = Math.min(props.maxHeight, props.photo.height);
  const zoomFactor = height / props.photo.height;
  const width = props.photo.width * zoomFactor;

  type style = Record<string, string | number | undefined>;
  var divStyle: style = {
    margin: props.margin,
    display: "block",
  };

  if (props.direction === "column") {
    divStyle.position = "absolute";
    divStyle.left = props.left;
    divStyle.top = props.top;
  }

  var handleClick = function handleClick(event: React.MouseEvent) {
    if (props.onClick) {
      props.onClick(event, { index: props.index });
    }
  };

  const video = props.photo.src.includes("stream");
  const ImagePreview = video ? "video" : "img";

  const { marker } = props.photo;
  const title = wallItemTitle(marker);
  const tagNames = marker.tags.map((p) => p.name);

  return (
    <div
      className={cx("wall-item", { "show-title": showTitle })}
      role="button"
      style={{
        ...divStyle,
        width,
        height,
      }}
    >
      <ImagePreview
        loading="lazy"
        loop={video}
        muted={!video || !playSound || !active}
        autoPlay={video}
        playsInline={video}
        key={props.photo.key}
        src={props.photo.src}
        width={width}
        height={height}
        alt={props.photo.alt}
        onMouseEnter={() => setActive(true)}
        onMouseLeave={() => setActive(false)}
        onClick={handleClick}
        onError={() => {
          props.photo.onError?.(props.photo);
        }}
      />
      <div className="lineargradient">
        <footer className="wall-item-footer">
          <Link to={props.photo.link} onClick={(e) => e.stopPropagation()}>
            {title && (
              <TruncatedText
                text={title}
                lineCount={1}
                className="wall-item-title"
              />
            )}
            <TruncatedText text={tagNames.join(", ")} />
          </Link>
        </footer>
      </div>
    </div>
  );
};

interface IMarkerWallProps {
  markers: GQL.SceneMarkerDataFragment[];
  zoomIndex: number;
}

// HACK: typescript doesn't allow Gallery to accept a parameter for some reason
const MarkerGallery = Gallery as unknown as GalleryI<IMarkerPhoto>;

function getFirstValidSrc(srcSet: string[], invalidSrcSet: string[]) {
  if (!srcSet.length) {
    return "";
  }

  return (
    srcSet.find((src) => !invalidSrcSet.includes(src)) ??
    ([...srcSet].pop() as string)
  );
}

interface IFile {
  width: number;
  height: number;
}

function getDimensions(file?: IFile) {
  const defaults = { width: 1280, height: 720 };

  if (!file) return defaults;

  return {
    width: file.width || defaults.width,
    height: file.height || defaults.height,
  };
}

const breakpointZoomHeights = [
  { minWidth: 576, heights: [100, 120, 240, 360] },
  { minWidth: 768, heights: [120, 160, 240, 480] },
  { minWidth: 1200, heights: [120, 160, 240, 300] },
  { minWidth: 1400, heights: [160, 240, 300, 480] },
];

const MarkerWall: React.FC<IMarkerWallProps> = ({ markers, zoomIndex }) => {
  const history = useHistory();

  const containerRef = React.useRef<HTMLDivElement>(null);

  const margin = 3;
  const direction = "row";

  const [erroredImgs, setErroredImgs] = useState<string[]>([]);

  const handleError = useCallback((photo: PhotoProps<IMarkerPhoto>) => {
    setErroredImgs((prev) => [...prev, photo.src]);
  }, []);

  useEffect(() => {
    setErroredImgs([]);
  }, [markers]);

  const photos: PhotoProps<IMarkerPhoto>[] = useMemo(() => {
    return markers.map((m, index) => {
      const { width = 1280, height = 720 } = getDimensions(m.scene.files[0]);

      return {
        marker: m,
        src: getFirstValidSrc([m.stream, m.preview, m.screenshot], erroredImgs),
        link: NavUtils.makeSceneMarkerUrl(m),
        width,
        height,
        tabIndex: index,
        key: m.id,
        loading: "lazy",
        alt: objectTitle(m),
        onError: handleError,
      };
    });
  }, [markers, erroredImgs, handleError]);

  const onClick = useCallback(
    (event, { index }) => {
      history.push(photos[index].link);
    },
    [history, photos]
  );

  function columns(containerWidth: number) {
    let preferredSize = 300;
    let columnCount = containerWidth / preferredSize;
    return Math.round(columnCount);
  }

  const targetRowHeight = useCallback(
    (containerWidth: number) => {
      let zoomHeight = 280;
      breakpointZoomHeights.forEach((e) => {
        if (containerWidth >= e.minWidth) {
          zoomHeight = e.heights[zoomIndex];
        }
      });
      return zoomHeight;
    },
    [zoomIndex]
  );

  // set the max height as a factor of the targetRowHeight
  // this allows some images to be taller than the target row height
  // but prevents images from becoming too tall when there is a small number of items
  const maxHeightFactor = 1.3;

  const renderImage = useCallback(
    (props: RenderImageProps<IMarkerPhoto>) => {
      return (
        <MarkerWallItem
          {...props}
          maxHeight={
            targetRowHeight(containerRef.current?.offsetWidth ?? 0) *
            maxHeightFactor
          }
        />
      );
    },
    [targetRowHeight]
  );

  return (
    <div className="marker-wall" ref={containerRef}>
      {photos.length ? (
        <MarkerGallery
          photos={photos}
          renderImage={renderImage}
          onClick={onClick}
          margin={margin}
          direction={direction}
          columns={columns}
          targetRowHeight={targetRowHeight}
        />
      ) : null}
    </div>
  );
};

interface IMarkerWallPanelProps {
  markers: GQL.SceneMarkerDataFragment[];
  zoomIndex: number;
}

export const MarkerWallPanel: React.FC<IMarkerWallPanelProps> = ({
  markers,
  zoomIndex,
}) => {
  return <MarkerWall markers={markers} zoomIndex={zoomIndex} />;
};
