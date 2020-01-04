import React from "react";
import { Button, ButtonGroup } from 'react-bootstrap';

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  onChangePage: (page: number) => void;
}

interface IPaginationState {
  pages: number[];
  totalPages: number;
}

export class Pagination extends React.Component<IPaginationProps, IPaginationState> {
  constructor(props: IPaginationProps) {
    super(props);
    this.state = {
      pages: [],
      totalPages: Number.MAX_SAFE_INTEGER,
    };
  }

  public componentWillMount() {
    this.setPage(this.props.currentPage, false);
  }

  public componentDidUpdate(prevProps: IPaginationProps) {
    if (this.props.totalItems !== prevProps.totalItems || this.props.itemsPerPage !== prevProps.itemsPerPage) {
      this.setPage(this.props.currentPage);
    }
  }

  public render() {
    if (!this.state || !this.state.pages || this.state.pages.length <= 1) { return null; }

    return (
      <ButtonGroup className="filter-container">
        <Button
          disabled={this.props.currentPage === 1}
          onClick={() => this.setPage(1)}
        >First</Button>
        <Button
          disabled={this.props.currentPage === 1}
          onClick={() => this.setPage(this.props.currentPage - 1)}
        >Previous</Button>
        {this.renderPageButtons()}
        <Button
          disabled={this.props.currentPage === this.state.totalPages}
          onClick={() => this.setPage(this.props.currentPage + 1)}
        >Next</Button>
        <Button
          disabled={this.props.currentPage === this.state.totalPages}
          onClick={() => this.setPage(this.state.totalPages)}
        >Last</Button>
      </ButtonGroup>
    );
  }

  private renderPageButtons() {
    return this.state.pages.map((page: number, index: number) => (
      <Button
        key={index}
        active={this.props.currentPage === page}
        onClick={() => this.setPage(page)}
      >{page}</Button>
    ));
  }

  private setPage(page?: number, propagate: boolean = true) {
    if (page === undefined) { return; }

    const pagerState = this.getPagerState(this.props.totalItems, page, this.props.itemsPerPage);

    if (page < 1) { page = 1; }
    if (page > pagerState.totalPages) { page = pagerState.totalPages; }

    this.setState(pagerState);
    if (propagate) { this.props.onChangePage(page); }
  }

  private getPagerState(totalItems: number, currentPage: number, pageSize: number) {
    const totalPages = Math.ceil(totalItems / pageSize);

    let startPage: number;
    let endPage: number;
    if (totalPages <= 10) {
      // less than 10 total pages so show all
      startPage = 1;
      endPage = totalPages;
    } else {
      // more than 10 total pages so calculate start and end pages
      if (currentPage <= 6) {
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

    // create an array of pages numbers
    const pages = [...Array((endPage + 1) - startPage).keys()].map((i) => startPage + i);

    return {
      pages,
      totalPages,
    };
  }
}
