import React from "react";
import cx from "classnames";
import { Button, Spinner } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { faBookmark } from "@fortawesome/free-solid-svg-icons";
import { faBookmark as faBookmarkRegular } from "@fortawesome/free-regular-svg-icons";

export interface IWatchQueueButtonProps {
  loading: boolean;
  inQueue: boolean;
  onClick: () => void;
}

export const WatchQueueButton: React.FC<IWatchQueueButtonProps> = ({
  loading,
  inQueue,
  onClick,
}) => {
  if (loading) return <Spinner animation="border" role="status" />;

  return (
    <Button
      variant="secondary"
      title={inQueue ? "Remove from Watch Queue" : "Add to Watch Queue"}
      className={cx(
        "minimal",
        "watch-queue-button",
        inQueue ? "in-queue" : "not-in-queue"
      )}
      onClick={onClick}
    >
      <Icon icon={inQueue ? faBookmark : faBookmarkRegular} />
    </Button>
  );
};
