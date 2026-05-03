/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import NavUtils from "src/utils/navigation";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import { useTagUpdate } from "src/core/StashService";
import { useTableColumns } from "src/hooks/useTableColumns";
import cx from "classnames";
import { IColumn, ListTable } from "../List/ListTable";

interface ITagListTableProps {
  tags: GQL.TagListDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "tags";

export const TagListTable: React.FC<ITagListTableProps> = (
  props: ITagListTableProps
) => {
  const intl = useIntl();

  const [updateTag] = useTagUpdate();

  function setFavorite(v: boolean, tagId: string) {
    if (tagId) {
      updateTag({
        variables: {
          input: {
            id: tagId,
            favorite: v,
          },
        },
      });
    }
  }

  const ImageCell = (tag: GQL.TagListDataFragment) => (
    <Link to={`/tags/${tag.id}`}>
      <img
        loading="lazy"
        className="image-thumbnail"
        alt={tag.name ?? ""}
        src={tag.image_path ?? ""}
      />
    </Link>
  );

  const NameCell = (tag: GQL.TagListDataFragment) => (
    <Link to={`/tags/${tag.id}`}>
      <div className="ellips-data" title={tag.name}>
        {tag.name}
      </div>
    </Link>
  );

  const AliasesCell = (tag: GQL.TagListDataFragment) => {
    let aliases = tag.aliases ? tag.aliases.join(", ") : "";
    return (
      <span className="ellips-data" title={aliases}>
        {aliases}
      </span>
    );
  };

  const FavoriteCell = (tag: GQL.TagListDataFragment) => (
    <Button
      className={cx("minimal", tag.favorite ? "favorite" : "not-favorite")}
      onClick={() => setFavorite(!tag.favorite, tag.id)}
    >
      <Icon icon={faHeart} />
    </Button>
  );

  const SceneCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagScenesUrl(tag)}>
      <span>{tag.scene_count}</span>
    </Link>
  );

  const GalleryCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagGalleriesUrl(tag)}>
      <span>{tag.gallery_count}</span>
    </Link>
  );

  const ImageCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagImagesUrl(tag)}>
      <span>{tag.image_count}</span>
    </Link>
  );

  const GroupCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagGroupsUrl(tag)}>
      <span>{tag.group_count}</span>
    </Link>
  );

  const StudioCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagStudiosUrl(tag)}>
      <span>{tag.studio_count}</span>
    </Link>
  );

  const PerformerCountCell = (tag: GQL.TagListDataFragment) => (
    <Link to={NavUtils.makeTagPerformersUrl(tag)}>
      <span>{tag.performer_count}</span>
    </Link>
  );

  interface IColumnSpec {
    value: string;
    label: string;
    defaultShow?: boolean;
    mandatory?: boolean;
    render?: (tag: GQL.TagListDataFragment, index: number) => React.ReactNode;
  }

  const allColumns: IColumnSpec[] = [
    {
      value: "image",
      label: intl.formatMessage({ id: "image" }),
      defaultShow: true,
      render: ImageCell,
    },
    {
      value: "name",
      label: intl.formatMessage({ id: "name" }),
      mandatory: true,
      defaultShow: true,
      render: NameCell,
    },
    {
      value: "aliases",
      label: intl.formatMessage({ id: "aliases" }),
      defaultShow: true,
      render: AliasesCell,
    },
    {
      value: "favourite",
      label: intl.formatMessage({ id: "favourite" }),
      defaultShow: true,
      render: FavoriteCell,
    },
    {
      value: "scene_count",
      label: intl.formatMessage({ id: "scenes" }),
      defaultShow: true,
      render: SceneCountCell,
    },
    {
      value: "gallery_count",
      label: intl.formatMessage({ id: "galleries" }),
      defaultShow: true,
      render: GalleryCountCell,
    },
    {
      value: "image_count",
      label: intl.formatMessage({ id: "images" }),
      defaultShow: true,
      render: ImageCountCell,
    },
    {
      value: "group_count",
      label: intl.formatMessage({ id: "groups" }),
      defaultShow: true,
      render: GroupCountCell,
    },
    {
      value: "performer_count",
      label: intl.formatMessage({ id: "performers" }),
      defaultShow: true,
      render: PerformerCountCell,
    },
    {
      value: "studio_count",
      label: intl.formatMessage({ id: "studios" }),
      defaultShow: true,
      render: StudioCountCell,
    },
  ];

  const defaultColumns = allColumns
    .filter((col) => col.defaultShow)
    .map((col) => col.value);

  const { selectedColumns, saveColumns } = useTableColumns(
    TABLE_NAME,
    defaultColumns
  );

  const columnRenderFuncs: Record<
    string,
    (tag: GQL.TagListDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    tag: GQL.TagListDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];

    if (render) return render(tag, index);
  }

  return (
    <ListTable
      className="tag-table"
      items={props.tags}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={props.selectedIds}
      onSelectChange={props.onSelectChange}
      renderCell={renderCell}
    />
  );
};
