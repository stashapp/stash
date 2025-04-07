import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
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
import { useIntl } from "react-intl";
import cx from "classnames";

interface IScenePhoto {
  scene: GQL.SlimSceneDataFragment;
  link: string;
  onError?: (photo: PhotoProps<IScenePhoto>) => void;
}

export const SceneWallItem: React.FC<RenderImageProps<IScenePhoto>> = (
  props: RenderImageProps<IScenePhoto>
) => {
  const intl = useIntl();

  const { configuration } = useContext(ConfigurationContext);
  const playSound = configuration?.interface.soundOnPreview ?? false;
  const showTitle = configuration?.interface.wallShowTitle ?? false;

  const [active, setActive] = useState(false);

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

  const video = props.photo.src.includes("preview");
  const ImagePreview = video ? "video" : "img";

  const { scene } = props.photo;
  const title = objectTitle(scene);
  const performerNames = scene.performers.map((p) => p.name);
  const performers =
    performerNames.length >= 2
      ? [...performerNames.slice(0, -2), performerNames.slice(-2).join(" & ")]
      : performerNames;

  return (
    <div
      className={cx("wall-item", { "show-title": showTitle })}
      role="button"
      style={{
        ...divStyle,
        width: props.photo.width,
        height: props.photo.height,
      }}
    >
      <ImagePreview
        loading="lazy"
        loop={video}
        muted={!video || !playSound || !active}
        autoPlay={video}
        key={props.photo.key}
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
            <TruncatedText text={performers.join(", ")} />
            <div>{scene.date && TextUtils.formatDate(intl, scene.date)}</div>
          </Link>
        </footer>
      </div>
    </div>
  );
};

function getDimensions(s: GQL.SlimSceneDataFragment) {
  const defaults = { width: 1280, height: 720 };

  if (!s.files.length) return defaults;

  return {
    width: s.files[0].width || defaults.width,
    height: s.files[0].height || defaults.height,
  };
}

interface ISceneWallProps {
  scenes: GQL.SlimSceneDataFragment[];
  sceneQueue?: SceneQueue;
}

// HACK: typescript doesn't allow Gallery to accept a parameter for some reason
const SceneGallery = Gallery as unknown as GalleryI<IScenePhoto>;

const SceneWall: React.FC<ISceneWallProps> = ({ scenes, sceneQueue }) => {
  const history = useHistory();

  const margin = 3;
  const direction = "row";

  const [erroredImgs, setErroredImgs] = useState<string[]>([]);

  const handleError = useCallback((photo: PhotoProps<IScenePhoto>) => {
    setErroredImgs((prev) => [...prev, photo.src]);
  }, []);

  useEffect(() => {
    setErroredImgs([]);
  }, [scenes]);

  const photos: PhotoProps<IScenePhoto>[] = useMemo(() => {
    return scenes.map((s, index) => {
      const { width, height } = getDimensions(s);

      return {
        scene: s,
        src:
          s.paths.preview && !erroredImgs.includes(s.paths.preview)
            ? s.paths.preview!
            : s.paths.screenshot!,
        link: sceneQueue
          ? sceneQueue.makeLink(s.id, { sceneIndex: index })
          : `/scenes/${s.id}`,
        width,
        height,
        tabIndex: index,
        key: s.id,
        loading: "lazy",
        alt: objectTitle(s),
        onError: handleError,
      };
    });
  }, [scenes, sceneQueue, erroredImgs, handleError]);

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

  const renderImage = useCallback((props: RenderImageProps<IScenePhoto>) => {
    return <SceneWallItem {...props} />;
  }, []);

  return (
    <div className="scene-wall">
      {photos.length ? (
        <SceneGallery
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

interface ISceneWallPanelProps {
  scenes: GQL.SlimSceneDataFragment[];
  sceneQueue?: SceneQueue;
}

export const SceneWallPanel: React.FC<ISceneWallPanelProps> = ({
  scenes,
  sceneQueue,
}) => {
  return <SceneWall scenes={scenes} sceneQueue={sceneQueue} />;
};
