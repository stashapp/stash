import React, { useContext, useMemo } from "react";
import { useIntl } from "react-intl";
import { FrontPageContent, ICustomFilter } from "src/core/config";
import * as GQL from "src/core/generated-graphql";
import { useFindSavedFilter } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { ListFilterModel } from "src/models/list-filter/filter";
import { GalleryRecommendationRow } from "../Galleries/GalleryRecommendationRow";
import { ImageRecommendationRow } from "../Images/ImageRecommendationRow";
import { GroupRecommendationRow } from "../Movies/MovieRecommendationRow";
import { PerformerRecommendationRow } from "../Performers/PerformerRecommendationRow";
import { SceneRecommendationRow } from "../Scenes/SceneRecommendationRow";
import { StudioRecommendationRow } from "../Studios/StudioRecommendationRow";
import { TagRecommendationRow } from "../Tags/TagRecommendationRow";

interface IFilter {
  mode: GQL.FilterMode;
  filter: ListFilterModel;
  header: string;
}

const RecommendationRow: React.FC<IFilter> = ({ mode, filter, header }) => {
  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  const isTouch = isTouchEnabled();

  switch (mode) {
    case GQL.FilterMode.Scenes:
      return (
        <SceneRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Studios:
      return (
        <StudioRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Movies:
      return (
        <GroupRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Performers:
      return (
        <PerformerRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Galleries:
      return (
        <GalleryRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Images:
      return (
        <ImageRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    case GQL.FilterMode.Tags:
      return (
        <TagRecommendationRow
          isTouch={isTouch}
          filter={filter}
          header={header}
        />
      );
    default:
      return <></>;
  }
};

interface ISavedFilterResults {
  savedFilterID: string;
}

const SavedFilterResults: React.FC<ISavedFilterResults> = ({
  savedFilterID,
}) => {
  const { configuration: config } = useContext(ConfigurationContext);
  const { loading, data } = useFindSavedFilter(savedFilterID.toString());

  const filter = useMemo(() => {
    if (!data?.findSavedFilter) return;

    const { mode } = data.findSavedFilter;

    const ret = new ListFilterModel(mode, config);
    ret.currentPage = 1;
    ret.configureFromSavedFilter(data.findSavedFilter);
    ret.randomSeed = -1;
    return ret;
  }, [data?.findSavedFilter, config]);

  if (loading || !data?.findSavedFilter || !filter) {
    return <></>;
  }

  const { name, mode } = data.findSavedFilter;

  return <RecommendationRow mode={mode} filter={filter} header={name} />;
};

interface ICustomFilterProps {
  customFilter: ICustomFilter;
}

const CustomFilterResults: React.FC<ICustomFilterProps> = ({
  customFilter,
}) => {
  const { configuration: config } = useContext(ConfigurationContext);
  const intl = useIntl();

  const filter = useMemo(() => {
    const itemsPerPage = 25;
    const ret = new ListFilterModel(customFilter.mode, config);
    ret.sortBy = customFilter.sortBy;
    ret.sortDirection = customFilter.direction;
    ret.itemsPerPage = itemsPerPage;
    ret.currentPage = 1;
    ret.randomSeed = -1;
    return ret;
  }, [customFilter, config]);

  const header = customFilter.message
    ? intl.formatMessage(
        { id: customFilter.message.id },
        customFilter.message.values
      )
    : customFilter.title ?? "";

  return (
    <RecommendationRow
      mode={customFilter.mode}
      filter={filter}
      header={header}
    />
  );
};

interface IProps {
  content: FrontPageContent;
}

export const Control: React.FC<IProps> = ({ content }) => {
  switch (content.__typename) {
    case "SavedFilter":
      if (!content.savedFilterId) {
        return <div>Error: missing savedFilterId</div>;
      }

      return (
        <SavedFilterResults savedFilterID={content.savedFilterId.toString()} />
      );
    case "CustomFilter":
      return <CustomFilterResults customFilter={content} />;
    default:
      return <></>;
  }
};
