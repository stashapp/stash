import React from "react";
import * as GQL from "src/core/generated-graphql";
import { StudioCard } from "./StudioCard";

interface IStudioCardGrid {
  studios: GQL.StudioDataFragment[];
  fromParent: boolean | undefined;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const StudioCardGrid: React.FC<IStudioCardGrid> = ({
  studios,
  fromParent,
  selectedIds,
  onSelectChange,
}) => {
  return (
    <div className="row justify-content-center">
      {studios.map((studio) => (
        <StudioCard
          key={studio.id}
          studio={studio}
          hideParent={fromParent}
          selecting={selectedIds.size > 0}
          selected={selectedIds.has(studio.id)}
          onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
            onSelectChange(studio.id, selected, shiftKey)
          }
        />
      ))}
    </div>
  );
};
