import { ApolloCache, DocumentNode } from "@apollo/client";
import { isField, resultKeyNameFromField, getQueryDefinition, getOperationName } from "@apollo/client/utilities";
import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";

import { createClient } from "./createClient";

const { client } = createClient();

const getQueryNames = (queries: DocumentNode[]) => {
  return queries.map(q => getOperationName(q) ?? '').filter(q => q !== '');
};

export const getClient = () => client;

export const useFindGalleries = (filter: ListFilterModel) =>
  GQL.useFindGalleriesQuery({
    variables: {
      filter: filter.makeFindFilter(),
      gallery_filter: filter.makeGalleryFilter(),
    },
  });

export const useFindScenes = (filter: ListFilterModel) =>
  GQL.useFindScenesQuery({
    variables: {
      filter: filter.makeFindFilter(),
      scene_filter: filter.makeSceneFilter(),
    },
  });

export const queryFindScenes = (filter: ListFilterModel) =>
  client.query<GQL.FindScenesQuery>({
    query: GQL.FindScenesDocument,
    variables: {
      filter: filter.makeFindFilter(),
      scene_filter: filter.makeSceneFilter(),
    },
  });

export const useFindSceneMarkers = (filter: ListFilterModel) =>
  GQL.useFindSceneMarkersQuery({
    variables: {
      filter: filter.makeFindFilter(),
      scene_marker_filter: filter.makeSceneMarkerFilter(),
    },
  });

export const queryFindSceneMarkers = (filter: ListFilterModel) =>
  client.query<GQL.FindSceneMarkersQuery>({
    query: GQL.FindSceneMarkersDocument,
    variables: {
      filter: filter.makeFindFilter(),
      scene_marker_filter: filter.makeSceneMarkerFilter(),
    },
  });

export const useFindStudios = (filter: ListFilterModel) =>
  GQL.useFindStudiosQuery({
    variables: {
      filter: filter.makeFindFilter(),
      studio_filter: filter.makeStudioFilter(),
    },
  });

export const useFindMovies = (filter: ListFilterModel) =>
  GQL.useFindMoviesQuery({
    variables: {
      filter: filter.makeFindFilter(),
      movie_filter: filter.makeMovieFilter(),
    },
  });

export const useFindPerformers = (filter: ListFilterModel) =>
  GQL.useFindPerformersQuery({
    variables: {
      filter: filter.makeFindFilter(),
      performer_filter: filter.makePerformerFilter(),
    },
  });

export const useFindTags = (filter: ListFilterModel) =>
  GQL.useFindTagsQuery({
    variables: {
      filter: filter.makeFindFilter(),
      tag_filter: filter.makeTagFilter(),
    },
  });

export const queryFindPerformers = (filter: ListFilterModel) =>
  client.query<GQL.FindPerformersQuery>({
    query: GQL.FindPerformersDocument,
    variables: {
      filter: filter.makeFindFilter(),
      performer_filter: filter.makePerformerFilter(),
    },
  });

export const useFindGallery = (id: string) =>
  GQL.useFindGalleryQuery({ variables: { id } });
export const useFindScene = (id: string) =>
  GQL.useFindSceneQuery({ variables: { id } });
export const useSceneStreams = (id: string) =>
  GQL.useSceneStreamsQuery({ variables: { id } });

export const useFindPerformer = (id: string) => {
  const skip = id === "new";
  return GQL.useFindPerformerQuery({ variables: { id }, skip });
};
export const useFindStudio = (id: string) => {
  const skip = id === "new";
  return GQL.useFindStudioQuery({ variables: { id }, skip });
};
export const useFindMovie = (id: string) => {
  const skip = id === "new";
  return GQL.useFindMovieQuery({ variables: { id }, skip });
};
export const useFindTag = (id: string) => {
  const skip = id === "new";
  return GQL.useFindTagQuery({ variables: { id }, skip });
};

