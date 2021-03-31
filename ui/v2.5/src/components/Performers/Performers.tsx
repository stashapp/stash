import React from "react";
import { Route, Switch } from "react-router-dom";
import { PersistanceLevel } from "src/hooks/ListHook";
import { Performer } from "./PerformerDetails/Performer";
import { PerformerList } from "./PerformerList";

const Performers = () => (
  <Switch>
    <Route
      exact
      path="/performers"
      render={(props) => (
        <PerformerList persistState={PersistanceLevel.ALL} {...props} />
      )}
    />
    <Route path="/performers/:id/:tab?" component={Performer} />
  </Switch>
);

export default Performers;
