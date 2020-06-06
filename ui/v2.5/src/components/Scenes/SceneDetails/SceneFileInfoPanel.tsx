import React from "react";
import { FormattedNumber } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
  function renderChecksum() {
    return (
      <div className="row">
        <span className="col-4">Checksum</span>
        <span className="col-8 text-truncate">{props.scene.checksum}</span>
      </div>
    );
  }

  function renderPath() {
    const {
      scene: { path },
    } = props;
    return (
      <div className="row">
        <span className="col-4">Path</span>
        <span className="col-8 text-truncate">
          <a href={`file://${path}`}>{`file://${props.scene.path}`}</a>{" "}
        </span>
      </div>
    );
  }

  function renderStream() {
    return (
      <div className="row">
        <span className="col-4">Stream</span>
        <span className="col-8 text-truncate">
          <a href={props.scene.paths.stream ?? ""}>
            {props.scene.paths.stream}
          </a>{" "}
        </span>
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
        <span className="col-4">File Size</span>
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
        <span className="col-4">Duration</span>
        <span className="col-8 text-truncate">
          {TextUtils.secondsToTimestamp(props.scene.file.duration ?? 0)}
        </span>
      </div>
    );
  }

  function renderDimensions() {
    if (props.scene.file.duration === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">Dimensions</span>
        <span className="col-8 text-truncate">
          {props.scene.file.width} x {props.scene.file.height}
        </span>
      </div>
    );
  }

  function renderFrameRate() {
    if (props.scene.file.framerate === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">Frame Rate</span>
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
        <span className="col-4">Bit Rate</span>
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
        <span className="col-4">Video Codec</span>
        <span className="col-8 text-truncate">
          {props.scene.file.video_codec}
        </span>
      </div>
    );
  }

  function renderAudioCodec() {
    if (props.scene.file.audio_codec === undefined) {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">Audio Codec</span>
        <span className="col-8 text-truncate">
          {props.scene.file.audio_codec}
        </span>
      </div>
    );
  }

  function renderUrl() {
    if (!props.scene.url || props.scene.url === "") {
      return;
    }
    return (
      <div className="row">
        <span className="col-4">Downloaded From</span>
        <span className="col-8 text-truncate">
          <a href={TextUtils.sanitiseURL(props.scene.url)}>{props.scene.url}</a>
        </span>
      </div>
    );
  }

  return (
    <div className="container scene-file-info">
      {renderChecksum()}
      {renderPath()}
      {renderStream()}
      {renderFileSize()}
      {renderDuration()}
      {renderDimensions()}
      {renderFrameRate()}
      {renderbitrate()}
      {renderVideoCodec()}
      {renderAudioCodec()}
      {renderUrl()}
    </div>
  );
};
