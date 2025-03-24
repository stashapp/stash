import React from "react";
import * as GQL from "src/core/generated-graphql";
import { IPerformerCardExtraCriteria, PerformerCard } from "./PerformerCard";
import {
  useCardWidth,
  useContainerDimensions,
} from "../Shared/GridCard/GridCard";

interface IPerformerCardGrid {
  performers: GQL.PerformerDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

const zoomWidths = [240, 300, 375, 470];

export const PerformerCardGrid: React.FC<IPerformerCardGrid> = ({
  performers,
  selectedIds,
  zoomIndex,
  onSelectChange,
  extraCriteria,
}) => {
  const [componentRef, { width: containerWidth }] = useContainerDimensions();
  const cardWidth = useCardWidth(containerWidth, zoomIndex, zoomWidths);

  return (
    <div className="row justify-content-center" ref={componentRef}>
      {performers.map((p) => (
        <PerformerCard
          key={p.id}
          cardWidth={cardWidth}
          performer={p}
          zoomIndex={zoomIndex}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(p.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(p.id, selected, shiftKey)
          }
          extraCriteria={extraCriteria}
        />
      ))}
    </div>
  );
};
