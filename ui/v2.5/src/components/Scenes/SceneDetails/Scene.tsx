import { Tab, Tabs } from "react-bootstrap";
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { useParams, useLocation, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { LoadingIndicator } from "src/components/Shared";
import { ScenePlayer } from "src/components/ScenePlayer";
import { ScenePerformerPanel } from "./ScenePerformerPanel";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneDetailPanel } from "./SceneDetailPanel";

export const Scene: React.FC = () => {
  const { id = "new" } = useParams();
  const location = useLocation();
  const history = useHistory();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [scene, setScene] = useState<GQL.SceneDataFragment | undefined>();
  const { data, error, loading } = StashService.useFindScene(id);

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

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  if (loading || !scene || !data?.findScene) {
    return <LoadingIndicator />;
  }

  if (error) return <div>{error.message}</div>;

  return (
    <>
      <ScenePlayer scene={scene} timestamp={timestamp} autoplay={autoplay} />
      <div id="scene-details-container" className="col col-sm-9 m-sm-auto">
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
          <Tab
            eventKey="scene-edit-panel"
            title="Edit"
            tabClassName="d-none d-sm-block"
          >
            <SceneEditPanel
              scene={scene}
              onUpdate={newScene => setScene(newScene)}
              onDelete={() => history.push("/scenes")}
            />
          </Tab>
        </Tabs>
      </div>
    </>
  );
};
