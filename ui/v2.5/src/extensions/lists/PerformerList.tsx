import React, { useCallback, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindPerformers,
  useFindPerformers,
  usePerformersDestroy,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { PerformerTagger } from "src/components/Tagger/performers/PerformerTagger";
import { ExportDialog } from "src/components/Shared/ExportDialog";
import { DeleteEntityDialog } from "src/components/Shared/DeleteEntityDialog";
import { IPerformerCardExtraCriteria } from "src/components/Performers/PerformerCard";
import { PerformerListTable } from "src/components/Performers/PerformerListTable";
import { EditPerformersDialog } from "src/components/Performers/EditPerformersDialog";
import { cmToImperial, cmToInches, kgToLbs } from "src/utils/units";
import TextUtils from "src/utils/text";
import { PerformerCardGrid } from "src/components/Performers/PerformerCardGrid";
import { View } from "src/components/List/views";
import { LoadedContent } from "src/components/List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "src/components/List/util";
import {
  OperationDropdown,
  OperationDropdownItem,
} from "src/components/List/ListOperationButtons";
import { useFilteredItemList } from "src/components/List/ItemList";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "src/components/Shared/Sidebar";
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import {
  SidebarRatingFilter,
  SidebarGenderFilter,
  SidebarTagsFilter,
  SidebarBooleanFilter,
  SidebarStashIDFilter,
  SidebarCountryFilter,
  SidebarStringFilter,
  SidebarNumberFilter,
  SidebarStudiosFilter,
  SidebarGroupsFilter,
  SidebarDateFilter,
  SidebarAgeFilter,
  SidebarFilterSelector,
  FilterWrapper,
  SidebarCircumcisedFilter,
  SidebarIsMissingFilter,
} from "src/extensions/filters";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "src/extensions/ui";
import { PatchContainerComponent } from "src/patch";
import { Pagination } from "src/components/List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "src/extensions/ui";
import { ListResultsHeader } from "src/extensions/ui";
import {
  createDateCriterionOption,
  createNumberCriterionOption,
  createStringCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { useFocus } from "src/extensions/hooks";
import {
  faFilter,
  faPencil,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { TaggerContext } from "src/components/Tagger/context";
import { GenderCriterionOption } from "src/models/list-filter/criteria/gender";
import { TagsCriterionOption } from "src/models/list-filter/criteria/tags";
import { FavoritePerformerCriterionOption } from "src/models/list-filter/criteria/favorite";
import { StashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { CountryCriterionOption } from "src/models/list-filter/criteria/country";
import { TattoosCriterionOption } from "src/models/list-filter/criteria/tattoos";
import { PiercingsCriterionOption } from "src/models/list-filter/criteria/piercings";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { SidebarFilterDefinition } from "src/extensions/hooks/useSidebarFilters";
import {
  createMandatoryTimestampCriterionOption,
  createMandatoryNumberCriterionOption,
  createBooleanCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { CircumcisedCriterionOption } from "src/models/list-filter/criteria/circumcised";
import { GroupsCriterionOption } from "src/models/list-filter/criteria/groups";
import { PerformerIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import {
  usePerformerFacetCounts,
  FacetCountsContext,
} from "src/extensions/hooks/useFacetCounts";

function getItems(result: GQL.FindPerformersQueryResult) {
  return result?.data?.findPerformers?.performers ?? [];
}

function getCount(result: GQL.FindPerformersQueryResult) {
  return result?.data?.findPerformers?.count ?? 0;
}

function useOpenRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const openRandom = useCallback(async () => {
    // query for a random performer
    if (count === 0) {
      return;
    }

    const pages = Math.ceil(count / filter.itemsPerPage);
    const page = Math.floor(Math.random() * pages) + 1;

    const indexMax = Math.min(filter.itemsPerPage, count);
    const index = Math.floor(Math.random() * indexMax);
    const filterCopy = cloneDeep(filter);
    filterCopy.currentPage = page;
    filterCopy.sortBy = "random";
    const queryResults = await queryFindPerformers(filterCopy);
    const performer = queryResults.data.findPerformers.performers[index];
    if (performer) {
      history.push(`/performers/${performer.id}`);
    }
  }, [filter, count, history]);

  return openRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const openRandom = useOpenRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      openRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [openRandom]);
}

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

const PerformerListContent: React.FC<{
  performers: GQL.SlimPerformerDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}> = ({ performers, filter, selectedIds, onSelectChange, extraCriteria }) => {
  if (performers.length === 0) {
    return null;
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <PerformerCardGrid
        performers={performers as any}
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
        performers={performers as any}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Tagger) {
    return <PerformerTagger performers={performers as any} />;
  }

  return null;
};

export const MyPerformersFilterSidebarSections = PatchContainerComponent(
  "MyFilteredPerformerList.SidebarSections"
);

// Define available filters for performers sidebar
const performerFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "rating", messageId: "rating", defaultVisible: true },
  { id: "favorite", messageId: "favourite", defaultVisible: true },
  { id: "tags", messageId: "tags", defaultVisible: true },
  { id: "studios", messageId: "studios", defaultVisible: true },
  { id: "groups", messageId: "groups", defaultVisible: false },
  { id: "gender", messageId: "gender", defaultVisible: true },
  { id: "age", messageId: "age", defaultVisible: true },
  { id: "country", messageId: "country", defaultVisible: true },

  // Tier 2: Physical attributes
  { id: "ethnicity", messageId: "ethnicity", defaultVisible: false },
  { id: "circumcised", messageId: "circumcised", defaultVisible: false },
  { id: "hair_color", messageId: "hair_color", defaultVisible: false },
  { id: "eye_color", messageId: "eye_color", defaultVisible: false },
  { id: "height_cm", messageId: "height", defaultVisible: false },
  { id: "weight", messageId: "weight", defaultVisible: false },
  { id: "penis_length", messageId: "penis_length", defaultVisible: false },
  { id: "measurements", messageId: "measurements", defaultVisible: false },
  { id: "fake_tits", messageId: "fake_tits", defaultVisible: false },
  { id: "tattoos", messageId: "tattoos", defaultVisible: false },
  { id: "piercings", messageId: "piercings", defaultVisible: false },

  // Tier 3: Dates
  { id: "birthdate", messageId: "birthdate", defaultVisible: false },
  { id: "death_date", messageId: "death_date", defaultVisible: false },
  { id: "birth_year", messageId: "birth_year", defaultVisible: false },
  { id: "death_year", messageId: "death_year", defaultVisible: false },
  { id: "career_length", messageId: "career_length", defaultVisible: false },

  // Tier 4: Library stats
  { id: "scene_count", messageId: "scene_count", defaultVisible: false },
  { id: "image_count", messageId: "image_count", defaultVisible: false },
  { id: "gallery_count", messageId: "gallery_count", defaultVisible: false },
  { id: "tag_count", messageId: "tag_count", defaultVisible: false },
  { id: "play_count", messageId: "play_count", defaultVisible: false },
  { id: "o_counter", messageId: "o_count", defaultVisible: false },

  // Tier 5: Metadata
  { id: "name", messageId: "name", defaultVisible: false },
  { id: "aliases", messageId: "aliases", defaultVisible: false },
  { id: "disambiguation", messageId: "disambiguation", defaultVisible: false },
  { id: "details", messageId: "details", defaultVisible: false },
  { id: "url", messageId: "url", defaultVisible: false },

  // Tier 6: System
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "ignore_auto_tag", messageId: "ignore_auto_tag", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
  { id: "stash_id", messageId: "stash_id", defaultVisible: false },
];

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  sidebarOpen: boolean;
  onClose?: () => void;
  showEditFilter: (editingCriterion?: string) => void;
  count?: number;
  focus?: ReturnType<typeof useFocus>;
  clearAllCriteria: () => void;
  onFilterEditModeChange?: (isEditMode: boolean) => void;
}> = ({
  filter,
  setFilter,
  filterHook,
  view,
  showEditFilter,
  sidebarOpen,
  onClose,
  count,
  focus,
  clearAllCriteria,
  onFilterEditModeChange,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  // Criterion options
  const UrlCriterionOption = createStringCriterionOption("url");
  const AgeCriterionOption = createNumberCriterionOption("age");
  const DeathDateCriterionOption = createDateCriterionOption("death_date");
  const BirthdateCriterionOption = createDateCriterionOption("birthdate");
  const BirthYearCriterionOption = createNumberCriterionOption("birth_year");
  const DeathYearCriterionOption = createNumberCriterionOption("death_year");
  const CareerLengthCriterionOption = createStringCriterionOption("career_length");
  const SceneCountCriterionOption = createNumberCriterionOption("scene_count");
  const ImageCountCriterionOption = createNumberCriterionOption("image_count");
  const GalleryCountCriterionOption = createNumberCriterionOption("gallery_count");
  const TagCountCriterionOption = createMandatoryNumberCriterionOption("tag_count");
  const PlayCountCriterionOption = createMandatoryNumberCriterionOption("play_count");
  const OCounterCriterionOption = createNumberCriterionOption("o_counter", "o_count");
  const NameCriterionOption = createStringCriterionOption("name");
  const AliasesCriterionOption = createStringCriterionOption("aliases");
  const DisambiguationCriterionOption = createStringCriterionOption("disambiguation");
  const DetailsCriterionOption = createStringCriterionOption("details");
  const EthnicityCriterionOption = createStringCriterionOption("ethnicity");
  const HairColorCriterionOption = createStringCriterionOption("hair_color");
  const EyeColorCriterionOption = createStringCriterionOption("eye_color");
  const HeightCriterionOption = createNumberCriterionOption("height_cm", "height");
  const WeightCriterionOption = createNumberCriterionOption("weight");
  const PenisLengthCriterionOption = createNumberCriterionOption("penis_length");
  const MeasurementsCriterionOption = createStringCriterionOption("measurements");
  const FakeTitsCriterionOption = createStringCriterionOption("fake_tits");
  const IgnoreAutoTagCriterionOption = createBooleanCriterionOption("ignore_auto_tag");
  const CreatedAtCriterionOption = createMandatoryTimestampCriterionOption("created_at");
  const UpdatedAtCriterionOption = createMandatoryTimestampCriterionOption("updated_at");

  return (
    <>
      <FilteredSidebarHeader
        sidebarOpen={sidebarOpen}
        showEditFilter={showEditFilter}
        filter={filter}
        setFilter={setFilter}
        view={view}
        focus={focus}
      />
      <MyPerformersFilterSidebarSections>
        <div className="sidebar-filters">
          <SidebarFilterSelector
            viewName="performers"
            filterDefinitions={performerFilterDefinitions}
            headerContent={
              <>
                <Icon icon={faFilter} />
                <FormattedMessage id="filters" />
              </>
            }
            onEditModeChange={onFilterEditModeChange}
          >
            {/* Tier 1: Primary Filters */}
            <FilterWrapper filterId="rating">
              <SidebarRatingFilter
                title={<FormattedMessage id="rating" />}
                option={RatingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="rating"
              />
            </FilterWrapper>
            <FilterWrapper filterId="favorite">
              <SidebarBooleanFilter
                title={<FormattedMessage id="favourite" />}
                option={FavoritePerformerCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="favourite"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tags">
              <SidebarTagsFilter
                title={<FormattedMessage id="tags" />}
                option={TagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="studios">
              <SidebarStudiosFilter
                title={<FormattedMessage id="studios" />}
                option={StudiosCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="studios"
              />
            </FilterWrapper>
            <FilterWrapper filterId="groups">
              <SidebarGroupsFilter
                title={<FormattedMessage id="groups" />}
                option={GroupsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="groups"
              />
            </FilterWrapper>
            <FilterWrapper filterId="gender">
              <SidebarGenderFilter
                title={<FormattedMessage id="gender" />}
                option={GenderCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="gender"
              />
            </FilterWrapper>
            <FilterWrapper filterId="age">
              <SidebarAgeFilter
                title={<FormattedMessage id="age" />}
                option={AgeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="age"
              />
            </FilterWrapper>
            <FilterWrapper filterId="country">
              <SidebarCountryFilter
                title={<FormattedMessage id="country" />}
                option={CountryCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="country"
              />
            </FilterWrapper>

            {/* Tier 2: Physical Attributes */}
            <FilterWrapper filterId="ethnicity">
              <SidebarStringFilter
                title={<FormattedMessage id="ethnicity" />}
                option={EthnicityCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="ethnicity"
              />
            </FilterWrapper>
            <FilterWrapper filterId="circumcised">
              <SidebarCircumcisedFilter
                title={<FormattedMessage id="circumcised" />}
                option={CircumcisedCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="circumcised"
              />
            </FilterWrapper>
            <FilterWrapper filterId="hair_color">
              <SidebarStringFilter
                title={<FormattedMessage id="hair_color" />}
                option={HairColorCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="hair_color"
              />
            </FilterWrapper>
            <FilterWrapper filterId="eye_color">
              <SidebarStringFilter
                title={<FormattedMessage id="eye_color" />}
                option={EyeColorCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="eye_color"
              />
            </FilterWrapper>
            <FilterWrapper filterId="height_cm">
              <SidebarNumberFilter
                title={<FormattedMessage id="height" />}
                option={HeightCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="height_cm"
              />
            </FilterWrapper>
            <FilterWrapper filterId="weight">
              <SidebarNumberFilter
                title={<FormattedMessage id="weight" />}
                option={WeightCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="weight"
              />
            </FilterWrapper>
            <FilterWrapper filterId="penis_length">
              <SidebarNumberFilter
                title={<FormattedMessage id="penis_length" />}
                option={PenisLengthCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="penis_length"
              />
            </FilterWrapper>
            <FilterWrapper filterId="measurements">
              <SidebarStringFilter
                title={<FormattedMessage id="measurements" />}
                option={MeasurementsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="measurements"
              />
            </FilterWrapper>
            <FilterWrapper filterId="fake_tits">
              <SidebarStringFilter
                title={<FormattedMessage id="fake_tits" />}
                option={FakeTitsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="fake_tits"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tattoos">
              <SidebarStringFilter
                title={<FormattedMessage id="tattoos" />}
                option={TattoosCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tattoos"
              />
            </FilterWrapper>
            <FilterWrapper filterId="piercings">
              <SidebarStringFilter
                title={<FormattedMessage id="piercings" />}
                option={PiercingsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="piercings"
              />
            </FilterWrapper>

            {/* Tier 3: Dates */}
            <FilterWrapper filterId="birthdate">
              <SidebarDateFilter
                title={<FormattedMessage id="birthdate" />}
                option={BirthdateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="birthdate"
              />
            </FilterWrapper>
            <FilterWrapper filterId="death_date">
              <SidebarDateFilter
                title={<FormattedMessage id="death_date" />}
                option={DeathDateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="death_date"
              />
            </FilterWrapper>
            <FilterWrapper filterId="birth_year">
              <SidebarNumberFilter
                title={<FormattedMessage id="birth_year" />}
                option={BirthYearCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="birth_year"
              />
            </FilterWrapper>
            <FilterWrapper filterId="death_year">
              <SidebarNumberFilter
                title={<FormattedMessage id="death_year" />}
                option={DeathYearCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="death_year"
              />
            </FilterWrapper>
            <FilterWrapper filterId="career_length">
              <SidebarStringFilter
                title={<FormattedMessage id="career_length" />}
                option={CareerLengthCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="career_length"
              />
            </FilterWrapper>

            {/* Tier 4: Library Stats */}
            <FilterWrapper filterId="scene_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="scene_count" />}
                option={SceneCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="scene_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="image_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="image_count" />}
                option={ImageCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="image_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="gallery_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="gallery_count" />}
                option={GalleryCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="gallery_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="tag_count" />}
                option={TagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tag_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="play_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="play_count" />}
                option={PlayCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="play_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="o_counter">
              <SidebarNumberFilter
                title={<FormattedMessage id="o_count" />}
                option={OCounterCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="o_counter"
              />
            </FilterWrapper>

            {/* Tier 5: Metadata */}
            <FilterWrapper filterId="name">
              <SidebarStringFilter
                title={<FormattedMessage id="name" />}
                option={NameCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="name"
              />
            </FilterWrapper>
            <FilterWrapper filterId="aliases">
              <SidebarStringFilter
                title={<FormattedMessage id="aliases" />}
                option={AliasesCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="aliases"
              />
            </FilterWrapper>
            <FilterWrapper filterId="disambiguation">
              <SidebarStringFilter
                title={<FormattedMessage id="disambiguation" />}
                option={DisambiguationCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="disambiguation"
              />
            </FilterWrapper>
            <FilterWrapper filterId="details">
              <SidebarStringFilter
                title={<FormattedMessage id="details" />}
                option={DetailsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="details"
              />
            </FilterWrapper>
            <FilterWrapper filterId="url">
              <SidebarStringFilter
                title={<FormattedMessage id="url" />}
                option={UrlCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="url"
              />
            </FilterWrapper>

            {/* Tier 6: System */}
            <FilterWrapper filterId="is_missing">
              <SidebarIsMissingFilter
                title={<FormattedMessage id="isMissing" />}
                option={PerformerIsMissingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="is_missing"
              />
            </FilterWrapper>
            <FilterWrapper filterId="ignore_auto_tag">
              <SidebarBooleanFilter
                title={<FormattedMessage id="ignore_auto_tag" />}
                option={IgnoreAutoTagCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="ignore_auto_tag"
              />
            </FilterWrapper>
            <FilterWrapper filterId="created_at">
              <SidebarDateFilter
                title={<FormattedMessage id="created_at" />}
                option={CreatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="created_at"
                isTime
              />
            </FilterWrapper>
            <FilterWrapper filterId="updated_at">
              <SidebarDateFilter
                title={<FormattedMessage id="updated_at" />}
                option={UpdatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="updated_at"
                isTime
              />
            </FilterWrapper>
            <FilterWrapper filterId="stash_id">
              <SidebarStashIDFilter
                title={<FormattedMessage id="stash_id" />}
                option={StashIDCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="stash_id"
              />
            </FilterWrapper>
          </SidebarFilterSelector>
        </div>
      </MyPerformersFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
      <div className="clear-all-filters">
        <Button
          className="clear-all-filters-button"
          variant="secondary"
          onClick={() => clearAllCriteria()}
          title="Clear All Filters"
        >
          <FormattedMessage id="Clear All Filters" />
        </Button>
      </div>
    </>
  );
};

interface IOperations {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  className?: string;
}

const PerformerListOperations: React.FC<{
  items: number;
  hasSelection: boolean;
  operations: IOperations[];
  onEdit: () => void;
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({ items, hasSelection, operations, onEdit, onDelete, onCreateNew }) => {
  const intl = useIntl();

  return (
    <div className="list-operations">
      <ButtonGroup>
        {!hasSelection && (
          <Button
            className="create-new-button"
            variant="secondary"
            onClick={() => onCreateNew()}
            title={intl.formatMessage(
              { id: "actions.create_entity" },
              { entityType: intl.formatMessage({ id: "performer" }) }
            )}
          >
            <Icon icon={faPlus} />
          </Button>
        )}

        {hasSelection && (
          <>
            <Button variant="secondary" onClick={() => onEdit()}>
              <Icon icon={faPencil} />
            </Button>
            <Button
              variant="danger"
              className="btn-danger-minimal"
              onClick={() => onDelete()}
            >
              <Icon icon={faTrash} />
            </Button>
          </>
        )}

        <OperationDropdown
          className="list-operations"
          menuPortalTarget={document.body}
        >
          {operations.map((o) => {
            if (o.isDisplayed && !o.isDisplayed()) {
              return null;
            }

            return (
              <OperationDropdownItem
                key={o.text}
                onClick={o.onClick}
                text={o.text}
                className={o.className}
              />
            );
          })}
        </OperationDropdown>
      </ButtonGroup>
    </div>
  );
};

interface IFilteredPerformers {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const MyFilteredPerformerList = (props: IFilteredPerformers) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery, extraCriteria } = props;

  // States
  const {
    showSidebar,
    setShowSidebar,
    loading: sidebarStateLoading,
    sectionOpen: baseSectionOpen,
    setSectionOpen: baseSetSectionOpen,
  } = useSidebarState(view);

  // Track filter customization edit mode
  const [isFilterEditMode, setIsFilterEditMode] = useState(false);

  // When in filter edit mode, close all sections and prevent opening
  const sectionOpen = useMemo(() => {
    if (isFilterEditMode) {
      const closedSections: Record<string, boolean> = {};
      Object.keys(baseSectionOpen).forEach((key) => {
        closedSections[key] = false;
      });
      return closedSections;
    }
    return baseSectionOpen;
  }, [isFilterEditMode, baseSectionOpen]);

  const setSectionOpen = useCallback(
    (section: string, open: boolean) => {
      if (isFilterEditMode) return;
      baseSetSectionOpen(section, open);
    },
    [isFilterEditMode, baseSetSectionOpen]
  );

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Performers,
        defaultSort,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindPerformers,
        getCount: (r) => r.data?.findPerformers.count ?? 0,
        getItems: (r) => r.data?.findPerformers.performers ?? [],
        filterHook,
      },
    });

  const { filter, setFilter, loading: filterLoading } = filterState;

  const { effectiveFilter, result, cachedResult, items, totalCount } =
    queryResult;

  const {
    selectedIds,
    selectedItems,
    onSelectChange,
    onSelectAll,
    onSelectNone,
    hasSelection,
  } = listSelect;

  const { modal, showModal, closeModal } = modalState;

  // Utility hooks
  const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
    filter,
    setFilter,
  });

  useAddKeybinds(filter, totalCount);
  useFilteredSidebarKeybinds({
    showSidebar,
    setShowSidebar,
  });

  // Fetch facet counts for sidebar filters
  const { counts: facetCounts, loading: facetLoading } = usePerformerFacetCounts(filter, {
    isOpen: showSidebar ?? false,
    debounceMs: 300,
  });

  useEffect(() => {
    Mousetrap.bind("e", () => {
      if (hasSelection) {
        onEdit?.();
      }
    });

    Mousetrap.bind("d d", () => {
      if (hasSelection) {
        onDelete?.();
      }
    });

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  const onCloseEditDelete = useCloseEditDelete({
    closeModal,
    onSelectNone,
    result,
  });

  const openRandom = useOpenRandom(filter, totalCount);

  function onCreateNew() {
    history.push("/performers/new");
  }

  function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          performers: {
            ids: Array.from(selectedIds.values()),
            all: all,
          },
        }}
        onClose={() => closeModal()}
      />
    );
  }

  function onEdit() {
    showModal(
      <EditPerformersDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
      />
    );
  }

  function onDelete() {
    showModal(
      <DeleteEntityDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
        singularEntity={intl.formatMessage({ id: "performer" })}
        pluralEntity={intl.formatMessage({ id: "performers" })}
        destroyMutation={usePerformersDestroy}
      />
    );
  }

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.open_random" }),
      onClick: openRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "performer" }) }
      ),
      onClick: () => onCreateNew(),
      isDisplayed: () => !hasSelection,
      className: "create-new-item",
    },
    {
      text: intl.formatMessage({ id: "actions.select_all" }),
      onClick: () => onSelectAll(),
      isDisplayed: () => totalCount > 0,
    },
    {
      text: intl.formatMessage({ id: "actions.select_none" }),
      onClick: () => onSelectNone(),
      isDisplayed: () => hasSelection,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: () => onExport(false),
      isDisplayed: () => hasSelection,
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: () => onExport(true),
    },
  ];

  // render
  if (filterLoading || sidebarStateLoading) return null;

  const operations = (
    <PerformerListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onEdit={onEdit}
      onDelete={onDelete}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <TaggerContext>
      <div
        className={cx("item-list-container performer-list", {
          "hide-sidebar": !showSidebar,
        })}
      >
        {modal}

        <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen, disabled: isFilterEditMode }}>
          <FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
          <SidebarPane hideSidebar={!showSidebar}>
            <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
              <SidebarContent
                filter={filter}
                setFilter={setFilter}
                onFilterEditModeChange={setIsFilterEditMode}
                filterHook={filterHook}
                showEditFilter={showEditFilter}
                view={view}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                focus={searchFocus}
                clearAllCriteria={() => clearAllCriteria(true)}
              />
            </Sidebar>
            <SidebarPaneContent>
              <FilteredListToolbar2
                className="performer-list-toolbar"
                hasSelection={hasSelection}
                filterSection={
                  <ToolbarFilterSection
                    filter={filter}
                    onSetFilter={setFilter}
                    onToggleSidebar={() => setShowSidebar(!showSidebar)}
                    onEditCriterion={(c) =>
                      showEditFilter(c?.criterionOption.type)
                    }
                    onRemoveCriterion={removeCriterion}
                    onRemoveAllCriterion={() => clearAllCriteria(true)}
                    onEditSearchTerm={() => {
                      setShowSidebar(true);
                      setSearchFocus(true);
                    }}
                    onRemoveSearchTerm={() =>
                      setFilter(filter.clearSearchTerm())
                    }
                  />
                }
                selectionSection={
                  <ToolbarSelectionSection
                    selected={selectedIds.size}
                    onToggleSidebar={() => setShowSidebar(!showSidebar)}
                    onSelectAll={() => onSelectAll()}
                    onSelectNone={() => onSelectNone()}
                    operations={operations}
                  />
                }
                operationSection={operations}
              />

              <ListResultsHeader
                loading={cachedResult.loading}
                filter={filter}
                totalCount={totalCount}
                onChangeFilter={(newFilter) => setFilter(newFilter)}
              />

              <LoadedContent loading={result.loading} error={result.error}>
                <PerformerListContent
                  filter={effectiveFilter}
                  performers={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
                  extraCriteria={extraCriteria}
                />
              </LoadedContent>

              {totalCount > filter.itemsPerPage && (
                <div className="pagination-footer">
                  <Pagination
                    itemsPerPage={filter.itemsPerPage}
                    currentPage={filter.currentPage}
                    totalItems={totalCount}
                    onChangePage={setPage}
                    pagePopupPlacement="top"
                  />
                </div>
              )}
            </SidebarPaneContent>
          </SidebarPane>
          </FacetCountsContext.Provider>
        </SidebarStateContext.Provider>
      </div>
    </TaggerContext>
  );
};

// Keep the old component for backward compatibility
export const PerformerList: React.FC<IFilteredPerformers> = (props) => {
  return <MyFilteredPerformerList {...props} />;
};

export default MyFilteredPerformerList;
