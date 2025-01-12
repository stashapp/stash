import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import { Button, Form, Spinner } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { useIntl } from "react-intl";
import {
  faChevronDown,
  faChevronUp,
  faRandom,
  faStepBackward,
  faStepForward,
} from "@fortawesome/free-solid-svg-icons";
import { objectTitle } from "src/core/files";
import {
  QueuedScene,
  QueuedSceneMarker,
  QueuedItem,
} from "src/models/sceneQueue";
import { markerTitle } from "src/core/markers";
import TextUtils from "src/utils/text";

export interface IPlaylistViewer {
  scenes: QueuedItem[];
  currentID?: string;
  currentMarkerSeconds?: number;
  start?: number;
  continue?: boolean;
  hasMoreScenes: boolean;
  setContinue: (v: boolean) => void;
  onSceneClicked: (scene: QueuedItem) => void;
  onNext: () => void;
  onPrevious: () => void;
  onRandom: () => void;
  onMoreScenes: () => void;
  onLessScenes: () => void;
}

export const QueueViewer: React.FC<IPlaylistViewer> = ({
  scenes,
  currentID,
  currentMarkerSeconds,
  start = 0,
  continue: continuePlaylist = false,
  hasMoreScenes,
  setContinue,
  onNext,
  onPrevious,
  onRandom,
  onSceneClicked,
  onMoreScenes,
  onLessScenes,
}) => {
  const intl = useIntl();
  const [lessLoading, setLessLoading] = useState(false);
  const [moreLoading, setMoreLoading] = useState(false);

  const currentIndex = scenes.findIndex((s) => {
    if (s.__typename === "SceneMarker") {
      return (
        s.scene.id === currentID &&
        Math.trunc(s.seconds) === currentMarkerSeconds
      );
    } else {
      return s.id === currentID;
    }
  });

  useEffect(() => {
    setLessLoading(false);
    setMoreLoading(false);
  }, [scenes]);

  function isCurrentScene(scene: QueuedScene) {
    return scene.id === currentID;
  }

  function isCurrentMarker(marker: QueuedSceneMarker) {
    return (
      marker.scene.id === currentID &&
      Math.trunc(marker.seconds) === currentMarkerSeconds
    );
  }

  function handleSceneClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>,
    scene: QueuedItem
  ) {
    onSceneClicked(scene);
    event.preventDefault();
  }

  function lessClicked() {
    setLessLoading(true);
    onLessScenes();
  }

  function moreClicked() {
    setMoreLoading(true);
    onMoreScenes();
  }

  function renderPlaylistEntryScene(scene: QueuedScene) {
    return (
      <li
        className={cx("my-2", { current: isCurrentScene(scene) })}
        key={scene.id}
      >
        <Link
          to={`/scenes/${scene.id}`}
          onClick={(e) => handleSceneClick(e, scene)}
        >
          <div className="ml-1 d-flex align-items-center">
            <div className="thumbnail-container">
              <img
                loading="lazy"
                alt={scene.title ?? ""}
                src={scene.paths.screenshot ?? ""}
              />
            </div>
            <div className="queue-scene-details">
              <span className="queue-marker-title">{objectTitle(scene)}</span>
              <span className="queue-marker-scene">{scene?.studio?.name}</span>
              <span className="queue-marker-performers">
                {scene?.performers
                  ?.map(function (performer) {
                    return performer.name;
                  })
                  .join(", ")}
              </span>
              <span className="queue-scene-date">{scene?.date}</span>
            </div>
          </div>
        </Link>
      </li>
    );
  }

  function renderPlaylistEntryMarker(marker: QueuedSceneMarker) {
    const tags = [marker.primary_tag, ...(marker.tags || [])];

    return (
      <li
        className={cx("my-2", { current: isCurrentMarker(marker) })}
        key={marker.id}
      >
        <Link
          to={`/scenes/${marker.scene.id}?t=${marker.seconds}`}
          onClick={(e) => handleSceneClick(e, marker)}
        >
          <div className="ml-1 d-flex align-items-center">
            <div className="thumbnail-container">
              <img loading="lazy" src={marker.screenshot ?? ""} />
            </div>
            <div className="queue-marker-details">
              <span className="queue-scene-title">
                {markerTitle(marker)}
                {" - "}
                {TextUtils.formatTimestampRange(
                  marker.seconds,
                  marker.end_seconds ?? undefined
                )}
              </span>
              <span className="queue-scene-studio">
                {objectTitle(marker.scene)}
              </span>
              <span className="queue-scene-performers">
                {marker.scene.performers
                  ?.map(function (performer) {
                    return performer.name;
                  })
                  .join(", ")}
              </span>
              <span className="queue-marker-tags">
                {tags.map((tag) => tag.name).join(", ")}
              </span>
            </div>
          </div>
        </Link>
      </li>
    );
  }

  return (
    <div id="queue-viewer">
      <div className="queue-controls">
        <div>
          <Form.Check
            id="continue-checkbox"
            checked={continuePlaylist}
            label={intl.formatMessage({ id: "actions.continue" })}
            onChange={() => {
              setContinue(!continuePlaylist);
            }}
          />
        </div>
        <div>
          {currentIndex > 0 || start > 1 ? (
            <Button
              className="minimal"
              variant="secondary"
              onClick={() => onPrevious()}
            >
              <Icon icon={faStepBackward} />
            </Button>
          ) : (
            ""
          )}
          {currentIndex < scenes.length - 1 || hasMoreScenes ? (
            <Button
              className="minimal"
              variant="secondary"
              onClick={() => onNext()}
            >
              <Icon icon={faStepForward} />
            </Button>
          ) : (
            ""
          )}
          <Button
            className="minimal"
            variant="secondary"
            onClick={() => onRandom()}
          >
            <Icon icon={faRandom} />
          </Button>
        </div>
      </div>
      <div id="queue-content">
        {start > 1 ? (
          <div className="d-flex justify-content-center">
            <Button onClick={() => lessClicked()} disabled={lessLoading}>
              {!lessLoading ? (
                <Icon icon={faChevronUp} />
              ) : (
                <Spinner animation="border" role="status" />
              )}
            </Button>
          </div>
        ) : undefined}
        <ol start={start}>
          {scenes.map((item) =>
            item.__typename == "Scene"
              ? renderPlaylistEntryScene(item)
              : item.__typename == "SceneMarker"
              ? renderPlaylistEntryMarker(item)
              : ""
          )}
        </ol>
        {hasMoreScenes ? (
          <div className="d-flex justify-content-center">
            <Button onClick={() => moreClicked()} disabled={moreLoading}>
              {!moreLoading ? (
                <Icon icon={faChevronDown} />
              ) : (
                <Spinner animation="border" role="status" />
              )}
            </Button>
          </div>
        ) : undefined}
      </div>
    </div>
  );
};

export default QueueViewer;
