import React from "react";
import * as GQL from "src/core/generated-graphql";
import { HoverPopover } from "../Shared/HoverPopover";
import { GalleryCard } from "./GalleryCard";
import { ConfigurationContext } from "../../hooks/Config";
import { Placement } from "react-bootstrap/esm/Overlay";

interface IGalleryPopoverCardProps {
  gallery: GQL.SlimGalleryDataFragment;
}

export const GalleryPopoverCard: React.FC<IGalleryPopoverCardProps> = ({
  gallery,
}) => {
  return (
    <div className="tag-popover-card">
      <GalleryCard gallery={gallery} zoomIndex={0} />
    </div>
  );
};

interface IGalleryPopoverProps {
  gallery: GQL.SlimGalleryDataFragment;
  hide?: boolean;
  placement?: Placement;
  target?: React.RefObject<HTMLElement>;
}

export const GalleryPopover: React.FC<IGalleryPopoverProps> = ({
  gallery,
  hide,
  children,
  placement = "top",
  target,
}) => {
  const { configuration: config } = React.useContext(ConfigurationContext);

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
      content={<GalleryPopoverCard gallery={gallery} />}
    >
      {children}
    </HoverPopover>
  );
};