const sceneMarkerMutationImpactedQueries = [
  GQL.refetchFindSceneQuery(),
  GQL.refetchFindScenesQuery(),
  GQL.refetchFindSceneMarkersQuery(),
  GQL.refetchMarkerStringsQuery(),
];

export const useSceneMarkerCreate = () =>
  GQL.useSceneMarkerCreateMutation({ refetchQueries: sceneMarkerMutationImpactedQueries });
export const useSceneMarkerUpdate = () =>
  GQL.useSceneMarkerUpdateMutation({ refetchQueries: sceneMarkerMutationImpactedQueries });
export const useSceneMarkerDestroy = () =>
  GQL.useSceneMarkerDestroyMutation({ refetchQueries: sceneMarkerMutationImpactedQueries });

export const useListPerformerScrapers = () =>
  GQL.useListPerformerScrapersQuery();
export const useScrapePerformerList = (scraperId: string, q: string) =>
  GQL.useScrapePerformerListQuery({
    variables: { scraper_id: scraperId, query: q },
    skip: q === "",
  });
export const useScrapePerformer = (
  scraperId: string,
  scrapedPerformer: GQL.ScrapedPerformerInput
) =>
  GQL.useScrapePerformerQuery({
    variables: { scraper_id: scraperId, scraped_performer: scrapedPerformer },
  });

export const useListSceneScrapers = () => GQL.useListSceneScrapersQuery();

export const useListMovieScrapers = () => GQL.useListMovieScrapersQuery();

export const useScrapeFreeonesPerformers = (q: string) =>
  GQL.useScrapeFreeonesPerformersQuery({ variables: { q } });

export const usePlugins = () => GQL.usePluginsQuery();
export const usePluginTasks = () => GQL.usePluginTasksQuery();

export const useMarkerStrings = () => GQL.useMarkerStringsQuery();
export const useAllTags = () => GQL.useAllTagsQuery();
export const useAllTagsForFilter = () => GQL.useAllTagsForFilterQuery();
export const useAllPerformersForFilter = () =>
  GQL.useAllPerformersForFilterQuery();
export const useAllStudiosForFilter = () => GQL.useAllStudiosForFilterQuery();
export const useAllMoviesForFilter = () => GQL.useAllMoviesForFilterQuery();
export const useValidGalleriesForScene = (sceneId: string) =>
  GQL.useValidGalleriesForSceneQuery({
    variables: { scene_id: sceneId },
  });
export const useStats = () => GQL.useStatsQuery();
export const useVersion = () => GQL.useVersionQuery();
export const useLatestVersion = () =>
  GQL.useLatestVersionQuery({
    notifyOnNetworkStatusChange: true,
    errorPolicy: "ignore",
  });

export const useConfiguration = () => GQL.useConfigurationQuery();
export const useDirectory = (path?: string) =>
  GQL.useDirectoryQuery({ variables: { path } });

const performerMutationImpactedQueries = [
  GQL.refetchFindPerformersQuery(),
  GQL.refetchFindSceneQuery(),
  GQL.refetchFindScenesQuery(),
  GQL.refetchFindSceneMarkersQuery(),
  GQL.refetchAllPerformersForFilterQuery(),
];

export const usePerformerCreate = () =>
  GQL.usePerformerCreateMutation({
    refetchQueries: performerMutationImpactedQueries,
  });
export const usePerformerUpdate = () =>
  GQL.usePerformerUpdateMutation({
    refetchQueries: performerMutationImpactedQueries,
  });
export const usePerformerDestroy = () =>
  GQL.usePerformerDestroyMutation({
    refetchQueries: performerMutationImpactedQueries,
  });

const sceneMutationImpactedQueries = [
  GQL.FindPerformerDocument,
  GQL.FindPerformersDocument,
  GQL.FindScenesDocument,
  GQL.FindSceneMarkersDocument,
  GQL.FindStudioDocument,
  GQL.FindStudiosDocument,
  GQL.FindMovieDocument,
  GQL.FindMoviesDocument,
  GQL.FindTagDocument,
  GQL.FindTagsDocument,
  GQL.AllTagsDocument,
];

