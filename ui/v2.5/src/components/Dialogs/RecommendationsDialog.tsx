import React, { useState, useMemo, useEffect } from "react";
import { Form } from "react-bootstrap";
import { Modal } from "src/components/Shared";
import { FiltersList, FiltersEditor } from "./Filters";
import { useIntl } from "react-intl";
import { useToast } from "src/hooks";
import { FilterMode, SavedFilter } from "src/core/generated-graphql";
import {
  useFindSavedFilters,
  useSaveFilter,
  useFindRecommendationFilters,
} from "src/core/StashService";

interface IIdentifyDialogProps {
  onClose: () => void;
}

export const RecommendationsDialog: React.FC<IIdentifyDialogProps> = ({
  onClose,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const [animation, setAnimation] = useState(true);

  const [filters, setFilters] = useState<SavedFilter[]>([]);
  const [editingFilter, setEditingFilter] = useState<SavedFilter | undefined>();
  const [saving, setSaving] = useState(false);
  const [saveFilter] = useSaveFilter();

  const savedSceneFilters = useFindSavedFilters(FilterMode.Scenes);
  const savedGalleryFilters = useFindSavedFilters(FilterMode.Galleries);
  const savedStudioFilters = useFindSavedFilters(FilterMode.Studios);
  const savedMovieFilters = useFindSavedFilters(FilterMode.Movies);
  const savedPerformerFilters = useFindSavedFilters(FilterMode.Performers);

  const filterDefaults = useFindRecommendationFilters();

  const allFilters = useMemo(() => {
    const ret: SavedFilter[] = [];

    savedSceneFilters.data?.findSavedFilters.forEach(function (filter) {
      ret.push(filter);
    });

    savedGalleryFilters.data?.findSavedFilters.forEach(function (filter) {
      ret.push(filter);
    });

    savedStudioFilters.data?.findSavedFilters.forEach(function (filter) {
      ret.push(filter);
    });

    savedMovieFilters.data?.findSavedFilters.forEach(function (filter) {
      ret.push(filter);
    });

    savedPerformerFilters.data?.findSavedFilters.forEach(function (filter) {
      ret.push(filter);
    });

    const sortedRet = ret.sort((f1, f2) => {
      if (f1.name > f2.name) {
        return 1;
      }

      if (f1.name < f2.name) {
        return -1;
      }

      return 0;
    });

    return sortedRet;
  }, [
    savedSceneFilters,
    savedGalleryFilters,
    savedStudioFilters,
    savedMovieFilters,
    savedPerformerFilters,
  ]);

  useEffect(() => {
    if (!allFilters) return;

    if (!!filterDefaults?.data?.findRecommendedFilters.length) {
      const mappedFilters = filterDefaults.data?.findRecommendedFilters
        .map((f) => {
          const found = allFilters.find((fs) => fs.id === f.id);

          if (!found) return;

          const ret: SavedFilter = {
            ...found,
          };

          return ret;
        })
        .filter((f) => f) as SavedFilter[];

      setFilters(mappedFilters);
    }
  }, [allFilters, filterDefaults]);

  if (!allFilters) return <div />;

  function getAvailableFilters() {
    // only include scrapers not already present
    return !editingFilter?.id === undefined
      ? []
      : allFilters?.filter((f) => {
          return !filters.some((fs) => fs.id === f.id);
        }) ?? [];
  }

  function onEditFilter(f?: SavedFilter) {
    setAnimation(false);

    // if undefined, then set a dummy filter to create a new one
    if (!f) {
      setEditingFilter(getAvailableFilters()[0]);
    } else {
      setEditingFilter(f);
    }
  }

  async function saveRecommendation(filter: SavedFilter) {
    try {
      setSaving(true);
      await saveFilter({
        variables: {
          input: {
            id: filter.id,
            mode: filter.mode,
            name: filter.name,
            filter: filter.filter,
            recommendation_index: filter.recommendation_index,
          },
        },
      });
    } catch (err) {
      Toast.error(err);
    } finally {
      setSaving(false);
    }
  }

  async function onSave() {
    try {
      allFilters.forEach(function (filter) {
        filter.recommendation_index = filters
          .map(function (f) {
            return f.id;
          })
          .indexOf(filter.id);
        saveRecommendation(filter);
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  function isNewFilter() {
    return !!editingFilter && !filters.includes(editingFilter);
  }

  function onSaveFilter(f?: SavedFilter) {
    if (f) {
      let found = false;
      const newFilters = filters.map((fs) => {
        if (fs.id === f.id) {
          found = true;
          return f;
        }
        return fs;
      });

      if (!found) {
        newFilters.push(f);
      }

      setFilters(newFilters);
    }
    setEditingFilter(undefined);
  }

  if (editingFilter) {
    return (
      <FiltersEditor
        availableFilters={getAvailableFilters()}
        filter={editingFilter}
        saveFilter={onSaveFilter}
        isNew={isNewFilter()}
      />
    );
  }

  return (
    <Modal
      modalProps={{ animation, size: "lg" }}
      show
      icon="cogs"
      header={intl.formatMessage({ id: "actions.manage_recommendations" })}
      accept={{
        text: "Accept",
        onClick: onSave,
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={editingFilter || saving}
    >
      <Form>
        <FiltersList
          filters={filters}
          setFilters={(f) => setFilters(f)}
          editFilter={onEditFilter}
          canAdd={filters.length < allFilters.length}
        />
      </Form>
    </Modal>
  );
};
