import React, { useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared";

interface IInstanceProps {
  instance: IPerformerURLInputInstance;
  onSave: (instance: IPerformerURLInputInstance) => void;
  onDelete: (id: number) => void;
  onScrape: (url: string) => void;
  urlScrapable(url: string): boolean;
}

const Instance: React.FC<IInstanceProps> = ({
  instance,
  onSave,
  onDelete,
  onScrape,
  urlScrapable,
}) => {
  const intl = useIntl();
  const handleInput = (key: string, value: string) => {
    const newObj = {
      ...instance,
      [key]: value,
    };
    onSave(newObj);
  };

  return (
    <Form.Group className="row no-gutters">
      <InputGroup className="col">
        <Form.Control
          placeholder="URL"
          className="text-input"
          value={instance?.url}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("url", e.currentTarget.value.trim())
          }
        />
        <InputGroup.Append>
          <Button
            className="scrape-url-button text-input"
            variant="secondary"
            onClick={() => onScrape(instance?.url ?? "")}
            disabled={!instance?.url || !urlScrapable(instance?.url)}
            title={intl.formatMessage({ id: "actions.scrape" })}
          >
            <Icon icon="file-download" />
          </Button>
          <Button
            className=""
            variant="danger"
            title={intl.formatMessage({ id: "actions.delete" })}
            onClick={() => onDelete(instance.index)}
          >
            <Icon icon="minus" />
          </Button>
        </InputGroup.Append>
      </InputGroup>
    </Form.Group>
  );
};

interface IPerformerURLInputProps {
  urls: IPerformerURLInputInstance[];
  saveURLs: (boxes: IPerformerURLInputInstance[]) => void;
  onScrapeClick(url: string): void;
  urlScrapable(url: string): boolean;
}

export interface IPerformerURLInputInstance {
  url?: string;
  index: number;
}

export const PerformerURLInput: React.FC<IPerformerURLInputProps> = ({
  urls,
  saveURLs,
  onScrapeClick,
  urlScrapable,
}) => {
  const intl = useIntl();
  const [index, setIndex] = useState(1000);

  const handleSave = (instance: IPerformerURLInputInstance) =>
    saveURLs(
      urls.map((url) => (url.index === instance.index ? instance : url))
    );
  const handleDelete = (id: number) =>
    saveURLs(urls.filter((url) => url.index !== id));
  const handleAdd = () => {
    saveURLs([...urls, { index }]);
    setIndex(index + 1);
  };

  return (
    <Form.Group>
      {urls.map((instance) => (
        <Instance
          instance={instance}
          onSave={handleSave}
          onDelete={handleDelete}
          onScrape={onScrapeClick}
          urlScrapable={urlScrapable}
          key={instance.index}
        />
      ))}
      <Button
        className="minimal"
        title={intl.formatMessage({ id: "actions.add" })}
        onClick={handleAdd}
      >
        <Icon icon="plus" />
      </Button>
    </Form.Group>
  );
};
