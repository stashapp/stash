import React from "react";
import { Route, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Movie from "./MovieDetails/Movie";
import MovieCreate from "./MovieDetails/MovieCreate";
import { MovieList } from "./MovieList";

const Movies: React.FC = () => {
  const titleProps = useTitleProps({ id: "movies" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/movies" component={MovieList} />
        <Route exact path="/movies/new" component={MovieCreate} />
        <Route path="/movies/:id/:tab?" component={Movie} />
      </Switch>
    </>
  );
};

export default Movies;
