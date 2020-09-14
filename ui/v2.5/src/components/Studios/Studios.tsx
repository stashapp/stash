import React from "react";
import { Route, Switch } from "react-router-dom";
import { Studio } from "./StudioDetails/Studio";
import { StudioList } from "./StudioList";

const Studios = () => (
  <Switch>
    <Route exact path="/studios" component={StudioList} />
    <Route path="/studios/:id/:tab?" component={Studio} />
  </Switch>
);

export default Studios;
