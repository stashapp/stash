import React, { useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { defineMessages, FormattedMessage, useIntl } from "react-intl";
import {
  useConfigureFrontPage,
  useFindFrontPageFiltersQuery,
} from "src/core/StashService";
import { SceneRecommendationRow } from "src/components/Scenes/SceneRecommendationRow";
import { StudioRecommendationRow } from "src/components/Studios/StudioRecommendationRow";
import { MovieRecommendationRow } from "src/components/Movies/MovieRecommendationRow";
import { PerformerRecommendationRow } from "src/components/Performers/PerformerRecommendationRow";
import { GalleryRecommendationRow } from "src/components/Galleries/GalleryRecommendationRow";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import { LoadingIndicator } from "src/components/Shared";
import { Button } from "react-bootstrap";
import { FrontPageConfig } from "./FrontPageConfig";
import { useToast } from "src/hooks";

interface IRowSpec {
  filter: ListFilterModel;
  header: string;
}

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

const Recommendations: React.FC = () => {
  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  const isTouch = isTouchEnabled();

  const intl = useIntl();
  const Toast = useToast();

  const [isEditing, setIsEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const { data, loading } = useFindFrontPageFiltersQuery();
  const [updateFrontPageConfig] = useConfigureFrontPage();

  async function onUpdateConfig(newIDs?: string[]) {
    setIsEditing(false);

    if (!newIDs) {
      return;
    }

    setSaving(true);
    try {
      updateFrontPageConfig({
        variables: {
          input: {
            savedFilterIDs: newIDs,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
    setSaving(false);
  }

  const filters = useMemo(() => {
    const ret: IRowSpec[] = [];

    if (!loading && data) {
      const { findFrontPageFilters } = data;
      if (!!findFrontPageFilters.length) {
        data.findFrontPageFilters.forEach(function (f) {
          const newFilter = new ListFilterModel(f.mode);
          newFilter.currentPage = 1;
          newFilter.configureFromQueryParameters(JSON.parse(f.filter));
          newFilter.randomSeed = -1;
          ret.push({
            filter: newFilter,
            header: f.name,
          });
        });
      } else {
        const itemsPerPage = 25;
        const scenefilter = new ListFilterModel(GQL.FilterMode.Scenes);
        scenefilter.sortBy = "date";
        scenefilter.sortDirection = GQL.SortDirectionEnum.Desc;
        scenefilter.itemsPerPage = itemsPerPage;
        ret.push({
          filter: scenefilter,
          header: intl.formatMessage(messages.recentlyReleasedScenes),
        });
        const studiofilter = new ListFilterModel(GQL.FilterMode.Studios);
        studiofilter.sortBy = "created_at";
        studiofilter.sortDirection = GQL.SortDirectionEnum.Desc;
        studiofilter.itemsPerPage = itemsPerPage;
        ret.push({
          filter: studiofilter,
          header: intl.formatMessage(messages.recentlyAddedStudios),
        });
        const moviefilter = new ListFilterModel(GQL.FilterMode.Movies);
        moviefilter.sortBy = "date";
        moviefilter.sortDirection = GQL.SortDirectionEnum.Desc;
        moviefilter.itemsPerPage = itemsPerPage;
        ret.push({
          filter: moviefilter,
          header: intl.formatMessage(messages.recentlyReleasedMovies),
        });
        const performerfilter = new ListFilterModel(GQL.FilterMode.Performers);
        performerfilter.sortBy = "created_at";
        performerfilter.sortDirection = GQL.SortDirectionEnum.Desc;
        performerfilter.itemsPerPage = itemsPerPage;
        ret.push({
          filter: performerfilter,
          header: intl.formatMessage(messages.recentlyAddedPerformers),
        });
        const galleryfilter = new ListFilterModel(GQL.FilterMode.Galleries);
        galleryfilter.sortBy = "date";
        galleryfilter.sortDirection = GQL.SortDirectionEnum.Desc;
        galleryfilter.itemsPerPage = itemsPerPage;
        ret.push({
          filter: galleryfilter,
          header: intl.formatMessage(messages.recentlyReleasedGalleries),
        });
      }
    }

    return ret;
  }, [loading, data, intl]);

  const rows = useMemo(() => {
    function renderRow(rowSpec: IRowSpec, index: number) {
      if (rowSpec.filter.mode == GQL.FilterMode.Scenes) {
        return (
          <SceneRecommendationRow
            isTouch={isTouch}
            filter={rowSpec.filter}
            queue={SceneQueue.fromListFilterModel(rowSpec.filter)}
            header={rowSpec.header}
            linkText={intl.formatMessage(messages.viewAll)}
            index={index}
          />
        );
      } else if (rowSpec.filter.mode == GQL.FilterMode.Studios) {
        return (
          <StudioRecommendationRow
            isTouch={isTouch}
            filter={rowSpec.filter}
            header={rowSpec.header}
            linkText={intl.formatMessage(messages.viewAll)}
            index={index}
          />
        );
      } else if (rowSpec.filter.mode == GQL.FilterMode.Movies) {
        return (
          <MovieRecommendationRow
            isTouch={isTouch}
            filter={rowSpec.filter}
            header={rowSpec.header}
            linkText={intl.formatMessage(messages.viewAll)}
            index={index}
          />
        );
      } else if (rowSpec.filter.mode == GQL.FilterMode.Performers) {
        return (
          <PerformerRecommendationRow
            isTouch={isTouch}
            filter={rowSpec.filter}
            header={rowSpec.header}
            linkText={intl.formatMessage(messages.viewAll)}
            index={index}
          />
        );
      } else if (rowSpec.filter.mode == GQL.FilterMode.Galleries) {
        return (
          <GalleryRecommendationRow
            isTouch={isTouch}
            filter={rowSpec.filter}
            header={rowSpec.header}
            linkText={intl.formatMessage(messages.viewAll)}
            index={index}
          />
        );
      }
    }

    return filters.map((filter, index) => renderRow(filter, index));
  }, [isTouch, filters, intl]);

  if (loading || saving) {
    return <LoadingIndicator />;
  }

  if (isEditing) {
    return (
      <FrontPageConfig onClose={(filterIDs) => onUpdateConfig(filterIDs)} />
    );
  }

  return (
    <div className="recommendations-container">
      <div>{rows}</div>
      <div className="recommendations-footer">
        <Button onClick={() => setIsEditing(true)}>
          <FormattedMessage id={"actions.customise"} />
        </Button>
      </div>
    </div>
  );
};

export default Recommendations;
