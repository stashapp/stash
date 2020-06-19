import { Card } from "react-bootstrap";
import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FormattedPlural } from "react-intl";

interface IProps {
  gallery: GQL.GalleryDataFragment;
  zoomIndex: number;
}

export const GalleryCard: React.FC<IProps> = ({ gallery, zoomIndex }) => {
  return (
    <Card className={`gallery-card zoom-${zoomIndex}`}> 
      <Link to={`/galleries/${gallery.id}`} className="gallery-card-header">
        {gallery.files.length > 0 ? 
        <img
          className="gallery-card-image"
          alt={gallery.path}
          src={`${gallery.files[0].path}?thumb=true`}
        />
        : undefined}
      </Link>
      <div className="card-section">
        <h5 className="card-section-title">{gallery.path}</h5>
        <span>
          {gallery.files.length}&nbsp;
          <FormattedPlural
            value={gallery.files.length ?? 0}
            one="image"
            other="images"
          />
          .
        </span>
      </div>
    </Card>
  );
};
