import type React from "react";
import { Helmet } from "react-helmet";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useTitleProps } from "src/hooks/title";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { FolderNav } from "./FolderNav";

export const RootFolders: React.FC = () => {
  const intl = useIntl();
  const titleProps = useTitleProps({ id: "root_folders" });
  const { loading, error, data } = GQL.useFindRootFoldersForBrowserQuery();

  if (loading) {
    return <LoadingIndicator />;
  }

  if (error) {
    return <ErrorMessage error={error.message} />;
  }

  return (
    <>
      <Helmet {...titleProps} />
      <h3>{intl.formatMessage({ id: "root_folders" })}</h3>
      <FolderNav folders={data?.findFolders.folders ?? []} />
    </>
  );
};
