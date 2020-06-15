import React from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedNumber, useIntl } from "react-intl";

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  onChangePage: (page: number) => void;
}

interface IPaginationIndexProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
}

export const Pagination: React.FC<IPaginationProps> = ({
  itemsPerPage,
  currentPage,
  totalItems,
  onChangePage,
}) => {
  const totalPages = Math.ceil(totalItems / itemsPerPage);

  let startPage: number;
  let endPage: number;
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

  const pages = [...Array(endPage + 1 - startPage).keys()].map(
    (i) => startPage + i
  );

  const calculatePageClass = (buttonPage: number) => {
    if (pages.length <= 4) return "";

    if (currentPage === 1 && buttonPage <= 4) return "";
    const maxPage = pages[pages.length - 1];
    if (currentPage === maxPage && buttonPage > maxPage - 3) return "";
    if (Math.abs(buttonPage - currentPage) <= 1) return "";
    return "d-none d-sm-block";
  };

  const pageButtons = pages.map((page: number) => (
    <Button
      variant="secondary"
      className={calculatePageClass(page)}
      key={page}
      active={currentPage === page}
      onClick={() => onChangePage(page)}
    >
      <FormattedNumber value={page} />
    </Button>
  ));

  if (pages.length <= 1) return <div />;

  return (
    <ButtonGroup className="filter-container pagination">
      <Button
        variant="secondary"
        disabled={currentPage === 1}
        onClick={() => onChangePage(1)}
      >
        <span className="d-none d-sm-inline">First</span>
        <span className="d-inline d-sm-none">&#x300a;</span>
      </Button>
      <Button
        className="d-none d-sm-block"
        variant="secondary"
        disabled={currentPage === 1}
        onClick={() => onChangePage(currentPage - 1)}
      >
        Previous
      </Button>
      {pageButtons}
      <Button
        className="d-none d-sm-block"
        variant="secondary"
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(currentPage + 1)}
      >
        Next
      </Button>
      <Button
        variant="secondary"
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(totalPages)}
      >
        <span className="d-none d-sm-inline">Last</span>
        <span className="d-inline d-sm-none">&#x300b;</span>
      </Button>
    </ButtonGroup>
  );
};

export const PaginationIndex: React.FC<IPaginationIndexProps> = ({
  itemsPerPage,
  currentPage,
  totalItems,
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
    <span className="filter-container text-muted paginationIndex">
      {indexText}
    </span>
  );
};
