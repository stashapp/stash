import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  CriterionModifier,
  FilterMode,
  FolderDataFragment,
  MultiCriterionInput,
  useFindFolderHierarchyForIDsQuery,
  useFindFoldersForQueryQuery,
  useFindRootFoldersForSelectQuery,
} from "src/core/generated-graphql";
import {
  ISidebarSectionProps,
  SidebarSection,
} from "src/components/Shared/Sidebar";
import {
  faChevronDown,
  faChevronRight,
  faMinus,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";
import { ExpandCollapseButton } from "src/components/Shared/CollapseButton";
import cx from "classnames";
import { queryFindSubFolders } from "src/core/StashService";
import { keyboardClickHandler } from "src/utils/keyboard";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  FolderCriterion,
  FolderCriterionOption,
} from "src/models/list-filter/criteria/folder";
import { Option, SelectedList } from "./SidebarListFilter";
import {
  defineMessages,
  FormattedMessage,
  MessageDescriptor,
  useIntl,
} from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { Button, Form } from "react-bootstrap";
import { DepthSelector } from "./SelectableFilter";
import ClearableInput from "src/components/Shared/ClearableInput";
import { useDebouncedState } from "src/hooks/debounce";
import { ModifierCriterionOption } from "src/models/list-filter/criteria/criterion";

interface IFolder extends FolderDataFragment {
  children?: IFolder[];
  expanded: boolean;
}

const FolderRow: React.FC<{
  folder: IFolder;
  level?: number;
  canExclude?: boolean;
  toggleExpanded: (folder: IFolder) => void;
  onSelect: (folder: IFolder, exclude?: boolean) => void;
}> = ({ folder, level, toggleExpanded, onSelect, canExclude }) => {
  return (
    <>
      <li
        className="folder-row unselected-object"
        style={{ paddingLeft: (level ?? 0) * 5 }}
      >
        <a
          onClick={() => onSelect(folder)}
          onKeyDown={keyboardClickHandler(() => onSelect(folder))}
          tabIndex={0}
        >
          <span>
            <span
              className={cx({
                empty: folder.children && folder.children.length === 0,
              })}
            >
              <ExpandCollapseButton
                collapsed={!folder.expanded}
                setCollapsed={() => toggleExpanded(folder)}
                collapsedIcon={faChevronRight}
                notCollapsedIcon={faChevronDown}
              />
            </span>
            {folder.basename}
          </span>
          {canExclude && (
            <Button
              onClick={(e) => {
                e.stopPropagation();
                onSelect(folder, true);
              }}
              onKeyDown={(e) => e.stopPropagation()}
              className="minimal exclude-button"
            >
              <span className="exclude-button-text">
                <FormattedMessage id="actions.exclude_lowercase" />
              </span>
              <Icon className="fa-fw exclude-icon" icon={faMinus} />
            </Button>
          )}
        </a>
      </li>
      {folder.expanded &&
        folder.children?.map((child) => (
          <FolderRow
            key={child.id}
            folder={child}
            level={(level ?? 0) + 1}
            toggleExpanded={toggleExpanded}
            onSelect={onSelect}
            canExclude={canExclude}
          />
        ))}
    </>
  );
};

function toggleExpandedFn(object: IFolder): (f: IFolder) => IFolder {
  return (f: IFolder) => {
    if (f.id === object.id) {
      return { ...f, expanded: !f.expanded };
    }

    if (f.children) {
      return {
        ...f,
        children: f.children.map(toggleExpandedFn(object)),
      };
    }

    return f;
  };
}

function replaceFolder(folder: IFolder): (f: IFolder) => IFolder {
  return (f: IFolder) => {
    if (f.id === folder.id) {
      return folder;
    }

    if (f.children) {
      return {
        ...f,
        children: f.children.map(replaceFolder(folder)),
      };
    }

    return f;
  };
}

function mergeFolderMaps(base: IFolder[], update: IFolder[]): IFolder[] {
  const ret = [...base];

  update.forEach((updateFolder) => {
    const existingIndex = ret.findIndex((f) => f.id === updateFolder.id);
    if (existingIndex === -1) {
      // not found, add to the end
      ret.push(updateFolder);
    } else {
      // found, replace
      ret[existingIndex] = updateFolder;
    }
  });

  return ret;
}

