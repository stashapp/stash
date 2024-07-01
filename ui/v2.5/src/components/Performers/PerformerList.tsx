import cloneDeep from "lodash-es/cloneDeep";
import React, { useState } from "react";
import { useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindPerformers,
  useFindPerformers,
  usePerformersDestroy,
} from "src/core/StashService";
import { makeItemList, showWhenSelected } from "../List/ItemList";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { PerformerTagger } from "../Tagger/performers/PerformerTagger";
import { ExportDialog } from "../Shared/ExportDialog";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { IPerformerCardExtraCriteria } from "./PerformerCard";
import { PerformerListTable } from "./PerformerListTable";
import { EditPerformersDialog } from "./EditPerformersDialog";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";
import TextUtils from "src/utils/text";
import { PerformerCardGrid } from "./PerformerCardGrid";
import { View } from "../List/views";

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

export const FormatHeight = (height?: number | null) => {
  const intl = useIntl();
  if (!height) {
    return "";
  }

  const [feet, inches] = cmToImperial(height);

  return (
    <span className="performer-height">
      <span className="height-metric">
        {intl.formatNumber(height, {
          style: "unit",
          unit: "centimeter",
          unitDisplay: "short",
        })}
      </span>
      <span className="height-imperial">
        {intl.formatNumber(feet, {
          style: "unit",
          unit: "foot",
          unitDisplay: "narrow",
        })}
        {intl.formatNumber(inches, {
          style: "unit",
          unit: "inch",
          unitDisplay: "narrow",
        })}
      </span>
    </span>
  );
};

export const FormatAge = (
  birthdate?: string | null,
  deathdate?: string | null
) => {
  if (!birthdate) {
    return "";
  }
  const age = TextUtils.age(birthdate, deathdate);

  return (
    <span className="performer-age">
      <span className="age">{age}</span>
      <span className="birthdate"> ({birthdate})</span>
    </span>
  );
};

export const FormatWeight = (weight?: number | null) => {
  const intl = useIntl();
  if (!weight) {
    return "";
  }

  const lbs = kgToLbs(weight);

  return (
    <span className="performer-weight">
      <span className="weight-metric">
        {intl.formatNumber(weight, {
          style: "unit",
          unit: "kilogram",
          unitDisplay: "short",
        })}
      </span>
      <span className="weight-imperial">
        {intl.formatNumber(lbs, {
          style: "unit",
          unit: "pound",
          unitDisplay: "short",
        })}
      </span>
    </span>
  );
};

export const FormatCircumcised = (circumcised?: GQL.CircumisedEnum | null) => {
  const intl = useIntl();
  if (!circumcised) {
    return "";
  }

  return (
    <span className="penis-circumcised">
      {intl.formatMessage({
        id: "circumcised_types." + circumcised,
      })}
    </span>
  );
};

export const FormatPenisLength = (penis_length?: number | null) => {
  const intl = useIntl();
  if (!penis_length) {
    return "";
  }

  const inches = cmToInches(penis_length);

  return (
    <span className="performer-penis-length">
      <span className="penis-length-metric">
        {intl.formatNumber(penis_length, {
          style: "unit",
          unit: "centimeter",
          unitDisplay: "short",
          maximumFractionDigits: 2,
        })}
      </span>
      <span className="penis-length-imperial">
        {intl.formatNumber(inches, {
          style: "unit",
          unit: "inch",
          unitDisplay: "narrow",
          maximumFractionDigits: 2,
        })}
      </span>
    </span>
  );
};

interface IPerformerList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const PerformerList: React.FC<IPerformerList> = ({
  filterHook,
  view,
  alterQuery,
  extraCriteria,
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

  function addKeybinds(
    result: GQL.FindPerformersQueryResult,
    filter: ListFilterModel
  ) {
    Mousetrap.bind("p r", () => {
      openRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }

  async function openRandom(
    result: GQL.FindPerformersQueryResult,
    filter: ListFilterModel
  ) {
    if (result.data?.findPerformers) {
      const { count } = result.data.findPerformers;
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
    result: GQL.FindPerformersQueryResult,
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
      if (!result.data?.findPerformers) return;

      if (filter.displayMode === DisplayMode.Grid) {
        return (
          <PerformerCardGrid
            performers={result.data.findPerformers.performers}
            zoomIndex={filter.zoomIndex}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
            extraCriteria={extraCriteria}
          />
        );
      }
      if (filter.displayMode === DisplayMode.List) {
        return (
          <PerformerListTable
            performers={result.data.findPerformers.performers}
            selectedIds={selectedIds}
            onSelectChange={onSelectChange}
          />
        );
      }
      if (filter.displayMode === DisplayMode.Tagger) {
        return (
          <PerformerTagger performers={result.data.findPerformers.performers} />
        );
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

  return (
    <PerformerItemList
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
