import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import Performer from "./PerformerDetails/Performer";
import PerformerCreate from "./PerformerDetails/PerformerCreate";
import { PerformerList } from "./PerformerList";

const Performers: React.FC = () => {
  const titleProps = useTitleProps({ id: "performers" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route
          exact
          path="/performers"
          render={(props) => (
            <PerformerList persistState={PersistanceLevel.ALL} {...props} />
          )}
        />
        <Route path="/performers/new" component={PerformerCreate} />
        <Route path="/performers/:id/:tab?" component={Performer} />
      </Switch>
    </>
  );
};
export default Performers;
