// @ts-nocheck
/* eslint-disable jsx-a11y/control-has-associated-label */
import React from "react";
import { Table, Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";
import { Icon } from "src/components/Shared";

interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
}

export const SceneListTable: React.FC<ISceneListTableProps> = (
  props: ISceneListTableProps
) => {
  const renderTags = (tags: GQL.Tag[]) =>
    tags.map((tag) => (
      <Link key={tag.id} to={NavUtils.makeTagScenesUrl(tag)}>
        <h6>{tag.name}</h6>
      </Link>
    ));

  const renderPerformers = (performers: Partial<GQL.Performer>[]) =>
    performers.map((performer) => (
      <Link key={performer.id} to={NavUtils.makePerformerScenesUrl(performer)}>
        <h6>{performer.name}</h6>
      </Link>
    ));

  const renderMovies = (scene: GQL.SlimSceneDataFragment) =>
    scene.movies.map((sceneMovie) =>
      !sceneMovie.movie ? undefined : (
        <Link to={NavUtils.makeMovieScenesUrl(sceneMovie.movie)}>
          <h6>{sceneMovie.movie.name}</h6>
        </Link>
      )
    );

  const renderSceneRow = (scene: GQL.SlimSceneDataFragment) => (
    <tr key={scene.id}>
      <td>
        <Link to={`/scenes/${scene.id}`}>
          <img
            className="image-thumbnail"
            alt={scene.title ?? ""}
            src={scene.paths.screenshot ?? ""}
          />
        </Link>
      </td>
      <td className="text-left">
        <Link to={`/scenes/${scene.id}`}>
          <h5 className="text-truncate">
            {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
          </h5>
        </Link>
      </td>
      <td>{scene.rating ? scene.rating : ""}</td>
      <td>
        {scene.file.duration &&
          TextUtils.secondsToTimestamp(scene.file.duration)}
      </td>
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
      <td>
        {scene.gallery && (
          <Button className="minimal">
            <Link to={`/galleries/${scene.gallery.id}`}>
              <Icon icon="image" />
            </Link>
          </Button>
        )}
      </td>
    </tr>
  );

  return (
    <div className="row table-list justify-content-center">
      <Table striped bordered>
        <thead>
          <tr>
            <th />
            <th className="text-left">Title</th>
            <th>Rating</th>
            <th>Duration</th>
            <th>Tags</th>
            <th>Performers</th>
            <th>Studio</th>
            <th>Movies</th>
            <th>Gallery</th>
          </tr>
        </thead>
        <tbody>{props.scenes.map(renderSceneRow)}</tbody>
      </Table>
    </div>
  );
};
