import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import * as GQL from "src/core/generated-graphql";
import Gallery, {
  GalleryI,
  PhotoProps,
  RenderImageProps,
} from "react-photo-gallery";
import { ConfigurationContext } from "src/hooks/Config";
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

export const MarkerWallItem: React.FC<RenderImageProps<IMarkerPhoto>> = (
  props: RenderImageProps<IMarkerPhoto>
) => {
  const { configuration } = useContext(ConfigurationContext);
  const playSound = configuration?.interface.soundOnPreview ?? false;
  const showTitle = configuration?.interface.wallShowTitle ?? false;

  const [active, setActive] = useState(false);

  type style = Record<string, string | number | undefined>;
  var imgStyle: style = {
    margin: props.margin,
    display: "block",
  };

  if (props.direction === "column") {
    imgStyle.position = "absolute";
    imgStyle.left = props.left;
    imgStyle.top = props.top;
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
      style={{ width: props.photo.width, height: props.photo.height }}
    >
      <ImagePreview
        loading="lazy"
        loop={video}
        muted={!video || !playSound || !active}
        autoPlay={video}
        key={props.photo.key}
        style={imgStyle}
        src={props.photo.src}
        width={props.photo.width}
        height={props.photo.height}
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

const MarkerWall: React.FC<IMarkerWallProps> = ({ markers }) => {
  const history = useHistory();

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
      const { width = 1280, height = 720 } = m.scene.files[0];

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

  const renderImage = useCallback((props: RenderImageProps<IMarkerPhoto>) => {
    return <MarkerWallItem {...props} />;
  }, []);

  return (
    <div className="marker-wall">
      {photos.length ? (
        <MarkerGallery
          photos={photos}
          renderImage={renderImage}
          onClick={onClick}
          margin={margin}
          direction={direction}
          columns={columns}
        />
      ) : null}
    </div>
  );
};

interface IMarkerWallPanelProps {
  markers: GQL.SceneMarkerDataFragment[];
}

export const MarkerWallPanel: React.FC<IMarkerWallPanelProps> = ({
  markers,
}) => {
  return <MarkerWall markers={markers} />;
};
