import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "./MovieCard";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";

interface IGroupCardGrid {
  groups: GQL.MovieDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const GroupCardGrid: React.FC<IGroupCardGrid> = ({
  groups,
  selectedIds,
  onSelectChange,
}) => {
  const [componentRef, { width }] = useContainerDimensions();
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {groups.map((p) => (
        <GroupCard
          key={p.id}
          containerWidth={width}
          group={p}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(p.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(p.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
