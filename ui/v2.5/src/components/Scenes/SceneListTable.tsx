import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import SceneQueue from "src/models/sceneQueue";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { useSceneUpdate } from "src/core/StashService";
import { IColumn, ListTable } from "../List/ListTable";
import { useTableColumns } from "src/hooks/useTableColumns";

interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "scenes";

export const SceneListTable: React.FC<ISceneListTableProps> = (
  props: ISceneListTableProps
) => {
  const intl = useIntl();

  const [updateScene] = useSceneUpdate();

  function setRating(v: number | null, sceneId: string) {
    if (sceneId) {
      updateScene({
        variables: {
          input: {
            id: sceneId,
            rating100: v,
          },
        },
      });
    }
  }

  const CoverImageCell = (scene: GQL.SlimSceneDataFragment, index: number) => {
    const title = objectTitle(scene);
    const sceneLink = props.queue
      ? props.queue.makeLink(scene.id, { sceneIndex: index })
      : `/scenes/${scene.id}`;

    return (
      <Link to={sceneLink}>
        <img
          loading="lazy"
          className="image-thumbnail"
          alt={title}
          src={scene.paths.screenshot ?? ""}
        />
      </Link>
    );
  };

  const TitleCell = (scene: GQL.SlimSceneDataFragment, index: number) => {
    const title = objectTitle(scene);
    const sceneLink = props.queue
      ? props.queue.makeLink(scene.id, { sceneIndex: index })
      : `/scenes/${scene.id}`;

    return (
      <Link to={sceneLink} title={title}>
        <span className="ellips-data">{title}</span>
      </Link>
    );
  };

  const DateCell = (scene: GQL.SlimSceneDataFragment) => <>{scene.date}</>;

  const RatingCell = (scene: GQL.SlimSceneDataFragment) => (
    <RatingSystem
      value={scene.rating100}
      onSetRating={(value) => setRating(value, scene.id)}
      clickToRate
    />
  );

  const DurationCell = (scene: GQL.SlimSceneDataFragment) => {
    const file = scene.files.length > 0 ? scene.files[0] : undefined;
    return file?.duration && TextUtils.secondsToTimestamp(file.duration);
  };

  const TagCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list overflowable">
      {scene.tags.map((tag) => (
        <li key={tag.id}>
          <Link to={NavUtils.makeTagScenesUrl(tag)}>
            <span>{tag.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PerformersCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list overflowable">
      {scene.performers.map((performer) => (
        <li key={performer.id}>
          <Link to={NavUtils.makePerformerScenesUrl(performer)}>
            <span>{performer.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const StudioCell = (scene: GQL.SlimSceneDataFragment) => {
    if (scene.studio) {
      return (
        <Link
          to={NavUtils.makeStudioScenesUrl(scene.studio)}
          title={scene.studio.name}
        >
          <span className="ellips-data">{scene.studio.name}</span>
        </Link>
      );
    }
  };

  const GroupCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list overflowable">
      {scene.groups.map((sceneGroup) => (
        <li key={sceneGroup.group.id}>
          <Link to={NavUtils.makeGroupScenesUrl(sceneGroup.group)}>
            <span className="ellips-data">{sceneGroup.group.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const GalleriesCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list overflowable">
      {scene.galleries.map((gallery) => (
        <li key={gallery.id}>
          <Link to={`/galleries/${gallery.id}`}>
            <span>{galleryTitle(gallery)}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PlayCountCell = (scene: GQL.SlimSceneDataFragment) => (
    <FormattedMessage
      id="plays"
      values={{ value: intl.formatNumber(scene.play_count ?? 0) }}
    />
  );

  const PlayDurationCell = (scene: GQL.SlimSceneDataFragment) => (
    <>{TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}</>
  );

  const ResolutionCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span> {TextUtils.resolution(file?.width, file?.height)}</span>
        </li>
      ))}
    </ul>
  );

  function renderFileSize(file: { size: number | undefined }) {
    const { size, unit } = TextUtils.fileSize(file.size);

    return (
      <FormattedNumber
        value={size}
        style="unit"
        unit={unit}
        unitDisplay="narrow"
        maximumFractionDigits={2}
      />
    );
  }

  const FileSizeCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list">
      {scene.files.map((file) => (
        <li key={file.id}>{renderFileSize(file)}</li>
      ))}
    </ul>
  );

  const FrameRateCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span>
            <FormattedMessage
              id="frames_per_second"
              values={{ value: intl.formatNumber(file.frame_rate ?? 0) }}
            />
          </span>
        </li>
      ))}
    </ul>
  );

  const BitRateCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span>
            <FormattedMessage
              id="megabits_per_second"
              values={{
                value: intl.formatNumber((file.bit_rate ?? 0) / 1000000, {
                  maximumFractionDigits: 2,
                }),
              }}
            />
          </span>
        </li>
      ))}
    </ul>
  );

  const AudioCodecCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list over">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span>{file.audio_codec}</span>
        </li>
      ))}
    </ul>
  );

  const VideoCodecCell = (scene: GQL.SlimSceneDataFragment) => (
    <ul className="comma-list">
      {scene.files.map((file) => (
        <li key={file.id}>
          <span>{file.video_codec}</span>
        </li>
      ))}
    </ul>
  );

  const PathCell = (scene: GQL.SlimSceneDataFragment) => (
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
      scene: GQL.SlimSceneDataFragment,
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
      value: "scene_code",
      label: intl.formatMessage({ id: "scene_code" }),
      render: (s) => <>{s.code}</>,
    },
    {
      value: "duration",
      label: intl.formatMessage({ id: "duration" }),
      defaultShow: true,
      render: DurationCell,
    },
    {
      value: "studio",
      label: intl.formatMessage({ id: "studio" }),
      defaultShow: true,
      render: StudioCell,
    },
    {
      value: "performers",
      label: intl.formatMessage({ id: "performers" }),
      defaultShow: true,
      render: PerformersCell,
    },
    {
      value: "tags",
      label: intl.formatMessage({ id: "tags" }),
      defaultShow: true,
      render: TagCell,
    },
    {
      value: "groups",
      label: intl.formatMessage({ id: "groups" }),
      defaultShow: true,
      render: GroupCell,
    },
    {
      value: "galleries",
      label: intl.formatMessage({ id: "galleries" }),
      defaultShow: true,
      render: GalleriesCell,
    },
    {
      value: "play_count",
      label: intl.formatMessage({ id: "play_count" }),
      render: PlayCountCell,
    },
    {
      value: "play_duration",
      label: intl.formatMessage({ id: "play_duration" }),
      render: PlayDurationCell,
    },
    {
      value: "o_counter",
      label: intl.formatMessage({ id: "o_count" }),
      render: (s) => <>{s.o_counter}</>,
    },
    {
      value: "resolution",
      label: intl.formatMessage({ id: "resolution" }),
      render: ResolutionCell,
    },
    {
      value: "path",
      label: intl.formatMessage({ id: "path" }),
      render: PathCell,
    },
    {
      value: "filesize",
      label: intl.formatMessage({ id: "filesize" }),
      render: FileSizeCell,
    },
    {
      value: "framerate",
      label: intl.formatMessage({ id: "framerate" }),
      render: FrameRateCell,
    },
    {
      value: "bitrate",
      label: intl.formatMessage({ id: "bitrate" }),
      render: BitRateCell,
    },
    {
      value: "video_codec",
      label: intl.formatMessage({ id: "video_codec" }),
      render: VideoCodecCell,
    },
    {
      value: "audio_codec",
      label: intl.formatMessage({ id: "audio_codec" }),
      render: AudioCodecCell,
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
    (scene: GQL.SlimSceneDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    scene: GQL.SlimSceneDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];

    if (render) return render(scene, index);
  }

  return (
    <ListTable
      className="scene-table"
      items={props.scenes}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={props.selectedIds}
      onSelectChange={props.onSelectChange}
      renderCell={renderCell}
    />
  );
};
