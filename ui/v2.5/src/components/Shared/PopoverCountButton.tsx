import {
  faFilm,
  faImage,
  faImages,
  faPlayCircle,
} from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Link } from "react-router-dom";
import Icon from "./Icon";

type PopoverLinkType = "scene" | "image" | "gallery" | "movie";

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
  const intl = useIntl();

  function getIcon() {
    switch (type) {
      case "scene":
        return faPlayCircle;
      case "image":
        return faImage;
      case "gallery":
        return faImages;
      case "movie":
        return faFilm;
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
      case "movie":
        return {
          one: "movie",
          other: "movies",
        };
    }
  }

  function getTitle() {
    const pluralCategory = intl.formatPlural(count);
    const options = getPluralOptions();
    const plural = options[pluralCategory as "one"] || options.other;
    return `${count} ${plural}`;
  }

  return (
    <Link className={className} to={url} title={getTitle()}>
      <Button className="minimal">
        <Icon icon={getIcon()} />
        <span>{count}</span>
      </Button>
    </Link>
  );
};
