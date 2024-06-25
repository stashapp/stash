import React, { useMemo, useState } from "react";
import {
  Button,
  ButtonGroup,
  Card,
  Col,
  Dropdown,
  Form,
  OverlayTrigger,
  Row,
  Table,
  Tooltip,
} from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import {
  GalleryLink,
  GroupLink,
  SceneMarkerLink,
  TagLink,
} from "../Shared/TagLink";
import { SweatDrops } from "../Shared/SweatDrops";
import { Pagination } from "src/components/List/Pagination";
import TextUtils from "src/utils/text";
import { DeleteScenesDialog } from "src/components/Scenes/DeleteScenesDialog";
import { EditScenesDialog } from "../Scenes/EditScenesDialog";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import {
  faBox,
  faExclamationTriangle,
  faFileAlt,
  faFilm,
  faImages,
  faMapMarkerAlt,
  faPencilAlt,
  faTag,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { SceneMergeModal } from "../Scenes/SceneMergeDialog";
import { objectTitle } from "src/core/files";

const CLASSNAME = "duplicate-checker";

const defaultDurationDiff = "1";

export const SceneDuplicateChecker: React.FC = () => {
  const intl = useIntl();
  const history = useHistory();
  const query = new URLSearchParams(history.location.search);
  const currentPage = Number.parseInt(query.get("page") ?? "1", 10);
  const pageSize = Number.parseInt(query.get("size") ?? "20", 10);
  const hashDistance = Number.parseInt(query.get("distance") ?? "0", 10);
  const durationDiff = Number.parseFloat(
    query.get("durationDiff") ?? defaultDurationDiff
  );

  const [currentPageSize, setCurrentPageSize] = useState(pageSize);
  const [isMultiDelete, setIsMultiDelete] = useState(false);
  const [deletingScenes, setDeletingScenes] = useState(false);
  const [editingScenes, setEditingScenes] = useState(false);
  const [chkSafeSelect, setChkSafeSelect] = useState(true);

  const [checkedScenes, setCheckedScenes] = useState<Record<string, boolean>>(
    {}
  );

  const { data, loading, refetch } = GQL.useFindDuplicateScenesQuery({
    fetchPolicy: "no-cache",
    variables: {
      distance: hashDistance,
      duration_diff: durationDiff,
    },
  });

  const scenes = data?.findDuplicateScenes ?? [];

  const { data: missingPhash } = GQL.useFindScenesQuery({
    variables: {
      filter: {
        per_page: 0,
      },
      scene_filter: {
        is_missing: "phash",
        file_count: {
          modifier: GQL.CriterionModifier.GreaterThan,
          value: 0,
        },
      },
    },
  });

  const [selectedScenes, setSelectedScenes] = useState<
    GQL.SlimSceneDataFragment[] | null
  >(null);

  const [mergeScenes, setMergeScenes] =
    useState<{ id: string; title: string }[]>();

  const pageOptions = useMemo(() => {
    const pageSizes = [
      10, 20, 30, 40, 50, 100, 150, 200, 250, 500, 750, 1000, 1250, 1500,
    ];

    const filteredSizes = pageSizes.filter((s, i) => {
      return scenes.length > s || i == 0 || scenes.length > pageSizes[i - 1];
    });

    return filteredSizes.map((size) => {
      return (
        <option key={size} value={size}>
          {size}
        </option>
      );
    });
  }, [scenes.length]);

  if (loading) return <LoadingIndicator />;
  if (!data) return <ErrorMessage error="Error searching for duplicates." />;

  const filteredScenes = scenes.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );
  const checkCount = Object.keys(checkedScenes).filter(
    (id) => checkedScenes[id]
  ).length;

  const setQuery = (q: Record<string, string | number | undefined>) => {
    const newQuery = new URLSearchParams(query);
    for (const key of Object.keys(q)) {
      const value = q[key];
      if (value !== undefined) {
        newQuery.set(key, String(value));
      } else {
        newQuery.delete(key);
      }
    }
    history.push({ search: newQuery.toString() });
  };

  const resetCheckboxSelection = () => {
    const updatedScenes: Record<string, boolean> = {};

    Object.keys(checkedScenes).forEach((sceneKey) => {
      updatedScenes[sceneKey] = false;
    });

    setCheckedScenes(updatedScenes);
  };

  function onDeleteDialogClosed(deleted: boolean) {
    setDeletingScenes(false);
    if (deleted) {
      setSelectedScenes(null);
      refetch();
      if (isMultiDelete) setCheckedScenes({});
    }
    resetCheckboxSelection();
  }

  const findLargestScene = (group: GQL.SlimSceneDataFragment[]) => {
    // Get maximum file size of a scene
    const totalSize = (scene: GQL.SlimSceneDataFragment) => {
      return scene.files.reduce((prev: number, f) => Math.max(prev, f.size), 0);
    };
    // Find scene object with maximum total size
    return group.reduce((largest, scene) => {
      const largestSize = totalSize(largest);
      const currentSize = totalSize(scene);
      return currentSize > largestSize ? scene : largest;
    });
  };

  const findLargestResolutionScene = (group: GQL.SlimSceneDataFragment[]) => {
    // Get maximum resolution of a scene
    const sceneResolution = (scene: GQL.SlimSceneDataFragment) => {
      return scene.files.reduce(
        (prev: number, f) => Math.max(prev, f.height * f.width),
        0
      );
    };
    // Find scene object with maximum resolution
    return group.reduce((largest, scene) => {
      const largestSize = sceneResolution(largest);
      const currentSize = sceneResolution(scene);
      return currentSize > largestSize ? scene : largest;
    });
  };

  // Helper to get file date

  const findFirstFileByAge = (
    oldest: boolean,
    compareScenes: GQL.SlimSceneDataFragment[]
  ) => {
    let selectedFile: GQL.VideoFileDataFragment;
    let oldestTimestamp: Date | undefined = undefined;

    // Loop through all files
    for (const file of compareScenes.flatMap((s) => s.files)) {
      // Get timestamp
      const timestamp: Date = new Date(file.mod_time);

      // Check if current file is oldest
      if (oldest) {
        if (oldestTimestamp === undefined || timestamp < oldestTimestamp) {
          oldestTimestamp = timestamp;
          selectedFile = file;
        }
      } else {
        if (oldestTimestamp === undefined || timestamp > oldestTimestamp) {
          oldestTimestamp = timestamp;
          selectedFile = file;
        }
      }
    }

    // Find scene with oldest file
    return compareScenes.find((s) =>
      s.files.some((f) => f.id === selectedFile.id)
    );
  };

  function checkSameCodec(codecGroup: GQL.SlimSceneDataFragment[]) {
    const codecs = codecGroup.map((s) => s.files[0]?.video_codec);
    return new Set(codecs).size === 1;
  }

  function checkSameResolution(dataGroup: GQL.SlimSceneDataFragment[]) {
    const resolutions = dataGroup.map(
      (s) => s.files[0]?.width * s.files[0]?.height
    );
    return new Set(resolutions).size === 1;
  }

  const onSelectLargestClick = () => {
    setSelectedScenes([]);
    const checkedArray: Record<string, boolean> = {};

    filteredScenes.forEach((group) => {
      if (chkSafeSelect && !checkSameCodec(group)) {
        return;
      }
      // Find largest scene in group a
      const largest = findLargestScene(group);
      group.forEach((scene) => {
        if (scene !== largest) {
          checkedArray[scene.id] = true;
        }
      });
    });

    setCheckedScenes(checkedArray);
  };

  const onSelectLargestResolutionClick = () => {
    setSelectedScenes([]);
    const checkedArray: Record<string, boolean> = {};

    filteredScenes.forEach((group) => {
      if (chkSafeSelect && !checkSameCodec(group)) {
        return;
      }
      // Don't select scenes where resolution is identical.
      if (checkSameResolution(group)) {
        return;
      }
      // Find the highest resolution scene in group.
      const highest = findLargestResolutionScene(group);
      group.forEach((scene) => {
        if (scene !== highest) {
          checkedArray[scene.id] = true;
        }
      });
    });

    setCheckedScenes(checkedArray);
  };

  const onSelectByAge = (oldest: boolean) => {
    setSelectedScenes([]);

    const checkedArray: Record<string, boolean> = {};

    filteredScenes.forEach((group) => {
      if (chkSafeSelect && !checkSameCodec(group)) {
        return;
      }

      const oldestScene = findFirstFileByAge(oldest, group);
      group.forEach((scene) => {
        if (scene !== oldestScene) {
          checkedArray[scene.id] = true;
        }
      });
    });

    setCheckedScenes(checkedArray);
  };

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
    resetCheckboxSelection();
  }

  const renderFilesize = (filesize: number | null | undefined) => {
    const { size: parsedSize, unit } = TextUtils.fileSize(filesize ?? 0);
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
          <Icon icon={faExclamationTriangle} className="text-warning" />
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

  function maybeRenderTagPopoverButton(scene: GQL.SlimSceneDataFragment) {
    if (scene.tags.length <= 0) return;

    const popoverContent = scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon={faTag} />
          <span>{scene.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton(scene: GQL.SlimSceneDataFragment) {
    if (scene.performers.length <= 0) return;

    return <PerformerPopoverButton performers={scene.performers} />;
  }

  function maybeRenderGroupPopoverButton(scene: GQL.SlimSceneDataFragment) {
    if (scene.movies.length <= 0) return;

    const popoverContent = scene.movies.map((sceneMovie) => (
      <div className="group-tag-container row" key={sceneMovie.movie.id}>
        <Link
          to={`/groups/${sceneMovie.movie.id}`}
          className="group-tag col m-auto zoom-2"
        >
          <img
            className="image-thumbnail"
            alt={sceneMovie.movie.name ?? ""}
            src={sceneMovie.movie.front_image_path ?? ""}
          />
        </Link>
        <GroupLink
          key={sceneMovie.movie.id}
          group={sceneMovie.movie}
          className="d-block"
        />
      </div>
    ));

    return (
      <HoverPopover
        placement="bottom"
        content={popoverContent}
        className="tag-tooltip"
      >
        <Button className="minimal">
          <Icon icon={faFilm} />
          <span>{scene.movies.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderSceneMarkerPopoverButton(
    scene: GQL.SlimSceneDataFragment
  ) {
    if (scene.scene_markers.length <= 0) return;

    const popoverContent = scene.scene_markers.map((marker) => {
      const markerWithScene = { ...marker, scene: { id: scene.id } };
      return <SceneMarkerLink key={marker.id} marker={markerWithScene} />;
    });

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon={faMapMarkerAlt} />
          <span>{scene.scene_markers.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderOCounter(scene: GQL.SlimSceneDataFragment) {
    if (scene.o_counter) {
      return (
        <div>
          <Button className="minimal">
            <span className="fa-icon">
              <SweatDrops />
            </span>
            <span>{scene.o_counter}</span>
          </Button>
        </div>
      );
    }
  }

  function maybeRenderGallery(scene: GQL.SlimSceneDataFragment) {
    if (scene.galleries.length <= 0) return;

    const popoverContent = scene.galleries.map((gallery) => (
      <GalleryLink key={gallery.id} gallery={gallery} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon={faImages} />
          <span>{scene.galleries.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderFileCount(scene: GQL.SlimSceneDataFragment) {
    if (scene.files.length <= 1) return;

    const popoverContent = (
      <FormattedMessage
        id="files_amount"
        values={{ value: intl.formatNumber(scene.files.length ?? 0) }}
      />
    );

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon={faFileAlt} />
          <span>{scene.files.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderOrganized(scene: GQL.SlimSceneDataFragment) {
    if (scene.organized) {
      return (
        <div>
          <Button className="minimal">
            <Icon icon={faBox} />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderPopoverButtonGroup(scene: GQL.SlimSceneDataFragment) {
    if (
      scene.tags.length > 0 ||
      scene.performers.length > 0 ||
      scene.movies.length > 0 ||
      scene.scene_markers.length > 0 ||
      scene?.o_counter ||
      scene.galleries.length > 0 ||
      scene.files.length > 1 ||
      scene.organized
    ) {
      return (
        <>
          <ButtonGroup className="flex-wrap">
            {maybeRenderTagPopoverButton(scene)}
            {maybeRenderPerformerPopoverButton(scene)}
            {maybeRenderGroupPopoverButton(scene)}
            {maybeRenderSceneMarkerPopoverButton(scene)}
            {maybeRenderOCounter(scene)}
            {maybeRenderGallery(scene)}
            {maybeRenderFileCount(scene)}
            {maybeRenderOrganized(scene)}
          </ButtonGroup>
        </>
      );
    }
  }

  function renderPagination() {
    return (
      <div className="d-flex mt-2 mb-2">
        <h6 className="mr-auto align-self-center">
          <FormattedMessage
            id="dupe_check.found_sets"
            values={{ setCount: scenes.length }}
          />
        </h6>
        {checkCount > 0 && (
          <ButtonGroup>
            <OverlayTrigger
              overlay={
                <Tooltip id="edit">
                  {intl.formatMessage({ id: "actions.edit" })}
                </Tooltip>
              }
            >
              <Button variant="secondary" onClick={onEdit}>
                <Icon icon={faPencilAlt} />
              </Button>
            </OverlayTrigger>
            <OverlayTrigger
              overlay={
                <Tooltip id="delete">
                  {intl.formatMessage({ id: "actions.delete" })}
                </Tooltip>
              }
            >
              <Button variant="danger" onClick={handleDeleteChecked}>
                <Icon icon={faTrash} />
              </Button>
            </OverlayTrigger>
          </ButtonGroup>
        )}
        <Pagination
          itemsPerPage={pageSize}
          currentPage={currentPage}
          totalItems={scenes.length}
          metadataByline={[]}
          onChangePage={(newPage) => {
            setQuery({ page: newPage === 1 ? undefined : newPage });
            resetCheckboxSelection();
          }}
        />
        <Form.Control
          as="select"
          className="w-auto ml-2 btn-secondary"
          defaultValue={pageSize}
          value={currentPageSize}
          onChange={(e) => {
            setCurrentPageSize(parseInt(e.currentTarget.value, 10));
            setQuery({
              size:
                e.currentTarget.value === "20"
                  ? undefined
                  : e.currentTarget.value,
            });
            resetCheckboxSelection();
          }}
        >
          {pageOptions}
        </Form.Control>
      </div>
    );
  }

  function renderMergeDialog() {
    if (mergeScenes) {
      return (
        <SceneMergeModal
          scenes={mergeScenes}
          onClose={(mergedID?: string) => {
            setMergeScenes(undefined);
            if (mergedID) {
              // refresh
              refetch();
            }
          }}
          show
        />
      );
    }
  }

  function onMergeClicked(
    sceneGroup: GQL.SlimSceneDataFragment[],
    scene: GQL.SlimSceneDataFragment
  ) {
    const selected = scenes.flat().filter((s) => checkedScenes[s.id]);

    // if scenes in this group other than this scene are selected, then only
    // the selected scenes will be selected as source. Otherwise all other
    // scenes will be source
    let srcScenes =
      selected.filter((s) => {
        if (s === scene) return false;
        return sceneGroup.includes(s);
      }) ?? [];

    if (!srcScenes.length) {
      srcScenes = sceneGroup.filter((s) => s !== scene);
    }

    // insert subject scene to the front so that it is considered the destination
    srcScenes.unshift(scene);

    setMergeScenes(
      srcScenes.map((s) => {
        return {
          id: s.id,
          title: objectTitle(s),
        };
      })
    );
  }

  return (
    <Card id="scene-duplicate-checker" className="col col-xl-12 mx-auto">
      <div className={CLASSNAME}>
        {deletingScenes && selectedScenes && (
          <DeleteScenesDialog
            selected={selectedScenes}
            onClose={onDeleteDialogClosed}
          />
        )}
        {renderMergeDialog()}
        {maybeRenderEdit()}
        <h4>
          <FormattedMessage id="dupe_check.title" />
        </h4>
        <Form>
          <Form.Group>
            <Row noGutters>
              <Form.Label>
                <FormattedMessage id="dupe_check.search_accuracy_label" />
              </Form.Label>
              <Col xs="auto">
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
                  defaultValue={hashDistance}
                  className="input-control ml-4"
                >
                  <option value={0}>
                    {intl.formatMessage({ id: "dupe_check.options.exact" })}
                  </option>
                  <option value={4}>
                    {intl.formatMessage({ id: "dupe_check.options.high" })}
                  </option>
                  <option value={8}>
                    {intl.formatMessage({ id: "dupe_check.options.medium" })}
                  </option>
                  <option value={10}>
                    {intl.formatMessage({ id: "dupe_check.options.low" })}
                  </option>
                </Form.Control>
              </Col>
            </Row>
            <Form.Text>
              <FormattedMessage id="dupe_check.description" />
            </Form.Text>
          </Form.Group>

          <Form.Group>
            <Row noGutters>
              <Form.Label>
                <FormattedMessage id="dupe_check.duration_diff" />
              </Form.Label>
              <Col xs="auto">
                <Form.Control
                  as="select"
                  onChange={(e) =>
                    setQuery({
                      durationDiff:
                        e.currentTarget.value === defaultDurationDiff
                          ? undefined
                          : e.currentTarget.value,
                      page: undefined,
                    })
                  }
                  defaultValue={durationDiff}
                  className="input-control ml-4"
                >
                  <option value={-1}>
                    {intl.formatMessage({
                      id: "dupe_check.duration_options.any",
                    })}
                  </option>
                  <option value={0}>
                    {intl.formatMessage({
                      id: "dupe_check.duration_options.equal",
                    })}
                  </option>
                  <option value={1}>
                    1 {intl.formatMessage({ id: "second" })}
                  </option>
                  <option value={5}>
                    5 {intl.formatMessage({ id: "seconds" })}
                  </option>
                  <option value={10}>
                    10 {intl.formatMessage({ id: "seconds" })}
                  </option>
                </Form.Control>
              </Col>
            </Row>
          </Form.Group>
          <Form.Group>
            <Row noGutters>
              <Col xs="12">
                <Dropdown className="">
                  <Dropdown.Toggle variant="secondary">
                    <FormattedMessage id="dupe_check.select_options" />
                  </Dropdown.Toggle>
                  <Dropdown.Menu className="bg-secondary text-white">
                    <Dropdown.Item onClick={() => resetCheckboxSelection()}>
                      {intl.formatMessage({ id: "dupe_check.select_none" })}
                    </Dropdown.Item>

                    <Dropdown.Item
                      onClick={() => onSelectLargestResolutionClick()}
                    >
                      {intl.formatMessage({
                        id: "dupe_check.select_all_but_largest_resolution",
                      })}
                    </Dropdown.Item>

                    <Dropdown.Item onClick={() => onSelectLargestClick()}>
                      {intl.formatMessage({
                        id: "dupe_check.select_all_but_largest_file",
                      })}
                    </Dropdown.Item>

                    <Dropdown.Item onClick={() => onSelectByAge(true)}>
                      {intl.formatMessage({
                        id: "dupe_check.select_oldest",
                      })}
                    </Dropdown.Item>

                    <Dropdown.Item onClick={() => onSelectByAge(false)}>
                      {intl.formatMessage({
                        id: "dupe_check.select_youngest",
                      })}
                    </Dropdown.Item>
                  </Dropdown.Menu>
                </Dropdown>
              </Col>
            </Row>
            <Row noGutters>
              <Form.Check
                type="checkbox"
                id="chkSafeSelect"
                label={intl.formatMessage({
                  id: "dupe_check.only_select_matching_codecs",
                })}
                checked={chkSafeSelect}
                onChange={(e) => {
                  setChkSafeSelect(e.target.checked);
                  resetCheckboxSelection();
                }}
              />
            </Row>
          </Form.Group>
        </Form>

        {maybeRenderMissingPhashWarning()}
        {renderPagination()}

        <Table responsive striped className={`${CLASSNAME}-table`}>
          <colgroup>
            <col className={`${CLASSNAME}-checkbox`} />
            <col className={`${CLASSNAME}-sprite`} />
            <col className={`${CLASSNAME}-title`} />
            <col className={`${CLASSNAME}-details`} />
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
              <th>{intl.formatMessage({ id: "details" })}</th>
              <th> </th>
              <th>{intl.formatMessage({ id: "duration" })}</th>
              <th>{intl.formatMessage({ id: "filesize" })}</th>
              <th>{intl.formatMessage({ id: "resolution" })}</th>
              <th>{intl.formatMessage({ id: "bitrate" })}</th>
              <th>{intl.formatMessage({ id: "media_info.video_codec" })}</th>
              <th>{intl.formatMessage({ id: "actions.delete" })}</th>
            </tr>
          </thead>
          <tbody>
            {filteredScenes.map((group, groupIndex) =>
              group.map((scene, i) => {
                const file =
                  scene.files.length > 0 ? scene.files[0] : undefined;

                return (
                  <>
                    {i === 0 && groupIndex !== 0 ? (
                      <tr className="separator" />
                    ) : undefined}
                    <tr
                      className={i === 0 ? "duplicate-group" : ""}
                      key={scene.id}
                    >
                      <td>
                        <Form.Check
                          checked={checkedScenes[scene.id]}
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
                          <img
                            src={scene.paths.sprite ?? ""}
                            alt=""
                            width={100}
                            style={{
                              border: checkedScenes[scene.id]
                                ? "2px solid red"
                                : "",
                            }}
                          />
                        </HoverPopover>
                      </td>
                      <td className="text-left">
                        <p>
                          <Link
                            to={`/scenes/${scene.id}`}
                            style={{
                              fontWeight: checkedScenes[scene.id]
                                ? "bold"
                                : "inherit",
                              textDecoration: checkedScenes[scene.id]
                                ? "line-through 3px"
                                : "inherit",
                              textDecorationColor: checkedScenes[scene.id]
                                ? "red"
                                : "inherit",
                            }}
                          >
                            {" "}
                            {scene.title
                              ? scene.title
                              : TextUtils.fileNameFromPath(
                                  file?.path ?? ""
                                )}{" "}
                          </Link>
                        </p>
                        <p className="scene-path">{file?.path ?? ""}</p>
                      </td>
                      <td className="scene-details">
                        {maybeRenderPopoverButtonGroup(scene)}
                      </td>
                      <td>
                        {file?.duration &&
                          TextUtils.secondsToTimestamp(file.duration)}
                      </td>
                      <td>{renderFilesize(file?.size ?? 0)}</td>
                      <td>{`${file?.width ?? 0}x${file?.height ?? 0}`}</td>
                      <td>
                        <FormattedNumber
                          value={(file?.bit_rate ?? 0) / 1000000}
                          maximumFractionDigits={2}
                        />
                        &nbsp;mbps
                      </td>
                      <td>{file?.video_codec ?? ""}</td>
                      <td>
                        <Button
                          className="edit-button"
                          variant="danger"
                          onClick={() => handleDeleteScene(scene)}
                        >
                          <FormattedMessage id="actions.delete" />
                        </Button>
                        <Button
                          className="edit-button"
                          onClick={() => onMergeClicked(group, scene)}
                        >
                          <FormattedMessage id="actions.merge" />
                        </Button>
                      </td>
                    </tr>
                  </>
                );
              })
            )}
          </tbody>
        </Table>
        {scenes.length === 0 && (
          <h4 className="text-center mt-4">No duplicates found.</h4>
        )}
        {renderPagination()}
      </div>
    </Card>
  );
};

export default SceneDuplicateChecker;
