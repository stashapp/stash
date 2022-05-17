import React from "react";
import * as GQL from "src/core/generated-graphql";
import { defineMessages, useIntl } from "react-intl";
import { useFindRecommendationFilters } from "src/core/StashService";
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

  const userRecommendations = useFindRecommendationFilters();

  const messages = defineMessages({
    emptyServer: {
      id: "empty_server",
      defaultMessage:
        "Add some scenes to your server to view recommendations on this page.",
    },
    recentlyAddedStudios: {
      id: "recently_added_studios",
      defaultMessage: "Recently Added Studios",
    },
    recentlyAddedPerformers: {
      id: "recently_added_performers",
      defaultMessage: "Recently Added Performers",
    },
    recentlyReleasedGalleries: {
      id: "recently_released_galleries",
      defaultMessage: "Recently Released Galleries",
    },
    recentlyReleasedMovies: {
      id: "recently_released_movies",
      defaultMessage: "Recently Released Movies",
    },
    recentlyReleasedScenes: {
      id: "recently_released_scenes",
      defaultMessage: "Recently Released Scenes",
    },
    viewAll: {
      id: "view_all",
      defaultMessage: "View All",
    },
  });

  const filters: ListFilterModel[] = [];
  const labels: string[] = [];
  let inProgress = true;
  if (!userRecommendations.loading) {
    if (!!userRecommendations?.data?.findRecommendedFilters?.length) {
      userRecommendations.data.findRecommendedFilters.forEach(function (f) {
        const newFilter = new ListFilterModel(f.mode);
        newFilter.currentPage = 1;
        newFilter.configureFromQueryParameters(JSON.parse(f.filter));
        newFilter.randomSeed = -1;
        filters.push(newFilter);
        labels.push(f.name);
      });
      inProgress = false;
    } else if (!userRecommendations.loading) {
      const itemsPerPage = 25;
      const scenefilter = new ListFilterModel(GQL.FilterMode.Scenes);
      scenefilter.sortBy = "date";
      scenefilter.sortDirection = GQL.SortDirectionEnum.Desc;
      scenefilter.itemsPerPage = itemsPerPage;
      filters.push(scenefilter);
      labels.push(intl.formatMessage(messages.recentlyReleasedScenes));
      const studiofilter = new ListFilterModel(GQL.FilterMode.Studios);
      studiofilter.sortBy = "created_at";
      studiofilter.sortDirection = GQL.SortDirectionEnum.Desc;
      studiofilter.itemsPerPage = itemsPerPage;
      filters.push(studiofilter);
      labels.push(intl.formatMessage(messages.recentlyAddedStudios));
      const moviefilter = new ListFilterModel(GQL.FilterMode.Movies);
      moviefilter.sortBy = "date";
      moviefilter.sortDirection = GQL.SortDirectionEnum.Desc;
      moviefilter.itemsPerPage = itemsPerPage;
      filters.push(moviefilter);
      labels.push(intl.formatMessage(messages.recentlyReleasedMovies));
      const performerfilter = new ListFilterModel(GQL.FilterMode.Performers);
      performerfilter.sortBy = "created_at";
      performerfilter.sortDirection = GQL.SortDirectionEnum.Desc;
      performerfilter.itemsPerPage = itemsPerPage;
      filters.push(performerfilter);
      labels.push(intl.formatMessage(messages.recentlyAddedPerformers));
      const galleryfilter = new ListFilterModel(GQL.FilterMode.Galleries);
      galleryfilter.sortBy = "date";
      galleryfilter.sortDirection = GQL.SortDirectionEnum.Desc;
      galleryfilter.itemsPerPage = itemsPerPage;
      filters.push(galleryfilter);
      labels.push(intl.formatMessage(messages.recentlyReleasedGalleries));
      inProgress = false;
    }
  }

  if (inProgress) {
    return <LoadingIndicator />;
  }

  return (
    <div className="recommendations-container">
      <div>
        {filters.map((filter, index) => {
          if (filter.mode == GQL.FilterMode.Scenes) {
            return (
              <SceneRecommendationRow
                isTouch={isTouch}
                filter={filter}
                queue={SceneQueue.fromListFilterModel(filter)}
                header={labels[index]!}
                linkText={intl.formatMessage(messages.viewAll)}
                index={index}
              />
            );
          } else if (filter.mode == GQL.FilterMode.Studios) {
            return (
              <StudioRecommendationRow
                isTouch={isTouch}
                filter={filter}
                header={labels[index]!}
                linkText={intl.formatMessage(messages.viewAll)}
                index={index}
              />
            );
          } else if (filter.mode == GQL.FilterMode.Movies) {
            return (
              <MovieRecommendationRow
                isTouch={isTouch}
                filter={filter}
                header={labels[index]!}
                linkText={intl.formatMessage(messages.viewAll)}
                index={index}
              />
            );
          } else if (filter.mode == GQL.FilterMode.Performers) {
            return (
              <PerformerRecommendationRow
                isTouch={isTouch}
                filter={filter}
                header={labels[index]!}
                linkText={intl.formatMessage(messages.viewAll)}
                index={index}
              />
            );
          } else if (filter.mode == GQL.FilterMode.Galleries) {
            return (
              <GalleryRecommendationRow
                isTouch={isTouch}
                filter={filter}
                header={labels[index]!}
                linkText={intl.formatMessage(messages.viewAll)}
                index={index}
              />
            );
          }
        })}
      </div>
    </div>
  );
};

export default Recommendations;
