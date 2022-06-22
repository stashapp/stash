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
import { useIntl } from "react-intl";
import { Icon } from "../Shared";
import {
  faList,
  faSquare,
  faTags,
  faThLarge,
} from "@fortawesome/free-solid-svg-icons";

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

  const intl = useIntl();

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
          return faThLarge;
        case DisplayMode.List:
          return faList;
        case DisplayMode.Wall:
          return faSquare;
        case DisplayMode.Tagger:
          return faTags;
      }
    }
    function getLabel(option: DisplayMode) {
      let displayModeId = "unknown";
      switch (option) {
        case DisplayMode.Grid:
          displayModeId = "grid";
          break;
        case DisplayMode.List:
          displayModeId = "list";
          break;
        case DisplayMode.Wall:
          displayModeId = "wall";
          break;
        case DisplayMode.Tagger:
          displayModeId = "tagger";
          break;
      }
      return intl.formatMessage({ id: `display_mode.${displayModeId}` });
    }

    if (displayModeOptions.length < 2) {
      return;
    }

    return (
      <ButtonGroup className="mb-2">
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
        <div className="ml-2 mb-2 d-none d-sm-inline-flex">
          <Form.Control
            className="zoom-slider ml-1"
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
      {maybeRenderDisplayModeOptions()}
      {maybeRenderZoom()}
    </>
  );
};
