import React from "react";
import { Route, Switch } from "react-router-dom";
import { Dvd } from "./DvdDetails/Dvd";
import { DvdList } from "./DvdList";

const Dvds = () => (
  <Switch>
    <Route exact={true} path="/dvds" component={DvdList} />
    <Route path="/dvds/:id" component={Dvd} />
  </Switch>
);

export default Dvds;
