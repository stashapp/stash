import React from "react";
import { Button, ButtonGroup } from "react-bootstrap";

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  onChangePage: (page: number) => void;
}

export const Pagination: React.FC<IPaginationProps> = ({
  itemsPerPage,
  currentPage,
  totalItems,
  onChangePage
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
    i => startPage + i
  );

  const pageButtons = pages.map((page: number) => (
    <Button
      key={page}
      active={currentPage === page}
      onClick={() => onChangePage(page)}
    >
      {page}
    </Button>
  ));

  return (
    <ButtonGroup className="filter-container">
      <Button disabled={currentPage === 1} onClick={() => onChangePage(1)}>
        First
      </Button>
      <Button
        disabled={currentPage === 1}
        onClick={() => onChangePage(currentPage - 1)}
      >
        Previous
      </Button>
      {pageButtons}
      <Button
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(currentPage + 1)}
      >
        Next
      </Button>
      <Button
        disabled={currentPage === totalPages}
        onClick={() => onChangePage(totalPages)}
      >
        Last
      </Button>
    </ButtonGroup>
  );
};
