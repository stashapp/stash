import React, { useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import { ScrapeDialogRow } from "src/components/Shared/ScrapeDialog/ScrapeDialogRow";
import { PerformerSelect } from "src/components/Performers/PerformerSelect";
import {
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import { TagIDSelect } from "src/components/Tags/TagSelect";
import { StudioSelect } from "src/components/Studios/StudioSelect";
import { GroupSelect } from "src/components/Groups/GroupSelect";
import { uniq } from "lodash-es";
import { CollapseButton } from "../CollapseButton";
import { Badge, Button } from "react-bootstrap";
import { Icon } from "../Icon";
import { faLink, faPlus } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";

interface INewScrapedObjects<T> {
  newValues: T[];
  onCreateNew: (value: T) => void;
  onLinkExisting?: (value: T) => void;
  getName: (value: T) => string;
}

export const NewScrapedObjects = <T,>(props: INewScrapedObjects<T>) => {
  const intl = useIntl();

  if (props.newValues.length === 0) {
    return null;
  }

  const ret = (
    <>
      {props.newValues.map((t) => (
        <Badge
          className="tag-item"
          variant="secondary"
          key={props.getName(t)}
          onClick={() => props.onCreateNew(t)}
        >
          {props.getName(t)}
          <Button className="minimal ml-2">
            <Icon className="fa-fw" icon={faPlus} />
          </Button>
          {props.onLinkExisting ? (
            <Button
              className="minimal"
              onClick={(e) => {
                props.onLinkExisting?.(t);
                e.stopPropagation();
              }}
            >
              <Icon className="fa-fw" icon={faLink} />
            </Button>
          ) : null}
        </Badge>
      ))}
    </>
  );

  const minCollapseLength = 10;

  if (props.newValues!.length >= minCollapseLength) {
    const missingText = intl.formatMessage({
      id: "dialogs.scrape_results_missing",
    });
    return (
      <CollapseButton text={`${missingText} (${props.newValues!.length})`}>
        {ret}
      </CollapseButton>
    );
  }

  return ret;
};

interface IScrapedStudioRow {
  title: string;
  field: string;
  result: ObjectScrapeResult<GQL.ScrapedStudio>;
  onChange: (value: ObjectScrapeResult<GQL.ScrapedStudio>) => void;
  newStudio?: GQL.ScrapedStudio;
  onCreateNew?: (value: GQL.ScrapedStudio) => void;
  onLinkExisting?: (value: GQL.ScrapedStudio) => void;
}

function getObjectName<T extends { name: string }>(value: T) {
  return value.name;
}

export const ScrapedStudioRow: React.FC<IScrapedStudioRow> = ({
  title,
  field,
  result,
  onChange,
  newStudio,
  onCreateNew,
  onLinkExisting,
}) => {
  function renderScrapedStudio(
    scrapeResult: ObjectScrapeResult<GQL.ScrapedStudio>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedStudio) => void
  ) {
    const resultValue = isNew
      ? scrapeResult.newValue
      : scrapeResult.originalValue;
    const value = resultValue ? [resultValue] : [];

    const selectValue = value.map((p) => {
      const aliases: string[] = p.aliases
        ? p.aliases.split(",").map((a) => a.trim())
        : [];
      return {
        id: p.stored_id ?? "",
        name: p.name ?? "",
        aliases,
      };
    });

    return (
      <StudioSelect
        className="form-control react-select"
        isDisabled={!isNew}
        onSelect={(items) => {
          if (onChangeFn) {
            const { id, aliases, ...data } = items[0];
            onChangeFn({
              ...data,
              stored_id: id,
              aliases: aliases?.join(", "),
            });
          }
        }}
        values={selectValue}
      />
    );
  }

  return (
    <ScrapeDialogRow
      title={title}
      field={field}
      result={result}
      originalField={renderScrapedStudio(result)}
      newField={renderScrapedStudio(result, true, (value) =>
        onChange(result.cloneWithValue(value))
      )}
      onChange={onChange}
      newValues={
        newStudio && onCreateNew ? (
          <NewScrapedObjects
            newValues={[newStudio]}
            onCreateNew={onCreateNew}
            getName={getObjectName}
            onLinkExisting={onLinkExisting}
          />
        ) : undefined
      }
    />
  );
};

interface IScrapedObjectsRow<T> {
  title: string;
  field: string;
  result: ScrapeResult<T[]>;
  onChange: (value: ScrapeResult<T[]>) => void;
  newObjects?: T[];
  onCreateNew?: (value: T) => void;
  onLinkExisting?: (value: T) => void;
  renderObjects: (
    result: ScrapeResult<T[]>,
    isNew?: boolean,
    onChange?: (value: T[]) => void
  ) => JSX.Element;
  getName: (value: T) => string;
}

export const ScrapedObjectsRow = <T,>(props: IScrapedObjectsRow<T>) => {
  const {
    title,
    field,
    result,
    onChange,
    newObjects = [],
    onCreateNew,
    onLinkExisting,
    renderObjects,
    getName,
  } = props;

  return (
    <ScrapeDialogRow
      title={title}
      field={field}
      result={result}
      originalField={renderObjects(result)}
      newField={renderObjects(result, true, (value) =>
        onChange(result.cloneWithValue(value))
      )}
      onChange={onChange}
      newValues={
        onCreateNew && newObjects.length > 0 ? (
          <NewScrapedObjects
            newValues={newObjects ?? []}
            onCreateNew={onCreateNew}
            onLinkExisting={onLinkExisting}
            getName={getName}
          />
        ) : undefined
      }
    />
  );
};

type IScrapedObjectRowImpl<T> = Omit<
  IScrapedObjectsRow<T>,
  "renderObjects" | "getName"
>;

export const ScrapedPerformersRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedPerformer> & { ageFromDate?: string | null }
> = ({
  title,
  field,
  result,
  onChange,
  newObjects,
  onCreateNew,
  ageFromDate,
  onLinkExisting,
}) => {
  const performersCopy = useMemo(() => {
    return (
      newObjects?.map((p) => {
        const name: string = p.name ?? "";
        return { ...p, name };
      }) ?? []
    );
  }, [newObjects]);

  function renderScrapedPerformers(
    scrapeResult: ScrapeResult<GQL.ScrapedPerformer[]>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedPerformer[]) => void
  ) {
    const resultValue = isNew
      ? scrapeResult.newValue
      : scrapeResult.originalValue;
    const value = resultValue ?? [];

    const selectValue = value.map((p) => {
      const alias_list: string[] = [];
      return {
        id: p.stored_id ?? "",
        name: p.name ?? "",
        alias_list,
      };
    });

    return (
      <PerformerSelect
        isMulti
        className="form-control"
        isDisabled={!isNew}
        onSelect={(items) => {
          if (onChangeFn) {
            // map the id back to stored_id
            onChangeFn(items.map((p) => ({ ...p, stored_id: p.id })));
          }
        }}
        values={selectValue}
        ageFromDate={ageFromDate}
      />
    );
  }

  return (
    <ScrapedObjectsRow<GQL.ScrapedPerformer>
      title={title}
      field={field}
      result={result}
      renderObjects={renderScrapedPerformers}
      onChange={onChange}
      newObjects={performersCopy}
      onCreateNew={onCreateNew}
      getName={(value) => value.name ?? ""}
      onLinkExisting={onLinkExisting}
    />
  );
};

