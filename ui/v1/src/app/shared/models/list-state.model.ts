import {
  SceneFilterType,
  ResolutionEnum,
  PerformerFilterType,
  SceneMarkerFilterType,
  SlimSceneData,
  PerformerData,
  StudioData,
  GalleryData,
  SceneMarkerData
} from '../../core/graphql-generated';

import { StashService } from '../../core/stash.service';

export enum DisplayMode {
  Grid,
  List,
  Wall
}

export enum FilterMode {
  Scenes,
  Performers,
  Studios,
  Galleries,
  SceneMarkers
}

export class CustomCriteria {
  key: string;
  value: string;
  constructor(key: string, value: string) {
    this.key = key;
    this.value = value;
  }
}

export enum CriteriaType {
  None,
  Rating,
  Resolution,
  Favorite,
  HasMarkers,
  IsMissing,
  Tags,
  SceneTags,
  Performers
}

export enum CriteriaValueType {
  Single,
  Multiple
}

export class CriteriaOption {
  name: string;
  value: CriteriaType;
  constructor(type: CriteriaType, name: string = CriteriaType[type]) {
    this.name = name;
    this.value = type;
  }
}

interface CriteriaConfig {
  valueType: CriteriaValueType;
  parameterName: string;
  options: any[];
}

export class Criteria {
  type: CriteriaType;
  valueType: CriteriaValueType;
  options: any[] = [];
  parameterName: string;
  value: string;
  values: string[];

  private stashService: StashService;

  async configure(type: CriteriaType, stashService: StashService) {
    this.type = type;
    this.stashService = stashService;

    let config: CriteriaConfig = {
      valueType: CriteriaValueType.Single,
      parameterName: '',
      options: []
    };

    switch (type) {
      case CriteriaType.Rating:
        config.parameterName = 'rating';
        config.options = [1, 2, 3, 4, 5];
        break;
      case CriteriaType.Resolution:
        config.parameterName = 'resolution';
        config.options = ['240p', '480p', '720p', '1080p', '4k'];
        break;
      case CriteriaType.Favorite:
        config.parameterName = 'filter_favorites';
        config.options = ['true', 'false'];
        break;
      case CriteriaType.HasMarkers:
        config.parameterName = 'has_markers';
        config.options = ['true', 'false'];
        break;
      case CriteriaType.IsMissing:
        config.parameterName = 'is_missing';
        config.options = ['title', 'url', 'date', 'gallery', 'studio', 'performers'];
        break;
      case CriteriaType.Tags:
        config = await this.configureTags('tags');
        break;
      case CriteriaType.SceneTags:
        config = await this.configureTags('scene_tags');
        break;
      case CriteriaType.Performers:
        config = await this.configurePerformers('performers');
        break;
      case CriteriaType.None:
      default: break;
    }

    this.valueType = config.valueType;
    this.parameterName = config.parameterName;
    this.options = config.options;

    this.value = ''; // Need this or else we send invalid value to the new filter
    // this.values = []; // TODO this seems to break the "Multiple" filters
  }

  private async configureTags(name: string) {
    const result = await this.stashService.allTagsForFilter().result();
    return {
      valueType: CriteriaValueType.Multiple,
      parameterName: name,
      options: result.data.allTags.map(item => {
        return { id: item.id, name: item.name };
      })
    };
  }

  private async configurePerformers(name: string) {
    const result = await this.stashService.allPerformersForFilter().result();
    return {
      valueType: CriteriaValueType.Multiple,
      parameterName: name,
      options: result.data.allPerformers.map(item => {
        return { id: item.id, name: item.name, image_path: item.image_path };
      })
    };
  }
}

export class ListFilter {
  searchTerm?: string;
  performers?: number[];
  currentPage = 1;
  itemsPerPage = 40;
  sortDirection = 'asc';
  sortBy: string;
  displayModeOptions: DisplayMode[] = [];
  displayMode: DisplayMode;
  filterMode: FilterMode;
  sortByOptions: string[];
  criteriaFilterOpen = false;
  criteriaOptions: CriteriaOption[];
  criterions: Criteria[] = [];
  customCriteria: CustomCriteria[] = [];

