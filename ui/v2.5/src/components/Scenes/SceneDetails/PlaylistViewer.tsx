import React from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export interface IPlaylistViewer {
  scenes?: GQL.SlimSceneDataFragment[];
  currentID?: string;
}

export const PlaylistViewer: React.FC<IPlaylistViewer> = ({
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
      <ol>{(scenes ?? []).map(renderPlaylistEntry)}</ol>
    </div>
  );
};
