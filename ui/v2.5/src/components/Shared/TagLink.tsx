import { Badge, OverlayTrigger, Tooltip } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import {
  PerformerDataFragment,
  TagDataFragment,
  MovieDataFragment,
  SceneDataFragment,
} from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import * as GQL from "src/core/generated-graphql";
import { TagPopover } from "../Tags/TagPopover";
import { markerTitle } from "src/core/markers";
import { Placement } from "react-bootstrap/esm/Overlay";
import { faFolderTree } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "../Shared/Icon";

interface IFile {
  path: string;
}
interface IGallery {
  id: string;
  files: IFile[];
  folder?: GQL.Maybe<IFile>;
  title: GQL.Maybe<string>;
}

type SceneMarkerFragment = Pick<GQL.SceneMarker, "id" | "title" | "seconds"> & {
  scene: Pick<GQL.Scene, "id">;
  primary_tag: Pick<GQL.Tag, "id" | "name">;
};

interface IProps {
  tag?: Partial<TagDataFragment>;
  linkType?: "performer" | "scene" | "gallery" | "image" | "details";
  performer?: Partial<PerformerDataFragment>;
  marker?: SceneMarkerFragment;
  movie?: Partial<MovieDataFragment>;
  scene?: Partial<Pick<SceneDataFragment, "id" | "title" | "files">>;
  gallery?: Partial<IGallery>;
  className?: string;
  hoverPlacement?: Placement;
  showHierarchyIcon?: boolean;
}

interface ICommonLinkProps {
  id: string;
  link: string;
  title: string;
  className?: string;
  hoverPlacement?: Placement;
  showHierarchyIcon?: boolean;
}

const CommonLinkComponent: React.FC<ICommonLinkProps> = ({
  id,
  link,
  title,
  className,
  hoverPlacement,
  showHierarchyIcon = false,
}) => {
  return (
    <Badge className={cx("tag-item", className)} variant="secondary">
      <TagPopover id={id} placement={hoverPlacement}>
        <Link to={link}>
          {title}
          {showHierarchyIcon && (
            <OverlayTrigger
              placement="top"
              overlay={
                <Tooltip id="tag-hierarchy-tooltip">
                  Explore tag hierarchy
                </Tooltip>
              }
            >
              <span className="icon-wrapper">
                <span className="vertical-line">|</span>
                <Icon icon={faFolderTree} className="tag-icon" />
              </span>
            </OverlayTrigger>
          )}
        </Link>
      </TagPopover>
    </Badge>
  );
};

function getLinkAndTitle(props: IProps) {
  let id: string = "";
  let link: string = "#";
  let title: string = "";
  if (props.tag) {
    id = props.tag.id || "";
    switch (props.linkType) {
      case "scene":
      case undefined:
        link = NavUtils.makeTagScenesUrl(props.tag);
        break;
      case "performer":
        link = NavUtils.makeTagPerformersUrl(props.tag);
        break;
      case "gallery":
        link = NavUtils.makeTagGalleriesUrl(props.tag);
        break;
      case "image":
        link = NavUtils.makeTagImagesUrl(props.tag);
        break;
      case "details":
        link = NavUtils.makeTagUrl(id);
        break;
    }
    title = props.tag.name || "";
  } else if (props.performer) {
    link = NavUtils.makePerformerScenesUrl(props.performer);
    title = props.performer.name || "";
  } else if (props.movie) {
    link = NavUtils.makeMovieScenesUrl(props.movie);
    title = props.movie.name || "";
  } else if (props.marker) {
    link = NavUtils.makeSceneMarkerUrl(props.marker);
    title = `${markerTitle(props.marker)} - ${TextUtils.secondsToTimestamp(
      props.marker.seconds || 0
    )}`;
  } else if (props.gallery) {
    link = `/galleries/${props.gallery.id}`;
    title = galleryTitle(props.gallery);
  } else if (props.scene) {
    link = `/scenes/${props.scene.id}`;
    title = objectTitle(props.scene);
  }

  return { id, link, title };
}

export const PerformerLink: React.FC<IProps> = (props: IProps) => {
  const { id, link, title } = getLinkAndTitle(props);

  return (
    <CommonLinkComponent
      id={id}
      link={link}
      title={title}
      className={props.className}
      hoverPlacement={props.hoverPlacement}
    />
  );
};

export const SceneLink: React.FC<IProps> = (props: IProps) => {
  const { id, link, title } = getLinkAndTitle(props);

  return (
    <CommonLinkComponent
      id={id}
      link={link}
      title={title}
      className={props.className}
      hoverPlacement={props.hoverPlacement}
    />
  );
};

export const GalleryLink: React.FC<IProps> = (props: IProps) => {
  const { id, link, title } = getLinkAndTitle(props);

  return (
    <CommonLinkComponent
      id={id}
      link={link}
      title={title}
      className={props.className}
      hoverPlacement={props.hoverPlacement}
    />
  );
};

export const ImageLink: React.FC<IProps> = (props: IProps) => {
  const { id, link, title } = getLinkAndTitle(props);

  return (
    <CommonLinkComponent
      id={id}
      link={link}
      title={title}
      className={props.className}
      hoverPlacement={props.hoverPlacement}
    />
  );
};

export const DetailsLink: React.FC<IProps> = (props: IProps) => {
  const { id, link, title } = getLinkAndTitle(props);

  return (
    <CommonLinkComponent
      id={id}
      link={link}
      title={title}
      className={props.className}
      hoverPlacement={props.hoverPlacement}
      showHierarchyIcon={props.showHierarchyIcon}
    />
  );
};