const deleteCache = (queries: DocumentNode[]) => {
  const names = queries.map(q => {
      const field = getQueryDefinition(q).selectionSet.selections[0];
      return (isField(field) && resultKeyNameFromField(field)) ?? "Unknown";
  }).reduce((fields, name) => ({ ...fields, [name ?? ""]: (_, { DELETE }) => DELETE }), {});


  return (cache: ApolloCache<any>) => (
    cache.modify({
      id: "ROOT_QUERY",
      fields
    })
  );
}

export const useSceneUpdate = (input: GQL.SceneUpdateInput) =>
  GQL.useSceneUpdateMutation({
    variables: input,
    update: (cache) => {
      cache.modify({
        id: "ROOT_QUERY",
        fields: {
          findPerformer(_items, { DELETE }) {
            return DELETE;
          },
          findPerformers(_items, { DELETE }) {
            return DELETE;
          },
          findScenes(_items, { DELETE }) {
            return DELETE;
          }
        }
      })
    }
  });

export const useBulkSceneUpdate = (input: GQL.BulkSceneUpdateInput) =>
  GQL.useBulkSceneUpdateMutation({
    variables: input,
    refetchQueries: getQueryNames(sceneMutationImpactedQueries),
  });

export const useScenesUpdate = (input: GQL.SceneUpdateInput[]) =>
  GQL.useScenesUpdateMutation({ variables: { input } });

export const useSceneIncrementO = (id: string) =>
  GQL.useSceneIncrementOMutation({
    variables: { id },
  });

export const useSceneDecrementO = (id: string) =>
  GQL.useSceneDecrementOMutation({
    variables: { id },
  });

export const useSceneResetO = (id: string) =>
  GQL.useSceneResetOMutation({
    variables: { id },
  });

export const useSceneDestroy = (input: GQL.SceneDestroyInput) =>
  GQL.useSceneDestroyMutation({
    variables: input,
    refetchQueries: getQueryNames(sceneMutationImpactedQueries),
  });

export const useScenesDestroy = (input: GQL.ScenesDestroyInput) =>
  GQL.useScenesDestroyMutation({
    variables: input,
    refetchQueries: getQueryNames(sceneMutationImpactedQueries),
  });

export const useSceneGenerateScreenshot = () =>
  GQL.useSceneGenerateScreenshotMutation({
    refetchQueries: [GQL.refetchFindScenesQuery()],
  });

export const studioMutationImpactedQueries = [
  GQL.refetchFindStudiosQuery(),
  GQL.refetchFindSceneQuery(),
  GQL.refetchFindScenesQuery(),
  GQL.refetchAllStudiosForFilterQuery(),
];

export const useStudioCreate = (input: GQL.StudioCreateInput) =>
  GQL.useStudioCreateMutation({
    variables: input,
    refetchQueries: studioMutationImpactedQueries,
  });

export const useStudioUpdate = (input: GQL.StudioUpdateInput) =>
  GQL.useStudioUpdateMutation({
    variables: input,
    refetchQueries: studioMutationImpactedQueries,
  });

export const useStudioDestroy = (input: GQL.StudioDestroyInput) =>
  GQL.useStudioDestroyMutation({
    variables: input,
    refetchQueries: studioMutationImpactedQueries,
  });

export const movieMutationImpactedQueries = [
  GQL.refetchFindSceneQuery(),
  GQL.refetchFindScenesQuery(),
  GQL.refetchFindMoviesQuery(),
  GQL.refetchAllMoviesForFilterQuery(),
];

export const useMovieCreate = (input: GQL.MovieCreateInput) =>
  GQL.useMovieCreateMutation({
    variables: input,
    refetchQueries: movieMutationImpactedQueries,
  });

