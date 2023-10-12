import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindPerformer } from "../../core/StashService";
import { PerformerCard } from "./PerformerCard";
import { ConfigurationContext } from "../../hooks/Config";
import { IUIConfig } from "src/core/config";
import { Placement } from "react-bootstrap/esm/Overlay";

interface IPerformerPopoverCardProps {
  id: string;
}

export const PerformerPopoverCard: React.FC<IPerformerPopoverCardProps> = ({ id }) => {
  const { data, loading, error } = useFindPerformer(id);

  if (loading)
    return (
      <div className="performer-popover-card-placeholder">
        <LoadingIndicator card={true} message={""} />
      </div>
    );
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findPerformer)
    return <ErrorMessage error={`No performer found with id ${id}.`} />;

  const performer = data.findPerformer;

  return (
    <div className="tag-popover-card">
      <PerformerCard performer={performer} />
    </div>
  );
};

interface IPerformerPopoverProps {
  id: string;
  hide?: boolean;
  placement?: Placement;
}

export const PerformerPopover: React.FC<IPerformerPopoverProps> = ({
  id,
  hide,
  children,
  placement = "top",
}) => {
  const { configuration: config } = React.useContext(ConfigurationContext);

  const showPerformerCardOnHover = true;
   // (config?.ui as IUIConfig)?.showPerformerCardOnHover ?? true;

  if (hide || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={<PerformerPopoverCard id={id} />}
    >
      {children}
    </HoverPopover>
  );
};
