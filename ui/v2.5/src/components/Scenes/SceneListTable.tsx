import React, { useContext, useState } from "react";
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
import { useConfigureUI, useSceneUpdate } from "src/core/StashService";
import { IUIConfig } from "src/core/config";
import { ConfigurationContext } from "src/hooks/Config";
import { useToast } from "src/hooks/Toast";
import { CheckBoxSelect } from "../Shared/Select";

interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const SceneListTable: React.FC<ISceneListTableProps> = (
  props: ISceneListTableProps
) => {
  const intl = useIntl();
  const { configuration } = useContext(ConfigurationContext);
  const uiConfig = configuration?.ui as IUIConfig | undefined;

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
  const Column = [
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
  const defaultColumn = uiConfig?.defaultSceneColumns
    ? uiConfig.defaultSceneColumns
    : [
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
      ];

  console.log(defaultColumn);

  const [column, setColumn] = useState(defaultColumn);
  const [updateScene] = useSceneUpdate();

  const [saveUI] = useConfigureUI();
  const Toast = useToast();

  async function onUpdateConfig(columns?: [{ value: string; label: string }]) {
    console.log(columns);
    if (!columns) {
      return;
    }

    try {
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            defaultSceneColumns: columns,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function maybeRenderColHead(col: { value: string; label: string }) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === col.value
      )
    ) {
      return <th className={`${col.value}-head`}>{col.label}</th>;
    }
  }

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

  function maybeRenderCoverImageCell(
    scene: GQL.SlimSceneDataFragment,
    sceneLink: string,
    title: string
  ) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === coverImageCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderTitleCell(sceneLink: string, title: string) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === titleCol.value
      )
    ) {
      return (
        <td className={`${titleCol.value}-data`} title={title}>
          <Link to={sceneLink}>
            <span>{title}</span>
          </Link>
        </td>
      );
    }
  }

  function maybeRenderDateCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === dateCol.value
      )
    ) {
      return <td className={`${dateCol.value}-data`}>{scene.date}</td>;
    }
  }

  function maybeRenderRatingCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === ratingCol.value
      )
    ) {
      return (
        <td className={`${ratingCol.value}-data`}>
          <RatingSystem
            value={scene.rating100}
            onSetRating={(value) => setRating(value, scene.id)}
          />
        </td>
      );
    }
  }

  function maybeRenderStudioCodeCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === studioCodeCol.value
      )
    ) {
      return <td className={`${studioCodeCol.value}-data`}>{scene.code}</td>;
    }
  }

  function maybeRenderDurationCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === durationCol.value
      )
    ) {
      const file = scene.files.length > 0 ? scene.files[0] : undefined;
      return (
        <td className={`${durationCol.value}-data`}>
          {file?.duration && TextUtils.secondsToTimestamp(file.duration)}
        </td>
      );
    }
  }

  function maybeRenderTagCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === tagsCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderPerformersCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === performersCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderStudioCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === studioCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderMovieCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === moviesCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderGalleriesCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === galleriesCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderPlayCountCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === playCountCol.value
      )
    ) {
      return (
        <td className={`${playCountCol.value}-data`}>
          {`${scene.play_count} plays`}
        </td>
      );
    }
  }

  function maybeRenderPlayDurationCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) =>
          e.value === playDurationCol.value
      )
    ) {
      return (
        <td className={`${playDurationCol.value}-data`}>
          {TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}
        </td>
      );
    }
  }

  function maybeRenderOCounterCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === oCounterCol.value
      )
    ) {
      return <td className={`${oCounterCol.value}-data`}>{scene.o_counter}</td>;
    }
  }

  function maybeRenderResolutionCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === resolutionCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderFrameRateCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === frameRateCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderBitRateCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === bitRateCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderAudioCodecCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === audioCodecCol.value
      )
    ) {
      return (
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
    }
  }

  function maybeRenderVideoCodecCell(scene: GQL.SlimSceneDataFragment) {
    if (
      column.some(
        (e: { value: string; label: string }) => e.value === videoCodecCol.value
      )
    ) {
      return (
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
    }
  }

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
        {maybeRenderCoverImageCell(scene, sceneLink, title)}
        {maybeRenderTitleCell(sceneLink, title)}
        {maybeRenderDateCell(scene)}
        {maybeRenderRatingCell(scene)}
        {maybeRenderStudioCodeCell(scene)}
        {maybeRenderDurationCell(scene)}
        {maybeRenderTagCell(scene)}
        {maybeRenderPerformersCell(scene)}
        {maybeRenderStudioCell(scene)}
        {maybeRenderMovieCell(scene)}
        {maybeRenderGalleriesCell(scene)}
        {maybeRenderPlayCountCell(scene)}
        {maybeRenderPlayDurationCell(scene)}
        {maybeRenderOCounterCell(scene)}
        {maybeRenderResolutionCell(scene)}
        {maybeRenderFrameRateCell(scene)}
        {maybeRenderBitRateCell(scene)}
        {maybeRenderVideoCodecCell(scene)}
        {maybeRenderAudioCodecCell(scene)}
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
                <CheckBoxSelect
                  options={Column}
                  value={column}
                  setOptions={setColumn}
                  onUpdateConfig={onUpdateConfig}
                />
              </div>
            </th>
            {maybeRenderColHead(coverImageCol)}
            {maybeRenderColHead(titleCol)}
            {maybeRenderColHead(dateCol)}
            {maybeRenderColHead(ratingCol)}
            {maybeRenderColHead(studioCodeCol)}
            {maybeRenderColHead(durationCol)}
            {maybeRenderColHead(tagsCol)}
            {maybeRenderColHead(performersCol)}
            {maybeRenderColHead(studioCol)}
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
