import React from "react";
import { useIntl } from "react-intl";
import { TagLink } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { genderToString } from "src/core/StashService";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

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
        <dt className="col-3 col-xl-2">Tags</dt>
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

  return (
    <>
      <TextField
        name="Gender"
        value={genderToString(performer.gender ?? undefined)}
      />
      <TextField
        name="Birthdate"
        value={TextUtils.formatDate(intl, performer.birthdate ?? undefined)}
      />
      <TextField name="Ethnicity" value={performer.ethnicity} />
      <TextField name="Eye Color" value={performer.eye_color} />
      <TextField name="Country" value={performer.country} />
      <TextField name="Height" value={formatHeight(performer.height)} />
      <TextField name="Measurements" value={performer.measurements} />
      <TextField name="Fake Tits" value={performer.fake_tits} />
      <TextField name="Career Length" value={performer.career_length} />
      <TextField name="Tattoos" value={performer.tattoos} />
      <TextField name="Piercings" value={performer.piercings} />
      <URLField
        name="URL"
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
      {renderTagsField()}
      {renderStashIDs()}
    </>
  );
};
