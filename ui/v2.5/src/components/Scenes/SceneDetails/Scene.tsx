import { Tab, Nav, Dropdown, Button, ButtonGroup } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { useParams, useLocation, useHistory, Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  useFindScene,
  useSceneIncrementO,
  useSceneDecrementO,
  useSceneResetO,
  useSceneStreams,
  useSceneGenerateScreenshot,
  useSceneUpdate
} from "src/core/StashService";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { ErrorMessage, LoadingIndicator, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import { ScenePlayer } from "src/components/ScenePlayer";
import { TextUtils, JWUtils } from "src/utils";
import * as Mousetrap from "mousetrap";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { OCounterButton } from "./OCounterButton";
import { SceneMoviePanel } from "./SceneMoviePanel";
import { DeleteScenesDialog } from "../DeleteScenesDialog";
import { SceneGenerateDialog } from "../SceneGenerateDialog";
import { SceneVideoFilterPanel } from "./SceneVideoFilterPanel";
import { OrganizedButton } from "./OrganizedButton";

interface ISceneParams {
  id?: string;
}

export const Scene: React.FC = () => {
  const { id = "new" } = useParams<ISceneParams>();
  const location = useLocation();
  const history = useHistory();
  const Toast = useToast();
  const [updateScene] = useSceneUpdate();
  const [generateScreenshot] = useSceneGenerateScreenshot();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [collapsed, setCollapsed] = useState(false);

  const { data, error, loading } = useFindScene(id);
  const scene = data?.findScene;
  const {
    data: sceneStreams,
    error: streamableError,
    loading: streamableLoading,
  } = useSceneStreams(id);
  const [oLoading, setOLoading] = useState(false);
  const [incrementO] = useSceneIncrementO(scene?.id ?? "0");
  const [decrementO] = useSceneDecrementO(scene?.id ?? "0");
  const [resetO] = useSceneResetO(scene?.id ?? "0");

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("scene-details-panel");

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [isGenerateDialogOpen, setIsGenerateDialogOpen] = useState(false);

  const queryParams = queryString.parse(location.search);
  const autoplay = queryParams?.autoplay === "true";

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
          id: scene?.id ?? "",
          organized: !scene?.organized,
        }
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

  async function onGenerateScreenshot(at?: number) {
    if (!scene) {
      return;
    }

    await generateScreenshot({
      variables: {
        id: scene.id,
        at,
      },
    });
    Toast.success({ content: "Generating screenshot" });
  }

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      history.push("/scenes");
    }
  }

  function maybeRenderDeleteDialog() {
    if (isDeleteAlertOpen && scene) {
      return (
        <DeleteScenesDialog selected={[scene]} onClose={onDeleteDialogClosed} />
      );
    }
  }

  function maybeRenderSceneGenerateDialog() {
    if (isGenerateDialogOpen && scene) {
      return (
        <SceneGenerateDialog
          selectedIds={[scene.id]}
          onClose={() => {
            setIsGenerateDialogOpen(false);
          }}
        />
      );
    }
  }

  function renderOperations() {
    return (
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
            key="generate"
            className="bg-secondary text-white"
            onClick={() => setIsGenerateDialogOpen(true)}
          >
            Generate...
          </Dropdown.Item>
          <Dropdown.Item
            key="generate-screenshot"
            className="bg-secondary text-white"
            onClick={() =>
              onGenerateScreenshot(JWUtils.getPlayer().getPosition())
            }
          >
            Generate thumbnail from current
          </Dropdown.Item>
          <Dropdown.Item
            key="generate-default"
            className="bg-secondary text-white"
            onClick={() => onGenerateScreenshot()}
          >
            Generate default thumbnail
          </Dropdown.Item>
          <Dropdown.Item
            key="delete-scene"
            className="bg-secondary text-white"
            onClick={() => setIsDeleteAlertOpen(true)}
          >
            Delete Scene
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  function renderTabs() {
    if (!scene) {
      return;
    }

    return (
      <Tab.Container
        activeKey={activeTabKey}
        onSelect={(k) => k && setActiveTabKey(k)}
      >
        <div>
          <Nav variant="tabs" className="mr-auto">
            <Nav.Item>
              <Nav.Link eventKey="scene-details-panel">Details</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="scene-markers-panel">Markers</Nav.Link>
            </Nav.Item>
            {scene.movies.length > 0 ? (
              <Nav.Item>
                <Nav.Link eventKey="scene-movie-panel">Movies</Nav.Link>
              </Nav.Item>
            ) : (
              ""
            )}
            {scene.gallery ? (
              <Nav.Item>
                <Nav.Link eventKey="scene-gallery-panel">Gallery</Nav.Link>
              </Nav.Item>
            ) : (
              ""
            )}
            <Nav.Item>
              <Nav.Link eventKey="scene-video-filter-panel">Filters</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="scene-file-info-panel">File Info</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="scene-edit-panel">Edit</Nav.Link>
            </Nav.Item>
            <ButtonGroup className="ml-auto">
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
          {scene.gallery ? (
            <Tab.Pane eventKey="scene-gallery-panel">
              <GalleryViewer gallery={scene.gallery} />
            </Tab.Pane>
          ) : (
            ""
          )}
          <Tab.Pane eventKey="scene-video-filter-panel">
            <SceneVideoFilterPanel scene={scene} />
          </Tab.Pane>
          <Tab.Pane
            className="file-info-panel"
            eventKey="scene-file-info-panel"
          >
            <SceneFileInfoPanel scene={scene} />
          </Tab.Pane>
          <Tab.Pane eventKey="scene-edit-panel">
            <SceneEditPanel
              isVisible={activeTabKey === "scene-edit-panel"}
              scene={scene}
              onDelete={() => setIsDeleteAlertOpen(true)}
            />
          </Tab.Pane>
        </Tab.Content>
      </Tab.Container>
    );
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("scene-details-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("scene-edit-panel"));
    Mousetrap.bind("k", () => setActiveTabKey("scene-markers-panel"));
    Mousetrap.bind("f", () => setActiveTabKey("scene-file-info-panel"));
    Mousetrap.bind("o", () => onIncrementClick());

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("k");
      Mousetrap.unbind("f");
      Mousetrap.unbind("o");
    };
  });

  function getCollapseButtonText() {
    return collapsed ? ">" : "<";
  }

  if (loading || streamableLoading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (streamableError) return <ErrorMessage error={streamableError.message} />;
  if (!scene) return <ErrorMessage error={`No scene found with id ${id}.`} />;

  return (
    <div className="row">
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
        <ScenePlayer
          className="w-100 m-sm-auto no-gutter"
          scene={scene}
          timestamp={timestamp}
          autoplay={autoplay}
          sceneStreams={sceneStreams?.sceneStreams ?? []}
        />
      </div>
    </div>
  );
};
