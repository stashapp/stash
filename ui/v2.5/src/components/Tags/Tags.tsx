import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import Tag from "./TagDetails/Tag";
import TagCreate from "./TagDetails/TagCreate";
import { TagList } from "./TagList";

const Tags: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "tags",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />

      <Switch>
        <Route exact path="/tags" component={TagList} />
        <Route exact path="/tags/new" component={TagCreate} />
        <Route path="/tags/:id/:tab?" component={Tag} />
      </Switch>
    </>
  );
};
export default Tags;
