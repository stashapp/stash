import React, { useEffect } from "react";
import Mousetrap from "mousetrap";
import {
  Button,
  ButtonGroup,
  Form,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import { DisplayMode } from "src/models/list-filter/types";
import { Icon } from "../Shared";

interface IListViewOptionsProps {
  zoomIndex?: number;
  onSetZoom?: (zoomIndex: number) => void;
  displayMode: DisplayMode;
  onSetDisplayMode: (m: DisplayMode) => void;
  displayModeOptions: DisplayMode[];
}

export const ListViewOptions: React.FC<IListViewOptionsProps> = ({
  zoomIndex,
  onSetZoom,
  displayMode,
  onSetDisplayMode,
  displayModeOptions,
}) => {
  const minZoom = 0;
  const maxZoom = 3;

  useEffect(() => {
    Mousetrap.bind("v g", () => {
      if (displayModeOptions.includes(DisplayMode.Grid)) {
        onSetDisplayMode(DisplayMode.Grid);
      }
    });
    Mousetrap.bind("v l", () => {
      if (displayModeOptions.includes(DisplayMode.List)) {
        onSetDisplayMode(DisplayMode.List);
      }
    });
    Mousetrap.bind("v w", () => {
      if (displayModeOptions.includes(DisplayMode.Wall)) {
        onSetDisplayMode(DisplayMode.Wall);
      }
    });
    Mousetrap.bind("+", () => {
      if (onSetZoom && zoomIndex !== undefined && zoomIndex < maxZoom) {
        onSetZoom(zoomIndex + 1);
      }
    });
    Mousetrap.bind("-", () => {
      if (onSetZoom && zoomIndex !== undefined && zoomIndex > minZoom) {
        onSetZoom(zoomIndex - 1);
      }
    });

    return () => {
      Mousetrap.unbind("v g");
      Mousetrap.unbind("v l");
      Mousetrap.unbind("v w");
      Mousetrap.unbind("+");
      Mousetrap.unbind("-");
    };
  });

  function maybeRenderDisplayModeOptions() {
    function getIcon(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid:
          return "th-large";
        case DisplayMode.List:
          return "list";
        case DisplayMode.Wall:
          return "square";
        case DisplayMode.Tagger:
          return "tags";
      }
    }
    function getLabel(option: DisplayMode) {
      switch (option) {
        case DisplayMode.Grid:
          return "Grid";
        case DisplayMode.List:
          return "List";
        case DisplayMode.Wall:
          return "Wall";
        case DisplayMode.Tagger:
          return "Tagger";
      }
    }

    if (displayModeOptions.length < 2) {
      return;
    }

    return (
      <ButtonGroup>
        {displayModeOptions.map((option) => (
          <OverlayTrigger
            key={option}
            overlay={
              <Tooltip id="display-mode-tooltip">{getLabel(option)}</Tooltip>
            }
          >
            <Button
              variant="secondary"
              active={displayMode === option}
              onClick={() => onSetDisplayMode(option)}
            >
              <Icon icon={getIcon(option)} />
            </Button>
          </OverlayTrigger>
        ))}
      </ButtonGroup>
    );
  }

  function onChangeZoom(v: number) {
    if (onSetZoom) {
      onSetZoom(v);
    }
  }

  function maybeRenderZoom() {
    if (onSetZoom && displayMode === DisplayMode.Grid) {
      return (
        <div className="align-middle">
          <Form.Control
            className="zoom-slider d-none d-sm-inline-flex ml-3"
            type="range"
            min={minZoom}
            max={maxZoom}
            value={zoomIndex}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              onChangeZoom(Number.parseInt(e.currentTarget.value, 10))
            }
          />
        </div>
      );
    }
  }

  return (
    <>
      <ButtonGroup>{maybeRenderDisplayModeOptions()}</ButtonGroup>
      {maybeRenderZoom()}
    </>
  );
};
