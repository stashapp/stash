import { ITypename } from "src/utils";
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

export type FrontPageContent = ISavedFilterRow | ICustomFilter;

export interface IUIConfig {
  frontPageContent?: FrontPageContent[];
}

function recentlyReleased(mode: FilterMode, objects: string): ICustomFilter {
  return {
    __typename: "CustomFilter",
    message: { id: "recently_released_objects", values: { objects } },
    mode,
    sortBy: "date",
    direction: SortDirectionEnum.Desc,
  };
}

function recentlyAdded(mode: FilterMode, objects: string): ICustomFilter {
  return {
    __typename: "CustomFilter",
    message: { id: "recently_added_objects", values: { objects } },
    mode,
    sortBy: "created_at",
    direction: SortDirectionEnum.Desc,
  };
}

export function generateDefaultFrontPageContent() {
  return [
    recentlyReleased(FilterMode.Scenes, "scenes"),
    recentlyAdded(FilterMode.Studios, "studios"),
    recentlyReleased(FilterMode.Movies, "movies"),
    recentlyAdded(FilterMode.Performers, "performers"),
    recentlyReleased(FilterMode.Galleries, "galleries"),
  ];
}

export function generatePremadeFrontPageContent() {
  return [
    recentlyReleased(FilterMode.Scenes, "scenes"),
    recentlyAdded(FilterMode.Scenes, "scenes"),
    recentlyReleased(FilterMode.Galleries, "galleries"),
    recentlyAdded(FilterMode.Galleries, "galleries"),
    recentlyAdded(FilterMode.Images, "images"),
    recentlyReleased(FilterMode.Movies, "movies"),
    recentlyAdded(FilterMode.Movies, "movies"),
    recentlyAdded(FilterMode.Studios, "studios"),
    recentlyAdded(FilterMode.Performers, "performers"),
  ];
}
