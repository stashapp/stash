import React from "react";
import { Route, Switch } from "react-router-dom";
import Tag from "./TagDetails/Tag";
import TagCreate from "./TagDetails/TagCreate";
import { TagList } from "./TagList";

const Tags = () => (
  <Switch>
    <Route exact path="/tags" component={TagList} />
    <Route exact path="/tags/new" component={TagCreate} />
    <Route path="/tags/:id/:tab?" component={Tag} />
  </Switch>
);

export default Tags;
