import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { FormattedMessage, useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { useAudioUpdate } from "src/core/StashService";
import { IColumn, ListTable } from "../List/ListTable";
import { useTableColumns } from "src/hooks/useTableColumns";
import { FileSize } from "../Shared/FileSize";

interface IAudioListTableProps {
  audios: GQL.SlimAudioDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "audios";

export const AudioListTable: React.FC<IAudioListTableProps> = (
  props: IAudioListTableProps
) => {
  const intl = useIntl();

  const [updateAudio] = useAudioUpdate();

  function setRating(v: number | null, audioId: string) {
    if (audioId) {
      updateAudio({
        variables: {
          input: {
            id: audioId,
            rating100: v,
          },
        },
      });
    }
  }

  const CoverImageCell = (audio: GQL.SlimAudioDataFragment) => {
    const title = objectTitle(audio);

    return (
      <Link to={`/audios/${audio.id}`}>
        <img
          loading="lazy"
          className="image-thumbnail"
          alt={title}
          src={audio.paths.screenshot ?? ""}
        />
      </Link>
    );
  };

  const TitleCell = (audio: GQL.SlimAudioDataFragment) => {
    const title = objectTitle(audio);

    return (
      <Link to={`/audios/${audio.id}`} title={title}>
        <span className="ellips-data">{title}</span>
      </Link>
    );
  };

  const DateCell = (audio: GQL.SlimAudioDataFragment) => <>{audio.date}</>;

  const RatingCell = (audio: GQL.SlimAudioDataFragment) => (
    <RatingSystem
      value={audio.rating100}
      onSetRating={(value) => setRating(value, audio.id)}
      clickToRate
    />
  );

  const DurationCell = (audio: GQL.SlimAudioDataFragment) => {
    const file = audio.files.length > 0 ? audio.files[0] : undefined;
    return file?.duration && TextUtils.secondsToTimestamp(file.duration);
  };

  const TagCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list overflowable">
      {audio.tags.map((tag) => (
        <li key={tag.id}>
          <Link to={NavUtils.makeTagAudiosUrl(tag)}>
            <span>{tag.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PerformersCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list overflowable">
      {audio.performers.map((performer) => (
        <li key={performer.id}>
          <Link to={NavUtils.makePerformerAudiosUrl(performer)}>
            <span>{performer.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const StudioCell = (audio: GQL.SlimAudioDataFragment) => {
    if (audio.studio) {
      return (
        <Link
          to={NavUtils.makeStudioAudiosUrl(audio.studio)}
          title={audio.studio.name}
        >
          <span className="ellips-data">{audio.studio.name}</span>
        </Link>
      );
    }
  };

  const GroupCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list overflowable">
      {audio.groups.map((audioGroup) => (
        <li key={audioGroup.group.id}>
          <Link to={NavUtils.makeGroupAudiosUrl(audioGroup.group)}>
            <span className="ellips-data">{audioGroup.group.name}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const GalleriesCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list overflowable">
      {audio.galleries.map((gallery) => (
        <li key={gallery.id}>
          <Link to={`/galleries/${gallery.id}`}>
            <span>{galleryTitle(gallery)}</span>
          </Link>
        </li>
      ))}
    </ul>
  );

  const PlayCountCell = (audio: GQL.SlimAudioDataFragment) => (
    <FormattedMessage
      id="plays"
      values={{ value: intl.formatNumber(audio.play_count ?? 0) }}
    />
  );

  const PlayDurationCell = (audio: GQL.SlimAudioDataFragment) => (
    <>{TextUtils.secondsToTimestamp(audio.play_duration ?? 0)}</>
  );

  const FileSizeCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list">
      {audio.files.map((file) => (
        <li key={file.id}>
          <FileSize size={file.size} />
        </li>
      ))}
    </ul>
  );

  const BitRateCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list">
      {audio.files.map((file) => (
        <li key={file.id}>
          <span>
            <FormattedMessage
              id="kilobits_per_second"
              values={{
                value: intl.formatNumber((file.bit_rate ?? 0) / 1000, {
                  maximumFractionDigits: 0,
                }),
              }}
            />
          </span>
        </li>
      ))}
    </ul>
  );

  const SampleRateCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list">
      {audio.files.map((file) => (
        <li key={file.id}>
          <span>
            <FormattedMessage
              id="hertz"
              values={{ value: intl.formatNumber(file.sample_rate ?? 0) }}
            />
          </span>
        </li>
      ))}
    </ul>
  );

  const AudioCodecCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="comma-list over">
      {audio.files.map((file) => (
        <li key={file.id}>
          <span>{file.audio_codec}</span>
        </li>
      ))}
    </ul>
  );

  const PathCell = (audio: GQL.SlimAudioDataFragment) => (
    <ul className="newline-list overflowable TruncatedText">
      {audio.files.map((file) => (
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
      audio: GQL.SlimAudioDataFragment,
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
      value: "audio_code",
      label: intl.formatMessage({ id: "audio_code" }),
      render: (a) => <>{a.code}</>,
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
      render: (a) => <>{a.o_counter}</>,
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
      value: "bitrate",
      label: intl.formatMessage({ id: "bitrate" }),
      render: BitRateCell,
    },
    {
      value: "sample_rate",
      label: intl.formatMessage({ id: "sample_rate" }),
      render: SampleRateCell,
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
    (audio: GQL.SlimAudioDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    audio: GQL.SlimAudioDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];

    if (render) return render(audio, index);
  }

  return (
    <ListTable
      className="audio-table"
      items={props.audios}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={props.selectedIds}
      onSelectChange={props.onSelectChange}
      renderCell={renderCell}
    />
  );
};
