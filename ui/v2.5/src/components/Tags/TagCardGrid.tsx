import React from "react";
import * as GQL from "src/core/generated-graphql";
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
  return (
    <div className="row justify-content-center">
      {tags.map((tag) => (
        <TagCard
          key={tag.id}
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
