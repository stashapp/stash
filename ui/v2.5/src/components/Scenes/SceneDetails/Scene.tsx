import { Tab, Nav, Dropdown, Button, ButtonGroup } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState, useMemo, useContext, lazy } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useParams, useLocation, useHistory, Link } from "react-router-dom";
import { Helmet } from "react-helmet";
import * as GQL from "src/core/generated-graphql";
import {
  mutateMetadataScan,
  useFindScene,
  useSceneIncrementO,
  useSceneDecrementO,
  useSceneResetO,
  useSceneGenerateScreenshot,
  useSceneUpdate,
  queryFindScenes,
  queryFindScenesByID,
} from "src/core/StashService";

import Icon from "src/components/Shared/Icon";
import { useToast } from "src/hooks";
import SceneQueue from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import Mousetrap from "mousetrap";
import { OCounterButton } from "./OCounterButton";
import { OrganizedButton } from "./OrganizedButton";
import { ConfigurationContext } from "src/hooks/Config";
import { getPlayerPosition } from "src/components/ScenePlayer/util";
import { faEllipsisV } from "@fortawesome/free-solid-svg-icons";

const SubmitStashBoxDraft = lazy(
  () => import("src/components/Dialogs/SubmitDraft")
);
const ScenePlayer = lazy(
  () => import("src/components/ScenePlayer/ScenePlayer")
);

const GalleryViewer = lazy(
  () => import("src/components/Galleries/GalleryViewer")
);
const ExternalPlayerButton = lazy(() => import("./ExternalPlayerButton"));

const QueueViewer = lazy(() => import("./QueueViewer"));
const SceneMarkersPanel = lazy(() => import("./SceneMarkersPanel"));
const SceneFileInfoPanel = lazy(() => import("./SceneFileInfoPanel"));
const SceneEditPanel = lazy(() => import("./SceneEditPanel"));
const SceneDetailPanel = lazy(() => import("./SceneDetailPanel"));
const SceneMoviePanel = lazy(() => import("./SceneMoviePanel"));
const SceneGalleriesPanel = lazy(() => import("./SceneGalleriesPanel"));
const DeleteScenesDialog = lazy(() => import("../DeleteScenesDialog"));
const GenerateDialog = lazy(() => import("../../Dialogs/GenerateDialog"));
const SceneVideoFilterPanel = lazy(() => import("./SceneVideoFilterPanel"));
import { objectPath, objectTitle } from "src/core/files";

interface IProps {
  scene: GQL.SceneDataFragment;
  refetch: () => void;
  setTimestamp: (num: number) => void;
  queueScenes: GQL.SceneDataFragment[];
  onQueueNext: () => void;
  onQueuePrevious: () => void;
  onQueueRandom: () => void;
  continuePlaylist: boolean;
  playScene: (sceneID: string, page?: number) => void;
  queueHasMoreScenes: () => boolean;
  onQueueMoreScenes: () => void;
  onQueueLessScenes: () => void;
  queueStart: number;
  collapsed: boolean;
  setCollapsed: (state: boolean) => void;
  setContinuePlaylist: (value: boolean) => void;
}

