import React from "react";
import { FormattedNumber } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TruncatedText } from "src/components/Shared";

interface IImageFileInfoPanelProps {
  image: GQL.ImageDataFragment;
}

export const ImageFileInfoPanel: React.FC<IImageFileInfoPanelProps> = (
  props: IImageFileInfoPanelProps
) => {
  function renderChecksum() {
    return (
      <div className="row">
        <span className="col-4">Checksum</span>
        <TruncatedText className="col-8" text={props.image.checksum} />
      </div>
    );
  }

  function renderPath() {
    const {
      image: { path },
    } = props;
    return (
      <div className="row">
        <span className="col-4">Path</span>
        <a href={`file://${path}`} className="col-8">
          <TruncatedText text={`file://${props.image.path}`} />
        </a>{" "}
      </div>
    );
  }

  function renderFileSize() {
    if (props.image.file.size === undefined) {
      return;
    }

    const { size, unit } = TextUtils.fileSize(props.image.file.size ?? 0);

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

  function renderDimensions() {
    if (props.image.file.height && props.image.file.width) {
      return (
        <div className="row">
          <span className="col-4">Dimensions</span>
          <span className="col-8 text-truncate">
            {props.image.file.width} x {props.image.file.height}
          </span>
        </div>
      );
    }
  }

  return (
    <div className="container image-file-info">
      {renderChecksum()}
      {renderPath()}
      {renderFileSize()}
      {renderDimensions()}
    </div>
  );
};
