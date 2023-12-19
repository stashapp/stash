import React from "react";
import { Table, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { FormattedMessage, useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import SceneQueue from "src/models/sceneQueue";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { useSceneUpdate } from "src/core/StashService";
import { useTableColumns } from "src/hooks/useTableColumns";
import { ColumnSelector, IColumn } from "../Shared/ColumnSelector";

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

  const coverImageCol = {
    value: "cover_image",
    label: intl.formatMessage({ id: "cover_image" }),
  };
  const titleCol = {
    value: "title",
    label: intl.formatMessage({ id: "title" }),
  };
  const dateCol = {
    value: "date",
    label: intl.formatMessage({ id: "date" }),
  };
  const ratingCol = {
    value: "rating",
    label: intl.formatMessage({ id: "rating" }),
  };
  const studioCodeCol = {
    value: "scene_code",
    label: intl.formatMessage({ id: "scene_code" }),
  };
  const durationCol = {
    value: "duration",
    label: intl.formatMessage({ id: "duration" }),
  };
  const tagsCol = { value: "tags", label: intl.formatMessage({ id: "tags" }) };
  const performersCol = {
    value: "performers",
    label: intl.formatMessage({ id: "performers" }),
  };
  const studioCol = {
    value: "studio",
    label: intl.formatMessage({ id: "studio" }),
  };
  const moviesCol = {
    value: "movies",
    label: intl.formatMessage({ id: "movies" }),
  };
  const galleriesCol = {
    value: "galleries",
    label: intl.formatMessage({ id: "galleries" }),
  };
  const playCountCol = {
    value: "play_count",
    label: intl.formatMessage({ id: "play_count" }),
  };
  const playDurationCol = {
    value: "play_duration",
    label: intl.formatMessage({ id: "play_duration" }),
  };
  const oCounterCol = {
    value: "o_counter",
    label: intl.formatMessage({ id: "o_counter" }),
  };
  const resolutionCol = {
    value: "resolution",
    label: intl.formatMessage({ id: "resolution" }),
  };
  const frameRateCol = {
    value: "framerate",
    label: intl.formatMessage({ id: "framerate" }),
  };
  const bitRateCol = {
    value: "bitrate",
    label: intl.formatMessage({ id: "bitrate" }),
  };
  const videoCodecCol = {
    value: "video_codec",
    label: intl.formatMessage({ id: "video_codec" }),
  };
  const audioCodecCol = {
    value: "audio_codec",
    label: intl.formatMessage({ id: "audio_codec" }),
  };
  const columns = [
    coverImageCol,
    titleCol,
    dateCol,
    ratingCol,
    studioCodeCol,
    durationCol,
    tagsCol,
    performersCol,
    studioCol,
    moviesCol,
    galleriesCol,
    playCountCol,
    playDurationCol,
    oCounterCol,
    resolutionCol,
    frameRateCol,
    bitRateCol,
    videoCodecCol,
    audioCodecCol,
  ];
  const defaultColumns = [
    coverImageCol,
    titleCol,
    dateCol,
    ratingCol,
    durationCol,
    tagsCol,
    performersCol,
    studioCol,
    moviesCol,
    galleriesCol,
  ].map((c) => c.value);

  const [updateScene] = useSceneUpdate();
  const selectedColumns = useTableColumns(TABLE_NAME, defaultColumns);

  function maybeRenderColHead(column: IColumn) {
    if (selectedColumns[column.value]) {
      return <th className={`${column.value}-head`}>{column.label}</th>;
    }
  }

  const maybeRenderCell = (column: IColumn, cell: React.ReactNode) => {
    if (selectedColumns[column.value]) return cell;
  };

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

  const CoverImageCell = (
    scene: GQL.SlimSceneDataFragment,
    sceneLink: string,
    title: string
  ) => (
    <td className={`${coverImageCol.value}-data`}>
      <Link to={sceneLink}>
        <img
          loading="lazy"
          className="image-thumbnail"
          alt={title}
          src={scene.paths.screenshot ?? ""}
        />
      </Link>
    </td>
  );

  const TitleCell = (sceneLink: string, title: string) => (
    <td className={`${titleCol.value}-data`} title={title}>
      <Link to={sceneLink}>
        <span>{title}</span>
      </Link>
    </td>
  );

  const DateCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${dateCol.value}-data`}>{scene.date}</td>
  );

  const RatingCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${ratingCol.value}-data`}>
      <RatingSystem
        value={scene.rating100}
        onSetRating={(value) => setRating(value, scene.id)}
      />
    </td>
  );

  const StudioCodeCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${studioCodeCol.value}-data`}>{scene.code}</td>
  );

  const DurationCell = (scene: GQL.SlimSceneDataFragment) => {
    const file = scene.files.length > 0 ? scene.files[0] : undefined;
    return (
      <td className={`${durationCol.value}-data`}>
        {file?.duration && TextUtils.secondsToTimestamp(file.duration)}
      </td>
    );
  };

  const TagCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${tagsCol.value}-data`}>
      <ul className="comma-list">
        {scene.tags.map((tag) => (
          <li key={tag.id}>
            <Link to={NavUtils.makeTagScenesUrl(tag)}>
              <span>{tag.name}</span>
            </Link>
          </li>
        ))}
      </ul>
    </td>
  );

  const PerformersCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${performersCol.value}-data`}>
      <ul className="comma-list">
        {scene.performers.map((performer) => (
          <li key={performer.id}>
            <Link to={NavUtils.makePerformerScenesUrl(performer)}>
              <span>{performer.name}</span>
            </Link>
          </li>
        ))}
      </ul>
    </td>
  );

  const StudioCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${studioCol.value}-data`}>
      {scene.studio && (
        <Link
          to={NavUtils.makeStudioScenesUrl(scene.studio)}
          title={scene.studio.name}
        >
          <span>{scene.studio.name}</span>
        </Link>
      )}
    </td>
  );

  const MovieCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${moviesCol.value}-data`}>
      <ul className="comma-list">
        {scene.movies.map((sceneMovie) => (
          <li key={sceneMovie.movie.id}>
            <Link to={NavUtils.makeMovieScenesUrl(sceneMovie.movie)}>
              <span>{sceneMovie.movie.name}</span>
            </Link>
          </li>
        ))}
      </ul>
    </td>
  );

  const GalleriesCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${galleriesCol.value}-data`}>
      <ul className="comma-list">
        {scene.galleries.map((gallery) => (
          <li key={gallery.id}>
            <Link to={`/galleries/${gallery.id}`}>
              <span>{galleryTitle(gallery)}</span>
            </Link>
          </li>
        ))}
      </ul>
    </td>
  );

  const PlayCountCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${playCountCol.value}-data`}>
      <FormattedMessage
        id="plays"
        values={{ value: intl.formatNumber(scene.play_count ?? 0) }}
      />
    </td>
  );

  const PlayDurationCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${playDurationCol.value}-data`}>
      {TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}
    </td>
  );

  const OCounterCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${oCounterCol.value}-data`}>{scene.o_counter}</td>
  );

  const ResolutionCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${resolutionCol.value}-data`}>
      <ul className="comma-list">
        {scene.files.map((file) => (
          <li key={file.id}>
            <span> {TextUtils.resolution(file?.width, file?.height)}</span>
          </li>
        ))}
      </ul>
    </td>
  );

  const FrameRateCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${frameRateCol.value}-data`}>
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
    </td>
  );

  const BitRateCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${bitRateCol.value}-data`}>
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
    </td>
  );

  const AudioCodecCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${audioCodecCol.value}-data`}>
      <ul className="comma-list">
        {scene.files.map((file) => (
          <li key={file.id}>
            <span>{file.audio_codec}</span>
          </li>
        ))}
      </ul>
    </td>
  );

  const VideoCodecCell = (scene: GQL.SlimSceneDataFragment) => (
    <td className={`${videoCodecCol.value}-data`}>
      <ul className="comma-list">
        {scene.files.map((file) => (
          <li key={file.id}>
            <span>{file.video_codec}</span>
          </li>
        ))}
      </ul>
    </td>
  );

  const renderSceneRow = (scene: GQL.SlimSceneDataFragment, index: number) => {
    const sceneLink = props.queue
      ? props.queue.makeLink(scene.id, { sceneIndex: index })
      : `/scenes/${scene.id}`;

    let shiftKey = false;

    const title = objectTitle(scene);
    return (
      <tr key={scene.id}>
        <td className="select-col">
          <label>
            <Form.Control
              type="checkbox"
              checked={props.selectedIds.has(scene.id)}
              onChange={() =>
                props.onSelectChange(
                  scene.id,
                  !props.selectedIds.has(scene.id),
                  shiftKey
                )
              }
              onClick={(
                event: React.MouseEvent<HTMLInputElement, MouseEvent>
              ) => {
                shiftKey = event.shiftKey;
                event.stopPropagation();
              }}
            />
          </label>
        </td>
        {maybeRenderCell(
          coverImageCol,
          CoverImageCell(scene, sceneLink, title)
        )}
        {maybeRenderCell(titleCol, TitleCell(sceneLink, title))}
        {maybeRenderCell(dateCol, DateCell(scene))}
        {maybeRenderCell(ratingCol, RatingCell(scene))}
        {maybeRenderCell(studioCol, StudioCell(scene))}
        {maybeRenderCell(studioCodeCol, StudioCodeCell(scene))}
        {maybeRenderCell(durationCol, DurationCell(scene))}
        {maybeRenderCell(performersCol, PerformersCell(scene))}
        {maybeRenderCell(tagsCol, TagCell(scene))}
        {maybeRenderCell(moviesCol, MovieCell(scene))}
        {maybeRenderCell(galleriesCol, GalleriesCell(scene))}
        {maybeRenderCell(playCountCol, PlayCountCell(scene))}
        {maybeRenderCell(playDurationCol, PlayDurationCell(scene))}
        {maybeRenderCell(oCounterCol, OCounterCell(scene))}
        {maybeRenderCell(resolutionCol, ResolutionCell(scene))}
        {maybeRenderCell(frameRateCol, FrameRateCell(scene))}
        {maybeRenderCell(bitRateCol, BitRateCell(scene))}
        {maybeRenderCell(videoCodecCol, VideoCodecCell(scene))}
        {maybeRenderCell(audioCodecCol, AudioCodecCell(scene))}
      </tr>
    );
  };

  return (
    <div className="row scene-table table-list justify-content-center">
      <Table striped bordered>
        <thead>
          <tr>
            <th className="select-col">
              <div
                className="d-inline-block"
                data-toggle="popover"
                data-trigger="focus"
              >
                <ColumnSelector
                  tableName={TABLE_NAME}
                  columns={columns}
                  defaultColumns={defaultColumns}
                />
              </div>
            </th>
            {maybeRenderColHead(coverImageCol)}
            {maybeRenderColHead(titleCol)}
            {maybeRenderColHead(dateCol)}
            {maybeRenderColHead(ratingCol)}
            {maybeRenderColHead(studioCol)}
            {maybeRenderColHead(studioCodeCol)}
            {maybeRenderColHead(durationCol)}
            {maybeRenderColHead(performersCol)}
            {maybeRenderColHead(tagsCol)}
            {maybeRenderColHead(moviesCol)}
            {maybeRenderColHead(galleriesCol)}
            {maybeRenderColHead(playCountCol)}
            {maybeRenderColHead(playDurationCol)}
            {maybeRenderColHead(oCounterCol)}
            {maybeRenderColHead(resolutionCol)}
            {maybeRenderColHead(frameRateCol)}
            {maybeRenderColHead(bitRateCol)}
            {maybeRenderColHead(videoCodecCol)}
            {maybeRenderColHead(audioCodecCol)}
          </tr>
        </thead>
        <tbody>{props.scenes.map(renderSceneRow)}</tbody>
      </Table>
    </div>
  );
};
