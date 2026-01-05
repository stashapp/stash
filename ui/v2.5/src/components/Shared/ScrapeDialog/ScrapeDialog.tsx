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
  className?: string;
  title: string;
  existingLabel?: React.ReactNode;
  scrapedLabel?: React.ReactNode;
  onClose: (apply?: boolean) => void;
}

export const ScrapeDialog: React.FC<
  React.PropsWithChildren<IScrapeDialogProps>
> = (props: React.PropsWithChildren<IScrapeDialogProps>) => {
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
        dialogClassName: `${props.className ?? ""} scrape-dialog ${
          sfwContentMode ? "sfw-mode" : ""
        }`,
      }}
    >
      <div className="dialog-container">
        <ScrapeDialogContext.Provider value={contextState}>
          <Form>
            <Row className="px-3 pt-3">
              <Col lg={{ span: 9, offset: 3 }}>
                <Row>
                  <Form.Label
                    column
                    lg="6"
                    className="d-lg-block d-none column-label"
                  >
                    {existingLabel}
                  </Form.Label>
                  <Form.Label
                    column
                    lg="6"
                    className="d-lg-block d-none column-label"
                  >
                    {scrapedLabel}
                  </Form.Label>
                </Row>
              </Col>
            </Row>

            {props.children}
          </Form>
        </ScrapeDialogContext.Provider>
      </div>
    </ModalComponent>
  );
};
