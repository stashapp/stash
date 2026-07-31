import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { Placement } from "react-bootstrap/esm/Overlay";
import { FormattedMessage } from "react-intl";
import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { Icon } from "src/components/Shared/Icon";
import { PopoverCard } from "src/components/Shared/HoverPopover";
import { ScrapedPerformerCard } from "./ScrapedPerformerCard";

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

const valuesCollide = (remoteValue: unknown, localValue: unknown) => {
  const remoteText = toStringOrNull(remoteValue);
  if (!remoteText) return false;

  return normalizeValue(remoteText) !== normalizeValue(localValue);
};

const getMismatchedStashID = (
  remote: GQL.ScrapedPerformer,
  local: GQL.PerformerDataFragment,
  endpoint?: string
) => {
  if (!endpoint || !remote.remote_site_id) return undefined;

  return local.stash_ids.find(
    (stashID) =>
      stashID.endpoint === endpoint &&
      stashID.stash_id !== remote.remote_site_id
  );
};

export const getPerformerCollisionMessageIds = (
  remote: GQL.ScrapedPerformer,
  local: GQL.PerformerDataFragment,
  endpoint?: string
): string[] => {
  const messageIds: string[] = [];

  if (getMismatchedStashID(remote, local, endpoint)) {
    messageIds.push("tagger.performers.stash_mismatch");
  }
  if (valuesCollide(remote.birthdate, local.birthdate)) {
    messageIds.push("tagger.performers.birthdate");
  }
  if (valuesCollide(remote.country, local.country)) {
    messageIds.push("tagger.performers.country");
  }
  if (valuesCollide(remote.gender, local.gender)) {
    messageIds.push("tagger.performers.gender");
  }
  if (valuesCollide(remote.ethnicity, local.ethnicity)) {
    messageIds.push("tagger.performers.ethnicity");
  }

  return messageIds;
};

export const hasPerformerCollision = (
  remote: GQL.ScrapedPerformer,
  local: GQL.PerformerDataFragment,
  endpoint?: string
) => getPerformerCollisionMessageIds(remote, local, endpoint).length > 0;

export const PerformerCollisionPopoverContent: React.FC<{
  messageIds: string[];
}> = ({ messageIds }) => (
  <PopoverCard className="performer-collision-popover">
    <div className="performer-collision-warnings">
      {messageIds.map((messageId) => (
        <div key={messageId} className="performer-collision-warning">
          <Icon
            className="text-warning performer-collision-warning-icon"
            icon={faTriangleExclamation}
          />
          <FormattedMessage id={messageId} />
        </div>
      ))}
    </div>
  </PopoverCard>
);

interface ITaggerPerformerPopoverProps {
  performer?: GQL.PerformerDataFragment;
  performerID?: string;
  scrapedPerformer?: GQL.ScrapedPerformer;
  endpoint?: string;
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
  placement = "bottom",
  triggerClassName = "d-inline-block",
  onOpen,
  onClose,
  children,
}) => {
  const [isOpened, setIsOpened] = useState(false);

  const { data: selectedPerformerData } = GQL.useFindPerformerQuery({
    variables: { id: performerID ?? "" },
    skip: !performerID || !!performer || !isOpened,
  });
  const localPerformer = performer ?? selectedPerformerData?.findPerformer;

  const wantsLocalCard = !!(performerID || performer);

  const cardContent = localPerformer ? (
    <div className="tag-popover-card">
      <PerformerCard performer={localPerformer} zoomIndex={0} />
    </div>
  ) : wantsLocalCard ? undefined : scrapedPerformer ? (
    <div className="tag-popover-card tagger-performer-popover-card">
      <ScrapedPerformerCard
        scrapedPerformer={scrapedPerformer}
        endpoint={endpoint ?? ""}
        zoomIndex={0}
      />
    </div>
  ) : undefined;

  return (
    <PerformerPopover
      id={cardContent ? undefined : performerID}
      cardContent={cardContent}
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
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
