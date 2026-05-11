import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioUpdate } from "src/core/StashService";
import { useTableColumns } from "src/hooks/useTableColumns";
import { ListTable, IColumn } from "../List/ListTable";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { Link } from "react-router-dom";
import NavUtils from "src/utils/navigation";

interface IStudioListTableProps {
  studios: GQL.StudioDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "studios";

export const StudioListTable: React.FC<IStudioListTableProps> = ({
  studios,
  selectedIds,
  onSelectChange,
}) => {
  const [updateStudio] = useStudioUpdate();

  function setRating(v: number | null, studioId: string) {
    if (studioId) {
      updateStudio({
        variables: {
          input: {
            id: studioId,
            rating100: v,
          },
        },
      });
    }
  }

  const ImageCell = (studio: GQL.StudioDataFragment) => (
    <Link to={`/studios/${studio.id}`}>
      <img
        loading="lazy"
        className="image-thumbnail"
        alt={studio.name ?? ""}
        src={studio.image_path ?? ""}
      />
    </Link>
  );

  const NameCell = (studio: GQL.StudioDataFragment) => (
    <Link to={`/studios/${studio.id}`}>
      <div className="ellips-data" title={studio.name ?? ""}>
        {studio.name}
      </div>
    </Link>
  );

  const AliasesCell = (studio: GQL.StudioDataFragment) => {
    const aliases = studio.aliases?.join(", ") ?? "";
    return (
      <span className="ellips-data" title={aliases}>
        {aliases}
      </span>
    );
  };

  const RatingCell = (studio: GQL.StudioDataFragment) => (
    <RatingSystem
      value={studio.rating100}
      onSetRating={(value) => setRating(value, studio.id)}
      clickToRate
    />
  );

  const SceneCountCell = (studio: GQL.StudioDataFragment) => (
    <Link to={NavUtils.makeStudioScenesUrl(studio)}>
      <span>{studio.scene_count}</span>
    </Link>
  );

  const ImageCountCell = (studio: GQL.StudioDataFragment) => (
    <Link to={NavUtils.makeStudioImagesUrl(studio)}>
      <span>{studio.image_count}</span>
    </Link>
  );

  const GalleryCountCell = (studio: GQL.StudioDataFragment) => (
    <Link to={NavUtils.makeStudioGalleriesUrl(studio)}>
      <span>{studio.gallery_count}</span>
    </Link>
  );

  const PerformerCountCell = (studio: GQL.StudioDataFragment) => (
    <Link to={NavUtils.makeStudioPerformersUrl(studio)}>
      <span>{studio.performer_count}</span>
    </Link>
  );

  const RelatedCell = (studio: GQL.StudioDataFragment) => {
    const parentLink = studio.parent_studio ? (
      <Link to={`/studios/${studio.parent_studio.id}`}>
        {studio.parent_studio.name}
      </Link>
    ) : null;
    const childLink =
      studio.child_studios && studio.child_studios.length > 0 ? (
        <Link to={NavUtils.makeChildStudiosUrl(studio)}>
          {studio.child_studios.length}
        </Link>
      ) : null;
    return (
      <div className="studio-related">
        {parentLink}
        {parentLink && childLink && " / "}
        {childLink}
      </div>
    );
  };

  interface IColumnSpec {
    value: string;
    label: string;
    defaultShow?: boolean;
    mandatory?: boolean;
    render?: (studio: GQL.StudioDataFragment, index: number) => React.ReactNode;
  }

  const allColumns: IColumnSpec[] = [
    {
      value: "image",
      label: "Logo",
      defaultShow: true,
      render: ImageCell,
    },
    {
      value: "name",
      label: "Name",
      mandatory: true,
      defaultShow: true,
      render: NameCell,
    },
    {
      value: "aliases",
      label: "Aliases",
      defaultShow: true,
      render: AliasesCell,
    },
    {
      value: "rating",
      label: "Rating",
      defaultShow: true,
      render: RatingCell,
    },
    {
      value: "scene_count",
      label: "Scenes",
      defaultShow: true,
      render: SceneCountCell,
    },
    {
      value: "image_count",
      label: "Images",
      defaultShow: true,
      render: ImageCountCell,
    },
    {
      value: "gallery_count",
      label: "Galleries",
      defaultShow: true,
      render: GalleryCountCell,
    },
    {
      value: "performer_count",
      label: "Performers",
      defaultShow: true,
      render: PerformerCountCell,
    },
    {
      value: "related",
      label: "Related",
      defaultShow: true,
      render: RelatedCell,
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
    (studio: GQL.StudioDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    studio: GQL.StudioDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];
    if (render) return render(studio, index);
  }

  return (
    <ListTable
      className="studio-table"
      items={studios}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={selectedIds}
      onSelectChange={onSelectChange}
      renderCell={renderCell}
    />
  );
};
