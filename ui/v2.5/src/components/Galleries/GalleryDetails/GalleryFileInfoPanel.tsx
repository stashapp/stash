import React from "react";
import * as GQL from "src/core/generated-graphql";

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
        <span className="col-8 text-truncate">{props.gallery.checksum}</span>
      </div>
    );
  }

  function renderPath() {
    const {
      gallery: { path },
    } = props;
    return (
      <div className="row">
        <span className="col-4">Path</span>
        <span className="col-8 text-truncate">
          <a href={`file://${path}`}>{`file://${props.gallery.path}`}</a>{" "}
        </span>
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
