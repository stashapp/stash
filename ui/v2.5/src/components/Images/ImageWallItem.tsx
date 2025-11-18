import React from "react";
import type { RenderImageProps } from "react-photo-gallery";

interface IExtraProps {
  maxHeight: number;
}

export const ImageWallItem: React.FC<RenderImageProps & IExtraProps> = (
  props: RenderImageProps & IExtraProps
) => {
  const height = Math.min(props.maxHeight, props.photo.height);
  const zoomFactor = height / props.photo.height;
  const width = props.photo.width * zoomFactor;

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
  const ImagePreview = video ? "video" : "img";

  return (
    <ImagePreview
      loop={video}
      muted={video}
      playsInline={video}
      autoPlay={video}
      key={props.photo.key}
      style={imgStyle}
      src={props.photo.src}
      width={width}
      height={height}
      alt={props.photo.alt}
      onClick={handleClick}
    />
  );
};
