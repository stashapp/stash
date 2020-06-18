import React from "react";
import * as GQL from "src/core/generated-graphql";
import { Button, Badge, Card } from "react-bootstrap";
import { TextUtils } from "src/utils";

interface IPrimaryTags {
  sceneMarkers: GQL.SceneMarkerDataFragment[];
  onClickMarker: (marker: GQL.SceneMarkerDataFragment) => void;
  onEdit: (marker: GQL.SceneMarkerDataFragment) => void;
}

export const PrimaryTags: React.FC<IPrimaryTags> = ({
  sceneMarkers,
  onClickMarker,
  onEdit,
}) => {
  if (!sceneMarkers?.length) return <div />;

  const primaries: Record<string, GQL.Tag> = {};
  const primaryTags: Record<string, GQL.SceneMarkerDataFragment[]> = {};
  sceneMarkers.forEach((m) => {
    if (primaryTags[m.primary_tag.id]) primaryTags[m.primary_tag.id].push(m);
    else {
      primaryTags[m.primary_tag.id] = [m];
      primaries[m.primary_tag.id] = m.primary_tag;
    }
  });

  const primaryCards = Object.keys(primaryTags).map((id) => {
    const markers = primaryTags[id].map((marker) => {
      const tags = marker.tags.map((tag) => (
        <Badge key={tag.id} variant="secondary" className="tag-item">
          {tag.name}
        </Badge>
      ));

      return (
        <div key={marker.id}>
          <hr />
          <div className="row">
            <Button variant="link" onClick={() => onClickMarker(marker)}>
              {marker.title}
            </Button>
            <Button
              variant="link"
              className="ml-auto"
              onClick={() => onEdit(marker)}
            >
              Edit
            </Button>
          </div>
          <div>{TextUtils.secondsToTimestamp(marker.seconds)}</div>
          <div className="card-section centered">{tags}</div>
        </div>
      );
    });

    return (
      <Card className="primary-card col-12 col-sm-3 col-xl-6" key={id}>
        <h3>{primaries[id].name}</h3>
        <Card.Body className="primary-card-body">{markers}</Card.Body>
      </Card>
    );
  });

  return <div className="primary-tag row">{primaryCards}</div>;
};
