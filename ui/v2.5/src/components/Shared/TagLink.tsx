import { Badge, OverlayTrigger, Tooltip } from "react-bootstrap";
import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import NavUtils, { INamedObject } from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { IFile, IObjectWithTitleFiles, objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import * as GQL from "src/core/generated-graphql";
import { TagPopover } from "../Tags/TagPopover";
import { markerTitle } from "src/core/markers";
import { Placement } from "react-bootstrap/esm/Overlay";
import { faFolderTree } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "../Shared/Icon";
import { FormattedMessage } from "react-intl";

type SceneMarkerFragment = Pick<GQL.SceneMarker, "id" | "title" | "seconds"> & {
  scene: Pick<GQL.Scene, "id">;
  primary_tag: Pick<GQL.Tag, "id" | "name">;
};

interface ICommonLinkProps {
  link: string;
  className?: string;
}

const CommonLinkComponent: React.FC<ICommonLinkProps> = ({
  link,
  className,
  children,
}) => {
  return (
    <Badge className={cx("tag-item", className)} variant="secondary">
      <Link to={link}>{children}</Link>
    </Badge>
  );
};

interface IPerformerLinkProps {
  performer: INamedObject & { disambiguation?: string | null };
  linkType?: "scene" | "gallery" | "image";
  className?: string;
}

export const PerformerLink: React.FC<IPerformerLinkProps> = ({
  performer,
  linkType = "scene",
  className,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "gallery":
        return NavUtils.makePerformerGalleriesUrl(performer);
      case "image":
        return NavUtils.makePerformerImagesUrl(performer);
      case "scene":
      default:
        return NavUtils.makePerformerScenesUrl(performer);
    }
  }, [performer, linkType]);

  const title = performer.name || "";

  return (
    <CommonLinkComponent link={link} className={className}>
      <span>{title}</span>
      {performer.disambiguation && (
        <span className="performer-disambiguation">{` (${performer.disambiguation})`}</span>
      )}
    </CommonLinkComponent>
  );
};

interface IMovieLinkProps {
  movie: INamedObject;
  linkType?: "scene";
  className?: string;
}

export const MovieLink: React.FC<IMovieLinkProps> = ({
  movie,
  linkType = "scene",
  className,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "scene":
        return NavUtils.makeMovieScenesUrl(movie);
    }
  }, [movie, linkType]);

  const title = movie.name || "";

  return (
    <CommonLinkComponent link={link} className={className}>
      {title}
    </CommonLinkComponent>
  );
};

interface ISceneMarkerLinkProps {
  marker: SceneMarkerFragment;
  linkType?: "scene";
  className?: string;
}

export const SceneMarkerLink: React.FC<ISceneMarkerLinkProps> = ({
  marker,
  linkType = "scene",
  className,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "scene":
        return NavUtils.makeSceneMarkerUrl(marker);
    }
  }, [marker, linkType]);

  const title = `${markerTitle(marker)} - ${TextUtils.secondsToTimestamp(
    marker.seconds || 0
  )}`;

  return (
    <CommonLinkComponent link={link} className={className}>
      {title}
    </CommonLinkComponent>
  );
};

interface IObjectWithIDTitleFiles extends IObjectWithTitleFiles {
  id: string;
}

interface ISceneLinkProps {
  scene: IObjectWithIDTitleFiles;
  linkType?: "details";
  className?: string;
}

export const SceneLink: React.FC<ISceneLinkProps> = ({
  scene,
  linkType = "details",
  className,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "details":
        return `/scenes/${scene.id}`;
    }
  }, [scene, linkType]);

  const title = objectTitle(scene);

  return (
    <CommonLinkComponent link={link} className={className}>
      {title}
    </CommonLinkComponent>
  );
};

interface IGallery extends IObjectWithIDTitleFiles {
  folder?: GQL.Maybe<IFile>;
}

interface IGalleryLinkProps {
  gallery: IGallery;
  linkType?: "details";
  className?: string;
}

export const GalleryLink: React.FC<IGalleryLinkProps> = ({
  gallery,
  linkType = "details",
  className,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "details":
        return `/galleries/${gallery.id}`;
    }
  }, [gallery, linkType]);

  const title = galleryTitle(gallery);

  return (
    <CommonLinkComponent link={link} className={className}>
      {title}
    </CommonLinkComponent>
  );
};

interface ITagLinkProps {
  tag: INamedObject;
  linkType?: "scene" | "gallery" | "image" | "details" | "performer" | "movie";
  className?: string;
  hoverPlacement?: Placement;
  showHierarchyIcon?: boolean;
  hierarchyTooltipID?: string;
}

export const TagLink: React.FC<ITagLinkProps> = ({
  tag,
  linkType = "scene",
  className,
  hoverPlacement,
  showHierarchyIcon = false,
  hierarchyTooltipID,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "scene":
        return NavUtils.makeTagScenesUrl(tag);
      case "performer":
        return NavUtils.makeTagPerformersUrl(tag);
      case "gallery":
        return NavUtils.makeTagGalleriesUrl(tag);
      case "image":
        return NavUtils.makeTagImagesUrl(tag);
      case "movie":
        return NavUtils.makeTagMoviesUrl(tag);
      case "details":
        return NavUtils.makeTagUrl(tag.id ?? "");
    }
  }, [tag, linkType]);

  const title = tag.name || "";

  const tooltip = useMemo(() => {
    if (!hierarchyTooltipID) {
      return <></>;
    }

    return (
      <Tooltip id="tag-hierarchy-tooltip">
        <FormattedMessage id={hierarchyTooltipID} />
      </Tooltip>
    );
  }, [hierarchyTooltipID]);

  return (
    <CommonLinkComponent link={link} className={className}>
      <TagPopover id={tag.id ?? ""} placement={hoverPlacement}>
        {title}
        {showHierarchyIcon && (
          <OverlayTrigger placement="top" overlay={tooltip}>
            <span className="icon-wrapper">
              <span className="vertical-line">|</span>
              <Icon icon={faFolderTree} className="tag-icon" />
            </span>
          </OverlayTrigger>
        )}
      </TagPopover>
    </CommonLinkComponent>
  );
};
