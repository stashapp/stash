import React from "react";
import type {
  RenderImageProps,
  renderImageClickHandler,
  PhotoProps,
} from "react-photo-gallery";

interface IImageWallProps {
  margin?: string;
  index: number;
  photo: PhotoProps;
  onClick: renderImageClickHandler | null;
  direction: "row" | "column";
  top?: number;
  left?: number;
}

export const ImageWallItem: React.FC<RenderImageProps> = (
  props: IImageWallProps
) => {
  const clip = Boolean(
    Number(props.photo.alt?.charAt(props.photo.alt?.length - 1))
  );
  props.photo.alt = props.photo.alt?.slice(0, -1);

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

  const ImagePreview = clip ? "video" : "img";

  return (
    <ImagePreview
      loop={clip}
      autoPlay={clip}
      key={props.photo.key}
      style={imgStyle}
      src={props.photo.src}
      width={props.photo.width}
      height={props.photo.height}
      alt={props.photo.alt}
      onClick={handleClick}
    />
  );
};
