import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindTag } from "../../core/StashService";
import { TagCard } from "./TagCard";
import { ConfigurationContext } from "../../hooks/Config";
import { IUIConfig } from "src/core/config";
import { Placement } from "react-bootstrap/esm/Overlay";

interface ITagPopoverProps {
  id?: string;
  placement?: Placement;
}

export const TagPopoverCard: React.FC<ITagPopoverCardProps> = ({ id }) => {
  const { data, loading, error } = useFindTag(id ?? "");

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

export const TagPopover: React.FC<ITagPopoverProps> = ({
  id,
  children,
  placement = "top",
}) => {
  const { configuration: config } = React.useContext(ConfigurationContext);

  const showTagCardOnHover =
    (config?.ui as IUIConfig)?.showTagCardOnHover ?? true;

  if (!id || !showTagCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={<TagPopoverCard id={id} />}
    >
      {children}
    </HoverPopover>
  );
};

interface ITagPopoverCardProps {
  id?: string;
}
