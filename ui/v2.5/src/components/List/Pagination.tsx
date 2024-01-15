import React, { useEffect, useMemo, useRef, useState } from "react";
import {
  Button,
  ButtonGroup,
  Form,
  InputGroup,
  Overlay,
  Popover,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import useFocus from "src/utils/focus";
import { Icon } from "../Shared/Icon";
import { faCheck, faChevronDown } from "@fortawesome/free-solid-svg-icons";

interface IPaginationProps {
  itemsPerPage: number;
  currentPage: number;
  totalItems: number;
  metadataByline?: React.ReactNode;
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
  onChangePage,
}) => {
  const intl = useIntl();

  const currentPageCtrl = useRef(null);
  const [pageInput, pageFocus] = useFocus();

  const [showSelectPage, setShowSelectPage] = useState(false);

  useEffect(() => {
    if (showSelectPage) {
      pageFocus();
    }
  }, [showSelectPage, pageFocus]);

  const totalPages = useMemo(
    () => Math.ceil(totalItems / itemsPerPage),
    [totalItems, itemsPerPage]
  );

  function onCustomChangePage() {
    const newPage = Number.parseInt(pageInput.current?.value ?? "0");
    if (newPage) {
      onChangePage(newPage);
    }
    setShowSelectPage(false);
  }

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
        className="d-none d-sm-block"
        variant="secondary"
        disabled={currentPage === 1}
        onClick={() => onChangePage(currentPage - 1)}
        title={intl.formatMessage({ id: "pagination.previous" })}
      >
        &lt;
      </Button>
      <div>
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
          />{" "}
          <Icon size="xs" icon={faChevronDown} />
        </Button>
        <Overlay
          target={currentPageCtrl.current}
          show={showSelectPage}
          placement="bottom"
          rootClose
          onHide={() => setShowSelectPage(false)}
        >
          <Popover id="select_page_popover">
            <Form inline>
              <InputGroup>
                <Form.Control
                  type="number"
                  min={1}
                  max={totalPages}
                  className="text-input"
                  ref={pageInput}
                  defaultValue={currentPage}
                  onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
                    if (e.key === "Enter") {
                      onCustomChangePage();
                      e.preventDefault();
                    }
                  }}
                />
                <InputGroup.Append>
                  <Button
                    variant="primary"
                    onClick={() => onCustomChangePage()}
                  >
                    <Icon icon={faCheck} />
                  </Button>
                </InputGroup.Append>
              </InputGroup>
            </Form>
          </Popover>
        </Overlay>
      </div>
      <Button
        className="d-none d-sm-block"
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