  configureForFilterMode(filterMode: FilterMode) {
    switch (filterMode) {
      case FilterMode.Scenes:
        if (!!this.sortBy === false) { this.sortBy = 'date'; }
        this.sortByOptions = ['title', 'rating', 'date', 'filesize', 'duration', 'framerate', 'bitrate', 'random'];
        this.displayModeOptions = [
          DisplayMode.Grid,
          DisplayMode.List,
          DisplayMode.Wall
        ];
        this.criteriaOptions = [
          new CriteriaOption(CriteriaType.None),
          new CriteriaOption(CriteriaType.Rating),
          new CriteriaOption(CriteriaType.Resolution),
          new CriteriaOption(CriteriaType.HasMarkers),
          new CriteriaOption(CriteriaType.IsMissing),
          new CriteriaOption(CriteriaType.Tags)
        ];
        break;
      case FilterMode.Performers:
        if (!!this.sortBy === false) { this.sortBy = 'name'; }
        this.sortByOptions = ['name', 'height', 'birthdate', 'scenes_count'];
        this.displayModeOptions = [
          DisplayMode.Grid,
          DisplayMode.List
        ];
        this.criteriaOptions = [
          new CriteriaOption(CriteriaType.None),
          new CriteriaOption(CriteriaType.Favorite)
        ];
        break;
      case FilterMode.Studios:
        if (!!this.sortBy === false) { this.sortBy = 'name'; }
        this.sortByOptions = ['name', 'scenes_count'];
        this.displayModeOptions = [
          DisplayMode.Grid
        ];
        this.criteriaOptions = [
          new CriteriaOption(CriteriaType.None)
        ];
        break;
      case FilterMode.Galleries:
        if (!!this.sortBy === false) { this.sortBy = 'title'; }
        this.sortByOptions = ['title', 'path'];
        this.displayModeOptions = [
          DisplayMode.Grid
        ];
        this.criteriaOptions = [
          new CriteriaOption(CriteriaType.None)
        ];
        break;
      case FilterMode.SceneMarkers:
        if (!!this.sortBy === false) { this.sortBy = 'title'; }
        this.sortByOptions = ['title', 'seconds', 'scene_id', 'random', 'scenes_updated_at'];
        this.displayModeOptions = [
          DisplayMode.Wall
        ];
        this.criteriaOptions = [
          new CriteriaOption(CriteriaType.None),
          new CriteriaOption(CriteriaType.Tags),
          new CriteriaOption(CriteriaType.SceneTags),
          new CriteriaOption(CriteriaType.Performers)
        ];
        break;
      default:
        this.sortByOptions = [];
        this.displayModeOptions = [];
        this.criteriaOptions = [new CriteriaOption(CriteriaType.None)];
        break;
    }
    if (!!this.displayMode === false) { this.displayMode = this.displayModeOptions[0]; }
  }

  configureFromQueryParameters(params, stashService: StashService) {
    if (params['sortby'] != null) {
      this.sortBy = params['sortby'];
    }
    if (params['sortdir'] != null) {
      this.sortDirection = params['sortdir'];
    }
    if (params['disp'] != null) {
      this.displayMode = params['disp'];
    }
    if (params['q'] != null) {
      this.searchTerm = params['q'];
    }
    if (params['p'] != null) {
      this.currentPage = Number(params['p']);
    }

    if (params['c'] != null) {
      this.criterions = [];

      let jsonParameters: any[];
      if (params['c'] instanceof Array) {
        jsonParameters = params['c'];
      } else {
        jsonParameters = [params['c']];
      }

      if (jsonParameters.length !== 0) {
        this.criteriaFilterOpen = true;
      }

      jsonParameters.forEach(jsonString => {
        const encodedCriteria = JSON.parse(jsonString);
        const criteria = new Criteria();
        criteria.configure(encodedCriteria.type, stashService);
        if (criteria.valueType === CriteriaValueType.Single) {
          criteria.value = encodedCriteria.value;
        } else {
          criteria.values = encodedCriteria.values;
        }
        this.criterions.push(criteria);
      });
    }
  }