function useFolderMap(props: {
  query: string;
  skip?: boolean;
  initialSelected?: string[];
  mode?: FilterMode;
}) {
  const { query, skip = false, initialSelected, mode } = props;

  const [cachedInitialSelected] = useState<string[]>(initialSelected ?? []);

  // exclude zip folders for scenes and galleries
  const excludeZipFolders =
    mode === FilterMode.Scenes || mode === FilterMode.Galleries;

  const zipFileFilter: MultiCriterionInput | undefined = useMemo(
    () =>
      excludeZipFolders
        ? {
            modifier: CriterionModifier.IsNull,
          }
        : undefined,
    [excludeZipFolders]
  );

  const folderFilterForQuery = useMemo(
    () => (zipFileFilter ? { zip_file: zipFileFilter } : undefined),
    [zipFileFilter]
  );

  const { data: rootFoldersResult } = useFindRootFoldersForSelectQuery({
    skip,
    variables: {
      zip_file_filter: zipFileFilter,
    },
  });

  const { data: queryFoldersResult } = useFindFoldersForQueryQuery({
    skip: !query,
    variables: {
      filter: { q: query, per_page: 200 },
      folder_filter: folderFilterForQuery,
    },
  });

  const { data: initialSelectedResult } = useFindFolderHierarchyForIDsQuery({
    skip: !initialSelected || cachedInitialSelected.length === 0,
    variables: {
      ids: cachedInitialSelected ?? [],
    },
  });

  const rootFolders: IFolder[] = useMemo(() => {
    const ret = rootFoldersResult?.findFolders.folders ?? [];
    return ret.map((f) => ({ ...f, expanded: false, children: undefined }));
  }, [rootFoldersResult]);

  const initialSelectedFolders: IFolder[] = useMemo(() => {
    const ret: IFolder[] = [];
    (initialSelectedResult?.findFolders.folders ?? []).forEach((folder) => {
      if (!folder.parent_folders.length) {
        // add root folder if not present
        if (!ret.find((f) => f.id === folder.id)) {
          ret.push({ ...folder, expanded: true, children: [] });
        }
        return;
      }

      let currentParent: IFolder | undefined;

      for (let i = folder.parent_folders.length - 1; i >= 0; i--) {
        const thisFolder = folder.parent_folders[i];
        let existing: IFolder | undefined;

        if (i === folder.parent_folders.length - 1) {
          // last parent, add the folder as root if not present
          existing = ret.find((f) => f.id === thisFolder.id);
          if (!existing) {
            existing = {
              ...folder.parent_folders[i],
              expanded: true,
              children: folder.parent_folders[i].sub_folders
                // filter out zip folders if needed
                .filter((f) => f.zip_file === null || !excludeZipFolders)
                .map((f) => ({
                  ...f,
                  expanded: false,
                  children: undefined,
                })),
            };
            ret.push(existing);
          }
          currentParent = existing;
          continue;
        }

        const existingIndex =
          currentParent!.children?.findIndex((f) => f.id === thisFolder.id) ??
          -1;
        if (existingIndex === -1) {
          // should be guaranteed
          throw new Error(
            `Parent folder ${thisFolder.id} not found in children of ${
              currentParent!.id
            }`
          );
        }

        existing = currentParent!.children![existingIndex];

        // replace children
        existing = {
          ...existing,
          expanded: true,
          // filter out zip folders if needed
          children: thisFolder.sub_folders
            .filter((f) => f.zip_file === null || !excludeZipFolders)
            .map((f) => ({
              ...f,
              expanded: false,
              children: undefined,
            })),
        };

        currentParent!.children![existingIndex] = existing;
        currentParent = existing;
      }
    });
    return ret;
  }, [initialSelectedResult, excludeZipFolders]);

  const mergedRootFolders = useMemo(() => {
    if (query) {
      return rootFolders;
    }

    return mergeFolderMaps(rootFolders, initialSelectedFolders);
  }, [rootFolders, initialSelectedFolders, query]);

  const queryFolders: IFolder[] = useMemo(() => {
    // construct the folder list from the query result
    const ret: IFolder[] = [];

    (queryFoldersResult?.findFolders.folders ?? []).forEach((folder) => {
      if (!folder.parent_folders.length) {
        // no parents, just add it if not present
        if (!ret.find((f) => f.id === folder.id)) {
          ret.push({ ...folder, expanded: true, children: [] });
        }
        return;
      }

      // expand the parent folders
      let currentParent: IFolder | undefined;
      for (let i = folder.parent_folders.length - 1; i >= 0; i--) {
        const thisFolder = folder.parent_folders[i];
        let existing: IFolder | undefined;

        if (i === folder.parent_folders.length - 1) {
          // last parent, add the folder as root
          existing = ret.find((f) => f.id === thisFolder.id);
          if (!existing) {
            existing = {
              ...folder.parent_folders[i],
              expanded: true,
              children: [],
            };
            ret.push(existing);
          }
          currentParent = existing;
          continue;
        }

        // find folder in current parent's children
        // currentParent is guaranteed to be defined here
        existing = currentParent!.children?.find((f) => f.id === thisFolder.id);
        if (!existing) {
          // add to current parent's children
          existing = {
            ...thisFolder,
            expanded: true,
            children: [],
          };
          currentParent!.children!.push(existing);
        }
        currentParent = existing;
      }

      if (!currentParent) {
        return;
      }

      if (!currentParent.children) {
        currentParent.children = [];
      }

      // currentParent is now the immediate parent folder
      currentParent!.children!.push({
        ...folder,
        expanded: false,
        children: undefined,
      });
    });
    return ret;
  }, [queryFoldersResult]);

  const [folderMap, setFolderMap] = React.useState<IFolder[]>([]);

  useEffect(() => {
    if (!query) {
      setFolderMap(mergedRootFolders);
    } else {
      setFolderMap(queryFolders);
    }
  }, [query, mergedRootFolders, queryFolders]);

  async function onToggleExpanded(folder: IFolder) {
    setFolderMap(folderMap.map(toggleExpandedFn(folder)));

    // query children folders if not already loaded
    if (folder.children === undefined) {
      const subFolderResult = await queryFindSubFolders(
        folder.id,
        excludeZipFolders
      );
      setFolderMap((current) =>
        current.map(
          replaceFolder({
            ...folder,
            expanded: true,
            children: subFolderResult.data.findFolders.folders.map((f) => ({
              ...f,
              expanded: false,
            })),
          })
        )
      );
    }
  }

  return { folderMap, onToggleExpanded };
}

