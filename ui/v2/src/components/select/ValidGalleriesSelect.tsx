import * as React from "react";

import { Button, MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer, Select } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";

const InternalSelect = Select.ofType<GQL.ValidGalleriesForSceneValidGalleriesForScene>();

interface IProps extends HTMLInputProps {
  initialId?: string;
  sceneId: string;
  onSelectItem: (item: GQL.ValidGalleriesForSceneValidGalleriesForScene | undefined) => void;
}

export const ValidGalleriesSelect: React.FunctionComponent<IProps> = (props: IProps) => {
  const { data } = StashService.useValidGalleriesForScene(props.sceneId);
  const items = !!data && !!data.validGalleriesForScene ? data.validGalleriesForScene : [];
  // Add a none option to clear the gallery
  if (!items.find((item) => item.id === "0")) { items.unshift({id: "0", path: "None"}); }

  const [selectedItem, setSelectedItem] = React.useState<GQL.ValidGalleriesForSceneValidGalleriesForScene | undefined>(undefined);
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);

  if (!!props.initialId && !selectedItem && !isInitialized) {
    const initialItem = items.find((item) => props.initialId === item.id);
    if (!!initialItem) {
      setSelectedItem(initialItem);
      setIsInitialized(true);
    }
  }

  const renderItem: ItemRenderer<GQL.ValidGalleriesForSceneValidGalleriesForScene> = (item, itemProps) => {
    if (!itemProps.modifiers.matchesPredicate) { return null; }
    return (
      <MenuItem
        active={itemProps.modifiers.active}
        disabled={itemProps.modifiers.disabled}
        key={item.id}
        onClick={itemProps.handleClick}
        text={item.path}
        shouldDismissPopover={false}
      />
    );
  };

  const filter: ItemPredicate<GQL.ValidGalleriesForSceneValidGalleriesForScene> = (query, item) => {
    return item.path!.toLowerCase().indexOf(query.toLowerCase()) >= 0;
  };

  function onItemSelect(item: GQL.ValidGalleriesForSceneValidGalleriesForScene | undefined) {
    if (item && item.id === "0") {
      item = undefined;
    }

    props.onSelectItem(item);
    setSelectedItem(item);
  }

  const buttonText = selectedItem ? selectedItem.path : "(No selection)";
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
