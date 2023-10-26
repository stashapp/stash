import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { Button, ButtonGroup, OverlayTrigger, Tooltip } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import {
  faSitemap,
  faLayerGroup,
  faVolumeXmark,
  faVolumeHigh,
} from "@fortawesome/free-solid-svg-icons";
import {
  useConfiguration,
  useConfigureInterface,
  useConfigureUI,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

interface IPage {
  page: string;
}
interface IButtonItem {
  name: string;
  pages: IPage[];
  tooltipDisabled: string;
  tooltipEnabled: string;
  iconDisabled: IconDefinition;
  iconEnabled: IconDefinition;
}

const allButtonItems: IButtonItem[] = [
  {
    name: "toggleAudio",
    pages: [{ page: "scenes" }, { page: "markers" }],
    tooltipDisabled: "audio-disabled",
    tooltipEnabled: "audio-enabled",
    iconDisabled: faVolumeXmark,
    iconEnabled: faVolumeHigh,
  },
  {
    name: "toggleChildStudios",
    pages: [{ page: "studios/" }],
    tooltipDisabled: "studios-disabled",
    tooltipEnabled: "studios-enabled",
    iconDisabled: faSitemap,
    iconEnabled: faSitemap,
  },
  {
    name: "toggleChildTags",
    pages: [{ page: "tags/" }],
    tooltipDisabled: "tags-disabled",
    tooltipEnabled: "tags-enabled",
    iconDisabled: faSitemap,
    iconEnabled: faSitemap,
  },
  {
    name: "toggleTagsHover",
    pages: [
      { page: "scenes" },
      { page: "images" },
      { page: "galleries" },
      { page: "performers" },
      { page: "studios/" },
      { page: "tags/" },
    ],
    tooltipDisabled: "tags-hover-disabled",
    tooltipEnabled: "tags-hover-enabled",
    iconDisabled: faLayerGroup,
    iconEnabled: faLayerGroup,
  },
];
export const ListToggleConfigButtons: React.FC = ({}) => {
  const intl = useIntl();
  const { configuration } = React.useContext(ConfigurationContext);
  const [buttonItems] = useState<IButtonItem[]>(allButtonItems);
  const location = useLocation();

  const config = useConfiguration();
  const [saveUI] = useConfigureUI();
  const [saveInterface] = useConfigureInterface();

  const [toggleAudio, setToggleAudio] = useState<boolean>(false);
  const [toggleChildStudios, setToggleChildStudios] = useState<boolean>(false);
  const [toggleChildTags, setToggleChildTags] = useState<boolean>(false);
  const [toggleTagsHover, setToggleTagsHover] = useState<boolean>(false);

  useEffect(() => {
    const audio = configuration?.interface?.soundOnPreview;
    const childStudio = configuration?.ui?.showChildStudioContent;
    const childTag = configuration?.ui?.showChildTagContent;
    const tagHover = configuration?.ui?.showTagCardOnHover;

    if (audio !== undefined && audio != null) {
      setToggleAudio(audio);
    }
    if (childStudio !== undefined) {
      setToggleChildStudios(childStudio);
    }
    if (childTag !== undefined) {
      setToggleChildTags(childTag);
    }
    if (tagHover !== undefined) {
      setToggleTagsHover(tagHover);
    }
  }, [configuration]);

  function setConfigure(mode: string, updatedConfig: object) {
    switch (mode) {
      case "interface":
        saveInterface({
          variables: {
            input: {
              ...config.data?.configuration.interface,
              ...updatedConfig,
            },
          },
        });
        break;
      case "ui":
        saveUI({
          variables: {
            input: {
              ...config.data?.configuration,
              ...updatedConfig,
            },
          },
        });
        break;
    }
  }

  function onSetToggleSetting(matchingPage: IButtonItem) {
    const updatedConfig = { ...config.data?.configuration.ui };
    const updatedConfigInt = { ...config.data?.configuration.interface };

    let mode: string = "";
    let shouldToggle: boolean = false;
    let shouldToggleInt: boolean = false;

    switch (matchingPage.name) {
      case "toggleAudio":
        mode = "interface";
        shouldToggleInt = true;
        updatedConfigInt.soundOnPreview = !toggleAudio;
        setToggleAudio((prevToggleAudio) => !prevToggleAudio);
        break;
      case "toggleChildStudios":
        mode = "ui";
        shouldToggle = true;
        updatedConfig.showChildStudioContent = !toggleChildStudios;
        setToggleChildStudios(!toggleChildStudios);
        break;
      case "toggleChildTags":
        mode = "ui";
        shouldToggle = true;
        updatedConfig.showChildTagContent = !toggleChildTags;
        // section = "showChildTagContent"
        setToggleChildTags((prevToggleChildTags) => !prevToggleChildTags);
        break;
      case "toggleTagsHover":
        mode = "ui";
        shouldToggle = true;
        updatedConfig.showTagCardOnHover = !toggleTagsHover;
        setToggleTagsHover(!toggleTagsHover);
        break;
      default:
        break;
    }
    if (shouldToggle) {
      setConfigure(mode, updatedConfig);
    }

    if (shouldToggleInt) {
      setConfigure(mode, updatedConfigInt);
    }
  }

  function maybeRenderButtons() {
    function evaluateEnabled(matchingPage: IButtonItem) {
      let enabled: boolean;

      switch (matchingPage.name) {
        case "toggleAudio":
          enabled = toggleAudio;
          break;
        case "toggleChildStudios":
          enabled = toggleChildStudios;
          break;
        case "toggleChildTags":
          enabled = toggleChildTags;
          break;
        case "toggleTagsHover":
          enabled = toggleTagsHover;
          break;
        default:
          enabled = false;
      }
      return enabled;
    }

    function returnTooltip(matchingPage: IButtonItem) {
      let enabled: boolean;
      let returnValue: string;

      switch (matchingPage.name) {
        case "toggleAudio":
          enabled = !toggleAudio;
          break;
        case "toggleChildStudios":
          enabled = !toggleChildStudios;
          break;
        case "toggleChildTags":
          enabled = !toggleChildTags;
          break;
        case "toggleTagsHover":
          enabled = !toggleTagsHover;
          break;
        default:
          enabled = false;
      }

      const tooltipKey = enabled
        ? matchingPage.tooltipEnabled
        : matchingPage.tooltipDisabled;

      returnValue = `config_mode.${tooltipKey}`;

      return returnValue;
    }

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
                  <Tooltip id={matchingPage.name}>
                    {intl.formatMessage({
                      id: returnTooltip(matchingPage),
                    })}
                  </Tooltip>
                }
              >
                <Button
                  variant="secondary"
                  active={evaluateEnabled(matchingPage)}
                  onClick={() => onSetToggleSetting(matchingPage)}
                >
                  {evaluateEnabled(matchingPage) ? (
                    <Icon icon={matchingPage.iconEnabled} color="lime" />
                  ) : (
                    <Icon icon={matchingPage.iconDisabled} />
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

  return <>{maybeRenderButtons()}</>;
};
