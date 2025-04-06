import React, { MouseEvent, useCallback, useContext, useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { WallItem, WallItemData, WallItemType } from "./WallItem";
import Gallery, { PhotoProps, renderImageClickHandler, RenderImageProps } from "react-photo-gallery";
import { ConfigurationContext } from "src/hooks/Config";
import { objectTitle } from "src/core/files";

interface IWallPanelProps<T extends WallItemType> {
  type: T;
  data: WallItemData[T][];
  sceneQueue?: SceneQueue;
  clickHandler?: (e: MouseEvent, item: WallItemData[T]) => void;
}

const calculateClass = (index: number, count: number) => {
  // First position and more than one row
  if (index === 0 && count > 5) return "transform-origin-top-left";
  // Fifth position and more than one row
  if (index === 4 && count > 5) return "transform-origin-top-right";
  // Top row
  if (index < 5) return "transform-origin-top";
  // Two or more rows, with full last row and index is last
  if (count > 9 && count % 5 === 0 && index + 1 === count)
    return "transform-origin-bottom-right";
  // Two or more rows, with full last row and index is fifth to last
  if (count > 9 && count % 5 === 0 && index + 5 === count)
    return "transform-origin-bottom-left";
  // Multiple of five minus one
  if (index % 5 === 4) return "transform-origin-right";
  // Multiple of five
  if (index % 5 === 0) return "transform-origin-left";
  // Position is equal or larger than first position in last row
  if (count - (count % 5 || 5) <= index + 1) return "transform-origin-bottom";
  // Default
  return "transform-origin-center";
};

const WallPanel = <T extends WallItemType>({
  type,
  data,
  sceneQueue,
  clickHandler,
}: IWallPanelProps<T>) => {
  function renderItems() {
    return data.map((item, index, arr) => (
      <WallItem
        type={type}
        key={item.id}
        index={index}
        data={item}
        sceneQueue={sceneQueue}
        clickHandler={clickHandler}
        className={calculateClass(index, arr.length)}
      />
    ));
  }

  return (
    <div className="row">
      <div className="wall w-100 row justify-content-center">
        {renderItems()}
      </div>
    </div>
  );
};

interface IMarkerWallPanelProps {
  markers: GQL.SceneMarkerDataFragment[];
  clickHandler?: (e: MouseEvent, item: GQL.SceneMarkerDataFragment) => void;
}

export const MarkerWallPanel: React.FC<IMarkerWallPanelProps> = ({
  markers,
  clickHandler,
}) => {
  return (
    <WallPanel type="sceneMarker" data={markers} clickHandler={clickHandler} />
  );
};

interface ISceneWallPanelProps {
  scenes: GQL.SlimSceneDataFragment[];
  sceneQueue?: SceneQueue;
  clickHandler?: (e: MouseEvent, item: GQL.SlimSceneDataFragment) => void;
}

export const SceneWallPanel: React.FC<ISceneWallPanelProps> = ({
  scenes,
  sceneQueue,
  clickHandler,
}) => {
  return (
    // <WallPanel
    //   type="scene"
    //   data={scenes}
    //   sceneQueue={sceneQueue}
    //   clickHandler={clickHandler}
    // />
    <SceneWall
      scenes={scenes}
      handleImageOpen={() => {}}
    />
  );
};

interface IBackupSrc {
  backupSrc?: string;
  onError?: (photo: PhotoProps<IBackupSrc>) => void;
}

interface ISceneWallItemProps {
  margin?: string;
  index: number;
  photo: PhotoProps<IBackupSrc>;
  onClick: renderImageClickHandler | null;
  direction: "row" | "column";
  top?: number;
  left?: number;
}

export const SceneWallItem: React.FC<RenderImageProps> = (
  props: ISceneWallItemProps
) => {
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

  var handleClick = function handleClick(
    event: React.MouseEvent<Element, MouseEvent>
  ) {
    if (props.onClick) {
      props.onClick(event, { index: props.index });
    }
  };

  const video = props.photo.src.includes("preview");
  const VideoPreview = "video";
  const ImagePreview = video ? "video" : "img";

  const previewProps = {
    loop: video,
    muted: video,
    autoPlay: video,
    key: props.photo.key,
    style: imgStyle,
    src: props.photo.src,
    width: props.photo.width,
    height: props.photo.height,
    alt: props.photo.alt,
    onClick: handleClick,
  };

  // if (isMissing) {
  //   // show the image if the video preview is unavailable
  //   if (props.photo.backupSrc) {
  //     return <img {...previewProps} src={props.photo.backupSrc} />;
  //   }

  //   return (
  //     <div className="wall-item-media wall-item-missing">
  //       Pending preview generation
  //     </div>
  //   );
  // }

  return (
    <ImagePreview
      loop={video}
      muted={video}
      autoPlay={video}
      key={props.photo.key}
      style={imgStyle}
      src={props.photo.src}
      width={props.photo.width}
      height={props.photo.height}
      alt={props.photo.alt}
      onClick={handleClick}
      onError={() => {
        props.photo.onError?.(props.photo);
      }}
    />
  );
};

interface ISceneWallProps {
  scenes: GQL.SlimSceneDataFragment[];
  handleImageOpen: (index: number) => void;
}

const SceneWall: React.FC<ISceneWallProps> = ({ scenes, handleImageOpen }) => {
  const { configuration } = useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;

  const [erroredImgs, setErroredImgs] = useState<string[]>([]);

  const handleError = useCallback(
    (photo: PhotoProps<IBackupSrc>) => {
      setErroredImgs((prev) => [...prev, photo.src]);
    },
    []
  );

  useEffect(() => {
    setErroredImgs([]);
  }, [scenes]);

  const photos: PhotoProps<IBackupSrc>[] = useMemo(() => {
    return scenes.filter((s) => {
      return s.files.length > 0 && s.files[0].width > 0 && s.files[0].height > 0;
    }).map((s, index) => {
      const width = s.files[0].width;
      const height = s.files[0].height;

      return {
        src:
          s.paths.preview && !erroredImgs.includes(s.paths.preview)
            ? s.paths.preview!
            : s.paths.screenshot!,
        width,
        height,
        tabIndex: index,
        key: s.id,
        loading: "lazy",
        className: "gallery-image",
        alt: objectTitle(s),
        onError: handleError,
      };
    });
  }, [scenes, erroredImgs, handleError]);

  const showLightboxOnClick = useCallback(
    (event, { index }) => {
      // handleImageOpen(index);
    },
    [handleImageOpen]
  );

  function columns(containerWidth: number) {
    let preferredSize = 300;
    let columnCount = containerWidth / preferredSize;
    return Math.round(columnCount);
  }

  return (
    <div className="gallery">
      {photos.length ? (
        <Gallery
          photos={photos}
          renderImage={SceneWallItem}
          onClick={showLightboxOnClick}
          margin={uiConfig?.imageWallOptions?.margin!}
          direction={uiConfig?.imageWallOptions?.direction!}
          columns={columns}
        />
      ) : null}
    </div>
  );
};