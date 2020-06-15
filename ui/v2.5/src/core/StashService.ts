import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";

import { createClient } from "./createClient";

const { client, cache } = createClient();

export const getClient = () => client;

// TODO: Invalidation should happen through apollo client, rather than rewriting cache directly
const invalidateQueries = (queries: string[]) => {
  if (cache) {
    const keyMatchers = queries.map((query) => {
      return new RegExp(`^${query}`);
    });

    // TODO: Hack to invalidate, manipulating private data
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const rootQuery = (cache as any).data.data.ROOT_QUERY;
    Object.keys(rootQuery).forEach((key) => {
      if (
        keyMatchers.some((matcher) => {
          return !!key.match(matcher);
        })
      ) {
        delete rootQuery[key];
      }
    });
  }
};

export const useFindGalleries = (filter: ListFilterModel) =>
  GQL.useFindGalleriesQuery({
    variables: {
      filter: filter.makeFindFilter(),
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

// TODO - scene marker manipulation functions are handled differently
export const sceneMarkerMutationImpactedQueries = [
  "findSceneMarkers",
  "findScenes",
  "markerStrings",
  "sceneMarkerTags",
];

export const useSceneMarkerCreate = () =>
  GQL.useSceneMarkerCreateMutation({ refetchQueries: ["FindScene"] });
export const useSceneMarkerUpdate = () =>
  GQL.useSceneMarkerUpdateMutation({ refetchQueries: ["FindScene"] });
export const useSceneMarkerDestroy = () =>
  GQL.useSceneMarkerDestroyMutation({ refetchQueries: ["FindScene"] });

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

export const useScrapeFreeonesPerformers = (q: string) =>
  GQL.useScrapeFreeonesPerformersQuery({ variables: { q } });
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

export const performerMutationImpactedQueries = [
  "findPerformers",
  "findScenes",
  "findSceneMarkers",
  "allPerformers",
];

export const usePerformerCreate = () =>
  GQL.usePerformerCreateMutation({
    update: () => invalidateQueries(performerMutationImpactedQueries),
  });
export const usePerformerUpdate = () =>
  GQL.usePerformerUpdateMutation({
    update: () => invalidateQueries(performerMutationImpactedQueries),
  });
export const usePerformerDestroy = () =>
  GQL.usePerformerDestroyMutation({
    update: () => invalidateQueries(performerMutationImpactedQueries),
  });

export const sceneMutationImpactedQueries = [
  "findPerformers",
  "findScenes",
  "findSceneMarkers",
  "findStudios",
  "findMovies",
  "allTags",
  // TODO - add "findTags" when it is implemented
];

export const useSceneUpdate = (input: GQL.SceneUpdateInput) =>
  GQL.useSceneUpdateMutation({
    variables: input,
    update: () => invalidateQueries(sceneMutationImpactedQueries),
    refetchQueries: ["AllTagsForFilter"],
  });

// remove findScenes for bulk scene update so that we don't lose
// existing results
export const sceneBulkMutationImpactedQueries = [
  "findPerformers",
  "findSceneMarkers",
  "findStudios",
  "findMovies",
  "allTags",
];

export const useBulkSceneUpdate = (input: GQL.BulkSceneUpdateInput) =>
  GQL.useBulkSceneUpdateMutation({
    variables: input,
    update: () => invalidateQueries(sceneBulkMutationImpactedQueries),
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
    update: () => invalidateQueries(sceneMutationImpactedQueries),
  });

export const useSceneGenerateScreenshot = () =>
  GQL.useSceneGenerateScreenshotMutation({
    update: () => invalidateQueries(["findScenes"]),
  });

export const studioMutationImpactedQueries = [
  "findStudios",
  "findScenes",
  "allStudios",
];

export const useStudioCreate = (input: GQL.StudioCreateInput) =>
  GQL.useStudioCreateMutation({
    variables: input,
    update: () => invalidateQueries(studioMutationImpactedQueries),
  });

export const useStudioUpdate = (input: GQL.StudioUpdateInput) =>
  GQL.useStudioUpdateMutation({
    variables: input,
    update: () => invalidateQueries(studioMutationImpactedQueries),
  });

export const useStudioDestroy = (input: GQL.StudioDestroyInput) =>
  GQL.useStudioDestroyMutation({
    variables: input,
    update: () => invalidateQueries(studioMutationImpactedQueries),
  });

export const movieMutationImpactedQueries = [
  "findMovies",
  "findScenes",
  "allMovies",
];

export const useMovieCreate = (input: GQL.MovieCreateInput) =>
  GQL.useMovieCreateMutation({
    variables: input,
    update: () => invalidateQueries(movieMutationImpactedQueries),
  });

export const useMovieUpdate = (input: GQL.MovieUpdateInput) =>
  GQL.useMovieUpdateMutation({
    variables: input,
    update: () => invalidateQueries(movieMutationImpactedQueries),
  });

export const useMovieDestroy = (input: GQL.MovieDestroyInput) =>
  GQL.useMovieDestroyMutation({
    variables: input,
    update: () => invalidateQueries(movieMutationImpactedQueries),
  });

export const tagMutationImpactedQueries = [
  "findScenes",
  "findSceneMarkers",
  "sceneMarkerTags",
  "allTags",
];

export const useTagCreate = (input: GQL.TagCreateInput) =>
  GQL.useTagCreateMutation({
    variables: input,
    refetchQueries: ["AllTags", "AllTagsForFilter"],
    // update: () => StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
  });
export const useTagUpdate = (input: GQL.TagUpdateInput) =>
  GQL.useTagUpdateMutation({
    variables: input,
    refetchQueries: ["AllTags", "AllTagsForFilter"],
  });
export const useTagDestroy = (input: GQL.TagDestroyInput) =>
  GQL.useTagDestroyMutation({
    variables: input,
    refetchQueries: ["AllTags", "AllTagsForFilter"],
    update: () => invalidateQueries(tagMutationImpactedQueries),
  });

export const useConfigureGeneral = (input: GQL.ConfigGeneralInput) =>
  GQL.useConfigureGeneralMutation({
    variables: { input },
    refetchQueries: ["Configuration"],
  });

export const useConfigureInterface = (input: GQL.ConfigInterfaceInput) =>
  GQL.useConfigureInterfaceMutation({
    variables: { input },
    refetchQueries: ["Configuration"],
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
  });

export const queryScrapePerformerURL = (url: string) =>
  client.query<GQL.ScrapePerformerUrlQuery>({
    query: GQL.ScrapePerformerUrlDocument,
    variables: {
      url,
    },
  });

export const queryScrapeSceneURL = (url: string) =>
  client.query<GQL.ScrapeSceneUrlQuery>({
    query: GQL.ScrapeSceneUrlDocument,
    variables: {
      url,
    },
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
  });

export const mutateReloadScrapers = () =>
  client.mutate<GQL.ReloadScrapersMutation>({
    mutation: GQL.ReloadScrapersDocument,
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
