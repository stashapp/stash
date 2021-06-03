import React from "react";
import * as GQL from "src/core/generated-graphql";
import { TruncatedText } from "src/components/Shared";
import { FormattedMessage } from "react-intl";

interface IGalleryFileInfoPanelProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryFileInfoPanel: React.FC<IGalleryFileInfoPanelProps> = (
  props: IGalleryFileInfoPanelProps
) => {
  function renderChecksum() {
    return (
      <div className="row">
        <span className="col-4">Checksum</span>
        <TruncatedText className="col-8" text={props.gallery.checksum} />
      </div>
    );
  }

  function renderPath() {
    const filePath = `file://${props.gallery.path}`;

    return (
      <div className="row">
        <span className="col-4"><FormattedMessage id="path"/></span>
        <a href={filePath} className="col-8">
          <TruncatedText text={filePath} />
        </a>
      </div>
    );
  }

  return (
    <div className="container gallery-file-info">
      {renderChecksum()}
      {renderPath()}
    </div>
  );
};
