import { Tab, Nav, Dropdown, Button, ButtonGroup } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState, useMemo } from "react";
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
  useSceneStreams,
  useSceneGenerateScreenshot,
  useSceneUpdate,
  queryFindScenes,
  queryFindScenesByID,
} from "src/core/StashService";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { ErrorMessage, LoadingIndicator, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import { ScenePlayer } from "src/components/ScenePlayer";
import { TextUtils, JWUtils } from "src/utils";
import Mousetrap from "mousetrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SceneQueue } from "src/models/sceneQueue";
import { QueueViewer } from "./QueueViewer";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { OCounterButton } from "./OCounterButton";
import { ExternalPlayerButton } from "./ExternalPlayerButton";
import { SceneMoviePanel } from "./SceneMoviePanel";
import { SceneGalleriesPanel } from "./SceneGalleriesPanel";
import { DeleteScenesDialog } from "../DeleteScenesDialog";
import { GenerateDialog } from "../../Dialogs/GenerateDialog";
import { SceneVideoFilterPanel } from "./SceneVideoFilterPanel";
import { OrganizedButton } from "./OrganizedButton";

interface IProps {
  scene: GQL.SceneDataFragment;
  refetch: () => void;
}

const ScenePage: React.FC<IProps> = ({ scene, refetch }) => {
  const location = useLocation();
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const [updateScene] = useSceneUpdate();
  const [generateScreenshot] = useSceneGenerateScreenshot();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [collapsed, setCollapsed] = useState(false);
  const [showScrubber, setShowScrubber] = useState(true);

  const {
    data: sceneStreams,
    error: streamableError,
    loading: streamableLoading,
  } = useSceneStreams(scene.id);

  const [oLoading, setOLoading] = useState(false);
  const [incrementO] = useSceneIncrementO(scene.id);
  const [decrementO] = useSceneDecrementO(scene.id);
  const [resetO] = useSceneResetO(scene.id);

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("scene-details-panel");

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);

  const [sceneQueue, setSceneQueue] = useState<SceneQueue>(new SceneQueue());
  const [queueScenes, setQueueScenes] = useState<GQL.SlimSceneDataFragment[]>(
    []
  );

  const [queueTotal, setQueueTotal] = useState(0);
  const [queueStart, setQueueStart] = useState(1);
  const [continuePlaylist, setContinuePlaylist] = useState(false);

  const [rerenderPlayer, setRerenderPlayer] = useState(false);

  const queryParams = useMemo(() => queryString.parse(location.search), [
    location.search,
  ]);
  const autoplay = queryParams?.autoplay === "true";
  const currentQueueIndex = queueScenes.findIndex((s) => s.id === scene.id);

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
    setContinuePlaylist(queryParams?.continue === "true");
  }, [queryParams]);

  // HACK - jwplayer doesn't handle re-rendering when scene changes, so force
  // a rerender by not drawing it
  useEffect(() => {
    if (rerenderPlayer) {
      setRerenderPlayer(false);
    }
  }, [rerenderPlayer]);

  useEffect(() => {
    setRerenderPlayer(true);
  }, [scene.id]);

  useEffect(() => {
    setSceneQueue(SceneQueue.fromQueryParameters(location.search));
  }, [location.search]);

  useEffect(() => {
    if (sceneQueue.query) {
      getQueueFilterScenes(sceneQueue.query);
    } else if (sceneQueue.sceneIDs) {
      getQueueScenes(sceneQueue.sceneIDs);
    }
  }, [sceneQueue]);

  function getInitialTimestamp() {
    const params = queryString.parse(location.search);
    const initialTimestamp = params?.t ?? "0";
    return Number.parseInt(
      Array.isArray(initialTimestamp) ? initialTimestamp[0] : initialTimestamp,
      10
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

  const onIncrementClick = async () => {
    try {
      setOLoading(true);
      await incrementO();
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  const onDecrementClick = async () => {
    try {
      setOLoading(true);
      await decrementO();
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  const onResetClick = async () => {
    try {
      setOLoading(true);
      await resetO();
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  async function onRescan() {
    await mutateMetadataScan({
      paths: [scene.path],
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

  function playScene(sceneID: string, page?: number) {
    sceneQueue.playScene(history, sceneID, {
      newPage: page,
      autoPlay: true,
      continue: continuePlaylist,
    });
  }

  function onQueueNext() {
    if (currentQueueIndex >= 0 && currentQueueIndex < queueScenes.length - 1) {
      playScene(queueScenes[currentQueueIndex + 1].id);
    }
  }

  function onQueuePrevious() {
    if (currentQueueIndex > 0) {
      playScene(queueScenes[currentQueueIndex - 1].id);
    }
  }

  async function onQueueRandom() {
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
        title="Operations"
      >
        <Icon icon="ellipsis-v" />
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
          onClick={() =>
            onGenerateScreenshot(JWUtils.getPlayer().getPosition())
          }
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
                loading={oLoading}
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
            setContinue={(v) => setContinuePlaylist(v)}
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
            scene={scene}
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
    Mousetrap.bind(".", () => setShowScrubber(!showScrubber));

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
      Mousetrap.unbind(".");
    };
  });

  function getCollapseButtonText() {
    return collapsed ? ">" : "<";
  }

  if (streamableLoading) return <LoadingIndicator />;
  if (streamableError) return <ErrorMessage error={streamableError.message} />;

  return (
    <div className="row">
      <Helmet>
        <title>{scene.title ?? TextUtils.fileNameFromPath(scene.path)}</title>
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
          <h3 className="scene-header">
            {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
          </h3>
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
      <div className={`scene-player-container ${collapsed ? "expanded" : ""}`}>
        {!rerenderPlayer ? (
          <ScenePlayer
            className={`w-100 m-sm-auto no-gutter ${
              !showScrubber ? "hide-scrubber" : ""
            }`}
            scene={scene}
            timestamp={timestamp}
            autoplay={autoplay}
            sceneStreams={sceneStreams?.sceneStreams ?? []}
            onComplete={onComplete}
          />
        ) : undefined}
      </div>
    </div>
  );
};

const SceneLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const { data, loading, error, refetch } = useFindScene(id ?? "");

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findScene)
    return <ErrorMessage error={`No scene found with id ${id}.`} />;

  return <ScenePage scene={data.findScene} refetch={refetch} />;
};

export default SceneLoader;
