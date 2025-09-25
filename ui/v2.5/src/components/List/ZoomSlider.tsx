import React, { useEffect } from "react";
import Mousetrap from "mousetrap";
import { Form } from "react-bootstrap";

export interface IZoomSelectProps {
  minZoom: number;
  maxZoom: number;
  zoomIndex: number;
  onChangeZoom: (v: number) => void;
}

export const ZoomSelect: React.FC<IZoomSelectProps> = ({
  minZoom,
  maxZoom,
  zoomIndex,
  onChangeZoom,
}) => {
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
