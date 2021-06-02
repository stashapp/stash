import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { TagLink } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { genderToString } from "src/core/StashService";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";

interface IPerformerDetails {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    if (!performer.tags?.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2"><FormattedMessage id="tags" /></dt>
        <dd className="col-9 col-xl-10">
          <ul className="pl-0">
            {(performer.tags ?? []).map((tag) => (
              <TagLink key={tag.id} tagType="performer" tag={tag} />
            ))}
          </ul>
        </dd>
      </dl>
    );
  }

  function renderRating() {
    if (!performer.rating) {
      return null;
    }

    return (
      <dl className="row mb-0">
        <dt className="col-3 col-xl-2"><FormattedMessage id="rating" />:</dt>
        <dd className="col-9 col-xl-10">
          <RatingStars value={performer.rating} />
        </dd>
      </dl>
    );
  }

  function renderStashIDs() {
    if (!performer.stash_ids?.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">StashIDs</dt>
        <dd className="col-9 col-xl-10">
          <ul className="pl-0">
            {performer.stash_ids.map((stashID) => {
              const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
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
      </dl>
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
    <>
      <TextField
        id="gender"
        name={intl.formatMessage({ id:'performer_bio.gender' })}
        value={genderToString(performer.gender ?? undefined)}
      />
      <TextField
        id="birthdate"
        name={intl.formatMessage({ id:'performer_bio.birthdate' })}
        value={TextUtils.formatDate(intl, performer.birthdate ?? undefined)}
      />
      <TextField
        id="death_date"
        name={intl.formatMessage({ id:'performer_bio.death_date' })}
        value={TextUtils.formatDate(intl, performer.death_date ?? undefined)}
      />
      <TextField id="ethnicity" name={intl.formatMessage({ id:'performer_bio.ethnicity' })} value={performer.ethnicity} />
      <TextField id="hair_color" name={intl.formatMessage({ id:'performer_bio.hair_color' })} value={performer.hair_color} />
      <TextField id="eye_color" name={intl.formatMessage({ id:'performer_bio.eye_color' })} value={performer.eye_color} />
      <TextField id="country" name={intl.formatMessage({ id:'country' })} value={performer.country} />
      <TextField id="height" name={intl.formatMessage({ id:'performer_bio.height' })} value={formatHeight(performer.height)} />
      <TextField id="weight" name={intl.formatMessage({ id:'performer_bio.weight' })} value={formatWeight(performer.weight)} />
      <TextField id="measurements" name={intl.formatMessage({ id:'performer_bio.measurements' })} value={performer.measurements} />
      <TextField id="fake_tits" name={intl.formatMessage({ id:'performer_bio.fake_tits' })} value={performer.fake_tits} />
      <TextField id="career_length" name={intl.formatMessage({ id:'performer_bio.career_length' })} value={performer.career_length} />
      <TextField id="tattoos" name={intl.formatMessage({ id:'performer_bio.tattoos' })} value={performer.tattoos} />
      <TextField id="piercings" name={intl.formatMessage({ id:'performer_bio.piercings' })} value={performer.piercings} />
      <TextField id="details" name={intl.formatMessage({ id:'details' })} value={performer.details} />
      <URLField
        id="url"
        name={intl.formatMessage({ id:'url' })}
        value={performer.url}
        url={TextUtils.sanitiseURL(performer.url ?? "")}
      />
      <URLField
        name="Twitter"
        value={performer.twitter}
        url={TextUtils.sanitiseURL(
          performer.twitter ?? "",
          TextUtils.twitterURL
        )}
      />
      <URLField
        name="Instagram"
        value={performer.instagram}
        url={TextUtils.sanitiseURL(
          performer.instagram ?? "",
          TextUtils.instagramURL
        )}
      />
      {renderRating()}
      {renderTagsField()}
      {renderStashIDs()}
    </>
  );
};
