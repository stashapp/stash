import { faFolder } from "@fortawesome/free-solid-svg-icons";
import type React from "react";
import { Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import type * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";

export const FolderNav: React.FC<{
  folders: GQL.FolderBrowserFolderDataFragment[];
}> = ({ folders }) => (
  <div className="folder-browser-folders mb-4">
    {folders.map((folder) => (
      <Link key={folder.id} to={`/folders/${folder.id}`}>
        <Button variant="secondary" className="folder-browser-folder">
          <Icon icon={faFolder} />
          <span className="ml-2">{folder.basename}</span>
        </Button>
      </Link>
    ))}
  </div>
);
