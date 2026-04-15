import React, { useEffect, useMemo, useRef, useState } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  Form,
  InputGroup,
  Overlay,
  Popover,
} from "react-bootstrap";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import useFocus from "src/utils/focus";
import { Icon } from "../Shared/Icon";
import { faCheck, faChevronDown } from "@fortawesome/free-solid-svg-icons";
import { useStopWheelScroll } from "src/utils/form";
import { Placement } from "react-bootstrap/esm/Overlay";
import { PatchComponent } from "src/patch";

const PageCount: React.FC<{
  totalPages: number;
  currentPage: number;
  onChangePage: (page: number) => void;
  pagePopupPlacement?: Placement;
}> = ({
  totalPages,
  currentPage,
  onChangePage,
  pagePopupPlacement = "bottom",
}) => {
  const intl = useIntl();
  const currentPageCtrl = useRef(null);
  const [pageInput, pageFocus] = useFocus();
  const [showSelectPage, setShowSelectPage] = useState(false);

  useEffect(() => {
    if (showSelectPage) {
      // delaying the focus to the next execution loop so that rendering takes place first and stops the page from resetting.
      setTimeout(() => {
        pageFocus();
      }, 0);
    }
  }, [showSelectPage, pageFocus]);

  useStopWheelScroll(pageInput);

  const pageOptions = useMemo(() => {
    const maxPagesToShow = 1000;
    const min = Math.max(1, currentPage - maxPagesToShow / 2);
    const max = Math.min(min + maxPagesToShow, totalPages);
    const pages = [];
    for (let i = min; i <= max; i++) {
      pages.push(i);
    }
    return pages;
  }, [totalPages, currentPage]);

  function onCustomChangePage() {
    const newPage = Number.parseInt(pageInput.current?.value ?? "0");
    if (newPage) {
      onChangePage(newPage);
    }
    setShowSelectPage(false);
  }

  return (
    <div className="page-count-container">
      <ButtonGroup>
        <Button
          variant="secondary"
          className="page-count"
          ref={currentPageCtrl}
          onClick={() => {
            setShowSelectPage(true);
            pageFocus();
          }}
        >
          <FormattedMessage
            id="pagination.current_total"
            values={{
              current: intl.formatNumber(currentPage),
              total: intl.formatNumber(totalPages),
            }}
          />
        </Button>
        <Dropdown>
          <Dropdown.Toggle variant="secondary" className="page-count-dropdown">
            <Icon size="xs" icon={faChevronDown} />
          </Dropdown.Toggle>
          <Dropdown.Menu>
            {pageOptions.map((s) => (
              <Dropdown.Item
                key={s}
                active={s === currentPage}
                onClick={() => onChangePage(s)}
              >
                {s}
              </Dropdown.Item>
            ))}
          </Dropdown.Menu>
        </Dropdown>
      </ButtonGroup>
      <Overlay
        target={currentPageCtrl.current}
        show={showSelectPage}
        placement={pagePopupPlacement}
        rootClose
        onHide={() => setShowSelectPage(false)}
      >
        <Popover id="select_page_popover">
          <Form inline>
            <InputGroup>
              {/* can't use NumberField because of the ref */}
              <Form.Control
                type="number"
                min={1}
                max={totalPages}
                className="text-input"
                ref={pageInput}
                defaultValue={currentPage}
                onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                  if (e.key === "Enter") {
                    onCustomChangePage();
                    e.preventDefault();
                  }
                }}
                onFocus={(e: React.FocusEvent<HTMLInputElement>) =>
                  e.target.select()
                }
              />
              <InputGroup.Append>
                <Button variant="primary" onClick={() => onCustomChangePage()}>
                  <Icon icon={faCheck} />
                </Button>
              </InputGroup.Append>
            </InputGroup>
          </Form>
        </Popover>
      </Overlay>
    </div>
  );
};

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  metadataByline?: React.ReactNode;
  onChangePage: (page: number) => void;
  pagePopupPlacement?: Placement;
}

interface IPaginationIndexProps {
  loading?: boolean;
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  metadataByline?: React.ReactNode;
}

const minPagesForCompact = 4;

export const Pagination: React.FC<IPaginationProps> = PatchComponent(
  "Pagination",
  ({
    itemsPerPage,
    currentPage,
    totalItems,
    onChangePage,
    pagePopupPlacement,
  }) => {
    const intl = useIntl();
    const totalPages = useMemo(
      () => Math.ceil(totalItems / itemsPerPage),
      [totalItems, itemsPerPage]
    );

    const pageButtons = useMemo(() => {
      if (totalPages >= minPagesForCompact)
        return (
          <PageCount
            totalPages={totalPages}
            currentPage={currentPage}
            onChangePage={onChangePage}
            pagePopupPlacement={pagePopupPlacement}
          />
        );

      const pages = [...Array(totalPages).keys()].map((i) => i + 1);

      return pages.map((page: number) => (
        <Button
          variant="secondary"
          key={page}
          active={currentPage === page}
          onClick={() => onChangePage(page)}
        >
          <FormattedNumber value={page} />
        </Button>
      ));
    }, [totalPages, currentPage, onChangePage, pagePopupPlacement]);

    if (totalPages <= 1) return <div />;

    return (
      <ButtonGroup className="pagination">
        <Button
          variant="secondary"
          disabled={currentPage === 1}
          onClick={() => onChangePage(1)}
          title={intl.formatMessage({ id: "pagination.first" })}
        >
          <span>«</span>
        </Button>
        <Button
          variant="secondary"
          disabled={currentPage === 1}
          onClick={() => onChangePage(currentPage - 1)}
          title={intl.formatMessage({ id: "pagination.previous" })}
        >
          &lt;
        </Button>
        {pageButtons}
        <Button
          variant="secondary"
          disabled={currentPage === totalPages}
          onClick={() => onChangePage(currentPage + 1)}
          title={intl.formatMessage({ id: "pagination.next" })}
        >
          &gt;
        </Button>
        <Button
          variant="secondary"
          disabled={currentPage === totalPages}
          onClick={() => onChangePage(totalPages)}
          title={intl.formatMessage({ id: "pagination.last" })}
        >
          <span>»</span>
        </Button>
      </ButtonGroup>
    );
  }
);

export const PaginationIndex: React.FC<IPaginationIndexProps> = PatchComponent(
  "PaginationIndex",
  ({ loading, itemsPerPage, currentPage, totalItems, metadataByline }) => {
    const intl = useIntl();

    if (loading) return null;

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
  }
);
