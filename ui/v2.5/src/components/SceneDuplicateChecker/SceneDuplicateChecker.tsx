import React, { useState } from "react";
import {
  Button,
  ButtonGroup,
  Card,
  Col,
  Form,
  OverlayTrigger,
  Row,
  Table,
  Tooltip,
} from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import { FormattedNumber } from "react-intl";
import querystring from "query-string";

import * as GQL from "src/core/generated-graphql";
import {
  LoadingIndicator,
  ErrorMessage,
  HoverPopover,
  Icon,
  TagLink,
} from "src/components/Shared";
import { Pagination } from "src/components/List/Pagination";
import { TextUtils } from "src/utils";
import { DeleteScenesDialog } from "src/components/Scenes/DeleteScenesDialog";
import { sortPerformers } from "src/core/performers";
import { EditScenesDialog } from "../Scenes/EditScenesDialog";

const CLASSNAME = "duplicate-checker";

export const SceneDuplicateChecker: React.FC = () => {
  const history = useHistory();
  const { page, size, distance } = querystring.parse(history.location.search);
  const currentPage = Number.parseInt(
    Array.isArray(page) ? page[0] : page ?? "1",
    10
  );
  const pageSize = Number.parseInt(
    Array.isArray(size) ? size[0] : size ?? "20",
    10
  );
  const hashDistance = Number.parseInt(
    Array.isArray(distance) ? distance[0] : distance ?? "0",
    10
  );
  const [isMultiDelete, setIsMultiDelete] = useState(false);
  const [deletingScenes, setDeletingScenes] = useState(false);
  const [editingScenes, setEditingScenes] = useState(false);
  const [checkedScenes, setCheckedScenes] = useState<Record<string, boolean>>(
    {}
  );
  const { data, loading, refetch } = GQL.useFindDuplicateScenesQuery({
    fetchPolicy: "no-cache",
    variables: { distance: hashDistance },
  });
  const { data: missingPhash } = GQL.useFindScenesQuery({
    variables: {
      filter: {
        per_page: 0,
      },
      scene_filter: {
        is_missing: "phash",
      },
    },
  });

  const [selectedScenes, setSelectedScenes] = useState<
    GQL.SlimSceneDataFragment[] | null
  >(null);

  if (loading) return <LoadingIndicator />;
  if (!data) return <ErrorMessage error="Error searching for duplicates." />;

  const scenes = data?.findDuplicateScenes ?? [];
  const filteredScenes = scenes.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );
  const checkCount = Object.keys(checkedScenes).filter(
    (id) => checkedScenes[id]
  ).length;

  const setQuery = (q: Record<string, string | number | undefined>) => {
    history.push({
      search: querystring.stringify({
        ...querystring.parse(history.location.search),
        ...q,
      }),
    });
  };

  function onDeleteDialogClosed(deleted: boolean) {
    setDeletingScenes(false);
    if (deleted) {
      setSelectedScenes(null);
      refetch();
      if (isMultiDelete) setCheckedScenes({});
    }
  }

  const handleCheck = (checked: boolean, sceneID: string) => {
    setCheckedScenes({ ...checkedScenes, [sceneID]: checked });
  };

  const handleDeleteChecked = () => {
    setSelectedScenes(scenes.flat().filter((s) => checkedScenes[s.id]));
    setDeletingScenes(true);
    setIsMultiDelete(true);
  };

  const handleDeleteScene = (scene: GQL.SlimSceneDataFragment) => {
    setSelectedScenes([scene]);
    setDeletingScenes(true);
    setIsMultiDelete(false);
  };

  function onEdit() {
    setSelectedScenes(scenes.flat().filter((s) => checkedScenes[s.id]));
    setEditingScenes(true);
  }

  const renderFilesize = (filesize: string | null | undefined) => {
    const { size: parsedSize, unit } = TextUtils.fileSize(
      Number.parseInt(filesize ?? "0", 10)
    );
    return (
      <FormattedNumber
        value={parsedSize}
        style="unit"
        unit={unit}
        unitDisplay="narrow"
        maximumFractionDigits={2}
      />
    );
  };

  function maybeRenderMissingPhashWarning() {
    const missingPhashes = missingPhash?.findScenes.count ?? 0;
    if (missingPhashes > 0) {
      return (
        <p className="lead">
          <Icon icon="exclamation-triangle" className="text-warning" />
          Missing phashes for {missingPhashes} scenes. Please run the phash
          generation task.
        </p>
      );
    }
  }

  function maybeRenderEdit() {
    if (editingScenes && selectedScenes) {
      return (
        <EditScenesDialog
          selected={selectedScenes}
          onClose={() => setEditingScenes(false)}
        />
      );
    }
  }

  const renderTags = (tags: GQL.Tag[]) =>
    tags.map((tag) => <TagLink key={tag.id} tag={tag} />);

  const renderPerformers = (performers: Partial<GQL.Performer>[]) => {
    const sorted = sortPerformers(performers);
    return sorted.map((performer) => (
      <div className="performer-tag-container row" key={performer.id}>
        <TagLink key={performer.id} performer={performer} className="d-block" />
      </div>
    ));
  };

  const renderMovies = (scene: GQL.SlimSceneDataFragment) =>
    scene.movies.map((sceneMovie) => (
      <div className="movie-tag-container row" key="movie">
        <TagLink
          key={sceneMovie.movie.id}
          movie={sceneMovie.movie}
          className="d-block"
        />
      </div>
    ));

  function maybeRenderSceneMarkerPopoverButton(
    scene: GQL.SlimSceneDataFragment
  ) {
    if (scene.scene_markers.length <= 0) return;

    const popoverContent = scene.scene_markers.map((marker) => {
      const markerPopover = { ...marker, scene: { id: scene.id } };
      return <TagLink key={marker.id} marker={markerPopover} />;
    });

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="map-marker-alt" />
          <span>{scene.scene_markers.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  return (
    <Card id="scene-duplicate-checker" className="col col-xl-10 mx-auto">
      <div className={CLASSNAME}>
        {deletingScenes && selectedScenes && (
          <DeleteScenesDialog
            selected={selectedScenes}
            onClose={onDeleteDialogClosed}
          />
        )}
        {maybeRenderEdit()}
        <h4>Duplicate Scenes</h4>
        <Form.Group>
          <Row noGutters>
            <Form.Label>Search Accuracy</Form.Label>
            <Col xs={2}>
              <Form.Control
                as="select"
                onChange={(e) =>
                  setQuery({
                    distance:
                      e.currentTarget.value === "0"
                        ? undefined
                        : e.currentTarget.value,
                    page: undefined,
                  })
                }
                defaultValue={distance ?? 0}
                className="input-control ml-4"
              >
                <option value={0}>Exact</option>
                <option value={4}>High</option>
                <option value={8}>Medium</option>
                <option value={10}>Low</option>
              </Form.Control>
            </Col>
          </Row>
          <Form.Text>
            Levels below &ldquo;Exact&rdquo; can take longer to calculate. False
            positives might also be returned on lower accuracy levels.
          </Form.Text>
        </Form.Group>
        {maybeRenderMissingPhashWarning()}
        <div className="d-flex mb-2">
          <h6 className="mr-auto align-self-center">
            {scenes.length} sets of duplicates found.
          </h6>
          {checkCount > 0 && (
            <ButtonGroup>
              <OverlayTrigger overlay={<Tooltip id="edit">Edit</Tooltip>}>
                <Button variant="secondary" onClick={onEdit}>
                  <Icon icon="pencil-alt" />
                </Button>
              </OverlayTrigger>
              <OverlayTrigger overlay={<Tooltip id="delete">Delete</Tooltip>}>
                <Button variant="danger" onClick={handleDeleteChecked}>
                  <Icon icon="trash" />
                </Button>
              </OverlayTrigger>
            </ButtonGroup>
          )}
          <Pagination
            itemsPerPage={pageSize}
            currentPage={currentPage}
            totalItems={scenes.length}
            onChangePage={(newPage) =>
              setQuery({ page: newPage === 1 ? undefined : newPage })
            }
          />
          <Form.Control
            as="select"
            className="w-auto ml-2 btn-secondary"
            defaultValue={pageSize}
            onChange={(e) =>
              setQuery({
                size:
                  e.currentTarget.value === "20"
                    ? undefined
                    : e.currentTarget.value,
              })
            }
          >
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={40}>40</option>
            <option value={60}>60</option>
            <option value={80}>80</option>
          </Form.Control>
        </div>
        <Table responsive striped className={`${CLASSNAME}-table`}>
          <colgroup>
            <col className={`${CLASSNAME}-checkbox`} />
            <col className={`${CLASSNAME}-sprite`} />
            <col className={`${CLASSNAME}-title`} />
            <col className={`${CLASSNAME}-duration`} />
            <col className={`${CLASSNAME}-filesize`} />
            <col className={`${CLASSNAME}-resolution`} />
            <col className={`${CLASSNAME}-bitrate`} />
            <col className={`${CLASSNAME}-codec`} />
            <col className={`${CLASSNAME}-operations`} />
          </colgroup>
          <thead>
            <tr>
              <th> </th>
              <th> </th>
              <th>Details</th>
              <th>Duration</th>
              <th>Filesize</th>
              <th>Resolution</th>
              <th>Bitrate</th>
              <th>Codec</th>
              <th>Delete</th>
            </tr>
          </thead>
          <tbody>
            {filteredScenes.map((group) =>
              group.map((scene, i) => (
                <tr className={i === 0 ? "duplicate-group" : ""} key={scene.id}>
                  <td>
                    <Form.Check
                      onChange={(e) =>
                        handleCheck(e.currentTarget.checked, scene.id)
                      }
                    />
                  </td>
                  <td>
                    <HoverPopover
                      content={
                        <img
                          src={scene.paths.sprite ?? ""}
                          alt=""
                          width={600}
                        />
                      }
                      placement="right"
                    >
                      <img src={scene.paths.sprite ?? ""} alt="" width={100} />
                    </HoverPopover>
                  </td>
                  <td className="text-left">
                    <p>
                      <Link to={`/scenes/${scene.id}`}>
                        {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
                      </Link>
                    </p>
                    <p className="scene-path">{scene.path}</p>
                    <p className="scene-metadata">
                      {renderPerformers(scene.performers)}
                      {renderTags(scene.tags)}
                      {renderMovies(scene)}
                      {maybeRenderSceneMarkerPopoverButton(scene)}
                    </p>
                  </td>
                  <td>
                    {scene.file.duration &&
                      TextUtils.secondsToTimestamp(scene.file.duration)}
                  </td>
                  <td>{renderFilesize(scene.file.size)}</td>
                  <td>{`${scene.file.width}x${scene.file.height}`}</td>
                  <td>
                    <FormattedNumber
                      value={(scene.file.bitrate ?? 0) / 1000000}
                      maximumFractionDigits={2}
                    />
                    &nbsp;mbps
                  </td>
                  <td>{scene.file.video_codec}</td>
                  <td>
                    <Button
                      className="edit-button"
                      variant="danger"
                      onClick={() => handleDeleteScene(scene)}
                    >
                      Delete
                    </Button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </Table>
        {scenes.length === 0 && (
          <h4 className="text-center mt-4">No duplicates found.</h4>
        )}
      </div>
    </Card>
  );
};
