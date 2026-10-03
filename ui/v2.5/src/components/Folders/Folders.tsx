import type React from "react";
import { Route, Switch } from "react-router-dom";
import { FolderDetails } from "./FolderDetails";
import { RootFolders } from "./RootFolders";

const FoldersRoute: React.FC = () => (
  <Switch>
    <Route exact path="/folders" component={RootFolders} />
    <Route exact path="/folders/:folderId" component={FolderDetails} />
  </Switch>
);

export default FoldersRoute;
