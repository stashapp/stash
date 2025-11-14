import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindPerformer } from "../../core/StashService";
import { PerformerCard } from "./PerformerCard";
import { useConfigurationContext } from "../../hooks/Config";
import { Placement } from "react-bootstrap/esm/Overlay";

interface IPeromerPopoverCardProps {
  id: string;
}

export const PerformerPopoverCard: React.FC<IPeromerPopoverCardProps> = ({
  id,
}) => {
  const { data, loading, error } = useFindPerformer(id);

  if (loading)
    return (
      <div className="tag-popover-card-placeholder">
        <LoadingIndicator card={true} message={""} />
      </div>
    );
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findPerformer)
    return <ErrorMessage error={`No tag found with id ${id}.`} />;

  const performer = data.findPerformer;

  return (
    <div className="tag-popover-card">
      <PerformerCard performer={performer} zoomIndex={0} />
    </div>
  );
};

interface IPeroformerPopoverProps {
  id: string;
  hide?: boolean;
  placement?: Placement;
  target?: React.RefObject<HTMLElement>;
}

export const PerformerPopover: React.FC<IPeroformerPopoverProps> = ({
  id,
  hide,
  children,
  placement = "top",
  target,
}) => {
  const { configuration: config } = useConfigurationContext();

  const showPerformerCardOnHover = config?.ui.showTagCardOnHover ?? true;

  if (hide || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      target={target}
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={<PerformerPopoverCard id={id} />}
    >
      {children}
    </HoverPopover>
  );
};
