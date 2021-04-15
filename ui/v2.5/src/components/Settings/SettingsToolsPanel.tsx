import React from "react";
import { Button, Form } from "react-bootstrap";
import { Link } from "react-router-dom";

export const SettingsToolsPanel: React.FC = () => {
  return (
    <>
      <h4>Scene Tools</h4>

      <Form.Group>
        <Link to="/sceneFilenameParser">
          <Button variant="secondary">Scene Filename Parser</Button>
        </Link>
      </Form.Group>

      <Form.Group>
        <Link to="/settings?tab=duplicates">
          <Button variant="secondary">Scene Duplicate Checker</Button>
        </Link>
      </Form.Group>
    </>
  );
};