function getMatchingFolders(folders: IFolder[], query: string): IFolder[] {
  let matches: IFolder[] = [];

  const queryLower = query.toLowerCase();

  folders.forEach((folder) => {
    if (
      folder.basename.toLowerCase().includes(queryLower) ||
      folder.path.toLowerCase() === queryLower
    ) {
      matches.push(folder);
    }

    if (folder.children) {
      matches = matches.concat(getMatchingFolders(folder.children, query));
    }
  });

  return matches;
}

export const FolderSelector: React.FC<{
  onSelect: (folder: IFolder, exclude?: boolean) => void;
  canExclude?: boolean;
  preListContent?: React.ReactNode;
  folderMap: IFolder[];
  onToggleExpanded: (folder: IFolder) => void;
}> = ({
  onSelect,
  preListContent,
  canExclude = false,
  folderMap,
  onToggleExpanded,
}) => {
  return (
    <ul className="selectable-list">
      {preListContent}
      {folderMap.map((folder) => (
        <FolderRow
          key={folder.id}
          folder={folder}
          onSelect={(f, exclude) => onSelect(f, exclude)}
          toggleExpanded={onToggleExpanded}
          canExclude={canExclude}
        />
      ))}
    </ul>
  );
};

interface IInputFilterProps {
  criterion: FolderCriterion;
  setCriterion: (c: FolderCriterion) => void;
  mode?: FilterMode;
}

