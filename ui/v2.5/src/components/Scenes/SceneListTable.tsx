import React, { useState } from "react";
import { Table, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { FormattedMessage, useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import SceneQueue from "src/models/sceneQueue";
import ReactSelect, { components } from "react-select";
import { Icon } from "../Shared/Icon";
import { faTableColumns } from "@fortawesome/free-solid-svg-icons";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { useSceneUpdate } from "src/core/StashService";

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
  const coverImageCol = {
    value: "cover_image",
    label: <FormattedMessage id="cover_image" />,
  };
  const titleCol = { value: "title", label: <FormattedMessage id="title" /> };
  const dateCol = {
    value: "date",
    label: <FormattedMessage id="date" />,
  };
  const ratingCol = {
    value: "rating",
    label: <FormattedMessage id="rating" />,
  };
  const studioCodeCol = {
    value: "scene_code",
    label: <FormattedMessage id="scene_code" />,
  };
  const durationCol = {
    value: "duration",
    label: <FormattedMessage id="duration" />,
  };
  const tagsCol = { value: "tags", label: <FormattedMessage id="tags" /> };
  const performersCol = {
    value: "performers",
    label: <FormattedMessage id="performers" />,
  };
  const studioCol = {
    value: "studio",
    label: <FormattedMessage id="studio" />,
  };
  const moviesCol = {
    value: "movies",
    label: <FormattedMessage id="movies" />,
  };
  const galleriesCol = {
    value: "galleries",
    label: <FormattedMessage id="galleries" />,
  };
  const playCountCol = {
    value: "play_count",
    label: <FormattedMessage id="play_count" />,
  };
  const playDurationCol = {
    value: "play_duration",
    label: <FormattedMessage id="play_duration" />,
  };
  const resolutionCol = {
    value: "resolution",
    label: <FormattedMessage id="resolution" />,
  };
  const bitRateCol = {
    value: "bitrate",
    label: <FormattedMessage id="bitrate" />,
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
    resolutionCol,
    bitRateCol,
  ];

  const defaultColumn = [
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

  const [column, setColumn] = useState(defaultColumn);
  const [updateScene] = useSceneUpdate();

  const handleChange = (selected: any) => {
    setColumn(selected);
  };

  function maybeRenderColHead(col: any) {
    if (column.some((e: any) => e.value === col.value)) {
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
    if (column.some((e: any) => e.value === coverImageCol.value)) {
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
    if (column.some((e: any) => e.value === titleCol.value)) {
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
    if (column.some((e: any) => e.value === dateCol.value)) {
      return (
        <td className={`${dateCol.value}-data`}>
          {scene.date}
        </td>
      );
    }
  }

  function maybeRenderRatingCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === ratingCol.value)) {
      return (
        <td className={`${ratingCol.value}-data`}>
          <RatingSystem
                value={scene.rating100}
                onSetRating={(value) => setRating(value, scene.id)}
              />
          {/* {scene.rating100 ? scene.rating100 : ""} */}
        </td>
      );
    }
  }

  function maybeRenderStudioCodeCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === studioCodeCol.value)) {
      return <td className={`${studioCodeCol.value}-data`}>{scene.code}</td>;
    }
  }

  function maybeRenderDurationCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === durationCol.value)) {
      const file = scene.files.length > 0 ? scene.files[0] : undefined;
      return (
        <td className={`${durationCol.value}-data`}>
          {file?.duration && TextUtils.secondsToTimestamp(file.duration)}
        </td>
      );
    }
  }

  function maybeRenderTagCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === tagsCol.value)) {
      return (
        <td className={`${tagsCol.value}-data`}>
          <ul className="comma-list">
            {scene.tags.map((tag) => (
              <li>
                <Link key={tag.id} to={NavUtils.makeTagScenesUrl(tag)}>
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
    if (column.some((e: any) => e.value === performersCol.value)) {
      return (
        <td className={`${performersCol.value}-data`}>
          <ul className="comma-list">
            {scene.performers.map((performer) => (
              <li>
                <Link
                  key={performer.id}
                  to={NavUtils.makePerformerScenesUrl(performer)}
                >
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
    if (column.some((e: any) => e.value === studioCol.value)) {
      return (
        <td className={`${studioCol.value}-data`}>
          {scene.studio && (
            <Link to={NavUtils.makeStudioScenesUrl(scene.studio)} title={scene.studio.name}>
              <span>{scene.studio.name}</span>
            </Link>
          )}
        </td>
      );
    }
  }

  function maybeRenderMovieCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === moviesCol.value)) {
      return (
        <td className={`${moviesCol.value}-data`}>
          <ul className="comma-list">
            {scene.movies.map((sceneMovie) => (
              <li>
                <Link
                  key={sceneMovie.movie.id}
                  to={NavUtils.makeMovieScenesUrl(sceneMovie.movie)}
                >
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
    if (column.some((e: any) => e.value === galleriesCol.value)) {
      return (
        <td className={`${galleriesCol.value}-data`}>
          <ul className="comma-list">
            {scene.galleries.map((gallery) => (
              <li>
                <Link key={gallery.id} to={`/galleries/${gallery.id}`}>
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
    if (column.some((e: any) => e.value === playCountCol.value)) {
      return (
        <td className={`${playCountCol.value}-data`}>
          {`${scene.play_count} plays`}
        </td>
      );
    }
  }

  function maybeRenderPlayDurationCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === playDurationCol.value)) {
      console.log(scene.play_duration)
      console.log(TextUtils.secondsToTimestamp(scene.play_duration ?? 0))
      return (
        <td className={`${playDurationCol.value}-data`}>
          {TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}
        </td>
      );
    }
  }
  
  function maybeRenderResolutionCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === resolutionCol.value)) {
      return (
        <td className={`${resolutionCol.value}-data`}>
          <ul className="comma-list">
            {scene.files.map((file) => (
              <li>
                <span key={file.id}>
                  {" "}
                  {TextUtils.resolution(file?.width, file?.height)}
                </span>
              </li>
            ))}
          </ul>
        </td>
      );
    }
  }

  function maybeRenderBitRateCell(scene: GQL.SlimSceneDataFragment) {
    if (column.some((e: any) => e.value === bitRateCol.value)) {
      return (
        <td className={`${bitRateCol.value}-data`}>
          <ul className="comma-list">
            {scene.files.map((file) => (
              <li>
                <span key={file.id}>
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
        {maybeRenderResolutionCell(scene)}
        {maybeRenderBitRateCell(scene)}
      </tr>
    );
  };

  const Option = (props: any) => {
    return (
      <div>
        <components.Option {...props}>
          <input
            type="checkbox"
            checked={props.isSelected}
            onChange={() => null}
          />{" "}
          <label>{props.label}</label>
        </components.Option>
      </div>
    );
  };

  const DropdownIndicator = (props: any) => {
    return (
      <div>
        <components.DropdownIndicator {...props}>
          <Icon icon={faTableColumns} className="column-select" />
        </components.DropdownIndicator>
      </div>
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
                <ReactSelect
                  options={Column}
                  value={column}
                  isMulti
                  closeMenuOnSelect={false}
                  hideSelectedOptions={false}
                  isSearchable={false}
                  isClearable={false}
                  components={{
                    DropdownIndicator,
                    Option,
                  }}
                  onChange={handleChange}
                  styles={{
                    container: (base) => ({
                      ...base,
                      display: "inline-block",
                    }),
                    control: (base) => ({
                      ...base,
                      height: "25px",
                      width: "25px",
                      backgroundColor: "none",
                      border: "none",
                      transition: "none",
                    }),
                    valueContainer: (base) => {
                      return {
                        ...base,
                        display: "none",
                      };
                    },
                    dropdownIndicator: (base) => {
                      return {
                        ...base,
                        color: "rgb(255, 255, 255)",
                        padding: "0",
                      };
                    },
                    indicatorSeparator: (base) => {
                      return {
                        ...base,
                        display: "none",
                      };
                    },
                    menu: (base) => {
                      return {
                        ...base,
                        width: "150px!important",
                        backgroundColor: "rgb(57, 75, 89)",
                      };
                    },
                    option: (base, props) => {
                      return {
                        ...base,
                        backgroundColor: props.isFocused
                          ? "rgb(37, 49, 58)"
                          : "rgb(57, 75, 89)",
                        padding: "0px 12px",
                      };
                    },
                    menuList: (base, props) => {
                      return {
                        ...base,
                        position: "fixed",
                      };
                    },
                  }}
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
            {maybeRenderColHead(resolutionCol)}
            {maybeRenderColHead(bitRateCol)}
          </tr>
        </thead>
        <tbody>{props.scenes.map(renderSceneRow)}</tbody>
      </Table>
    </div>
  );
};
