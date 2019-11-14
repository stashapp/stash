import ApolloClient from "apollo-client";
import { WebSocketLink } from 'apollo-link-ws';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink, split } from "apollo-boost";
import _ from "lodash";
import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";
import { getMainDefinition } from "apollo-utilities";

export class StashService {
  public static client: ApolloClient<any>;
  private static cache: InMemoryCache;

  public static initialize() {
    const platformUrl = new URL(window.location.origin);
    const wsPlatformUrl = new URL(window.location.origin);
    wsPlatformUrl.protocol = "ws:";

    if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
      platformUrl.port = "9999"; // TODO: Hack. Development expects port 9999
      wsPlatformUrl.port = "9999";

      if (process.env.REACT_APP_HTTPS === "true") {
        platformUrl.protocol = "https:";
      }
    }

    if (platformUrl.protocol === "https:") {
      wsPlatformUrl.protocol = "wss:";
    }
    const url = platformUrl.toString().slice(0, -1) + "/graphql";
    const wsUrl = wsPlatformUrl.toString().slice(0, -1) + "/graphql";

    const httpLink = new HttpLink({
      uri: url,
    });

    const wsLink = new WebSocketLink({
      uri: wsUrl,
      options: {
        reconnect: true
      },
    });
    
    const link = split(
      ({ query }) => {
        const { kind, operation } = getMainDefinition(query);
        return kind === 'OperationDefinition' && operation === 'subscription';
      },
      wsLink,
      httpLink,
    );

    StashService.cache = new InMemoryCache();
    StashService.client = new ApolloClient({
      link: link,
      cache: StashService.cache
    });

