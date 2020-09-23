import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import FsLightbox from "fslightbox-react";
import "flexbin/flexbin.css";

interface IProps {
  gallery: Partial<GQL.GalleryDataFragment>;
}

export const GalleryViewer: React.FC<IProps> = ({ gallery }) => {
  const [lightboxToggle, setLightboxToggle] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);

  const openImage = (index: number) => {
    setCurrentIndex(index);
    setLightboxToggle(!lightboxToggle);
  };

  const photos = !gallery.images ? [] : gallery.images.map((file) => file.paths.image ?? "");
  const thumbs = !gallery.images ? [] : gallery.images.map((file, index) => (
    <div
      role="link"
      tabIndex={index}
      key={file.checksum ?? index}
      onClick={() => openImage(index)}
      onKeyPress={() => openImage(index)}
    >
      <img
        src={file.paths.thumbnail ?? ''}
        loading="lazy"
        className="gallery-image"
        alt={file.title ?? index.toString()}
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
        key={gallery.id!}
      />
    </div>
  );
};
