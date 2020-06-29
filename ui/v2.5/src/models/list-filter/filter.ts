import queryString, { ParsedQuery } from "query-string";
import {
  FindFilterType,
  PerformerFilterType,
  ResolutionEnum,
  SceneFilterType,
  SceneMarkerFilterType,
  SortDirectionEnum,
  MovieFilterType,
  StudioFilterType,
  GalleryFilterType,
  TagFilterType,
} from "src/core/generated-graphql";
import { stringToGender } from "src/core/StashService";
import {
  Criterion,
  ICriterionOption,
  CriterionType,
  CriterionOption,
  NumberCriterion,
  StringCriterion,
  DurationCriterion,
} from "./criteria/criterion";
import {
  FavoriteCriterion,
  FavoriteCriterionOption,
} from "./criteria/favorite";
import {
  HasMarkersCriterion,
  HasMarkersCriterionOption,
} from "./criteria/has-markers";
import {
  IsMissingCriterion,
  PerformerIsMissingCriterionOption,
  SceneIsMissingCriterionOption,
  GalleryIsMissingCriterionOption,
  TagIsMissingCriterionOption,
} from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import {
  PerformersCriterion,
  PerformersCriterionOption,
} from "./criteria/performers";
import { RatingCriterion, RatingCriterionOption } from "./criteria/rating";
import {
  ResolutionCriterion,
  ResolutionCriterionOption,
} from "./criteria/resolution";
import {
  StudiosCriterion,
  StudiosCriterionOption,
  ParentStudiosCriterion,
  ParentStudiosCriterionOption,
} from "./criteria/studios";
import {
  SceneTagsCriterionOption,
  TagsCriterion,
  TagsCriterionOption,
} from "./criteria/tags";
import { makeCriteria } from "./criteria/utils";
import { DisplayMode, FilterMode } from "./types";
import { GenderCriterionOption, GenderCriterion } from "./criteria/gender";
import { MoviesCriterionOption, MoviesCriterion } from "./criteria/movies";

interface IQueryParameters {
  perPage?: string;
  sortby?: string;
  sortdir?: string;
  disp?: string;
  q?: string;
  p?: string;
  c?: string[];
}

const DEFAULT_PARAMS = {
  sortDirection: SortDirectionEnum.Asc,
  displayMode: DisplayMode.Grid,
  currentPage: 1,
  itemsPerPage: 40,
};

// TODO: handle customCriteria
export class ListFilterModel {
  public filterMode: FilterMode = FilterMode.Scenes;
  public searchTerm?: string;
  public currentPage = DEFAULT_PARAMS.currentPage;
  public itemsPerPage = DEFAULT_PARAMS.itemsPerPage;
  public sortDirection: SortDirectionEnum = SortDirectionEnum.Asc;
  public sortBy?: string;
  public sortByOptions: string[] = [];
  public displayMode: DisplayMode = DEFAULT_PARAMS.displayMode;
  public displayModeOptions: DisplayMode[] = [];
  public criterionOptions: ICriterionOption[] = [];
  public criteria: Array<Criterion> = [];
  public randomSeed = -1;

  private static createCriterionOption(criterion: CriterionType) {
    return new CriterionOption(Criterion.getLabel(criterion), criterion);
  }

