import React from "react";
import { Form, Col, Row } from "react-bootstrap";
import { ModalComponent } from "../Modal";
import { FormattedMessage, useIntl } from "react-intl";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { useConfigurationContext } from "src/hooks/Config";

interface IScrapeDialogProps {
  title: string;
  existingLabel?: string;
  scrapedLabel?: string;
  renderScrapeRows: () => JSX.Element;
  onClose: (apply?: boolean) => void;
}

export const ScrapeDialog: React.FC<IScrapeDialogProps> = (
  props: IScrapeDialogProps
) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

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
        <Form>
          <Row className="px-3 pt-3">
            <Col lg={{ span: 9, offset: 3 }}>
              <Row>
                <Form.Label column xs="6">
                  {props.existingLabel ?? (
                    <FormattedMessage id="dialogs.scrape_results_existing" />
                  )}
                </Form.Label>
                <Form.Label column xs="6">
                  {props.scrapedLabel ?? (
                    <FormattedMessage id="dialogs.scrape_results_scraped" />
                  )}
                </Form.Label>
              </Row>
            </Col>
          </Row>

          {props.renderScrapeRows()}
        </Form>
      </div>
    </ModalComponent>
  );
};
