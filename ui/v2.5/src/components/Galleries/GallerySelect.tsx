import React, { useEffect, useMemo, useState } from "react";
import {
  OptionProps,
  components as reactSelectComponents,
  MultiValueGenericProps,
  SingleValueProps,
} from "react-select";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import {
  queryFindGalleriesForSelect,
  queryFindGalleriesByIDForSelect,
  useGalleryCreate,
} from "src/core/StashService";
import { useConfigurationContext } from "src/hooks/Config";
import { useIntl } from "react-intl";
import { defaultMaxOptionsShown } from "src/core/config";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  FilterSelectComponent,
  IFilterIDProps,
  IFilterProps,
  IFilterValueProps,
  Option as SelectOption,
} from "../Shared/FilterSelect";
import { useCompare } from "src/hooks/state";
import { Placement } from "react-bootstrap/esm/Overlay";
import { sortByRelevance } from "src/utils/query";
import { galleryTitle } from "src/core/galleries";
import { PatchComponent, PatchFunction } from "src/patch";
import {
  ModifierCriterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { PathCriterion } from "src/models/list-filter/criteria/path";
import { TruncatedText } from "../Shared/TruncatedText";

export type Gallery = Pick<GQL.Gallery, "id" | "title" | "date" | "code"> & {
  studio?: Pick<GQL.Studio, "name"> | null;
  files: Pick<GQL.GalleryFile, "path">[];
  folder?: Pick<GQL.Folder, "path"> | null;
  cover?: Pick<GQL.Image, "paths"> | null;
};
type Option = SelectOption<Gallery>;

type ExtraGalleryProps = {
  hoverPlacement?: Placement;
  excludeIds?: string[];
  extraCriteria?: Array<ModifierCriterion<CriterionValue>>;
};

type FindGalleriesResult = Awaited<
  ReturnType<typeof queryFindGalleriesForSelect>
>["data"]["findGalleries"]["galleries"];

function sortGalleriesByRelevance(
  input: string,
  galleries: FindGalleriesResult
) {
  return sortByRelevance(input, galleries, galleryTitle, (g) => {
    return g.files.map((f) => f.path).concat(g.folder?.path ?? []);
  });
}

const gallerySelectSort = PatchFunction(
  "GallerySelect.sort",
  sortGalleriesByRelevance
);

const _GallerySelect: React.FC<
  IFilterProps & IFilterValueProps<Gallery> & ExtraGalleryProps
> = (props) => {
  const [createGallery] = useGalleryCreate();

  const { configuration } = useConfigurationContext();
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.gallery;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadGalleries(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Galleries);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "title";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;

    if (props.extraCriteria) {
      filter.criteria = [...props.extraCriteria];
    }

    const query = await queryFindGalleriesForSelect(filter);
    let ret = query.data.findGalleries.galleries.filter((gallery) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(gallery.id.toString());
    });

    return gallerySelectSort(input, ret).map((gallery) => ({
      value: gallery.id,
      object: gallery,
    }));
  }

  const GalleryOption: React.FC<OptionProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    const title = galleryTitle(object);

    // if title does not match the input value but the path does, show the path
    const { inputValue } = optionProps.selectProps;
    let matchedPath: string | undefined = "";
    if (!title.toLowerCase().includes(inputValue.toLowerCase())) {
      matchedPath = object.files?.find((a) =>
        a.path.toLowerCase().includes(inputValue.toLowerCase())
      )?.path;

      if (
        !matchedPath &&
        object.folder?.path.toLowerCase().includes(inputValue.toLowerCase())
      ) {
        matchedPath = object.folder?.path;
      }
    }

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="gallery-select-option">
          <span className="gallery-select-row">
            {object.cover?.paths?.thumbnail && (
              <img
                className="gallery-select-image"
                src={object.cover.paths.thumbnail}
                loading="lazy"
              />
            )}

            <span className="gallery-select-details">
              <TruncatedText
                className="gallery-select-title"
                text={title}
                lineCount={1}
              />

              {object.studio?.name && (
                <span className="gallery-select-studio">
                  {object.studio?.name}
                </span>
              )}

              {object.date && (
                <span className="gallery-select-date">{object.date}</span>
              )}

              {object.code && (
                <span className="gallery-select-code">{object.code}</span>
              )}
            </span>
          </span>

          {matchedPath && (
            <span className="gallery-select-alias">{`(${matchedPath})`}</span>
          )}
        </span>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const GalleryMultiValueLabel: React.FC<
    MultiValueGenericProps<Option, boolean>
  > = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: galleryTitle(object),
    };

    return <reactSelectComponents.MultiValueLabel {...thisOptionProps} />;
  };

  const GalleryValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: <>{galleryTitle(object)}</>,
    };

    return <reactSelectComponents.SingleValue {...thisOptionProps} />;
  };

  const onCreate = async (name: string) => {
    const result = await createGallery({
      variables: { input: { title: name } },
    });
    return {
      value: result.data!.galleryCreate!.id,
      item: result.data!.galleryCreate!,
      message: "Created gallery",
    };
  };

  const getNamedObject = (id: string, name: string): Gallery => {
    return {
      id,
      title: name,
      files: [],
      folder: null,
    };
  };

  const isValidNewOption = (inputValue: string, options: Gallery[]) => {
    if (!inputValue) {
      return false;
    }

    if (
      options.some((o) => {
        return galleryTitle(o).toLowerCase() === inputValue.toLowerCase();
      })
    ) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent<Gallery, boolean>
      {...props}
      className={cx(
        "gallery-select",
        {
          "gallery-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadGalleries}
      getNamedObject={getNamedObject}
      isValidNewOption={isValidNewOption}
      components={{
        Option: GalleryOption,
        MultiValueLabel: GalleryMultiValueLabel,
        SingleValue: GalleryValueLabel,
      }}
      isMulti={props.isMulti ?? false}
      creatable={props.creatable ?? defaultCreatable}
      onCreate={onCreate}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          {
            entityType: intl.formatMessage({
              id: props.isMulti ? "galleries" : "gallery",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const GallerySelect = PatchComponent("GallerySelect", _GallerySelect);

const _GalleryIDSelect: React.FC<
  IFilterProps & IFilterIDProps<Gallery> & ExtraGalleryProps
> = (props) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Gallery[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Gallery[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Gallery[]> {
    const query = await queryFindGalleriesByIDForSelect(idsToLoad);
    const { galleries: loadedGalleries } = query.data.findGalleries;

    return loadedGalleries;
  }

  useEffect(() => {
    if (!idsChanged) {
      return;
    }

    if (!ids || ids?.length === 0) {
      setValues([]);
      return;
    }

    // load the values if we have ids and they haven't been loaded yet
    const filteredValues = values.filter((v) => ids.includes(v.id.toString()));
    if (filteredValues.length === ids.length) {
      return;
    }

    const load = async () => {
      const items = await loadObjectsByID(ids);
      setValues(items);
    };

    load();
  }, [ids, idsChanged, values]);

  return <GallerySelect {...props} values={values} onSelect={onSelect} />;
};

export const GalleryIDSelect = PatchComponent(
  "GalleryIDSelect",
  _GalleryIDSelect
);

function getExcludeFilebaseGalleriesFilter() {
  const ret = new PathCriterion();
  ret.modifier = GQL.CriterionModifier.IsNull;
  return ret;
}

export const excludeFileBasedGalleries = [getExcludeFilebaseGalleriesFilter()];
