import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { TagLink } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { TextUtils, getStashboxBase } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IPerformerDetails {
  performer: GQL.PerformerDataFragment;
}

export const PerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    if (!performer.tags.length) {
      return;
    }

    return (
      <>
        <dt>
          <FormattedMessage id="tags" />
        </dt>
        <dd>
          <ul className="pl-0">
            {(performer.tags ?? []).map((tag) => (
              <TagLink key={tag.id} tagType="performer" tag={tag} />
            ))}
          </ul>
        </dd>
      </>
    );
  }

  function renderStashIDs() {
    if (!performer.stash_ids.length) {
      return;
    }

    return (
      <>
        <dt>StashIDs</dt>
        <dd>
          <ul className="pl-0">
            {performer.stash_ids.map((stashID) => {
              const base = getStashboxBase(stashID.endpoint);
              const link = base ? (
                <a
                  href={`${base}performers/${stashID.stash_id}`}
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

  const formatHeight = (height?: string | null) => {
    if (!height) {
      return "";
    }
    return intl.formatNumber(Number.parseInt(height, 10), {
      style: "unit",
      unit: "centimeter",
      unitDisplay: "narrow",
    });
  };

  const formatWeight = (weight?: number | null) => {
    if (!weight) {
      return "";
    }
    return intl.formatNumber(weight, {
      style: "unit",
      unit: "kilogram",
      unitDisplay: "narrow",
    });
  };

  return (
    <dl className="details-list">
      <TextField
        id="gender"
        value={
          performer.gender
            ? intl.formatMessage({ id: "gender_types." + performer.gender })
            : undefined
        }
      />
      <TextField
        id="birthdate"
        value={TextUtils.formatDate(intl, performer.birthdate ?? undefined)}
      />
      <TextField
        id="death_date"
        value={TextUtils.formatDate(intl, performer.death_date ?? undefined)}
      />
      <TextField id="ethnicity" value={performer.ethnicity} />
      <TextField id="hair_color" value={performer.hair_color} />
      <TextField id="eye_color" value={performer.eye_color} />
      <TextField id="country" value={performer.country} />
      <TextField id="height" value={formatHeight(performer.height)} />
      <TextField id="weight" value={formatWeight(performer.weight)} />
      <TextField id="measurements" value={performer.measurements} />
      <TextField id="fake_tits" value={performer.fake_tits} />
      <TextField id="career_length" value={performer.career_length} />
      <TextField id="tattoos" value={performer.tattoos} />
      <TextField id="piercings" value={performer.piercings} />
      <TextField id="details" value={performer.details} />
      <URLField
        id="url"
        value={performer.url}
        url={TextUtils.sanitiseURL(performer.url ?? "")}
      />
      <URLField
        id="twitter"
        value={performer.twitter}
        url={TextUtils.sanitiseURL(
          performer.twitter ?? "",
          TextUtils.twitterURL
        )}
      />
      <URLField
        id="instagram"
        value={performer.instagram}
        url={TextUtils.sanitiseURL(
          performer.instagram ?? "",
          TextUtils.instagramURL
        )}
      />
      {renderTagsField()}
      {renderStashIDs()}
    </dl>
  );
};
