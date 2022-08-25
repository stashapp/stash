import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { RatingSystem } from "src/components/Scenes/SceneDetails/RatingSystem";
import { TextField, URLField } from "src/utils/field";

interface IStudioDetailsPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
}) => {
  function renderRatingField() {
    if (!studio.rating) {
      return;
    }

    return <RatingSystem value={studio.rating} disabled />;
  }

  function renderTagsList() {
    if (!studio.aliases?.length) {
      return;
    }

    return (
      <>
        <dt>
          <FormattedMessage id="aliases" />
        </dt>
        <dd>
          {studio.aliases.map((a) => (
            <Badge className="tag-item" variant="secondary" key={a}>
              {a}
            </Badge>
          ))}
        </dd>
      </>
    );
  }

  return (
    <div className="studio-details">
      <div>
        <h2>{studio.name}</h2>
      </div>

      <dl className="details-list">
        <URLField
          id="url"
          value={studio.url}
          url={TextUtils.sanitiseURL(studio.url ?? "")}
        />

        <TextField id="details" value={studio.details} />

        <URLField
          id="parent_studios"
          value={studio.parent_studio?.name}
          url={`/studios/${studio.parent_studio?.id}`}
          trusted
          target="_self"
        />

        {renderRatingField()}
        {renderTagsList()}
      </dl>
    </div>
  );
};
