import * as React from "react";

import { MenuItem } from "@blueprintjs/core";
import { IMultiSelectProps, ItemPredicate, ItemRenderer, MultiSelect } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";

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
  var createNewFunc = undefined;
  
  const [newTagName, setNewTagName] = React.useState<string>("");
  const createTag = StashService.useTagCreate(getTagInput() as GQL.TagCreateInput);

  function getTagInput() {
    const tagInput: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = { name: newTagName };
    return tagInput;
  }

  async function onCreateNewObject(item: ValidTypes) {
    var created : any;
    if (props.type === "tags") {
      try {
        created = await createTag();
        
        addSelectedItem(created.data.tagCreate);

        ToastUtils.success("Created tag");
      } catch (e) {
        ErrorUtils.handle(e);
      }
    }
  }

  function createNewTag(query : string) {
    setNewTagName(query);
    return {
      name : query
    };
  }

  function createNewRenderer(query: string, active: boolean, handleClick: React.MouseEventHandler<HTMLElement>) {
    // if tag already exists with that name, then don't return anything
    if (items.find((item) => {
      return item.name === query;
    })) {
      return undefined;
    }

    return (
      <MenuItem
        icon="add"
        text={`Create "${query}"`}
        active={active}
        onClick={handleClick}
      />
    );
  }

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
      createNewFunc = createNewTag;
      break;
    }
    default: {
      console.error("Unhandled case in FilterMultiSelect");
      return <>Unhandled case in FilterMultiSelect</>;
    }
  }

  /* eslint-disable react-hooks/rules-of-hooks */
  const [selectedItems, setSelectedItems] = React.useState<ValidTypes[]>([]);
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);
  /* eslint-enable */

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

  function addSelectedItem(item: ValidTypes) {
    selectedItems.push(item);
    setSelectedItems(selectedItems);
    props.onUpdate(selectedItems);
  }

  function onItemSelect(item: ValidTypes) {
    if (item.id === undefined) {
      // create the new item, if applicable
      onCreateNewObject(item);
    } else {
      addSelectedItem(item);
    }
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
      popoverProps={{position: "bottom"}}
      createNewItemFromQuery={createNewFunc}
      createNewItemRenderer={createNewRenderer}
      {...props}
    />
  );
};
