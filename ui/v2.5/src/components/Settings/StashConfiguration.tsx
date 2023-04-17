import { faEllipsisV } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Form, Row, Col, Dropdown } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import * as GQL from "src/core/generated-graphql";
import { FolderSelectDialog } from "../Shared/FolderSelect/FolderSelectDialog";
import { BooleanSetting } from "./Inputs";
import { SettingSection } from "./SettingSection";

interface IStashProps {
  index: number;
  stash: GQL.StashConfig;
  onSave: (instance: GQL.StashConfig) => void;
  onEdit: () => void;
  onDelete: () => void;
}

const Stash: React.FC<IStashProps> = ({
  index,
  stash,
  onSave,
  onEdit,
  onDelete,
}) => {
  // eslint-disable-next-line
  const handleInput = (key: string, value: any) => {
    const newObj = {
      ...stash,
      [key]: value,
    };
    onSave(newObj);
  };

  const classAdd = index % 2 === 1 ? "bg-dark" : "";

  return (
    <Row className={`stash-row align-items-center ${classAdd}`}>
      <Form.Label column md={7}>
        {stash.path}
      </Form.Label>
      <Col md={2} xs={4} className="col form-label">
        {/* NOTE - language is opposite to meaning:
        internally exclude flags, displayed as include */}
        <div>
          <h6 className="d-md-none">
            <FormattedMessage id="videos" />
          </h6>
          <BooleanSetting
            id={`stash-exclude-video-${index}`}
            checked={!stash.excludeVideo}
            onChange={(v) => handleInput("excludeVideo", !v)}
          />
        </div>
      </Col>

      <Col md={2} xs={4} className="col-form-label">
        <div>
          <h6 className="d-md-none">
            <FormattedMessage id="images" />
          </h6>
          <BooleanSetting
            id={`stash-exclude-image-${index}`}
            checked={!stash.excludeImage}
            onChange={(v) => handleInput("excludeImage", !v)}
          />
        </div>
      </Col>
      <Col className="justify-content-end" xs={4} md={1}>
        <Dropdown className="text-right">
          <Dropdown.Toggle
            variant="minimal"
            id={`stash-menu-${index}`}
            className="minimal"
          >
            <Icon icon={faEllipsisV} />
          </Dropdown.Toggle>
          <Dropdown.Menu className="bg-secondary text-white">
            <Dropdown.Item onClick={() => onEdit()}>
              <FormattedMessage id="actions.edit" />
            </Dropdown.Item>
            <Dropdown.Item onClick={() => onDelete()}>
              <FormattedMessage id="actions.delete" />
            </Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown>
      </Col>
    </Row>
  );
};

interface IStashConfigurationProps {
  stashes: GQL.StashConfig[];
  setStashes: (v: GQL.StashConfig[]) => void;
}

const StashConfiguration: React.FC<IStashConfigurationProps> = ({
  stashes,
  setStashes,
}) => {
  const [isCreating, setIsCreating] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | undefined>();

  function onEdit(index: number) {
    setEditingIndex(index);
  }

  function onDelete(index: number) {
    setStashes(stashes.filter((v, i) => i !== index));
  }

  function onNew() {
    setIsCreating(true);
  }

  const handleSave = (index: number, stash: GQL.StashConfig) =>
    setStashes(stashes.map((s, i) => (i === index ? stash : s)));

  return (
    <>
      {isCreating ? (
        <FolderSelectDialog
          onClose={(v) => {
            if (v)
              setStashes([
                ...stashes,
                {
                  path: v,
                  excludeVideo: false,
                  excludeImage: false,
                },
              ]);
            setIsCreating(false);
          }}
        />
      ) : undefined}

      {editingIndex !== undefined ? (
        <FolderSelectDialog
          defaultValue={stashes[editingIndex].path}
          onClose={(v) => {
            if (v)
              setStashes(
                stashes.map((vv, index) => {
                  if (index === editingIndex) {
                    return {
                      ...vv,
                      path: v,
                    };
                  }
                  return vv;
                })
              );
            setEditingIndex(undefined);
          }}
        />
      ) : undefined}

      <div className="content" id="stash-table">
        {stashes.length > 0 && (
          <Row className="d-none d-md-flex">
            <h6 className="col-md-7">
              <FormattedMessage id="path" />
            </h6>
            <h6 className="col-md-2 col-4">
              <FormattedMessage id="videos" />
            </h6>
            <h6 className="col-md-2 col-4">
              <FormattedMessage id="images" />
            </h6>
          </Row>
        )}
        {stashes.map((stash, index) => (
          <Stash
            index={index}
            stash={stash}
            onSave={(s) => handleSave(index, s)}
            onEdit={() => onEdit(index)}
            onDelete={() => onDelete(index)}
            key={stash.path}
          />
        ))}
        <Button className="mt-2" variant="secondary" onClick={() => onNew()}>
          <FormattedMessage id="actions.add_directory" />
        </Button>
      </div>
    </>
  );
};

interface IStashSetting {
  value: GQL.StashConfigInput[];
  onChange: (v: GQL.StashConfigInput[]) => void;
}

export const StashSetting: React.FC<IStashSetting> = ({ value, onChange }) => {
  return (
    <SettingSection
      id="stashes"
      headingID="library"
      subHeadingID="config.general.directory_locations_to_your_content"
    >
      <StashConfiguration stashes={value} setStashes={(v) => onChange(v)} />
    </SettingSection>
  );
};

export default StashConfiguration;
