import React from "react";
import * as GQL from "src/core/generated-graphql";
import { TextField, URLField } from "src/utils/field";

interface IGalleryFileInfoPanelProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryFileInfoPanel: React.FC<IGalleryFileInfoPanelProps> = (
  props: IGalleryFileInfoPanelProps
) => {
  return (
    <dl className="container gallery-file-info details-list">
      <TextField
        id="media_info.checksum"
        value={props.gallery.checksum}
        truncate
      />
      <URLField
        id="path"
        url={`file://${props.gallery.path}`}
        value={`file://${props.gallery.path}`}
        truncate
      />
      <URLField
        id="media_info.downloaded_from"
        url={props.gallery.url}
        value={props.gallery.url}
        truncate
      />
    </dl>
  );
};
