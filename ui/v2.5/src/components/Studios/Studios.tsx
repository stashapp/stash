import React from "react";
import { Route, Switch } from "react-router-dom";
import { Studio } from "./StudioDetails/Studio";
import { StudioList } from "./StudioList";

const Studios = () => (
  <Switch>
    <Route exact={true} path="/studios" component={StudioList} />
    <Route path="/studios/:id" component={Studio} />
  </Switch>
);

export default Studios;
