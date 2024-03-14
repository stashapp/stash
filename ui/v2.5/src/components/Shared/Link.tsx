import { useMemo } from "react";
import { Link } from "react-router-dom";
import NavUtils from "src/utils/navigation";

// common link components

export const DirectorLink: React.FC<{
  director: string;
  linkType: "scene" | "movie";
}> = ({ director: director, linkType = "scene" }) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "scene":
        return NavUtils.makeDirectorScenesUrl(director);
      case "movie":
        return NavUtils.makeDirectorMoviesUrl(director);
    }
  }, [director, linkType]);

  return <Link to={link}>{director}</Link>;
};

export const PhotographerLink: React.FC<{
  photographer: string;
  linkType: "gallery" | "image";
}> = ({ photographer, linkType = "image" }) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "gallery":
        return NavUtils.makePhotographerGalleriesUrl(photographer);
      case "image":
        return NavUtils.makePhotographerImagesUrl(photographer);
    }
  }, [photographer, linkType]);

  return <Link to={link}>{photographer}</Link>;
};
