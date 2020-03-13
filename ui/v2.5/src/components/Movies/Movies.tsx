import React from "react";
import { Route, Switch } from "react-router-dom";
import { Movie } from "./MovieDetails/Movie";
import { MovieList } from "./MovieList";

const Movies = () => (
  <Switch>
    <Route exact={true} path="/movies" component={MovieList} />
    <Route path="/movies/:id" component={Movie} />
  </Switch>
);

export default Movies;
