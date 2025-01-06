import React, { PropsWithChildren, useMemo } from "react";
import { QueryResult } from "@apollo/client";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Pagination, PaginationIndex } from "./Pagination";
import { LoadingIndicator } from "../Shared/LoadingIndicator";

export const PagedList: React.FC<
  PropsWithChildren<{
    result: QueryResult;
    cachedResult: QueryResult;
    filter: ListFilterModel;
    totalCount: number;
    onChangePage: (page: number) => void;
    metadataByline?: React.ReactNode;
  }>
> = ({
  result,
  cachedResult,
  filter,
  totalCount,
  onChangePage,
  metadataByline,
  children,
}) => {
  const pages = Math.ceil(totalCount / filter.itemsPerPage);

  const pagination = useMemo(() => {
    return (
      <Pagination
        itemsPerPage={filter.itemsPerPage}
        currentPage={filter.currentPage}
        totalItems={totalCount}
        metadataByline={metadataByline}
        onChangePage={onChangePage}
      />
    );
  }, [
    filter.itemsPerPage,
    filter.currentPage,
    totalCount,
    metadataByline,
    onChangePage,
  ]);

  const paginationIndex = useMemo(() => {
    if (cachedResult.loading) return;
    return (
      <PaginationIndex
        itemsPerPage={filter.itemsPerPage}
        currentPage={filter.currentPage}
        totalItems={totalCount}
        metadataByline={metadataByline}
      />
    );
  }, [
    cachedResult.loading,
    filter.itemsPerPage,
    filter.currentPage,
    totalCount,
    metadataByline,
  ]);

  const content = useMemo(() => {
    if (result.loading) {
      return <LoadingIndicator />;
    }
    if (result.error) {
      return <h1>{result.error.message}</h1>;
    }

    return (
      <>
        {children}
        {!!pages && (
          <>
            {paginationIndex}
            {pagination}
          </>
        )}
      </>
    );
  }, [
    result.loading,
    result.error,
    pages,
    children,
    pagination,
    paginationIndex,
  ]);

  return (
    <>
      {pagination}
      {paginationIndex}
      {content}
    </>
  );
};
