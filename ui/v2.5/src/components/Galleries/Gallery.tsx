import React from "react";
import { useParams } from "react-router-dom";
import { StashService } from "src/core/StashService";
import { LoadingIndicator } from 'src/components/Shared';
import { GalleryViewer } from "./GalleryViewer";

export const Gallery: React.FC = () => {
  const { id = "" } = useParams();

  const { data, error, loading } = StashService.useFindGallery(id);
  const gallery = data?.findGallery;

  if (loading || !gallery)
    return <LoadingIndicator />;
  if (error) return <div>{error.message}</div>;

  return (
    <div style={{ width: "75vw", margin: "0 auto" }}>
      <GalleryViewer gallery={gallery as any} />
    </div>
  );
};
