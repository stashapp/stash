import React from "react";
import { Route, Switch } from "react-router-dom";
import { Image } from "./ImageDetails/Image";
import { ImageList } from "./ImageList";

const Images = () => (
  <Switch>
    <Route
      exact
      path="/images"
      render={(props) => <ImageList persistState {...props} />}
    />
    <Route path="/images/:id" component={Image} />
  </Switch>
);

export default Images;
