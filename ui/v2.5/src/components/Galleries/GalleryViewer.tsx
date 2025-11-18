import React, { useCallback, useMemo } from "react";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import Gallery, { PhotoClickHandler } from "react-photo-gallery";
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
  const showLightboxOnClick: PhotoClickHandler = useCallback(
    (event, { index }) => {
      showLightbox({ initialIndex: index });
    },
    [showLightbox]
  );

  if (loading) return <LoadingIndicator />;

  let photos: {
    src: string;
    srcSet?: string | string[] | undefined;
    sizes?: string | string[] | undefined;
    width: number;
    height: number;
    alt?: string | undefined;
    key?: string | undefined;
  }[] = [];

  images.forEach((image, index) => {
    let imageData = {
      src: image.paths.thumbnail!,
      width: image.visual_files[0]?.width ?? 0,
      height: image.visual_files[0]?.height ?? 0,
      tabIndex: index,
      key: image.id ?? index,
      loading: "lazy",
      className: "gallery-image",
      alt: image.title ?? index.toString(),
    };
    photos.push(imageData);
  });

  return (
    <div className="gallery">
      <Gallery photos={photos} onClick={showLightboxOnClick} margin={2.5} />
    </div>
  );
};

export default GalleryViewer;