const ScenePage: React.FC<IProps> = ({
  scene,
  refetch,
  setTimestamp,
  queueScenes,
  onQueueNext,
  onQueuePrevious,
  onQueueRandom,
  continuePlaylist,
  playScene,
  queueHasMoreScenes,
  onQueueMoreScenes,
  onQueueLessScenes,
  queueStart,
  collapsed,
  setCollapsed,
  setContinuePlaylist,
}) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const [updateScene] = useSceneUpdate();
  const [generateScreenshot] = useSceneGenerateScreenshot();
  const { configuration } = useContext(ConfigurationContext);

  const [showDraftModal, setShowDraftModal] = useState(false);
  const boxes = configuration?.general?.stashBoxes ?? [];

  const [incrementO] = useSceneIncrementO(scene.id);
  const [decrementO] = useSceneDecrementO(scene.id);
  const [resetO] = useSceneResetO(scene.id);

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("scene-details-panel");

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);

  const onIncrementClick = async () => {
    try {
      await incrementO();
    } catch (e) {
      Toast.error(e);
    }
  };

  const onDecrementClick = async () => {
    try {
      await decrementO();
    } catch (e) {
      Toast.error(e);
    }
  };

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("scene-details-panel"));
    Mousetrap.bind("q", () => setActiveTabKey("scene-queue-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("scene-edit-panel"));
    Mousetrap.bind("k", () => setActiveTabKey("scene-markers-panel"));
    Mousetrap.bind("i", () => setActiveTabKey("scene-file-info-panel"));
    Mousetrap.bind("o", () => onIncrementClick());
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
      Mousetrap.unbind("o");
      Mousetrap.unbind("p n");
      Mousetrap.unbind("p p");
      Mousetrap.unbind("p r");
      Mousetrap.unbind(",");
    };
  });

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

  const onResetClick = async () => {
    try {
      await resetO();
    } catch (e) {
      Toast.error(e);
    }
  };

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  async function onRescan() {
    await mutateMetadataScan({
      paths: [objectPath(scene)],
    });

    Toast.success({
      content: intl.formatMessage(
        { id: "toast.rescanning_entity" },
        {
          count: 1,
          singularEntity: intl
            .formatMessage({ id: "scene" })
            .toLocaleLowerCase(),
        }
      ),
    });
  }

  async function onGenerateScreenshot(at?: number) {
    await generateScreenshot({
      variables: {
        id: scene.id,
        at,
      },
    });
    Toast.success({
      content: intl.formatMessage({ id: "toast.generating_screenshot" }),
    });
  }

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      history.push("/scenes");
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
        <Dropdown.Item
          key="rescan"
          className="bg-secondary text-white"
          onClick={() => onRescan()}
        >
          <FormattedMessage id="actions.rescan" />
        </Dropdown.Item>
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
          {(queueScenes ?? []).length > 0 ? (
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
          {scene.movies.length > 0 ? (
            <Nav.Item>
              <Nav.Link eventKey="scene-movie-panel">
                <FormattedMessage
                  id="countables.movies"
                  values={{ count: scene.movies.length }}
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
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link eventKey="scene-edit-panel">
              <FormattedMessage id="actions.edit" />
            </Nav.Link>
          </Nav.Item>
          <ButtonGroup className="ml-auto">
            <Nav.Item className="ml-auto">
              <ExternalPlayerButton scene={scene} />
            </Nav.Item>
            <Nav.Item className="ml-auto">
              <OCounterButton
                value={scene.o_counter || 0}
                onIncrement={onIncrementClick}
                onDecrement={onDecrementClick}
                onReset={onResetClick}
              />
            </Nav.Item>
            <Nav.Item>
              <OrganizedButton
                loading={organizedLoading}
                organized={scene.organized}
                onClick={onOrganizedClick}
              />
            </Nav.Item>
            <Nav.Item>{renderOperations()}</Nav.Item>
          </ButtonGroup>
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
            onSceneClicked={(sceneID) => playScene(sceneID)}
            onNext={onQueueNext}
            onPrevious={onQueuePrevious}
            onRandom={onQueueRandom}
            start={queueStart}
            hasMoreScenes={queueHasMoreScenes()}
            onLessScenes={() => onQueueLessScenes()}
            onMoreScenes={() => onQueueMoreScenes()}
          />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-markers-panel">
          <SceneMarkersPanel
            sceneId={scene.id}
            onClickMarker={onClickMarker}
            isVisible={activeTabKey === "scene-markers-panel"}
          />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-movie-panel">
          <SceneMoviePanel scene={scene} />
        </Tab.Pane>
        {scene.galleries.length === 1 && (
          <Tab.Pane eventKey="scene-galleries-panel">
            <GalleryViewer galleryId={scene.galleries[0].id} />
          </Tab.Pane>
        )}
        {scene.galleries.length > 1 && (
          <Tab.Pane eventKey="scene-galleries-panel">
            <SceneGalleriesPanel galleries={scene.galleries} />
          </Tab.Pane>
        )}
        <Tab.Pane eventKey="scene-video-filter-panel">
          <SceneVideoFilterPanel scene={scene} />
        </Tab.Pane>
        <Tab.Pane className="file-info-panel" eventKey="scene-file-info-panel">
          <SceneFileInfoPanel scene={scene} />
        </Tab.Pane>
        <Tab.Pane eventKey="scene-edit-panel">
          <SceneEditPanel
            isVisible={activeTabKey === "scene-edit-panel"}
            scene={scene}
            onDelete={() => setIsDeleteAlertOpen(true)}
            onUpdate={() => refetch()}
          />
        </Tab.Pane>
      </Tab.Content>
    </Tab.Container>
  );

  function getCollapseButtonText() {
    return collapsed ? ">" : "<";
  }

  const title = objectTitle(scene);

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
        <div className="d-none d-xl-block">
          {scene.studio && (
            <h1 className="text-center">
              <Link to={`/studios/${scene.studio.id}`}>
                <img
                  src={scene.studio.image_path ?? ""}
                  alt={`${scene.studio.name} logo`}
                  className="studio-logo"
                />
              </Link>
            </h1>
          )}
          <h3 className="scene-header">{title}</h3>
        </div>
        {renderTabs()}
      </div>
      <div className="scene-divider d-none d-xl-block">
        <Button
          onClick={() => {
            setCollapsed(!collapsed);
          }}
        >
          {getCollapseButtonText()}
        </Button>
      </div>
      <SubmitStashBoxDraft
        boxes={boxes}
        entity={scene}
        query={GQL.SubmitStashBoxSceneDraftDocument}
        show={showDraftModal}
        onHide={() => setShowDraftModal(false)}
      />
    </>
  );
};

const SceneLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const location = useLocation();
  const history = useHistory();
  const { configuration } = useContext(ConfigurationContext);
  const { data, loading, refetch } = useFindScene(id ?? "");
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [collapsed, setCollapsed] = useState(false);
  const [continuePlaylist, setContinuePlaylist] = useState(false);
  const [showScrubber, setShowScrubber] = useState(
    configuration?.interface.showScrubber ?? true
  );

  const sceneQueue = useMemo(
    () => SceneQueue.fromQueryParameters(location.search),
    [location.search]
  );
  const [queueScenes, setQueueScenes] = useState<GQL.SceneDataFragment[]>([]);

  const [queueTotal, setQueueTotal] = useState(0);
  const [queueStart, setQueueStart] = useState(1);

  const queryParams = useMemo(() => queryString.parse(location.search), [
    location.search,
  ]);

  function getInitialTimestamp() {
    const params = queryString.parse(location.search);
    const initialTimestamp = params?.t ?? "0";
    return Number.parseInt(
      Array.isArray(initialTimestamp) ? initialTimestamp[0] : initialTimestamp,
      10
    );
  }

  const autoplay = queryParams?.autoplay === "true";
  const currentQueueIndex = queueScenes
    ? queueScenes.findIndex((s) => s.id === id)
    : -1;

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind(".", () => setShowScrubber(!showScrubber));

    return () => {
      Mousetrap.unbind(".");
    };
  });

  useEffect(() => {
    // reset timestamp after notifying player
    if (timestamp !== -1) setTimestamp(-1);
  }, [timestamp]);

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
    const newScenes = scenes.concat(queueScenes);
    setQueueScenes(newScenes);
    setQueueStart(newStart);
  }

  function queueHasMoreScenes() {
    return queueStart + queueScenes.length - 1 < queueTotal;
  }

  async function onQueueMoreScenes() {
    if (!sceneQueue.query || !queueHasMoreScenes()) {
      return;
    }

    const filterCopy = sceneQueue.query.clone();
    const newStart = queueStart + queueScenes.length;
    filterCopy.currentPage = Math.ceil(newStart / filterCopy.itemsPerPage);
    const query = await queryFindScenes(filterCopy);
    const { scenes } = query.data.findScenes;

    // append scenes to scene list
    const newScenes = scenes.concat(queueScenes);
    setQueueScenes(newScenes);
    // don't change queue start
  }

  function playScene(sceneID: string, newPage?: number) {
    sceneQueue.playScene(history, sceneID, {
      newPage,
      autoPlay: true,
      continue: continuePlaylist,
    });
  }

  function onQueueNext() {
    if (!queueScenes) return;
    if (currentQueueIndex >= 0 && currentQueueIndex < queueScenes.length - 1) {
      playScene(queueScenes[currentQueueIndex + 1].id);
    }
  }

  function onQueuePrevious() {
    if (!queueScenes) return;
    if (currentQueueIndex > 0) {
      playScene(queueScenes[currentQueueIndex - 1].id);
    }
  }

  async function onQueueRandom() {
    if (!queueScenes) return;

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
        const { id: sceneID } = queryResults!.data!.findScenes!.scenes[index];
        // navigate to the image player page
        playScene(sceneID, page);
      }
    } else {
      const index = Math.floor(Math.random() * queueTotal);
      playScene(queueScenes[index].id);
    }
  }

  function onComplete() {
    // load the next scene if we're autoplaying
    if (continuePlaylist) {
      onQueueNext();
    }
  }

  /*
  if (error) return <ErrorMessage error={error.message} />;
  if (!loading && !data?.findScene)
    return <ErrorMessage error={`No scene found with id ${id}.`} />;
     */

  const scene = data?.findScene;

  return (
    <div className="row">
      {!loading && scene ? (
        <ScenePage
          scene={scene}
          refetch={refetch}
          setTimestamp={setTimestamp}
          queueScenes={queueScenes ?? []}
          queueStart={queueStart}
          onQueueNext={onQueueNext}
          onQueuePrevious={onQueuePrevious}
          onQueueRandom={onQueueRandom}
          continuePlaylist={continuePlaylist}
          playScene={playScene}
          queueHasMoreScenes={queueHasMoreScenes}
          onQueueLessScenes={onQueueLessScenes}
          onQueueMoreScenes={onQueueMoreScenes}
          collapsed={collapsed}
          setCollapsed={setCollapsed}
          setContinuePlaylist={setContinuePlaylist}
        />
      ) : (
        <div className="scene-tabs" />
      )}
      <div
        className={`scene-player-container ${collapsed ? "expanded" : ""} ${
          !showScrubber ? "hide-scrubber" : ""
        }`}
      >
        <ScenePlayer
          key="ScenePlayer"
          className="w-100 m-sm-auto no-gutter"
          scene={scene}
          timestamp={timestamp}
          autoplay={autoplay}
          onComplete={onComplete}
          onNext={
            currentQueueIndex >= 0 && currentQueueIndex < queueScenes.length - 1
              ? onQueueNext
              : undefined
          }
          onPrevious={currentQueueIndex > 0 ? onQueuePrevious : undefined}
        />
      </div>
    </div>
  );
};

export default SceneLoader;
