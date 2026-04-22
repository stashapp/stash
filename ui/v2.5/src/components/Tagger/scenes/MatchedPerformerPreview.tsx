import { ReactNode, useState } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { FormattedMessage, IntlShape, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { useConfigurationContext } from "src/hooks/Config";
import { LocalPerformerCard } from "./ScrapedPerformerPreview";

interface IPerformerDeltaRow {
  label: string;
  value: string;
}

interface IMatchedPerformerPreviewProps {
  performerID?: string | null;
  scrapedPerformer: GQL.ScrapedPerformer;
  endpoint?: string;
  placement?: Placement;
  children?: ReactNode;
}

const normalizeValue = (value: unknown) =>
  (() => {
    const text = String(value ?? "").trim();
    if (!text) return "";

    const isNumericLike = /^[-+]?(?:\d+\.?\d*|\.\d+)$/.test(text);
    if (isNumericLike) {
      const numeric = Number(text);
      if (!Number.isNaN(numeric)) {
        return String(numeric);
      }
    }

    return text.toLowerCase();
  })();

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
  local: GQL.PerformerDataFragment,
  intl: IntlShape
): IPerformerDeltaRow[] => {
  const rows: IPerformerDeltaRow[] = [];

  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "birthdate", defaultMessage: "Birthdate" }),
    remote.birthdate,
    local.birthdate
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "death_date", defaultMessage: "Death Date" }),
    remote.death_date,
    local.death_date
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "ethnicity", defaultMessage: "Ethnicity" }),
    remote.ethnicity,
    local.ethnicity
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "hair_color", defaultMessage: "Hair Color" }),
    remote.hair_color,
    local.hair_color
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "eye_color", defaultMessage: "Eye Color" }),
    remote.eye_color,
    local.eye_color
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "height", defaultMessage: "Height" }),
    remote.height,
    local.height_cm
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "weight", defaultMessage: "Weight" }),
    remote.weight,
    local.weight
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "penis_length", defaultMessage: "Penis Length" }),
    remote.penis_length,
    local.penis_length
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "circumcised", defaultMessage: "Circumcised" }),
    remote.circumcised,
    local.circumcised
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "measurements", defaultMessage: "Measurements" }),
    remote.measurements,
    local.measurements
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "fake_tits", defaultMessage: "Fake Tits" }),
    remote.fake_tits,
    local.fake_tits
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "tattoos", defaultMessage: "Tattoos" }),
    remote.tattoos,
    local.tattoos
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "piercings", defaultMessage: "Piercings" }),
    remote.piercings,
    local.piercings
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "career_start", defaultMessage: "Career Start" }),
    remote.career_start,
    local.career_start
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "career_end", defaultMessage: "Career End" }),
    remote.career_end,
    local.career_end
  );

  const remoteAliasesCount = remote.aliases
    ? remote.aliases
        .split(",")
        .map((a) => a.trim())
        .filter(Boolean).length
    : 0;
  const localAliasesCount = local.alias_list?.length ?? 0;
  if (remoteAliasesCount > localAliasesCount) {
    rows.push({
      label: intl.formatMessage({ id: "aliases", defaultMessage: "Aliases" }),
      value: String(remoteAliasesCount),
    });
  }

  const remoteUrlsCount = remote.urls?.length ?? 0;
  const localUrlsCount = local.urls?.length ?? 0;
  if (remoteUrlsCount > localUrlsCount) {
    rows.push({
      label: intl.formatMessage({ id: "urls", defaultMessage: "URLs" }),
      value: String(remoteUrlsCount),
    });
  }

  return rows;
};

export const MatchedPerformerPreview = ({
  performerID,
  scrapedPerformer,
  endpoint,
  placement = "bottom",
  children,
}: IMatchedPerformerPreviewProps) => {
  const intl = useIntl();
  const { configuration: config } = useConfigurationContext();
  const showPerformerCardOnHover = config?.ui.showTagCardOnHover ?? true;
  const [isOpened, setIsOpened] = useState(false);
  const { data: selectedPerformerData, loading: selectedPerformerLoading } =
    GQL.useFindPerformerQuery({
      variables: { id: performerID ?? "" },
      skip: !performerID || !isOpened,
    });
  const performer = selectedPerformerData?.findPerformer;
  const warningStashID =
    endpoint && scrapedPerformer.remote_site_id && performer
      ? performer.stash_ids.find(
          (stashID) =>
            stashID.endpoint === endpoint &&
            stashID.stash_id !== scrapedPerformer.remote_site_id
        )
      : undefined;
  const deltaRows = performer
    ? buildPerformerDeltaRows(scrapedPerformer, performer, intl)
    : [];
  const warningEndpointName = warningStashID
    ? config?.general.stashBoxes.find(
        (sb) => sb.endpoint === warningStashID.endpoint
      )?.name ?? warningStashID.endpoint
    : null;

  if (!performerID || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      className="d-inline-block"
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      onOpen={() => setIsOpened(true)}
      content={
        <div>
          {performer ? (
            <LocalPerformerCard performer={performer} />
          ) : (
            <div className="tag-popover-card tagger-matched-performer-popover p-3">
              {selectedPerformerLoading ? (
                <FormattedMessage
                  id="ui.loading.generic"
                  defaultMessage="Loading..."
                />
              ) : null}
            </div>
          )}
          {(warningStashID || deltaRows.length > 0) && (
            <div className="tagger-matched-performer-popover-extra">
              {warningStashID && (
                <div className="tagger-performer-stashid-warning">
                  <span className="stash-id-pill">
                    <span
                      className="tagger-performer-stashid-warning-chip"
                      title={warningStashID.stash_id}
                    >
                      {warningEndpointName}
                    </span>
                  </span>
                </div>
              )}
              {deltaRows.length > 0 && (
                <div className="tagger-performer-delta-rows mt-2">
                  {deltaRows.map((row) => (
                    <div key={row.label}>
                      <span>{row.label}:</span> <span>{row.value}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      }
    >
      {children}
    </HoverPopover>
  );
};
