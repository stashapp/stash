import queryString from "query-string";
import {
  FindFilterType,
  PerformerFilterType,
  ResolutionEnum,
  SceneFilterType,
  SceneMarkerFilterType,
  SortDirectionEnum,
} from "../../core/generated-graphql";
import { Criterion, ICriterionOption, CriterionType, CriterionOption, NumberCriterion, StringCriterion } from "./criteria/criterion";
import { FavoriteCriterion, FavoriteCriterionOption } from "./criteria/favorite";
import { HasMarkersCriterion, HasMarkersCriterionOption } from "./criteria/has-markers";
import { IsMissingCriterion, IsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { PerformersCriterion, PerformersCriterionOption } from "./criteria/performers";
import { RatingCriterion, RatingCriterionOption } from "./criteria/rating";
import { ResolutionCriterion, ResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterion, StudiosCriterionOption } from "./criteria/studios";
import { SceneTagsCriterionOption, TagsCriterion, TagsCriterionOption } from "./criteria/tags";
import { makeCriteria } from "./criteria/utils";
import {
  DisplayMode,
  FilterMode,
} from "./types";

interface IQueryParameters {
  sortby?: string;
  sortdir?: string;
  disp?: string;
  q?: string;
  p?: string;
  c?: string[];
}

// TODO: handle customCriteria
export class ListFilterModel {
  public filterMode: FilterMode = FilterMode.Scenes;
  public searchTerm?: string;
  public currentPage = 1;
  public itemsPerPage = 40;
  public sortDirection: "asc" | "desc" = "asc";
  public sortBy?: string;
  public sortByOptions: string[] = [];
  public displayMode: DisplayMode = DisplayMode.Grid;
  public displayModeOptions: DisplayMode[] = [];
  public criterionOptions: ICriterionOption[] = [];
  public criteria: Array<Criterion<any, any>> = [];
  public totalCount: number = 0;

  public constructor(filterMode: FilterMode) {
    switch (filterMode) {
      case FilterMode.Scenes:
        if (!!this.sortBy === false) { this.sortBy = "date"; }
        this.sortByOptions = ["title", "path", "rating", "date", "filesize", "duration", "framerate", "bitrate", "random"];
        this.displayModeOptions = [
          DisplayMode.Grid,
          DisplayMode.List,
          DisplayMode.Wall,
        ];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new RatingCriterionOption(),
          new ResolutionCriterionOption(),
          new HasMarkersCriterionOption(),
          new IsMissingCriterionOption(),
          new TagsCriterionOption(),
          new PerformersCriterionOption(),
          new StudiosCriterionOption(),
        ];
        break;
      case FilterMode.Performers:
        if (!!this.sortBy === false) { this.sortBy = "name"; }
        this.sortByOptions = ["name", "height", "birthdate", "scenes_count"];
        this.displayModeOptions = [
          DisplayMode.Grid,
          DisplayMode.List,
        ];

        var numberCriteria : CriterionType[] = ["birth_year", "age"];
        var stringCriteria : CriterionType[] = [
          "ethnicity",
          "country",
          "eye_color",
          "height",
          "measurements",
          "fake_tits",
          "career_length",
          "tattoos",
          "piercings",
          "aliases"
        ];

        this.criterionOptions = [
          new NoneCriterionOption(),
          new FavoriteCriterionOption()
        ];

        this.criterionOptions = this.criterionOptions.concat(numberCriteria.concat(stringCriteria).map((c) => {
          return new CriterionOption(Criterion.getLabel(c), c);
        }));
        break;
      case FilterMode.Studios:
        if (!!this.sortBy === false) { this.sortBy = "name"; }
        this.sortByOptions = ["name", "scenes_count"];
        this.displayModeOptions = [
          DisplayMode.Grid,
        ];
        this.criterionOptions = [
          new NoneCriterionOption(),
        ];
        break;
      case FilterMode.Galleries:
        if (!!this.sortBy === false) { this.sortBy = "path"; }
        this.sortByOptions = ["path"];
        this.displayModeOptions = [
          DisplayMode.List,
        ];
        this.criterionOptions = [
          new NoneCriterionOption(),
        ];
        break;
      case FilterMode.SceneMarkers:
        if (!!this.sortBy === false) { this.sortBy = "title"; }
        this.sortByOptions = ["title", "seconds", "scene_id", "random", "scenes_updated_at"];
        this.displayModeOptions = [
          DisplayMode.Wall,
        ];
        this.criterionOptions = [
          new NoneCriterionOption(),
          new TagsCriterionOption(),
          new SceneTagsCriterionOption(),
          new PerformersCriterionOption(),
        ];
        break;
      default:
        this.sortByOptions = [];
        this.displayModeOptions = [];
        this.criterionOptions = [
          new NoneCriterionOption(),
        ];
        break;
    }
    if (!!this.displayMode === false) { this.displayMode = this.displayModeOptions[0]; }
    this.sortByOptions = [...this.sortByOptions, "created_at", "updated_at"];
  }

  public configureFromQueryParameters(rawParms: any) {
    const params = rawParms as IQueryParameters;
    if (params.sortby !== undefined) {
      this.sortBy = params.sortby;
    }
    if (params.sortdir === "asc" || params.sortdir === "desc") {
      this.sortDirection = params.sortdir;
    }
    if (params.disp !== undefined) {
      this.displayMode = parseInt(params.disp, 10);
    }
    if (params.q !== undefined) {
      this.searchTerm = params.q;
    }
    if (params.p !== undefined) {
      this.currentPage = Number(params.p);
    }

    if (params.c !== undefined) {
      this.criteria = [];

      let jsonParameters: any[];
      if (params.c instanceof Array) {
        jsonParameters = params.c;
      } else {
        jsonParameters = [params.c];
      }

      for (const jsonString of jsonParameters) {
        const encodedCriterion = JSON.parse(jsonString);
        const criterion = makeCriteria(encodedCriterion.type);
        criterion.value = encodedCriterion.value;
        criterion.modifier = encodedCriterion.modifier;
        this.criteria.push(criterion);
      }
    }
  }

  public makeQueryParameters(): string {
    const encodedCriteria: string[] = [];
    this.criteria.forEach((criterion) => {
      const encodedCriterion: any = {};
      encodedCriterion.type = criterion.type;
      encodedCriterion.value = criterion.value;
      encodedCriterion.modifier = criterion.modifier;
      const jsonCriterion = JSON.stringify(encodedCriterion);
      encodedCriteria.push(jsonCriterion);
    });

    const result = {
      sortby: this.sortBy,
      sortdir: this.sortDirection,
      disp: this.displayMode,
      q: this.searchTerm,
      p: this.currentPage,
      c: encodedCriteria,
    };
    return queryString.stringify(result, {encode: false});
  }

  // TODO: These don't support multiple of the same criteria, only the last one set is used.

  public makeFindFilter(): FindFilterType {
    return {
      q: this.searchTerm,
      page: this.currentPage,
      per_page: this.itemsPerPage,
      sort: this.sortBy,
      direction: this.sortDirection === "asc" ? SortDirectionEnum.Asc : SortDirectionEnum.Desc,
    };
  }

  public makeSceneFilter(): SceneFilterType {
    const result: SceneFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "rating":
          const ratingCrit = criterion as RatingCriterion;
          result.rating = { value: ratingCrit.value, modifier: ratingCrit.modifier };
          break;
        case "resolution": {
          switch ((criterion as ResolutionCriterion).value) {
            case "240p": result.resolution = ResolutionEnum.Low; break;
            case "480p": result.resolution = ResolutionEnum.Standard; break;
            case "720p": result.resolution = ResolutionEnum.StandardHd; break;
            case "1080p": result.resolution = ResolutionEnum.FullHd; break;
            case "4k": result.resolution = ResolutionEnum.FourK; break;
          }
          break;
        }
        case "hasMarkers":
          result.has_markers = (criterion as HasMarkersCriterion).value;
          break;
        case "isMissing":
          result.is_missing = (criterion as IsMissingCriterion).value;
          break;
        case "tags":
          const tagsCrit = criterion as TagsCriterion;
          result.tags = { value: tagsCrit.value.map((tag) => tag.id), modifier: tagsCrit.modifier };
          break;
        case "performers":
          const perfCrit = criterion as PerformersCriterion;
          result.performers = { value: perfCrit.value.map((perf) => perf.id), modifier: perfCrit.modifier };
          break;
        case "studios":
          const studCrit = criterion as StudiosCriterion;
          result.studios = { value: studCrit.value.map((studio) => studio.id), modifier: studCrit.modifier };
          break;
      }
    });
    return result;
  }

  public makePerformerFilter(): PerformerFilterType {
    const result: PerformerFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "favorite":
          result.filter_favorites = (criterion as FavoriteCriterion).value === "true";
          break;
        case "birth_year":
          const byCrit = criterion as NumberCriterion;
          result.birth_year = { value: byCrit.value, modifier: byCrit.modifier };
          break;
        case "age":
          const ageCrit = criterion as NumberCriterion;
          result.age = { value: ageCrit.value, modifier: ageCrit.modifier };
          break;
        case "ethnicity":
          const ethCrit = criterion as StringCriterion;
          result.ethnicity = { value: ethCrit.value, modifier: ethCrit.modifier };
          break;
        case "country":
          const cntryCrit = criterion as StringCriterion;
          result.country = { value: cntryCrit.value, modifier: cntryCrit.modifier };
          break;
        case "eye_color":
          const ecCrit = criterion as StringCriterion;
          result.eye_color = { value: ecCrit.value, modifier: ecCrit.modifier };
          break;
        case "height":
          const hCrit = criterion as StringCriterion;
          result.height = { value: hCrit.value, modifier: hCrit.modifier };
          break;
        case "measurements":
          const mCrit = criterion as StringCriterion;
          result.measurements = { value: mCrit.value, modifier: mCrit.modifier };
          break;
        case "fake_tits":
          const ftCrit = criterion as StringCriterion;
          result.fake_tits = { value: ftCrit.value, modifier: ftCrit.modifier };
          break;
        case "career_length":
          const clCrit = criterion as StringCriterion;
          result.career_length = { value: clCrit.value, modifier: clCrit.modifier };
          break;
        case "tattoos":
          const tCrit = criterion as StringCriterion;
          result.tattoos = { value: tCrit.value, modifier: tCrit.modifier };
          break;
        case "piercings":
          const pCrit = criterion as StringCriterion;
          result.piercings = { value: pCrit.value, modifier: pCrit.modifier };
          break;
        case "aliases":
          const aCrit = criterion as StringCriterion;
          result.aliases = { value: aCrit.value, modifier: aCrit.modifier };
          break;
      }
    });
    return result;
  }

  public makeSceneMarkerFilter(): SceneMarkerFilterType {
    const result: SceneMarkerFilterType = {};
    this.criteria.forEach((criterion) => {
      switch (criterion.type) {
        case "tags":
          const tagsCrit = criterion as TagsCriterion;
          result.tags = { value: tagsCrit.value.map((tag) => tag.id), modifier: tagsCrit.modifier };
          break;
        case "sceneTags":
          const sceneTagsCrit = criterion as TagsCriterion;
          result.scene_tags = { value: sceneTagsCrit.value.map((tag) => tag.id), modifier: sceneTagsCrit.modifier };
          break;
        case "performers":
          const performersCrit = criterion as PerformersCriterion;
          result.performers = { value: performersCrit.value.map((performer) => performer.id), modifier: performersCrit.modifier };
          break;
      }
    });
    return result;
  }
}
