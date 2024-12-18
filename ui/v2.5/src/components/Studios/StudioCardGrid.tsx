import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useContainerDimensions } from "../Shared/GridCard/GridCard";
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
  const [componentRef, { width }] = useContainerDimensions();
  return (
    <div className="row justify-content-center" ref={componentRef}>
      {studios.map((studio) => (
        <StudioCard
          key={studio.id}
          containerWidth={width}
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
