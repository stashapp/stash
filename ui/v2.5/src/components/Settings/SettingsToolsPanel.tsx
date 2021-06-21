import React from "react";
import { Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";

export const SettingsToolsPanel: React.FC = () => {
  return (
    <>
      <h4>
        <FormattedMessage id="config.tools.scene_tools" />
      </h4>

      <Form.Group>
        <Link to="/sceneFilenameParser">
          <FormattedMessage id="config.tools.scene_filename_parser.title" />
        </Link>
      </Form.Group>

      <Form.Group>
        <Link to="/sceneDuplicateChecker">
          <FormattedMessage id="config.tools.scene_duplicate_checker" />
        </Link>
      </Form.Group>
    </>
  );
};
