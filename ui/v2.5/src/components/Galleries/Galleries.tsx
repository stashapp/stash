import React from "react";
import { Route, Switch } from "react-router-dom";
import { Gallery } from "./Gallery";
import { GalleryList } from "./GalleryList";

const Galleries = () => (
  <Switch>
    <Route exact path="/galleries" component={GalleryList} />
    <Route path="/galleries/:id" component={Gallery} />
  </Switch>
);

export default Galleries;
