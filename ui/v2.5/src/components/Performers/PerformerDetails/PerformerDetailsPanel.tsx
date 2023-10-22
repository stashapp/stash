import React from "react";
import { useIntl } from "react-intl";
import { TagLink } from "src/components/Shared/TagLink";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";
import { DetailItem } from "src/components/Shared/DetailItem";
import { CountryFlag } from "src/components/Shared/CountryFlag";
import { StashIDPill } from "src/components/Shared/StashID";

interface IPerformerDetails {
  performer: GQL.PerformerDataFragment;
  collapsed?: boolean;
  fullWidth?: boolean;
}

export const PerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
  collapsed,
  fullWidth,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    if (!performer.tags.length) {
      return;
    }
    return (
      <ul className="pl-0">
        {(performer.tags ?? []).map((tag) => (
          <TagLink key={tag.id} linkType="performer" tag={tag} />
        ))}
      </ul>
    );
  }

  function renderStashIDs() {
    if (!performer.stash_ids.length) {
      return;
    }

    return (
      <ul className="pl-0">
        {performer.stash_ids.map((stashID) => (
          <li key={stashID.stash_id} className="row no-gutters">
            <StashIDPill stashID={stashID} linkType="performers" />
          </li>
        ))}
      </ul>
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

  const formatAge = (birthdate?: string | null, deathdate?: string | null) => {
    if (!birthdate) {
      return "";
    }

    const age = TextUtils.age(birthdate, deathdate);

    return (
      <span className="performer-age">
        <span className="age">{age}</span>
        <span className="birthdate"> ({birthdate})</span>
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

  const formatCircumcised = (circumcised?: GQL.CircumisedEnum | null) => {
    if (!circumcised) {
      return "";
    }

    return (
      <span className="penis-circumcised">
        {intl.formatMessage({
          id: "circumcised_types." + performer.circumcised,
        })}
      </span>
    );
  };

  function maybeRenderExtraDetails() {
    if (!collapsed) {
      /* Remove extra urls provided in details since they will be present by perfomr name */
      /* This code can be removed once multple urls are supported for performers */
      let details = performer?.details
        ?.replace(/\[((?:http|www\.)[^\n\]]+)\]/gm, "")
        .trim();
      return (
        <>
          <DetailItem
            id="tattoos"
            value={performer?.tattoos}
            fullWidth={fullWidth}
          />
          <DetailItem
            id="piercings"
            value={performer?.piercings}
            fullWidth={fullWidth}
          />
          <DetailItem
            id="career_length"
            value={performer?.career_length}
            fullWidth={fullWidth}
          />
          <DetailItem id="details" value={details} fullWidth={fullWidth} />
          <DetailItem
            id="tags"
            value={renderTagsField()}
            fullWidth={fullWidth}
          />
          <DetailItem
            id="stash_ids"
            value={renderStashIDs()}
            fullWidth={fullWidth}
          />
        </>
      );
    }
  }

  return (
    <div className="detail-group">
      {performer.gender ? (
        <DetailItem
          id="gender"
          value={intl.formatMessage({ id: "gender_types." + performer.gender })}
          fullWidth={fullWidth}
        />
      ) : (
        ""
      )}
      <DetailItem
        id="age"
        value={
          !fullWidth
            ? TextUtils.age(performer.birthdate, performer.death_date)
            : formatAge(performer.birthdate, performer.death_date)
        }
        title={
          !fullWidth
            ? TextUtils.formatDate(intl, performer.birthdate ?? undefined)
            : ""
        }
        fullWidth={fullWidth}
      />
      <DetailItem id="death_date" value={performer.death_date} />
      {performer.country ? (
        <DetailItem
          id="country"
          value={
            <CountryFlag
              country={performer.country}
              className="mr-2"
              includeName={true}
            />
          }
          fullWidth={fullWidth}
        />
      ) : (
        ""
      )}
      <DetailItem
        id="ethnicity"
        value={performer?.ethnicity}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="hair_color"
        value={performer?.hair_color}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="eye_color"
        value={performer?.eye_color}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="height"
        value={formatHeight(performer.height_cm)}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="weight"
        value={formatWeight(performer.weight)}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="penis_length"
        value={formatPenisLength(performer.penis_length)}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="circumcised"
        value={formatCircumcised(performer.circumcised)}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="measurements"
        value={performer?.measurements}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="fake_tits"
        value={performer?.fake_tits}
        fullWidth={fullWidth}
      />
      {maybeRenderExtraDetails()}
    </div>
  );
};

export const CompressedPerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
}) => {
  // Network state
  const intl = useIntl();

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="performer-name" onClick={() => scrollToTop()}>
          {performer.name}
        </a>
        {performer.gender ? (
          <>
            <span className="detail-divider">/</span>
            <span className="performer-gender">
              {intl.formatMessage({ id: "gender_types." + performer.gender })}
            </span>
          </>
        ) : (
          ""
        )}
        {performer.birthdate ? (
          <>
            <span className="detail-divider">/</span>
            <span
              className="performer-age"
              title={TextUtils.formatDate(
                intl,
                performer.birthdate ?? undefined
              )}
            >
              {TextUtils.age(performer.birthdate, performer.death_date)}
            </span>
          </>
        ) : (
          ""
        )}
        {performer.country ? (
          <>
            <span className="detail-divider">/</span>
            <span className="performer-country">
              <CountryFlag
                country={performer.country}
                className="mr-2"
                includeName={true}
              />
            </span>
          </>
        ) : (
          ""
        )}
      </div>
    </div>
  );
};
