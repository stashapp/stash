import { IntlShape } from "react-intl";
import { ITypename } from "src/utils/data";
import { ImageWallOptions } from "src/utils/imageWall";
import { RatingSystemOptions } from "src/utils/rating";
import { FilterMode, SortDirectionEnum } from "./generated-graphql";

// NOTE: double capitals aren't converted correctly in the backend

export interface ISavedFilterRow extends ITypename {
  __typename: "SavedFilter";
  savedFilterId: number;
}

export interface IMessage {
  id: string;
  values: { [key: string]: string };
}

export interface ICustomFilter extends ITypename {
  __typename: "CustomFilter";
  message?: IMessage;
  title?: string;
  mode: FilterMode;
  sortBy: string;
  direction: SortDirectionEnum;
}

// NOTE: This value cannot be more defined, because the generated enum it depends upon is UpperCase, which leads to errors on saving
export type PinnedFilters = Record<string, Array<string>>;

export type FrontPageContent = ISavedFilterRow | ICustomFilter;

export const defaultMaxOptionsShown = 200;

export interface IUIConfig {
  frontPageContent?: FrontPageContent[];

  showChildTagContent?: boolean;
  showChildStudioContent?: boolean;
  showTagCardOnHover?: boolean;

  abbreviateCounters?: boolean;

  ratingSystemOptions?: RatingSystemOptions;

  // if true a background image will be display on header
  enableMovieBackgroundImage?: boolean;
  // if true a background image will be display on header
  enablePerformerBackgroundImage?: boolean;
  // if true a background image will be display on header
  enableStudioBackgroundImage?: boolean;
  // if true a background image will be display on header
  enableTagBackgroundImage?: boolean;
  // if true view expanded details compact
  compactExpandedDetails?: boolean;
  // if true show all content details by default
  showAllDetails?: boolean;
  // if true the chromecast option will enabled
  enableChromecast?: boolean;
  // if true continue scene will always play from the beginning
  alwaysStartFromBeginning?: boolean;
  // if true enable activity tracking
  trackActivity?: boolean;
  // the minimum percentage of scene duration which a scene must be played
  // before the play count is incremented
  minimumPlayPercent?: number;

  // maximum number of items to shown in the dropdown list - defaults to 200
  // upper limit of 1000
  maxOptionsShown?: number;

  imageWallOptions?: ImageWallOptions;

  lastNoteSeen?: number;

  vrTag?: string;
  pinnedFilters?: PinnedFilters;
}

function recentlyReleased(
  intl: IntlShape,
  mode: FilterMode,
  objectsID: string
): ICustomFilter {
  return {
    __typename: "CustomFilter",
    message: {
      id: "recently_released_objects",
      values: { objects: intl.formatMessage({ id: objectsID }) },
    },
    mode,
    sortBy: "date",
    direction: SortDirectionEnum.Desc,
  };
}

function recentlyAdded(
  intl: IntlShape,
  mode: FilterMode,
  objectsID: string
): ICustomFilter {
  return {
    __typename: "CustomFilter",
    message: {
      id: "recently_added_objects",
      values: { objects: intl.formatMessage({ id: objectsID }) },
    },
    mode,
    sortBy: "created_at",
    direction: SortDirectionEnum.Desc,
  };
}

export function generateDefaultFrontPageContent(intl: IntlShape) {
  return [
    recentlyReleased(intl, FilterMode.Scenes, "scenes"),
    recentlyAdded(intl, FilterMode.Studios, "studios"),
    recentlyReleased(intl, FilterMode.Movies, "movies"),
    recentlyAdded(intl, FilterMode.Performers, "performers"),
    recentlyReleased(intl, FilterMode.Galleries, "galleries"),
  ];
}

export function generatePremadeFrontPageContent(intl: IntlShape) {
  return [
    recentlyReleased(intl, FilterMode.Scenes, "scenes"),
    recentlyAdded(intl, FilterMode.Scenes, "scenes"),
    recentlyReleased(intl, FilterMode.Galleries, "galleries"),
    recentlyAdded(intl, FilterMode.Galleries, "galleries"),
    recentlyAdded(intl, FilterMode.Images, "images"),
    recentlyReleased(intl, FilterMode.Movies, "movies"),
    recentlyAdded(intl, FilterMode.Movies, "movies"),
    recentlyAdded(intl, FilterMode.Studios, "studios"),
    recentlyAdded(intl, FilterMode.Performers, "performers"),
  ];
}
