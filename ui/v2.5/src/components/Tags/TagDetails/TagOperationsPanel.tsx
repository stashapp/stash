import { Button, Form } from "react-bootstrap";
import React from "react";
import * as GQL from "src/core/generated-graphql";
import { mutateMetadataAutoTag } from "src/core/StashService";
import { useToast } from "src/hooks";

interface ITagOperationsProps {
  tag: Partial<GQL.TagDataFragment>;
}

export const TagOperationsPanel: React.FC<ITagOperationsProps> = ({ tag }) => {
  const Toast = useToast();

  async function onAutoTag() {
    if (!tag?.id) return;
    try {
      await mutateMetadataAutoTag({ tags: [tag.id] });
      Toast.success({ content: "Started auto tagging" });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <>
      <h5>Auto Tagging</h5>
      <Form.Group>
        <Button onClick={onAutoTag}>Auto Tag</Button>
        <Form.Text className="text-muted">
          Auto-tag content based on filenames.
        </Form.Text>
      </Form.Group>
    </>
  );
};
