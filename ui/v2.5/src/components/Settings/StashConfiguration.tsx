import React, { useState } from "react";
import { Button, Form, Row, Col } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { FolderSelectDialog } from "../Shared/FolderSelect/FolderSelectDialog";

interface IStashProps {
  index: number;
  stash: GQL.StashConfig;
  onSave: (instance: GQL.StashConfig) => void;
  onDelete: () => void;
}

const Stash: React.FC<IStashProps> = ({ index, stash, onSave, onDelete }) => {
  // eslint-disable-next-line
  const handleInput = (key: string, value: any) => {
    const newObj = {
      ...stash,
      [key]: value,
    };
    onSave(newObj);
  };

  const intl = useIntl();
  const classAdd = index % 2 === 1 ? "bg-dark" : "";

  return (
    <Row className={`align-items-center ${classAdd}`}>
      <Form.Label column xs={4}>
        {stash.path}
      </Form.Label>
      <Col xs={3}>
        <Form.Check
          id="stash-exclude-video"
          checked={stash.excludeVideo}
          onChange={() => handleInput("excludeVideo", !stash.excludeVideo)}
        />
      </Col>

      <Col xs={3}>
        <Form.Check
          id="stash-exclude-image"
          checked={stash.excludeImage}
          onChange={() => handleInput("excludeImage", !stash.excludeImage)}
        />
      </Col>
      <Col xs={2}>
        <Button
          size="sm"
          variant="danger"
          title={intl.formatMessage({ id: 'actions.delete' })}
          onClick={() => onDelete()}
        >
          <Icon icon="minus" />
        </Button>
      </Col>
    </Row>
  );
};

interface IStashConfigurationProps {
  stashes: GQL.StashConfig[];
  setStashes: (v: GQL.StashConfig[]) => void;
}

export const StashConfiguration: React.FC<IStashConfigurationProps> = ({
  stashes,
  setStashes,
}) => {
  const [isDisplayingDialog, setIsDisplayingDialog] = useState<boolean>(false);

  const handleSave = (index: number, stash: GQL.StashConfig) =>
    setStashes(stashes.map((s, i) => (i === index ? stash : s)));
  const handleDelete = (index: number) =>
    setStashes(stashes.filter((s, i) => i !== index));
  const handleAdd = (folder?: string) => {
    setIsDisplayingDialog(false);

    if (!folder) {
      return;
    }

    setStashes([
      ...stashes,
      {
        path: folder,
        excludeImage: false,
        excludeVideo: false,
      },
    ]);
  };

  function maybeRenderDialog() {
    if (!isDisplayingDialog) {
      return;
    }

    return <FolderSelectDialog onClose={handleAdd} />;
  }

  return (
    <>
      {maybeRenderDialog()}
      <Form.Group>
        {stashes.length > 0 && (
          <Row>
            <h6 className="col-4">Path</h6>
            <h6 className="col-3">Exclude Video</h6>
            <h6 className="col-3">Exclude Image</h6>
          </Row>
        )}
        {stashes.map((stash, index) => (
          <Stash
            index={index}
            stash={stash}
            onSave={(s) => handleSave(index, s)}
            onDelete={() => handleDelete(index)}
            key={stash.path}
          />
        ))}
        <Button
          className="mt-2"
          variant="secondary"
          onClick={() => setIsDisplayingDialog(true)}
        >
          Add Directory
        </Button>
      </Form.Group>
    </>
  );
};

export default StashConfiguration;
