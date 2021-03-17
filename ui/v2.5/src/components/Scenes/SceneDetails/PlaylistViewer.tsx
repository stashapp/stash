import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export interface IPlaylistViewer {
  scenes?: GQL.SlimSceneDataFragment[];
}

export const PlaylistViewer: React.FC<IPlaylistViewer> = ({ scenes }) => {
  function renderPlaylistEntry(scene: GQL.SlimSceneDataFragment) {
    return (
      <li className="my-3">
        <div className="ml-1 d-flex align-items-center">
          <div className="thumbnail-container">
            <Link to={`/scenes/${scene.id}`}>
              <img alt={scene.title ?? ""} src={scene.paths.screenshot ?? ""} />
            </Link>
          </div>
          <div>
            <Link className="align-middle" to={`/scenes/${scene.id}`}>
              {/* <TruncatedText
                text={scene.title ?? TextUtils.fileNameFromPath(scene.path)}
              /> */}
              {scene.title ?? TextUtils.fileNameFromPath(scene.path)}
            </Link>
          </div>
        </div>
      </li>
    );
  }

  return (
    <div id="playlist-viewer">
      <ol>{(scenes ?? []).map(renderPlaylistEntry)}</ol>
    </div>
  );
};
