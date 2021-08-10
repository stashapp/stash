import React from "react";
import { FormattedNumber } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IImageFileInfoPanelProps {
  image: GQL.ImageDataFragment;
}

export const ImageFileInfoPanel: React.FC<IImageFileInfoPanelProps> = (
  props: IImageFileInfoPanelProps
) => {
  function renderFileSize() {
    if (props.image.file.size === undefined) {
      return;
    }

    const { size, unit } = TextUtils.fileSize(props.image.file.size ?? 0);

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

  return (
    <dl className="container image-file-info details-list">
      <TextField
        id="media_info.checksum"
        value={props.image.checksum}
        truncate
      />
      <URLField
        id="path"
        url={`file://${props.image.path}`}
        value={`file://${props.image.path}`}
        truncate
      />
      {renderFileSize()}
      <TextField
        id="dimensions"
        value={`${props.image.file.width} x ${props.image.file.height}`}
        truncate
      />
    </dl>
  );
};
