import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { Placement } from "react-bootstrap/esm/Overlay";
import { IntlShape, useIntl } from "react-intl";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { useConfigurationContext } from "src/hooks/Config";
import { ScrapedPerformerCard } from "./ScrapedPerformerCard";

interface IPerformerDeltaRow {
  label: string;
  value: string;
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

interface ITaggerPerformerPopoverProps {
  performer?: GQL.PerformerDataFragment;
  performerID?: string;
  scrapedPerformer?: GQL.ScrapedPerformer;
  endpoint?: string;
  cardExtras?: React.ReactNode;
  includeMatchExtras?: boolean;
  placement?: Placement;
  triggerClassName?: string;
  onOpen?: () => void;
  onClose?: () => void;
}

export const TaggerPerformerPopover: React.FC<
  React.PropsWithChildren<ITaggerPerformerPopoverProps>
> = ({
  performer,
  performerID,
  scrapedPerformer,
  endpoint,
  cardExtras,
  includeMatchExtras = false,
  placement = "bottom",
  triggerClassName = "d-inline-block",
  onOpen,
  onClose,
  children,
}) => {
  const intl = useIntl();
  const { configuration: config } = useConfigurationContext();
  const [isOpened, setIsOpened] = useState(false);

  const { data: selectedPerformerData } = GQL.useFindPerformerQuery({
    variables: { id: performerID ?? "" },
    skip: !performerID || !!performer || !isOpened,
  });
  const localPerformer = performer ?? selectedPerformerData?.findPerformer;

  const cardContent = localPerformer ? (
    <div className="tag-popover-card tagger-performer-popover-card">
      <PerformerCard performer={localPerformer} zoomIndex={0} />
    </div>
  ) : scrapedPerformer ? (
    <div className="tag-popover-card tagger-performer-popover-card">
      <ScrapedPerformerCard
        scrapedPerformer={scrapedPerformer}
        endpoint={endpoint ?? ""}
        zoomIndex={0}
      />
    </div>
  ) : undefined;

  const warningStashID =
    includeMatchExtras &&
    endpoint &&
    scrapedPerformer?.remote_site_id &&
    localPerformer
      ? localPerformer.stash_ids.find(
          (stashID) =>
            stashID.endpoint === endpoint &&
            stashID.stash_id !== scrapedPerformer.remote_site_id
        )
      : undefined;

  const deltaRows =
    includeMatchExtras && scrapedPerformer && localPerformer
      ? buildPerformerDeltaRows(scrapedPerformer, localPerformer, intl)
      : [];

  const warningEndpointName = warningStashID
    ? config?.general.stashBoxes.find(
        (sb) => sb.endpoint === warningStashID.endpoint
      )?.name ?? warningStashID.endpoint
    : null;

  const matchExtras =
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
    ) : null;

  return (
    <PerformerPopover
      id={cardContent ? undefined : performerID}
      cardContent={cardContent}
      cardExtras={
        matchExtras || cardExtras ? (
          <>
            {matchExtras}
            {cardExtras}
          </>
        ) : null
      }
      placement={placement}
      enterDelay={1000}
      leaveDelay={500}
      triggerClassName={triggerClassName}
      onOpen={() => {
        setIsOpened(true);
        onOpen?.();
      }}
      onClose={() => {
        setIsOpened(false);
        onClose?.();
      }}
    >
      {children}
    </PerformerPopover>
  );
};
