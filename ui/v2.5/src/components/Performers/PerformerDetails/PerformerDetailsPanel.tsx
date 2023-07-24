import React from "react";
import { useIntl } from "react-intl";
import { TagLink } from "src/components/Shared/TagLink";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { getStashboxBase } from "src/utils/stashbox";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";
import { DetailItem } from "src/components/Shared/DetailItem";
import { CountryFlag } from "src/components/Shared/CountryFlag";

interface IPerformerDetails {
  performer: GQL.PerformerDataFragment;
}

export const PerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    return (
      <ul className="pl-0">
        {(performer.tags ?? []).map((tag) => (
          <TagLink key={tag.id} tagType="performer" tag={tag} />
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

  return (
    <div className="detail-group">
      <DetailItem
        id="gender"
        value={intl.formatMessage({ id: "gender_types." + performer.gender })}
      />
      <DetailItem
        id="age"
        value={TextUtils.age(performer.birthdate, performer.death_date)}
        title={TextUtils.formatDate(intl, performer.birthdate ?? undefined)}
      />
      <DetailItem id="death_date" value={performer.death_date} />
      <DetailItem
        id="country"
        value={
          <CountryFlag
            country={performer.country}
            className="mr-2"
            includeName={true}
          />
        }
      />
      <DetailItem id="ethnicity" value={performer?.ethnicity} />
      <DetailItem id="hair_color" value={performer?.hair_color} />
      <DetailItem id="eye_color" value={performer?.eye_color} />
      <DetailItem id="height" value={formatHeight(performer.height_cm)} />
      <DetailItem id="weight" value={formatWeight(performer.weight)} />
      <DetailItem
        id="penis_length"
        value={formatPenisLength(performer.penis_length)}
      />
      <DetailItem
        id="circumcised"
        value={formatCircumcised(performer.circumcised)}
      />
      <DetailItem id="measurements" value={performer?.measurements} />
      <DetailItem id="fake_tits" value={performer?.fake_tits} />
      <DetailItem id="tattoos" value={performer?.tattoos} />
      <DetailItem id="piercings" value={performer?.piercings} />
      <DetailItem id="details" value={performer?.details} />
      <DetailItem id="tags" value={renderTagsField()} />
      <DetailItem id="StashIDs" value={renderStashIDs()} />
    </div>
  );
};

export const CompressedPerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
}) => {
  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="performer-name" onClick={() => scrollToTop()}>
          {performer.name}
        </a>
        <span className="performer-gender">{performer?.gender}</span>
        <span className="performer-age">
          {TextUtils.age(performer.birthdate, performer.death_date)}
        </span>
        <span className="performer-country">
          <CountryFlag
            country={performer.country}
            className="mr-2"
            includeName={true}
          />
        </span>
      </div>
    </div>
  );
};
