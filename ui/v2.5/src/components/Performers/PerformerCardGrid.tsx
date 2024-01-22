import React from "react";
import * as GQL from "src/core/generated-graphql";
import { IPerformerCardExtraCriteria, PerformerCard } from "./PerformerCard";

interface IPerformerCardGrid {
  performers: GQL.PerformerDataFragment[];
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const PerformerCardGrid: React.FC<IPerformerCardGrid> = ({
  performers,
  selectedIds,
  onSelectChange,
  extraCriteria,
}) => {
  return (
    <div className="row justify-content-center">
      {performers.map((performer) => (
        <PerformerCard
          key={performer.id}
          performer={performer}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(performer.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(performer.id, selected, shiftKey)
          }
          extraCriteria={extraCriteria}
        />
      ))}
    </div>
  );
};
