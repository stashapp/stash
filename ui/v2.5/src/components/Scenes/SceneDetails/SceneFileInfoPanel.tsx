import React from "react";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils, getStashboxBase } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
  const intl = useIntl();

  function renderFileSize() {
    if (props.scene.file.size === undefined) {
      return;
    }

    const { size, unit } = TextUtils.fileSize(
      Number.parseInt(props.scene.file.size ?? "0", 10)
    );

    return (
      <TextField id="filesize">
        <span className="text-truncate">
          <FormattedNumber
            value={size}
            // eslint-disable-next-line react/style-prop-object
            style="unit"
            unit={unit}
            unitDisplay="narrow"
            maximumFractionDigits={2}
          />
        </span>
      </TextField>
    );
  }

  function renderStashIDs() {
    if (!props.scene.stash_ids.length) {
      return;
    }

    return (
      <>
        <dt>StashIDs</dt>
        <dd>
          <ul>
            {props.scene.stash_ids.map((stashID) => {
              const base = getStashboxBase(stashID.endpoint);
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
        </dd>
      </>
    );
  }

  function renderFunscript() {
    if (props.scene.interactive) {
      return (
        <URLField
          name="Funscript"
          url={props.scene.paths.funscript}
          value={props.scene.paths.funscript}
          truncate
        />
      );
    }
  }

  function renderInteractiveSpeed() {
    if (props.scene.interactive_speed) {
      return (
        <TextField id="media_info.interactive_speed">
          <FormattedNumber value={props.scene.interactive_speed} />
        </TextField>
      );
    }
  }

  return (
    <dl className="container scene-file-info details-list">
      <TextField id="media_info.hash" value={props.scene.oshash} truncate />
      <TextField
        id="media_info.checksum"
        value={props.scene.checksum}
        truncate
      />
      <URLField
        id="media_info.phash"
        abbr="Perceptual hash"
        value={props.scene.phash}
        url={NavUtils.makeScenesPHashMatchUrl(props.scene.phash)}
        target="_self"
        truncate
        trusted
      />
      <URLField
        id="path"
        url={`file://${props.scene.path}`}
        value={`file://${props.scene.path}`}
        truncate
      />
      <URLField
        id="media_info.stream"
        url={props.scene.paths.stream}
        value={props.scene.paths.stream}
        truncate
      />
      {renderFunscript()}
      {renderInteractiveSpeed()}
      {renderFileSize()}
      <TextField
        id="duration"
        value={TextUtils.secondsToTimestamp(props.scene.file.duration ?? 0)}
        truncate
      />
      <TextField
        id="dimensions"
        value={`${props.scene.file.width} x ${props.scene.file.height}`}
        truncate
      />
      <TextField id="framerate">
        <FormattedMessage
          id="frames_per_second"
          values={{ value: intl.formatNumber(props.scene.file.framerate ?? 0) }}
        />
      </TextField>
      <TextField id="bitrate">
        <FormattedMessage
          id="megabits_per_second"
          values={{
            value: intl.formatNumber(
              (props.scene.file.bitrate ?? 0) / 1000000,
              { maximumFractionDigits: 2 }
            ),
          }}
        />
      </TextField>
      <TextField
        id="media_info.video_codec"
        value={props.scene.file.video_codec}
        truncate
      />
      <TextField
        id="media_info.audio_codec"
        value={props.scene.file.audio_codec}
        truncate
      />
      <URLField
        id="media_info.downloaded_from"
        url={props.scene.url}
        value={props.scene.url}
        truncate
      />
      {renderStashIDs()}
    </dl>
  );
};
