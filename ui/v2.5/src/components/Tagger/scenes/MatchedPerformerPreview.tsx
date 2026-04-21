import React from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";

interface IMatchedPerformerPreviewProps {
  performerID?: string | null;
  placement?: Placement;
}

export const MatchedPerformerPreview: React.FC<IMatchedPerformerPreviewProps> = ({
  performerID,
  placement = "right",
  children,
}) => {
  if (!performerID) {
    return <>{children}</>;
  }

  return (
    <PerformerPopover
      id={performerID}
      placement={placement}
      cardClassName="tagger-matched-performer-popover"
    >
      {children}
    </PerformerPopover>
  );
};

