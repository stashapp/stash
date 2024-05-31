import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Movie from "./MovieDetails/Movie";
import MovieCreate from "./MovieDetails/MovieCreate";
import { MovieList } from "./MovieList";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { View } from "../List/views";

const Movies: React.FC = () => {
  useScrollToTopOnMount();

  return <MovieList view={View.Movies} />;
};

const MovieRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "groups" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/groups" component={Movies} />
        <Route exact path="/groups/new" component={MovieCreate} />
        <Route path="/groups/:id/:tab?" component={Movie} />
      </Switch>
    </>
  );
};

export default MovieRoutes;
