import React, { useEffect, useMemo, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { OptionalField } from "../IncludeButton";
import {
  Performer,
  PerformerSelect,
} from "src/components/Performers/PerformerSelect";
import { getStashboxBase } from "src/utils/stashbox";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { Link } from "react-router-dom";
import { LinkButton } from "../LinkButton";
import { MatchedPerformerPreview } from "./MatchedPerformerPreview";
import { ScrapedPerformerPreview } from "./ScrapedPerformerPreview";

interface IPerformerDeltaRow {
  label: string;
  value: string;
}

const normalizeValue = (value: unknown) =>
  String(value ?? "")
    .trim()
    .toLowerCase();

const toStringOrNull = (value: unknown) => {
  if (value === null || value === undefined) return null;
  const text = String(value).trim();
  return text.length > 0 ? text : null;
};

const pushDeltaIfDifferent = (
  rows: IPerformerDeltaRow[],
  label: string,
  remoteValue: unknown,
  localValue: unknown
) => {
  const remoteText = toStringOrNull(remoteValue);
  if (!remoteText) return;

  if (normalizeValue(remoteText) === normalizeValue(localValue)) return;
  rows.push({ label, value: remoteText });
};

const buildPerformerDeltaRows = (
  remote: GQL.ScrapedPerformer,
  local: GQL.PerformerDataFragment
): IPerformerDeltaRow[] => {
  const rows: IPerformerDeltaRow[] = [];

  pushDeltaIfDifferent(rows, "Birthdate", remote.birthdate, local.birthdate);
  pushDeltaIfDifferent(rows, "Death Date", remote.death_date, local.death_date);
  pushDeltaIfDifferent(rows, "Ethnicity", remote.ethnicity, local.ethnicity);
  pushDeltaIfDifferent(rows, "Hair Color", remote.hair_color, local.hair_color);
  pushDeltaIfDifferent(rows, "Eye Color", remote.eye_color, local.eye_color);
  pushDeltaIfDifferent(rows, "Height", remote.height, local.height_cm);
  pushDeltaIfDifferent(rows, "Weight", remote.weight, local.weight);
  pushDeltaIfDifferent(
    rows,
    "Penis Length",
    remote.penis_length,
    local.penis_length
  );
  pushDeltaIfDifferent(
    rows,
    "Circumcised",
    remote.circumcised,
    local.circumcised
  );
  pushDeltaIfDifferent(
    rows,
    "Measurements",
    remote.measurements,
    local.measurements
  );
  pushDeltaIfDifferent(rows, "Fake Tits", remote.fake_tits, local.fake_tits);
  pushDeltaIfDifferent(rows, "Tattoos", remote.tattoos, local.tattoos);
  pushDeltaIfDifferent(rows, "Piercings", remote.piercings, local.piercings);
  pushDeltaIfDifferent(
    rows,
    "Career Start",
    remote.career_start,
    local.career_start
  );
  pushDeltaIfDifferent(rows, "Career End", remote.career_end, local.career_end);

  const remoteAliasesCount = remote.aliases
    ? remote.aliases
        .split(",")
        .map((a) => a.trim())
        .filter(Boolean).length
    : 0;
  const localAliasesCount = local.alias_list?.length ?? 0;
  if (remoteAliasesCount > localAliasesCount) {
    rows.push({ label: "Aliases", value: String(remoteAliasesCount) });
  }

  const remoteUrlsCount = remote.urls?.length ?? 0;
  const localUrlsCount = local.urls?.length ?? 0;
  if (remoteUrlsCount > localUrlsCount) {
    rows.push({ label: "URLs", value: String(remoteUrlsCount) });
  }

  return rows;
};

const PerformerLink: React.FC<{
  performer: GQL.ScrapedPerformer | Performer;
  url: string | undefined;
  internal?: boolean;
}> = ({ performer, url, internal = false }) => {
  const name = useMemo(() => {
    if (!url) return performer.name;

    return internal ? (
      <Link to={url} target="_blank">
        {performer.name}
      </Link>
    ) : (
      <ExternalLink href={url}>{performer.name}</ExternalLink>
    );
  }, [url, performer.name, internal]);

  return (
    <>
      <span>{name}</span>
      {performer.disambiguation && (
        <span className="performer-disambiguation">
          {` (${performer.disambiguation})`}
        </span>
      )}
    </>
  );
};

interface IPerformerResultProps {
  performer: GQL.ScrapedPerformer;
  selectedID: string | undefined;
  setSelectedID: (id: string | undefined) => void;
  onCreate: () => void;
  onLink?: () => Promise<void>;
  endpoint?: string;
  ageFromDate?: string | null;
}

const PerformerResult: React.FC<IPerformerResultProps> = ({
  performer,
  selectedID,
  setSelectedID,
  onCreate,
  onLink,
  endpoint,
  ageFromDate,
}) => {
  const { data: performerData, loading: stashLoading } =
    GQL.useFindPerformerQuery({
      variables: { id: performer.stored_id ?? "" },
      skip: !performer.stored_id,
    });

  const matchedPerformer = performerData?.findPerformer;
  const matchedStashID = matchedPerformer?.stash_ids.some(
    (stashID) =>
      stashID.endpoint === endpoint &&
      stashID.stash_id === performer.remote_site_id
  );
  const [selectedPerformer, setSelectedPerformer] = useState<Performer>();
  const { data: selectedPerformerData, loading: selectedPerformerLoading } =
    GQL.useFindPerformerQuery({
      variables: { id: selectedID ?? "" },
      skip: !selectedID,
    });
  const selectedPerformerDetails = selectedPerformerData?.findPerformer;

  const stashboxPerformerPrefix = endpoint
    ? `${getStashboxBase(endpoint)}performers/`
    : undefined;
  const performerURLPrefix = "/performers/";

  function selectPerformer(selected: Performer | undefined) {
    setSelectedPerformer(selected);
    setSelectedID(selected?.id);
  }

  useEffect(() => {
    if (
      performerData?.findPerformer &&
      selectedID === performerData?.findPerformer?.id
    ) {
      setSelectedPerformer(performerData.findPerformer);
    }
  }, [performerData?.findPerformer, selectedID]);

  const handleSelect = (performers: Performer[]) => {
    if (performers.length) {
      selectPerformer(performers[0]);
    } else {
      selectPerformer(undefined);
    }
  };

  const handleSkip = () => {
    selectPerformer(undefined);
  };

  if (stashLoading || selectedPerformerLoading)
    return <div>Loading performer</div>;

  if (matchedPerformer && matchedStashID) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          <FormattedMessage id="countables.performers" values={{ count: 1 }} />:
          <b className="ml-2">
            <ScrapedPerformerPreview performer={performer}>
              <PerformerLink
                performer={performer}
                url={`${stashboxPerformerPrefix}${performer.remote_site_id}`}
              />
            </ScrapedPerformerPreview>
          </b>
        </div>
        <span className="ml-auto">
          <OptionalField
            exclude={selectedID === undefined}
            setExclude={(v) =>
              v ? handleSkip() : setSelectedID(matchedPerformer.id)
            }
          >
            <div>
              <span className="mr-2">
                <FormattedMessage id="component_tagger.verb_matched" />:
              </span>
              <b className="col-3 text-right">
                <PerformerLink
                  performer={matchedPerformer}
                  url={`${performerURLPrefix}${matchedPerformer.id}`}
                  internal
                />
              </b>
            </div>
          </OptionalField>
        </span>
      </div>
    );
  }

  const selectedSource = !selectedID ? "skip" : "existing";
  const selectedPerformerConflictStashID =
    endpoint && performer.remote_site_id && selectedPerformerDetails
      ? selectedPerformerDetails.stash_ids.find(
          (stashID) =>
            stashID.endpoint === endpoint &&
            stashID.stash_id !== performer.remote_site_id
        )
      : undefined;
  const selectedPerformerDeltaRows = selectedPerformerDetails
    ? buildPerformerDeltaRows(performer, selectedPerformerDetails)
    : [];

  const safeBuildPerformerScraperLink = (id: string | null | undefined) => {
    return stashboxPerformerPrefix && id
      ? `${stashboxPerformerPrefix}${id}`
      : undefined;
  };

  return (
    <div className="row no-gutters align-items-center mt-2">
      <div className="entity-name">
        <FormattedMessage id="countables.performers" values={{ count: 1 }} />:
        <b className="ml-2">
          <ScrapedPerformerPreview performer={performer}>
            <PerformerLink
              performer={performer}
              url={safeBuildPerformerScraperLink(performer.remote_site_id)}
            />
          </ScrapedPerformerPreview>
        </b>
      </div>
      <ButtonGroup>
        <Button variant="secondary" onClick={() => onCreate()}>
          <FormattedMessage id="actions.create" />
        </Button>
        <Button
          variant={selectedSource === "skip" ? "primary" : "secondary"}
          onClick={() => handleSkip()}
        >
          <FormattedMessage id="actions.skip" />
        </Button>
        <MatchedPerformerPreview
          performerID={selectedPerformer?.id}
          performer={selectedPerformerDetails}
          warningStashID={selectedPerformerConflictStashID}
          deltaRows={selectedPerformerDeltaRows}
        >
          <PerformerSelect
            values={selectedPerformer ? [selectedPerformer] : []}
            onSelect={handleSelect}
            active={selectedSource === "existing"}
            isClearable={false}
            ageFromDate={ageFromDate}
          />
        </MatchedPerformerPreview>
        {endpoint && onLink && (
          <LinkButton disabled={selectedID === undefined} onLink={onLink} />
        )}
      </ButtonGroup>
    </div>
  );
};

export default PerformerResult;
