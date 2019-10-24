import {
  Card,
  Spinner,
  Tab,
  Tabs,
} from "@blueprintjs/core";
import queryString from "query-string";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { IBaseProps } from "../../../models";
import { GalleryViewer } from "../../Galleries/GalleryViewer";
import { ScenePlayer } from "../ScenePlayer/ScenePlayer";
import { SceneDetailPanel } from "./SceneDetailPanel";
import { SceneEditPanel } from "./SceneEditPanel";
import { SceneFileInfoPanel } from "./SceneFileInfoPanel";
import { SceneMarkersPanel } from "./SceneMarkersPanel";
import { ScenePerformerPanel } from "./ScenePerformerPanel";

interface ISceneProps extends IBaseProps {}

export const Scene: FunctionComponent<ISceneProps> = (props: ISceneProps) => {
  const [timestamp, setTimestamp] = useState<number>(0);
  const [scene, setScene] = useState<Partial<GQL.SceneDataFragment>>({});
  const [isLoading, setIsLoading] = useState(false);
  const { data, error, loading } = StashService.useFindScene(props.match.params.id);

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findScene || !!error) { return; }
    setScene(StashService.nullToUndefined(data.findScene));
  }, [data]);

  useEffect(() => {
    const queryParams = queryString.parse(props.location.search);
    if (!!queryParams.t && typeof queryParams.t === "string" && timestamp === 0) {
      const newTimestamp = parseInt(queryParams.t, 10);
      setTimestamp(newTimestamp);
    }
  });

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    setTimestamp(marker.seconds);
  }

  if (!data || !data.findScene || isLoading || Object.keys(scene).length === 0) {
    return <Spinner size={Spinner.SIZE_LARGE} />;
  }
  const modifiedScene =
    Object.assign({scene_marker_tags: data.sceneMarkerTags}, scene) as GQL.SceneDataFragment; // TODO Hack from angular
  if (!!error) { return <>error...</>; }

  return (
    <>
      <ScenePlayer scene={modifiedScene} timestamp={timestamp} />
      <Card id="details-container">
        <Tabs
          renderActiveTabPanelOnly={true}
          large={true}
        >
            <Tab id="scene-details-panel" title="Details" panel={<SceneDetailPanel scene={modifiedScene} />} />
            <Tab
              id="scene-markers-panel"
              title="Markers"
              panel={<SceneMarkersPanel scene={modifiedScene} onClickMarker={onClickMarker} />}
            />
            {modifiedScene.performers.length > 0 ?
              <Tab
                id="scene-performer-panel"
                title="Performers"
                panel={<ScenePerformerPanel scene={modifiedScene} />}
              /> : undefined
            }
            {!!modifiedScene.gallery ?
              <Tab
                id="scene-gallery-panel"
                title="Gallery"
                panel={<GalleryViewer gallery={modifiedScene.gallery} />}
              /> : undefined
            }
            <Tab id="scene-file-info-panel" title="File Info" panel={<SceneFileInfoPanel scene={modifiedScene} />} />
            <Tab
              id="scene-edit-panel"
              title="Edit"
              panel={
                <SceneEditPanel 
                  scene={modifiedScene} 
                  onUpdate={(newScene) => setScene(newScene)} 
                  onDelete={() => props.history.push("/scenes")}
                />}
            />
        </Tabs>
      </Card>
    </>
  );
};
