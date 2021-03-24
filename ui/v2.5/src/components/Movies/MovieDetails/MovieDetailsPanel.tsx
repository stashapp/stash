import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  DurationUtils,
  TextUtils,
} from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { TextField, URLField } from "src/utils/field";

interface IMovieDetailsPanel {
  movie: Partial<GQL.MovieDataFragment>;
}

export const MovieDetailsPanel: React.FC<IMovieDetailsPanel> = ({movie}) => {
  // Network state
  const intl = useIntl();

  function maybeRenderAliases() {
    if (movie.aliases) {
      return (
        <div>
          <span className="alias-head">Also known as </span>
          <span className="alias">{movie.aliases}</span>
        </div>
      );
    }
  }

  function renderRatingField() {
    if (!movie.rating) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">Rating</dt>
        <dd className="col-9 col-xl-10">
          <RatingStars
            value={movie.rating}
            disabled={true}
          />
        </dd>
      </dl>
    );
  }

  // TODO: CSS class
  return (
    <div className="movie-details">
      <div>
        <h2>{movie.name}</h2>
      </div>

      {maybeRenderAliases()}

      <div>
        <TextField
          name="Duration"
          value={movie.duration ? DurationUtils.secondsToString(movie.duration) : ""}
        />
        <TextField
          name="Date"
          value={movie.date ? TextUtils.formatDate(intl, movie.date) : ""}
        />
        <URLField
          name="Studio"
          value={movie.studio?.name}
          url={`/studios/${movie.studio?.id}`}
        />
        <TextField
          name="Director"
          value={movie.director}
        />

        {renderRatingField()}

        <URLField
          name="URL"
          value={movie.url}
          url={TextUtils.sanitiseURL(movie.url ?? "")}
        />

        <TextField
          name="Synopsis"
          value={movie.synopsis}
        />
      </div>
    </div>
  );
};
