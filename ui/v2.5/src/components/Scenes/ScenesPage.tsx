import React, { useMemo, useState } from "react";
import { DisplayMode } from "src/models/list-filter/types";
import { FilterMode, FindScenesQueryResult } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { queryFindScenes, useFindScenes } from "src/core/StashService";
import { SceneCardsGrid } from "./SceneCardsGrid";
import SceneQueue, { IPlaySceneOptions } from "src/models/sceneQueue";
import { SceneListTable } from "./SceneListTable";
import { SceneWallPanel } from "../Wall/WallPanel";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { TaggerContext } from "../Tagger/context";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import TextUtils from "src/utils/text";
import { useListSelect } from "src/hooks/listSelect";
import { IItemListOperation } from "../List/ItemList";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { ListOperationButtons } from "../List/ListOperationButtons";
import { ListOperationDropdown } from "../List/ListOperationDropdown";
import { faPlay, faPlus, faShuffle } from "@fortawesome/free-solid-svg-icons";
import { ConfigurationContext } from "src/hooks/Config";
import { useHistory } from "react-router-dom";
import { objectTitle } from "src/core/files";
import { useModal } from "src/hooks/modal";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import { IdentifyDialog } from "../Dialogs/IdentifyDialog/IdentifyDialog";
import { SceneMergeModal } from "./SceneMergeDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { getFromIds } from "src/utils/data";
import { EditScenesDialog } from "./EditScenesDialog";
import { DeleteScenesDialog } from "./DeleteScenesDialog";
import { ListPage } from "../List/ListPage";
import { useFilterURL, useInitialFilter, useResultCount } from "../List/util";
import DropdownItem from "react-bootstrap/esm/DropdownItem";

const filterMode = FilterMode.Scenes;

