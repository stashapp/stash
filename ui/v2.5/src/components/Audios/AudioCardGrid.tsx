import React from "react";
import * as GQL from "src/core/generated-graphql";
import { AudioQueue } from "src/models/audioQueue";
import { AudioCard } from "./AudioCard";
import {
  useCardWidth,
  useContainerDimensions,
} from "../Shared/GridCard/GridCard";
import { PatchComponent } from "src/patch";

interface IAudioCardGrid {
  audios: GQL.SlimAudioDataFragment[];
  queue?: AudioQueue;
  selectedIds: Set<string>;
  zoomIndex: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
}

const zoomWidths = [280, 340, 480, 640];

export const AudioCardGrid: React.FC<IAudioCardGrid> = PatchComponent(
  "AudioCardGrid",
  ({ audios, queue, selectedIds, zoomIndex, onSelectChange, fromGroupId }) => {
    const [componentRef, { width: containerWidth }] = useContainerDimensions();

    const cardWidth = useCardWidth(containerWidth, zoomIndex, zoomWidths);

    return (
      <div className="row justify-content-center" ref={componentRef}>
        {audios.map((audio, index) => (
          <AudioCard
            key={audio.id}
            width={cardWidth}
            audio={audio}
            queue={queue}
            index={index}
            zoomIndex={zoomIndex}
            selecting={selectedIds.size > 0}
            selected={selectedIds.has(audio.id)}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(audio.id, selected, shiftKey)
            }
            fromGroupId={fromGroupId}
          />
        ))}
      </div>
    );
  }
);
