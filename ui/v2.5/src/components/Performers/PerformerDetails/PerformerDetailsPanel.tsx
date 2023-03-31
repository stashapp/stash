import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { TagLink } from "src/components/Shared/TagLink";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { getStashboxBase } from "src/utils/stashbox";
import { getCountryByISO } from "src/utils/country";
import { TextField, URLField } from "src/utils/field";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";

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

  const formatHeight = (height?: number | null) => {
    if (!height) {
      return "";
    }

    const [feet, inches] = cmToImperial(height);

    return (
      <span className="performer-height">
        <span className="height-metric">
          {intl.formatNumber(height, {
            style: "unit",
            unit: "centimeter",
            unitDisplay: "short",
          })}
        </span>
        <span className="height-imperial">
          {intl.formatNumber(feet, {
            style: "unit",
            unit: "foot",
            unitDisplay: "narrow",
          })}
          {intl.formatNumber(inches, {
            style: "unit",
            unit: "inch",
            unitDisplay: "narrow",
          })}
        </span>
      </span>
    );
  };

  const formatWeight = (weight?: number | null) => {
    if (!weight) {
      return "";
    }

    const lbs = kgToLbs(weight);

    return (
      <span className="performer-weight">
        <span className="weight-metric">
          {intl.formatNumber(weight, {
            style: "unit",
            unit: "kilogram",
            unitDisplay: "short",
          })}
        </span>
        <span className="weight-imperial">
          {intl.formatNumber(lbs, {
            style: "unit",
            unit: "pound",
            unitDisplay: "short",
          })}
        </span>
      </span>
    );
  };

  const formatPenisLength = (penis_length?: number | null) => {
    if (!penis_length) {
      return "";
    }

    const inches = cmToInches(penis_length);

    return (
      <span className="performer-penis-length">
        <span className="penis-length-metric">
          {intl.formatNumber(penis_length, {
            style: "unit",
            unit: "centimeter",
            unitDisplay: "short",
            maximumFractionDigits: 2,
          })}
        </span>
        <span className="penis-length-imperial">
          {intl.formatNumber(inches, {
            style: "unit",
            unit: "inch",
            unitDisplay: "narrow",
            maximumFractionDigits: 2,
          })}
        </span>
      </span>
    );
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
      <TextField
        id="country"
        value={
          getCountryByISO(performer.country, intl.locale) ?? performer.country
        }
      />

      {!!performer.height_cm && (
        <>
          <dt>
            <FormattedMessage id="height" />
          </dt>
          <dd>{formatHeight(performer.height_cm)}</dd>
        </>
      )}

      {!!performer.weight && (
        <>
          <dt>
            <FormattedMessage id="weight" />
          </dt>
          <dd>{formatWeight(performer.weight)}</dd>
        </>
      )}

      {!!performer.penis_length && (
        <>
          <dt>
            <FormattedMessage id="penis_length" />
          </dt>
          <dd>{formatPenisLength(performer.penis_length)}</dd>
        </>
      )}
      <TextField
        id="circumcised"
        value={
          performer.circumcised
            ? intl.formatMessage({
                id: "circumcised_types." + performer.circumcised,
              })
            : undefined
        }
      />
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
