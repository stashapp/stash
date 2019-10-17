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
  onSelectItem: (item: ValidTypes | undefined) => void;
}

function addNoneOption(items: ValidTypes[]) {
  // Add a none option to clear the gallery
  if (!items.find((item) => item.id === "0")) { items.unshift({id: "0", name: "None"}); }
}

export const FilterSelect: React.FunctionComponent<IProps> = (props: IProps) => {
  let items: ValidTypes[];
  let InternalSelect: new (props: ISelectProps<any>) => Select<any>;
  switch (props.type) {
    case "performers": {
      const { data } = StashService.useAllPerformersForFilter();
      items = !!data && !!data.allPerformers ? data.allPerformers : [];
      addNoneOption(items);
      InternalSelect = InternalPerformerSelect;
      break;
    }
    case "studios": {
      const { data } = StashService.useAllStudiosForFilter();
      items = !!data && !!data.allStudios ? data.allStudios : [];
      addNoneOption(items);
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

  /* eslint-disable react-hooks/rules-of-hooks */
  const [selectedItem, setSelectedItem] = React.useState<ValidTypes | undefined>(undefined);

  React.useEffect(() => {
    if (!!items) {
      const initialItem = items.find((item) => props.initialId === item.id);
      if (!!initialItem) {
        setSelectedItem(initialItem);
      } else {
        setSelectedItem(undefined);
      }
    }
  }, [props.initialId, items]);
  /* eslint-enable */

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

  function onItemSelect(item: ValidTypes | undefined) {
    if (item && item.id == "0") {
      item = undefined;
    }

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
