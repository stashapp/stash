import cloneDeep from "lodash-es/cloneDeep";
import React, { useState } from "react";
import { useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindPerformers,
  useFindPerformers,
  queryFindStudioPerformers,
  useFindStudioPerformers,
  usePerformersDestroy,
} from "src/core/StashService";
import {
  makeItemList,
  PersistanceLevel,
  showWhenSelected,
  IQueryParameters,
} from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { PerformerTagger } from "../Tagger/performers/PerformerTagger";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { IPerformerCardExtraCriteria, PerformerCard } from "./PerformerCard";
import { PerformerListTable } from "./PerformerListTable";
import { EditPerformersDialog } from "./EditPerformersDialog";
import { QueryResult } from "@apollo/client";

type FindStudioPerformersQueryResult = QueryResult<GQL.FindStudioPerformersQuery, Omit<GQL.FindStudioPerformersQueryVariables, 'id'> & { id?: string } >;

const StudioPerformerItemList = makeItemList({
  filterMode: GQL.FilterMode.Performers,
  useResult: useFindStudioPerformers,
  getItems(result: FindStudioPerformersQueryResult) {
    return result?.data?.findStudio?.findPerformers?.performers ?? [];
  },
  getCount(result: FindStudioPerformersQueryResult) {
    return result?.data?.findStudio?.findPerformers?.count ?? 0;
  },
});

const PerformerItemList = makeItemList({
  filterMode: GQL.FilterMode.Performers,
  useResult: useFindPerformers,
  getItems(result: GQL.FindPerformersQueryResult) {
    return result?.data?.findPerformers?.performers ?? [];
  },
  getCount(result: GQL.FindPerformersQueryResult) {
    return result?.data?.findPerformers?.count ?? 0;
  },
});

interface IPerformerList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: PersistanceLevel;
  alterQuery?: boolean;
  extraCriteria?: IPerformerCardExtraCriteria;
  queryArgs?: IQueryParameters;
}

export const PerformerList: React.FC<IPerformerList> = ({
  filterHook,
  persistState,
  alterQuery,
  extraCriteria,
  queryArgs,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.open_random" }),
      onClick: openRandom,
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

  function isFindStudioPerformersQueryResult(
    obj: QueryResult
  ): obj is FindStudioPerformersQueryResult {
    return (
      typeof obj === "object" &&
      obj != null &&
      "data" in obj &&
      "findStudio" in obj.data &&
      "findPerformers" in obj.data.findStudio
    );
  }

  function isFindPerformersQueryResult(
    obj: QueryResult
  ): obj is GQL.FindPerformersQueryResult {
    return (
      typeof obj === "object" &&
      obj !== null &&
      "data" in obj &&
      "findPerformers" in obj.data
    );
  }

  function addKeybinds(result: QueryResult, filter: ListFilterModel) {
    Mousetrap.bind("p r", () => {
      openRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }

  async function openRandom(result: QueryResult, filter: ListFilterModel) {
    const studioPerformersQuery =
      isFindStudioPerformersQueryResult(result) == true
        ? result.data.findStudio.findPerformers
        : undefined;
    const performersQuery =
      isFindPerformersQueryResult(result) == true
        ? result.data.findPerformers
        : undefined;

    if (studioPerformersQuery && queryArgs) {
      const { count } = studioPerformersQuery;
      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindStudioPerformers(
        filterCopy,
        queryArgs.id,
        queryArgs.depth
      );
      if (
        singleResult.data.findStudio?.findPerformers.performers.length === 1
      ) {
        const { id } =
          singleResult.data.findStudio?.findPerformers.performers[0]!;
        history.push(`/performers/${id}`);
      }
    } else if (performersQuery) {
      const { count } = performersQuery;
      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindPerformers(filterCopy);
      if (singleResult.data.findPerformers.performers.length === 1) {
        const { id } = singleResult.data.findPerformers.performers[0]!;
        history.push(`/performers/${id}`);
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
    result: QueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
  ) {
    function maybeRenderPerformerExportDialog() {
      if (isExportDialogOpen) {
        return (
          <>
            <ExportDialog
              exportInput={{
                performers: {
                  ids: Array.from(selectedIds.values()),
                  all: isExportAll,
                },
              }}
              onClose={() => setIsExportDialogOpen(false)}
            />
          </>
        );
      }
    }

    function renderPerformers() {
      const performers =
        isFindStudioPerformersQueryResult(result) == true
          ? result.data.findStudio.findPerformers.performers
          : isFindPerformersQueryResult(result) == true
          ? result.data.findPerformers.performers
          : undefined;

      if (!performers) {
        return;
      }

      if (filter.displayMode === DisplayMode.Grid) {
        return (
          <div className="row justify-content-center">
            {performers.map((p: GQL.Performer) => (
              <PerformerCard
                key={p.id}
                performer={p}
                selecting={selectedIds.size > 0}
                selected={selectedIds.has(p.id)}
                onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                  onSelectChange(p.id, selected, shiftKey)
                }
                extraCriteria={extraCriteria}
              />
            ))}
          </div>
        );
      }
      if (filter.displayMode === DisplayMode.List) {
        return <PerformerListTable performers={performers} />;
      }
      if (filter.displayMode === DisplayMode.Tagger) {
        return <PerformerTagger performers={performers} />;
      }
    }

    return (
      <>
        {maybeRenderPerformerExportDialog()}
        {renderPerformers()}
      </>
    );
  }

  function renderEditDialog(
    selectedPerformers: GQL.SlimPerformerDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <EditPerformersDialog selected={selectedPerformers} onClose={onClose} />
    );
  }

  function renderDeleteDialog(
    selectedPerformers: GQL.SlimPerformerDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <DeleteEntityDialog
        selected={selectedPerformers}
        onClose={onClose}
        singularEntity={intl.formatMessage({ id: "performer" })}
        pluralEntity={intl.formatMessage({ id: "performers" })}
        destroyMutation={usePerformersDestroy}
      />
    );
  }

  if (queryArgs?.type === "STUDIO") {
    return (
      <StudioPerformerItemList
        selectable
        persistState={persistState}
        alterQuery={alterQuery}
        otherOperations={otherOperations}
        addKeybinds={addKeybinds}
        renderContent={renderContent}
        renderEditDialog={renderEditDialog}
        renderDeleteDialog={renderDeleteDialog}
        queryArgs={queryArgs}
      />
    );
  } else {
    return (
      <PerformerItemList
        selectable
        filterHook={filterHook}
        persistState={persistState}
        alterQuery={alterQuery}
        otherOperations={otherOperations}
        addKeybinds={addKeybinds}
        renderContent={renderContent}
        renderEditDialog={renderEditDialog}
        renderDeleteDialog={renderDeleteDialog}
      />
    );
  }
};
