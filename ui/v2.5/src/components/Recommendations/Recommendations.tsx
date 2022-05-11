import * as GQL from "src/core/generated-graphql";
import { defineMessages, useIntl } from "react-intl";
import React from "react";
import {
  useFindScenes,
  useFindMovies,
  useFindStudios,
  useFindGalleries,
  useFindPerformers,
} from "src/core/StashService";
import { SceneRecommendationRow } from "src/components/Scenes/SceneRecommendationRow";
import { StudioRecommendationRow } from "src/components/Studios/StudioRecommendationRow";
import { MovieRecommendationRow } from "src/components/Movies/MovieRecommendationRow";
import { PerformerRecommendationRow } from "src/components/Performers/PerformerRecommendationRow";
import { GalleryRecommendationRow } from "src/components/Galleries/GalleryRecommendationRow";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import { LoadingIndicator } from "src/components/Shared";

const Recommendations: React.FC = () => {
  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  const isTouch = isTouchEnabled();

  const intl = useIntl();
  const itemsPerPage = 25;
  const scenefilter = new ListFilterModel(GQL.FilterMode.Scenes);
  scenefilter.sortBy = "date";
  scenefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  scenefilter.itemsPerPage = itemsPerPage;
  const sceneResult = useFindScenes(scenefilter);
  const hasScenes = !!sceneResult?.data?.findScenes?.count;

  const studiofilter = new ListFilterModel(GQL.FilterMode.Studios);
  studiofilter.sortBy = "scenes_count";
  studiofilter.sortDirection = GQL.SortDirectionEnum.Desc;
  studiofilter.itemsPerPage = itemsPerPage;
  const studioResult = useFindStudios(studiofilter);
  const hasStudios = !!studioResult?.data?.findStudios?.count;

  const moviefilter = new ListFilterModel(GQL.FilterMode.Movies);
  moviefilter.sortBy = "date";
  moviefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  moviefilter.itemsPerPage = itemsPerPage;
  const movieResult = useFindMovies(moviefilter);
  const hasMovies = !!movieResult?.data?.findMovies?.count;

  const performerfilter = new ListFilterModel(GQL.FilterMode.Performers);
  performerfilter.sortBy = "created_at";
  performerfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  performerfilter.itemsPerPage = itemsPerPage;
  const performerResult = useFindPerformers(performerfilter);
  const hasPerformers = !!performerResult?.data?.findPerformers?.count;

  const galleryfilter = new ListFilterModel(GQL.FilterMode.Galleries);
  galleryfilter.sortBy = "date";
  galleryfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  galleryfilter.itemsPerPage = itemsPerPage;
  const galleryResult = useFindGalleries(galleryfilter);
  const hasGalleries = !!galleryResult?.data?.findGalleries?.count;

  const messages = defineMessages({
    emptyServer: {
      id: "empty_server",
      defaultMessage:
        "Add some scenes to your server to view recommendations on this page.",
    },
    latestScenes: {
      id: "latest_scenes",
      defaultMessage: "Latest Scenes",
    },
    mostActiveStudios: {
      id: "most_active_studios",
      defaultMessage: "Most Active Studios",
    },
    latestMovies: {
      id: "latest_movies",
      defaultMessage: "Latest Movies",
    },
    latestPerformers: {
      id: "latest_performers",
      defaultMessage: "Latest Performers",
    },
    latestGalleries: {
      id: "latest_galleries",
      defaultMessage: "Latest Galleries",
    },
  });

  if (
    sceneResult.loading ||
    studioResult.loading ||
    movieResult.loading ||
    performerResult.loading ||
    galleryResult.loading
  ) {
    return <LoadingIndicator />;
  } else {
    return (
      <div className="recommendations-container">
        {!hasScenes &&
        !hasStudios &&
        !hasMovies &&
        !hasPerformers &&
        !hasGalleries ? (
          <div className="no-recommendations">
            {intl.formatMessage(messages.emptyServer)}
          </div>
        ) : (
          <div>
            {hasScenes && (
              <SceneRecommendationRow
                isTouch={isTouch}
                result={sceneResult}
                queue={SceneQueue.fromListFilterModel(scenefilter)}
                header={intl.formatMessage(messages.latestScenes)}
              />
            )}

            {hasStudios && (
              <StudioRecommendationRow
                isTouch={isTouch}
                result={studioResult}
                header={intl.formatMessage(messages.mostActiveStudios)}
              />
            )}

            {hasMovies && (
              <MovieRecommendationRow
                isTouch={isTouch}
                result={movieResult}
                header={intl.formatMessage(messages.latestMovies)}
              />
            )}

            {hasPerformers && (
              <PerformerRecommendationRow
                isTouch={isTouch}
                result={performerResult}
                header={intl.formatMessage(messages.latestPerformers)}
              />
            )}

            {hasGalleries && (
              <GalleryRecommendationRow
                isTouch={isTouch}
                result={galleryResult}
                header={intl.formatMessage(messages.latestGalleries)}
              />
            )}
          </div>
        )}
      </div>
    );
  }
};

export default Recommendations;
