import React, { useRef } from "react";
import * as GQL from "src/core/generated-graphql";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";
import { TagCard } from "./TagCard";

interface ITagCardGrid {
  tags: GQL.TagDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const TagCardGrid: React.FC<ITagCardGrid> = ({
  tags,
  selectedIds,
  zoomIndex,
  onSelectChange,
}) => {
  const componentRef = useRef<HTMLDivElement>(null);
  const { width } = useContainerDimensions(componentRef);
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {tags.map((tag) => (
        <TagCard
          key={tag.id}
          containerWidth={width}
          tag={tag}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(tag.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(tag.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
