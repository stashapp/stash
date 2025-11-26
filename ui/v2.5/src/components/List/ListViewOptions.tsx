import React, { useEffect, useRef, useState } from "react";
import Mousetrap from "mousetrap";
import {
  Button,
  ButtonGroup,
  Dropdown,
  Overlay,
  OverlayTrigger,
  Popover,
  Tooltip,
} from "react-bootstrap";
import { DisplayMode } from "src/models/list-filter/types";
import { IntlShape, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import {
  faChevronDown,
  faList,
  faSquare,
  faTags,
  faThLarge,
} from "@fortawesome/free-solid-svg-icons";
import { ZoomSelect } from "./ZoomSlider";

interface IListViewOptionsProps {
  zoomIndex?: number;
  onSetZoom?: (zoomIndex: number) => void;
  displayMode: DisplayMode;
  onSetDisplayMode: (m: DisplayMode) => void;
  displayModeOptions: DisplayMode[];
}

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

function getLabelId(option: DisplayMode) {
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
  return `display_mode.${displayModeId}`;
}

function getLabel(intl: IntlShape, option: DisplayMode) {
  return intl.formatMessage({ id: getLabelId(option) });
}

export const ListViewOptions: React.FC<IListViewOptionsProps> = ({
  zoomIndex,
  onSetZoom,
  displayMode,
  onSetDisplayMode,
  displayModeOptions,
}) => {
  const intl = useIntl();

  const overlayTarget = useRef(null);
  const [showOptions, setShowOptions] = useState(false);

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
    Mousetrap.bind("v t", () => {
      if (displayModeOptions.includes(DisplayMode.Tagger)) {
        onSetDisplayMode(DisplayMode.Tagger);
      }
    });

    return () => {
      Mousetrap.unbind("v g");
      Mousetrap.unbind("v l");
      Mousetrap.unbind("v w");
      Mousetrap.unbind("v t");
    };
  });

  function onChangeZoom(v: number) {
    if (onSetZoom) {
      onSetZoom(v);
    }
  }

  return (
    <>
      <Button
        className="display-mode-select"
        ref={overlayTarget}
        variant="secondary"
        title={intl.formatMessage(
          { id: "display_mode.label_current" },
          { current: getLabel(intl, displayMode) }
        )}
        onClick={() => setShowOptions(!showOptions)}
      >
        <Icon icon={getIcon(displayMode)} />
        <Icon size="xs" icon={faChevronDown} />
      </Button>
      <Overlay
        target={overlayTarget.current}
        show={showOptions}
        placement="bottom"
        rootClose
        onHide={() => setShowOptions(false)}
      >
        {({ placement, arrowProps, show: _show, ...props }) => (
          <div className="popover" {...props} style={{ ...props.style }}>
            <Popover.Content className="display-mode-popover">
              <div className="display-mode-menu">
                {onSetZoom &&
                zoomIndex !== undefined &&
                (displayMode === DisplayMode.Grid ||
                  displayMode === DisplayMode.Wall) ? (
                  <div className="zoom-slider-container">
                    <ZoomSelect
                      zoomIndex={zoomIndex}
                      onChangeZoom={onChangeZoom}
                    />
                  </div>
                ) : null}
                {displayModeOptions.map((option) => (
                  <Dropdown.Item
                    key={option}
                    active={displayMode === option}
                    onClick={() => {
                      setShowOptions(false);
                      onSetDisplayMode(option);
                    }}
                  >
                    <Icon icon={getIcon(option)} /> {getLabel(intl, option)}
                  </Dropdown.Item>
                ))}
              </div>
            </Popover.Content>
          </div>
        )}
      </Overlay>
    </>
  );
};

export const ListViewButtonGroup: React.FC<IListViewOptionsProps> = ({
  zoomIndex,
  onSetZoom,
  displayMode,
  onSetDisplayMode,
  displayModeOptions,
}) => {
  const intl = useIntl();

  return (
    <>
      {displayModeOptions.length > 1 && (
        <ButtonGroup>
          {displayModeOptions.map((option) => (
            <OverlayTrigger
              key={option}
              overlay={
                <Tooltip id="display-mode-tooltip">
                  {getLabel(intl, option)}
                </Tooltip>
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
      )}
      <div className="zoom-slider-container">
        {onSetZoom &&
        zoomIndex !== undefined &&
        (displayMode === DisplayMode.Grid ||
          displayMode === DisplayMode.Wall) ? (
          <ZoomSelect zoomIndex={zoomIndex} onChangeZoom={onSetZoom} />
        ) : null}
      </div>
    </>
  );
};
