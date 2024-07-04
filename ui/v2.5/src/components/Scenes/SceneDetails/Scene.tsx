import { Tab, Nav, Dropdown, Button } from "react-bootstrap";
import React, {
  useEffect,
  useState,
  useMemo,
  useContext,
  useRef,
  useLayoutEffect,
} from "react";
import { FormattedDate, FormattedMessage, useIntl } from "react-intl";
import { Link, RouteComponentProps } from "react-router-dom";
import { Helmet } from "react-helmet";
import * as GQL from "src/core/generated-graphql";
import {
  mutateMetadataScan,
  useFindScene,
  useSceneIncrementO,
  useSceneGenerateScreenshot,
  useSceneUpdate,
  queryFindScenes,
  queryFindScenesByID,
  useSceneIncrementPlayCount,
} from "src/core/StashService";

import { SceneEditPanel } from "./SceneEditPanel";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { Icon } from "src/components/Shared/Icon";
import { Counter } from "src/components/Shared/Counter";
import { useToast } from "src/hooks/Toast";
import SceneQueue, { QueuedScene } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import Mousetrap from "mousetrap";
import { OrganizedButton } from "./OrganizedButton";
import { ConfigurationContext } from "src/hooks/Config";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import {
  faEllipsisV,
  faChevronRight,
  faChevronLeft,
} from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import TextUtils from "src/utils/text";
import {
  OCounterButton,
  ViewCountButton,
} from "src/components/Shared/CountButton";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { lazyComponent } from "src/utils/lazyComponent";
import cx from "classnames";
import { TruncatedText } from "src/components/Shared/TruncatedText";

const SubmitStashBoxDraft = lazyComponent(
  () => import("src/components/Dialogs/SubmitDraft")
);
const ScenePlayer = lazyComponent(
  () => import("src/components/ScenePlayer/ScenePlayer")
);

const GalleryViewer = lazyComponent(
  () => import("src/components/Galleries/GalleryViewer")
);
const ExternalPlayerButton = lazyComponent(
  () => import("./ExternalPlayerButton")
);

const QueueViewer = lazyComponent(() => import("./QueueViewer"));
const SceneMarkersPanel = lazyComponent(() => import("./SceneMarkersPanel"));
const SceneFileInfoPanel = lazyComponent(() => import("./SceneFileInfoPanel"));
const SceneDetailPanel = lazyComponent(() => import("./SceneDetailPanel"));
const SceneHistoryPanel = lazyComponent(() => import("./SceneHistoryPanel"));
const SceneGroupPanel = lazyComponent(() => import("./SceneGroupPanel"));
const SceneGalleriesPanel = lazyComponent(
  () => import("./SceneGalleriesPanel")
);
const DeleteScenesDialog = lazyComponent(() => import("../DeleteScenesDialog"));
const GenerateDialog = lazyComponent(
  () => import("../../Dialogs/GenerateDialog")
);
const SceneVideoFilterPanel = lazyComponent(
  () => import("./SceneVideoFilterPanel")
);

const VideoFrameRateResolution: React.FC<{
  width?: number;
  height?: number;
  frameRate?: number;
}> = ({ width, height, frameRate }) => {
  const intl = useIntl();

  const resolution = useMemo(() => {
    if (width && height) {
      const r = TextUtils.resolution(width, height);
      return (
        <span className="resolution" data-value={r}>
          {r}
        </span>
      );
    }
    return undefined;
  }, [width, height]);

  const frameRateDisplay = useMemo(() => {
    if (frameRate) {
      return (
        <span className="frame-rate" data-value={frameRate}>
          <FormattedMessage
            id="frames_per_second"
            values={{ value: intl.formatNumber(frameRate ?? 0) }}
          />
        </span>
      );
    }
    return undefined;
  }, [intl, frameRate]);

  const divider = useMemo(() => {
    return resolution && frameRateDisplay ? (
      <span className="divider"> | </span>
    ) : undefined;
  }, [resolution, frameRateDisplay]);

  return (
    <span>
      {frameRateDisplay}
      {divider}
      {resolution}
    </span>
  );
};

