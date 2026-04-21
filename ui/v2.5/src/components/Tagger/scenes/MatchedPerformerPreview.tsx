import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";

interface IMatchedPerformerPreviewProps {
  performerID?: string | null;
  placement?: Placement;
  children?: ReactNode;
}

export const MatchedPerformerPreview = ({
  performerID,
  placement = "right",
  children,
}: IMatchedPerformerPreviewProps) => {
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
