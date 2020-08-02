import React from "react";
import { Route, Switch } from "react-router-dom";
import { Tag } from "./TagDetails/Tag";
import { TagList } from "./TagList";

const Tags = () => (
  <Switch>
    <Route exact path="/tags" component={TagList} />
    <Route path="/tags/:id" component={Tag} />
  </Switch>
);

export default Tags;
