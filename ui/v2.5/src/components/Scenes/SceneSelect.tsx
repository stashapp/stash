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
  queryFindScenesForSelect,
  queryFindScenesByIDForSelect,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
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
import { objectTitle } from "src/core/files";
import { PatchComponent, PatchFunction } from "src/patch";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { TruncatedText } from "../Shared/TruncatedText";

export type Scene = Pick<GQL.Scene, "id" | "title" | "date" | "code"> & {
  studio?: Pick<GQL.Studio, "name"> | null;
  files?: Pick<GQL.VideoFile, "path">[];
  paths?: Pick<GQL.ScenePathsType, "screenshot">;
};

type Option = SelectOption<Scene>;

type ExtraSceneProps = {
  hoverPlacement?: Placement;
  excludeIds?: string[];
  extraCriteria?: Array<Criterion<CriterionValue>>;
};

type FindScenesResult = Awaited<
  ReturnType<typeof queryFindScenesForSelect>
>["data"]["findScenes"]["scenes"];

function sortScenesByRelevance(input: string, scenes: FindScenesResult) {
  return sortByRelevance(input, scenes, objectTitle, (s) => {
    return s.files.map((f) => f.path);
  });
}

const sceneSelectSort = PatchFunction(
  "SceneSelect.sort",
  sortScenesByRelevance
);

const _SceneSelect: React.FC<
  IFilterProps & IFilterValueProps<Scene> & ExtraSceneProps
> = (props) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadScenes(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Scenes);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "title";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;

    if (props.extraCriteria) {
      filter.criteria = [...props.extraCriteria];
    }

    const query = await queryFindScenesForSelect(filter);
    let ret = query.data.findScenes.scenes.filter((scene) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(scene.id.toString());
    });

    return sceneSelectSort(input, ret).map((scene) => ({
      value: scene.id,
      object: scene,
    }));
  }

  const SceneOption: React.FC<OptionProps<Option, boolean>> = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    const title = objectTitle(object);

    // if title does not match the input value but the path does, show the path
    const { inputValue } = optionProps.selectProps;
    let matchedPath: string | undefined = "";
    if (!title.toLowerCase().includes(inputValue.toLowerCase())) {
      matchedPath = object.files?.find((a) =>
        a.path.toLowerCase().includes(inputValue.toLowerCase())
      )?.path;
    }

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="scene-select-option">
          <span className="scene-select-row">
            {object.paths?.screenshot && (
              <img
                className="scene-select-image"
                src={object.paths.screenshot}
                loading="lazy"
              />
            )}

            <span className="scene-select-details">
              <TruncatedText
                className="scene-select-title"
                text={title}
                lineCount={1}
              />

              {object.studio?.name && (
                <span className="scene-select-studio">
                  {object.studio?.name}
                </span>
              )}

              {object.date && (
                <span className="scene-select-date">{object.date}</span>
              )}

              {object.code && (
                <span className="scene-select-code">{object.code}</span>
              )}
            </span>
          </span>

          {matchedPath && (
            <span className="scene-select-alias">{`(${matchedPath})`}</span>
          )}
        </span>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const SceneMultiValueLabel: React.FC<
    MultiValueGenericProps<Option, boolean>
  > = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: objectTitle(object),
    };

    return <reactSelectComponents.MultiValueLabel {...thisOptionProps} />;
  };

  const SceneValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: <>{objectTitle(object)}</>,
    };

    return <reactSelectComponents.SingleValue {...thisOptionProps} />;
  };

  return (
    <FilterSelectComponent<Scene, boolean>
      {...props}
      className={cx(
        "scene-select",
        {
          "scene-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadScenes}
      components={{
        Option: SceneOption,
        MultiValueLabel: SceneMultiValueLabel,
        SingleValue: SceneValueLabel,
      }}
      isMulti={props.isMulti ?? false}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          {
            entityType: intl.formatMessage({
              id: props.isMulti ? "scenes" : "scene",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const SceneSelect = PatchComponent("SceneSelect", _SceneSelect);

const _SceneIDSelect: React.FC<
  IFilterProps & IFilterIDProps<Scene> & ExtraSceneProps
> = (props) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Scene[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Scene[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Scene[]> {
    const query = await queryFindScenesByIDForSelect(idsToLoad);
    const { scenes: loadedScenes } = query.data.findScenes;

    return loadedScenes;
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

  return <SceneSelect {...props} values={values} onSelect={onSelect} />;
};

export const SceneIDSelect = PatchComponent("SceneIDSelect", _SceneIDSelect);
