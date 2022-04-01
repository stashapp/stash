import React, { useState } from "react";
import { useIntl } from "react-intl";
import _ from "lodash";
import Mousetrap from "mousetrap";
import { useHistory } from "react-router-dom";
import {
  FindMoviesQueryResult,
  SlimMovieDataFragment,
  MovieDataFragment,
} from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { queryFindMovies, useMoviesDestroy } from "src/core/StashService";
import {
  showWhenSelected,
  useMoviesList,
  PersistanceLevel,
} from "src/hooks/ListHook";
import { ExportDialog, DeleteEntityDialog } from "src/components/Shared";
import { MovieCard } from "./MovieCard";
import { EditMoviesDialog } from "./EditMoviesDialog";

interface IMovieList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const MovieList: React.FC<IMovieList> = ({ filterHook }) => {
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

  const addKeybinds = (
    result: FindMoviesQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  function renderEditDialog(
    selectedMovies: MovieDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <>
        <EditMoviesDialog selected={selectedMovies} onClose={onClose} />
      </>
    );
  }

  const renderDeleteDialog = (
    selectedMovies: SlimMovieDataFragment[],
    onClose: (confirmed: boolean) => void
  ) => (
    <DeleteEntityDialog
      selected={selectedMovies}
      onClose={onClose}
      singularEntity={intl.formatMessage({ id: "movie" })}
      pluralEntity={intl.formatMessage({ id: "movies" })}
      destroyMutation={useMoviesDestroy}
    />
  );

  const listData = useMoviesList({
    renderContent,
    addKeybinds,
    otherOperations,
    selectable: true,
    persistState: PersistanceLevel.ALL,
    renderEditDialog,
    renderDeleteDialog,
    filterHook,
  });

  async function viewRandom(
    result: FindMoviesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random image
    if (result.data && result.data.findMovies) {
      const { count } = result.data.findMovies;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindMovies(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findMovies &&
        singleResult.data.findMovies.movies.length === 1
      ) {
        const { id } = singleResult!.data!.findMovies!.movies[0];
        // navigate to the movie page
        history.push(`/movies/${id}`);
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

  function maybeRenderMovieExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ExportDialog
            exportInput={{
              movies: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => {
              setIsExportDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function renderContent(
    result: FindMoviesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    if (!result.data?.findMovies) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <>
          {maybeRenderMovieExportDialog(selectedIds)}
          <div className="row justify-content-center">
            {result.data.findMovies.movies.map((p) => (
              <MovieCard
                key={p.id}
                movie={p}
                selecting={selectedIds.size > 0}
                selected={selectedIds.has(p.id)}
                onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                  listData.onSelectChange(p.id, selected, shiftKey)
                }
              />
            ))}
          </div>
        </>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
