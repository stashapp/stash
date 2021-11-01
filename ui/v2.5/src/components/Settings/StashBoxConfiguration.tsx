import React, { useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared";

interface IInstanceProps {
  instance: IStashBoxInstance;
  onSave: (instance: IStashBoxInstance) => void;
  onDelete: (id: number) => void;
  isMulti: boolean;
}

const Instance: React.FC<IInstanceProps> = ({
  instance,
  onSave,
  onDelete,
  isMulti,
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
          placeholder={intl.formatMessage({ id: "config.stashbox.name" })}
          className="text-input col-3 stash-box-name"
          value={instance?.name}
          isValid={!isMulti || (instance?.name?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("name", e.currentTarget.value)
          }
        />
        <Form.Control
          placeholder={intl.formatMessage({
            id: "config.stashbox.graphql_endpoint",
          })}
          className="text-input col-3 stash-box-endpoint"
          value={instance?.endpoint}
          isValid={(instance?.endpoint?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("endpoint", e.currentTarget.value.trim())
          }
        />
        <Form.Control
          placeholder={intl.formatMessage({ id: "config.stashbox.api_key" })}
          className="text-input col-3 stash-box-apikey"
          value={instance?.api_key}
          isValid={(instance?.api_key?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("api_key", e.currentTarget.value.trim())
          }
        />
        <InputGroup.Append>
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

interface IStashBoxConfigurationProps {
  boxes: IStashBoxInstance[];
  saveBoxes: (boxes: IStashBoxInstance[]) => void;
}

export interface IStashBoxInstance {
  name?: string;
  endpoint?: string;
  api_key?: string;
  index: number;
}

export const StashBoxConfiguration: React.FC<IStashBoxConfigurationProps> = ({
  boxes,
  saveBoxes,
}) => {
  const intl = useIntl();
  const [index, setIndex] = useState(1000);

  const handleSave = (instance: IStashBoxInstance) =>
    saveBoxes(
      boxes.map((box) => (box.index === instance.index ? instance : box))
    );
  const handleDelete = (id: number) =>
    saveBoxes(boxes.filter((box) => box.index !== id));
  const handleAdd = () => {
    saveBoxes([...boxes, { index }]);
    setIndex(index + 1);
  };

  return (
    <Form.Group>
      <h6>{intl.formatMessage({ id: "config.stashbox.title" })}</h6>
      {boxes.length > 0 && (
        <div className="row no-gutters">
          <h6 className="col-3 ml-1">
            {intl.formatMessage({ id: "config.stashbox.name" })}
          </h6>
          <h6 className="col-3 ml-1">
            {intl.formatMessage({ id: "config.stashbox.endpoint" })}
          </h6>
          <h6 className="col-3 ml-1">
            {intl.formatMessage({ id: "config.general.auth.api_key" })}
          </h6>
        </div>
      )}
      {boxes.map((instance) => (
        <Instance
          instance={instance}
          onSave={handleSave}
          onDelete={handleDelete}
          key={instance.index}
          isMulti={boxes.length > 1}
        />
      ))}
      <Button
        className="minimal"
        title={intl.formatMessage({ id: "config.stashbox.add_instance" })}
        onClick={handleAdd}
      >
        <Icon icon="plus" />
      </Button>
      <Form.Text className="text-muted">
        {intl.formatMessage({ id: "config.stashbox.description" })}
      </Form.Text>
    </Form.Group>
  );
};
