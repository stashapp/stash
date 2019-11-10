import * as React from "react";

import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer, Suggest } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalSuggest = Suggest.ofType<string>();

interface IProps extends HTMLInputProps {
  scraperId: string;
  onQueryChange: (query: string) => void;
}

export const FreeOnesPerformerSuggest: React.FunctionComponent<IProps> = (props: IProps) => {
  const [query, setQuery] = React.useState<string>("");
  const { data } = StashService.useScrapePerformerList(props.scraperId, query);
  const performerNames = !!data && !!data.scrapePerformerList ? data.scrapePerformerList : [];

  const renderInputValue = (performerName: string) => performerName;

  const renderItem: ItemRenderer<string> = (performerName, itemProps) => {
    if (!itemProps.modifiers.matchesPredicate) { return null; }
    return (
      <MenuItem
        active={itemProps.modifiers.active}
        disabled={itemProps.modifiers.disabled}
        key={performerName}
        onClick={itemProps.handleClick}
        text={performerName}
      />
    );
  };

  return (
    <InternalSuggest
      inputValueRenderer={renderInputValue}
      items={performerNames}
      itemRenderer={renderItem}
      onItemSelect={(item) => { props.onQueryChange(item); setQuery(item); }}
      onQueryChange={(newQuery) => { props.onQueryChange(newQuery); setQuery(newQuery); }}
      activeItem={null}
      selectedItem={query}
      popoverProps={{position: "bottom"}}
    />
  );
};