const ScenesPageImpl: React.FC<{
  defaultFilter: ListFilterModel;
}> = ({ defaultFilter }) => {
  const intl = useIntl();
  const history = useHistory();

  const config = React.useContext(ConfigurationContext);

  const [filter, setFilterState] = useState<ListFilterModel>(defaultFilter);

  const { setFilter } = useFilterURL(filter, setFilterState, defaultFilter);

  const result = useFindScenes(filter);
  const { loading } = result;
  const items = useMemo(
    () => result.data?.findScenes.scenes ?? [],
    [result.data?.findScenes.scenes]
  );

  const listSelect = useListSelect(items);
  const { selectedIds, onSelectChange } = listSelect;

  const { modal, showModal, closeModal } = useModal();

  const totalCount = useResultCount(
    filter,
    loading,
    result.data?.findScenes.count ?? 0
  );

  const metadataByline = useMemo(() => {
    const { data } = result;
    const duration = data?.findScenes?.duration;
    const size = data?.findScenes?.filesize;
    const filesize = size ? TextUtils.fileSize(size) : undefined;

    if (loading || (!duration && !size)) {
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
        {size && filesize ? (
          <span className="scenes-size">
            <FormattedNumber
              value={filesize.size}
              maximumFractionDigits={TextUtils.fileSizeFractionalDigits(
                filesize.unit
              )}
            />
            {` ${TextUtils.formatFileSizeUnit(filesize.unit)}`}
          </span>
        ) : undefined}
        )
      </span>
    );
  }, [result, loading]);

  function renderScenes() {
    if (!result.data?.findScenes) return;

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneCardsGrid
          scenes={items}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <SceneListTable
          scenes={items}
          queue={queue}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <SceneWallPanel scenes={items} sceneQueue={queue} />;
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return (
        <TaggerContext>
          <Tagger scenes={items} queue={queue} />
        </TaggerContext>
      );
    }
  }

  function playScene(
    queue: SceneQueue,
    sceneID: string,
    options: IPlaySceneOptions
  ) {
    history.push(queue.makeLink(sceneID, options));
  }

  async function playAll() {
    if (items.length === 0) return;

    const queue = SceneQueue.fromListFilterModel(filter);
    const autoPlay =
      config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
    playScene(queue, items[0].id, { autoPlay });
  }

  async function playSelected() {
    // populate queue and go to first scene
    const sceneIDs = Array.from(selectedIds.values());
    const queue = SceneQueue.fromSceneIDList(sceneIDs);
    const autoPlay =
      config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
    playScene(queue, sceneIDs[0], { autoPlay });
  }

  async function playRandom() {
    if (items.length === 0) return;

    // query for a random scene
    if (result.data?.findScenes) {
      const { count } = result.data.findScenes;
      const pages = Math.ceil(count / filter.itemsPerPage);
      const page = Math.floor(Math.random() * pages) + 1;
      const indexMax = Math.min(filter.itemsPerPage, count);
      const index = Math.floor(Math.random() * indexMax);
      const filterCopy = filter.clone();
      filterCopy.currentPage = page;
      filterCopy.sortBy = "random";
      const queryResults = await queryFindScenes(filterCopy);
      const scene = queryResults.data.findScenes.scenes[index];
      if (scene) {
        // navigate to the image player page
        const queue = SceneQueue.fromListFilterModel(filterCopy);
        const autoPlay =
          config.configuration?.interface.autostartVideoOnPlaySelected ?? false;
        playScene(queue, scene.id, { sceneIndex: index, autoPlay });
      }
    }
  }

  async function onMerge() {
    const selected =
      result.data?.findScenes.scenes
        .filter((s) => selectedIds.has(s.id))
        .map((s) => {
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

  async function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          scenes: {
            ids: Array.from(selectedIds.values()),
            all,
          },
        }}
        onClose={closeModal}
      />
    );
  }

  const otherOperations: IItemListOperation<FindScenesQueryResult>[] = [
    {
      text: `${intl.formatMessage({ id: "actions.generate" })}…`,
      onClick: async () => {
        showModal(
          <GenerateDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={closeModal}
          />
        );
      },
    },
    {
      text: `${intl.formatMessage({ id: "actions.identify" })}…`,
      onClick: async () => {
        showModal(
          <IdentifyDialog
            selectedIds={Array.from(selectedIds.values())}
            onClose={closeModal}
          />
        );
      },
    },
    {
      text: `${intl.formatMessage({ id: "actions.merge" })}…`,
      onClick: onMerge,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: () => onExport(false),
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: () => onExport(true),
    },
  ];

  async function onOperationClicked(
    o: IItemListOperation<FindScenesQueryResult>
  ) {
    await o.onClick(result, filter, selectedIds);
    if (o.postRefetch) {
      result.refetch();
    }
  }

  const operations = otherOperations?.map((o) => ({
    text: o.text,
    onClick: () => {
      onOperationClicked(o);
    },
    isDisplayed: () => {
      if (o.isDisplayed) {
        return o.isDisplayed(result, filter, selectedIds);
      }

      return true;
    },
    icon: o.icon,
    buttonVariant: o.buttonVariant,
  }));

  function renderButtons() {
    return (
      <div>
        <div>
          <Button
            className="play-scenes-button"
            variant="secondary"
            onClick={() => playSelected()}
          >
            <Icon icon={faPlay} />
          </Button>
        </div>

        <ListOperationButtons
          itemsSelected
          onEdit={() =>
            showModal(
              <EditScenesDialog
                selected={getFromIds(items, selectedIds)}
                onClose={closeModal}
              />
            )
          }
          onDelete={() =>
            showModal(
              <DeleteScenesDialog
                selected={getFromIds(items, selectedIds)}
                onClose={closeModal}
              />
            )
          }
        />

        <ListOperationDropdown operations={operations} />
      </div>
    );
  }

  return (
    <>
      <ListPage
        id="scenes-page"
        loading={loading}
        filter={filter}
        setFilter={(f) => setFilter(f)}
        listSelect={listSelect}
        actionButtons={
          <>
            <DropdownItem href="/scenes/new">
              <Icon icon={faPlus} />
              <FormattedMessage id="new" defaultMessage="New" />
            </DropdownItem>
            {items.length !== 0 && (
              <>
                <DropdownItem
                  className="play-scenes-button"
                  onClick={() => playAll()}
                >
                  <Icon icon={faPlay} />
                  <span>
                    <FormattedMessage id="actions.play" />
                  </span>
                </DropdownItem>
                <DropdownItem
                  className="shuffle-scenes-button"
                  onClick={() => playRandom()}
                >
                  <Icon icon={faShuffle} />
                  <span>
                    <FormattedMessage id="actions.shuffle" />
                  </span>
                </DropdownItem>
              </>
            )}
          </>
        }
        selectedButtons={renderButtons}
        metadataByline={metadataByline}
        totalCount={totalCount}
      >
        {renderScenes()}
      </ListPage>
      {modal}
    </>
  );
};

export const ScenesPage: React.FC = () => {
  const initialFilter = useInitialFilter(filterMode);

  if (!initialFilter) {
    return null;
  }

  return <ScenesPageImpl defaultFilter={initialFilter} />;
};

export default ScenesPage;