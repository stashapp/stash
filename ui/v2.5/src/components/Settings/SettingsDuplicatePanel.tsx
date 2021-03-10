import React, { useState } from "react";
import { Button, Col, Form, Row, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import { FormattedNumber } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import {
  LoadingIndicator,
  ErrorMessage,
  HoverPopover,
} from "src/components/Shared";
import { Pagination } from "src/components/List/Pagination";
import { TextUtils } from "src/utils";
import { DeleteScenesDialog } from "src/components/Scenes/DeleteScenesDialog";

const CLASSNAME = "DuplicateFinder";

export const SettingsDuplicatePanel: React.FC = () => {
  const [page, setPage] = useState(1);
  const [distance, setDistance] = useState(0);
  const [isMultiDelete, setIsMultiDelete] = useState(false);
  const [checkedScenes, setCheckedScenes] = useState<Record<string, boolean>>(
    {}
  );
  const { data, loading, refetch } = GQL.useFindDuplicateScenesQuery({
    fetchPolicy: "no-cache",
    variables: { distance },
  });
  const [deletingScene, setDeletingScene] = useState<
    GQL.SlimSceneDataFragment[] | null
  >(null);

  if (loading) return <LoadingIndicator />;
  if (!data) return <ErrorMessage error="Error searching for duplicates." />;
  const scenes = data?.findDuplicateScenes ?? [];
  const filteredScenes = scenes.slice((page - 1) * 20, page * 20);
  const checkCount = Object.keys(checkedScenes).filter(
    (id) => checkedScenes[id]
  ).length;

  function onDeleteDialogClosed(deleted: boolean) {
    setDeletingScene(null);
    if (deleted) {
      refetch();
      if (isMultiDelete) setCheckedScenes({});
    }
  }

  const handleCheck = (checked: boolean, sceneID: string) => {
    setCheckedScenes({ ...checkedScenes, [sceneID]: checked });
  };

  const handleDeleteChecked = () => {
    setDeletingScene(scenes.flat().filter((s) => checkedScenes[s.id]));
    setIsMultiDelete(true);
  };

  const handleDeleteScene = (scene: GQL.SlimSceneDataFragment) => {
    setDeletingScene([scene]);
    setIsMultiDelete(false);
  };

  const renderFilesize = (filesize: string | null | undefined) => {
    const { size, unit } = TextUtils.fileSize(
      Number.parseInt(filesize ?? "0", 10)
    );
    return (
      <FormattedNumber
        value={size}
        style="unit"
        unit={unit}
        unitDisplay="narrow"
        maximumFractionDigits={2}
      />
    );
  };

  return (
    <div className={CLASSNAME}>
      {deletingScene && (
        <DeleteScenesDialog
          selected={deletingScene}
          onClose={onDeleteDialogClosed}
        />
      )}
      <h4>Duplicate Scenes</h4>
      <Form.Group>
        <Row noGutters>
          <Form.Label>Search Accuracy</Form.Label>
          <Col xs={2}>
            <Form.Control
              as="select"
              onChange={(e) =>
                setDistance(Number.parseInt(e.currentTarget.value, 10) || 0)
              }
              value={distance}
              className="ml-4"
            >
              <option value={0}>Accurate</option>
              <option value={1}>High</option>
              <option value={2}>Medium</option>
              <option value={3}>Low</option>
            </Form.Control>
          </Col>
        </Row>
        <Form.Text>
          Levels below &ldquo;Accurate&rdquo; take longer to calculate and can
          yield larger amounts of false positives.
        </Form.Text>
      </Form.Group>
      <div className="d-flex mb-2">
        <h6 className="mr-auto">{scenes.length} sets of duplicates found.</h6>
        {checkCount > 0 && (
          <Button
            className="edit-button"
            variant="danger"
            onClick={handleDeleteChecked}
          >
            Delete {checkCount} scene{checkCount > 1 && "s"}
          </Button>
        )}
        <Pagination
          itemsPerPage={20}
          currentPage={page}
          totalItems={scenes.length}
          onChangePage={setPage}
        />
      </div>
      <Table striped className={`${CLASSNAME}-table`}>
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
            <th>Title</th>
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
                      <img src={scene.paths.sprite ?? ""} alt="" width={600} />
                    }
                    placement="right"
                  >
                    <img src={scene.paths.sprite ?? ""} alt="" width={100} />
                  </HoverPopover>
                </td>
                <td className="text-left">
                  <Link to={`/scenes/${scene.id}`}>
                    {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
                  </Link>
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
        <h4 className="text-center mt-4">
          No duplicates found. Make sure the phash task has been run.
        </h4>
      )}
    </div>
  );
};
