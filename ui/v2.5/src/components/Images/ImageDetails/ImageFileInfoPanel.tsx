import React, { useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import { FormattedMessage, FormattedTime } from "react-intl";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { DeleteFilesDialog } from "src/components/Shared/DeleteFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateImageSetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import TextUtils from "src/utils/text";
import { TextField, URLField, URLsField } from "src/utils/field";
import { FileSize } from "src/components/Shared/FileSize";

interface IFileInfoPanelProps {
  file: GQL.ImageFileDataFragment | GQL.VideoFileDataFragment;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  loading?: boolean;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const checksum = props.file.fingerprints.find((f) => f.type === "md5");

  return (
    <div>
      <dl className="container image-file-info details-list">
        {props.primary && (
          <>
            <dt></dt>
            <dd className="primary-file">
              <FormattedMessage id="primary_file" />
            </dd>
          </>
        )}
        <TextField id="media_info.checksum" value={checksum?.value} truncate />
        <URLField
          id="path"
          url={`file://${props.file.path}`}
          value={`file://${props.file.path}`}
          truncate
        />
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
          id="dimensions"
          value={`${props.file.width} x ${props.file.height}`}
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
interface IImageFileInfoPanelProps {
  image: GQL.ImageDataFragment;
}

export const ImageFileInfoPanel: React.FC<IImageFileInfoPanelProps> = (
  props: IImageFileInfoPanelProps
) => {
  const Toast = useToast();

  const [loading, setLoading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<
    GQL.ImageFileDataFragment | GQL.VideoFileDataFragment | undefined
  >();

  if (props.image.visual_files.length === 0) {
    return <></>;
  }

  if (props.image.visual_files.length === 1) {
    return (
      <>
        <dl className="container image-file-info details-list">
          <URLsField id="urls" urls={props.image.urls} truncate />
        </dl>

        <FileInfoPanel file={props.image.visual_files[0]} />
      </>
    );
  }

  async function onSetPrimaryFile(fileID: string) {
    try {
      setLoading(true);
      await mutateImageSetPrimaryFile(props.image.id, fileID);
    } catch (e) {
      Toast.error(e);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Accordion defaultActiveKey={props.image.visual_files[0].id}>
      {deletingFile && (
        <DeleteFilesDialog
          onClose={() => setDeletingFile(undefined)}
          selected={[deletingFile]}
        />
      )}
      {props.image.visual_files.map((file, index) => (
        <Card key={file.id} className="image-file-card">
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
                loading={loading}
              />
            </Card.Body>
          </Accordion.Collapse>
        </Card>
      ))}
    </Accordion>
  );
};
