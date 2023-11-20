import React from "react";
import { Table, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { FormattedMessage } from "react-intl";
import { objectTitle } from "src/core/files";
import { galleryTitle } from "src/core/galleries";
import SceneQueue from "src/models/sceneQueue";

interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const SceneListTable: React.FC<ISceneListTableProps> = (
  props: ISceneListTableProps
) => {
  const renderTags = (tags: Partial<GQL.TagDataFragment>[]) =>
    tags.map((tag) => (
      <Link key={tag.id} to={NavUtils.makeTagScenesUrl(tag)}>
        <h6>{tag.name}</h6>
      </Link>
    ));

  const renderPerformers = (performers: Partial<GQL.PerformerDataFragment>[]) =>
    performers.map((performer) => (
      <Link key={performer.id} to={NavUtils.makePerformerScenesUrl(performer)}>
        <h6>{performer.name}</h6>
      </Link>
    ));

  const renderMovies = (scene: GQL.SlimSceneDataFragment) =>
    scene.movies.map((sceneMovie) => (
      <Link
        key={sceneMovie.movie.id}
        to={NavUtils.makeMovieScenesUrl(sceneMovie.movie)}
      >
        <h6>{sceneMovie.movie.name}</h6>
      </Link>
    ));

  const renderGalleries = (scene: GQL.SlimSceneDataFragment) =>
    scene.galleries.map((gallery) => (
      <Link key={gallery.id} to={`/galleries/${gallery.id}`}>
        <h6>{galleryTitle(gallery)}</h6>
      </Link>
    ));

  const renderSceneRow = (scene: GQL.SlimSceneDataFragment, index: number) => {
    const sceneLink = props.queue
      ? props.queue.makeLink(scene.id, { sceneIndex: index })
      : `/scenes/${scene.id}`;

    let shiftKey = false;

    const file = scene.files.length > 0 ? scene.files[0] : undefined;

    const title = objectTitle(scene);
    return (
      <tr key={scene.id}>
        <td>
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
        <td>
          <Link to={sceneLink}>
            <img
              loading="lazy"
              className="image-thumbnail"
              alt={title}
              src={scene.paths.screenshot ?? ""}
            />
          </Link>
        </td>
        <td className="text-left">
          <Link to={sceneLink}>
            <h5>{title}</h5>
          </Link>
        </td>
        <td>{scene.rating100 ? scene.rating100 : ""}</td>
        <td>{file?.duration && TextUtils.secondsToTimestamp(file.duration)}</td>
        <td>{renderTags(scene.tags)}</td>
        <td>{renderPerformers(scene.performers)}</td>
        <td>
          {scene.studio && (
            <Link to={NavUtils.makeStudioScenesUrl(scene.studio)}>
              <h6>{scene.studio.name}</h6>
            </Link>
          )}
        </td>
        <td>{renderMovies(scene)}</td>
        <td>{renderGalleries(scene)}</td>
      </tr>
    );
  };

  return (
    <div className="row scene-table table-list justify-content-center">
      <Table striped bordered>
        <thead>
          <tr>
            <th />
            <th />
            <th className="text-left">
              <FormattedMessage id="title" />
            </th>
            <th>
              <FormattedMessage id="rating" />
            </th>
            <th>
              <FormattedMessage id="duration" />
            </th>
            <th>
              <FormattedMessage id="tags" />
            </th>
            <th>
              <FormattedMessage id="performers" />
            </th>
            <th>
              <FormattedMessage id="studio" />
            </th>
            <th>
              <FormattedMessage id="movies" />
            </th>
            <th>
              <FormattedMessage id="galleries" />
            </th>
          </tr>
        </thead>
        <tbody>{props.scenes.map(renderSceneRow)}</tbody>
      </Table>
    </div>
  );
};
