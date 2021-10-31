import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared";
import { PersistanceLevel } from "src/hooks/ListHook";
import { Image } from "./ImageDetails/Image";
import { ImageList } from "./ImageList";

const Images: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage(
    {
      id: "countables.images",
    },
    { count: 100 }
  )} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />
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
