import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import { Image } from "./ImageDetails/Image";
import { ImageList } from "./ImageList";

const Images: React.FC = () => {
  const titleProps = useTitleProps({ id: "images" });
  return (
    <>
      <Helmet {...titleProps} />
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
    </>
  );
};

export default Images;
