import React, { useState } from "react";
import { useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindStudios,
  useFindStudios,
  useStudiosDestroy,
} from "src/core/StashService";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { StudioTagger } from "../Tagger/studios/StudioTagger";
import { StudioCardGrid } from "./StudioCardGrid";
import { View } from "../List/views";

function getItems(result: GQL.FindStudiosQueryResult) {
  return result?.data?.findStudios?.studios ?? [];
}

function getCount(result: GQL.FindStudiosQueryResult) {
  return result?.data?.findStudios?.count ?? 0;
}

interface IStudioList {
  fromParent?: boolean;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
}

export const StudioList: React.FC<IStudioList> = ({
  fromParent,
  filterHook,
  view,
  alterQuery,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const filterMode = GQL.FilterMode.Studios;

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
    result: GQL.FindStudiosQueryResult,
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
    result: GQL.FindStudiosQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random studio
    if (result.data?.findStudios) {
      const { count } = result.data.findStudios;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindStudios(filterCopy);
      if (singleResult.data.findStudios.studios.length === 1) {
        const { id } = singleResult.data.findStudios.studios[0];
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

  function renderContent(
    result: GQL.FindStudiosQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
  ) {
    function maybeRenderExportDialog() {
      if (isExportDialogOpen) {
        return (
          <ExportDialog
            exportInput={{
              studios: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => setIsExportDialogOpen(false)}
          />
        );
      }
    }

    function renderStudios() {
      if (!result.data?.findStudios) return;

      if (filter.displayMode === DisplayMode.Grid) {
        return (
          <StudioCardGrid
            studios={result.data.findStudios.studios}
            fromParent={fromParent}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
          />
        );
      }
      if (filter.displayMode === DisplayMode.List) {
        return <h1>TODO</h1>;
      }
      if (filter.displayMode === DisplayMode.Wall) {
        return <h1>TODO</h1>;
      }
      if (filter.displayMode === DisplayMode.Tagger) {
        return <StudioTagger studios={result.data.findStudios.studios} />;
      }
    }

    return (
      <>
        {maybeRenderExportDialog()}
        {renderStudios()}
      </>
    );
  }

  function renderDeleteDialog(
    selectedStudios: GQL.SlimStudioDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <DeleteEntityDialog
        selected={selectedStudios}
        onClose={onClose}
        singularEntity={intl.formatMessage({ id: "studio" })}
        pluralEntity={intl.formatMessage({ id: "studios" })}
        destroyMutation={useStudiosDestroy}
      />
    );
  }

  return (
    <ItemListContext
      filterMode={filterMode}
      useResult={useFindStudios}
      getItems={getItems}
      getCount={getCount}
      alterQuery={alterQuery}
      filterHook={filterHook}
      view={view}
      selectable
    >
      <ItemList
        view={view}
        otherOperations={otherOperations}
        addKeybinds={addKeybinds}
        renderContent={renderContent}
        renderDeleteDialog={renderDeleteDialog}
      />
    </ItemListContext>
  );
};
