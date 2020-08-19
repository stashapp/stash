import React, { useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
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
  const handleInput = (key: string, value: string) => {
    const newObj = {
      ...instance,
      [key]: value,
    };
    onSave(newObj);
  };

  return (
    <Form.Group className="row no-gutters">
      <InputGroup className="col-6">
        <Form.Control
          placeholder="Name"
          className="text-input"
          value={instance?.name}
          isValid={!isMulti || (instance?.name?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("name", e.currentTarget.value)
          }
        />
        <Form.Control
          placeholder="GraphQL endpoint"
          className="text-input"
          value={instance?.endpoint}
          isValid={(instance?.endpoint?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("endpoint", e.currentTarget.value)
          }
        />
        <Form.Control
          placeholder="API key"
          className="text-input"
          value={instance?.api_key}
          isValid={(instance?.api_key?.length ?? 0) > 0}
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleInput("api_key", e.currentTarget.value)
          }
        />
        <InputGroup.Append>
          <Button
            className=""
            variant="danger"
            title="Delete"
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
      <h4>Stash-box integration</h4>
      <div className="">
        {boxes.map((instance) => (
          <Instance
            instance={instance}
            onSave={handleSave}
            onDelete={handleDelete}
            key={instance.index}
            isMulti={boxes.length > 1}
          />
        ))}
      </div>
      <Button
        className="minimal"
        title="Add stash-box instance"
        onClick={handleAdd}
      >
        <Icon icon="plus" />
      </Button>
      <Form.Text className="text-muted">
        Stash-box facilitates automated tagging of scenes and performers based
        on fingerprints and filenames.
      </Form.Text>
    </Form.Group>
  );
};

export default StashBoxConfiguration;
