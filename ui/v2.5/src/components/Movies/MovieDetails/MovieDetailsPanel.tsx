import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { DurationUtils, TextUtils } from "src/utils";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { TextField, URLField } from "src/utils/field";

interface IMovieDetailsPanel {
  movie: GQL.MovieDataFragment;
}

export const MovieDetailsPanel: React.FC<IMovieDetailsPanel> = ({ movie }) => {
  // Network state
  const intl = useIntl();

  function maybeRenderAliases() {
    if (movie.aliases) {
      return (
        <div>
          <span className="alias-head">
            {intl.formatMessage({ id: "also_known_as" })}{" "}
          </span>
          <span className="alias">{movie.aliases}</span>
        </div>
      );
    }
  }

  function renderRatingField() {
    if (!movie.rating100) {
      return;
    }

    return (
      <>
        <dt>{intl.formatMessage({ id: "rating" })}</dt>
        <dd>
          <RatingSystem value={movie.rating100} disabled />
        </dd>
      </>
    );
  }

  // TODO: CSS class
  return (
    <div className="movie-details">
      <div>
        <h2>{movie.name}</h2>
      </div>

      {maybeRenderAliases()}

      <dl className="details-list">
        <TextField
          id="duration"
          value={
            movie.duration ? DurationUtils.secondsToString(movie.duration) : ""
          }
        />
        <TextField
          id="date"
          value={movie.date ? TextUtils.formatDate(intl, movie.date) : ""}
        />
        <URLField
          id="studio"
          value={movie.studio?.name}
          url={`/studios/${movie.studio?.id}`}
        />
        <TextField id="director" value={movie.director} />

        {renderRatingField()}

        <URLField
          id="url"
          value={movie.url}
          url={TextUtils.sanitiseURL(movie.url ?? "")}
        />

        <TextField id="synopsis" value={movie.synopsis} />
      </dl>
    </div>
  );
};
