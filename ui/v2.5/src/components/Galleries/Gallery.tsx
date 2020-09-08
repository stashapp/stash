import React from "react";
import { useParams } from "react-router-dom";
import { useFindGallery } from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared";
import { GalleryViewer } from "./GalleryViewer";

interface IGalleryParams {
  id: string;
}

export const Gallery: React.FC = () => {
  const { id } = useParams<IGalleryParams>();

  const { data, error, loading } = useFindGallery(id);
  const gallery = data?.findGallery;

  if (loading || !gallery) return <LoadingIndicator />;
  if (error) return <div>{error.message}</div>;

  return (
    <div className="col col-lg-9 m-auto">
      <GalleryViewer gallery={gallery} />
    </div>
  );
};
