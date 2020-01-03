import * as React from "react";

import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer, Suggest } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalSuggest = Suggest.ofType<GQL.MarkerStringsMarkerStrings>();

interface IProps extends HTMLInputProps {
  initialMarkerString?: string;
  onQueryChange: (query: string) => void;
}

export const MarkerTitleSuggest: React.FunctionComponent<IProps> = (props: IProps) => {
  const { data } = StashService.useMarkerStrings();
  const markerStrings = !!data && !!data.markerStrings ? data.markerStrings : [];
  const [selectedItem, setSelectedItem] = React.useState<GQL.MarkerStringsMarkerStrings | null>(null);

  if (!!props.initialMarkerString && !selectedItem) {
    const initialItem = markerStrings.find((item) => {
      return props.initialMarkerString!.toLowerCase() === item!.title.toLowerCase();
    });
    if (!!initialItem) { setSelectedItem(initialItem); }
  }

  const renderInputValue = (markerString: GQL.MarkerStringsMarkerStrings) => markerString.title;

  const renderItem: ItemRenderer<GQL.MarkerStringsMarkerStrings> = (markerString, itemProps) => {
    if (!itemProps.modifiers.matchesPredicate) { return null; }
    return (
      <MenuItem
        active={itemProps.modifiers.active}
        disabled={itemProps.modifiers.disabled}
        label={markerString.count.toString()}
        key={markerString.id}
        onClick={itemProps.handleClick}
        text={markerString.title}
      />
    );
  };

  const filter: ItemPredicate<GQL.MarkerStringsMarkerStrings> = (query, item) => {
    return item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0;
  };

  return (
    <InternalSuggest
      inputValueRenderer={renderInputValue}
      items={markerStrings as any}
      itemRenderer={renderItem}
      itemPredicate={filter}
      onItemSelect={(item) => { props.onQueryChange(item.title); setSelectedItem(item); }}
      onQueryChange={(query) => { props.onQueryChange(query); setSelectedItem(null); }}
      activeItem={null}
      selectedItem={selectedItem}
      popoverProps={{position: "bottom"}}
    />
  );
};
