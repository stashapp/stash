import React, { useState } from "react";
import { useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import {
  FindStudiosQueryResult,
  SlimStudioDataFragment,
} from "src/core/generated-graphql";
import { useStudiosList } from "src/hooks";
import { showWhenSelected, PersistanceLevel } from "src/hooks/ListHook";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { queryFindStudios, useStudiosDestroy } from "src/core/StashService";
import { ExportDialog, DeleteEntityDialog } from "src/components/Shared";
import { StudioCard } from "./StudioCard";

interface IStudioList {
  fromParent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const StudioList: React.FC<IStudioList> = ({
  fromParent,
  filterHook,
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

  const addKeybinds = (
    result: FindStudiosQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  async function viewRandom(
    result: FindStudiosQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random studio
    if (result.data && result.data.findStudios) {
      const { count } = result.data.findStudios;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindStudios(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findStudios &&
        singleResult.data.findStudios.studios.length === 1
      ) {
        const { id } = singleResult!.data!.findStudios!.studios[0];
        // navigate to the studio page
        history.push(`/studios/${id}`);
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

  function maybeRenderExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ExportDialog
            exportInput={{
              studios: {
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

  const renderDeleteDialog = (
    selectedStudios: SlimStudioDataFragment[],
    onClose: (confirmed: boolean) => void
  ) => (
    <DeleteEntityDialog
      selected={selectedStudios}
      onClose={onClose}
      singularEntity={intl.formatMessage({ id: "studio" })}
      pluralEntity={intl.formatMessage({ id: "studios" })}
      destroyMutation={useStudiosDestroy}
    />
  );

  const listData = useStudiosList({
    renderContent,
    filterHook,
    addKeybinds,
    otherOperations,
    selectable: true,
    persistState: !fromParent ? PersistanceLevel.ALL : PersistanceLevel.NONE,
    renderDeleteDialog,
  });

  function renderStudios(
    result: FindStudiosQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    if (!result.data?.findStudios) return;

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row px-xl-5 justify-content-center">
          {result.data.findStudios.studios.map((studio) => (
            <StudioCard
              key={studio.id}
              studio={studio}
              hideParent={fromParent}
              selecting={selectedIds.size > 0}
              selected={selectedIds.has(studio.id)}
              onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                listData.onSelectChange(studio.id, selected, shiftKey)
              }
            />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  function renderContent(
    result: FindStudiosQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    return (
      <>
        {maybeRenderExportDialog(selectedIds)}
        {renderStudios(result, filter, selectedIds)}
      </>
    );
  }

  return listData.template;
};
