import * as React from "react";

import { Button, MenuItem } from "@blueprintjs/core";
import { ISelectProps, ItemPredicate, ItemRenderer, Select } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalPerformerSelect = Select.ofType<GQL.AllPerformersForFilterAllPerformers>();
const InternalTagSelect = Select.ofType<GQL.AllTagsForFilterAllTags>();
const InternalStudioSelect = Select.ofType<GQL.AllStudiosForFilterAllStudios>();

type ValidTypes =
  GQL.AllPerformersForFilterAllPerformers |
  GQL.AllTagsForFilterAllTags |
  GQL.AllStudiosForFilterAllStudios;

interface IProps extends HTMLInputProps {
  type: "performers" | "studios" | "tags";
  initialId?: string;
  onSelectItem: (item: ValidTypes) => void;
}

export const FilterSelect: React.FunctionComponent<IProps> = (props: IProps) => {
  let items: ValidTypes[];
  let InternalSelect: new (props: ISelectProps<any>) => Select<any>;
  switch (props.type) {
    case "performers": {
      const { data } = StashService.useAllPerformersForFilter();
      items = !!data && !!data.allPerformers ? data.allPerformers : [];
      InternalSelect = InternalPerformerSelect;
      break;
    }
    case "studios": {
      const { data } = StashService.useAllStudiosForFilter();
      items = !!data && !!data.allStudios ? data.allStudios : [];
      InternalSelect = InternalStudioSelect;
      break;
    }
    case "tags": {
      const { data } = StashService.useAllTagsForFilter();
      items = !!data && !!data.allTags ? data.allTags : [];
      InternalSelect = InternalTagSelect;
      break;
    }
    default: {
      console.error("Unhandled case in FilterSelect");
      return <>Unhandled case in FilterSelect</>;
    }
  }

  const [selectedItem, setSelectedItem] = React.useState<ValidTypes | null>(null);
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);

  if (!!props.initialId && !selectedItem && !isInitialized) {
    const initialItem = items.find((item) => props.initialId === item.id);
    if (!!initialItem) {
      setSelectedItem(initialItem);
      setIsInitialized(true);
    }
  }

  const renderItem: ItemRenderer<ValidTypes> = (item, itemProps) => {
    if (!itemProps.modifiers.matchesPredicate) { return null; }
    return (
      <MenuItem
        active={itemProps.modifiers.active}
        disabled={itemProps.modifiers.disabled}
        key={item.id}
        onClick={itemProps.handleClick}
        text={item.name}
        shouldDismissPopover={false}
      />
    );
  };

  const filter: ItemPredicate<ValidTypes> = (query, item) => {
    return item.name!.toLowerCase().indexOf(query.toLowerCase()) >= 0;
  };

  function onItemSelect(item: ValidTypes) {
    props.onSelectItem(item);
    setSelectedItem(item);
  }

  const buttonText = selectedItem ? selectedItem.name : "(No selection)";
  return (
    <InternalSelect
      items={items}
      itemRenderer={renderItem}
      itemPredicate={filter}
      noResults={<MenuItem disabled={true} text="No results." />}
      onItemSelect={onItemSelect}
      popoverProps={{position: "bottom"}}
      {...props}
    >
      <Button fill={true} text={buttonText} />
    </InternalSelect>
  );
};
