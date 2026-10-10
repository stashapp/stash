import React from "react";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";
import type * as GQL from "src/core/generated-graphql";

export const Breadcrumbs: React.FC<{
  folder: GQL.FindFolderForBrowserQuery["findFolder"];
}> = ({ folder }) => (
  <div className="mb-3">
    <Link to="/folders">
      <FormattedMessage id="root_folders" />
    </Link>
    {folder.parent_folders.toReversed().map((parent) => (
      <React.Fragment key={parent.id}>
        <span className="mx-2">/</span>
        <Link to={`/folders/${parent.id}`}>{parent.basename}</Link>
      </React.Fragment>
    ))}
    <span className="mx-2">/</span>
    <span>{folder.basename}</span>
  </div>
);
