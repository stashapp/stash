import React from "react";
import { Route, Switch } from "react-router-dom";
import Movie from "./MovieDetails/Movie";
import MovieCreate from "./MovieDetails/MovieCreate";
import { MovieList } from "./MovieList";

const Movies = () => (
  <Switch>
    <Route exact path="/movies" component={MovieList} />
    <Route exact path="/movies/new" component={MovieCreate} />
    <Route path="/movies/:id/:tab?" component={Movie} />
  </Switch>
);

export default Movies;
