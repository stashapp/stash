import { Button, ButtonGroup } from "@blueprintjs/core";
import React from "react";

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
      <ButtonGroup large={true} className="filter-container">
        <Button
          text="First"
          disabled={this.props.currentPage === 1}
          onClick={() => this.setPage(1)}
        />
        <Button
          text="Previous"
          disabled={this.props.currentPage === 1}
          onClick={() => this.setPage(this.props.currentPage - 1)}
        />
        {this.renderPageButtons()}
        <Button
          text="Next"
          disabled={this.props.currentPage === this.state.totalPages}
          onClick={() => this.setPage(this.props.currentPage + 1)}
        />
        <Button
          text="Last"
          disabled={this.props.currentPage === this.state.totalPages}
          onClick={() => this.setPage(this.state.totalPages)}
        />
      </ButtonGroup>
    );
  }

  private renderPageButtons() {
    return this.state.pages.map((page: number, index: number) => (
      <Button
        key={index}
        text={page}
        active={this.props.currentPage === page}
        onClick={() => this.setPage(page)}
      />
    ));
  }

  private setPage(page?: number, propagate: boolean = true) {
    if (page === undefined) { return; }

    const pagerState = this.getPagerState(this.props.totalItems, page, this.props.itemsPerPage);

    // rearranged this so that the minimum page number is 1, not 0
    if (page > pagerState.totalPages) { page = pagerState.totalPages; }
    if (page < 1) { page = 1; }

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
