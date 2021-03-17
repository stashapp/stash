import React from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";

export interface IPlaylistViewer {
  scenes?: GQL.SlimSceneDataFragment[];
  currentID?: string;
}

export const QueueViewer: React.FC<IPlaylistViewer> = ({
  scenes,
  currentID,
}) => {
  function isCurrentScene(scene: GQL.SlimSceneDataFragment) {
    return scene.id === currentID;
  }

  function renderPlaylistEntry(scene: GQL.SlimSceneDataFragment) {
    return (
      <li className={cx("my-2", { current: isCurrentScene(scene) })}>
        <Link to={`/scenes/${scene.id}`}>
          <div className="ml-1 d-flex align-items-center">
            <div className="thumbnail-container">
              <img alt={scene.title ?? ""} src={scene.paths.screenshot ?? ""} />
            </div>
            <div>
              <span className="align-middle">
                {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
              </span>
            </div>
          </div>
        </Link>
      </li>
    );
  }

  return (
    <div id="playlist-viewer">
      <div id="playlist-controls" className="d-flex justify-content-end">
        <div>
          <Button className="minimal" variant="secondary" >
            <Icon icon="step-backward" />
          </Button>
          <Button className="minimal" variant="secondary">
            <Icon icon="step-forward" />
          </Button>
          <Button className="minimal" variant="secondary">
            <Icon icon="random" />
          </Button>
        </div>
      </div>
      <ol>{(scenes ?? []).map(renderPlaylistEntry)}</ol>
    </div>
  );
};
