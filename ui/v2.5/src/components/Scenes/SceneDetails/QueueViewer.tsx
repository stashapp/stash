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
  onSceneClicked: (id: string) => void;
  onNext: () => void;
  onPrevious: () => void;
  onRandom: () => void;
}

export const QueueViewer: React.FC<IPlaylistViewer> = ({
  scenes,
  currentID,
  onNext,
  onPrevious,
  onRandom,
  onSceneClicked,
}) => {
  const currentIndex = scenes?.findIndex(s => s.id === currentID);

  function isCurrentScene(scene: GQL.SlimSceneDataFragment) {
    return scene.id === currentID;
  }

  function handleSceneClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>,
    id: string
  ) {
    onSceneClicked(id);
    event.preventDefault();
  }

  function renderPlaylistEntry(scene: GQL.SlimSceneDataFragment) {
    return (
      <li className={cx("my-2", { current: isCurrentScene(scene) })}>
        <Link to={`/scenes/${scene.id}`} onClick={(e) => handleSceneClick(e, scene.id)}>
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
          {currentIndex ?? 0 > 0 ? (
            <Button className="minimal" variant="secondary" onClick={() => onPrevious()}>
              <Icon icon="step-backward" />
            </Button>
          ) : ""}
          {(currentIndex ?? 0) < (scenes ?? []).length - 1 ? (
            <Button className="minimal" variant="secondary" onClick={() => onNext()}>
              <Icon icon="step-forward" />
            </Button>
          ) : ""}
          <Button className="minimal" variant="secondary" onClick={() => onRandom()}>
            <Icon icon="random" />
          </Button>
        </div>
      </div>
      <ol>{(scenes ?? []).map(renderPlaylistEntry)}</ol>
    </div>
  );
};
