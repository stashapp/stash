import React from "react";
import { Accordion, Card } from "react-bootstrap";
import { FormattedNumber } from "react-intl";
import { TruncatedText } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IFileInfoPanelProps {
  file: GQL.ImageFileDataFragment;
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
    <dl className="container image-file-info details-list">
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
  );
};
interface IImageFileInfoPanelProps {
  image: GQL.ImageDataFragment;
}

export const ImageFileInfoPanel: React.FC<IImageFileInfoPanelProps> = (
  props: IImageFileInfoPanelProps
) => {
  if (props.image.files.length === 0) {
    return <></>;
  }

  if (props.image.files.length === 1) {
    return <FileInfoPanel file={props.image.files[0]} />;
  }

  return (
    <Accordion defaultActiveKey="0">
      {props.image.files.map((file, index) => (
        <Card key={index} className="image-file-card">
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
};
