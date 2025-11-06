import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import {
  faCheckCircle,
  faMinus,
  faPlus,
  faTimesCircle,
} from "@fortawesome/free-solid-svg-icons";
import { faTimesCircle as faTimesCircleRegular } from "@fortawesome/free-regular-svg-icons";
import { ClearableInput } from "src/components/Shared/ClearableInput";
import { useIntl } from "react-intl";
import { keyboardClickHandler } from "src/utils/keyboard";
import { useDebounce } from "src/hooks/debounce";
import useFocus from "src/utils/focus";
import cx from "classnames";
import ScreenUtils from "src/utils/screen";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { TruncatedInlineText } from "src/components/Shared/TruncatedText";

interface ISelectedItem {
  className?: string;
  label: string;
  excluded?: boolean;
  onClick: () => void;
  // true if the object is a special modifier value
  modifier?: boolean;
}

const SelectedItem: React.FC<ISelectedItem> = ({
  className,
  label,
  excluded = false,
  onClick,
  modifier = false,
}) => {
  const iconClassName = excluded ? "exclude-icon" : "include-button";
  const spanClassName = excluded
    ? "excluded-object-label"
    : "selected-object-label";
  const [hovered, setHovered] = useState(false);

  const icon = useMemo(() => {
    if (!hovered) {
      return excluded ? faTimesCircle : faCheckCircle;
    }

    return faTimesCircleRegular;
  }, [hovered, excluded]);

  function onMouseOver() {
    setHovered(true);
  }

  function onMouseOut() {
    setHovered(false);
  }

  return (
    <li
      className={cx("selected-object", className, {
        "modifier-object": modifier,
      })}
    >
      <a
        onClick={() => onClick()}
        onKeyDown={keyboardClickHandler(onClick)}
        onMouseEnter={() => onMouseOver()}
        onMouseLeave={() => onMouseOut()}
        onFocus={() => onMouseOver()}
        onBlur={() => onMouseOut()}
        tabIndex={0}
      >
        <div className="label-group">
          <Icon className={`fa-fw ${iconClassName}`} icon={icon} />
          <TruncatedInlineText className={spanClassName} text={label} />
        </div>
      </a>
    </li>
  );
};

const CandidateItem: React.FC<{
  className?: string;
  onSelect: (exclude: boolean) => void;
  label: string;
  canExclude?: boolean;
  modifier?: boolean;
  singleValue?: boolean;
}> = ({
  onSelect,
  label,
  canExclude,
  modifier = false,
  singleValue = false,
  className,
}) => {
  const singleValueClass = singleValue ? "single-value" : "";
  const includeIcon = (
    <Icon
      className={`fa-fw include-button ${singleValueClass}`}
      icon={faPlus}
    />
  );
  const excludeIcon = (
    <Icon className={`fa-fw exclude-icon ${singleValueClass}`} icon={faMinus} />
  );

  return (
    <li
      className={cx("unselected-object", className, {
        "modifier-object": modifier,
      })}
    >
      <a
        onClick={() => onSelect(false)}
        onKeyDown={keyboardClickHandler(() => onSelect(false))}
        tabIndex={0}
      >
        <div className="label-group">
          {includeIcon}
          <TruncatedInlineText
            className="unselected-object-label"
            text={label}
          />
        </div>
        <div>
          {/* TODO item count */}
          {/* <span className="object-count">{p.id}</span> */}
          {canExclude && (
            <Button
              onClick={(e) => {
                e.stopPropagation();
                onSelect(true);
              }}
              onKeyDown={(e) => e.stopPropagation()}
              className="minimal exclude-button"
            >
              <span className="exclude-button-text">exclude</span>
              {excludeIcon}
            </Button>
          )}
        </div>
      </a>
    </li>
  );
};

export type Option<T = unknown> = {
  id: string;
  className?: string;
  value?: T;
  label: string;
  canExclude?: boolean; // defaults to true
};

export const SelectedList: React.FC<{
  items: Option[];
  onUnselect: (item: Option) => void;
  excluded?: boolean;
}> = ({ items, onUnselect, excluded }) => {
  if (items.length === 0) {
    return null;
  }

  return (
    <ul className={cx("selected-list", { "excluded-list": excluded })}>
      {items.map((p) => (
        <SelectedItem
          key={p.id}
          className={p.className}
          label={p.label}
          excluded={excluded}
          onClick={() => onUnselect(p)}
        />
      ))}
    </ul>
  );
};

