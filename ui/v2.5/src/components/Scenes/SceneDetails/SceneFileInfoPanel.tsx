import React from "react";
import { FormattedMessage, FormattedNumber } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TruncatedText } from "src/components/Shared";

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
  function renderOSHash() {
    if (props.scene.oshash) {
      return (
        <div className="row">
          <span className="col-4">
            <FormattedMessage id="media_info.hash" />
          </span>
          <TruncatedText className="col-8" text={props.scene.oshash} />
        </div>
      );
    }
  }

  function renderChecksum() {
    if (props.scene.checksum) {
      return (
        <div className="row">
          <span className="col-4">
            <FormattedMessage id="media_info.checksum" />
          </span>
          <TruncatedText className="col-8" text={props.scene.checksum} />
        </div>
      );
    }
  }

  function renderPath() {
    const {
      scene: { path },
    } = props;
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="path" />
        </span>
        <a href={`file://${path}`} className="col-8" target="_blank">
          <TruncatedText text={`file://${props.scene.path}`} />
        </a>{" "}
      </div>
    );
  }

  function renderStream() {
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="media_info.stream" />
        </span>
        <a href={props.scene.paths.stream ?? ""} className="col-8" target="_blank">
          <TruncatedText text={props.scene.paths.stream} />
        </a>{" "}
      </div>
    );
  }

  function renderFileSize() {
    if (props.scene.file.size === undefined) {
      return;
    }

    const { size, unit } = TextUtils.fileSize(
      Number.parseInt(props.scene.file.size ?? "0", 10)
    );

    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="filesize" />
        </span>
        <span className="col-8 text-truncate">
          <FormattedNumber
            value={size}
            // eslint-disable-next-line react/style-prop-object
            style="unit"
            unit={unit}
            unitDisplay="narrow"
            maximumFractionDigits={2}
          />
        </span>
      </div>
    );
  }

  function renderDuration() {
    if (props.scene.file.duration === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="duration" />
        </span>
        <TruncatedText
          className="col-8"
          text={TextUtils.secondsToTimestamp(props.scene.file.duration ?? 0)}
        />
      </div>
    );
  }

  function renderDimensions() {
    if (props.scene.file.duration === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="dimensions" />
        </span>
        <TruncatedText
          className="col-8"
          text={`${props.scene.file.width} x ${props.scene.file.height}`}
        />
      </div>
    );
  }

  function renderFrameRate() {
    if (props.scene.file.framerate === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="framerate" />
        </span>
        <span className="col-8 text-truncate">
          <FormattedNumber value={props.scene.file.framerate ?? 0} /> frames per
          second
        </span>
      </div>
    );
  }

  function renderbitrate() {
    // TODO: An upcoming react-intl version will support compound units, megabits-per-second
    if (props.scene.file.bitrate === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="bitrate" />
        </span>
        <span className="col-8 text-truncate">
          <FormattedNumber
            value={(props.scene.file.bitrate ?? 0) / 1000000}
            maximumFractionDigits={2}
          />
          &nbsp;megabits per second
        </span>
      </div>
    );
  }

  function renderVideoCodec() {
    if (props.scene.file.video_codec === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="media_info.video_codec" />
        </span>
        <TruncatedText className="col-8" text={props.scene.file.video_codec} />
      </div>
    );
  }

  function renderAudioCodec() {
    if (props.scene.file.audio_codec === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="media_info.audio_codec" />
        </span>
        <TruncatedText className="col-8" text={props.scene.file.audio_codec} />
      </div>
    );
  }

  function renderUrl() {
    if (!props.scene.url || props.scene.url === "") {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">
          <FormattedMessage id="media_info.downloaded_from" />
        </span>
        <a href={TextUtils.sanitiseURL(props.scene.url)} className="col-8" target="_blank">
          <TruncatedText text={props.scene.url} />
        </a>
      </div>
    );
  }

  function renderStashIDs() {
    if (!props.scene.stash_ids.length) {
      return;
    }

    return (
      <div className="row">
        <span className="col-4">StashIDs</span>
        <ul className="col-8">
          {props.scene.stash_ids.map((stashID) => {
            const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
            const link = base ? (
              <a
                href={`${base}scenes/${stashID.stash_id}`}
                target="_blank"
                rel="noopener noreferrer"
              >
                {stashID.stash_id}
              </a>
            ) : (
              stashID.stash_id
            );
            return (
              <li key={stashID.stash_id} className="row no-gutters">
                {link}
              </li>
            );
          })}
        </ul>
      </div>
    );
  }

  function renderPhash() {
    if (props.scene.phash) {
      return (
        <div className="row">
          <abbr className="col-4" title="Perceptual hash">
            <FormattedMessage id="media_info.phash" />
          </abbr>
          <TruncatedText className="col-8" text={props.scene.phash} />
        </div>
      );
    }
  }

  function renderFunscript() {
    if (props.scene.interactive) {
      return (
        <div className="row">
          <span className="col-4">Funscript</span>
          <a href={props.scene.paths.funscript ?? ""} className="col-8">
            <TruncatedText text={props.scene.paths.funscript} />
          </a>{" "}
        </div>
      );
    }
  }

  return (
    <div className="container scene-file-info">
      {renderOSHash()}
      {renderChecksum()}
      {renderPhash()}
      {renderPath()}
      {renderStream()}
      {renderFunscript()}
      {renderFileSize()}
      {renderDuration()}
      {renderDimensions()}
      {renderFrameRate()}
      {renderbitrate()}
      {renderVideoCodec()}
      {renderAudioCodec()}
      {renderUrl()}
      {renderStashIDs()}
    </div>
  );
};
