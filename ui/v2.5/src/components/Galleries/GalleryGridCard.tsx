import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryCard } from "./GalleryCard";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";

interface IGalleryCardGrid {
  galleries: GQL.SlimGalleryDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const GalleryCardGrid: React.FC<IGalleryCardGrid> = ({
  galleries,
  selectedIds,
  zoomIndex,
  onSelectChange,
}) => {
  const [componentRef, { width }] = useContainerDimensions();
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {galleries.map((gallery) => (
        <GalleryCard
          key={gallery.id}
          containerWidth={width}
          gallery={gallery}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(gallery.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(gallery.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
