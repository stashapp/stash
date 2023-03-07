import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "../Shared/constants";
import { PersistanceLevel } from "../List/ItemList";
import Gallery from "./GalleryDetails/Gallery";
import GalleryCreate from "./GalleryDetails/GalleryCreate";
import { GalleryList } from "./GalleryList";

const Galleries: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "galleries",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />
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
