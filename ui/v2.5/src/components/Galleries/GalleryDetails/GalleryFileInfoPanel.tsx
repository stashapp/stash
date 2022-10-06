import React, { useMemo } from "react";
import { Accordion, Card } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IFileInfoPanelProps {
  folder?: Pick<GQL.Folder, "id" | "path">;
  file?: GQL.GalleryFileDataFragment;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const checksum = props.file?.fingerprints.find((f) => f.type === "md5");
  const path = props.folder ? props.folder.path : props.file?.path ?? "";
  const id = props.folder ? "folder" : "path";

  return (
    <dl className="container gallery-file-info details-list">
      <TextField id="media_info.checksum" value={checksum?.value} truncate />
      <URLField
        id={id}
        url={`file://${path}`}
        value={`file://${path}`}
        truncate
      />
    </dl>
  );
};
interface IGalleryFileInfoPanelProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryFileInfoPanel: React.FC<IGalleryFileInfoPanelProps> = (
  props: IGalleryFileInfoPanelProps
) => {
  const filesPanel = useMemo(() => {
    if (props.gallery.folder) {
      return <FileInfoPanel folder={props.gallery.folder} />;
    }

    if (props.gallery.files.length === 0) {
      return <></>;
    }

    if (props.gallery.files.length === 1) {
      return <FileInfoPanel file={props.gallery.files[0]} />;
    }

    return (
      <Accordion defaultActiveKey="0">
        {props.gallery.files.map((file, index) => (
          <Card key={index} className="gallery-file-card">
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
  }, [props.gallery]);

  return (
    <>
      <dl className="container gallery-file-info details-list">
        <URLField
          id="media_info.downloaded_from"
          url={props.gallery.url}
          value={props.gallery.url}
          truncate
        />
      </dl>

      {filesPanel}
    </>
  );
};
