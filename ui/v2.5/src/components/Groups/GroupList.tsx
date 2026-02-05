import React, { PropsWithChildren, useState } from "react";
import { useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import Mousetrap from "mousetrap";
import { useHistory } from "react-router-dom";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindGroups,
  useFindGroups,
  useGroupsDestroy,
} from "src/core/StashService";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { GroupCardGrid } from "./GroupCardGrid";
import { EditGroupsDialog } from "./EditGroupsDialog";
import { View } from "../List/views";
import {
  IFilteredListToolbar,
  IItemListOperation,
} from "../List/FilteredListToolbar";
import { PatchComponent } from "src/patch";

const GroupExportDialog: React.FC<{
  open?: boolean;
  selectedIds: Set<string>;
  isExportAll?: boolean;
  onClose: () => void;
}> = ({ open = false, selectedIds, isExportAll = false, onClose }) => {
  if (!open) {
    return null;
  }

  return (
    <ExportDialog
      exportInput={{
        groups: {
          ids: Array.from(selectedIds.values()),
          all: isExportAll,
        },
      }}
      onClose={onClose}
    />
  );
};

const filterMode = GQL.FilterMode.Groups;

function getItems(result: GQL.FindGroupsQueryResult) {
  return result?.data?.findGroups?.groups ?? [];
}

function getCount(result: GQL.FindGroupsQueryResult) {
  return result?.data?.findGroups?.count ?? 0;
}

interface IGroupListContext {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultFilter?: ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  selectable?: boolean;
}

export const GroupListContext: React.FC<
  PropsWithChildren<IGroupListContext>
> = ({ alterQuery, filterHook, defaultFilter, view, selectable, children }) => {
  return (
    <ItemListContext
      filterMode={filterMode}
      defaultFilter={defaultFilter}
      useResult={useFindGroups}
      getItems={getItems}
      getCount={getCount}
      alterQuery={alterQuery}
      filterHook={filterHook}
      view={view}
      selectable={selectable}
    >
      {children}
    </ItemListContext>
  );
};

interface IGroupList extends IGroupListContext {
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
  renderToolbar?: (props: IFilteredListToolbar) => React.ReactNode;
  otherOperations?: IItemListOperation<GQL.FindGroupsQueryResult>[];
}

export const GroupList: React.FC<IGroupList> = PatchComponent(
  "GroupList",
  ({
    filterHook,
    alterQuery,
    defaultFilter,
    view,
    fromGroupId,
    onMove,
    selectable,
    renderToolbar,
    otherOperations: providedOperations = [],
  }) => {
    const intl = useIntl();
    const history = useHistory();
    const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
    const [isExportAll, setIsExportAll] = useState(false);

    const otherOperations = [
      {
        text: intl.formatMessage({ id: "actions.view_random" }),
        onClick: viewRandom,
      },
      {
        text: intl.formatMessage({ id: "actions.export" }),
        onClick: onExport,
        isDisplayed: showWhenSelected,
      },
      {
        text: intl.formatMessage({ id: "actions.export_all" }),
        onClick: onExportAll,
      },
      ...providedOperations,
    ];

    function addKeybinds(
      result: GQL.FindGroupsQueryResult,
      filter: ListFilterModel
    ) {
      Mousetrap.bind("p r", () => {
        viewRandom(result, filter);
      });

      return () => {
        Mousetrap.unbind("p r");
      };
    }

    async function viewRandom(
      result: GQL.FindGroupsQueryResult,
      filter: ListFilterModel
    ) {
      // query for a random image
      if (result.data?.findGroups) {
        const { count } = result.data.findGroups;

        const index = Math.floor(Math.random() * count);
        const filterCopy = cloneDeep(filter);
        filterCopy.itemsPerPage = 1;
        filterCopy.currentPage = index + 1;
        const singleResult = await queryFindGroups(filterCopy);
        if (singleResult.data.findGroups.groups.length === 1) {
          const { id } = singleResult.data.findGroups.groups[0];
          // navigate to the group page
          history.push(`/groups/${id}`);
        }
      }
    }

    async function onExport() {
      setIsExportAll(false);
      setIsExportDialogOpen(true);
    }

    async function onExportAll() {
      setIsExportAll(true);
      setIsExportDialogOpen(true);
    }

    function renderContent(
      result: GQL.FindGroupsQueryResult,
      filter: ListFilterModel,
      selectedIds: Set<string>,
      onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
    ) {
      return (
        <>
          <GroupExportDialog
            open={isExportDialogOpen}
            selectedIds={selectedIds}
            isExportAll={isExportAll}
            onClose={() => setIsExportDialogOpen(false)}
          />
          {filter.displayMode === DisplayMode.Grid && (
            <GroupCardGrid
              groups={result.data?.findGroups.groups ?? []}
              zoomIndex={filter.zoomIndex}
              selectedIds={selectedIds}
              onSelectChange={onSelectChange}
              fromGroupId={fromGroupId}
              onMove={onMove}
            />
          )}
        </>
      );
    }

    function renderEditDialog(
      selectedGroups: GQL.ListGroupDataFragment[],
      onClose: (applied: boolean) => void
    ) {
      return <EditGroupsDialog selected={selectedGroups} onClose={onClose} />;
    }

    function renderDeleteDialog(
      selectedGroups: GQL.SlimGroupDataFragment[],
      onClose: (confirmed: boolean) => void
    ) {
      return (
        <DeleteEntityDialog
          selected={selectedGroups}
          onClose={onClose}
          singularEntity={intl.formatMessage({ id: "group" })}
          pluralEntity={intl.formatMessage({ id: "groups" })}
          destroyMutation={useGroupsDestroy}
        />
      );
    }

    return (
      <GroupListContext
        alterQuery={alterQuery}
        filterHook={filterHook}
        view={view}
        defaultFilter={defaultFilter}
        selectable={selectable}
      >
        <ItemList
          view={view}
          otherOperations={otherOperations}
          addKeybinds={addKeybinds}
          renderContent={renderContent}
          renderEditDialog={renderEditDialog}
          renderDeleteDialog={renderDeleteDialog}
          renderToolbar={renderToolbar}
        />
      </GroupListContext>
    );
  }
);
