import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import {
  Button,
  ButtonGroup,
  Form,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import { ConfigMode } from "src/models/list-filter/types";
import { useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import {
  faEye,
  faEyeSlash,
  faTree,
  faCircleCheck,
  faSitemap,
  faIoxhost,
  faLayerGroup,
  faTags,
  faVolumeXmark,
  faVolumeHigh,
} from "@fortawesome/free-solid-svg-icons";
import { useConfiguration, useConfigureUI } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

interface Page {
  page: string;
}
interface IButtonItem {
  name: string;
  pages: Page[];
  tooltipDisabled: string;
  tooltipEnabled: string;
  iconDisabled: IconDefinition;
  iconEnabled: IconDefinition;
  mode: string;
  configString: [key: string];
}

const allButtonItems: IButtonItem[] = [
  {
    name: "audio",
    pages: [
      { page: "scenes" },
      { page: "markers" },
    ],
    tooltipDisabled: "studios-disabled",
    tooltipEnabled: "studios-enabled",
    iconDisabled: faVolumeXmark,
    iconEnabled: faVolumeHigh,
    mode: "audio",
    configString: "config",
  },
  {
    name: "studios",
    pages: [{ page: "studios/" }],
    tooltipDisabled: "studios-disabled",
    tooltipEnabled: "studios-enabled",
    iconDisabled: faSitemap,
    iconEnabled: faSitemap,
    mode: "studios",
    configString: ["config.data?.configuration?.ui?.showChildStudioContent"],
  },
  {
    name: "tags",
    pages: [{ page: "tags/" }],
    tooltipDisabled: "tags-disabled",
    tooltipEnabled: "tags-enabled",
    iconDisabled: faSitemap,
    iconEnabled: faSitemap,
    mode: "tags",
    configString: "config",
  },
  {
    name: "hover",
    pages: [
      { page: "scenes" },
      { page: "images" },
      { page: "galleries" },
      { page: "performers" },
      { page: "tags" },
    ],
    tooltipDisabled: "tags-hover-disabled",
    tooltipEnabled: "tags-hover-enabled",
    iconDisabled: faLayerGroup,
    iconEnabled: faLayerGroup,
    mode: "tagsHover",
    configString: "config",
  },
];

interface IListToggleConfigSettingsProps {
  activePage: string;
  configMode: ConfigMode;
  settings: {
    showChildStudioContent: boolean | undefined;
    showChildTagContent: boolean | undefined;
    showTagCardOnHover: boolean | undefined;
  };
  onSetConfigMode: (m: ConfigMode) => void;
  configModeOptions: ConfigMode[];
}

export const ListToggleConfigButtons: React.FC<
  IListToggleConfigSettingsProps
> = ({
  activePage,
  configMode,
  settings,
  onSetConfigMode,
  configModeOptions,
}) => {
  const intl = useIntl();
  const { configuration, loading } = React.useContext(ConfigurationContext);
  const [buttonItems, setButtonItems] = useState<IButtonItem[]>(allButtonItems);
  const location = useLocation();
  const pathStudios = location.pathname.includes("studios");
  const pathTags = location.pathname.includes("tags");

  useEffect(() => {}, [configuration]);

  //const [activePage, setActivePage] = useState<string>();
  const config = useConfiguration();
  const [saveUI] = useConfigureUI();

  const [childActive, setChildActive] = useState<boolean>(false);
  const [toggleAudio, setToggleAudio] = useState<boolean>(false);
  const [toggleChildStudios, setToggleChildStudios] = useState<boolean>(false);
  const [toggleChildTags, setToggleChildTags] = useState<boolean>(false);
  const [toggleTagHoverActive, setToggleTagHoverActive] = useState<boolean>(false);

  const audio = config.data?.configuration?.interface?.soundOnPreview;
  const childStudio = config.data?.configuration?.ui?.showChildStudioContent;

  const childTag = config.data?.configuration?.ui?.showChildTagContent;

  const tagHover = config.data?.configuration?.ui?.showTagCardOnHover;

  useEffect(() => {
      if (audio) {
        
      }
      if (childStudio) {
        setChildActive(true);
      }
      //   setActivePage("tags")
      if (childTag) {
        setChildActive(true);
    }
  }, [childStudio, childTag]);

  function oTs(mode: number, updatedConfig: string) {
    switch (mode) {
      case 1:
        updatedConfig.ui.showTagCardOnHover = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      case 2:
        updatedConfig.ui.showTagCardOnHover = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      default:
        break;
    }
  }
  function onSetToggleChildren(mode: string) {
    // Clone the config object to avoid mutating the original
    const updatedConfig = { ...config.data?.configuration.ui };

    const updatedConfigs = { ...config.data?.configuration.interface };

    switch (mode) {
      case "audio":
        updatedConfigs.soundOnPreview = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      case "studios":
        alert("dd")
        updatedConfig.showChildStudioContent = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      case "tags":
        updatedConfig.ui.showChildTagContent = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      case "tagsHover":
        updatedConfig.ui.showTagCardOnHover = !childActive;
        setChildActive(!childActive);
        onSetConfigMode();
        break;
      default:
        break;
    }
    // Now, save the updated config object
    saveUI({
      variables: {
        input: {
          ...config.data?.configuration.ui,
          ...updatedConfig,
        },
      },
    });
  }
  const ss = JSON.stringify(config.data);
  function maybeRenderChildButtons(mode: string) {
    let setMode: string = "";
    let childToolTip: string = "";

    switch (mode) {
      case "studios":
        setMode = "studios";
        childToolTip = "Toggle display of child studios";
        break;
      case "tags":
        setMode = "tags";
        childToolTip = "Toggle display child tags";
        break;
      default:
        // Handle the default case
        setMode = "";
        childToolTip = "";
    }

    function evaluateVariable(key: [key: string]) {
      if (config[key] !== undefined && typeof config[key] === 'boolean') {
        return config[key];
      } else {
        return false; // Return false for any other cases
      }
    }

    // Filter the buttonItems to include only those with matching pages
    const matchingPages = buttonItems.filter((item) =>
      item.pages.some((page) => location.pathname.includes(page.page))
    );

    if (matchingPages.length > 0) {
      return (
        <>
          <ButtonGroup className="mb-2 btn-group-actions">
            {matchingPages.map((matchingPage) => (
              <OverlayTrigger
                key={"showChildren"} // This key should be unique for each OverlayTrigger
                overlay={
                  <Tooltip id="showChildren">
                    {intl.formatMessage({
                      id: `config_mode.${matchingPage.tooltipEnabled}`,
                    })}
                  </Tooltip>
                }
              >
                <Button
                  variant="secondary"
                  active={evaluateVariable(matchingPage.configString)}
                  onClick={() => onSetToggleChildren(matchingPage.mode)}
                >
                  {childActive ? (
                    <Icon icon={matchingPage.iconEnabled} color="lime" />
                  ) : (
                    <Icon icon={matchingPage.iconEnabled} />
                  )}
                </Button>
              </OverlayTrigger>
            ))}
          </ButtonGroup>
        </>
      );
    } else {
      return null;
    }
  }

  function maybeRenderConfigModeOptions() {
    function getIcon(option: ConfigMode) {
      switch (option) {
        case ConfigMode.Studios:
          return faTree;
        case ConfigMode.Tags:
          return faTree;
        case ConfigMode.Hover:
          return faTags;
      }
    }

    function getLabel(option: string) {
      let configModeId = "unknown";
      switch (option) {
        case "studios":
          configModeId = "studios";
          break;
        case "tags":
          configModeId = "tags";
          break;
        case "hover":
          configModeId = "hover";
          break;
      }
      return intl.formatMessage({ id: `config_mode.${configModeId}` });
    }

    if (configModeOptions.length < 2) {
      return;
    }

    return (
      <ButtonGroup className="mb-2">
        {configModeOptions.map((option) => (
          <OverlayTrigger
            key={option}
            overlay={
              <Tooltip id="config-mode-tooltip">{getLabel(option)}</Tooltip>
            }
          >
            <Button
              variant="secondary"
              active={configMode === option}
              onClick={() => onSetConfigMode(option)}
            >
              <Icon icon={getIcon(option)} />
              <i style={{ color: "#33d21e;" }}>
                <Icon icon={faCircleCheck} />
              </i>
            </Button>
          </OverlayTrigger>
        ))}
      </ButtonGroup>
    );
  }

  return <>{maybeRenderChildButtons(activePage || "")}</>;
};
