import React, { useMemo, useState } from "react";
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

type Conflict = GQL.ScraperUrlConflictsQuery["scraperURLConflicts"][number];
type ConflictPackage = NonNullable<Conflict["scraper_package"]>;

// a conflict as seen from the perspective of one particular scraper/pattern -
// "mine" is the side matching the scraper/url being rendered, "other" is the
// side it conflicts with
export interface ILocalConflict {
  mineID: string;
  minePattern: string;
  minePackage?: ConflictPackage | null;
  otherID: string;
  otherPattern: string;
  otherPackage?: ConflictPackage | null;
}

export function findConflict(
  conflicts: Conflict[] | undefined,
  scraperID: string,
  url: string
): ILocalConflict | undefined {
  const c = conflicts?.find(
    (conflict) =>
      (conflict.scraper_id === scraperID && conflict.pattern === url) ||
      (conflict.other_scraper_id === scraperID &&
        conflict.other_pattern === url)
  );

  if (!c) return undefined;

  if (c.scraper_id === scraperID) {
    return {
      mineID: c.scraper_id,
      minePattern: c.pattern,
      minePackage: c.scraper_package,
      otherID: c.other_scraper_id,
      otherPattern: c.other_pattern,
      otherPackage: c.other_scraper_package,
    };
  }

  return {
    mineID: c.other_scraper_id,
    minePattern: c.other_pattern,
    minePackage: c.other_scraper_package,
    otherID: c.scraper_id,
    otherPattern: c.pattern,
    otherPackage: c.scraper_package,
  };
}

function packageIsOlder(
  a?: ConflictPackage | null,
  b?: ConflictPackage | null
) {
  if (!a?.package.date || !b?.package.date) return false;
  return new Date(a.package.date) < new Date(b.package.date);
}

const ConflictSide: React.FC<{
  name: string;
  pattern: string;
  pkg?: ConflictPackage | null;
  recommended: boolean;
  onUninstall: () => void;
}> = ({ name, pattern, pkg, recommended, onUninstall }) => {
  const intl = useIntl();

  return (
    <div className="scraper-conflict-side">
      <div className="scraper-conflict-side-heading">
        <strong>{name}</strong>
        {recommended && (
          <Badge variant="warning">
            <FormattedMessage id="config.scraping.url_conflict.older" />
          </Badge>
        )}
      </div>
      <div>
        <code>{pattern}</code>
      </div>
      {pkg ? (
        <>
          <div className="text-muted">
            {pkg.package.version ?? pkg.package.package_id}
            {pkg.package.date && (
              <>
                {" "}
                (
                {intl.formatDate(new Date(pkg.package.date), {
                  timeZone: "utc",
                })}
                )
              </>
            )}
          </div>
          <Button
            size="sm"
            variant={recommended ? "danger" : "secondary"}
            onClick={onUninstall}
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
  conflict: ILocalConflict;
  mineName: string;
  otherName: string;
  onHide: () => void;
}> = ({ conflict, mineName, otherName, onHide }) => {
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

  const mineIsOlder = packageIsOlder(
    conflict.minePackage,
    conflict.otherPackage
  );
  const otherIsOlder = packageIsOlder(
    conflict.otherPackage,
    conflict.minePackage
  );

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
          <ConflictSide
            name={mineName}
            pattern={conflict.minePattern}
            pkg={conflict.minePackage}
            recommended={mineIsOlder}
            onUninstall={() =>
              conflict.minePackage &&
              setUninstallTarget({
                name: mineName,
                pkg: conflict.minePackage,
              })
            }
          />
          <ConflictSide
            name={otherName}
            pattern={conflict.otherPattern}
            pkg={conflict.otherPackage}
            recommended={otherIsOlder}
            onUninstall={() =>
              conflict.otherPackage &&
              setUninstallTarget({
                name: otherName,
                pkg: conflict.otherPackage,
              })
            }
          />
        </div>
      </ModalComponent>

      <AlertModal
        show={!!uninstallTarget}
        text={
          <FormattedMessage
            id="config.scraping.url_conflict.confirm_uninstall"
            values={{
              name: uninstallTarget?.name,
              siblings: uninstallTarget?.pkg.sibling_scraper_ids.join(", "),
              hasSiblings:
                (uninstallTarget?.pkg.sibling_scraper_ids.length ?? 0) > 0
                  ? 1
                  : 0,
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
  conflict: ILocalConflict;
  mineName: string;
  otherName: string;
}> = ({ conflict, mineName, otherName }) => {
  const intl = useIntl();
  const [showModal, setShowModal] = useState(false);

  const tooltip = useMemo(
    () =>
      intl.formatMessage(
        { id: "config.scraping.url_conflict.tooltip" },
        { name: otherName }
      ),
    [intl, otherName]
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
          conflict={conflict}
          mineName={mineName}
          otherName={otherName}
          onHide={() => setShowModal(false)}
        />
      )}
    </>
  );
};
