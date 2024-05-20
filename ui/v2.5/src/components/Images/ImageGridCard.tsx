import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ImageCard } from "./ImageCard";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";

interface IImageCardGrid {
  images: GQL.SlimImageDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  onPreview: (index: number, ev: React.MouseEvent<Element, MouseEvent>) => void;
}

export const ImageGridCard: React.FC<IImageCardGrid> = ({
  images,
  selectedIds,
  zoomIndex,
  onSelectChange,
  onPreview,
}) => {
  const [componentRef, { width }] = useContainerDimensions();
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {images.map((image, index) => (
        <ImageCard
          key={image.id}
          containerWidth={width}
          image={image}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(image.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(image.id, selected, shiftKey)
          }
          onPreview={
            selectedIds.size < 1 ? (ev) => onPreview(index, ev) : undefined
          }
        />
      ))}
    </div>
  );
};
