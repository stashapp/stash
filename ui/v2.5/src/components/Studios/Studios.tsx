import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Studio from "./StudioDetails/Studio";
import StudioCreate from "./StudioDetails/StudioCreate";
import { StudioList } from "./StudioList";
import { View } from "../List/views";

const Studios: React.FC = () => {
  return <StudioList view={View.Studios} />;
};

const StudioRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "studios" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/studios" component={Studios} />
        <Route exact path="/studios/new" component={StudioCreate} />
        <Route path="/studios/:id/:tab?" component={Studio} />
      </Switch>
    </>
  );
};

export default StudioRoutes;
