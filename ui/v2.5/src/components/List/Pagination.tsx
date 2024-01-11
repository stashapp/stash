import React, { useMemo } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  metadataByline?: React.ReactNode;
  pagesToShow?: number;
  onChangePage: (page: number) => void;
}

interface IPaginationIndexProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  metadataByline?: React.ReactNode;
}

export const Pagination: React.FC<IPaginationProps> = ({
  itemsPerPage,
  currentPage,
  totalItems,
  pagesToShow,
  onChangePage,
}) => {
  const totalPages = useMemo(
    () => Math.ceil(totalItems / itemsPerPage),
    [totalItems, itemsPerPage]
  );

  const pages = useMemo(() => {
    let startPage: number;
    let endPage: number;

    if (pagesToShow !== undefined) {
      startPage = Math.max(1, currentPage - Math.floor(pagesToShow / 2));
      endPage = Math.min(totalPages, startPage + pagesToShow - 1);

      if (endPage - startPage + 1 < pagesToShow) {
        startPage = Math.max(1, endPage - pagesToShow + 1);
      }
    } else {
      if (totalPages <= 10) {
        // less than 10 total pages so show all
        startPage = 1;
        endPage = totalPages;
      } else if (currentPage <= 6) {
        startPage = 1;
        endPage = 10;
      } else if (currentPage + 4 >= totalPages) {
        startPage = totalPages - 9;
        endPage = totalPages;
      } else {
        startPage = currentPage - 5;
        endPage = currentPage + 4;
      }
    }

    return [...Array(endPage + 1 - startPage).keys()].map((i) => startPage + i);
  }, [totalPages, currentPage, pagesToShow]);

  const pageButtons = useMemo(
    () =>
      pages.map((page: number) => {
        const calculatePageClass = (buttonPage: number) => {
          if (pages.length <= 4) return "";

          if (currentPage === 1 && buttonPage <= 4) return "";
          const maxPage = pages[pages.length - 1];
          if (currentPage === maxPage && buttonPage > maxPage - 3) return "";
          if (Math.abs(buttonPage - currentPage) <= 1) return "";
          return "d-none d-sm-block";
        };

        return (
          <Button
            variant="secondary"
            className={calculatePageClass(page)}
            key={page}
            active={currentPage === page}
            onClick={() => onChangePage(page)}
          >
            <FormattedNumber value={page} />
          </Button>
        );
      }),
    [pages, currentPage, onChangePage]
  );

  if (totalPages <= 1) return <div />;

  return (
    <ButtonGroup className="pagination">
      <Button
        variant="secondary"
        disabled={currentPage === 1}
        onClick={() => onChangePage(1)}
      >
        <span>«</span>
      </Button>
      <Button
        className="d-none d-sm-block"
        variant="secondary"
        disabled={currentPage === 1}
        onClick={() => onChangePage(currentPage - 1)}
      >
        &lt;
      </Button>
      {pageButtons}
      <Button
        className="d-none d-sm-block"
        variant="secondary"
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(currentPage + 1)}
      >
        &gt;
      </Button>
      <Button
        variant="secondary"
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(totalPages)}
      >
        <span>»</span>
      </Button>
    </ButtonGroup>
  );
};

export const PaginationIndex: React.FC<IPaginationIndexProps> = ({
  itemsPerPage,
  currentPage,
  totalItems,
  metadataByline,
}) => {
  const intl = useIntl();

  // Build the pagination index string
  const firstItemCount: number = Math.min(
    (currentPage - 1) * itemsPerPage + 1,
    totalItems
  );
  const lastItemCount: number = Math.min(
    firstItemCount + (itemsPerPage - 1),
    totalItems
  );
  const indexText: string = `${intl.formatNumber(
    firstItemCount
  )}-${intl.formatNumber(lastItemCount)} of ${intl.formatNumber(totalItems)}`;

  return (
    <span className="filter-container text-muted paginationIndex center-text">
      {indexText}
      <br />
      {metadataByline}
    </span>
  );
};
