import type React from "react";
import { Helmet } from "react-helmet";
import { FormattedMessage, useIntl } from "react-intl";
import { useParams } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { useTitleProps } from "src/hooks/title";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { Breadcrumbs } from "./Breadcrumbs";
import { FolderNav } from "./FolderNav";
import { ImageResults } from "./ImageResults";
import { SceneResults } from "./SceneResults";

export const FolderDetails: React.FC = () => {
  const intl = useIntl();
  const { folderId } = useParams<{ folderId: string }>();

  const { loading, error, data } = GQL.useFindFolderForBrowserQuery({
    variables: { id: folderId },
  });

  const titleProps = useTitleProps(
    data?.findFolder.basename ?? { id: "folders" }
  );

  if (loading) {
    return <LoadingIndicator />;
  }

  if (error) {
    return <ErrorMessage error={error.message} />;
  }

  if (data === undefined) {
    // should never happen
    return (
      <ErrorMessage error="Unexpected error: No data returned from query." />
    );
  }

  const folder = data.findFolder;

  if (folder.zip_file) {
    return (
      <ErrorMessage
        error={intl.formatMessage({ id: "folder_browser.zip_not_supported" })}
      />
    );
  }

  return (
    <>
      <Helmet {...titleProps} />
      <Breadcrumbs folder={folder} />
      <h3>{folder.basename}</h3>

      <FolderNav
        folders={folder.sub_folders.filter((folder) => !folder.zip_file)}
      />

      <h4>
        <FormattedMessage id="scenes" />
      </h4>
      <SceneResults
        key={`scenes-${folder.id}`}
        folderId={folder.id}
        scenesPerPage={40}
      />

      <h4>
        <FormattedMessage id="images" />
      </h4>
      <ImageResults
        key={`images-${folder.id}`}
        folderId={folder.id}
        imagesPerPage={40}
      />
    </>
  );
};
