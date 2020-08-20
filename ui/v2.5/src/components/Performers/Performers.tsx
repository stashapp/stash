import React from "react";
import { Route, Switch } from "react-router-dom";
import { Performer } from "./PerformerDetails/Performer";
import { PerformerList } from "./PerformerList";

const Performers = () => (
  <Switch>
    <Route exact path="/performers" component={PerformerList} />
    <Route path="/performers/:id/:tab?" component={Performer} />
  </Switch>
);

export default Performers;
