import React, { useMemo } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { ModalComponent } from "../Modal";
import { FormattedMessage, useIntl } from "react-intl";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { useConfigurationContext } from "src/hooks/Config";

export interface IScrapeDialogContextState {
  existingLabel?: React.ReactNode;
  scrapedLabel?: React.ReactNode;
}

export const ScrapeDialogContext =
  React.createContext<IScrapeDialogContextState>({});

interface IScrapeDialogProps {
  title: string;
  existingLabel?: React.ReactNode;
  scrapedLabel?: React.ReactNode;
  renderScrapeRows: () => JSX.Element;
  onClose: (apply?: boolean) => void;
}

export const ScrapeDialog: React.FC<IScrapeDialogProps> = (
  props: IScrapeDialogProps
) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const existingLabel = useMemo(
    () =>
      props.existingLabel ?? (
        <FormattedMessage id="dialogs.scrape_results_existing" />
      ),
    [props.existingLabel]
  );
  const scrapedLabel = useMemo(
    () =>
      props.scrapedLabel ?? (
        <FormattedMessage id="dialogs.scrape_results_scraped" />
      ),
    [props.scrapedLabel]
  );

  const contextState = useMemo(
    () => ({
      existingLabel: existingLabel,
      scrapedLabel: scrapedLabel,
    }),
    [existingLabel, scrapedLabel]
  );

  return (
    <ModalComponent
      show
      icon={faPencilAlt}
      header={props.title}
      accept={{
        onClick: () => {
          props.onClose(true);
        },
        text: intl.formatMessage({ id: "actions.apply" }),
      }}
      cancel={{
        onClick: () => props.onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      modalProps={{
        size: "lg",
        dialogClassName: `scrape-dialog ${sfwContentMode ? "sfw-mode" : ""}`,
      }}
    >
      <div className="dialog-container">
        <ScrapeDialogContext.Provider value={contextState}>
          <Form>
            <Row className="px-3 pt-3">
              <Col lg={{ span: 9, offset: 3 }}>
                <Row>
                  <Form.Label column xs="6">
                    {existingLabel}
                  </Form.Label>
                  <Form.Label column xs="6">
                    {scrapedLabel}
                  </Form.Label>
                </Row>
              </Col>
            </Row>

            {props.renderScrapeRows()}
          </Form>
        </ScrapeDialogContext.Provider>
      </div>
    </ModalComponent>
  );
};