export const useMovieUpdate = (input: GQL.MovieUpdateInput) =>
  GQL.useMovieUpdateMutation({
    variables: input,
    refetchQueries: movieMutationImpactedQueries,
  });

export const useMovieDestroy = (input: GQL.MovieDestroyInput) =>
  GQL.useMovieDestroyMutation({
    variables: input,
    refetchQueries: movieMutationImpactedQueries,
  });

export const tagMutationImpactedQueries = [
  GQL.refetchFindSceneQuery(),
  GQL.refetchFindScenesQuery(),
  GQL.refetchFindSceneMarkersQuery(),
  GQL.refetchAllTagsQuery(),
  GQL.refetchAllTagsForFilterQuery(),
  GQL.refetchFindTagsQuery(),
];

export const useTagCreate = (input: GQL.TagCreateInput) =>
  GQL.useTagCreateMutation({
    variables: input,
    refetchQueries: tagMutationImpactedQueries,
  });
export const useTagUpdate = (input: GQL.TagUpdateInput) =>
  GQL.useTagUpdateMutation({
    variables: input,
    refetchQueries: tagMutationImpactedQueries,
  });
export const useTagDestroy = (input: GQL.TagDestroyInput) =>
  GQL.useTagDestroyMutation({
    variables: input,
    refetchQueries: tagMutationImpactedQueries,
  });

export const useConfigureGeneral = (input: GQL.ConfigGeneralInput) =>
  GQL.useConfigureGeneralMutation({
    variables: { input },
    refetchQueries: [GQL.refetchConfigurationQuery()],
  });

export const useConfigureInterface = (input: GQL.ConfigInterfaceInput) =>
  GQL.useConfigureInterfaceMutation({
    variables: { input },
    refetchQueries: [GQL.refetchConfigurationQuery()],
  });

export const useMetadataUpdate = () => GQL.useMetadataUpdateSubscription();

export const useLoggingSubscribe = () => GQL.useLoggingSubscribeSubscription();

export const useLogs = () =>
  GQL.useLogsQuery({
    fetchPolicy: "no-cache",
  });

export const useJobStatus = () =>
  GQL.useJobStatusQuery({
    fetchPolicy: "no-cache",
  });

export const mutateStopJob = () =>
  client.mutate<GQL.StopJobMutation>({
    mutation: GQL.StopJobDocument,
  });

export const queryScrapeFreeones = (performerName: string) =>
  client.query<GQL.ScrapeFreeonesQuery>({
    query: GQL.ScrapeFreeonesDocument,
    variables: {
      performer_name: performerName,
    },
    fetchPolicy: "network-only",
  });

export const queryScrapePerformer = (
  scraperId: string,
  scrapedPerformer: GQL.ScrapedPerformerInput
) =>
  client.query<GQL.ScrapePerformerQuery>({
    query: GQL.ScrapePerformerDocument,
    variables: {
      scraper_id: scraperId,
      scraped_performer: scrapedPerformer,
    },
    fetchPolicy: "network-only",
  });

