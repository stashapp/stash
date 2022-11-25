import React, { useMemo } from "react";
import { useLightbox } from "src/hooks";
import { LoadingIndicator } from "src/components/Shared";
import "flexbin/flexbin.css";
import {
  CriterionModifier,
  useFindImagesQuery,
} from "src/core/generated-graphql";

interface IProps {
  galleryId: string;
}

export const GalleryViewer: React.FC<IProps> = ({ galleryId }) => {
  // TODO - add paging - don't load all images at once
  const pageSize = -1;

  const currentFilter = useMemo(() => {
    return {
      per_page: pageSize,
      sort: "path",
    };
  }, [pageSize]);

  const { data, loading } = useFindImagesQuery({
    variables: {
      filter: currentFilter,
      image_filter: {
        galleries: {
          modifier: CriterionModifier.Includes,
          value: [galleryId],
        },
      },
    },
  });

  const images = useMemo(() => data?.findImages?.images ?? [], [data]);

  const lightboxState = useMemo(() => {
    return {
      images,
      showNavigation: false,
    };
  }, [images]);

  const showLightbox = useLightbox(lightboxState);

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