interface IProps {
  scene: GQL.SceneDataFragment;
  setTimestamp: (num: number) => void;
  queueScenes: QueuedScene[];
  onQueueNext: () => void;
  onQueuePrevious: () => void;
  onQueueRandom: () => void;
  onQueueSceneClicked: (sceneID: string) => void;
  onDelete: () => void;
  continuePlaylist: boolean;
  queueHasMoreScenes: boolean;
  onQueueMoreScenes: () => void;
  onQueueLessScenes: () => void;
  queueStart: number;
  collapsed: boolean;
  setCollapsed: (state: boolean) => void;
  setContinuePlaylist: (value: boolean) => void;
}

interface ISceneParams {
  id: string;
}

const ScenePage: React.FC<IProps> = ({
  scene,
  setTimestamp,
  queueScenes,
  onQueueNext,
  onQueuePrevious,
  onQueueRandom,
  onQueueSceneClicked,
  onDelete,
  continuePlaylist,
  queueHasMoreScenes,
  onQueueMoreScenes,
  onQueueLessScenes,
  queueStart,
  collapsed,
  setCollapsed,
  setContinuePlaylist,
}) => {
  const Toast = useToast();
  const intl = useIntl();
  const [updateScene] = useSceneUpdate();
  const [generateScreenshot] = useSceneGenerateScreenshot();
  const { configuration } = useContext(ConfigurationContext);

  const [showDraftModal, setShowDraftModal] = useState(false);
  const boxes = configuration?.general?.stashBoxes ?? [];

  const [incrementO] = useSceneIncrementO(scene.id);

  const [incrementPlay] = useSceneIncrementPlayCount();

  function incrementPlayCount() {
    incrementPlay({
      variables: {
        id: scene.id,
      },
    });
  }

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("scene-details-panel");

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);

  const onIncrementOClick = async () => {
    try {
      await incrementO();
    } catch (e) {
      Toast.error(e);
    }
  };

  function setRating(v: number | null) {
    updateScene({
      variables: {
        input: {
          id: scene.id,
          rating100: v,
        },
      },
    });
  }

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("scene-details-panel"));
    Mousetrap.bind("q", () => setActiveTabKey("scene-queue-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("scene-edit-panel"));
    Mousetrap.bind("k", () => setActiveTabKey("scene-markers-panel"));
    Mousetrap.bind("i", () => setActiveTabKey("scene-file-info-panel"));
    Mousetrap.bind("h", () => setActiveTabKey("scene-history-panel"));
    Mousetrap.bind("o", () => {
      onIncrementOClick();
    });
    Mousetrap.bind("p n", () => onQueueNext());
    Mousetrap.bind("p p", () => onQueuePrevious());
    Mousetrap.bind("p r", () => onQueueRandom());
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("q");
      Mousetrap.unbind("e");
      Mousetrap.unbind("k");
      Mousetrap.unbind("i");
      Mousetrap.unbind("h");
      Mousetrap.unbind("o");
      Mousetrap.unbind("p n");
      Mousetrap.unbind("p p");
      Mousetrap.unbind("p r");
      Mousetrap.unbind(",");
    };
  });

  async function onSave(input: GQL.SceneCreateInput) {
    await updateScene({
      variables: {
        input: {
          id: scene.id,
          ...input,
        },
      },
    });
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "scene" }).toLocaleLowerCase() }
      )
    );
  }

  const onOrganizedClick = async () => {
    try {
      setOrganizedLoading(true);
      await updateScene({
        variables: {
          input: {
            id: scene.id,
            organized: !scene.organized,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setOrganizedLoading(false);
    }
  };

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  async function onRescan() {
    await mutateMetadataScan({
      paths: [objectPath(scene)],
      rescan: true,
    });

    Toast.success(
      intl.formatMessage(
        { id: "toast.rescanning_entity" },
        {
          count: 1,
          singularEntity: intl
            .formatMessage({ id: "scene" })
            .toLocaleLowerCase(),
        }
      )
    );
  }

  async function onGenerateScreenshot(at?: number) {
    await generateScreenshot({
      variables: {
        id: scene.id,
        at,
      },
    });
    Toast.success(intl.formatMessage({ id: "toast.generating_screenshot" }));
  }

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      onDelete();
    }
  }

  function maybeRenderDeleteDialog() {
    if (isDeleteAlertOpen) {
      return (
        <DeleteScenesDialog selected={[scene]} onClose={onDeleteDialogClosed} />
      );
    }
  }

  function maybeRenderSceneGenerateDialog() {
    if (isGenerateDialogOpen) {
      return (
        <GenerateDialog
          selectedIds={[scene.id]}
          onClose={() => {
            setIsGenerateDialogOpen(false);
          }}
          type="scene"
        />
      );
    }
  }

  const renderOperations = () => (
    <Dropdown>
      <Dropdown.Toggle
        variant="secondary"
        id="operation-menu"
        className="minimal"
        title={intl.formatMessage({ id: "operations" })}
      >
        <Icon icon={faEllipsisV} />
      </Dropdown.Toggle>
      <Dropdown.Menu className="bg-secondary text-white">
        {!!scene.files.length && (
          <Dropdown.Item
            key="rescan"
            className="bg-secondary text-white"
            onClick={() => onRescan()}
          >
            <FormattedMessage id="actions.rescan" />
          </Dropdown.Item>
        )}
        <Dropdown.Item
          key="generate"
          className="bg-secondary text-white"
          onClick={() => setIsGenerateDialogOpen(true)}
        >
          <FormattedMessage id="actions.generate" />
        </Dropdown.Item>
        <Dropdown.Item
          key="generate-screenshot"
          className="bg-secondary text-white"
          onClick={() => onGenerateScreenshot(getPlayerPosition())}
        >
          <FormattedMessage id="actions.generate_thumb_from_current" />
        </Dropdown.Item>
        <Dropdown.Item
          key="generate-default"
          className="bg-secondary text-white"
          onClick={() => onGenerateScreenshot()}
        >
          <FormattedMessage id="actions.generate_thumb_default" />
        </Dropdown.Item>
        {boxes.length > 0 && (
          <Dropdown.Item
            key="submit"
            className="bg-secondary text-white"
            onClick={() => setShowDraftModal(true)}
          >
            <FormattedMessage id="actions.submit_stash_box" />
          </Dropdown.Item>
        )}
        <Dropdown.Item
          key="delete-scene"
          className="bg-secondary text-white"
          onClick={() => setIsDeleteAlertOpen(true)}
        >
          <FormattedMessage
            id="actions.delete_entity"
            values={{ entityType: intl.formatMessage({ id: "scene" }) }}
          />
        </Dropdown.Item>
      </Dropdown.Menu>
    </Dropdown>
  );

  const renderTabs = () => (
    <Tab.Container
      activeKey={activeTabKey}
      onSelect={(k) => k && setActiveTabKey(k)}
    >
      <div>
        <Nav variant="tabs" className="mr-auto">
          <Nav.Item>
            <Nav.Link eventKey="scene-details-panel">
              <FormattedMessage id="details" />
            </Nav.Link>
          </Nav.Item>
          {queueScenes.length > 0 ? (
            <Nav.Item>
              <Nav.Link eventKey="scene-queue-panel">
                <FormattedMessage id="queue" />
              </Nav.Link>
            </Nav.Item>
          ) : (
            ""
          )}
          <Nav.Item>
            <Nav.Link eventKey="scene-markers-panel">
              <FormattedMessage id="markers" />
            </Nav.Link>
          </Nav.Item>
          {scene.groups.length > 0 ? (
            <Nav.Item>
              <Nav.Link eventKey="scene-group-panel">
                <FormattedMessage
                  id="countables.groups"
                  values={{ count: scene.groups.length }}
                />
              </Nav.Link>
            </Nav.Item>
          ) : (
            ""
          )}
          {scene.galleries.length >= 1 ? (
            <Nav.Item>
              <Nav.Link eventKey="scene-galleries-panel">
                <FormattedMessage
                  id="countables.galleries"
                  values={{ count: scene.galleries.length }}
                />
              </Nav.Link>
            </Nav.Item>
          ) : undefined}
          <Nav.Item>
            <Nav.Link eventKey="scene-video-filter-panel">
              <FormattedMessage id="effect_filters.name" />
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link eventKey="scene-file-info-panel">
              <FormattedMessage id="file_info" />
              <Counter count={scene.files.length} hideZero hideOne />
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link eventKey="scene-history-panel">
              <FormattedMessage id="history" />
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link eventKey="scene-edit-panel">
              <FormattedMessage id="actions.edit" />
            </Nav.Link>
          </Nav.Item>
        </Nav>
      </div>

      <Tab.Content>
        <Tab.Pane eventKey="scene-details-panel">
          <SceneDetailPanel scene={scene} />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-queue-panel">
          <QueueViewer
            scenes={queueScenes}
            currentID={scene.id}
            continue={continuePlaylist}
            setContinue={setContinuePlaylist}
            onSceneClicked={onQueueSceneClicked}
            onNext={onQueueNext}
            onPrevious={onQueuePrevious}
            onRandom={onQueueRandom}
            start={queueStart}
            hasMoreScenes={queueHasMoreScenes}
            onLessScenes={onQueueLessScenes}
            onMoreScenes={onQueueMoreScenes}
          />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-markers-panel">
          <SceneMarkersPanel
            sceneId={scene.id}
            onClickMarker={onClickMarker}
            isVisible={activeTabKey === "scene-markers-panel"}
          />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-group-panel">
          <SceneGroupPanel scene={scene} />
        </Tab.Pane>
        {scene.galleries.length >= 1 && (
          <Tab.Pane eventKey="scene-galleries-panel">
            <SceneGalleriesPanel galleries={scene.galleries} />
            {scene.galleries.length === 1 && (
              <GalleryViewer galleryId={scene.galleries[0].id} />
            )}
          </Tab.Pane>
        )}
        <Tab.Pane eventKey="scene-video-filter-panel">
          <SceneVideoFilterPanel scene={scene} />
        </Tab.Pane>
        <Tab.Pane className="file-info-panel" eventKey="scene-file-info-panel">
          <SceneFileInfoPanel scene={scene} />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-edit-panel" mountOnEnter>
          <SceneEditPanel
            isVisible={activeTabKey === "scene-edit-panel"}
            scene={scene}
            onSubmit={onSave}
            onDelete={() => setIsDeleteAlertOpen(true)}
          />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-history-panel">
          <SceneHistoryPanel scene={scene} />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );

  function getCollapseButtonIcon() {
    return collapsed ? faChevronRight : faChevronLeft;
  }

  const title = objectTitle(scene);

  const file = useMemo(
    () => (scene.files.length > 0 ? scene.files[0] : undefined),
    [scene]
  );

  return (
    <>
      <Helmet>
        <title>{title}</title>
      </Helmet>
      {maybeRenderSceneGenerateDialog()}
      {maybeRenderDeleteDialog()}
      <div
        className={`scene-tabs order-xl-first order-last ${
          collapsed ? "collapsed" : ""
        }`}
      >
        <div>
          <div className="scene-header-container">
            {scene.studio && (
              <h1 className="text-center scene-studio-image">
                <Link to={`/studios/${scene.studio.id}`}>
                  <img
                    src={scene.studio.image_path ?? ""}
                    alt={`${scene.studio.name} logo`}
                    className="studio-logo"
                  />
                </Link>
              </h1>
            )}
            <h3 className={cx("scene-header", { "no-studio": !scene.studio })}>
              <TruncatedText lineCount={2} text={title} />
            </h3>
          </div>

          <div className="scene-subheader">
            <span className="date" data-value={scene.date}>
              {!!scene.date && (
                <FormattedDate
                  value={scene.date}
                  format="long"
                  timeZone="utc"
                />
              )}
            </span>
            <VideoFrameRateResolution
              width={file?.width}
              height={file?.height}
              frameRate={file?.frame_rate}
            />
          </div>

          <div className="scene-toolbar">
            <span className="scene-toolbar-group">
              <RatingSystem
                value={scene.rating100}
                onSetRating={setRating}
                clickToRate
                withoutContext
              />
            </span>
            <span className="scene-toolbar-group">
              <span>
                <ExternalPlayerButton scene={scene} />
              </span>
              <span>
                <ViewCountButton
                  value={scene.play_count ?? 0}
                  onIncrement={() => incrementPlayCount()}
                />
              </span>
              <span>
                <OCounterButton
                  value={scene.o_counter ?? 0}
                  onIncrement={() => onIncrementOClick()}
                />
              </span>
              <span>
                <OrganizedButton
                  loading={organizedLoading}
                  organized={scene.organized}
                  onClick={onOrganizedClick}
                />
              </span>
              <span>{renderOperations()}</span>
            </span>
          </div>
        </div>
        {renderTabs()}
      </div>
      <div className="scene-divider d-none d-xl-block">
        <Button onClick={() => setCollapsed(!collapsed)}>
          <Icon className="fa-fw" icon={getCollapseButtonIcon()} />
        </Button>
      </div>
      <SubmitStashBoxDraft
        type="scene"
        boxes={boxes}
        entity={scene}
        show={showDraftModal}
        onHide={() => setShowDraftModal(false)}
      />
    </>
  );
};

const SceneLoader: React.FC<RouteComponentProps<ISceneParams>> = ({
  location,
  history,
  match,
}) => {
  const { id } = match.params;
  const { configuration } = useContext(ConfigurationContext);
  const { data, loading, error } = useFindScene(id);

  const [scene, setScene] = useState<GQL.SceneDataFragment>();

  // useLayoutEffect to update before paint
  useLayoutEffect(() => {
    // only update scene when loading is done
    if (!loading) {
      setScene(data?.findScene ?? undefined);
    }
  }, [data, loading]);

  const queryParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search]
  );
  const sceneQueue = useMemo(
    () => SceneQueue.fromQueryParameters(queryParams),
    [queryParams]
  );
  const queryContinue = useMemo(() => {
    let cont = queryParams.get("continue");
    if (cont) {
      return cont === "true";
    } else {
      return !!configuration?.interface.continuePlaylistDefault;
    }
  }, [configuration?.interface.continuePlaylistDefault, queryParams]);

  const [queueScenes, setQueueScenes] = useState<QueuedScene[]>([]);

  const [collapsed, setCollapsed] = useState(false);
  const [continuePlaylist, setContinuePlaylist] = useState(queryContinue);
  const [hideScrubber, setHideScrubber] = useState(
    !(configuration?.interface.showScrubber ?? true)
  );

  const _setTimestamp = useRef<(value: number) => void>();
  const initialTimestamp = useMemo(() => {
    return Number.parseInt(queryParams.get("t") ?? "0", 10);
  }, [queryParams]);

  const [queueTotal, setQueueTotal] = useState(0);
  const [queueStart, setQueueStart] = useState(1);

  const autoplay = queryParams.get("autoplay") === "true";
  const autoPlayOnSelected =
    configuration?.interface.autostartVideoOnPlaySelected ?? false;

  const currentQueueIndex = useMemo(
    () => queueScenes.findIndex((s) => s.id === id),
    [queueScenes, id]
  );

  function getSetTimestamp(fn: (value: number) => void) {
    _setTimestamp.current = fn;
  }

  function setTimestamp(value: number) {
    if (_setTimestamp.current) {
      _setTimestamp.current(value);
    }
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind(".", () => setHideScrubber((value) => !value));

    return () => {
      Mousetrap.unbind(".");
    };
  }, []);

  async function getQueueFilterScenes(filter: ListFilterModel) {
    const query = await queryFindScenes(filter);
    const { scenes, count } = query.data.findScenes;
    setQueueScenes(scenes);
    setQueueTotal(count);
    setQueueStart((filter.currentPage - 1) * filter.itemsPerPage + 1);
  }

  async function getQueueScenes(sceneIDs: number[]) {
    const query = await queryFindScenesByID(sceneIDs);
    const { scenes, count } = query.data.findScenes;
    setQueueScenes(scenes);
    setQueueTotal(count);
    setQueueStart(1);
  }

  useEffect(() => {
    if (sceneQueue.query) {
      getQueueFilterScenes(sceneQueue.query);
    } else if (sceneQueue.sceneIDs) {
      getQueueScenes(sceneQueue.sceneIDs);
    }
  }, [sceneQueue]);

  async function onQueueLessScenes() {
    if (!sceneQueue.query || queueStart <= 1) {
      return;
    }

    const filterCopy = sceneQueue.query.clone();
    const newStart = queueStart - filterCopy.itemsPerPage;
    filterCopy.currentPage = Math.ceil(newStart / filterCopy.itemsPerPage);
    const query = await queryFindScenes(filterCopy);
    const { scenes } = query.data.findScenes;

    // prepend scenes to scene list
    const newScenes = (scenes as QueuedScene[]).concat(queueScenes);
    setQueueScenes(newScenes);
    setQueueStart(newStart);

    return scenes;
  }

  const queueHasMoreScenes = useMemo(() => {
    return queueStart + queueScenes.length - 1 < queueTotal;
  }, [queueStart, queueScenes, queueTotal]);

  async function onQueueMoreScenes() {
    if (!sceneQueue.query || !queueHasMoreScenes) {
      return;
    }

    const filterCopy = sceneQueue.query.clone();
    const newStart = queueStart + queueScenes.length;
    filterCopy.currentPage = Math.ceil(newStart / filterCopy.itemsPerPage);
    const query = await queryFindScenes(filterCopy);
    const { scenes } = query.data.findScenes;

    // append scenes to scene list
    const newScenes = queueScenes.concat(scenes);
    setQueueScenes(newScenes);
    // don't change queue start
    return scenes;
  }

  function loadScene(sceneID: string, autoPlay?: boolean, newPage?: number) {
    const sceneLink = sceneQueue.makeLink(sceneID, {
      newPage,
      autoPlay,
      continue: continuePlaylist,
    });
    history.replace(sceneLink);
  }

  async function queueNext(autoPlay: boolean) {
    if (currentQueueIndex === -1) return;

    if (currentQueueIndex < queueScenes.length - 1) {
      loadScene(queueScenes[currentQueueIndex + 1].id, autoPlay);
    } else {
      // if we're at the end of the queue, load more scenes
      if (currentQueueIndex === queueScenes.length - 1 && queueHasMoreScenes) {
        const loadedScenes = await onQueueMoreScenes();
        if (loadedScenes && loadedScenes.length > 0) {
          // set the page to the next page
          const newPage = (sceneQueue.query?.currentPage ?? 0) + 1;
          loadScene(loadedScenes[0].id, autoPlay, newPage);
        }
      }
    }
  }

  async function queuePrevious(autoPlay: boolean) {
    if (currentQueueIndex === -1) return;

    if (currentQueueIndex > 0) {
      loadScene(queueScenes[currentQueueIndex - 1].id, autoPlay);
    } else {
      // if we're at the beginning of the queue, load the previous page
      if (queueStart > 1) {
        const loadedScenes = await onQueueLessScenes();
        if (loadedScenes && loadedScenes.length > 0) {
          const newPage = (sceneQueue.query?.currentPage ?? 0) - 1;
          loadScene(
            loadedScenes[loadedScenes.length - 1].id,
            autoPlay,
            newPage
          );
        }
      }
    }
  }

  async function queueRandom(autoPlay: boolean) {
    if (sceneQueue.query) {
      const { query } = sceneQueue;
      const pages = Math.ceil(queueTotal / query.itemsPerPage);
      const page = Math.floor(Math.random() * pages) + 1;
      const index = Math.floor(
        Math.random() * Math.min(query.itemsPerPage, queueTotal)
      );
      const filterCopy = sceneQueue.query.clone();
      filterCopy.currentPage = page;
      const queryResults = await queryFindScenes(filterCopy);
      if (queryResults.data.findScenes.scenes.length > index) {
        const { id: sceneID } = queryResults.data.findScenes.scenes[index];
        // navigate to the image player page
        loadScene(sceneID, autoPlay, page);
      }
    } else if (queueTotal !== 0) {
      const index = Math.floor(Math.random() * queueTotal);
      loadScene(queueScenes[index].id, autoPlay);
    }
  }

  function onComplete() {
    // load the next scene if we're continuing
    if (continuePlaylist) {
      queueNext(true);
    }
  }

  function onDelete() {
    if (
      continuePlaylist &&
      currentQueueIndex >= 0 &&
      currentQueueIndex < queueScenes.length - 1
    ) {
      loadScene(queueScenes[currentQueueIndex + 1].id);
    } else {
      history.push("/scenes");
    }
  }

  function getScenePage(sceneID: string) {
    if (!sceneQueue.query) return;

    // find the page that the scene is on
    const index = queueScenes.findIndex((s) => s.id === sceneID);

    if (index === -1) return;

    const perPage = sceneQueue.query.itemsPerPage;
    return Math.floor((index + queueStart - 1) / perPage) + 1;
  }

  function onQueueSceneClicked(sceneID: string) {
    loadScene(sceneID, autoPlayOnSelected, getScenePage(sceneID));
  }

  if (!scene) {
    if (loading) return <LoadingIndicator />;
    if (error) return <ErrorMessage error={error.message} />;
    return <ErrorMessage error={`No scene found with id ${id}.`} />;
  }

  return (
    <div className="row">
      <ScenePage
        scene={scene}
        setTimestamp={setTimestamp}
        queueScenes={queueScenes}
        queueStart={queueStart}
        onDelete={onDelete}
        onQueueNext={() => queueNext(autoPlayOnSelected)}
        onQueuePrevious={() => queuePrevious(autoPlayOnSelected)}
        onQueueRandom={() => queueRandom(autoPlayOnSelected)}
        onQueueSceneClicked={onQueueSceneClicked}
        continuePlaylist={continuePlaylist}
        queueHasMoreScenes={queueHasMoreScenes}
        onQueueLessScenes={onQueueLessScenes}
        onQueueMoreScenes={onQueueMoreScenes}
        collapsed={collapsed}
        setCollapsed={setCollapsed}
        setContinuePlaylist={setContinuePlaylist}
      />
      <div className={`scene-player-container ${collapsed ? "expanded" : ""}`}>
        <ScenePlayer
          key="ScenePlayer"
          scene={scene}
          hideScrubberOverride={hideScrubber}
          autoplay={autoplay}
          permitLoop={!continuePlaylist}
          initialTimestamp={initialTimestamp}
          sendSetTimestamp={getSetTimestamp}
          onComplete={onComplete}
          onNext={() => queueNext(true)}
          onPrevious={() => queuePrevious(true)}
        />
      </div>
    </div>
  );
};

export default SceneLoader;
