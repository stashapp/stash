import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Group from "./GroupDetails/Group";
import GroupCreate from "./GroupDetails/GroupCreate";
import { GroupList } from "./GroupList";
import { View } from "../List/views";

const Groups: React.FC = () => {
  return <GroupList view={View.Groups} />;
};

const GroupRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "groups" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/groups" component={Groups} />
        <Route exact path="/groups/new" component={GroupCreate} />
        <Route path="/groups/:id/:tab?" component={Group} />
      </Switch>
    </>
  );
};

export default GroupRoutes;
