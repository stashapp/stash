import {
  HTMLTable,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { TextUtils } from "../../../utils/text";

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneFileInfoPanel: FunctionComponent<ISceneFileInfoPanelProps> = (props: ISceneFileInfoPanelProps) => {
  function renderChecksum() {
    return (
      <tr>
        <td>Checksum</td>
        <td>{props.scene.checksum}</td>
      </tr>
    );
  }

  function renderPath() {
    return (
      <tr>
        <td>Path</td>
        <td><a href={"file://"+props.scene.path}>{"file://"+props.scene.path}</a> </td>
      </tr>
    );
  }
  
  function renderStream() {
    return (
      <tr>
        <td>Stream</td>
        <td><a href={props.scene.paths.stream}>{props.scene.paths.stream}</a> </td>
      </tr>
    );
  }

  function renderFileSize() {
    if (props.scene.file.size === undefined) { return; }
    return (
      <tr>
        <td>File Size</td>
        <td>{TextUtils.fileSize(parseInt(props.scene.file.size, 10))}</td>
      </tr>
    );
  }

  function renderDuration() {
    if (props.scene.file.duration === undefined) { return; }
    return (
      <tr>
        <td>Duration</td>
        <td>{TextUtils.secondsToTimestamp(props.scene.file.duration)}</td>
      </tr>
    );
  }

  function renderDimensions() {
    if (props.scene.file.duration === undefined) { return; }
    return (
      <tr>
        <td>Dimensions</td>
        <td>{props.scene.file.width} x {props.scene.file.height}</td>
      </tr>
    );
  }

  function renderFrameRate() {
    if (props.scene.file.framerate === undefined) { return; }
    return (
      <tr>
        <td>Frame Rate</td>
        <td>{props.scene.file.framerate} frames per second</td>
      </tr>
    );
  }

  function renderBitRate() {
    if (props.scene.file.bitrate === undefined) { return; }
    return (
      <tr>
        <td>Bit Rate</td>
        <td>{TextUtils.bitRate(props.scene.file.bitrate)}</td>
      </tr>
    );
  }

  function renderVideoCodec() {
    if (props.scene.file.video_codec === undefined) { return; }
    return (
      <tr>
        <td>Video Codec</td>
        <td>{props.scene.file.video_codec}</td>
      </tr>
    );
  }

  function renderAudioCodec() {
    if (props.scene.file.audio_codec === undefined) { return; }
    return (
      <tr>
        <td>Audio Codec</td>
        <td>{props.scene.file.audio_codec}</td>
      </tr>
    );
  }

  function renderUrl() {
    if (!props.scene.url || props.scene.url === "") { return; }
    return (
      <tr>
        <td>Downloaded From</td>
        <td>{props.scene.url}</td>
      </tr>
    );
  }

  return (
    <>
      <HTMLTable>
        <tbody>
          {renderChecksum()}
          {renderPath()}
          {renderStream()}
          {renderFileSize()}
          {renderDuration()}
          {renderDimensions()}
          {renderFrameRate()}
          {renderBitRate()}
          {renderVideoCodec()}
          {renderAudioCodec()}
          {renderUrl()}
        </tbody>
      </HTMLTable>
    </>
  );
};
