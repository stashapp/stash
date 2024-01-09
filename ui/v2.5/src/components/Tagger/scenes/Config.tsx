import { faTimes } from "@fortawesome/free-solid-svg-icons";
import React, { useRef, useContext } from "react";
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

interface IConfigProps {
  show: boolean;
}

const Config: React.FC<IConfigProps> = ({ show }) => {
  const { config, setConfig } = useContext(TaggerStateContext);
  const intl = useIntl();
  const blacklistRef = useRef<HTMLInputElement | null>(null);

  function addBlacklistItem() {
    if (!blacklistRef.current) return;

    const input = blacklistRef.current.value;
    if (!input) return;

    // don't add duplicate items
    if (!config.blacklist.includes(input)) {
      setConfig({
        ...config,
        blacklist: [...config.blacklist, input],
      });
    }

    blacklistRef.current.value = "";
  }

  function removeBlacklistItem(index: number) {
    const newBlacklist = [...config.blacklist];
    newBlacklist.splice(index, 1);
    setConfig({
      ...config,
      blacklist: newBlacklist,
    });
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
            <Form.Group controlId="tag-males" className="align-items-center">
              <Form.Check
                label={
                  <FormattedMessage id="component_tagger.config.show_male_label" />
                }
                checked={config.showMales}
                onChange={(e) =>
                  setConfig({ ...config, showMales: e.currentTarget.checked })
                }
              />
              <Form.Text>
                <FormattedMessage id="component_tagger.config.show_male_desc" />
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
            <h5>
              <FormattedMessage id="component_tagger.config.blacklist_label" />
            </h5>
            <InputGroup>
              <Form.Control
                className="text-input"
                ref={blacklistRef}
                onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
                  if (e.key === "Enter") {
                    addBlacklistItem();
                    e.preventDefault();
                  }
                }}
              />
              <InputGroup.Append>
                <Button onClick={() => addBlacklistItem()}>
                  <FormattedMessage id="actions.add" />
                </Button>
              </InputGroup.Append>
            </InputGroup>
            <div>
              {intl.formatMessage(
                { id: "component_tagger.config.blacklist_desc" },
                { chars_require_escape: <code>[\^$.|?*+()</code> }
              )}
            </div>
            {config.blacklist.map((item, index) => (
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
        </div>
      </Card>
    </Collapse>
  );
};

export default Config;
