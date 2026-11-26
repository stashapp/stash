import React, { useState } from "react";
import { Badge, Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  mutateUninstallScraperPackages,
  scraperMutationImpactedQueries,
} from "src/core/StashService";
import { useMonitorJob } from "src/utils/job";
import { Icon } from "../Shared/Icon";
import { WarningHoverPopover } from "../Shared/HoverPopover";
import { AlertModal } from "../Shared/Alert";
import { ModalComponent } from "../Shared/Modal";
import {
  faExclamationTriangle,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";

export type ScraperURLConflict =
  GQL.ScraperUrlConflictsQuery["scraperURLConflicts"][number];
type ConflictSide = GQL.ScraperUrlConflictSideDataFragment;
type ConflictPackage = NonNullable<ConflictSide["package"]>;

// all conflicts involving one scraper's URL pattern
export interface IPatternConflicts {
  pattern: ConflictSide;
  conflictsWith: ConflictSide[];
}

export function findConflicts(
  conflicts: ScraperURLConflict[] | undefined,
  scraperID: string,
  url: string
): IPatternConflicts | undefined {
  const isPattern = (s: ConflictSide) =>
    s.scraper_id === scraperID && s.pattern === url;

  let pattern: ConflictSide | undefined;
  const conflictsWith: ConflictSide[] = [];
  for (const c of conflicts ?? []) {
    for (const [side, other] of [
      [c.scraper, c.other],
      [c.other, c.scraper],
    ]) {
      if (isPattern(side)) {
        pattern = side;
        conflictsWith.push(other);
      }
    }
  }

  return pattern && { pattern, conflictsWith };
}

// returns the scrapers whose package is older than the newest package in the
// group. Returns none unless every side has a dated package, since the
// suggestion would otherwise be misleading.
// shortcut: package date is when the installed version was published, not
// when the scraper was superseded; an updated but deprecated scraper can look
// newer than its replacement
function olderScraperIDs(sides: ConflictSide[]) {
  const dates = sides.map((s) => s.package?.package.date);
  if (dates.some((d) => !d)) return new Set<string>();

  const times = dates.map((d) => new Date(d!).getTime());
  const newest = Math.max(...times);
  return new Set(
    sides.filter((_, i) => times[i] < newest).map((s) => s.scraper_id)
  );
}

const ConflictSideDetails: React.FC<{
  side: ConflictSide;
  older: boolean;
  onUninstall: (pkg: ConflictPackage) => void;
}> = ({ side, older, onUninstall }) => {
  const intl = useIntl();
  const { package: pkg, scraper_name: name } = side;

  return (
    <div className="scraper-conflict-side">
      <div className="scraper-conflict-side-heading">
        <strong>{name}</strong>
        {older && (
          <Badge variant="warning">
            <FormattedMessage id="config.scraping.url_conflict.older" />
          </Badge>
        )}
      </div>
      <div>
        <code>{side.pattern}</code>
      </div>
      {pkg ? (
        <>
          <div className="text-muted">
            {pkg.package.version ?? pkg.package.package_id}
            {pkg.package.date &&
              ` (${intl.formatDate(new Date(pkg.package.date), { timeZone: "utc" })})`}
          </div>
          <Button
            size="sm"
            variant={older ? "danger" : "secondary"}
            onClick={() => onUninstall(pkg)}
          >
            <Icon icon={faTrash} />
            <FormattedMessage
              id="config.scraping.url_conflict.uninstall"
              values={{ name }}
            />
          </Button>
        </>
      ) : (
        <div className="text-muted">
          <FormattedMessage id="config.scraping.url_conflict.no_package" />
        </div>
      )}
    </div>
  );
};

const ScraperConflictModal: React.FC<{
  conflicts: IPatternConflicts;
  onHide: () => void;
}> = ({ conflicts, onHide }) => {
  const intl = useIntl();
  const [uninstallTarget, setUninstallTarget] = useState<{
    name: string;
    pkg: ConflictPackage;
  }>();
  const [jobID, setJobID] = useState<string>();
  const { job } = useMonitorJob(jobID, () => {
    const ac = getClient();
    evictQueries(ac.cache, scraperMutationImpactedQueries);
    onHide();
  });

  const sides = [conflicts.pattern, ...conflicts.conflictsWith];
  const older = olderScraperIDs(sides);
  const siblings = uninstallTarget?.pkg.sibling_scraper_ids ?? [];

  async function confirmUninstall() {
    if (!uninstallTarget) return;

    const { package_id, sourceURL } = uninstallTarget.pkg.package;
    const r = await mutateUninstallScraperPackages([
      { id: package_id, sourceURL },
    ]);
    setJobID(r.data?.uninstallPackages);
    setUninstallTarget(undefined);
  }

  return (
    <>
      <ModalComponent
        show
        icon={faExclamationTriangle}
        header={intl.formatMessage({
          id: "config.scraping.url_conflict.heading",
        })}
        onHide={onHide}
        accept={{ onClick: onHide, variant: "secondary" }}
        disabled={!!job}
        modalProps={{ size: "lg" }}
      >
        <p>
          <FormattedMessage id="config.scraping.url_conflict.description" />
        </p>
        <div className="scraper-conflict-sides">
          {sides.map((side) => (
            <ConflictSideDetails
              key={`${side.scraper_id}:${side.pattern}`}
              side={side}
              older={older.has(side.scraper_id)}
              onUninstall={(pkg) =>
                setUninstallTarget({ name: side.scraper_name, pkg })
              }
            />
          ))}
        </div>
      </ModalComponent>

      <AlertModal
        show={!!uninstallTarget}
        text={
          <FormattedMessage
            id="config.scraping.url_conflict.confirm_uninstall"
            values={{
              name: uninstallTarget?.name,
              siblings: siblings.join(", "),
              siblingCount: siblings.length,
            }}
          />
        }
        onConfirm={confirmUninstall}
        onCancel={() => setUninstallTarget(undefined)}
      />
    </>
  );
};

export const ConflictBadge: React.FC<{
  conflicts: IPatternConflicts;
}> = ({ conflicts }) => {
  const intl = useIntl();
  const [showModal, setShowModal] = useState(false);

  const names = new Set(conflicts.conflictsWith.map((s) => s.scraper_name));
  const tooltip = intl.formatMessage(
    { id: "config.scraping.url_conflict.tooltip" },
    { names: intl.formatList([...names]) }
  );

  return (
    <>
      <button
        type="button"
        className="scraper-url-conflict-badge"
        onClick={() => setShowModal(true)}
      >
        <WarningHoverPopover content={tooltip} placement="top">
          <span />
        </WarningHoverPopover>
      </button>
      {showModal && (
        <ScraperConflictModal
          conflicts={conflicts}
          onHide={() => setShowModal(false)}
        />
      )}
    </>
  );
};