export const ScrapedGroupsRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedGroup>
> = ({
  title,
  field,
  result,
  onChange,
  newObjects,
  onCreateNew,
  onLinkExisting,
}) => {
  const groupsCopy = useMemo(() => {
    return (
      newObjects?.map((p) => {
        const name: string = p.name ?? "";
        return { ...p, name };
      }) ?? []
    );
  }, [newObjects]);

  function renderScrapedGroups(
    scrapeResult: ScrapeResult<GQL.ScrapedGroup[]>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedGroup[]) => void
  ) {
    const resultValue = isNew
      ? scrapeResult.newValue
      : scrapeResult.originalValue;
    const value = resultValue ?? [];

    const selectValue = value.map((p) => {
      const aliases: string = "";
      return {
        id: p.stored_id ?? "",
        name: p.name ?? "",
        aliases,
      };
    });

    return (
      <GroupSelect
        isMulti
        className="form-control react-select"
        isDisabled={!isNew}
        onSelect={(items) => {
          if (onChangeFn) {
            // map the id back to stored_id
            onChangeFn(items.map((p) => ({ ...p, stored_id: p.id })));
          }
        }}
        values={selectValue}
      />
    );
  }

  return (
    <ScrapedObjectsRow<GQL.ScrapedGroup>
      title={title}
      field={field}
      result={result}
      renderObjects={renderScrapedGroups}
      onChange={onChange}
      newObjects={groupsCopy}
      onCreateNew={onCreateNew}
      getName={(value) => value.name ?? ""}
      onLinkExisting={onLinkExisting}
    />
  );
};

export const ScrapedTagsRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedTag>
> = ({
  title,
  field,
  result,
  onChange,
  newObjects,
  onCreateNew,
  onLinkExisting,
}) => {
  function renderScrapedTags(
    scrapeResult: ScrapeResult<GQL.ScrapedTag[]>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedTag[]) => void
  ) {
    const resultValue = isNew
      ? scrapeResult.newValue
      : scrapeResult.originalValue;
    const value = resultValue ?? [];

    const selectValue = uniq(value.map((p) => p.stored_id ?? ""));

    // we need to use TagIDSelect here because we want to use the local name
    // of the tag instead of the name from the source
    return (
      <TagIDSelect
        isMulti
        className="form-control"
        isDisabled={!isNew}
        onSelect={(items) => {
          if (onChangeFn) {
            // map the id back to stored_id
            onChangeFn(items.map((p) => ({ ...p, stored_id: p.id })));
          }
        }}
        ids={selectValue}
      />
    );
  }

  return (
    <ScrapedObjectsRow<GQL.ScrapedTag>
      title={title}
      field={field}
      result={result}
      renderObjects={renderScrapedTags}
      onChange={onChange}
      newObjects={newObjects}
      onCreateNew={onCreateNew}
      onLinkExisting={onLinkExisting}
      getName={getObjectName}
    />
  );
};
