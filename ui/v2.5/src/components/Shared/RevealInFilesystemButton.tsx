import React from "react";
import { Button } from "react-bootstrap";
import { faFolderOpen } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";
import {
  mutateRevealFileInFileManager,
  mutateRevealFolderInFileManager,
} from "src/core/StashService";
import { getPlatformURL } from "src/core/createClient";

interface IRevealInFilesystemButtonProps {
  fileId?: string;
  folderId?: string;
}

function isLocalhost(): boolean {
  const { hostname } = getPlatformURL();
  return (
    hostname === "localhost" || hostname === "127.0.0.1" || hostname === "::1"
  );
}

export const RevealInFilesystemButton: React.FC<
  IRevealInFilesystemButtonProps
> = ({ fileId, folderId }) => {
  const intl = useIntl();

  if (!isLocalhost()) return null;

  function onClick() {
    if (folderId) {
      mutateRevealFolderInFileManager(folderId);
    } else if (fileId) {
      mutateRevealFileInFileManager(fileId);
    }
  }

  return (
    <Button
      className="minimal reveal-in-filesystem-button"
      title={intl.formatMessage({ id: "actions.reveal_in_file_manager" })}
      onClick={onClick}
    >
      <Icon icon={faFolderOpen} />
    </Button>
  );
};
