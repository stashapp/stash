import React, { useCallback, useEffect, useMemo } from "react";
import {
  FolderDataFragment,
  useFindRootFoldersForSelectQuery,
} from "src/core/generated-graphql";
import {
  ISidebarSectionProps,
  SidebarSection,
} from "src/components/Shared/Sidebar";
import {
  faChevronDown,
  faChevronRight,
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
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";

interface IFolder extends FolderDataFragment {
  children?: IFolder[];
  expanded: boolean;
}

const FolderRow: React.FC<{
  folder: IFolder;
  level?: number;
  toggleExpanded: (folder: IFolder) => void;
  onSelect: (folder: IFolder) => void;
}> = ({ folder, level, toggleExpanded, onSelect }) => {
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

export const FolderSelector: React.FC<{
  onSelect: (folder: IFolder) => void;
  preListContent?: React.ReactNode;
  skip?: boolean;
}> = ({ onSelect, preListContent, skip = false }) => {
  const { data: rootFoldersResult } = useFindRootFoldersForSelectQuery({
    skip,
  });

  const rootFolders: IFolder[] = useMemo(() => {
    const ret = rootFoldersResult?.findFolders.folders ?? [];
    return ret.map((f) => ({ ...f, expanded: false, children: undefined }));
  }, [rootFoldersResult]);

  const [folderMap, setFolderMap] = React.useState<IFolder[]>([]);

  useEffect(() => {
    setFolderMap(rootFolders);
  }, [rootFolders]);

  async function onToggleExpanded(folder: IFolder) {
    setFolderMap(folderMap.map(toggleExpandedFn(folder)));

    // query children folders if not already loaded
    if (folder.children === undefined) {
      const subFolderResult = await queryFindSubFolders(folder.id);
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

  return (
    <ul>
      {preListContent}
      {folderMap.map((folder) => (
        <FolderRow
          key={folder.id}
          folder={folder}
          onSelect={(f) => onSelect(f)}
          toggleExpanded={onToggleExpanded}
        />
      ))}
    </ul>
  );
};

export const SidebarFolderFilter: React.FC<
  ISidebarSectionProps & {
    filter: ListFilterModel;
    setFilter: (f: ListFilterModel) => void;
  }
> = (props) => {
  const intl = useIntl();
  const [skip, setSkip] = React.useState(true);

  function onOpen() {
    setSkip(false);
    props.onOpen?.();
  }

  const option = FolderCriterionOption;
  const { filter, setFilter } = props;

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as FolderCriterion;

    const newCriterion = filter.makeCriterion(option.type) as FolderCriterion;
    return newCriterion;
  }, [option.type, filter]);

  function onSelect(folder: IFolder) {
    const c = criterion.clone() as FolderCriterion;
    c.value = {
      items: [{ id: folder.id, label: folder.path }],
      depth: 0,
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

  const subDirsSelected = criterion.value?.depth === -1;

  const selectedList = useMemo(() => {
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
  }, [intl, subDirsSelected, criterion, onUnselect]);

  const modifierItem = criterion.value.items.length > 0 && !subDirsSelected && (
    <li className="unselected-object modifier-object">
      <a onClick={onSelectSubfolders}>
        <Icon className={`fa-fw include-button`} icon={faPlus} />
        (<FormattedMessage id="sub_folders" />)
      </a>
    </li>
  );

  return (
    <SidebarSection
      {...props}
      outsideCollapse={selectedList}
      onOpen={onOpen}
      className="sidebar-list-filter sidebar-path-filter"
    >
      {/* query input goes here */}
      <FolderSelector
        skip={skip}
        preListContent={modifierItem}
        onSelect={(f) => onSelect(f)}
      />
    </SidebarSection>
  );
};
