import ApolloClient from "apollo-boost";
import _ from "lodash";
import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";

export class StashService {
  public static client: ApolloClient<any>;

  public static initialize() {
    const platformUrl = new URL(window.location.origin);
    if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
      platformUrl.port = "9999"; // TODO: Hack. Development expects port 9999

      if (process.env.REACT_APP_HTTPS === "true") {
        platformUrl.protocol = "https:";
      }
    }
    const url = platformUrl.toString().slice(0, -1);

    StashService.client = new ApolloClient({
      uri: `${url}/graphql`,
    });

    (window as any).StashService = StashService;
    return StashService.client;
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

  public static useSceneMarkerCreate() { return GQL.useSceneMarkerCreate({ refetchQueries: ["FindScene"] }); }
  public static useSceneMarkerUpdate() { return GQL.useSceneMarkerUpdate({ refetchQueries: ["FindScene"] }); }
  public static useSceneMarkerDestroy() { return GQL.useSceneMarkerDestroy({ refetchQueries: ["FindScene"] }); }

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

  public static useConfiguration() { return GQL.useConfiguration(); }
  public static useDirectories(path?: string) { return GQL.useDirectories({ variables: { path }}); }

  public static usePerformerCreate(input: GQL.PerformerCreateInput) {
    return GQL.usePerformerCreate({ variables: input });
  }
  public static usePerformerUpdate(input: GQL.PerformerUpdateInput) {
    return GQL.usePerformerUpdate({ variables: input });
  }
  public static usePerformerDestroy(input: GQL.PerformerDestroyInput) {
    return GQL.usePerformerDestroy({ variables: input });
  }

  public static useSceneUpdate(input: GQL.SceneUpdateInput) {
    return GQL.useSceneUpdate({ variables: input });
  }

  public static useStudioCreate(input: GQL.StudioCreateInput) {
    return GQL.useStudioCreate({ variables: input });
  }
  public static useStudioUpdate(input: GQL.StudioUpdateInput) {
    return GQL.useStudioUpdate({ variables: input });
  }
  public static useStudioDestroy(input: GQL.StudioDestroyInput) {
    return GQL.useStudioDestroy({ variables: input });
  }

  public static useTagCreate(input: GQL.TagCreateInput) {
    return GQL.useTagCreate({ variables: input, refetchQueries: ["AllTags"] });
  }
  public static useTagUpdate(input: GQL.TagUpdateInput) {
    return GQL.useTagUpdate({ variables: input, refetchQueries: ["AllTags"] });
  }
  public static useTagDestroy(input: GQL.TagDestroyInput) {
    return GQL.useTagDestroy({ variables: input, refetchQueries: ["AllTags"] });
  }

  public static useConfigureGeneral(input: GQL.ConfigGeneralInput) {
    return GQL.useConfigureGeneral({ variables: { input }, refetchQueries: ["Configuration"] });
  }

  public static useConfigureInterface(input: GQL.ConfigInterfaceInput) {
    return GQL.useConfigureInterface({ variables: { input }, refetchQueries: ["Configuration"] });
  }

  public static queryScrapeFreeones(performerName: string) {
    return StashService.client.query<GQL.ScrapeFreeonesQuery>({
      query: GQL.ScrapeFreeonesDocument,
      variables: {
        performer_name: performerName,
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
