import { faTimes } from "@fortawesome/free-solid-svg-icons";
import React, { useContext, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Collapse,
  Form,
  InputGroup,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import { Icon } from "src/components/Shared/Icon";
import { ParseMode, TagOperation } from "../constants";
import { TaggerStateContext } from "../context";
import { GenderEnum } from "src/core/generated-graphql";
import { genderList } from "src/utils/gender";

const Blacklist: React.FC<{
  list: string[];
  setList: (blacklist: string[]) => void;
}> = ({ list, setList }) => {
  const intl = useIntl();

  const [currentValue, setCurrentValue] = useState("");
  const [error, setError] = useState<string>();

  function addBlacklistItem() {
    if (!currentValue) return;

    // don't add duplicate items
    if (list.includes(currentValue)) {
      setError(
        intl.formatMessage({
          id: "component_tagger.config.errors.blacklist_duplicate",
        })
      );
      return;
    }

    // validate regex
    try {
      new RegExp(currentValue);
    } catch (e) {
      setError((e as SyntaxError).message);
      return;
    }

    setList([...list, currentValue]);

    setCurrentValue("");
  }

  function removeBlacklistItem(index: number) {
    const newBlacklist = [...list];
    newBlacklist.splice(index, 1);
    setList(newBlacklist);
  }

  return (
    <div>
      <h5>
        <FormattedMessage id="component_tagger.config.blacklist_label" />
      </h5>
      <Form.Group>
        <InputGroup hasValidation>
          <Form.Control
            className="text-input"
            value={currentValue}
            onChange={(e) => {
              setCurrentValue(e.currentTarget.value);
              setError(undefined);
            }}
            onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
              if (e.key === "Enter") {
                addBlacklistItem();
                e.preventDefault();
              }
            }}
            isInvalid={!!error}
          />
          <InputGroup.Append>
            <Button onClick={() => addBlacklistItem()}>
              <FormattedMessage id="actions.add" />
            </Button>
          </InputGroup.Append>
          <Form.Control.Feedback type="invalid">{error}</Form.Control.Feedback>
        </InputGroup>
      </Form.Group>
      <div>
        {intl.formatMessage(
          { id: "component_tagger.config.blacklist_desc" },
          { chars_require_escape: <code>[\^$.|?*+()</code> }
        )}
      </div>
      {list.map((item, index) => (
        <Badge
          className="tag-item d-inline-block"
          variant="secondary"
          key={item}
        >
          {item.toString()}
          <Button
            className="minimal ml-2"
            onClick={() => removeBlacklistItem(index)}
          >
            <Icon icon={faTimes} />
          </Button>
        </Badge>
      ))}
    </div>
  );
};

interface IConfigProps {
  show: boolean;
}

