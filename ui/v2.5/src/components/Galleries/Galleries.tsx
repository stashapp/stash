import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import Gallery from "./GalleryDetails/Gallery";
import GalleryCreate from "./GalleryDetails/GalleryCreate";
import { GalleryList } from "./GalleryList";

const Galleries: React.FC = () => {
  const titleProps = useTitleProps({ id: "galleries" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route
          exact
          path="/galleries"
          render={(props) => (
            <GalleryList {...props} persistState={PersistanceLevel.ALL} />
          )}
        />
        <Route exact path="/galleries/new" component={GalleryCreate} />
        <Route path="/galleries/:id/:tab?" component={Gallery} />
      </Switch>
    </>
  );
};

export default Galleries;