  public constructor(filterMode: FilterMode, rawParms?: ParsedQuery<string>) {
    const params = rawParms as IQueryParameters;
    switch (filterMode) {
      case FilterMode.Scenes:
        this.sortBy = "date";
        this.sortByOptions = [
          "title",
          "path",
          "rating",
          "o_counter",
          "date",
          "filesize",
          "duration",
          "framerate",
          "bitrate",
          "random",
        ];
        this.displayModeOptions = [
          DisplayMode.Grid,
          DisplayMode.List,
          DisplayMode.Wall,
        ];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new RatingCriterionOption(),
          ListFilterModel.createCriterionOption("o_counter"),
          new ResolutionCriterionOption(),
          ListFilterModel.createCriterionOption("duration"),
          new HasMarkersCriterionOption(),
          new SceneIsMissingCriterionOption(),
          new TagsCriterionOption(),
          new PerformersCriterionOption(),
          new StudiosCriterionOption(),
          new MoviesCriterionOption(),
        ];
        break;
      case FilterMode.Performers: {
        this.sortBy = "name";
        this.sortByOptions = ["name", "height", "birthdate", "scenes_count"];
        this.displayModeOptions = [DisplayMode.Grid, DisplayMode.List];

        const numberCriteria: CriterionType[] = ["birth_year", "age"];
        const stringCriteria: CriterionType[] = [
          "ethnicity",
          "country",
          "eye_color",
          "height",
          "measurements",
          "fake_tits",
          "career_length",
          "tattoos",
          "piercings",
          "aliases",
        ];

        this.criterionOptions = [
          new NoneCriterionOption(),
          new FavoriteCriterionOption(),
          new GenderCriterionOption(),
          new PerformerIsMissingCriterionOption(),
          ...numberCriteria
            .concat(stringCriteria)
            .map((c) => ListFilterModel.createCriterionOption(c)),
        ];

        break;
      }
      case FilterMode.Studios:
        this.sortBy = "name";
        this.sortByOptions = ["name", "scenes_count"];
        this.displayModeOptions = [DisplayMode.Grid];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new ParentStudiosCriterionOption(),
        ];
        break;
      case FilterMode.Movies:
        this.sortBy = "name";
        this.sortByOptions = ["name", "scenes_count"];
        this.displayModeOptions = [DisplayMode.Grid];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new StudiosCriterionOption(),
        ];
        break;
      case FilterMode.Galleries:
        this.sortBy = "path";
        this.sortByOptions = ["path"];
        this.displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new GalleryIsMissingCriterionOption(),
        ];
        break;
      case FilterMode.SceneMarkers:
        this.sortBy = "title";
        this.sortByOptions = [
          "title",
          "seconds",
          "scene_id",
          "random",
          "scenes_updated_at",
        ];
        this.displayModeOptions = [DisplayMode.Wall];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new TagsCriterionOption(),
          new SceneTagsCriterionOption(),
          new PerformersCriterionOption(),
        ];
        break;
      case FilterMode.Tags:
        this.sortBy = "name";
        this.sortByOptions = ["name", "scenes_count", "scene_markers_count"];
        this.displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new TagIsMissingCriterionOption(),
          ListFilterModel.createCriterionOption("scene_count"),
          ListFilterModel.createCriterionOption("marker_count"),
        ];
        break;
      default:
        this.sortByOptions = [];
        this.displayModeOptions = [];
        this.criterionOptions = [new NoneCriterionOption()];
        break;
    }
    if (!!this.displayMode === false) {
      this.displayMode = this.displayModeOptions[0];
    }
    this.sortByOptions = [...this.sortByOptions, "created_at", "updated_at"];
    if (params) this.configureFromQueryParameters(params);
  }

  public configureFromQueryParameters(params: IQueryParameters) {
    if (params.sortby !== undefined) {
      this.sortBy = params.sortby;

      // parse the random seed if provided
      const randomPrefix = "random_";
      if (this.sortBy && this.sortBy.startsWith(randomPrefix)) {
        const seedStr = this.sortBy.substring(randomPrefix.length);

        this.sortBy = "random";
        try {
          this.randomSeed = Number.parseInt(seedStr, 10);
        } catch (err) {
          // ignore
        }
      }
    }
    this.sortDirection =
      params.sortdir === "desc"
        ? SortDirectionEnum.Desc
        : SortDirectionEnum.Asc;
    if (params.disp) {
      this.displayMode = Number.parseInt(params.disp, 10);
    }
    if (params.q) {
      this.searchTerm = params.q;
    }
    if (params.p) {
      this.currentPage = Number.parseInt(params.p, 10);
    }
    if (params.perPage) this.itemsPerPage = Number.parseInt(params.perPage, 10);

    if (params.c !== undefined) {
      this.criteria = [];

      let jsonParameters: string[];
      if (params.c instanceof Array) {
        jsonParameters = params.c;
      } else {
        jsonParameters = [params.c];
      }

      jsonParameters.forEach((jsonString) => {
        const encodedCriterion = JSON.parse(jsonString);
        const criterion = makeCriteria(encodedCriterion.type);
        // it's possible that we have unsupported criteria. Just skip if so.
        if (criterion) {
          criterion.value = encodedCriterion.value;
          criterion.modifier = encodedCriterion.modifier;
          this.criteria.push(criterion);
        }
      });
    }
  }

  private setRandomSeed() {
    if (this.sortBy === "random") {
      // #321 - set the random seed if it is not set
      if (this.randomSeed === -1) {
        // generate 8-digit seed
        this.randomSeed = Math.floor(Math.random() * 10 ** 8);
      }
    } else {
      this.randomSeed = -1;
    }
  }

  private getSortBy(): string | undefined {
    this.setRandomSeed();

    if (this.sortBy === "random") {
      return `${this.sortBy}_${this.randomSeed.toString()}`;
    }

    return this.sortBy;
  }

  public makeQueryParameters(): string {
    const encodedCriteria: string[] = [];
    this.criteria.forEach((criterion) => {
      const encodedCriterion: Partial<Criterion> = {
        type: criterion.type,
        // #394 - the presence of a # symbol results in the query URL being
        // malformed. We could set encode: true in the queryString.stringify
        // call below, but this results in a URL that gets pretty long and ugly.
        // Instead, we'll encode the criteria values.
        value: criterion.encodeValue(),
        modifier: criterion.modifier,
      };
      const jsonCriterion = JSON.stringify(encodedCriterion);
      encodedCriteria.push(jsonCriterion);
    });

    const result = {
      perPage:
        this.itemsPerPage !== DEFAULT_PARAMS.itemsPerPage
          ? this.itemsPerPage
          : undefined,
      sortby: this.sortBy !== "date" ? this.getSortBy() : undefined,
      sortdir:
        this.sortDirection === SortDirectionEnum.Desc ? "desc" : undefined,
      disp:
        this.displayMode !== DEFAULT_PARAMS.displayMode
          ? this.displayMode
          : undefined,
      q: this.searchTerm,
      p:
        this.currentPage !== DEFAULT_PARAMS.currentPage
          ? this.currentPage
          : undefined,
      c: encodedCriteria,
    };
    return queryString.stringify(result, { encode: false });
  }

  // TODO: These don't support multiple of the same criteria, only the last one set is used.

  public makeFindFilter(): FindFilterType {
    return {
      q: this.searchTerm,
      page: this.currentPage,
      per_page: this.itemsPerPage,
      sort: this.getSortBy(),
      direction: this.sortDirection,
    };
  }

  public makeSceneFilter(): SceneFilterType {
    const result: SceneFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "rating": {
          const ratingCrit = criterion as RatingCriterion;
          result.rating = {
            value: ratingCrit.value,
            modifier: ratingCrit.modifier,
          };
          break;
        }
        case "o_counter": {
          const oCounterCrit = criterion as NumberCriterion;
          result.o_counter = {
            value: oCounterCrit.value,
            modifier: oCounterCrit.modifier,
          };
          break;
        }
        case "resolution": {
          switch ((criterion as ResolutionCriterion).value) {
            case "240p":
              result.resolution = ResolutionEnum.Low;
              break;
            case "480p":
              result.resolution = ResolutionEnum.Standard;
              break;
            case "720p":
              result.resolution = ResolutionEnum.StandardHd;
              break;
            case "1080p":
              result.resolution = ResolutionEnum.FullHd;
              break;
            case "4k":
              result.resolution = ResolutionEnum.FourK;
              break;
            // no default
          }
          break;
        }
        case "duration": {
          const durationCrit = criterion as DurationCriterion;
          result.duration = {
            value: durationCrit.value,
            modifier: durationCrit.modifier,
          };
          break;
        }
        case "hasMarkers":
          result.has_markers = (criterion as HasMarkersCriterion).value;
          break;
        case "sceneIsMissing":
          result.is_missing = (criterion as IsMissingCriterion).value;
          break;
        case "tags": {
          const tagsCrit = criterion as TagsCriterion;
          result.tags = {
            value: tagsCrit.value.map((tag) => tag.id),
            modifier: tagsCrit.modifier,
          };
          break;
        }
        case "performers": {
          const perfCrit = criterion as PerformersCriterion;
          result.performers = {
            value: perfCrit.value.map((perf) => perf.id),
            modifier: perfCrit.modifier,
          };
          break;
        }
        case "studios": {
          const studCrit = criterion as StudiosCriterion;
          result.studios = {
            value: studCrit.value.map((studio) => studio.id),
            modifier: studCrit.modifier,
          };
          break;
        }
        case "movies": {
          const movCrit = criterion as MoviesCriterion;
          result.movies = {
            value: movCrit.value.map((movie) => movie.id),
            modifier: movCrit.modifier,
          };
          break;
        }
        // no default
      }
    });
    return result;
  }

  public makePerformerFilter(): PerformerFilterType {
    const result: PerformerFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "favorite":
          result.filter_favorites =
            (criterion as FavoriteCriterion).value === "true";
          break;
        case "birth_year": {
          const byCrit = criterion as NumberCriterion;
          result.birth_year = {
            value: byCrit.value,
            modifier: byCrit.modifier,
          };
          break;
        }
        case "age": {
          const ageCrit = criterion as NumberCriterion;
          result.age = { value: ageCrit.value, modifier: ageCrit.modifier };
          break;
        }
        case "ethnicity": {
          const ethCrit = criterion as StringCriterion;
          result.ethnicity = {
            value: ethCrit.value,
            modifier: ethCrit.modifier,
          };
          break;
        }
        case "country": {
          const cntryCrit = criterion as StringCriterion;
          result.country = {
            value: cntryCrit.value,
            modifier: cntryCrit.modifier,
          };
          break;
        }
        case "eye_color": {
          const ecCrit = criterion as StringCriterion;
          result.eye_color = { value: ecCrit.value, modifier: ecCrit.modifier };
          break;
        }
        case "height": {
          const hCrit = criterion as StringCriterion;
          result.height = { value: hCrit.value, modifier: hCrit.modifier };
          break;
        }
        case "measurements": {
          const mCrit = criterion as StringCriterion;
          result.measurements = {
            value: mCrit.value,
            modifier: mCrit.modifier,
          };
          break;
        }
        case "fake_tits": {
          const ftCrit = criterion as StringCriterion;
          result.fake_tits = { value: ftCrit.value, modifier: ftCrit.modifier };
          break;
        }
        case "career_length": {
          const clCrit = criterion as StringCriterion;
          result.career_length = {
            value: clCrit.value,
            modifier: clCrit.modifier,
          };
          break;
        }
        case "tattoos": {
          const tCrit = criterion as StringCriterion;
          result.tattoos = { value: tCrit.value, modifier: tCrit.modifier };
          break;
        }
        case "piercings": {
          const pCrit = criterion as StringCriterion;
          result.piercings = { value: pCrit.value, modifier: pCrit.modifier };
          break;
        }
        case "aliases": {
          const aCrit = criterion as StringCriterion;
          result.aliases = { value: aCrit.value, modifier: aCrit.modifier };
          break;
        }
        case "gender": {
          const gCrit = criterion as GenderCriterion;
          result.gender = {
            value: stringToGender(gCrit.value),
            modifier: gCrit.modifier,
          };
          break;
        }
        case "performerIsMissing":
          result.is_missing = (criterion as IsMissingCriterion).value;
        // no default
      }
    });
    return result;
  }

  public makeSceneMarkerFilter(): SceneMarkerFilterType {
    const result: SceneMarkerFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "tags": {
          const tagsCrit = criterion as TagsCriterion;
          result.tags = {
            value: tagsCrit.value.map((tag) => tag.id),
            modifier: tagsCrit.modifier,
          };
          break;
        }
        case "sceneTags": {
          const sceneTagsCrit = criterion as TagsCriterion;
          result.scene_tags = {
            value: sceneTagsCrit.value.map((tag) => tag.id),
            modifier: sceneTagsCrit.modifier,
          };
          break;
        }
        case "performers": {
          const performersCrit = criterion as PerformersCriterion;
          result.performers = {
            value: performersCrit.value.map((performer) => performer.id),
            modifier: performersCrit.modifier,
          };
          break;
        }
        // no default
      }
    });
    return result;
  }

  public makeMovieFilter(): MovieFilterType {
    const result: MovieFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "studios": {
          const studCrit = criterion as StudiosCriterion;
          result.studios = {
            value: studCrit.value.map((studio) => studio.id),
            modifier: studCrit.modifier,
          };
          break;
        }
        // no default
      }
    });
    return result;
  }

  public makeStudioFilter(): StudioFilterType {
    const result: StudioFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "parent_studios": {
          const studCrit = criterion as ParentStudiosCriterion;
          result.parents = {
            value: studCrit.value.map((studio) => studio.id),
            modifier: studCrit.modifier,
          };
          break;
        }
        // no default
      }
    });

    return result;
  }

  public makeGalleryFilter(): GalleryFilterType {
    const result: GalleryFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "galleryIsMissing":
          result.is_missing = (criterion as IsMissingCriterion).value;
          break;
        // no default
      }
    });

    return result;
  }

  public makeTagFilter(): TagFilterType {
    const result: TagFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "tagIsMissing":
          result.is_missing = (criterion as IsMissingCriterion).value;
          break;
        case "scene_count": {
          const countCrit = criterion as NumberCriterion;
          result.scene_count = {
            value: countCrit.value,
            modifier: countCrit.modifier,
          };
          break;
        }
        case "marker_count": {
          const countCrit = criterion as NumberCriterion;
          result.marker_count = {
            value: countCrit.value,
            modifier: countCrit.modifier,
          };
          break;
        }
        // no default
      }
    });

    return result;
  }
}
