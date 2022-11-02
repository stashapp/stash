import React, { useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import { FormattedMessage, FormattedNumber } from "react-intl";
import { TruncatedText } from "src/components/Shared";
import DeleteFilesDialog from "src/components/Shared/DeleteFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateImageSetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IFileInfoPanelProps {
  file: GQL.ImageFileDataFragment;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  loading?: boolean;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  function renderFileSize() {
    if (props.file.size === undefined) {
      return;
    }

    const { size, unit } = TextUtils.fileSize(props.file.size ?? 0);

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
        {renderFileSize()}
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
    GQL.ImageFileDataFragment | undefined
  >();

  if (props.image.files.length === 0) {
    return <></>;
  }

  if (props.image.files.length === 1) {
    return <FileInfoPanel file={props.image.files[0]} />;
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
    <Accordion defaultActiveKey={props.image.files[0].id}>
      {deletingFile && (
        <DeleteFilesDialog
          onClose={() => setDeletingFile(undefined)}
          selected={[deletingFile]}
        />
      )}
      {props.image.files.map((file, index) => (
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
