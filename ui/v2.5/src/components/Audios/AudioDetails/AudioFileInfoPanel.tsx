import React, { useMemo, useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import { FormattedMessage, FormattedTime, useIntl } from "react-intl";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { DeleteFilesDialog } from "src/components/Shared/DeleteFilesDialog";
import { RevealInFilesystemButton } from "src/components/Shared/RevealInFilesystemButton";
import { ReassignFilesDialog } from "src/components/Shared/ReassignFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateAudioSetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import TextUtils from "src/utils/text";
import { TextField, URLField, URLsField } from "src/utils/field";
import { PatchComponent } from "../../../patch";
import { FileSize } from "src/components/Shared/FileSize";

interface IFileInfoPanelProps {
  file: GQL.AudioFileDataFragment;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  onReassign?: () => void;
  loading?: boolean;
}

// NOTE - audio files have no resolution, framerate, video codec or phash
const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const intl = useIntl();

  const oshash = props.file.fingerprints.find((f) => f.type === "oshash");
  const checksum = props.file.fingerprints.find((f) => f.type === "md5");

  return (
    <div>
      <dl className="container audio-file-info details-list">
        {props.primary && (
          <>
            <dt></dt>
            <dd className="primary-file">
              <FormattedMessage id="primary_file" />
            </dd>
          </>
        )}
        <TextField
          id="media_info.oshash"
          abbr={intl.formatMessage({ id: "media_info.oshash_meaning" })}
          value={oshash?.value}
          truncate
        />
        <TextField id="media_info.md5" value={checksum?.value} truncate />
        <TextField id="path">
          <span className="d-flex align-items-center">
            <TruncatedText text={props.file.path} />
            <RevealInFilesystemButton fileId={props.file.id} />
          </span>
        </TextField>
        <TextField id="filesize">
          <span className="text-truncate">
            <FileSize size={props.file.size} />
          </span>
        </TextField>
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
        <TextField id="bitrate">
          <FormattedMessage
            id="kilobits_per_second"
            values={{
              value: intl.formatNumber((props.file.bit_rate ?? 0) / 1000, {
                maximumFractionDigits: 0,
              }),
            }}
          />
        </TextField>
        <TextField id="sample_rate">
          <FormattedMessage
            id="hertz"
            values={{
              value: intl.formatNumber(props.file.sample_rate ?? 0),
            }}
          />
        </TextField>
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

interface IAudioFileInfoPanelProps {
  audio: GQL.AudioDataFragment;
}

const _AudioFileInfoPanel: React.FC<IAudioFileInfoPanelProps> = (
  props: IAudioFileInfoPanelProps
) => {
  const Toast = useToast();

  const [loading, setLoading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<GQL.AudioFileDataFragment>();
  const [reassigningFile, setReassigningFile] =
    useState<GQL.AudioFileDataFragment>();

  const filesPanel = useMemo(() => {
    if (props.audio.files.length === 0) {
      return;
    }

    if (props.audio.files.length === 1) {
      return <FileInfoPanel file={props.audio.files[0]} />;
    }

    async function onSetPrimaryFile(fileID: string) {
      try {
        setLoading(true);
        await mutateAudioSetPrimaryFile(props.audio.id, fileID);
      } catch (e) {
        Toast.error(e);
      } finally {
        setLoading(false);
      }
    }

    return (
      <Accordion defaultActiveKey={props.audio.files[0].id}>
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
        {props.audio.files.map((file, index) => (
          <Card key={file.id} className="audio-file-card">
            <Accordion.Toggle as={Card.Header} eventKey={file.id}>
              <TruncatedText text={TextUtils.fileNameFromPath(file.path)} />
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={file.id}>
              <Card.Body>
                <FileInfoPanel
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
  }, [props.audio, loading, Toast, deletingFile, reassigningFile]);

  return (
    <>
      <dl className="container audio-file-info details-list">
        {props.audio.files.length > 0 && (
          <URLField
            id="media_info.stream"
            url={props.audio.paths.stream}
            value={props.audio.paths.stream}
            truncate
          />
        )}
        <URLsField id="urls" urls={props.audio.urls} truncate />
      </dl>

      {filesPanel}
    </>
  );
};

export const AudioFileInfoPanel = PatchComponent(
  "AudioFileInfoPanel",
  _AudioFileInfoPanel
);
export default AudioFileInfoPanel;
