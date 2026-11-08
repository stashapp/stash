import { faBookOpen } from "@fortawesome/free-solid-svg-icons";
import type {
  CriterionModifier,
  FilterMode,
  SortDirectionEnum,
  SavedFilterDataFragment,
} from "../../core/generated-graphql";
import {
  buildMainNavbarMenuItems,
  buildSavedFilterMenuItems,
  type INavbarMenuItem,
} from "./menuItems";
import type { ISavedFilterMenuItem } from "../../core/config";

const mangaFilter: SavedFilterDataFragment = {
  __typename: "SavedFilter",
  id: "7",
  mode: "GALLERIES" as FilterMode,
  name: "Manga",
  find_filter: {
    __typename: "SavedFindFilterType",
    q: "",
    page: 1,
    per_page: 40,
    sort: "date",
    direction: "ASC" as SortDirectionEnum,
  },
  object_filter: {
    tags: {
      value: ["12"],
      modifier: "INCLUDES_ALL" as CriterionModifier,
    },
  },
  ui_options: {
    display_mode: 2,
    zoom_index: 0,
  },
};

const configuredItems: ISavedFilterMenuItem[] = [
  {
    id: "manga-section",
    label: "Manga",
    filterId: "7",
    mode: "GALLERIES" as FilterMode,
  },
];

const [item] = buildSavedFilterMenuItems(configuredItems, [mangaFilter]);

assertEqual(item.name, "saved-filter:manga-section");
assertEqual(item.label, "Manga");
assertEqual(item.message, undefined);
assertEqual(item.icon, faBookOpen);
assertEqual(item.hotkey, undefined);
assertMatch(item.href, /^\/galleries\?/);
assertMatch(item.href, /sortby=date/);
assertMatch(item.href, /sortdir=asc/);
assertMatch(item.href, /z=0/);
assertMatch(item.href, /c=/);
assertEqual(buildSavedFilterMenuItems(configuredItems, []).length, 0);

const builtInItems: INavbarMenuItem[] = [
  {
    name: "galleries",
    label: "Galleries",
    href: "/galleries",
    icon: faBookOpen,
  },
];

const combinedItems = buildMainNavbarMenuItems(
  builtInItems,
  undefined,
  configuredItems,
  [mangaFilter]
);

assertEqual(combinedItems.length, 2);
assertEqual(combinedItems[0].name, "galleries");
assertEqual(combinedItems[1].name, "saved-filter:manga-section");
function assertEqual<T>(actual: T, expected: T) {
  if (actual !== expected) {
    throw new Error(`Expected ${String(expected)}, got ${String(actual)}`);
  }
}

function assertMatch(actual: string, expected: RegExp) {
  if (!expected.test(actual)) {
    throw new Error(`Expected ${actual} to match ${expected.toString()}`);
  }
}
