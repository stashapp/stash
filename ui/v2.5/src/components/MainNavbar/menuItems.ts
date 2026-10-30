import type { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { faBookOpen } from "@fortawesome/free-solid-svg-icons";
import type { MessageDescriptor } from "react-intl";
import { FilterMode } from "../../core/generated-graphql";
import type {
  ConfigDataFragment,
  SavedFilterDataFragment,
} from "../../core/generated-graphql";
import type { ISavedFilterMenuItem } from "../../core/config";
import { ListFilterModel } from "../../models/list-filter/filter";

export interface INavbarMenuItem {
  name: string;
  href: string;
  icon: IconDefinition;
  message?: MessageDescriptor;
  label?: string;
  hotkey?: string;
  userCreatable?: boolean;
}

function makeSavedFilterHref(
  savedFilter: SavedFilterDataFragment,
  config?: ConfigDataFragment
) {
  const pathname =
    savedFilter.mode === FilterMode.Galleries ? "/galleries" : undefined;
  if (!pathname) {
    return;
  }

  const filter = new ListFilterModel(savedFilter.mode, config);
  filter.configureFromSavedFilter(savedFilter);
  const query = filter.makeQueryParameters();

  return query ? `${pathname}?${query}` : pathname;
}

export function buildSavedFilterMenuItems(
  configuredItems: ISavedFilterMenuItem[] | undefined,
  savedFilters: SavedFilterDataFragment[],
  config?: ConfigDataFragment
): INavbarMenuItem[] {
  if (!configuredItems?.length) {
    return [];
  }

  return configuredItems.flatMap((item) => {
    const savedFilter = savedFilters.find(
      (filter) => filter.id === item.filterId && filter.mode === item.mode
    );
    if (!savedFilter) {
      return [];
    }

    const href = makeSavedFilterHref(savedFilter, config);
    if (!href) {
      return [];
    }

    return [
      {
        name: `saved-filter:${item.id}`,
        label: item.label || savedFilter.name,
        href,
        icon: faBookOpen,
      },
    ];
  });
}

export function buildMainNavbarMenuItems(
  builtInItems: INavbarMenuItem[],
  configuredBuiltInNames: string[] | null | undefined,
  configuredSavedFilterItems: ISavedFilterMenuItem[] | undefined,
  savedFilters: SavedFilterDataFragment[],
  config?: ConfigDataFragment
) {
  const builtInNames = configuredBuiltInNames?.map((item) =>
    item === "movies" ? "groups" : item
  );
  const enabledBuiltInItems = builtInNames
    ? builtInItems.filter((menuItem) => builtInNames.includes(menuItem.name))
    : builtInItems;

  return [
    ...enabledBuiltInItems,
    ...buildSavedFilterMenuItems(
      configuredSavedFilterItems,
      savedFilters,
      config
    ),
  ];
}
