import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "./GroupCard";
import {
  useCardWidth,
  useContainerDimensions,
} from "../Shared/GridCard/GridCard";

interface IGroupCardGrid {
  groups: GQL.ListGroupDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}

const zoomWidths = [210, 250, 300, 375];

export const GroupCardGrid: React.FC<IGroupCardGrid> = ({
  groups,
  selectedIds,
  zoomIndex,
  onSelectChange,
  fromGroupId,
  onMove,
}) => {
  const [componentRef, { width: containerWidth }] = useContainerDimensions();
  const cardWidth = useCardWidth(containerWidth, zoomIndex, zoomWidths);

  return (
    <div className="row justify-content-center" ref={componentRef}>
      {groups.map((p) => (
        <GroupCard
          key={p.id}
          cardWidth={cardWidth}
          group={p}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(p.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(p.id, selected, shiftKey)
          }
          fromGroupId={fromGroupId}
          onMove={onMove}
        />
      ))}
    </div>
  );
};
