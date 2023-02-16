import React from "react";
import { useIntl } from "react-intl";
import { Button, InputGroup, Form } from "react-bootstrap";
import { Icon } from "./Icon";
import { FormikHandlers } from "formik";
import { faFileDownload } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  value: string;
  name: string;
  onChange: FormikHandlers["handleChange"];
  onBlur: FormikHandlers["handleBlur"];
  onScrapeClick(): void;
  urlScrapable(url: string): boolean;
  isInvalid?: boolean;
}

export const URLField: React.FC<IProps> = (props: IProps) => {
  const intl = useIntl();

  return (
    <InputGroup className="mr-2 flex-grow-1">
      <Form.Control
        className="text-input"
        placeholder={intl.formatMessage({ id: "url" })}
        value={props.value}
        name={props.name}
        onChange={props.onChange}
        onBlur={props.onBlur}
        isInvalid={props.isInvalid}
      />
      <InputGroup.Append>
        <Button
          className="scrape-url-button text-input"
          variant="secondary"
          onClick={props.onScrapeClick}
          disabled={!props.value || !props.urlScrapable(props.value)}
          title={intl.formatMessage({ id: "actions.scrape" })}
        >
          <Icon icon={faFileDownload} />
        </Button>
      </InputGroup.Append>
    </InputGroup>
  );
};
