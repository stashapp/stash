import React from "react";
import { Route, Switch } from "react-router-dom";
import { PersistanceLevel } from "src/hooks/ListHook";
import { Image } from "./ImageDetails/Image";
import { ImageList } from "./ImageList";

const Images = () => (
  <Switch>
    <Route
      exact
      path="/images"
      render={(props) => (
        <ImageList persistState={PersistanceLevel.ALL} {...props} />
      )}
    />
    <Route path="/images/:id" component={Image} />
  </Switch>
);

export default Images;
