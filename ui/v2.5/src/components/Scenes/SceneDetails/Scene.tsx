import { Tab, Tabs } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { useParams, useLocation, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  useFindScene,
  useSceneIncrementO,
  useSceneDecrementO,
  useSceneResetO,
  useIsSceneStreamable,
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
import jwplayer from "src/utils/jwplayer";

export const Scene: React.FC = () => {
  const { id = "new" } = useParams();
  const location = useLocation();
  const history = useHistory();
  const Toast = useToast();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [scene, setScene] = useState<GQL.SceneDataFragment | undefined>();
  const { data, error, loading } = useFindScene(id);
  const { data: isSceneStreamable, error: streamableError, loading: streamableLoading } = useIsSceneStreamable(id, jwplayer.getSupportedFormats());
  const [oLoading, setOLoading] = useState(false);
  const [incrementO] = useSceneIncrementO(scene?.id ?? "0");
  const [decrementO] = useSceneDecrementO(scene?.id ?? "0");
  const [resetO] = useSceneResetO(scene?.id ?? "0");

  const queryParams = queryString.parse(location.search);
  const autoplay = queryParams?.autoplay === "true";
  const isStreamable = isSceneStreamable?.isSceneStreamable ?? false;

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

  if (loading || streamableLoading || !scene || !data?.findScene) {
    return <LoadingIndicator />;
  }

  if (error) return <div>{error.message}</div>;
  if (streamableError) return <div>{streamableError.message}</div>;

  return (
    <>
      <ScenePlayer scene={scene} timestamp={timestamp} autoplay={autoplay} streamable={isStreamable}/>
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
      </div>
    </>
  );
};