export const queryScrapePerformerURL = (url: string) =>
  client.query<GQL.ScrapePerformerUrlQuery>({
    query: GQL.ScrapePerformerUrlDocument,
    variables: {
      url,
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeSceneURL = (url: string) =>
  client.query<GQL.ScrapeSceneUrlQuery>({
    query: GQL.ScrapeSceneUrlDocument,
    variables: {
      url,
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeMovieURL = (url: string) =>
  client.query<GQL.ScrapeMovieUrlQuery>({
    query: GQL.ScrapeMovieUrlDocument,
    variables: {
      url,
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeScene = (
  scraperId: string,
  scene: GQL.SceneUpdateInput
) =>
  client.query<GQL.ScrapeSceneQuery>({
    query: GQL.ScrapeSceneDocument,
    variables: {
      scraper_id: scraperId,
      scene,
    },
    fetchPolicy: "network-only",
  });

export const mutateReloadScrapers = () =>
  client.mutate<GQL.ReloadScrapersMutation>({
    mutation: GQL.ReloadScrapersDocument,
  });

const reloadPluginsMutationImpactedQueries = [
  GQL.refetchPluginsQuery(),
  GQL.refetchPluginTasksQuery(),
];

export const mutateReloadPlugins = () =>
  client.mutate<GQL.ReloadPluginsMutation>({
    mutation: GQL.ReloadPluginsDocument,
    refetchQueries: reloadPluginsMutationImpactedQueries,
  });

export const mutateRunPluginTask = (
  pluginId: string,
  taskName: string,
  args?: GQL.PluginArgInput[]
) =>
  client.mutate<GQL.RunPluginTaskMutation>({
    mutation: GQL.RunPluginTaskDocument,
    variables: { plugin_id: pluginId, task_name: taskName, args },
  });

export const mutateMetadataScan = (input: GQL.ScanMetadataInput) =>
  client.mutate<GQL.MetadataScanMutation>({
    mutation: GQL.MetadataScanDocument,
    variables: { input },
  });

export const mutateMetadataAutoTag = (input: GQL.AutoTagMetadataInput) =>
  client.mutate<GQL.MetadataAutoTagMutation>({
    mutation: GQL.MetadataAutoTagDocument,
    variables: { input },
  });

export const mutateMetadataGenerate = (input: GQL.GenerateMetadataInput) =>
  client.mutate<GQL.MetadataGenerateMutation>({
    mutation: GQL.MetadataGenerateDocument,
    variables: { input },
  });

export const mutateMetadataClean = () =>
  client.mutate<GQL.MetadataCleanMutation>({
    mutation: GQL.MetadataCleanDocument,
  });

export const mutateMigrateHashNaming = () =>
  client.mutate<GQL.MigrateHashNamingMutation>({
    mutation: GQL.MigrateHashNamingDocument,
  });

export const mutateMetadataExport = () =>
  client.mutate<GQL.MetadataExportMutation>({
    mutation: GQL.MetadataExportDocument,
  });

export const mutateMetadataImport = () =>
  client.mutate<GQL.MetadataImportMutation>({
    mutation: GQL.MetadataImportDocument,
  });

export const querySceneByPathRegex = (filter: GQL.FindFilterType) =>
  client.query<GQL.FindScenesByPathRegexQuery>({
    query: GQL.FindScenesByPathRegexDocument,
    variables: { filter },
  });

export const queryParseSceneFilenames = (
  filter: GQL.FindFilterType,
  config: GQL.SceneParserInput
) =>
  client.query<GQL.ParseSceneFilenamesQuery>({
    query: GQL.ParseSceneFilenamesDocument,
    variables: { filter, config },
    fetchPolicy: "network-only",
  });

export const stringGenderMap = new Map<string, GQL.GenderEnum>([
  ["Male", GQL.GenderEnum.Male],
  ["Female", GQL.GenderEnum.Female],
  ["Transgender Male", GQL.GenderEnum.TransgenderMale],
  ["Transgender Female", GQL.GenderEnum.TransgenderFemale],
  ["Intersex", GQL.GenderEnum.Intersex],
  ["Non-Binary", GQL.GenderEnum.NonBinary],
]);

export const genderToString = (value?: GQL.GenderEnum) => {
  if (!value) {
    return undefined;
  }

  const foundEntry = Array.from(stringGenderMap.entries()).find((e) => {
    return e[1] === value;
  });

  if (foundEntry) {
    return foundEntry[0];
  }
};

export const stringToGender = (value?: string, caseInsensitive?: boolean) => {
  if (!value) {
    return undefined;
  }

  const ret = stringGenderMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }

  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringGenderMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const getGenderStrings = () => Array.from(stringGenderMap.keys());
