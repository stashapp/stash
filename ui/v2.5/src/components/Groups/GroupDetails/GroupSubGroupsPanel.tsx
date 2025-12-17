import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupList } from "../GroupList";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  ContainingGroupsCriterionOption,
  GroupsCriterion,
} from "src/models/list-filter/criteria/groups";
import {
  useRemoveSubGroups,
  useReorderSubGroupsMutation,
} from "src/core/StashService";
import { ButtonToolbar } from "react-bootstrap";
import { ListOperationButtons } from "src/components/List/ListOperationButtons";
import { useListContext } from "src/components/List/ListProvider";
import {
  PageSizeSelector,
  SearchTermInput,
} from "src/components/List/ListFilter";
import { useFilter } from "src/components/List/FilterProvider";
import {
  IFilteredListToolbar,
  IItemListOperation,
} from "src/components/List/FilteredListToolbar";
import {
  showWhenNoneSelected,
  showWhenSelected,
} from "src/components/List/ItemList";
import { faMinus, faPlus } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { useModal } from "src/hooks/modal";
import { AddSubGroupsDialog } from "./AddGroupsDialog";
import { PatchComponent } from "src/patch";

const useContainingGroupFilterHook = (
  group: Pick<GQL.StudioDataFragment, "id" | "name">,
  showSubGroupContent?: boolean
) => {
  return (filter: ListFilterModel) => {
    const groupValue = { id: group.id, label: group.name };
    // if studio is already present, then we modify it, otherwise add
    let groupCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "containing_groups";
    }) as GroupsCriterion | undefined;

    if (groupCriterion) {
      // add the group if not present
      if (
        !groupCriterion.value.items.find((p) => {
          return p.id === group.id;
        })
      ) {
        groupCriterion.value.items.push(groupValue);
      }
    } else {
      groupCriterion = new GroupsCriterion(ContainingGroupsCriterionOption);
      groupCriterion.value = {
        items: [groupValue],
        excluded: [],
        depth: showSubGroupContent ? -1 : 0,
      };
      groupCriterion.modifier = GQL.CriterionModifier.Includes;
      filter.criteria.push(groupCriterion);
    }

    filter.sortBy = "sub_group_order";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;

    return filter;
  };
};

const Toolbar: React.FC<IFilteredListToolbar> = ({
  onEdit,
  onDelete,
  operations,
}) => {
  const { getSelected, onSelectAll, onSelectNone } = useListContext();
  const { filter, setFilter } = useFilter();

  return (
    <ButtonToolbar className="filtered-list-toolbar">
      <div>
        <SearchTermInput filter={filter} onFilterUpdate={setFilter} />
      </div>
      <PageSizeSelector
        pageSize={filter.itemsPerPage}
        setPageSize={(size) => setFilter(filter.setPageSize(size))}
      />
      <ListOperationButtons
        onSelectAll={onSelectAll}
        onSelectNone={onSelectNone}
        itemsSelected={getSelected().length > 0}
        otherOperations={operations}
        onEdit={onEdit}
        onDelete={onDelete}
      />
    </ButtonToolbar>
  );
};

interface IGroupSubGroupsPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
  extraOperations?: IItemListOperation<GQL.FindGroupsQueryResult>[];
}

const defaultFilter = (() => {
  const sortBy = "sub_group_order";
  const ret = new ListFilterModel(GQL.FilterMode.Groups, undefined, {
    defaultSortBy: sortBy,
  });

  // unset the sort by so that its not included in the URL
  ret.sortBy = undefined;

  return ret;
})();

export const GroupSubGroupsPanel: React.FC<IGroupSubGroupsPanel> =
  PatchComponent(
    "GroupSubGroupsPanel",
    ({ active, group, extraOperations = [] }) => {
      const intl = useIntl();
      const Toast = useToast();
      const { modal, showModal, closeModal } = useModal();

      const [reorderSubGroups] = useReorderSubGroupsMutation();
      const mutateRemoveSubGroups = useRemoveSubGroups();

      const filterHook = useContainingGroupFilterHook(group);

      async function removeSubGroups(
        result: GQL.FindGroupsQueryResult,
        filter: ListFilterModel,
        selectedIds: Set<string>
      ) {
        try {
          await mutateRemoveSubGroups(
            group.id,
            Array.from(selectedIds.values())
          );

          Toast.success(
            intl.formatMessage(
              { id: "toast.removed_entity" },
              {
                count: selectedIds.size,
                singularEntity: intl.formatMessage({ id: "group" }),
                pluralEntity: intl.formatMessage({ id: "groups" }),
              }
            )
          );
        } catch (e) {
          Toast.error(e);
        }
      }

      async function onAddSubGroups() {
        showModal(
          <AddSubGroupsDialog containingGroup={group} onClose={closeModal} />
        );
      }

      const otherOperations = [
        ...extraOperations,
        {
          text: intl.formatMessage({ id: "actions.add_sub_groups" }),
          onClick: onAddSubGroups,
          isDisplayed: showWhenNoneSelected,
          postRefetch: true,
          icon: faPlus,
          buttonVariant: "secondary",
        },
        {
          text: intl.formatMessage({
            id: "actions.remove_from_containing_group",
          }),
          onClick: removeSubGroups,
          isDisplayed: showWhenSelected,
          postRefetch: true,
          icon: faMinus,
          buttonVariant: "danger",
        },
      ];

      function onMove(srcIds: string[], targetId: string, after: boolean) {
        reorderSubGroups({
          variables: {
            input: {
              group_id: group.id,
              sub_group_ids: srcIds,
              insert_at_id: targetId,
              insert_after: after,
            },
          },
        });
      }

      return (
        <>
          {modal}
          <GroupList
            defaultFilter={defaultFilter}
            filterHook={filterHook}
            alterQuery={active}
            fromGroupId={group.id}
            otherOperations={otherOperations}
            onMove={onMove}
            renderToolbar={(props) => <Toolbar {...props} />}
          />
        </>
      );
    }
  );
