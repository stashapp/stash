import React from "react";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { HoverPopover } from "../Shared/HoverPopover";
import { useFindPerformer } from "../../core/StashService";
import { useConfigurationContext } from "../../hooks/Config";
import { Placement } from "react-bootstrap/esm/Overlay";
import { IPerformerPreviewData, PerformerPreviewCard } from "./PerformerPreviewCard";

interface IPeromerPopoverCardProps {
  id?: string;
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
      <PerformerPreviewCard
        name={performer.name}
        image={performer.image_path}
        country={performer.country}
        gender={performer.gender}
        disambiguation={performer.disambiguation}
      />
      {cardExtras}
    </>
  );
};

export const PerformerPopoverCard: React.FC<IPeromerPopoverCardProps> = ({
  id,
  previewData,
  loading,
  loadingText = "",
  cardExtras,
}) => {
  if (previewData || loading) {
    return (
      <>
        {previewData ? (
          <PerformerPreviewCard {...previewData} />
        ) : (
          <div className="tag-popover-card tagger-performer-popover p-3">
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
  previewData?: IPerformerPreviewData;
  loading?: boolean;
  loadingText?: string;
  cardExtras?: React.ReactNode;
  hide?: boolean;
  placement?: Placement;
  target?: React.RefObject<HTMLElement>;
  triggerClassName?: string;
  onOpen?: () => void;
  onClose?: () => void;
}

export const PerformerPopover: React.FC<IPeroformerPopoverProps> = ({
  id,
  previewData,
  loading,
  loadingText,
  cardExtras,
  hide,
  children,
  placement = "top",
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
      enterDelay={500}
      leaveDelay={100}
      onOpen={onOpen}
      onClose={onClose}
      content={
        <PerformerPopoverCard
          id={id}
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
