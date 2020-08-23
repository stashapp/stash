import React from "react";
import { Route, Switch } from "react-router-dom";
import { Movie } from "./MovieDetails/Movie";
import { MovieList } from "./MovieList";

const Movies = () => (
  <Switch>
    <Route exact path="/movies" component={MovieList} />
    <Route path="/movies/:id/:tab?" component={Movie} />
  </Switch>
);

export default Movies;
