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
  let MultiSelectImpl = getMultiSelectImpl();
  let InternalMultiSelect = MultiSelectImpl.getInternalMultiSelect();
  const data = MultiSelectImpl.getData();
  
  const [selectedItems, setSelectedItems] = React.useState<ValidTypes[]>([]);
  const [items, setItems] = React.useState<ValidTypes[]>([]);
  const [newTagName, setNewTagName] = React.useState<string>("");
  const createTag = StashService.useTagCreate(getTagInput() as GQL.TagCreateInput);

  React.useEffect(() => {
    if (!!data) {
      MultiSelectImpl.translateData();
    }
  }, [data]);
      
  function getTagInput() {
    const tagInput: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = { name: newTagName };
    return tagInput;
  }

  async function onCreateNewObject(item: ValidTypes) {
    var created : any;
    if (props.type === "tags") {
      try {
        created = await createTag();
        
        items.push(created.data.tagCreate);
        setItems(items.slice());
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

  React.useEffect(() => {
    if (!!props.initialIds && !!items) {
      const initialItems = items.filter((item) => props.initialIds!.includes(item.id));
      setSelectedItems(initialItems);
    }
  }, [props.initialIds, items]);

  function getMultiSelectImpl() {
    let getInternalMultiSelect: () => new (props: IMultiSelectProps<any>) => MultiSelect<any>;
    let getData: () => GQL.AllPerformersForFilterQuery | GQL.AllStudiosForFilterQuery | GQL.AllTagsForFilterQuery | undefined;
    let translateData: () => void;
    let createNewObject: ((query : string) => void) | undefined = undefined; 

    switch (props.type) {
      case "performers": {
        getInternalMultiSelect = () => { return InternalPerformerMultiSelect; };
        getData = () => { const { data } = StashService.useAllPerformersForFilter(); return data; }
        translateData = () => { let perfData = data as GQL.AllPerformersForFilterQuery; setItems(!!perfData && !!perfData.allPerformers ? perfData.allPerformers : []); };
        break;
      }
      case "studios": {
        getInternalMultiSelect = () => { return InternalStudioMultiSelect; };
        getData = () => { const { data } = StashService.useAllStudiosForFilter(); return data; }
        translateData = () => { let studioData = data as GQL.AllStudiosForFilterQuery; setItems(!!studioData && !!studioData.allStudios ? studioData.allStudios : []); };
        break;
      }
      case "tags": {
        getInternalMultiSelect = () => { return InternalTagMultiSelect; };
        getData = () => { const { data } = StashService.useAllTagsForFilter(); return data; }
        translateData = () => { let tagData = data as GQL.AllTagsForFilterQuery; setItems(!!tagData && !!tagData.allTags ? tagData.allTags : []); };
        createNewObject = createNewTag;
        break;
      }
      default: {
        throw "Unhandled case in FilterMultiSelect";
      }
    }

    return {
      getInternalMultiSelect: getInternalMultiSelect,
      getData: getData,
      translateData: translateData,
      createNewObject: createNewObject
    };
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
      createNewItemFromQuery={MultiSelectImpl.createNewObject}
      createNewItemRenderer={createNewRenderer}
      {...props}
    />
  );
};
