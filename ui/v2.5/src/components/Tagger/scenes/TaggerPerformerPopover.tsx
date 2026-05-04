import React from "react";
import * as GQL from "src/core/generated-graphql";
import { Placement } from "react-bootstrap/esm/Overlay";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { ScrapedPerformerCard } from "./ScrapedPerformerCard";

interface ITaggerPerformerPopoverProps {
  performer?: GQL.PerformerDataFragment;
  performerID?: string;
  scrapedPerformer?: GQL.ScrapedPerformer;
  endpoint?: string;
  cardExtras?: React.ReactNode;
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
  placement = "bottom",
  triggerClassName = "d-inline-block",
  onOpen,
  onClose,
  children,
}) => {
  const cardContent = performer ? (
    <div className="tag-popover-card">
      <PerformerCard performer={performer} zoomIndex={0} />
    </div>
  ) : scrapedPerformer ? (
    <div className="tag-popover-card">
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
      cardExtras={cardExtras}
      placement={placement}
      enterDelay={1000}
      leaveDelay={500}
      triggerClassName={triggerClassName}
      onOpen={onOpen}
      onClose={onClose}
    >
      {children}
    </PerformerPopover>
  );
};
