import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import { PersistanceLevel } from "../List/ItemList";
import Performer from "./PerformerDetails/Performer";
import PerformerCreate from "./PerformerDetails/PerformerCreate";
import { PerformerList } from "./PerformerList";

const Performers: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "performers",
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
          path="/performers"
          render={(props) => (
            <PerformerList persistState={PersistanceLevel.ALL} {...props} />
          )}
        />
        <Route path="/performers/new" component={PerformerCreate} />
        <Route path="/performers/:id/:tab?" component={Performer} />
      </Switch>
    </>
  );
};
export default Performers;
