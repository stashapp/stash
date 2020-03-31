import ApolloClient from "apollo-client";
import { WebSocketLink } from "apollo-link-ws";
import { InMemoryCache, NormalizedCacheObject } from "apollo-cache-inmemory";
import { HttpLink } from "apollo-link-http";
import { split } from "apollo-link";
import { getMainDefinition } from "apollo-utilities";
import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";

export class StashService {
  public static client: ApolloClient<NormalizedCacheObject>;
  private static cache: InMemoryCache;

  public static getPlatformURL(ws?: boolean) {
    const platformUrl = new URL(window.location.origin);

    if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
      platformUrl.port = "9999"; // TODO: Hack. Development expects port 9999

      if (process.env.REACT_APP_HTTPS === "true") {
        platformUrl.protocol = "https:";
      }
    }

    if (ws) {
      platformUrl.protocol = "ws:";
    }

    return platformUrl;
  }

  public static initialize() {
    const platformUrl = StashService.getPlatformURL();
    const wsPlatformUrl = StashService.getPlatformURL(true);

    if (platformUrl.protocol === "https:") {
      wsPlatformUrl.protocol = "wss:";
    }

    const url = `${platformUrl.toString().slice(0, -1)}/graphql`;
    const wsUrl = `${wsPlatformUrl.toString().slice(0, -1)}/graphql`;

    const httpLink = new HttpLink({
      uri: url
    });

    const wsLink = new WebSocketLink({
      uri: wsUrl,
      options: {
        reconnect: true
      }
    });

    const link = split(
      ({ query }) => {
        const definition = getMainDefinition(query);
        return (
          definition.kind === "OperationDefinition" &&
          definition.operation === "subscription"
        );
      },
      wsLink,
      httpLink
    );

    StashService.cache = new InMemoryCache();
    StashService.client = new ApolloClient({
      link,
      cache: StashService.cache
    });

    return StashService.client;
  }

  // TODO: Invalidation should happen through apollo client, rather than rewriting cache directly
  private static invalidateQueries(queries: string[]) {
    if (StashService.cache) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const cache = StashService.cache as any;
      const keyMatchers = queries.map(query => {
        return new RegExp(`^${query}`);
      });

      const rootQuery = cache.data.data.ROOT_QUERY;
      Object.keys(rootQuery).forEach(key => {
        if (
          keyMatchers.some(matcher => {
            return !!key.match(matcher);
          })
        ) {
          delete rootQuery[key];
        }
      });
    }
  }

  public static useFindGalleries(filter: ListFilterModel) {
    return GQL.useFindGalleriesQuery({
      variables: {
        filter: filter.makeFindFilter()
      }
    });
  }

  public static useFindScenes(filter: ListFilterModel) {
    let sceneFilter = {};
    // if (!!filter && filter.criteriaFilterOpen) {
    sceneFilter = filter.makeSceneFilter();
    // }
    // if (filter.customCriteria) {
    //   filter.customCriteria.forEach(criteria => {
    //     scene_filter[criteria.key] = criteria.value;
    //   });
    // }

    return GQL.useFindScenesQuery({
      variables: {
        filter: filter.makeFindFilter(),
        scene_filter: sceneFilter
      }
    });
  }

  public static queryFindScenes(filter: ListFilterModel) {
    let sceneFilter = {};
    sceneFilter = filter.makeSceneFilter();

    return StashService.client.query<GQL.FindScenesQuery>({
      query: GQL.FindScenesDocument,
      variables: {
        filter: filter.makeFindFilter(),
        scene_filter: sceneFilter
      }
    });
  }

  public static useFindSceneMarkers(filter: ListFilterModel) {
    let sceneMarkerFilter = {};
    // if (!!filter && filter.criteriaFilterOpen) {
    sceneMarkerFilter = filter.makeSceneMarkerFilter();
    // }
    // if (filter.customCriteria) {
    //   filter.customCriteria.forEach(criteria => {
    //     scene_filter[criteria.key] = criteria.value;
    //   });
    // }

    return GQL.useFindSceneMarkersQuery({
      variables: {
        filter: filter.makeFindFilter(),
        scene_marker_filter: sceneMarkerFilter
      }
    });
  }

  public static queryFindSceneMarkers(filter: ListFilterModel) {
    let sceneMarkerFilter = {};
    sceneMarkerFilter = filter.makeSceneMarkerFilter();

    return StashService.client.query<GQL.FindSceneMarkersQuery>({
      query: GQL.FindSceneMarkersDocument,
      variables: {
        filter: filter.makeFindFilter(),
        scene_marker_filter: sceneMarkerFilter
      }
    });
  }

  public static useFindStudios(filter: ListFilterModel) {
    return GQL.useFindStudiosQuery({
      variables: {
        filter: filter.makeFindFilter()
      }
    });
  }

  public static useFindMovies(filter: ListFilterModel) {
    return GQL.useFindMoviesQuery({
      variables: {
        filter: filter.makeFindFilter()
      }
    });
  }

  public static useFindPerformers(filter: ListFilterModel) {
    let performerFilter = {};
    // if (!!filter && filter.criteriaFilterOpen) {
    performerFilter = filter.makePerformerFilter();
    // }
    // if (filter.customCriteria) {
    //   filter.customCriteria.forEach(criteria => {
    //     scene_filter[criteria.key] = criteria.value;
    //   });
    // }

    return GQL.useFindPerformersQuery({
      variables: {
        filter: filter.makeFindFilter(),
        performer_filter: performerFilter
      }
    });
  }

  public static queryFindPerformers(filter: ListFilterModel) {
    let performerFilter = {};
    performerFilter = filter.makePerformerFilter();

    return StashService.client.query<GQL.FindPerformersQuery>({
      query: GQL.FindPerformersDocument,
      variables: {
        filter: filter.makeFindFilter(),
        performer_filter: performerFilter
      }
    });
  }

  public static useFindGallery(id: string) {
    return GQL.useFindGalleryQuery({ variables: { id } });
  }
  public static useFindScene(id: string) {
    return GQL.useFindSceneQuery({ variables: { id } });
  }
  public static useFindPerformer(id: string) {
    const skip = id === "new";
    return GQL.useFindPerformerQuery({ variables: { id }, skip });
  }
  public static useFindStudio(id: string) {
    const skip = id === "new";
    return GQL.useFindStudioQuery({ variables: { id }, skip });
  }
  public static useFindMovie(id: string) {
    const skip = id === "new";
    return GQL.useFindMovieQuery({ variables: { id }, skip });
  }

  // TODO - scene marker manipulation functions are handled differently
  private static sceneMarkerMutationImpactedQueries = [
    "findSceneMarkers",
    "findScenes",
    "markerStrings",
    "sceneMarkerTags"
  ];

  public static useSceneMarkerCreate() {
    return GQL.useSceneMarkerCreateMutation({ refetchQueries: ["FindScene"] });
  }
  public static useSceneMarkerUpdate() {
    return GQL.useSceneMarkerUpdateMutation({ refetchQueries: ["FindScene"] });
  }
  public static useSceneMarkerDestroy() {
    return GQL.useSceneMarkerDestroyMutation({ refetchQueries: ["FindScene"] });
  }

  public static useListPerformerScrapers() {
    return GQL.useListPerformerScrapersQuery();
  }
  public static useScrapePerformerList(scraperId: string, q: string) {
    return GQL.useScrapePerformerListQuery({
      variables: { scraper_id: scraperId, query: q },
      skip: q === ""
    });
  }
  public static useScrapePerformer(
    scraperId: string,
    scrapedPerformer: GQL.ScrapedPerformerInput
  ) {
    return GQL.useScrapePerformerQuery({
      variables: { scraper_id: scraperId, scraped_performer: scrapedPerformer }
    });
  }

  public static useListSceneScrapers() {
    return GQL.useListSceneScrapersQuery();
  }

  public static useScrapeFreeonesPerformers(q: string) {
    return GQL.useScrapeFreeonesPerformersQuery({ variables: { q } });
  }
  public static useMarkerStrings() {
    return GQL.useMarkerStringsQuery();
  }
  public static useAllTags() {
    return GQL.useAllTagsQuery();
  }
  public static useAllTagsForFilter() {
    return GQL.useAllTagsForFilterQuery();
  }
  public static useAllPerformersForFilter() {
    return GQL.useAllPerformersForFilterQuery();
  }
  public static useAllStudiosForFilter() {
    return GQL.useAllStudiosForFilterQuery();
  }
  public static useAllMoviesForFilter() {
    return GQL.useAllMoviesForFilterQuery();
  }
  public static useValidGalleriesForScene(sceneId: string) {
    return GQL.useValidGalleriesForSceneQuery({
      variables: { scene_id: sceneId }
    });
  }
  public static useStats() {
    return GQL.useStatsQuery();
  }
  public static useVersion() {
    return GQL.useVersionQuery();
  }
  public static useLatestVersion() {
    return GQL.useLatestVersionQuery({
      notifyOnNetworkStatusChange: true,
      errorPolicy: "ignore"
    });
  }

  public static useConfiguration() {
    return GQL.useConfigurationQuery();
  }
  public static useDirectories(path?: string) {
    return GQL.useDirectoriesQuery({ variables: { path } });
  }

  private static performerMutationImpactedQueries = [
    "findPerformers",
    "findScenes",
    "findSceneMarkers",
    "allPerformers"
  ];

  public static usePerformerCreate() {
    return GQL.usePerformerCreateMutation({
      update: () =>
        StashService.invalidateQueries(
          StashService.performerMutationImpactedQueries
        )
    });
  }
  public static usePerformerUpdate() {
    return GQL.usePerformerUpdateMutation({
      update: () =>
        StashService.invalidateQueries(
          StashService.performerMutationImpactedQueries
        )
    });
  }
  public static usePerformerDestroy() {
    return GQL.usePerformerDestroyMutation({
      update: () =>
        StashService.invalidateQueries(
          StashService.performerMutationImpactedQueries
        )
    });
  }

  private static sceneMutationImpactedQueries = [
    "findPerformers",
    "findScenes",
    "findSceneMarkers",
    "findStudios",
    "findMovies",
    "allTags"
    // TODO - add "findTags" when it is implemented
  ];

  public static useSceneUpdate(input: GQL.SceneUpdateInput) {
    return GQL.useSceneUpdateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.sceneMutationImpactedQueries
        ),
      refetchQueries: ["AllTagsForFilter"]
    });
  }

  // remove findScenes for bulk scene update so that we don't lose
  // existing results
  private static sceneBulkMutationImpactedQueries = [
    "findPerformers",
    "findSceneMarkers",
    "findStudios",
    "findMovies",
    "allTags"
  ];

  public static useBulkSceneUpdate(input: GQL.BulkSceneUpdateInput) {
    return GQL.useBulkSceneUpdateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.sceneBulkMutationImpactedQueries
        )
    });
  }

  public static useScenesUpdate(input: GQL.SceneUpdateInput[]) {
    return GQL.useScenesUpdateMutation({ variables: { input } });
  }

  public static useSceneIncrementO(id: string) {
    return GQL.useSceneIncrementOMutation({
      variables: { id }
    });
  }

  public static useSceneDecrementO(id: string) {
    return GQL.useSceneDecrementOMutation({
      variables: { id }
    });
  }

  public static useSceneResetO(id: string) {
    return GQL.useSceneResetOMutation({
      variables: { id }
    });
  }

  public static useSceneDestroy(input: GQL.SceneDestroyInput) {
    return GQL.useSceneDestroyMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.sceneMutationImpactedQueries
        )
    });
  }

  public static useSceneGenerateScreenshot() {
    return GQL.useSceneGenerateScreenshotMutation({
      update: () => StashService.invalidateQueries(["findScenes"])
    });
  }

  private static studioMutationImpactedQueries = [
    "findStudios",
    "findScenes",
    "allStudios"
  ];

  public static useStudioCreate(input: GQL.StudioCreateInput) {
    return GQL.useStudioCreateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.studioMutationImpactedQueries
        )
    });
  }

  public static useStudioUpdate(input: GQL.StudioUpdateInput) {
    return GQL.useStudioUpdateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.studioMutationImpactedQueries
        )
    });
  }

  public static useStudioDestroy(input: GQL.StudioDestroyInput) {
    return GQL.useStudioDestroyMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.studioMutationImpactedQueries
        )
    });
  }

  private static movieMutationImpactedQueries = [
    "findMovies",
    "findScenes",
    "allMovies"
  ];

  public static useMovieCreate(input: GQL.MovieCreateInput) {
    return GQL.useMovieCreateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.movieMutationImpactedQueries
        )
    });
  }

  public static useMovieUpdate(input: GQL.MovieUpdateInput) {
    return GQL.useMovieUpdateMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.movieMutationImpactedQueries
        )
    });
  }

  public static useMovieDestroy(input: GQL.MovieDestroyInput) {
    return GQL.useMovieDestroyMutation({
      variables: input,
      update: () =>
        StashService.invalidateQueries(
          StashService.movieMutationImpactedQueries
        )
    });
  }

  private static tagMutationImpactedQueries = [
    "findScenes",
    "findSceneMarkers",
    "sceneMarkerTags",
    "allTags"
  ];

  public static useTagCreate(input: GQL.TagCreateInput) {
    return GQL.useTagCreateMutation({
      variables: input,
      refetchQueries: ["AllTags", "AllTagsForFilter"]
      // update: () => StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
    });
  }
  public static useTagUpdate(input: GQL.TagUpdateInput) {
    return GQL.useTagUpdateMutation({
      variables: input,
      refetchQueries: ["AllTags", "AllTagsForFilter"]
    });
  }
  public static useTagDestroy(input: GQL.TagDestroyInput) {
    return GQL.useTagDestroyMutation({
      variables: input,
      refetchQueries: ["AllTags", "AllTagsForFilter"],
      update: () =>
        StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
    });
  }

  public static useConfigureGeneral(input: GQL.ConfigGeneralInput) {
    return GQL.useConfigureGeneralMutation({
      variables: { input },
      refetchQueries: ["Configuration"]
    });
  }

  public static useConfigureInterface(input: GQL.ConfigInterfaceInput) {
    return GQL.useConfigureInterfaceMutation({
      variables: { input },
      refetchQueries: ["Configuration"]
    });
  }

  public static useMetadataUpdate() {
    return GQL.useMetadataUpdateSubscription();
  }

  public static useLoggingSubscribe() {
    return GQL.useLoggingSubscribeSubscription();
  }

  public static useLogs() {
    return GQL.useLogsQuery({
      fetchPolicy: "no-cache"
    });
  }

  public static useJobStatus() {
    return GQL.useJobStatusQuery({
      fetchPolicy: "no-cache"
    });
  }

  public static mutateStopJob() {
    return StashService.client.mutate<GQL.StopJobMutation>({
      mutation: GQL.StopJobDocument
    });
  }

  public static queryScrapeFreeones(performerName: string) {
    return StashService.client.query<GQL.ScrapeFreeonesQuery>({
      query: GQL.ScrapeFreeonesDocument,
      variables: {
        performer_name: performerName
      }
    });
  }

  public static queryScrapePerformer(
    scraperId: string,
    scrapedPerformer: GQL.ScrapedPerformerInput
  ) {
    return StashService.client.query<GQL.ScrapePerformerQuery>({
      query: GQL.ScrapePerformerDocument,
      variables: {
        scraper_id: scraperId,
        scraped_performer: scrapedPerformer
      }
    });
  }

  public static queryScrapePerformerURL(url: string) {
    return StashService.client.query<GQL.ScrapePerformerUrlQuery>({
      query: GQL.ScrapePerformerUrlDocument,
      variables: {
        url
      }
    });
  }

  public static queryScrapeSceneURL(url: string) {
    return StashService.client.query<GQL.ScrapeSceneUrlQuery>({
      query: GQL.ScrapeSceneUrlDocument,
      variables: {
        url
      }
    });
  }

  public static queryScrapeScene(
    scraperId: string,
    scene: GQL.SceneUpdateInput
  ) {
    return StashService.client.query<GQL.ScrapeSceneQuery>({
      query: GQL.ScrapeSceneDocument,
      variables: {
        scraper_id: scraperId,
        scene
      }
    });
  }

  public static mutateMetadataScan(input: GQL.ScanMetadataInput) {
    return StashService.client.mutate<GQL.MetadataScanMutation>({
      mutation: GQL.MetadataScanDocument,
      variables: { input }
    });
  }

  public static mutateMetadataAutoTag(input: GQL.AutoTagMetadataInput) {
    return StashService.client.mutate<GQL.MetadataAutoTagMutation>({
      mutation: GQL.MetadataAutoTagDocument,
      variables: { input }
    });
  }

  public static mutateMetadataGenerate(input: GQL.GenerateMetadataInput) {
    return StashService.client.mutate<GQL.MetadataGenerateMutation>({
      mutation: GQL.MetadataGenerateDocument,
      variables: { input }
    });
  }

  public static mutateMetadataClean() {
    return StashService.client.mutate<GQL.MetadataCleanMutation>({
      mutation: GQL.MetadataCleanDocument
    });
  }

  public static mutateMetadataExport() {
    return StashService.client.mutate<GQL.MetadataExportMutation>({
      mutation: GQL.MetadataExportDocument
    });
  }

  public static mutateMetadataImport() {
    return StashService.client.mutate<GQL.MetadataImportMutation>({
      mutation: GQL.MetadataImportDocument
    });
  }

  public static querySceneByPathRegex(filter: GQL.FindFilterType) {
    return StashService.client.query<GQL.FindScenesByPathRegexQuery>({
      query: GQL.FindScenesByPathRegexDocument,
      variables: { filter }
    });
  }

  public static queryParseSceneFilenames(
    filter: GQL.FindFilterType,
    config: GQL.SceneParserInput
  ) {
    return StashService.client.query<GQL.ParseSceneFilenamesQuery>({
      query: GQL.ParseSceneFilenamesDocument,
      variables: { filter, config },
      fetchPolicy: "network-only"
    });
  }

  private static stringGenderMap = new Map<string, GQL.GenderEnum>(
    [["Male", GQL.GenderEnum.Male],
    ["Female", GQL.GenderEnum.Female],
    ["Transgender Male", GQL.GenderEnum.TransgenderMale],
    ["Transgender Female", GQL.GenderEnum.TransgenderFemale],
    ["Intersex", GQL.GenderEnum.Intersex]]
  );

  public static genderToString(value?: GQL.GenderEnum) {
    if (!value) {
      return undefined;
    }

    const foundEntry = Array.from(StashService.stringGenderMap.entries()).find((e) => {
      return e[1] === value;
    });

    if (foundEntry) {
      return foundEntry[0];
    }
  }

  public static stringToGender(value?: string) {
    if (!value) {
      return undefined;
    }

    return StashService.stringGenderMap.get(value);
  }

  public static getGenderStrings() {
    return Array.from(StashService.stringGenderMap.keys());
  }

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  private constructor() {}
}
