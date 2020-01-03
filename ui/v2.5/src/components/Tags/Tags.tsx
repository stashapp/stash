import React from "react";
import { Route, Switch } from "react-router-dom";
import { TagList } from "./TagList";

const Tags = () => (
  <Switch>
    <Route exact={true} path="/tags" component={TagList} />
  </Switch>
);

export default Tags;
