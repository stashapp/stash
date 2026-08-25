import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { lazyComponent } from "src/utils/lazyComponent";
import { View } from "../List/views";

const AudioList = lazyComponent(() => import("./AudioList"));
const Audio = lazyComponent(() => import("./AudioDetails/Audio"));

const Audios: React.FC = () => {
  return <AudioList view={View.Audios} />;
};

const AudioRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "audios" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/audios" component={Audios} />
        <Route path="/audios/:id" component={Audio} />
      </Switch>
    </>
  );
};

export default AudioRoutes;
