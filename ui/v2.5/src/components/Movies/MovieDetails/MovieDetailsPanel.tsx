import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { DetailItem } from "src/components/Shared/DetailItem";
import { Link } from "react-router-dom";
import { DirectorLink } from "src/components/Shared/Link";
import { TagLink } from "src/components/Shared/TagLink";

interface IMovieDetailsPanel {
  movie: GQL.MovieDataFragment;
  collapsed?: boolean;
  fullWidth?: boolean;
}

export const MovieDetailsPanel: React.FC<IMovieDetailsPanel> = ({
  movie,
  collapsed,
  fullWidth,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    if (!movie.tags.length) {
      return;
    }
    return (
      <ul className="pl-0">
        {(movie.tags ?? []).map((tag) => (
          <TagLink key={tag.id} linkType="movie" tag={tag} />
        ))}
      </ul>
    );
  }

  function maybeRenderExtraDetails() {
    if (!collapsed) {
      return (
        <>
          <DetailItem
            id="synopsis"
            value={movie.synopsis}
            fullWidth={fullWidth}
          />
          <DetailItem
            id="tags"
            value={renderTagsField()}
            fullWidth={fullWidth}
          />
        </>
      );
    }
  }

  return (
    <div className="detail-group">
      <DetailItem
        id="duration"
        value={
          movie.duration ? TextUtils.secondsToTimestamp(movie.duration) : ""
        }
        fullWidth={fullWidth}
      />
      <DetailItem
        id="date"
        value={movie.date ? TextUtils.formatDate(intl, movie.date) : ""}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="studio"
        value={
          movie.studio?.id ? (
            <Link to={`/studios/${movie.studio?.id}`}>
              {movie.studio?.name}
            </Link>
          ) : (
            ""
          )
        }
        fullWidth={fullWidth}
      />

      <DetailItem
        id="director"
        value={
          movie.director ? (
            <DirectorLink director={movie.director} linkType="movie" />
          ) : (
            ""
          )
        }
        fullWidth={fullWidth}
      />
      {maybeRenderExtraDetails()}
    </div>
  );
};

export const CompressedMovieDetailsPanel: React.FC<IMovieDetailsPanel> = ({
  movie,
}) => {
  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="movie-name" onClick={() => scrollToTop()}>
          {movie.name}
        </a>
        {movie?.studio?.name ? (
          <>
            <span className="detail-divider">/</span>
            <span className="movie-studio">{movie?.studio?.name}</span>
          </>
        ) : (
          ""
        )}
      </div>
    </div>
  );
};
