import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryCard } from "./GalleryCard";
import {
  useCardWidth,
  useContainerDimensions,
} from "../Shared/GridCard/GridCard";

interface IGalleryCardGrid {
  galleries: GQL.SlimGalleryDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const zoomWidths = [280, 340, 480, 640];

export const GalleryCardGrid: React.FC<IGalleryCardGrid> = ({
  galleries,
  selectedIds,
  zoomIndex,
  onSelectChange,
}) => {
  const [componentRef, { width: containerWidth }] = useContainerDimensions();
  const cardWidth = useCardWidth(containerWidth, zoomIndex, zoomWidths);

  return (
    <div className="row justify-content-center" ref={componentRef}>
      {galleries.map((gallery) => (
        <GalleryCard
          key={gallery.id}
          cardWidth={cardWidth}
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
