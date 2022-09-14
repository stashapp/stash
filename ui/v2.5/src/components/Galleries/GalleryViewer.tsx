import React from "react";
import { useFindGallery } from "src/core/StashService";
import { useLightbox } from "src/hooks";
import { LoadingIndicator } from "src/components/Shared";
import "flexbin/flexbin.css";

interface IProps {
  galleryId: string;
}

export const GalleryViewer: React.FC<IProps> = ({ galleryId }) => {
  const { data, loading } = useFindGallery(galleryId);
  const images = data?.findGallery?.images ?? [];
  const showLightbox = useLightbox({ images, showNavigation: false });

  if (loading) return <LoadingIndicator />;

  const thumbs = images.map((file, index) => (
    <div
      role="link"
      tabIndex={index}
      key={file.id ?? index}
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

export default GalleryViewer;
