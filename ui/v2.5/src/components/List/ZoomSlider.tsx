import React, { useEffect } from "react";
import Mousetrap from "mousetrap";
import { Form } from "react-bootstrap";

const minZoom = 0;
const maxZoom = 3;

export function useZoomKeybinds(props: {
  zoomIndex: number | undefined;
  onChangeZoom: (v: number) => void;
}) {
  const { zoomIndex, onChangeZoom } = props;
  useEffect(() => {
    Mousetrap.bind("+", () => {
      if (zoomIndex !== undefined && zoomIndex < maxZoom) {
        onChangeZoom(zoomIndex + 1);
      }
    });
    Mousetrap.bind("-", () => {
      if (zoomIndex !== undefined && zoomIndex > minZoom) {
        onChangeZoom(zoomIndex - 1);
      }
    });

    return () => {
      Mousetrap.unbind("+");
      Mousetrap.unbind("-");
    };
  });
}

export interface IZoomSelectProps {
  zoomIndex: number;
  onChangeZoom: (v: number) => void;
}

export const ZoomSelect: React.FC<IZoomSelectProps> = ({
  zoomIndex,
  onChangeZoom,
}) => {
  return (
    <Form.Control
      className="zoom-slider"
      type="range"
      min={minZoom}
      max={maxZoom}
      value={zoomIndex}
      onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
        onChangeZoom(Number.parseInt(e.currentTarget.value, 10));
        e.preventDefault();
        e.stopPropagation();
      }}
    />
  );
};
