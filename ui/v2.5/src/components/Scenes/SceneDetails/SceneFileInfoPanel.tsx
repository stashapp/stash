import React, { useMemo, useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import {
  FormattedMessage,
  FormattedNumber,
  FormattedTime,
  useIntl,
} from "react-intl";
import { useHistory } from "react-router-dom";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { DeleteFilesDialog } from "src/components/Shared/DeleteFilesDialog";
import { ReassignFilesDialog } from "src/components/Shared/ReassignFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateSceneSetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { TextField, URLField, URLsField } from "src/utils/field";
import { StashIDPill } from "src/components/Shared/StashID";
import { PatchComponent } from "../../../patch";

interface IFileInfoPanelProps {
  sceneID: string;
  file: GQL.VideoFileDataFragment;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  onReassign?: () => void;
  loading?: boolean;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const intl = useIntl();
  const history = useHistory();

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

  function onSplit() {
    history.push(
      `/scenes/new?from_scene_id=${props.sceneID}&file_id=${props.file.id}`
    );
  }

  return (
    <div>
      <dl className="container scene-file-info details-list">
        {props.primary && (
          <>
            <dt></dt>
            <dd className="primary-file">
              <FormattedMessage id="primary_file" />
            </dd>
          </>
        )}
        <TextField id="media_info.hash" value={oshash?.value} truncate />
        <TextField id="media_info.checksum" value={checksum?.value} truncate />
        <URLField
          id="media_info.phash"
          abbr="Perceptual hash"
          value={phash?.value}
          url={NavUtils.makeScenesPHashMatchUrl(phash?.value)}
          target="_self"
          truncate
          internal
        />
        <URLField
          id="path"
          url={`file://${props.file.path}`}
          value={`file://${props.file.path}`}
          truncate
        />
        {renderFileSize()}
        <TextField id="file_mod_time">
          <FormattedTime
            dateStyle="medium"
            timeStyle="medium"
            value={props.file.mod_time ?? 0}
          />
        </TextField>
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
      {props.ofMany && props.onSetPrimaryFile && !props.primary && (
        <div>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onSetPrimaryFile}
          >
            <FormattedMessage id="actions.make_primary" />
          </Button>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onReassign}
          >
            <FormattedMessage id="actions.reassign" />
          </Button>
          <Button className="edit-button" onClick={onSplit}>
            <FormattedMessage id="actions.split" />
          </Button>
          <Button
            variant="danger"
            disabled={props.loading}
            onClick={props.onDeleteFile}
          >
            <FormattedMessage id="actions.delete_file" />
          </Button>
        </div>
      )}
    </div>
  );
};

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

const _SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
  const Toast = useToast();

  const [loading, setLoading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<GQL.VideoFileDataFragment>();
  const [reassigningFile, setReassigningFile] =
    useState<GQL.VideoFileDataFragment>();

  function renderStashIDs() {
    if (!props.scene.stash_ids.length) {
      return;
    }

    return (
      <>
        <dt>
          <FormattedMessage id="stash_ids" />
        </dt>
        <dd>
          <dl>
            {props.scene.stash_ids.map((stashID) => {
              return (
                <dd key={stashID.stash_id} className="row no-gutters">
                  <StashIDPill stashID={stashID} linkType="scenes" />
                </dd>
              );
            })}
          </dl>
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
      return (
        <FileInfoPanel sceneID={props.scene.id} file={props.scene.files[0]} />
      );
    }

    async function onSetPrimaryFile(fileID: string) {
      try {
        setLoading(true);
        await mutateSceneSetPrimaryFile(props.scene.id, fileID);
      } catch (e) {
        Toast.error(e);
      } finally {
        setLoading(false);
      }
    }

    return (
      <Accordion defaultActiveKey={props.scene.files[0].id}>
        {deletingFile && (
          <DeleteFilesDialog
            onClose={() => setDeletingFile(undefined)}
            selected={[deletingFile]}
          />
        )}
        {reassigningFile && (
          <ReassignFilesDialog
            onClose={() => setReassigningFile(undefined)}
            selected={reassigningFile}
          />
        )}
        {props.scene.files.map((file, index) => (
          <Card key={file.id} className="scene-file-card">
            <Accordion.Toggle as={Card.Header} eventKey={file.id}>
              <TruncatedText text={TextUtils.fileNameFromPath(file.path)} />
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={file.id}>
              <Card.Body>
                <FileInfoPanel
                  sceneID={props.scene.id}
                  file={file}
                  primary={index === 0}
                  ofMany
                  onSetPrimaryFile={() => onSetPrimaryFile(file.id)}
                  onDeleteFile={() => setDeletingFile(file)}
                  onReassign={() => setReassigningFile(file)}
                  loading={loading}
                />
              </Card.Body>
            </Accordion.Collapse>
          </Card>
        ))}
      </Accordion>
    );
  }, [props.scene, loading, Toast, deletingFile, reassigningFile]);

  return (
    <>
      <dl className="container scene-file-info details-list">
        {props.scene.files.length > 0 && (
          <URLField
            id="media_info.stream"
            url={props.scene.paths.stream}
            value={props.scene.paths.stream}
            truncate
          />
        )}
        {renderFunscript()}
        {renderInteractiveSpeed()}
        <URLsField id="urls" urls={props.scene.urls} truncate />
        {renderStashIDs()}
      </dl>

      {filesPanel}
    </>
  );
};

export const SceneFileInfoPanel = PatchComponent(
  "SceneFileInfoPanel",
  _SceneFileInfoPanel
);
export default SceneFileInfoPanel;
