import { ReactNode, useState } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { IntlShape, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";
import { useConfigurationContext } from "src/hooks/Config";
import TextUtils from "src/utils/text";
import { localPerformerToPreviewData } from "./ScrapedPerformerPreview";

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
    intl.formatMessage({ id: "birthdate" }),
    remote.birthdate,
    local.birthdate
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "death_date" }),
    remote.death_date,
    local.death_date
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "ethnicity" }),
    remote.ethnicity,
    local.ethnicity
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "hair_color" }),
    remote.hair_color,
    local.hair_color
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "eye_color" }),
    remote.eye_color,
    local.eye_color
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "height" }),
    remote.height,
    local.height_cm
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "weight" }),
    remote.weight,
    local.weight
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "penis_length" }),
    remote.penis_length,
    local.penis_length
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "circumcised" }),
    remote.circumcised,
    local.circumcised
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "measurements" }),
    remote.measurements,
    local.measurements
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "fake_tits" }),
    remote.fake_tits,
    local.fake_tits
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "tattoos" }),
    remote.tattoos,
    local.tattoos
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "piercings" }),
    remote.piercings,
    local.piercings
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "career_start" }),
    remote.career_start,
    local.career_start
  );
  pushDeltaIfDifferent(
    rows,
    intl.formatMessage({ id: "career_end" }),
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
      label: intl.formatMessage({ id: "aliases" }),
      value: String(remoteAliasesCount),
    });
  }

  const remoteUrlsCount = remote.urls?.length ?? 0;
  const localUrlsCount = local.urls?.length ?? 0;
  if (remoteUrlsCount > localUrlsCount) {
    rows.push({
      label: intl.formatMessage({ id: "urls" }),
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
  const loadingText = intl.formatMessage({
    id: "loading.generic",
    defaultMessage: "Loading...",
  });
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
  const matchedAge = performer
    ? TextUtils.age(performer.birthdate, performer.death_date)
    : 0;
  const matchedAgeString =
    performer && matchedAge !== 0
      ? intl.formatMessage(
          { id: "media_info.performer_card.age" },
          {
            age: matchedAge,
            years_old: intl.formatMessage({
              id: "years_old",
              defaultMessage: "years old",
            }),
          }
        )
      : null;
  const warningEndpointName = warningStashID
    ? config?.general.stashBoxes.find(
        (sb) => sb.endpoint === warningStashID.endpoint
      )?.name ?? warningStashID.endpoint
    : null;

  if (!performerID || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <PerformerPopover
      previewData={
        performer
          ? localPerformerToPreviewData(performer, matchedAgeString)
          : undefined
      }
      loading={selectedPerformerLoading}
      loadingText={loadingText}
      cardExtras={
        warningStashID || deltaRows.length > 0 ? (
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
        ) : null
      }
      triggerClassName="d-inline-block"
      placement={placement}
      onOpen={() => setIsOpened(true)}
      onClose={() => setIsOpened(false)}
    >
      {children}
    </PerformerPopover>
  );
};
