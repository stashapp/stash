import React from "react";
import { Route, Switch } from "react-router-dom";
import { useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import { TITLE_SUFFIX } from "src/components/Shared/constants";
import Movie from "./MovieDetails/Movie";
import MovieCreate from "./MovieDetails/MovieCreate";
import { MovieList } from "./MovieList";

const Movies: React.FC = () => {
  const intl = useIntl();

  const title_template = `${intl.formatMessage({
    id: "movies",
  })} ${TITLE_SUFFIX}`;
  return (
    <>
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />
      <Switch>
        <Route exact path="/movies" component={MovieList} />
        <Route exact path="/movies/new" component={MovieCreate} />
        <Route path="/movies/:id/:tab?" component={Movie} />
      </Switch>
    </>
  );
};

export default Movies;
