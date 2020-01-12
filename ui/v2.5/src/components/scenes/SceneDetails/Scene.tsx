import { Card, Spinner, Tab, Tabs } from 'react-bootstrap';
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import { useParams, useLocation, useHistory } from 'react-router-dom';
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { ScenePlayer } from "../ScenePlayer/ScenePlayer";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { ScenePerformerPanel } from "./ScenePerformerPanel";

export const Scene: React.FC = () => {
  const { id = 'new' } = useParams();
  const location = useLocation();
  const history = useHistory();
  const [timestamp, setTimestamp] = useState<number>(getInitialTimestamp());
  const [scene, setScene] = useState<Partial<GQL.SceneDataFragment>>({});
  const { data, error, loading } = StashService.useFindScene(id);

  const queryParams = queryString.parse(location.search);
  const autoplay = queryParams?.autoplay === 'true';

  useEffect(() => (
    setScene(data?.findScene ?? {})
  ), [data]);

  function getInitialTimestamp() {
    const params = queryString.parse(location.search);
    const timestamp = params?.t;
    return Number.parseInt(Array.isArray(timestamp) ? timestamp[0] : timestamp ?? '0', 10);
  }

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  if (!data?.findScene || loading || Object.keys(scene).length === 0) {
    return <Spinner animation="border"/>;
  }

  if (error)
    return <div>{error.message}</div>

  const modifiedScene =
    Object.assign({scene_marker_tags: data.sceneMarkerTags}, scene) as GQL.SceneDataFragment; // TODO Hack from angular

  return (
    <>
      <ScenePlayer scene={modifiedScene} timestamp={timestamp} autoplay={autoplay}/>
      <Card id="details-container">
        <Tabs id="scene-tabs" mountOnEnter={true}>
            <Tab eventKey="scene-details-panel" title="Details">
              <SceneDetailPanel scene={modifiedScene} />
            </Tab>
            <Tab
              eventKey="scene-markers-panel"
              title="Markers">
              <SceneMarkersPanel scene={modifiedScene} onClickMarker={onClickMarker} />
            </Tab>
            {modifiedScene.performers.length > 0 ?
              <Tab
                eventKey="scene-performer-panel"
                title="Performers">
                <ScenePerformerPanel scene={modifiedScene} />
              </Tab> : ''
            }
            {!!modifiedScene.gallery ?
              <Tab
                eventKey="scene-gallery-panel"
                title="Gallery">
                <GalleryViewer gallery={modifiedScene.gallery} />
              </Tab> : ''
            }
            <Tab eventKey="scene-file-info-panel" title="File Info">
                <SceneFileInfoPanel scene={modifiedScene} />
            </Tab>
            <Tab
              eventKey="scene-edit-panel"
              title="Edit">
              <SceneEditPanel
                scene={modifiedScene}
                onUpdate={(newScene) => setScene(newScene)}
                onDelete={() => history.push("/scenes")}
              />
            </Tab>
        </Tabs>
      </Card>
    </>
  );
};
