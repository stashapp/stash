import React, { useCallback, useContext, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindScenes, useFindScenes } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { Tagger } from "src/components/Tagger/scenes/SceneTagger";
import { IPlaySceneOptions, SceneQueue } from "src/models/sceneQueue";
import { SceneWallPanel } from "src/components/Scenes/SceneWallPanel";
import { SceneListTable } from "src/components/Scenes/SceneListTable";
import { EditScenesDialog } from "src/components/Scenes/EditScenesDialog";
import { DeleteScenesDialog } from "src/components/Scenes/DeleteScenesDialog";
import { GenerateDialog } from "src/components/Dialogs/GenerateDialog";
import { ExportDialog } from "src/components/Shared/ExportDialog";
import { SceneCardsGrid } from "src/components/Scenes/SceneCardsGrid";
import { TaggerContext } from "src/components/Tagger/context";
import { IdentifyDialog } from "src/components/Dialogs/IdentifyDialog/IdentifyDialog";
import { ConfigurationContext } from "src/hooks/Config";
import {
  faFilter,
  faPencil,
  faPlay,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import { SceneMergeModal } from "src/components/Scenes/SceneMergeDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";
import { View } from "src/components/List/views";
import { FileSize } from "src/components/Shared/FileSize";
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
import {
  SidebarPerformersFilter,
  SidebarStudiosFilter,
  SidebarGroupsFilter,
  SidebarTagsFilter,
  SidebarRatingFilter,
  SidebarBooleanFilter,
  SidebarPhashFilter,
  SidebarPathFilter,
  SidebarNumberFilter,
  SidebarStashIDFilter,
  SidebarStringFilter,
  SidebarDateFilter,
  SidebarOrientationFilter,
  SidebarPerformerTagsFilter,
  SidebarCaptionsFilter,
  SidebarResolutionFilter,
  SidebarFilterSelector,
  FilterWrapper,
  SidebarDurationFilter,
  SidebarIsMissingFilter,
} from "src/extensions/filters";
import { PerformersCriterionOption } from "src/models/list-filter/criteria/performers";
import { StudiosCriterionOption } from "src/models/list-filter/criteria/studios";
import { GroupsCriterionOption } from "src/models/list-filter/criteria/groups";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "src/models/list-filter/criteria/tags";
import cx from "classnames";
import { RatingCriterionOption } from "src/models/list-filter/criteria/rating";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "src/extensions/ui";
import { PatchContainerComponent } from "src/patch";
import { Pagination } from "src/components/List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  createDateCriterionOption,
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { useFocus } from "src/extensions/hooks";
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "src/extensions/ui";
import { ListResultsHeader } from "src/extensions/ui";
import {
  DuplicatedCriterionOption,
  PhashCriterionOption,
} from "src/models/list-filter/criteria/phash";
import { PathCriterionOption } from "src/models/list-filter/criteria/path";
import { StashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { OrientationCriterionOption } from "src/models/list-filter/criteria/orientation";
import { CaptionsCriterionOption } from "src/models/list-filter/criteria/captions";
import { ResolutionCriterionOption } from "src/models/list-filter/criteria/resolution";
import { SidebarFilterDefinition } from "src/hooks/useSidebarFilters";
import { InteractiveCriterionOption } from "src/models/list-filter/criteria/interactive";
import { HasMarkersCriterionOption } from "src/models/list-filter/criteria/has-markers";
import { PerformerFavoriteCriterionOption } from "src/models/list-filter/criteria/favorite";
import {
  createDurationCriterionOption,
  createMandatoryTimestampCriterionOption,
  createMandatoryStringCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { SceneIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import {
  useSceneFacetCounts,
  FacetCountsContext,
} from "src/hooks/useFacetCounts";

function renderMetadataByline(result: GQL.FindScenesQueryResult) {
  const duration = result?.data?.findScenes?.duration;
  const size = result?.data?.findScenes?.filesize;

  if (!duration && !size) {
    return;
  }

  const separator = duration && size ? " - " : "";

  return (
    <span className="scenes-stats">
      &nbsp;(
      {duration ? (
        <span className="scenes-duration">
          {TextUtils.secondsAsTimeString(duration, 3)}
        </span>
      ) : undefined}
      {separator}
      {size ? (
        <span className="scenes-size">
          <FileSize size={size} />
        </span>
      ) : undefined}
      )
    </span>
  );
}

function usePlayScene() {
  const history = useHistory();

  const { configuration: config } = useContext(ConfigurationContext);
  const cont = config?.interface.continuePlaylistDefault ?? false;
  const autoPlay = config?.interface.autostartVideoOnPlaySelected ?? false;

  const playScene = useCallback(
    (queue: SceneQueue, sceneID: string, options?: IPlaySceneOptions) => {
      history.push(
        queue.makeLink(sceneID, { autoPlay, continue: cont, ...options })
      );
    },
    [history, cont, autoPlay]
  );

  return playScene;
}

function usePlaySelected(selectedIds: Set<string>) {
  const playScene = usePlayScene();

  const playSelected = useCallback(() => {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);

    playScene(queue, sceneIDs[0]);
  }, [selectedIds, playScene]);

  return playSelected;
}

function usePlayFirst() {
  const playScene = usePlayScene();

  const playFirst = useCallback(
    (queue: SceneQueue, sceneID: string, index: number) => {
      // populate queue and go to first scene
      playScene(queue, sceneID, { sceneIndex: index });
    },
    [playScene]
  );

  return playFirst;
}

function usePlayRandom(filter: ListFilterModel, count: number) {
  const playScene = usePlayScene();

  const playRandom = useCallback(async () => {
    // query for a random scene
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
    const queryResults = await queryFindScenes(filterCopy);
    const scene = queryResults.data.findScenes.scenes[index];
    if (scene) {
      // navigate to the image player page
      const queue = SceneQueue.fromListFilterModel(filterCopy);
      playScene(queue, scene.id, { sceneIndex: index });
    }
  }, [filter, count, playScene]);

  return playRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const playRandom = usePlayRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      playRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [playRandom]);
}

const SceneList: React.FC<{
  scenes: GQL.SlimSceneDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
}> = ({ scenes, filter, selectedIds, onSelectChange, fromGroupId }) => {
  const queue = useMemo(() => SceneQueue.fromListFilterModel(filter), [filter]);

  if (scenes.length === 0) {
    return null;
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <SceneCardsGrid
        scenes={scenes}
        queue={queue}
        zoomIndex={filter.zoomIndex}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
        fromGroupId={fromGroupId}
      />
    );
  }
  if (filter.displayMode === DisplayMode.List) {
    return (
      <SceneListTable
        scenes={scenes}
        queue={queue}
        selectedIds={selectedIds}
        onSelectChange={onSelectChange}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Wall) {
    return (
      <SceneWallPanel
        scenes={scenes}
        sceneQueue={queue}
        zoomIndex={filter.zoomIndex}
      />
    );
  }
  if (filter.displayMode === DisplayMode.Tagger) {
    return <Tagger scenes={scenes} queue={queue} />;
  }

  return null;
};

export const ScenesFilterSidebarSections = PatchContainerComponent(
  "FilteredSceneList.SidebarSections"
);

// Define available filters for scenes sidebar
const sceneFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "rating", messageId: "rating", defaultVisible: true },
  { id: "date", messageId: "date", defaultVisible: true },
  { id: "tags", messageId: "tags", defaultVisible: true },
  { id: "performers", messageId: "performers", defaultVisible: true },
  { id: "studios", messageId: "studios", defaultVisible: true },
  { id: "groups", messageId: "groups", defaultVisible: false },
  { id: "organized", messageId: "organized", defaultVisible: true },
  { id: "duration", messageId: "duration", defaultVisible: true },

  // Tier 2: Discovery filters
  { id: "performer_tags", messageId: "performer_tags", defaultVisible: false },
  { id: "performer_favorite", messageId: "performer_favorite", defaultVisible: false },
  { id: "resolution", messageId: "resolution", defaultVisible: false },
  { id: "orientation", messageId: "orientation", defaultVisible: false },
  { id: "captions", messageId: "captions", defaultVisible: false },

  // Tier 3: Activity/History filters
  { id: "play_count", messageId: "play_count", defaultVisible: false },
  { id: "last_played_at", messageId: "last_played_at", defaultVisible: false },
  { id: "o_counter", messageId: "o_count", defaultVisible: false },
  { id: "play_duration", messageId: "play_duration", defaultVisible: false },
  { id: "resume_time", messageId: "resume_time", defaultVisible: false },
  { id: "has_markers", messageId: "has_markers", defaultVisible: false },

  // Tier 4: Technical/Media Info filters
  { id: "framerate", messageId: "framerate", defaultVisible: false },
  { id: "bitrate", messageId: "bitrate", defaultVisible: false },
  { id: "video_codec", messageId: "video_codec", defaultVisible: false },
  { id: "audio_codec", messageId: "audio_codec", defaultVisible: false },
  { id: "file_count", messageId: "file_count", defaultVisible: false },
  { id: "interactive", messageId: "interactive", defaultVisible: false },
  { id: "interactive_speed", messageId: "interactive_speed", defaultVisible: false },

  // Tier 5: Metadata Search filters
  { id: "title", messageId: "title", defaultVisible: false },
  { id: "code", messageId: "scene_code", defaultVisible: false },
  { id: "details", messageId: "details", defaultVisible: false },
  { id: "director", messageId: "director", defaultVisible: false },
  { id: "url", messageId: "url", defaultVisible: false },
  { id: "path", messageId: "path", defaultVisible: false },

  // Tier 6: Library Management filters
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
  { id: "tag_count", messageId: "tag_count", defaultVisible: false },
  { id: "performer_count", messageId: "performer_count", defaultVisible: false },
  { id: "performer_age", messageId: "performer_age", defaultVisible: false },

  // Tier 7: Technical/Duplicate Detection filters
  { id: "phash", messageId: "media_info.phash", defaultVisible: false },
  { id: "duplicated_phash", messageId: "duplicated_phash", defaultVisible: false },
  { id: "oshash", messageId: "media_info.hash", defaultVisible: false },
  { id: "checksum", messageId: "media_info.checksum", defaultVisible: false },
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

  const hideStudios = view === View.StudioScenes;
  
  // Criterion options
  const fileCountCriterionOption = createMandatoryNumberCriterionOption("file_count");
  const UrlCriterionOption = createStringCriterionOption("url");
  const DateCriterionOption = createDateCriterionOption("date");
  const TitleCriterionOption = createStringCriterionOption("title");
  const CodeCriterionOption = createStringCriterionOption("code", "scene_code");
  const DetailsCriterionOption = createStringCriterionOption("details");
  const DirectorCriterionOption = createStringCriterionOption("director");
  const VideoCodecCriterionOption = createStringCriterionOption("video_codec");
  const AudioCodecCriterionOption = createStringCriterionOption("audio_codec");
  const OshashCriterionOption = createMandatoryStringCriterionOption("oshash", "media_info.hash");
  const ChecksumCriterionOption = createStringCriterionOption("checksum", "media_info.checksum");
  const OCounterCriterionOption = createMandatoryNumberCriterionOption("o_counter", "o_count");
  const FramerateCriterionOption = createMandatoryNumberCriterionOption("framerate");
  const BitrateCriterionOption = createMandatoryNumberCriterionOption("bitrate");
  const PlayCountCriterionOption = createMandatoryNumberCriterionOption("play_count");
  const TagCountCriterionOption = createMandatoryNumberCriterionOption("tag_count");
  const PerformerCountCriterionOption = createMandatoryNumberCriterionOption("performer_count");
  const PerformerAgeCriterionOption = createMandatoryNumberCriterionOption("performer_age");
  const InteractiveSpeedCriterionOption = createMandatoryNumberCriterionOption("interactive_speed");
  const DurationCriterionOption = createDurationCriterionOption("duration");
  const ResumeTimeCriterionOption = createDurationCriterionOption("resume_time");
  const PlayDurationCriterionOption = createDurationCriterionOption("play_duration");
  const LastPlayedAtCriterionOption = createMandatoryTimestampCriterionOption("last_played_at");
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

      <ScenesFilterSidebarSections>
      <div className="sidebar-filters">
        <SidebarFilterSelector
          viewName="scenes"
          filterDefinitions={sceneFilterDefinitions}
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
            <FilterWrapper filterId="date">
              <SidebarDateFilter
                title={<FormattedMessage id="date" />}
                option={DateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="date"
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
            <FilterWrapper filterId="performers">
              <SidebarPerformersFilter
                title={<FormattedMessage id="performers" />}
                option={PerformersCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="performers"
              />
            </FilterWrapper>
            {!hideStudios && (
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
            )}
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
            <FilterWrapper filterId="organized">
              <SidebarBooleanFilter
                title={<FormattedMessage id="organized" />}
                option={OrganizedCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="organized"
              />
            </FilterWrapper>
            <FilterWrapper filterId="duration">
              <SidebarDurationFilter
                title={<FormattedMessage id="duration" />}
                option={DurationCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="duration"
              />
            </FilterWrapper>

            {/* Tier 2: Discovery Filters */}
            <FilterWrapper filterId="performer_tags">
              <SidebarPerformerTagsFilter
                title={<FormattedMessage id="performer_tags" />}
                option={PerformerTagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                sectionID="performer_tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_favorite">
              <SidebarBooleanFilter
                title={<FormattedMessage id="performer_favorite" />}
                option={PerformerFavoriteCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_favorite"
              />
            </FilterWrapper>
            <FilterWrapper filterId="resolution">
              <SidebarResolutionFilter
                title={<FormattedMessage id="resolution" />}
                option={ResolutionCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="resolution"
              />
            </FilterWrapper>
            <FilterWrapper filterId="orientation">
              <SidebarOrientationFilter
                title={<FormattedMessage id="orientation" />}
                option={OrientationCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="orientation"
              />
            </FilterWrapper>
            <FilterWrapper filterId="captions">
              <SidebarCaptionsFilter
                title={<FormattedMessage id="captions" />}
                option={CaptionsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="captions"
              />
            </FilterWrapper>

            {/* Tier 3: Activity/History Filters */}
            <FilterWrapper filterId="play_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="play_count" />}
                option={PlayCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="play_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="last_played_at">
              <SidebarDateFilter
                title={<FormattedMessage id="last_played_at" />}
                option={LastPlayedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="last_played_at"
                isTime
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
            <FilterWrapper filterId="play_duration">
              <SidebarDurationFilter
                title={<FormattedMessage id="play_duration" />}
                option={PlayDurationCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="play_duration"
              />
            </FilterWrapper>
            <FilterWrapper filterId="resume_time">
              <SidebarDurationFilter
                title={<FormattedMessage id="resume_time" />}
                option={ResumeTimeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="resume_time"
              />
            </FilterWrapper>
            <FilterWrapper filterId="has_markers">
              <SidebarBooleanFilter
                title={<FormattedMessage id="has_markers" />}
                option={HasMarkersCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="has_markers"
              />
            </FilterWrapper>

            {/* Tier 4: Technical/Media Info Filters */}
            <FilterWrapper filterId="framerate">
              <SidebarNumberFilter
                title={<FormattedMessage id="framerate" />}
                option={FramerateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="framerate"
              />
            </FilterWrapper>
            <FilterWrapper filterId="bitrate">
              <SidebarNumberFilter
                title={<FormattedMessage id="bitrate" />}
                option={BitrateCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="bitrate"
              />
            </FilterWrapper>
            <FilterWrapper filterId="video_codec">
              <SidebarStringFilter
                title={<FormattedMessage id="video_codec" />}
                option={VideoCodecCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="video_codec"
              />
            </FilterWrapper>
            <FilterWrapper filterId="audio_codec">
              <SidebarStringFilter
                title={<FormattedMessage id="audio_codec" />}
                option={AudioCodecCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="audio_codec"
              />
            </FilterWrapper>
            <FilterWrapper filterId="file_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="file_count" />}
                option={fileCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="file_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="interactive">
              <SidebarBooleanFilter
                title={<FormattedMessage id="interactive" />}
                option={InteractiveCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="interactive"
              />
            </FilterWrapper>
            <FilterWrapper filterId="interactive_speed">
              <SidebarNumberFilter
                title={<FormattedMessage id="interactive_speed" />}
                option={InteractiveSpeedCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="interactive_speed"
              />
            </FilterWrapper>

            {/* Tier 5: Metadata Search Filters */}
            <FilterWrapper filterId="title">
              <SidebarStringFilter
                title={<FormattedMessage id="title" />}
                option={TitleCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="title"
              />
            </FilterWrapper>
            <FilterWrapper filterId="code">
              <SidebarStringFilter
                title={<FormattedMessage id="scene_code" />}
                option={CodeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="code"
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
            <FilterWrapper filterId="director">
              <SidebarStringFilter
                title={<FormattedMessage id="director" />}
                option={DirectorCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="director"
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
            <FilterWrapper filterId="path">
              <SidebarPathFilter
                title={<FormattedMessage id="path" />}
                option={PathCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="path"
              />
            </FilterWrapper>

            {/* Tier 6: Library Management Filters */}
            <FilterWrapper filterId="is_missing">
              <SidebarIsMissingFilter
                title={<FormattedMessage id="isMissing" />}
                option={SceneIsMissingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="is_missing"
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
            <FilterWrapper filterId="tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="tag_count" />}
                option={TagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="tag_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="performer_count" />}
                option={PerformerCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_age">
              <SidebarNumberFilter
                title={<FormattedMessage id="performer_age" />}
                option={PerformerAgeCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_age"
              />
            </FilterWrapper>

            {/* Tier 7: Technical/Duplicate Detection Filters */}
            <FilterWrapper filterId="phash">
              <SidebarPhashFilter
                title={<FormattedMessage id="media_info.phash" />}
                option={PhashCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="phash"
              />
            </FilterWrapper>
            <FilterWrapper filterId="duplicated_phash">
              <SidebarBooleanFilter
                title={<FormattedMessage id="duplicated_phash" />}
                option={DuplicatedCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="duplicated_phash"
              />
            </FilterWrapper>
            <FilterWrapper filterId="oshash">
              <SidebarStringFilter
                title={<FormattedMessage id="media_info.hash" />}
                option={OshashCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="oshash"
              />
            </FilterWrapper>
            <FilterWrapper filterId="checksum">
              <SidebarStringFilter
                title={<FormattedMessage id="media_info.checksum" />}
                option={ChecksumCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="checksum"
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
      </ScenesFilterSidebarSections>

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
          // TODO: add message
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

const SceneListOperations: React.FC<{
  items: number;
  hasSelection: boolean;
  operations: IOperations[];
  onEdit: () => void;
  onDelete: () => void;
  onPlay: () => void;
  onCreateNew: () => void;
}> = ({
  items,
  hasSelection,
  operations,
  onEdit,
  onDelete,
  onPlay,
  onCreateNew,
}) => {
  const intl = useIntl();

  return (
    <div className="list-operations">
      <ButtonGroup>
        {!!items && (
          <Button
            className="play-button"
            variant="secondary"
            onClick={() => onPlay()}
            title={intl.formatMessage({ id: "actions.play" })}
          >
            <Icon icon={faPlay} />
          </Button>
        )}
        {!hasSelection && (
          <Button
            className="create-new-button"
            variant="secondary"
            onClick={() => onCreateNew()}
            title={intl.formatMessage(
              { id: "actions.create_entity" },
              { entityType: intl.formatMessage({ id: "scene" }) }
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

interface IFilteredScenes {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  fromGroupId?: string;
}

export const FilteredSceneList = (props: IFilteredScenes) => {
  const intl = useIntl();
  const history = useHistory();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery, fromGroupId } = props;

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
      // Return all sections as closed
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
      // Prevent opening sections while in filter edit mode
      if (isFilterEditMode) return;
      baseSetSectionOpen(section, open);
    },
    [isFilterEditMode, baseSetSectionOpen]
  );

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Scenes,
        defaultSort,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindScenes,
        getCount: (r) => r.data?.findScenes.count ?? 0,
        getItems: (r) => r.data?.findScenes.scenes ?? [],
        filterHook,
      },
    });

  const { filter, setFilter, loading: filterLoading } = filterState;

  // Fetch facet counts for sidebar filters
  // Note: showSidebar can be undefined initially, so we default to false
  // Expensive facets (performer_tags, captions) are lazy-loaded only when their section is expanded
  const { counts: facetCounts, loading: facetLoading } = useSceneFacetCounts(filter, { 
    isOpen: showSidebar ?? false,
    debounceMs: 300, // Faster response for filter changes
    includePerformerTags: baseSectionOpen["performer_tags"] ?? false,
    includeCaptions: baseSectionOpen["captions"] ?? false,
  });

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

  const metadataByline = useMemo(() => {
    if (cachedResult.loading) return null;

    return renderMetadataByline(cachedResult) ?? null;
  }, [cachedResult]);

  const queue = useMemo(() => SceneQueue.fromListFilterModel(filter), [filter]);

  const playRandom = usePlayRandom(effectiveFilter, totalCount);
  const playSelected = usePlaySelected(selectedIds);
  const playFirst = usePlayFirst();

  function onCreateNew() {
    history.push("/scenes/new");
  }

  function onPlay() {
    if (items.length === 0) {
      return;
    }

    // if there are selected items, play those
    if (hasSelection) {
      playSelected();
      return;
    }

    // otherwise, play the first item in the list
    const sceneID = items[0].id;
    playFirst(queue, sceneID, 0);
  }

  function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          scenes: {
            ids: Array.from(selectedIds.values()),
            all: all,
          },
        }}
        onClose={() => closeModal()}
      />
    );
  }

  function onMerge() {
    const selected =
      selectedItems.map((s) => {
        return {
          id: s.id,
          title: objectTitle(s),
        };
      }) ?? [];
    showModal(
      <SceneMergeModal
        scenes={selected}
        onClose={(mergedID?: string) => {
          closeModal();
          if (mergedID) {
            history.push(`/scenes/${mergedID}`);
          }
        }}
        show
      />
    );
  }

  function onEdit() {
    showModal(
      <EditScenesDialog selected={selectedItems} onClose={onCloseEditDelete} />
    );
  }

  function onDelete() {
    showModal(
      <DeleteScenesDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
      />
    );
  }

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.play" }),
      onClick: () => onPlay(),
      isDisplayed: () => items.length > 0,
      className: "play-item",
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "scene" }) }
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
      text: intl.formatMessage({ id: "actions.play_random" }),
      onClick: playRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: `${intl.formatMessage({ id: "actions.generate" })}…`,
      onClick: () =>
        showModal(
          <GenerateDialog
            type="scene"
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => closeModal()}
          />
        ),
      isDisplayed: () => hasSelection,
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: () =>
        showModal(
          <IdentifyDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={() => closeModal()}
          />
        ),
      isDisplayed: () => hasSelection,
    },
    {
      text: `${intl.formatMessage({ id: "actions.merge" })}…`,
      onClick: () => onMerge(),
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
    <SceneListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onEdit={onEdit}
      onDelete={onDelete}
      onPlay={onPlay}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <TaggerContext>
      <div
        className={cx("item-list-container scene-list", {
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
                filterHook={filterHook}
                showEditFilter={showEditFilter}
                view={view}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                focus={searchFocus}
                clearAllCriteria={() => clearAllCriteria(true)}
                onFilterEditModeChange={setIsFilterEditMode}
              />
            </Sidebar>
            <SidebarPaneContent>
              <FilteredListToolbar2
                className="scene-list-toolbar"
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
                    // view={view}
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
                metadataByline={metadataByline}
                onChangeFilter={(newFilter) => setFilter(newFilter)}
              />

              <LoadedContent loading={result.loading} error={result.error}>
                <SceneList
                  filter={effectiveFilter}
                  scenes={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
                  fromGroupId={fromGroupId}
                />
              </LoadedContent>

              {totalCount > filter.itemsPerPage && (
                <div className="pagination-footer">
                  <Pagination
                    itemsPerPage={filter.itemsPerPage}
                    currentPage={filter.currentPage}
                    totalItems={totalCount}
                    metadataByline={metadataByline}
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

export default FilteredSceneList;
