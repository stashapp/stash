import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useLightbox } from "src/hooks";
import "flexbin/flexbin.css";

interface IProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryViewer: React.FC<IProps> = ({ gallery }) => {
  const images = gallery?.images ?? [];
  const showLightbox = useLightbox({ images, showNavigation: false });

  const thumbs = images.map((file, index) => (
    <div
      role="link"
      tabIndex={index}
      key={file.checksum ?? index}
      onClick={() => showLightbox(index)}
      onKeyPress={() => showLightbox(index)}
    >
      <img
        src={file.paths.thumbnail ?? ""}
        loading="lazy"
        className="gallery-image"
        alt={file.title ?? index.toString()}
      />
    </div>
  ));

  return (
    <div className="gallery">
      <div className="flexbin">{thumbs}</div>
    </div>
  );
};
