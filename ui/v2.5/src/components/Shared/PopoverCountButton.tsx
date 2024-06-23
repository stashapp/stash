import {
  faFilm,
  faImage,
  faImages,
  faPlayCircle,
  faUser,
  faVideo,
  faMapMarkerAlt,
} from "@fortawesome/free-solid-svg-icons";
import React, { useMemo } from "react";
import { Button, OverlayTrigger, Tooltip } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";
import { Link } from "react-router-dom";
import { ConfigurationContext } from "src/hooks/Config";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";

type PopoverLinkType =
  | "scene"
  | "image"
  | "gallery"
  | "marker"
  | "movie"
  | "performer"
  | "studio";

interface IProps {
  className?: string;
  url: string;
  type: PopoverLinkType;
  count: number;
}

export const PopoverCountButton: React.FC<IProps> = ({
  className,
  url,
  type,
  count,
}) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const abbreviateCounter = configuration?.ui.abbreviateCounters ?? false;

  const intl = useIntl();

  function getIcon() {
    switch (type) {
      case "scene":
        return faPlayCircle;
      case "image":
        return faImage;
      case "gallery":
        return faImages;
      case "marker":
        return faMapMarkerAlt;
      case "movie":
        return faFilm;
      case "performer":
        return faUser;
      case "studio":
        return faVideo;
    }
  }

  function getPluralOptions() {
    switch (type) {
      case "scene":
        return {
          one: "scene",
          other: "scenes",
        };
      case "image":
        return {
          one: "image",
          other: "images",
        };
      case "gallery":
        return {
          one: "gallery",
          other: "galleries",
        };
      case "marker":
        return {
          one: "marker",
          other: "markers",
        };
      case "movie":
        return {
          one: "movie",
          other: "movies",
        };
      case "performer":
        return {
          one: "performer",
          other: "performers",
        };
      case "studio":
        return {
          one: "studio",
          other: "studios",
        };
    }
  }

  function getTitle() {
    const pluralCategory = intl.formatPlural(count);
    const options = getPluralOptions();
    const plural = options[pluralCategory as "one"] || options.other;
    return `${count} ${plural}`;
  }

  const countEl = useMemo(() => {
    if (!abbreviateCounter) {
      return count;
    }

    const formatted = TextUtils.abbreviateCounter(count);
    return (
      <span>
        <FormattedNumber
          value={formatted.size}
          maximumFractionDigits={formatted.digits}
        />
        {formatted.unit}
      </span>
    );
  }, [count, abbreviateCounter]);

  return (
    <>
      <OverlayTrigger
        overlay={<Tooltip id={`${type}-count-tooltip`}>{getTitle()}</Tooltip>}
        placement="bottom"
      >
        <Link className={className} to={url}>
          <Button className="minimal">
            <Icon icon={getIcon()} />
            <span>{countEl}</span>
          </Button>
        </Link>
      </OverlayTrigger>
    </>
  );
};
