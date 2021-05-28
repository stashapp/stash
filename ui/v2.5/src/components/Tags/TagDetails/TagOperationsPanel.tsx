import { Button, Form, Col, Row } from "react-bootstrap";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { mutateMetadataAutoTag, useTagsMerge } from "src/core/StashService";
import { Modal, TagSelect } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";

interface ITagOperationsProps {
  tag: Partial<GQL.TagDataFragment>;
}

export const TagOperationsPanel: React.FC<ITagOperationsProps> = ({ tag }) => {
  const Toast = useToast();
  const history = useHistory();

  const [mergeTags] = useTagsMerge();

  const [mergeModalOpen, setMergeModalOpen] = useState<string | null>(null);
  const [mergeSourceIds, setMergeSourceIds] = useState<string[]>([]);
  const [mergeDestinationId, setMergeDestinationId] = useState<string | null>(
    null
  );

  async function onAutoTag() {
    if (!tag?.id) return;
    try {
      await mutateMetadataAutoTag({ tags: [tag.id] });
      Toast.success({ content: "Started auto tagging" });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onMerge() {
    try {
      const result = await mergeTags({
        variables: {
          source: mergeSourceIds,
          destination: mergeDestinationId ?? "",
        },
      });
      if (result.data?.tagsMerge) {
        Toast.success({ content: "Merged tags" });
        setMergeModalOpen(null);
        history.push(`/tags/${mergeDestinationId}`);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderMergeModal() {
    return (
      <>
        <Modal
          show={mergeModalOpen !== null}
          icon={mergeModalOpen === "from" ? "sign-in-alt" : "sign-out-alt"}
          accept={{ text: "Merge", onClick: onMerge }}
          cancel={{
            variant: "secondary",
            onClick: () => setMergeModalOpen(null),
          }}
        >
          <div className="form-container row px-3">
            <div className="col-12 col-lg-6 col-xl-12">
              {mergeModalOpen === "from" && (
                <Form.Group controlId="source" as={Row}>
                  {FormUtils.renderLabel({
                    title: "Source",
                    labelProps: {
                      column: true,
                      sm: 3,
                      xl: 12,
                    },
                  })}
                  <Col sm={9} xl={12}>
                    <TagSelect
                      isMulti
                      creatable={false}
                      onSelect={(items) =>
                        setMergeSourceIds(items.map((item) => item.id))
                      }
                      ids={mergeSourceIds}
                      excludeIds={tag?.id ? [tag.id] : []}
                    />
                  </Col>
                </Form.Group>
              )}
              {mergeModalOpen === "into" && (
                <Form.Group controlId="destination" as={Row}>
                  {FormUtils.renderLabel({
                    title: "Destination",
                    labelProps: {
                      column: true,
                      sm: 3,
                      xl: 12,
                    },
                  })}
                  <Col sm={9} xl={12}>
                    <TagSelect
                      creatable={false}
                      onSelect={(items) => setMergeDestinationId(items[0]?.id)}
                      ids={
                        mergeDestinationId ? [mergeDestinationId] : undefined
                      }
                      excludeIds={tag?.id ? [tag.id] : []}
                    />
                  </Col>
                </Form.Group>
              )}
            </div>
          </div>
        </Modal>
      </>
    );
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

      <hr />
      <h5>Merging</h5>

      <Form.Group>
        <Button
          onClick={() => {
            setMergeSourceIds([]);
            setMergeDestinationId(tag?.id ?? null);
            setMergeModalOpen("from");
          }}
        >
          Merge From
        </Button>
        <Form.Text className="text-muted">
          Merge other tags into this tag.
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          onClick={() => {
            setMergeSourceIds([tag.id as string]);
            setMergeDestinationId(null);
            setMergeModalOpen("into");
          }}
        >
          Merge Into
        </Button>
        <Form.Text className="text-muted">
          Merge this tag into another tag.
        </Form.Text>
      </Form.Group>
      {renderMergeModal()}
    </>
  );
};
