import React, { useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import { ScrapeDialogRow } from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import { PerformerSelect } from "src/components/Performers/PerformerSelect";
import {
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import { TagSelect } from "src/components/Tags/TagSelect";
import { StudioSelect } from "src/components/Studios/StudioSelect";
import { GroupSelect } from "src/components/Movies/MovieSelect";

interface IScrapedStudioRow {
  title: string;
  result: ObjectScrapeResult<GQL.ScrapedStudio>;
  onChange: (value: ObjectScrapeResult<GQL.ScrapedStudio>) => void;
  newStudio?: GQL.ScrapedStudio;
  onCreateNew?: (value: GQL.ScrapedStudio) => void;
}

function getObjectName<T extends { name: string }>(value: T) {
  return value.name;
}

export const ScrapedStudioRow: React.FC<IScrapedStudioRow> = ({
  title,
  result,
  onChange,
  newStudio,
  onCreateNew,
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
      const aliases: string[] = [];
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
            const { id, ...data } = items[0];
            onChangeFn({
              ...data,
              stored_id: id,
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
      result={result}
      renderOriginalField={() => renderScrapedStudio(result)}
      renderNewField={() =>
        renderScrapedStudio(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
      newValues={newStudio ? [newStudio] : undefined}
      onCreateNew={() => {
        if (onCreateNew && newStudio) onCreateNew(newStudio);
      }}
      getName={getObjectName}
    />
  );
};

interface IScrapedObjectsRow<T> {
  title: string;
  result: ScrapeResult<T[]>;
  onChange: (value: ScrapeResult<T[]>) => void;
  newObjects?: T[];
  onCreateNew?: (value: T) => void;
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
    result,
    onChange,
    newObjects,
    onCreateNew,
    renderObjects,
    getName,
  } = props;

  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderObjects(result)}
      renderNewField={() =>
        renderObjects(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
      newValues={newObjects}
      onCreateNew={(i) => {
        if (onCreateNew) onCreateNew(newObjects![i]);
      }}
      getName={getName}
    />
  );
};

type IScrapedObjectRowImpl<T> = Omit<
  IScrapedObjectsRow<T>,
  "renderObjects" | "getName"
>;

export const ScrapedPerformersRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedPerformer>
> = ({ title, result, onChange, newObjects, onCreateNew }) => {
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
      />
    );
  }

  return (
    <ScrapedObjectsRow<GQL.ScrapedPerformer>
      title={title}
      result={result}
      renderObjects={renderScrapedPerformers}
      onChange={onChange}
      newObjects={performersCopy}
      onCreateNew={onCreateNew}
      getName={(value) => value.name ?? ""}
    />
  );
};

export const ScrapedGroupsRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedMovie>
> = ({ title, result, onChange, newObjects, onCreateNew }) => {
  const groupsCopy = useMemo(() => {
    return (
      newObjects?.map((p) => {
        const name: string = p.name ?? "";
        return { ...p, name };
      }) ?? []
    );
  }, [newObjects]);

  function renderScrapedGroups(
    scrapeResult: ScrapeResult<GQL.ScrapedMovie[]>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedMovie[]) => void
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
    <ScrapedObjectsRow<GQL.ScrapedMovie>
      title={title}
      result={result}
      renderObjects={renderScrapedGroups}
      onChange={onChange}
      newObjects={groupsCopy}
      onCreateNew={onCreateNew}
      getName={(value) => value.name ?? ""}
    />
  );
};

export const ScrapedTagsRow: React.FC<
  IScrapedObjectRowImpl<GQL.ScrapedTag>
> = ({ title, result, onChange, newObjects, onCreateNew }) => {
  function renderScrapedTags(
    scrapeResult: ScrapeResult<GQL.ScrapedTag[]>,
    isNew?: boolean,
    onChangeFn?: (value: GQL.ScrapedTag[]) => void
  ) {
    const resultValue = isNew
      ? scrapeResult.newValue
      : scrapeResult.originalValue;
    const value = resultValue ?? [];

    const selectValue = value.map((p) => {
      const aliases: string[] = [];
      return {
        id: p.stored_id ?? "",
        name: p.name ?? "",
        aliases,
      };
    });

    return (
      <TagSelect
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
      />
    );
  }

  return (
    <ScrapedObjectsRow<GQL.ScrapedTag>
      title={title}
      result={result}
      renderObjects={renderScrapedTags}
      onChange={onChange}
      newObjects={newObjects}
      onCreateNew={onCreateNew}
      getName={getObjectName}
    />
  );
};