export const FolderFilter: React.FC<IInputFilterProps> = ({
  criterion,
  setCriterion,
  mode,
}) => {
  const intl = useIntl();
  const [query, setQuery] = useState("");
  const [displayQuery, onQueryChange] = useDebouncedState(query, setQuery, 250);

  const { folderMap, onToggleExpanded } = useFolderMap({ query, mode });

  const messages = defineMessages({
    sub_folder_depth: {
      id: "sub_folder_depth",
      defaultMessage: "Levels (empty for all)",
    },
  });

  function criterionOptionTypeToIncludeID(): string {
    return "include-sub-folders";
  }

  function criterionOptionTypeToIncludeUIString(): MessageDescriptor {
    const optionType = "include_sub_folders";

    return {
      id: optionType,
    };
  }

  function onDepthChanged(depth: number) {
    // this could be ParentFolderCriterion, but the types are the same
    const newValue = criterion.clone() as FolderCriterion;
    newValue.value.depth = depth;
    setCriterion(newValue);
  }

  function onSelect(folder: IFolder, exclude: boolean = false) {
    // toggle selection
    const newValue = criterion.clone() as FolderCriterion;

    if (!exclude) {
      if (newValue.value.items.find((i) => i.id === folder.id)) {
        return;
      }

      newValue.value.items.push({ id: folder.id, label: folder.path });
    } else {
      if (newValue.value.excluded.find((i) => i.id === folder.id)) {
        return;
      }

      newValue.value.excluded.push({ id: folder.id, label: folder.path });
    }

    setCriterion(newValue);
  }

  const onUnselect = useCallback(
    (i: Option, excluded?: boolean) => {
      const newValue = criterion.clone() as FolderCriterion;

      if (!excluded) {
        newValue.value.items = newValue.value.items.filter(
          (item) => item.id !== i.id
        );
      } else {
        newValue.value.excluded = newValue.value.excluded.filter(
          (item) => item.id !== i.id
        );
      }
      setCriterion(newValue);
    },
    [criterion, setCriterion]
  );

  function onEnter() {
    if (!query) return;

    // if there is a single folder that matches the query, select it
    const matchingFolders = getMatchingFolders(folderMap, query);
    if (matchingFolders.length === 1) {
      onSelect(matchingFolders[0]);
    }
  }

  const selectedList = useMemo(() => {
    const selected: Option[] =
      criterion.value?.items.map((item) => ({
        id: item.id,
        label: item.label,
      })) ?? [];

    return <SelectedList items={selected} onUnselect={onUnselect} />;
  }, [criterion, onUnselect]);

  const excludedList = useMemo(() => {
    const selected: Option[] =
      criterion.value?.excluded.map((item) => ({
        id: item.id,
        label: item.label,
      })) ?? [];

    return (
      <SelectedList
        excluded
        items={selected}
        onUnselect={(i) => onUnselect(i, true)}
      />
    );
  }, [criterion, onUnselect]);

  return (
    <div className="folder-filter">
      <DepthSelector
        depth={criterion.value.depth}
        onDepthChanged={onDepthChanged}
        id={criterionOptionTypeToIncludeID()}
        label={intl.formatMessage(criterionOptionTypeToIncludeUIString())}
        placeholder={intl.formatMessage(messages.sub_folder_depth)}
      />

      <Form.Group>
        {selectedList}
        {excludedList}
        <ClearableInput
          value={displayQuery}
          setValue={(v) => onQueryChange(v)}
          placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
          onEnter={onEnter}
        />
        <FolderSelector
          folderMap={folderMap}
          onToggleExpanded={onToggleExpanded}
          onSelect={onSelect}
          canExclude
        />
      </Form.Group>
    </div>
  );
};

export const SidebarFolderFilter: React.FC<
  ISidebarSectionProps & {
    filter: ListFilterModel;
    setFilter: (f: ListFilterModel) => void;
    criterionOption?: ModifierCriterionOption;
  }
