import React from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Button, Badge, Card } from "react-bootstrap";
import TextUtils from "src/utils/text";
import { markerTitle } from "src/core/markers";

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

  const primaryTagNames: Record<string, string> = {};
  const markersByTag: Record<string, GQL.SceneMarkerDataFragment[]> = {};
  sceneMarkers.forEach((m) => {
    if (primaryTagNames[m.primary_tag.id]) {
      markersByTag[m.primary_tag.id].push(m);
    } else {
      primaryTagNames[m.primary_tag.id] = m.primary_tag.name;
      markersByTag[m.primary_tag.id] = [m];
    }
  });

  const primaryCards = Object.keys(markersByTag).map((id) => {
    const markers = markersByTag[id].map((marker) => {
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
              {markerTitle(marker)}
            </Button>
            <Button
              variant="link"
              className="ml-auto"
              onClick={() => onEdit(marker)}
            >
              <FormattedMessage id="actions.edit" />
            </Button>
          </div>
          <div>
            {TextUtils.formatTimestampRange(
              marker.seconds,
              marker.end_seconds ?? undefined
            )}
          </div>
          <div className="card-section centered">{tags}</div>
        </div>
      );
    });

    return (
      <Card className="primary-card col-12 col-sm-6 col-xl-6" key={id}>
        <h3>{primaryTagNames[id]}</h3>
        <Card.Body className="primary-card-body">{markers}</Card.Body>
      </Card>
    );
  });

  return <div className="primary-tag row">{primaryCards}</div>;
};
