import React, { useState } from "react";
import { useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import Mousetrap from "mousetrap";
import { useHistory } from "react-router-dom";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindMovies,
  useFindMovies,
  useMoviesDestroy,
} from "src/core/StashService";
import { makeItemList, showWhenSelected } from "../List/ItemList";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { GroupCardGrid } from "./MovieCardGrid";
import { EditGroupsDialog } from "./EditMoviesDialog";
import { View } from "../List/views";

const GroupItemList = makeItemList({
  filterMode: GQL.FilterMode.Movies,
  useResult: useFindMovies,
  getItems(result: GQL.FindMoviesQueryResult) {
    return result?.data?.findMovies?.movies ?? [];
  },
  getCount(result: GQL.FindMoviesQueryResult) {
    return result?.data?.findMovies?.count ?? 0;
  },
});

interface IGroupList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
}

export const GroupList: React.FC<IGroupList> = ({
  filterHook,
  alterQuery,
  view,
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
  ];

  function addKeybinds(
    result: GQL.FindMoviesQueryResult,
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
    result: GQL.FindMoviesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random image
    if (result.data?.findMovies) {
      const { count } = result.data.findMovies;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindMovies(filterCopy);
      if (singleResult.data.findMovies.movies.length === 1) {
        const { id } = singleResult.data.findMovies.movies[0];
        // navigate to the movie page
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
    result: GQL.FindMoviesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
  ) {
    function maybeRenderGroupExportDialog() {
      if (isExportDialogOpen) {
        return (
          <ExportDialog
            exportInput={{
              movies: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => setIsExportDialogOpen(false)}
          />
        );
      }
    }

    function renderGroups() {
      if (!result.data?.findMovies) return;

      if (filter.displayMode === DisplayMode.Grid) {
        return (
          <GroupCardGrid
            groups={result.data.findMovies.movies}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
          />
        );
      }
      if (filter.displayMode === DisplayMode.List) {
        return <h1>TODO</h1>;
      }
    }
    return (
      <>
        {maybeRenderGroupExportDialog()}
        {renderGroups()}
      </>
    );
  }

  function renderEditDialog(
    selectedGroups: GQL.MovieDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return <EditGroupsDialog selected={selectedGroups} onClose={onClose} />;
  }

  function renderDeleteDialog(
    selectedGroups: GQL.SlimMovieDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <DeleteEntityDialog
        selected={selectedGroups}
        onClose={onClose}
        singularEntity={intl.formatMessage({ id: "group" })}
        pluralEntity={intl.formatMessage({ id: "groups" })}
        destroyMutation={useMoviesDestroy}
      />
    );
  }

  return (
    <GroupItemList
      selectable
      filterHook={filterHook}
      view={view}
      alterQuery={alterQuery}
      otherOperations={otherOperations}
      addKeybinds={addKeybinds}
      renderContent={renderContent}
      renderEditDialog={renderEditDialog}
      renderDeleteDialog={renderDeleteDialog}
    />
  );
};
