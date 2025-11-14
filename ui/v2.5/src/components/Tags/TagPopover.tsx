import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindTag } from "../../core/StashService";
import { TagCard } from "./TagCard";
import { useConfigurationContext } from "../../hooks/Config";
import { Placement } from "react-bootstrap/esm/Overlay";

interface ITagPopoverCardProps {
  id: string;
}

export const TagPopoverCard: React.FC<ITagPopoverCardProps> = ({ id }) => {
  const { data, loading, error } = useFindTag(id);

  if (loading)
    return (
      <div className="tag-popover-card-placeholder">
        <LoadingIndicator card={true} message={""} />
      </div>
    );
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findTag)
    return <ErrorMessage error={`No tag found with id ${id}.`} />;

  const tag = data.findTag;

  return (
    <div className="tag-popover-card">
      <TagCard tag={tag} zoomIndex={0} />
    </div>
  );
};

interface ITagPopoverProps {
  id: string;
  hide?: boolean;
  placement?: Placement;
  target?: React.RefObject<HTMLElement>;
}

export const TagPopover: React.FC<ITagPopoverProps> = ({
  id,
  hide,
  children,
  placement = "top",
  target,
}) => {
  const { configuration: config } = useConfigurationContext();

  const showTagCardOnHover = config?.ui.showTagCardOnHover ?? true;

  if (hide || !showTagCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      target={target}
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={<TagPopoverCard id={id} />}
    >
      {children}
    </HoverPopover>
  );
};