  makeQueryParameters(): any {
    const encodedCriterion = [];
    this.criterions.forEach(criteria => {
      const encodedCriteria: any = {};
      encodedCriteria.type = criteria.type;
      if (criteria.valueType === CriteriaValueType.Single) {
        encodedCriteria.value = criteria.value;
      } else {
        encodedCriteria.values = criteria.values;
      }
      const jsonCriteria = JSON.stringify(encodedCriteria);
      encodedCriterion.push(jsonCriteria);
    });

    const result = {
      queryParams: {
        sortby: this.sortBy,
        sortdir: this.sortDirection,
        disp: this.displayMode,
        q: this.searchTerm,
        p: this.currentPage,
        c: encodedCriterion
      },
      queryParamsHandling: 'merge'
    };
    return result;
  }

  // TODO: These don't support multiple of the same criteria, only the last one set is used.

  makeSceneFilter(): SceneFilterType {
    const result: SceneFilterType = {};
    this.criterions.forEach(criteria => {
      switch (criteria.type) {
        case CriteriaType.Rating:
          result.rating = Number(criteria.value);
          break;
        case CriteriaType.Resolution: {
          switch (criteria.value) {
            case '240p': result.resolution = ResolutionEnum.Low; break;
            case '480p': result.resolution = ResolutionEnum.Standard; break;
            case '720p': result.resolution = ResolutionEnum.StandardHd; break;
            case '1080p': result.resolution = ResolutionEnum.FullHd; break;
            case '4k': result.resolution = ResolutionEnum.FourK; break;
          }
          break;
        }
        case CriteriaType.HasMarkers:
          result.has_markers = criteria.value;
          break;
        case CriteriaType.IsMissing:
          result.is_missing = criteria.value;
          break;
        case CriteriaType.Tags:
          result.tags = criteria.values;
          break;
      }
    });
    return result;
  }

  makePerformerFilter(): PerformerFilterType {
    const result: PerformerFilterType = {};
    this.criterions.forEach(criteria => {
      switch (criteria.type) {
        case CriteriaType.Favorite:
          result.filter_favorites = criteria.value === 'true';
          break;
      }
    });
    return result;
  }

  makeSceneMarkerFilter(): SceneMarkerFilterType {
    const result: SceneMarkerFilterType = {};
    this.criterions.forEach(criteria => {
      switch (criteria.type) {
        case CriteriaType.Tags:
          result.tags = criteria.values;
          break;
        case CriteriaType.SceneTags:
          result.scene_tags = criteria.values;
          break;
        case CriteriaType.Performers:
          result.performers = criteria.values;
          break;
      }
    });
    return result;
  }
}

export class ListState<T> {
  totalCount: number;
  scrollY: number;
  filter: ListFilter = new ListFilter();
  data: T[];

  reset() {
    this.data = null;
    this.totalCount = null;
  }
}

export class SceneListState extends ListState<SlimSceneData.Fragment> {
  constructor() {
    super();
    this.filter.filterMode = FilterMode.Scenes;
  }
}

export class PerformerListState extends ListState<PerformerData.Fragment> {
  constructor() {
    super();
    this.filter.filterMode = FilterMode.Performers;
  }
}

export class StudioListState extends ListState<StudioData.Fragment> {
  constructor() {
    super();
    this.filter.filterMode = FilterMode.Studios;
  }
}

export class GalleryListState extends ListState<GalleryData.Fragment> {
  constructor() {
    super();
    this.filter.filterMode = FilterMode.Galleries;
  }
}

export class SceneMarkerListState extends ListState<SceneMarkerData.Fragment> {
  constructor() {
    super();
    this.filter.filterMode = FilterMode.SceneMarkers;
  }
}