const QueryField: React.FC<{
  focus: ReturnType<typeof useFocus>;
  value: string;
  setValue: (query: string) => void;
}> = ({ focus, value, setValue }) => {
  const intl = useIntl();

  const [displayQuery, setDisplayQuery] = useState(value);
  const debouncedSetQuery = useDebounce(setValue, 250);

  useEffect(() => {
    setDisplayQuery(value);
  }, [value]);

  const onQueryChange = useCallback(
    (input: string) => {
      setDisplayQuery(input);
      debouncedSetQuery(input);
    },
    [debouncedSetQuery, setDisplayQuery]
  );

  return (
    <ClearableInput
      focus={focus}
      value={displayQuery}
      setValue={(v) => onQueryChange(v)}
      placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
    />
  );
};

interface IQueryableProps {
  inputFocus?: ReturnType<typeof useFocus>;
  query?: string;
  setQuery?: (query: string) => void;
}

export const CandidateList: React.FC<
  {
    items: Option[];
    onSelect: (item: Option, exclude: boolean) => void;
    canExclude?: boolean;
    singleValue?: boolean;
  } & IQueryableProps
> = ({
  inputFocus,
  query,
  setQuery,
  items,
  onSelect,
  canExclude,
  singleValue,
}) => {
  const showQueryField =
    inputFocus !== undefined && query !== undefined && setQuery !== undefined;

  return (
    <div className="queryable-candidate-list">
      {showQueryField && (
        <QueryField
          focus={inputFocus}
          value={query}
          setValue={(v) => setQuery(v)}
        />
      )}
      <ul>
        {items.map((p) => (
          <CandidateItem
            key={p.id}
            className={p.className}
            onSelect={(exclude) => onSelect(p, exclude)}
            label={p.label}
            canExclude={canExclude && (p.canExclude ?? true)}
            singleValue={singleValue}
          />
        ))}
      </ul>
    </div>
  );
};

export const SidebarListFilter: React.FC<{
  title: React.ReactNode;
  selected: Option[];
  excluded?: Option[];
  candidates: Option[];
  singleValue?: boolean;
  onSelect: (item: Option, exclude: boolean) => void;
  onUnselect: (item: Option, exclude: boolean) => void;
  canExclude?: boolean;
  query?: string;
  setQuery?: (query: string) => void;
  preSelected?: React.ReactNode;
  postSelected?: React.ReactNode;
  preCandidates?: React.ReactNode;
  postCandidates?: React.ReactNode;
  onOpen?: () => void;
  // used to store open/closed state in SidebarStateContext
  sectionID?: string;
}> = ({
  title,
  selected,
  excluded,
  candidates,
  onSelect,
  onUnselect,
  canExclude,
  query,
  setQuery,
  singleValue = false,
  preCandidates,
  postCandidates,
  preSelected,
  postSelected,
  onOpen,
  sectionID,
}) => {
  // TODO - sort items?

  const inputFocus = useFocus();
  const [, setInputFocus] = inputFocus;

  function unselectHook(item: Option, exclude: boolean) {
    onUnselect(item, exclude);

    // focus the input box
    // don't do this on touch devices, as it's annoying
    if (!ScreenUtils.isTouch()) {
      setInputFocus();
    }
  }

  function selectHook(item: Option, exclude: boolean) {
    onSelect(item, exclude);

    // reset filter query after selecting
    setQuery?.("");

    // focus the input box
    // don't do this on touch devices, as it's annoying
    if (!ScreenUtils.isTouch()) {
      setInputFocus();
    }
  }

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      sectionID={sectionID}
      outsideCollapse={
        <>
          {preSelected ? <div className="extra">{preSelected}</div> : null}
          <SelectedList
            items={selected}
            onUnselect={(i) => unselectHook(i, false)}
          />
          {excluded && (
            <SelectedList
              items={excluded}
              onUnselect={(i) => unselectHook(i, true)}
              excluded
            />
          )}
          {postSelected ? <div className="extra">{postSelected}</div> : null}
        </>
      }
      onOpen={onOpen}
    >
      {preCandidates ? <div className="extra">{preCandidates}</div> : null}
      <CandidateList
        items={candidates}
        onSelect={selectHook}
        canExclude={canExclude}
        inputFocus={inputFocus}
        query={query}
        setQuery={setQuery}
        singleValue={singleValue}
      />
      {postCandidates ? <div className="extra">{postCandidates}</div> : null}
    </SidebarSection>
  );
};

export function useStaticResults<T>(r: T) {
  return () => ({ results: r, loading: false });
}
