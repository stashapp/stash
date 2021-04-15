import React from "react";
import { Button, Form } from "react-bootstrap";
import { Link } from "react-router-dom";

export const SettingsToolsPanel: React.FC = () => {
  return (
    <>
      <h4>Scene Tools</h4>

      <Form.Group>
        <Link to="/sceneFilenameParser">
          Scene Filename Parser
        </Link>
      </Form.Group>

      <Form.Group>
        <Link to="/sceneDuplicateChecker">
          Scene Duplicate Checker
        </Link>
      </Form.Group>
    </>
  );
};
