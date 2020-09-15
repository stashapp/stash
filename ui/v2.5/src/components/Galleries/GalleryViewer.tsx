import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import FsLightbox from "fslightbox-react";
import "flexbin/flexbin.css";

interface IProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryViewer: React.FC<IProps> = ({ gallery }) => {
  const [lightboxToggle, setLightboxToggle] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);

  const openImage = (index: number) => {
    setCurrentIndex(index);
    setLightboxToggle(!lightboxToggle);
  };

  const photos = gallery.files.map((file) => file.path ?? "");
  const thumbs = gallery.files.map((file, index) => (
    <div
      role="link"
      tabIndex={index}
      key={file.index}
      onClick={() => openImage(index)}
      onKeyPress={() => openImage(index)}
    >
      <img
        src={`${file.path}?thumb=600` || ""}
        loading="lazy"
        className="gallery-image"
        alt={file.name ?? index.toString()}
      />
    </div>
  ));

  return (
    <div className="gallery">
      <div className="flexbin">{thumbs}</div>
      <FsLightbox
        sourceIndex={currentIndex}
        toggler={lightboxToggle}
        sources={photos}
      />
    </div>
  );
};
