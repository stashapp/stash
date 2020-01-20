import React from "react";
import { Spinner } from "react-bootstrap";
import { useParams } from "react-router-dom";
import { StashService } from "src/core/StashService";
import { GalleryViewer } from "./GalleryViewer";

export const Gallery: React.FC = () => {
  const { id = "" } = useParams();

  const { data, error, loading } = StashService.useFindGallery(id);
  const gallery = data?.findGallery;

  if (loading || !gallery)
    return <Spinner animation="border" variant="light" />;
  if (error) return <div>{error.message}</div>;

  return (
    <div style={{ width: "75vw", margin: "0 auto" }}>
      <GalleryViewer gallery={gallery as any} />
    </div>
  );
};
