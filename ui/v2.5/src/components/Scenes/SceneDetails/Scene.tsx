import { Tab, Tabs, Nav } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { useParams, useLocation, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  useFindScene,
  useSceneIncrementO,
  useSceneDecrementO,
  useSceneResetO,
} from "src/core/StashService";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { LoadingIndicator } from "src/components/Shared";
import { useToast } from "src/hooks";
import { ScenePlayer } from "src/components/ScenePlayer";
import { ScenePerformerPanel } from "./ScenePerformerPanel";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { OCounterButton } from "./OCounterButton";
import { SceneOperationsPanel } from "./SceneOperationsPanel";
import { SceneMoviePanel } from "./SceneMoviePanel";

export const Scene: React.FC = () => {
  const { id = "new" } = useParams();
  const location = useLocation();
  const history = useHistory();
  const Toast = useToast();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [scene, setScene] = useState<GQL.SceneDataFragment | undefined>();
  const { data, error, loading } = useFindScene(id);
  const [oLoading, setOLoading] = useState(false);
  const [incrementO] = useSceneIncrementO(scene?.id ?? "0");
  const [decrementO] = useSceneDecrementO(scene?.id ?? "0");
  const [resetO] = useSceneResetO(scene?.id ?? "0");

  const queryParams = queryString.parse(location.search);
  const autoplay = queryParams?.autoplay === "true";

  useEffect(() => {
    if (data?.findScene) setScene(data.findScene);
  }, [data]);

  function getInitialTimestamp() {
    const params = queryString.parse(location.search);
    const initialTimestamp = params?.t ?? "0";
    return Number.parseInt(
      Array.isArray(initialTimestamp) ? initialTimestamp[0] : initialTimestamp,
      10
    );
  }

  const updateOCounter = (newValue: number) => {
    const modifiedScene = { ...scene } as GQL.SceneDataFragment;
    modifiedScene.o_counter = newValue;
    setScene(modifiedScene);
  };

  const onIncrementClick = async () => {
    try {
      setOLoading(true);
      const result = await incrementO();
      if (result.data) updateOCounter(result.data.sceneIncrementO);
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  const onDecrementClick = async () => {
    try {
      setOLoading(true);
      const result = await decrementO();
      if (result.data) updateOCounter(result.data.sceneDecrementO);
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  const onResetClick = async () => {
    try {
      setOLoading(true);
      const result = await resetO();
      if (result.data) updateOCounter(result.data.sceneResetO);
    } catch (e) {
      Toast.error(e);
    } finally {
      setOLoading(false);
    }
  };

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  function renderTabs() {
    if (!scene) {
      return;
    }

    return (
      <Tabs id="scene-tabs" mountOnEnter>
        <Tab eventKey="scene-details-panel" title="Details">
          <SceneDetailPanel scene={scene} />
        </Tab>
        <Tab eventKey="scene-markers-panel" title="Markers">
          <SceneMarkersPanel scene={scene} onClickMarker={onClickMarker} />
        </Tab>
        {scene.performers.length > 0 ? (
          <Tab eventKey="scene-performer-panel" title="Performers">
            <ScenePerformerPanel scene={scene} />
          </Tab>
        ) : (
          ""
        )}
        {scene.movies.length > 0 ? (
          <Tab eventKey="scene-movie-panel" title="Movies">
            <SceneMoviePanel scene={scene} />
          </Tab>
        ) : (
          ""
        )}
        {scene.gallery ? (
          <Tab eventKey="scene-gallery-panel" title="Gallery">
            <GalleryViewer gallery={scene.gallery} />
          </Tab>
        ) : (
          ""
        )}
        <Tab
          className="file-info-panel"
          eventKey="scene-file-info-panel"
          title="File Info"
        >
          <SceneFileInfoPanel scene={scene} />
        </Tab>
        <Tab eventKey="scene-edit-panel" title="Edit">
          <SceneEditPanel
            scene={scene}
            onUpdate={(newScene) => setScene(newScene)}
            onDelete={() => history.push("/scenes")}
          />
        </Tab>
        <Tab eventKey="scene-operations-panel" title="Operations">
          <SceneOperationsPanel scene={scene} />
        </Tab>
      </Tabs>
    );

    /*return (
      <Tab.Container defaultActiveKey="scene-details-panel">
        <Tab.Content>
          <Tab.Pane eventKey="scene-details-panel" title="Details">
            <SceneDetailPanel scene={scene} />
          </Tab.Pane>
          <Tab.Pane eventKey="scene-markers-panel" title="Markers">
            <SceneMarkersPanel scene={scene} onClickMarker={onClickMarker} />
          </Tab.Pane>
          {scene.performers.length > 0 ? (
            <Tab.Pane eventKey="scene-performer-panel" title="Performers">
              <ScenePerformerPanel scene={scene} />
            </Tab.Pane>
          ) : (
            ""
          )}
          {scene.movies.length > 0 ? (
            <Tab.Pane eventKey="scene-movie-panel" title="Movies">
              <SceneMoviePanel scene={scene} />
            </Tab.Pane>
          ) : (
            ""
          )}
          {scene.gallery ? (
            <Tab.Pane eventKey="scene-gallery-panel" title="Gallery">
              <GalleryViewer gallery={scene.gallery} />
            </Tab.Pane>
          ) : (
            ""
          )}
          <Tab.Pane
            className="file-info-panel"
            eventKey="scene-file-info-panel"
            title="File Info"
          >
            <SceneFileInfoPanel scene={scene} />
          </Tab.Pane>
          <Tab.Pane eventKey="scene-edit-panel" title="Edit">
            <SceneEditPanel
              scene={scene}
              onUpdate={(newScene) => setScene(newScene)}
              onDelete={() => history.push("/scenes")}
            />
          </Tab.Pane>
          <Tab.Pane eventKey="scene-operations-panel" title="Operations">
            <SceneOperationsPanel scene={scene} />
          </Tab.Pane>
        </Tab.Content>
        <div>
          <Nav variant="tabs">
            <Nav.Item>
              <Nav.Link eventKey="scene-details-panel">Details</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="scene-markers-panel">Markers</Nav.Link>
            </Nav.Item>
          </Nav>
        </div>
      </Tab.Container>
    )*/
  }

  if (loading || !scene || !data?.findScene) {
    return <LoadingIndicator />;
  }

  if (error) return <div>{error.message}</div>;

  let layout = "default";
  layout = "compact";

  if (layout === "compact") {
    return (
      <div className="row">
        <div className="col col-sm-3">
          {renderTabs()}
        </div>
        <div className="col col-sm-9">
          <ScenePlayer className="w-100 m-sm-auto no-gutter" scene={scene} timestamp={timestamp} autoplay={autoplay} />
          <div className="float-right">
            <OCounterButton
              loading={oLoading}
              value={scene.o_counter || 0}
              onIncrement={onIncrementClick}
              onDecrement={onDecrementClick}
              onReset={onResetClick}
            />
          </div>
        </div>
      </div>
    )
  }

  else { //if (layout === "default") {
    return (
      <>
        <ScenePlayer scene={scene} timestamp={timestamp} autoplay={autoplay} />
        <div id="scene-details-container" className="col col-sm-9 m-sm-auto">
          <div className="float-right">
            <OCounterButton
              loading={oLoading}
              value={scene.o_counter || 0}
              onIncrement={onIncrementClick}
              onDecrement={onDecrementClick}
              onReset={onResetClick}
            />
          </div>
          {renderTabs()}
        </div>
      </>
    )
  };
};
