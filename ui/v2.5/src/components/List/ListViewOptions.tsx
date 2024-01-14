import React, { useEffect } from "react";
import Mousetrap from "mousetrap";
import { Dropdown, Form } from "react-bootstrap";
import { DisplayMode } from "src/models/list-filter/types";
import { useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import {
  faChevronDown,
  faList,
  faSquare,
  faTags,
  faThLarge,
} from "@fortawesome/free-solid-svg-icons";

export const ZoomSelect: React.FC<{
  minZoom: number;
  maxZoom: number;
  zoomIndex: number;
  onChangeZoom: (v: number) => void;
}> = ({ minZoom, maxZoom, zoomIndex, onChangeZoom }) => {
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
      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
        onChangeZoom(Number.parseInt(e.currentTarget.value, 10))
      }
    />
  );
};

interface IDisplayModeSelectProps {
  displayMode: DisplayMode;
  onSetDisplayMode: (m: DisplayMode) => void;
  displayModeOptions: DisplayMode[];
}

export const DisplayModeSelect: React.FC<IDisplayModeSelectProps> = ({
  displayMode,
  onSetDisplayMode,
  displayModeOptions,
}) => {
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

    return () => {
      Mousetrap.unbind("v g");
      Mousetrap.unbind("v l");
      Mousetrap.unbind("v w");
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
      <Dropdown className="display-mode-select">
        <Dropdown.Toggle variant="secondary" title={getLabel(displayMode)}>
          <Icon icon={getIcon(displayMode)} />
          <Icon size="xs" icon={faChevronDown} />
        </Dropdown.Toggle>
        <Dropdown.Menu>
          {displayModeOptions.map((option) => (
            <Dropdown.Item
              key={option}
              active={displayMode === option}
              onClick={() => onSetDisplayMode(option)}
            >
              <Icon icon={getIcon(option)} /> {getLabel(option)}
            </Dropdown.Item>
          ))}
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  return <>{maybeRenderDisplayModeOptions()}</>;
};
