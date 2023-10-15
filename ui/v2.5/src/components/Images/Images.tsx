import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { PersistanceLevel } from "../List/ItemList";
import Image from "./ImageDetails/Image";
import { ImageList } from "./ImageList";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";

const Images: React.FC = () => {
  useScrollToTopOnMount();

  return <ImageList persistState={PersistanceLevel.ALL} />;
};

const ImageRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "images" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/images" component={Images} />
        <Route path="/images/:id" component={Image} />
      </Switch>
    </>
  );
};

export default ImageRoutes;
