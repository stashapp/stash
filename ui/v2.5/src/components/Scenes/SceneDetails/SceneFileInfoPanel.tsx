import React, { useMemo } from "react";
import { Accordion, Card } from "react-bootstrap";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import { TruncatedText } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils, getStashboxBase } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IFileInfoPanelProps {
  file: GQL.VideoFileDataFragment;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const intl = useIntl();

  function renderFileSize() {
    const { size, unit } = TextUtils.fileSize(props.file.size);

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

  // TODO - generalise fingerprints
  const oshash = props.file.fingerprints.find((f) => f.type === "oshash");
  const phash = props.file.fingerprints.find((f) => f.type === "phash");
  const checksum = props.file.fingerprints.find((f) => f.type === "md5");

  return (
    <dl className="container scene-file-info details-list">
      <TextField id="media_info.hash" value={oshash?.value} truncate />
      <TextField id="media_info.checksum" value={checksum?.value} truncate />
      <URLField
        id="media_info.phash"
        abbr="Perceptual hash"
        value={phash?.value}
        url={NavUtils.makeScenesPHashMatchUrl(phash?.value)}
        target="_self"
        truncate
        trusted
      />
      <URLField
        id="path"
        url={`file://${props.file.path}`}
        value={`file://${props.file.path}`}
        truncate
      />
      {renderFileSize()}
      <TextField
        id="duration"
        value={TextUtils.secondsToTimestamp(props.file.duration ?? 0)}
        truncate
      />
      <TextField
        id="dimensions"
        value={`${props.file.width} x ${props.file.height}`}
        truncate
      />
      <TextField id="framerate">
        <FormattedMessage
          id="frames_per_second"
          values={{ value: intl.formatNumber(props.file.frame_rate ?? 0) }}
        />
      </TextField>
      <TextField id="bitrate">
        <FormattedMessage
          id="megabits_per_second"
          values={{
            value: intl.formatNumber((props.file.bit_rate ?? 0) / 1000000, {
              maximumFractionDigits: 2,
            }),
          }}
        />
      </TextField>
      <TextField
        id="media_info.video_codec"
        value={props.file.video_codec ?? ""}
        truncate
      />
      <TextField
        id="media_info.audio_codec"
        value={props.file.audio_codec ?? ""}
        truncate
      />
    </dl>
  );
};

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
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

  const filesPanel = useMemo(() => {
    if (props.scene.files.length === 0) {
      return;
    }

    if (props.scene.files.length === 1) {
      return <FileInfoPanel file={props.scene.files[0]} />;
    }

    return (
      <Accordion defaultActiveKey="0">
        {props.scene.files.map((file, index) => (
          <Card key={index} className="scene-file-card">
            <Accordion.Toggle as={Card.Header} eventKey={index.toString()}>
              <TruncatedText text={TextUtils.fileNameFromPath(file.path)} />
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={index.toString()}>
              <Card.Body>
                <FileInfoPanel file={file} />
              </Card.Body>
            </Accordion.Collapse>
          </Card>
        ))}
      </Accordion>
    );
  }, [props.scene]);

  return (
    <>
      <dl className="container scene-file-info details-list">
        <URLField
          id="media_info.stream"
          url={props.scene.paths.stream}
          value={props.scene.paths.stream}
          truncate
        />
        {renderFunscript()}
        {renderInteractiveSpeed()}
        <URLField
          id="media_info.downloaded_from"
          url={props.scene.url}
          value={props.scene.url}
          truncate
        />
        {renderStashIDs()}
      </dl>

      {filesPanel}
    </>
  );
};

export default SceneFileInfoPanel;
