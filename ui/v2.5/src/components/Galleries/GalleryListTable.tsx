import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import { useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { useGalleryUpdate } from "src/core/StashService";
import { IColumn, ListTable } from "../List/ListTable";
import { useTableColumns } from "src/hooks/useTableColumns";

interface IGalleryListTableProps {
  galleries: GQL.SlimGalleryDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "galleries";

export const GalleryListTable: React.FC<IGalleryListTableProps> = (
  props: IGalleryListTableProps
) => {
  const intl = useIntl();

  const [updateGallery] = useGalleryUpdate();

  function setRating(v: number | null, galleryId: string) {
    if (galleryId) {
      updateGallery({
        variables: {
          input: {
            id: galleryId,
            rating100: v,
          },
        },
      });
    }
  }

  const CoverImageCell = (gallery: GQL.SlimGalleryDataFragment) => {
    const title = galleryTitle(gallery);

    return (
      <Link to={`/galleries/${gallery.id}`}>
        <img
          loading="lazy"
          alt={title}
          className="image-thumbnail"
          src={gallery.paths.cover}
        />
      </Link>
    );
  };

  const TitleCell = (gallery: GQL.SlimGalleryDataFragment) => {
    const title = galleryTitle(gallery);

    return (
      <Link to={`/galleries/${gallery.id}`}>
        <span className="ellips-data">{title}</span>
      </Link>
    );
  };

  const DateCell = (gallery: GQL.SlimGalleryDataFragment) => (
    <>{gallery.date}</>
  );

  const RatingCell = (gallery: GQL.SlimGalleryDataFragment) => (
    <RatingSystem
      value={gallery.rating100}
      onSetRating={(value) => setRating(value, gallery.id)}
      clickToRate
    />
  );

  const ImagesCell = (gallery: GQL.SlimGalleryDataFragment) => {
    return (
      <Link to={NavUtils.makeGalleryImagesUrl(gallery)}>
        <span>{gallery.image_count}</span>
      </Link>
    );
  };

  const TagCell = (gallery: GQL.SlimGalleryDataFragment) => (
    <ul className="comma-list overflowable">
      {gallery.tags.map((tag) => (
        <li key={tag.id}>
          <Link to={NavUtils.makeTagGalleriesUrl(tag)}>
            <span>{tag.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PerformersCell = (gallery: GQL.SlimGalleryDataFragment) => (
    <ul className="comma-list overflowable">
      {gallery.performers.map((performer) => (
        <li key={performer.id}>
          <Link to={NavUtils.makePerformerGalleriesUrl(performer)}>
            <span>{performer.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const StudioCell = (gallery: GQL.SlimGalleryDataFragment) => {
    if (gallery.studio) {
      return (
        <Link
          to={NavUtils.makeStudioGalleriesUrl(gallery.studio)}
          title={gallery.studio.name}
        >
          <span className="ellips-data">{gallery.studio.name}</span>
        </Link>
      );
    }
  };

  const SceneCell = (gallery: GQL.SlimGalleryDataFragment) => (
    <ul className="comma-list">
      {gallery.scenes.map((galleryScene) => (
        <li key={galleryScene.id}>
          <Link to={`/scenes/${galleryScene.id}`}>
            <span className="ellips-data">{objectTitle(galleryScene)}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PathCell = (scene: GQL.SlimGalleryDataFragment) => (
    <ul className="newline-list overflowable TruncatedText">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span>{file.path}</span>
        </li>
      ))}
    </ul>
  );

  interface IColumnSpec {
    value: string;
    label: string;
    defaultShow?: boolean;
    mandatory?: boolean;
    render?: (
      gallery: GQL.SlimGalleryDataFragment,
      index: number
    ) => React.ReactNode;
  }

  const allColumns: IColumnSpec[] = [
    {
      value: "cover_image",
      label: intl.formatMessage({ id: "cover_image" }),
      defaultShow: true,
      render: CoverImageCell,
    },
    {
      value: "title",
      label: intl.formatMessage({ id: "title" }),
      defaultShow: true,
      mandatory: true,
      render: TitleCell,
    },
    {
      value: "date",
      label: intl.formatMessage({ id: "date" }),
      defaultShow: true,
      render: DateCell,
    },
    {
      value: "rating",
      label: intl.formatMessage({ id: "rating" }),
      defaultShow: true,
      render: RatingCell,
    },
    {
      value: "code",
      label: intl.formatMessage({ id: "scene_code" }),
      render: (s) => <>{s.code}</>,
    },
    {
      value: "images",
      label: intl.formatMessage({ id: "images" }),
      defaultShow: true,
      render: ImagesCell,
    },
    {
      value: "tags",
      label: intl.formatMessage({ id: "tags" }),
      defaultShow: true,
      render: TagCell,
    },
    {
      value: "performers",
      label: intl.formatMessage({ id: "performers" }),
      defaultShow: true,
      render: PerformersCell,
    },
    {
      value: "studio",
      label: intl.formatMessage({ id: "studio" }),
      defaultShow: true,
      render: StudioCell,
    },
    {
      value: "scenes",
      label: intl.formatMessage({ id: "scenes" }),
      defaultShow: true,
      render: SceneCell,
    },
    {
      value: "photographer",
      label: intl.formatMessage({ id: "photographer" }),
      render: (s) => <>{s.photographer}</>,
    },
    {
      value: "path",
      label: intl.formatMessage({ id: "path" }),
      render: PathCell,
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
    (gallery: GQL.SlimGalleryDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    gallery: GQL.SlimGalleryDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];

    if (render) return render(gallery, index);
  }

  return (
    <ListTable
      className="gallery-table"
      items={props.galleries}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={props.selectedIds}
      onSelectChange={props.onSelectChange}
      renderCell={renderCell}
    />
  );
};
