import React, { PropsWithChildren, useMemo } from "react";
import { ApolloError, QueryResult } from "@apollo/client";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Pagination, PaginationIndex } from "./Pagination";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { FormattedMessage } from "react-intl";

export const LoadedContent: React.FC<
  PropsWithChildren<{
    loading?: boolean;
    error?: ApolloError;
  }>
> = ({ loading, error, children }) => {
  if (loading) {
    return <LoadingIndicator />;
  }
  if (error) {
    return (
      <ErrorMessage
        message={
          <FormattedMessage
            id="errors.loading_type"
            values={{ type: "items" }}
          />
        }
        error={error.message}
      />
    );
  }

  return <>{children}</>;
};

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
    return (
      <LoadedContent loading={result.loading} error={result.error}>
        {children}
        {!!pages && (
          <>
            {paginationIndex}
            {pagination}
          </>
        )}
      </LoadedContent>
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
