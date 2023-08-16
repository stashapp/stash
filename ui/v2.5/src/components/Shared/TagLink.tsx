import { Badge } from "react-bootstrap";
import React, { useRef } from "react";
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
import {
  ListFilterModel,
  useDefaultFilter,
} from "src/models/list-filter/filter";

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
  tagType?: "performer" | "scene" | "gallery" | "image" | "details";
  performer?: Partial<PerformerDataFragment>;
  marker?: SceneMarkerFragment;
  movie?: Partial<MovieDataFragment>;
  scene?: Partial<Pick<SceneDataFragment, "id" | "title" | "files">>;
  gallery?: Partial<IGallery>;
  className?: string;
  hoverPlacement?: Placement;
}

export const TagLink: React.FC<IProps> = (props: IProps) => {
  const id = useRef("");
  const link = useRef("#");
  const title = useRef("");
  let modeMap = new Map<string, GQL.FilterMode>([
    ["scene", GQL.FilterMode.Scenes],
    ["performer", GQL.FilterMode.Performers],
    ["gallery", GQL.FilterMode.Galleries],
    ["image", GQL.FilterMode.Images],
    ["details", GQL.FilterMode.Tags],
  ]);
  const mode: GQL.FilterMode =
    props.tag && props.tagType != undefined
      ? modeMap.get(props.tagType) || GQL.FilterMode.Scenes
      : GQL.FilterMode.Scenes;
  const defaultFilter: ListFilterModel = useDefaultFilter(mode);

  if (props.tag) {
    id.current = props.tag.id || "";
    switch (props.tagType) {
      case "scene":
      case undefined:
        link.current = NavUtils.makeTagScenesUrl(props.tag, defaultFilter);
        break;
      case "performer":
        link.current = NavUtils.makeTagPerformersUrl(props.tag, defaultFilter);
        break;
      case "gallery":
        link.current = NavUtils.makeTagGalleriesUrl(props.tag, defaultFilter);
        break;
      case "image":
        link.current = NavUtils.makeTagImagesUrl(props.tag, defaultFilter);
        break;
      case "details":
        link.current = NavUtils.makeTagUrl(id.current);
        break;
    }
    title.current = props.tag.name || "";
  } else if (props.performer) {
    link.current = NavUtils.makePerformerScenesUrl(
      props.performer,
      undefined,
      undefined,
      defaultFilter
    );
    title.current = props.performer.name || "";
  } else if (props.movie) {
    link.current = NavUtils.makeMovieScenesUrl(props.movie, defaultFilter);
    title.current = props.movie.name || "";
  } else if (props.marker) {
    link.current = NavUtils.makeSceneMarkerUrl(props.marker);
    title.current = `${markerTitle(
      props.marker
    )} - ${TextUtils.secondsToTimestamp(props.marker.seconds || 0)}`;
  } else if (props.gallery) {
    link.current = `/galleries/${props.gallery.id}`;
    title.current = galleryTitle(props.gallery);
  } else if (props.scene) {
    link.current = `/scenes/${props.scene.id}`;
    title.current = objectTitle(props.scene);
  }

  return (
    <Badge className={cx("tag-item", props.className)} variant="secondary">
      <TagPopover id={id.current} placement={props.hoverPlacement}>
        <Link to={link.current}>{title.current}</Link>
      </TagPopover>
    </Badge>
  );
};