> = (props) => {
  const intl = useIntl();
  const [skip, setSkip] = useState(true);
  const [query, setQuery] = useState("");
  const [displayQuery, onQueryChange] = useDebouncedState(query, setQuery, 250);

  function onOpen() {
    setSkip(false);
    props.onOpen?.();
  }

  const option = props.criterionOption ?? FolderCriterionOption;
  const { filter, setFilter } = props;

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as FolderCriterion;

    const newCriterion = filter.makeCriterion(option.type) as FolderCriterion;
    return newCriterion;
  }, [option.type, filter]);

  const subDirsSelected = criterion.value?.depth === -1;

  // if there are multiple values or excluded values, then we show none of the
  // current values
  const multipleSelected =
    criterion.value.items.length > 1 || criterion.value.excluded.length > 0;

  const { folderMap, onToggleExpanded } = useFolderMap({
    query,
    skip,
    initialSelected: criterion.value.items.map((i) => i.id),
    mode: filter.mode,
  });

  function onSelect(folder: IFolder) {
    // maintain sub-folder select if present
    const depth = subDirsSelected ? -1 : 0;

    const c = criterion.clone() as FolderCriterion;
    c.value = {
      items: [{ id: folder.id, label: folder.path }],
      depth,
      excluded: [],
    };

    const newCriteria = props.filter.criteria.filter(
      (cc) => cc.criterionOption.type !== option.type
    );

    if (c.isValid()) newCriteria.push(c);

    setFilter(props.filter.setCriteria(newCriteria));
  }

  function onSelectSubfolders() {
    const c = criterion.clone() as FolderCriterion;
    c.value = {
      items: c.value?.items ?? [],
      depth: -1,
      excluded: c.value?.excluded ?? [],
    };

    setFilter(props.filter.replaceCriteria(option.type, [c]));
  }

  const onUnselect = useCallback(
    (i: Option) => {
      if (i.className === "modifier-object") {
        // subfolders option
        const c = criterion.clone() as FolderCriterion;
        c.value = {
          items: c.value?.items ?? [],
          depth: 0,
          excluded: c.value?.excluded ?? [],
        };

        setFilter(props.filter.replaceCriteria(option.type, [c]));
        return;
      }

      setFilter(props.filter.removeCriterion(option.type));
    },
    [props.filter, setFilter, option.type, criterion]
  );

  function onEnter() {
    if (!query) return;

    // if there is a single folder that matches the query, select it
    const matchingFolders = getMatchingFolders(folderMap, query);
    if (matchingFolders.length === 1) {
      onSelect(matchingFolders[0]);
    }
  }

  const selectedList = useMemo(() => {
    if (multipleSelected) {
      return null;
    }

    const selected: Option[] =
      criterion.value?.items.map((item) => ({
        id: item.id,
        label: item.label,
      })) ?? [];

    if (subDirsSelected) {
      selected.push({
        id: "subfolders",
        label: "(" + intl.formatMessage({ id: "sub_folders" }) + ")",
        className: "modifier-object",
      });
    }

    return <SelectedList items={selected} onUnselect={onUnselect} />;
  }, [intl, multipleSelected, subDirsSelected, criterion, onUnselect]);

  const modifierItem = criterion.value.items.length > 0 &&
    !multipleSelected &&
    !subDirsSelected && (
      <li className="unselected-object modifier-object">
        <a onClick={onSelectSubfolders}>
          <span>
            <Icon className={`fa-fw include-button`} icon={faPlus} />
            (<FormattedMessage id="sub_folders" />)
          </span>
        </a>
      </li>
    );

  return (
    <SidebarSection
      {...props}
      outsideCollapse={selectedList}
      onOpen={onOpen}
      className="sidebar-list-filter sidebar-folder-filter"
    >
      <ClearableInput
        value={displayQuery}
        setValue={(v) => onQueryChange(v)}
        placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
        onEnter={onEnter}
      />

      <FolderSelector
        folderMap={folderMap}
        onToggleExpanded={onToggleExpanded}
        preListContent={modifierItem}
        onSelect={(f) => onSelect(f)}
      />
    </SidebarSection>
  );
};
