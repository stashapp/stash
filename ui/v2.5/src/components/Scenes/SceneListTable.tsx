// @ts-nocheck
/* eslint-disable jsx-a11y/control-has-associated-label */
import React from "react";
import { Table, Button, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";
import { Icon } from "src/components/Shared";
import { FormattedMessage } from "react-intl";
import { objectTitle } from "src/core/files";

interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

export const SceneListTable: React.FC<ISceneListTableProps> = (
  props: ISceneListTableProps
) => {
  const renderTags = (tags: GQL.SlimTagDataFragment[]) =>
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
          <Form.Control
            type="checkbox"
            checked={props.selectedIds.has(scene.id)}
            onChange={() =>
              props.onSelectChange!(
                scene.id,
                !props.selectedIds.has(scene.id),
                shiftKey
              )
            }
            onClick={(
              event: React.MouseEvent<HTMLInputElement, MouseEvent>
            ) => {
              // eslint-disable-next-line prefer-destructuring
              shiftKey = event.shiftKey;
              event.stopPropagation();
            }}
          />
        </td>
        <td>
          <Link to={sceneLink}>
            <img
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
        <td>{scene.rating ? scene.rating : ""}</td>
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
        <td>
          {scene.gallery && (
            <Button className="minimal">
              <Link to={`/galleries/${scene.gallery.id}`}>
                <Icon icon={faImage} />
              </Link>
            </Button>
          )}
        </td>
      </tr>
    );
  };

  return (
    <div className="row table-list justify-content-center">
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
              <FormattedMessage id="gallery" />
            </th>
          </tr>
        </thead>
        <tbody>{props.scenes.map(renderSceneRow)}</tbody>
      </Table>
    </div>
  );
};
