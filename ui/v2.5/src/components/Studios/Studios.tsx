import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import Studio from "./StudioDetails/Studio";
import StudioCreate from "./StudioDetails/StudioCreate";
import { StudioList } from "./StudioList";

const Studios: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "studios",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />
      <Switch>
        <Route exact path="/studios" component={StudioList} />
        <Route exact path="/studios/new" component={StudioCreate} />
        <Route path="/studios/:id/:tab?" component={Studio} />
      </Switch>
    </>
  );
};
export default Studios;
