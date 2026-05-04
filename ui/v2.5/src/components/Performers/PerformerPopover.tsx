import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindPerformer } from "../../core/StashService";
import { useConfigurationContext } from "../../hooks/Config";
import { Placement } from "react-bootstrap/esm/Overlay";
import { PerformerCard } from "./PerformerCard";
import {
  IPerformerPreviewData,
  PerformerPreviewCard,
} from "./PerformerPreviewCard";

interface IPeromerPopoverCardProps {
  id?: string;
  cardContent?: React.ReactNode;
  previewData?: IPerformerPreviewData;
  loading?: boolean;
  loadingText?: string;
  cardExtras?: React.ReactNode;
}

const PerformerPopoverCardByID: React.FC<{
  id: string;
  cardExtras?: React.ReactNode;
}> = ({ id, cardExtras }) => {
  const { data, loading: isLoading, error } = useFindPerformer(id);

  if (isLoading)
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
    <>
      <div className="tag-popover-card">
        <PerformerCard performer={performer} zoomIndex={0} />
      </div>
      {cardExtras}
    </>
  );
};

export const PerformerPopoverCard: React.FC<IPeromerPopoverCardProps> = ({
  id,
  cardContent,
  previewData,
  loading,
  loadingText = "",
  cardExtras,
}) => {
  if (cardContent) {
    return (
      <>
        {cardContent}
        {cardExtras}
      </>
    );
  }

  if (previewData || loading) {
    return (
      <>
        {previewData ? (
          <PerformerPreviewCard {...previewData} />
        ) : (
          <div className="tag-popover-card performer-preview-popover p-3">
            {loading ? loadingText : null}
          </div>
        )}
        {cardExtras}
      </>
    );
  }

  if (!id) return null;
  return <PerformerPopoverCardByID id={id} cardExtras={cardExtras} />;
};

interface IPeroformerPopoverProps {
  id?: string;
  cardContent?: React.ReactNode;
  previewData?: IPerformerPreviewData;
  loading?: boolean;
  loadingText?: string;
  cardExtras?: React.ReactNode;
  hide?: boolean;
  placement?: Placement;
  enterDelay?: number;
  leaveDelay?: number;
  target?: React.RefObject<HTMLElement>;
  triggerClassName?: string;
  onOpen?: () => void;
  onClose?: () => void;
}

export const PerformerPopover: React.FC<IPeroformerPopoverProps> = ({
  id,
  cardContent,
  previewData,
  loading,
  loadingText,
  cardExtras,
  hide,
  children,
  placement = "top",
  enterDelay = 500,
  leaveDelay = 100,
  target,
  triggerClassName,
  onOpen,
  onClose,
}) => {
  const { configuration: config } = useConfigurationContext();

  const showPerformerCardOnHover = config?.ui.showTagCardOnHover ?? true;

  if (hide || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      className={triggerClassName}
      target={target}
      placement={placement}
      enterDelay={enterDelay}
      leaveDelay={leaveDelay}
      onOpen={onOpen}
      onClose={onClose}
      content={
        <PerformerPopoverCard
          id={id}
          cardContent={cardContent}
          previewData={previewData}
          loading={loading}
          loadingText={loadingText}
          cardExtras={cardExtras}
        />
      }
    >
      {children}
    </HoverPopover>
  );
};