    (window as any).StashService = StashService;
    return StashService.client;
  }

  private static invalidateCache() {
    StashService.client.resetStore();
  }

  private static invalidateQueries(queries : string[]) {
    if (!!StashService.cache) {
      const cache = StashService.cache as any;
      const keyMatchers = queries.map(query => {
        return new RegExp("^" + query);
      });

      const rootQuery = cache.data.data.ROOT_QUERY;
      Object.keys(rootQuery).forEach(key => {
        if (keyMatchers.some(matcher => {
          return !!key.match(matcher);
        })) {
          delete rootQuery[key];
        }
      });
    }
  }

  public static useFindGalleries(filter: ListFilterModel) {
    return GQL.useFindGalleries({
      variables: {
        filter: filter.makeFindFilter(),
      },
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

    return GQL.useFindScenes({
      variables: {
        filter: filter.makeFindFilter(),
        scene_filter: sceneFilter,
      },
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

    return GQL.useFindSceneMarkers({
      variables: {
        filter: filter.makeFindFilter(),
        scene_marker_filter: sceneMarkerFilter,
      },
    });
  }

  public static useFindStudios(filter: ListFilterModel) {
    return GQL.useFindStudios({
      variables: {
        filter: filter.makeFindFilter(),
      },
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

    return GQL.useFindPerformers({
      variables: {
        filter: filter.makeFindFilter(),
        performer_filter: performerFilter,
      },
    });
  }

  public static useFindGallery(id: string) { return GQL.useFindGallery({variables: {id}}); }
  public static useFindScene(id: string) { return GQL.useFindScene({variables: {id}}); }
  public static useFindPerformer(id: string) {
    const skip = id === "new" ? true : false;
    return GQL.useFindPerformer({variables: {id}, skip});
  }
  public static useFindStudio(id: string) {
    const skip = id === "new" ? true : false;
    return GQL.useFindStudio({variables: {id}, skip});
  }

  // TODO - scene marker manipulation functions are handled differently
  private static sceneMarkerMutationImpactedQueries = [
    "findSceneMarkers",
    "findScenes",
    "markerStrings",
    "sceneMarkerTags"
  ];

  public static useSceneMarkerCreate() {
    return GQL.useSceneMarkerCreate(); 
  }
  public static useSceneMarkerUpdate() { 
    return GQL.useSceneMarkerUpdate(); 
  }
  public static useSceneMarkerDestroy() {
    return GQL.useSceneMarkerDestroy(); 
  }

  public static useListScrapers(scraperType: GQL.ScraperType) { 
    return GQL.useListScrapers({
      variables: {
        scraper_type: scraperType
      }
    }); 
  }
  public static useScrapePerformerList(scraperId: string, q : string) { 
    return GQL.useScrapePerformerList({ variables: { scraper_id: scraperId, query: q }}); 
  }
  public static useScrapePerformer(scraperId: string, scrapedPerformer : GQL.ScrapedPerformerInput) {
     return GQL.useScrapePerformer({ variables: { scraper_id: scraperId, scraped_performer: scrapedPerformer }});
  }

  public static useScrapeFreeonesPerformers(q: string) { return GQL.useScrapeFreeonesPerformers({ variables: { q } }); }
  public static useMarkerStrings() { return GQL.useMarkerStrings(); }
  public static useAllTags() { return GQL.useAllTags(); }
  public static useAllTagsForFilter() { return GQL.useAllTagsForFilter(); }
  public static useAllPerformersForFilter() { return GQL.useAllPerformersForFilter(); }
  public static useAllStudiosForFilter() { return GQL.useAllStudiosForFilter(); }
  public static useValidGalleriesForScene(sceneId: string) {
    return GQL.useValidGalleriesForScene({variables: {scene_id: sceneId}});
  }
  public static useStats() { return GQL.useStats(); }
  public static useVersion() { return GQL.useVersion(); }

  public static useConfiguration() { return GQL.useConfiguration(); }
  public static useDirectories(path?: string) { return GQL.useDirectories({ variables: { path }}); }

  private static performerMutationImpactedQueries = [
    "findPerformers",
    "findScenes",
    "findSceneMarkers",
    "allPerformers"
  ];

  public static usePerformerCreate(input: GQL.PerformerCreateInput) {
    return GQL.usePerformerCreate({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.performerMutationImpactedQueries)
    });
  }
  public static usePerformerUpdate(input: GQL.PerformerUpdateInput) {
    return GQL.usePerformerUpdate({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.performerMutationImpactedQueries)
    });
  }
  public static usePerformerDestroy(input: GQL.PerformerDestroyInput) {
    return GQL.usePerformerDestroy({
      variables: input,
      update: () => StashService.invalidateQueries(StashService.performerMutationImpactedQueries)
    });
  }

  private static sceneMutationImpactedQueries = [
    "findPerformers",
    "findScenes",
    "findSceneMarkers",
    "findStudios",
    "allTags"
  ];

  public static useSceneUpdate(input: GQL.SceneUpdateInput) {
    return GQL.useSceneUpdate({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.sceneMutationImpactedQueries)
    });
  }

  // remove findScenes for bulk scene update so that we don't lose
  // existing results
  private static sceneBulkMutationImpactedQueries = [
    "findPerformers",
    "findSceneMarkers",
    "findStudios",
    "allTags"
  ];

  public static useBulkSceneUpdate(input: GQL.BulkSceneUpdateInput) {
    return GQL.useBulkSceneUpdate({ 
      variables: input, 
      update: () => StashService.invalidateQueries(StashService.sceneBulkMutationImpactedQueries)
    });
  }
  
  public static useScenesUpdate(input: GQL.SceneUpdateInput[]) {
    return GQL.useScenesUpdate({ variables: { input : input }});
  }

  public static useSceneDestroy(input: GQL.SceneDestroyInput) {
    return GQL.useSceneDestroy({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.sceneMutationImpactedQueries)
    });
  }

  private static studioMutationImpactedQueries = [
    "findStudios",
    "findScenes",
    "allStudios"
  ];

  public static useStudioCreate(input: GQL.StudioCreateInput) {
    return GQL.useStudioCreate({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.studioMutationImpactedQueries)
    });
  }

  public static useStudioUpdate(input: GQL.StudioUpdateInput) {
    return GQL.useStudioUpdate({ 
      variables: input,
      update: () => StashService.invalidateQueries(StashService.studioMutationImpactedQueries)
    });
  }

  public static useStudioDestroy(input: GQL.StudioDestroyInput) {
    return GQL.useStudioDestroy({ 
      variables: input, 
      update: () => StashService.invalidateQueries(StashService.studioMutationImpactedQueries)
    });
  }

  private static tagMutationImpactedQueries = [
    "findScenes",
    "findSceneMarkers",
    "sceneMarkerTags",
    "allTags"
  ];

  public static useTagCreate(input: GQL.TagCreateInput) {
    return GQL.useTagCreate({ 
      variables: input, 
      refetchQueries: ["AllTags"],
      update: () => StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
    });
  }
  public static useTagUpdate(input: GQL.TagUpdateInput) {
    return GQL.useTagUpdate({ 
      variables: input, 
      refetchQueries: ["AllTags"],
      update: () => StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
    });
  }
  public static useTagDestroy(input: GQL.TagDestroyInput) {
    return GQL.useTagDestroy({ 
      variables: input, 
      refetchQueries: ["AllTags"],
      update: () => StashService.invalidateQueries(StashService.tagMutationImpactedQueries)
    });
  }

  public static useConfigureGeneral(input: GQL.ConfigGeneralInput) {
    return GQL.useConfigureGeneral({ variables: { input }, refetchQueries: ["Configuration"] });
  }

  public static useConfigureInterface(input: GQL.ConfigInterfaceInput) {
    return GQL.useConfigureInterface({ variables: { input }, refetchQueries: ["Configuration"] });
  }

  public static useMetadataUpdate() {
    return GQL.useMetadataUpdate();
  }

  public static useLoggingSubscribe() {
    return GQL.useLoggingSubscribe();
  }

  public static useLogs() {
    return GQL.useLogs({
      fetchPolicy: 'no-cache'
    });
  }

  public static useJobStatus() {
    return GQL.useJobStatus({
      fetchPolicy: 'no-cache'
    });
  }

  public static queryStopJob() {
    return StashService.client.query<GQL.StopJobQuery>({
      query: GQL.StopJobDocument,
      fetchPolicy: "network-only",
    });
  }

  public static queryScrapeFreeones(performerName: string) {
    return StashService.client.query<GQL.ScrapeFreeonesQuery>({
      query: GQL.ScrapeFreeonesDocument,
      variables: {
        performer_name: performerName,
      },
    });
  }

  public static queryScrapePerformer(scraperId: string, scrapedPerformer: GQL.ScrapedPerformerInput) {
    return StashService.client.query<GQL.ScrapePerformerQuery>({
      query: GQL.ScrapePerformerDocument,
      variables: {
        scraper_id: scraperId,
        scraped_performer: scrapedPerformer,
      },
    });
  }

  public static queryScrapePerformerURL(url: string) {
    return StashService.client.query<GQL.ScrapePerformerUrlQuery>({
      query: GQL.ScrapePerformerUrlDocument,
      variables: {
        url: url,
      },
    });
  }

  public static queryMetadataScan(input: GQL.ScanMetadataInput) {
    return StashService.client.query<GQL.MetadataScanQuery>({
      query: GQL.MetadataScanDocument,
      variables: { input },
      fetchPolicy: "network-only",
    });
  }

  public static queryMetadataGenerate(input: GQL.GenerateMetadataInput) {
    return StashService.client.query<GQL.MetadataGenerateQuery>({
      query: GQL.MetadataGenerateDocument,
      variables: { input },
      fetchPolicy: "network-only",
    });
  }

  public static queryMetadataClean() {
    return StashService.client.query<GQL.MetadataCleanQuery>({
      query: GQL.MetadataCleanDocument,
      fetchPolicy: "network-only",
    });
  }

  public static queryMetadataExport() {
    return StashService.client.query<GQL.MetadataExportQuery>({
      query: GQL.MetadataExportDocument,
      fetchPolicy: "network-only",
    });
  }

  public static queryMetadataImport() {
    return StashService.client.query<GQL.MetadataImportQuery>({
      query: GQL.MetadataImportDocument,
      fetchPolicy: "network-only",
    });
  }

  public static querySceneByPathRegex(filter: GQL.FindFilterType) {
    return StashService.client.query<GQL.FindScenesByPathRegexQuery>({
      query: GQL.FindScenesByPathRegexDocument,
      variables: {filter: filter},
    });
  }

  public static queryParseSceneFilenames(filter: GQL.FindFilterType, config: GQL.SceneParserInput) {
    return StashService.client.query<GQL.ParseSceneFilenamesQuery>({
      query: GQL.ParseSceneFilenamesDocument,
      variables: {filter: filter, config: config},
      fetchPolicy: "network-only",
    });
  }

  public static nullToUndefined(value: any): any {
    if (_.isPlainObject(value)) {
      return _.mapValues(value, StashService.nullToUndefined);
    }
    if (_.isArray(value)) {
      return value.map(StashService.nullToUndefined);
    }
    if (value === null) {
      return undefined;
    }
    return value;
  }

  private constructor() {}
}
