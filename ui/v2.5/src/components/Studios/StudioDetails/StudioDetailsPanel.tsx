import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { TextField, URLField } from "src/utils/field";

interface IStudioDetailsPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
}) => {
  const intl = useIntl();

  function renderRatingField() {
    if (!studio.rating) {
      return;
    }

    return (
      <>
        <dt>{intl.formatMessage({ id: "rating" })}</dt>
        <dd>
          <RatingStars value={studio.rating} disabled />
        </dd>
      </>
    );
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

  function renderStashIDs() {
    if (!studio.stash_ids?.length) {
      return;
    }

    return (
      <>
        <dt>
          <FormattedMessage id="StashIDs" />
        </dt>
        <dd>
          <ul className="pl-0">
            {studio.stash_ids.map((stashID) => {
              const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
              const link = base ? (
                <a
                  href={`${base}studios/${stashID.stash_id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {stashID.stash_id}
                </a>
              ) : (
                stashID.stash_id
              );
              return (
                <li key={stashID.stash_id} className="row no-gutters">
                  {link}
                </li>
              );
            })}
          </ul>
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
        {renderStashIDs()}
      </dl>
    </div>
  );
};