const Config: React.FC<IConfigProps> = ({ show }) => {
  const { config, setConfig } = useContext(TaggerStateContext);
  const intl = useIntl();

  function renderGenderCheckbox(gender: GenderEnum) {
    const performerGenders = config.performerGenders || genderList.slice();
    return (
      <Form.Check
        id={`gender-${gender}`}
        key={gender}
        label={<FormattedMessage id={`gender_types.${gender}`} />}
        checked={performerGenders.includes(gender)}
        onChange={(e) => {
          const isChecked = e.currentTarget.checked;
          setConfig({
            ...config,
            performerGenders: isChecked
              ? [...performerGenders, gender]
              : performerGenders.filter((g) => g !== gender),
          });
        }}
      />
    );
  }

  return (
    <Collapse in={show}>
      <Card>
        <div className="row">
          <h4 className="col-12">
            <FormattedMessage id="configuration" />
          </h4>
          <hr className="w-100" />
          <Form className="col-md-6">
            <Form.Group
              controlId="performer-genders"
              className="align-items-center"
            >
              <Form.Label>
                <FormattedMessage id="component_tagger.config.performer_genders.heading" />
              </Form.Label>
              {genderList.map(renderGenderCheckbox)}
              <Form.Text>
                <FormattedMessage id="component_tagger.config.performer_genders.description" />
              </Form.Text>
            </Form.Group>
            <Form.Group controlId="set-cover" className="align-items-center">
              <Form.Check
                label={
                  <FormattedMessage id="component_tagger.config.set_cover_label" />
                }
                checked={config.setCoverImage}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    setCoverImage: e.currentTarget.checked,
                  })
                }
              />
              <Form.Text>
                <FormattedMessage id="component_tagger.config.set_cover_desc" />
              </Form.Text>
            </Form.Group>
            <Form.Group className="align-items-center">
              <div className="d-flex align-items-center">
                <Form.Check
                  id="tag-mode"
                  label={
                    <FormattedMessage id="component_tagger.config.set_tag_label" />
                  }
                  className="mr-4"
                  checked={config.setTags}
                  onChange={(e) =>
                    setConfig({ ...config, setTags: e.currentTarget.checked })
                  }
                />
                <Form.Control
                  id="tag-operation"
                  className="col-md-2 col-3 input-control"
                  as="select"
                  value={config.tagOperation}
                  onChange={(e) =>
                    setConfig({
                      ...config,
                      tagOperation: e.currentTarget.value as TagOperation,
                    })
                  }
                  disabled={!config.setTags}
                >
                  <option value="merge">
                    {intl.formatMessage({ id: "actions.merge" })}
                  </option>
                  <option value="overwrite">
                    {intl.formatMessage({ id: "actions.overwrite" })}
                  </option>
                </Form.Control>
              </div>
              <Form.Text>
                <FormattedMessage id="component_tagger.config.set_tag_desc" />
              </Form.Text>
            </Form.Group>

            <Form.Group controlId="mode-select">
              <div className="row no-gutters">
                <Form.Label className="mr-4 mt-1">
                  <FormattedMessage id="component_tagger.config.query_mode_label" />
                  :
                </Form.Label>
                <Form.Control
                  as="select"
                  className="col-md-2 col-3 input-control"
                  value={config.mode}
                  onChange={(e) =>
                    setConfig({
                      ...config,
                      mode: e.currentTarget.value as ParseMode,
                    })
                  }
                >
                  <option value="auto">
                    {intl.formatMessage({
                      id: "component_tagger.config.query_mode_auto",
                    })}
                  </option>
                  <option value="filename">
                    {intl.formatMessage({
                      id: "component_tagger.config.query_mode_filename",
                    })}
                  </option>
                  <option value="dir">
                    {intl.formatMessage({
                      id: "component_tagger.config.query_mode_dir",
                    })}
                  </option>
                  <option value="path">
                    {intl.formatMessage({
                      id: "component_tagger.config.query_mode_path",
                    })}
                  </option>
                  <option value="metadata">
                    {intl.formatMessage({
                      id: "component_tagger.config.query_mode_metadata",
                    })}
                  </option>
                </Form.Control>
              </div>
              <Form.Text>
                {intl.formatMessage({
                  id: `component_tagger.config.query_mode_${config.mode}_desc`,
                  defaultMessage: "Unknown query mode",
                })}
              </Form.Text>
            </Form.Group>
            <Form.Group controlId="toggle-organized">
              <Form.Check
                label={
                  <FormattedMessage id="component_tagger.config.mark_organized_label" />
                }
                checked={config.markSceneAsOrganizedOnSave}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    markSceneAsOrganizedOnSave: e.currentTarget.checked,
                  })
                }
              />
              <Form.Text>
                <FormattedMessage id="component_tagger.config.mark_organized_desc" />
              </Form.Text>
            </Form.Group>
          </Form>
          <div className="col-md-6">
            <Blacklist
              list={config.blacklist}
              setList={(blacklist) => setConfig({ ...config, blacklist })}
            />
          </div>
        </div>
      </Card>
    </Collapse>
  );
};

export default Config;
