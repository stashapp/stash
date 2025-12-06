import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Tag from "./TagDetails/Tag";
import TagCreate from "./TagDetails/TagCreate";
import { EnhancedTagList as TagList } from "src/extensions/facets/enhanced";
import { View } from "../List/views";

const Tags: React.FC = () => {
  return <TagList view={View.Tags} />;
};

const TagRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "tags" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/tags" component={Tags} />
        <Route exact path="/tags/new" component={TagCreate} />
        <Route path="/tags/:id/:tab?" component={Tag} />
      </Switch>
    </>
  );
};

export default TagRoutes;
