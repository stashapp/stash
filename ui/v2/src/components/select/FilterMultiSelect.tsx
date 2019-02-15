import * as React from "react";

import { MenuItem } from "@blueprintjs/core";
import { IMultiSelectProps, ItemPredicate, ItemRenderer, MultiSelect } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalPerformerMultiSelect = MultiSelect.ofType<GQL.AllPerformersForFilterAllPerformers>();
const InternalTagMultiSelect = MultiSelect.ofType<GQL.AllTagsForFilterAllTags>();
const InternalStudioMultiSelect = MultiSelect.ofType<GQL.AllStudiosForFilterAllStudios>();

type ValidTypes =
  GQL.AllPerformersForFilterAllPerformers |
  GQL.AllTagsForFilterAllTags |
  GQL.AllStudiosForFilterAllStudios;

interface IProps extends HTMLInputProps, Partial<IMultiSelectProps<ValidTypes>> {
  type: "performers" | "studios" | "tags";
  initialIds?: string[];
  onUpdate: (items: ValidTypes[]) => void;
}

export const FilterMultiSelect: React.FunctionComponent<IProps> = (props: IProps) => {
  let items: ValidTypes[];
  let InternalMultiSelect: new (props: IMultiSelectProps<any>) => MultiSelect<any>;
  switch (props.type) {
    case "performers": {
      const { data } = StashService.useAllPerformersForFilter();
      items = !!data && !!data.allPerformers ? data.allPerformers : [];
      InternalMultiSelect = InternalPerformerMultiSelect;
      break;
    }
    case "studios": {
      const { data } = StashService.useAllStudiosForFilter();
      items = !!data && !!data.allStudios ? data.allStudios : [];
      InternalMultiSelect = InternalStudioMultiSelect;
      break;
    }
    case "tags": {
      const { data } = StashService.useAllTagsForFilter();
      items = !!data && !!data.allTags ? data.allTags : [];
      InternalMultiSelect = InternalTagMultiSelect;
      break;
    }
    default: {
      console.error("Unhandled case in FilterMultiSelect");
      return <>Unhandled case in FilterMultiSelect</>;
    }
  }

  const [selectedItems, setSelectedItems] = React.useState<ValidTypes[]>([]);
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);

  if (!!props.initialIds && selectedItems.length === 0 && !isInitialized) {
    const initialItems = items.filter((item) => props.initialIds!.includes(item.id));
    if (initialItems.length > 0) {
      setSelectedItems(initialItems);
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
      />
    );
  };

  const filter: ItemPredicate<ValidTypes> = (query, item) => {
    if (selectedItems.includes(item)) { return false; }
    return item.name!.toLowerCase().indexOf(query.toLowerCase()) >= 0;
  };

  function onItemSelect(item: ValidTypes) {
    selectedItems.push(item);
    setSelectedItems(selectedItems);
    props.onUpdate(selectedItems);
  }

  function onItemRemove(value: string, index: number) {
    const newSelectedItems = selectedItems.filter((_, i) => i !== index);
    setSelectedItems(newSelectedItems);
    props.onUpdate(newSelectedItems);
  }

  return (
    <InternalMultiSelect
      items={items}
      selectedItems={selectedItems}
      itemRenderer={renderItem}
      itemPredicate={filter}
      tagRenderer={(tag) => tag.name}
      tagInputProps={{ onRemove: onItemRemove }}
      onItemSelect={onItemSelect}
      resetOnSelect={true}
      activeItem={null}
      popoverProps={{position: "bottom"}}
      {...props}
    />
  );
};
