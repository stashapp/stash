import * as React from "react";

import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer, Suggest } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalSuggest = Suggest.ofType<GQL.ScrapePerformerListScrapePerformerList>();

interface IProps extends HTMLInputProps {
  scraperId: string;
  onSelectPerformer: (query: GQL.ScrapePerformerListScrapePerformerList) => void;
}

export const ScrapePerformerSuggest: React.FunctionComponent<IProps> = (props: IProps) => {
  const [query, setQuery] = React.useState<string>("");
  const [selectedItem, setSelectedItem] = React.useState<GQL.ScrapePerformerListScrapePerformerList | undefined>();
  const [debouncedQuery, setDebouncedQuery] = React.useState<string>("");
  const { data } = StashService.useScrapePerformerList(props.scraperId, debouncedQuery);

  React.useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedQuery(query);
    }, 500);

    return () => {
      clearTimeout(handler);
    };
  }, [query])

  const performerNames = !!data && !!data.scrapePerformerList ? data.scrapePerformerList : [];

  const renderInputValue = (performer: GQL.ScrapePerformerListScrapePerformerList) => performer.name || "";

  const renderItem: ItemRenderer<GQL.ScrapePerformerListScrapePerformerList> = (performer, itemProps) => {
    if (!itemProps.modifiers.matchesPredicate) { return null; }
    return (
      <MenuItem
        active={itemProps.modifiers.active}
        disabled={itemProps.modifiers.disabled}
        key={performer.name}
        onClick={itemProps.handleClick}
        text={performer.name}
      />
    );
  };

  return (
    <InternalSuggest
      inputValueRenderer={renderInputValue}
      items={performerNames}
      itemRenderer={renderItem}
      onItemSelect={(item) => { props.onSelectPerformer(item); setSelectedItem(item); }}
      onQueryChange={(newQuery) => { setQuery(newQuery); }}
      activeItem={null}
      selectedItem={selectedItem}
      popoverProps={{position: "bottom"}}
    />
  );
};
