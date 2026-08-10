import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import { Button, Form, Spinner } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { useIntl } from "react-intl";
import {
  faChevronDown,
  faChevronUp,
  faMicrophone,
  faRandom,
  faStepBackward,
  faStepForward,
} from "@fortawesome/free-solid-svg-icons";
import { objectTitle } from "src/core/files";
import { QueuedAudio } from "src/models/audioQueue";

export interface IAudioPlaylistViewer {
  audios: QueuedAudio[];
  currentID?: string;
  start?: number;
  continue?: boolean;
  hasMoreAudios: boolean;
  setContinue: (v: boolean) => void;
  onAudioClicked: (id: string) => void;
  onNext: () => void;
  onPrevious: () => void;
  onRandom: () => void;
  onMoreAudios: () => void;
  onLessAudios: () => void;
}

export const AudioQueueViewer: React.FC<IAudioPlaylistViewer> = ({
  audios,
  currentID,
  start = 0,
  continue: continuePlaylist = false,
  hasMoreAudios,
  setContinue,
  onNext,
  onPrevious,
  onRandom,
  onAudioClicked,
  onMoreAudios,
  onLessAudios,
}) => {
  const intl = useIntl();
  const [lessLoading, setLessLoading] = useState(false);
  const [moreLoading, setMoreLoading] = useState(false);

  const currentIndex = audios.findIndex((a) => a.id === currentID);

  // HACK - this isn't a great way to handle loading state. Requires a larger refactor.
  // TODO - refactor this
  // biome-ignore lint/correctness/useExhaustiveDependencies: explicitly want to set loading to false when audios change
  useEffect(() => {
    setLessLoading(false);
    setMoreLoading(false);
  }, [audios]);

  function isCurrentAudio(audio: QueuedAudio) {
    return audio.id === currentID;
  }

  function handleAudioClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>,
    id: string
  ) {
    onAudioClicked(id);
    event.preventDefault();
  }

  function renderPlaylistEntry(audio: QueuedAudio) {
    return (
      <li
        className={cx("my-2", { current: isCurrentAudio(audio) })}
        key={audio.id}
      >
        <Link
          to={`/audios/${audio.id}`}
          onClick={(e) => handleAudioClick(e, audio.id)}
        >
          <div className="ml-1 d-flex align-items-center">
            <div className="thumbnail-container">
              {audio.paths.screenshot ? (
                <img
                  loading="lazy"
                  alt={audio.title ?? ""}
                  src={audio.paths.screenshot}
                />
              ) : (
                <div className="queue-audio-placeholder">
                  <Icon icon={faMicrophone} />
                </div>
              )}
            </div>
            <div className="queue-audio-details">
              <span className="queue-audio-title">{objectTitle(audio)}</span>
              <span className="queue-audio-studio">{audio?.studio?.name}</span>
              <span className="queue-audio-performers">
                {audio?.performers
                  ?.map((performer) => {
                    return performer.name;
                  })
                  .join(", ")}
              </span>
              <span className="queue-audio-date">{audio?.date}</span>
            </div>
          </div>
        </Link>
      </li>
    );
  }

  function lessClicked() {
    setLessLoading(true);
    onLessAudios();
  }

  function moreClicked() {
    setMoreLoading(true);
    onMoreAudios();
  }

  return (
    <div id="queue-viewer" className="audio-queue-viewer">
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
          {currentIndex < audios.length - 1 || hasMoreAudios ? (
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
        <ol start={start}>{audios.map(renderPlaylistEntry)}</ol>
        {hasMoreAudios ? (
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

export default AudioQueueViewer;
