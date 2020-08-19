import React, { useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared";

interface IInstanceProps {
  id?: number;
  instance?: GQL.StashBox;
  isCreate?: boolean;
  onSave: (instance: GQL.StashBoxInput, id?: number) => void;
  onDelete?: (id: number) => void;
  onCancel: () => void;
}

const Instance: React.FC<IInstanceProps> = ({
  id,
  instance,
  onSave,
  onCancel,
  onDelete,
  isCreate = false,
}) => {
  const [isEditing, setIsEditing] = useState(isCreate);
  const [endpoint, setEndpoint] = useState(instance?.endpoint);
  const [apiKey, setApiKey] = useState(instance?.api_key);

  const handleCancel = () => {
    if (isCreate) {
      onCancel();
      setEndpoint("");
      setApiKey("");
    } else {
      setIsEditing(false);
    }
  };

  const handleSave = () => {
    if (!endpoint || !apiKey) return;
    setIsEditing(false);
    onSave(
      {
        api_key: apiKey,
        endpoint,
      },
      id
    );
    if (id === undefined) {
      setEndpoint("");
      setApiKey("");
    }
  };

  return (
    <Form.Group className="row">
      <InputGroup className="col-6">
        <InputGroup.Prepend>
          {!isEditing && (
            <>
              <Button
                className=""
                variant="primary"
                title="Edit"
                onClick={() => setIsEditing(true)}
              >
                <Icon icon="edit" />
              </Button>
              {id !== undefined && (
                <Button
                  className=""
                  variant="danger"
                  title="Delete"
                  onClick={() => onDelete?.(id)}
                >
                  <Icon icon="minus" />
                </Button>
              )}
            </>
          )}
          {isEditing && (
            <>
              <Button
                className=""
                variant="primary"
                title="Save"
                onClick={handleSave}
              >
                <Icon icon="save" />
              </Button>
              <Button
                className=""
                variant="danger"
                title="Cancel"
                onClick={handleCancel}
              >
                <Icon icon="times" />
              </Button>
            </>
          )}
        </InputGroup.Prepend>
        <Form.Control
          placeholder="GraphQL endpoint"
          className="text-input"
          value={endpoint}
          disabled={!isEditing}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setEndpoint(e.currentTarget.value)
          }
        />
        <Form.Control
          placeholder="API key"
          className="text-input"
          value={apiKey}
          disabled={!isEditing}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setApiKey(e.currentTarget.value)
          }
        />
      </InputGroup>
    </Form.Group>
  );
};

interface IStashBoxConfigurationProps {
  boxes: GQL.StashBox[];
  saveBoxes: (boxes: GQL.StashBoxInput[]) => void;
}

export const StashBoxConfiguration: React.FC<IStashBoxConfigurationProps> = ({
  boxes,
  saveBoxes,
}) => {
  const [showCreate, setShowCreate] = useState(false);

  const handleCancel = () => setShowCreate(false);
  const handleSave = (instance: GQL.StashBoxInput, id?: number) => {
    if (!instance.api_key || !instance.endpoint) return;

    const newBoxes =
      id !== undefined
        ? boxes.map((box, index) => (index === id ? instance : box))
        : [...boxes, instance];

    if (id === undefined) setShowCreate(false);
    saveBoxes(newBoxes);
  };
  const handleDelete = (id: number) => {
    const newBoxes = boxes.filter((_, index) => index !== id);
    saveBoxes(newBoxes);
  };

  return (
    <Form.Group>
      <h4>Stash-box integration</h4>
      <div className="">
        {boxes.map((instance, index) => (
          <Instance
            instance={instance}
            onSave={handleSave}
            onCancel={handleCancel}
            onDelete={handleDelete}
            key={instance.endpoint}
            id={index}
          />
        ))}
        {showCreate && (
          <Instance onSave={handleSave} onCancel={handleCancel} isCreate />
        )}
      </div>
      <Button
        className="minimal"
        title="Add stash-box instance"
        onClick={() => setShowCreate(true)}
        disabled={showCreate}
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
