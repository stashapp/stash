import { Card, Spinner, Tab, Tabs } from 'react-bootstrap';
import queryString from "query-string";
import React, { useEffect, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { IBaseProps } from "src/models";
import { GalleryViewer } from "src/components/Galleries/GalleryViewer";
import { ScenePlayer } from "../ScenePlayer/ScenePlayer";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { ScenePerformerPanel } from "./ScenePerformerPanel";

interface ISceneProps extends IBaseProps {}

export const Scene: React.FC<ISceneProps> = (props: ISceneProps) => {
  const [timestamp, setTimestamp] = useState<number>(0);
  const [autoplay, setAutoplay] = useState<boolean>(false);
  const [scene, setScene] = useState<Partial<GQL.SceneDataFragment>>({});
  const { data, error, loading } = StashService.useFindScene(props.match.params.id);

  useEffect(() => (
    setScene(data?.findScene ?? {})
  ), [data]);

  useEffect(() => {
    const queryParams = queryString.parse(props.location.search);
    if (!!queryParams.t && typeof queryParams.t === "string" && timestamp === 0) {
      const newTimestamp = parseInt(queryParams.t, 10);
      setTimestamp(newTimestamp);
    }
    if (queryParams.autoplay && typeof queryParams.autoplay === "string") {
      setAutoplay(queryParams.autoplay === "true");
    }
  });

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
                onDelete={() => props.history.push("/scenes")}
              />
            </Tab>
        </Tabs>
      </Card>
    </>
  );
};
