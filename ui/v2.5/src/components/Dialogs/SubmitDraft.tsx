import React, { useState } from "react";
import { useMutation, DocumentNode } from "@apollo/client";
import { Button, Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { getStashboxBase } from "src/utils";

interface IProps {
  show: boolean;
  entity: { name?: string | null; id: string; title?: string | null };
  boxes: Pick<GQL.StashBox, "name" | "endpoint">[];
  query: DocumentNode;
  onHide: () => void;
}

type Variables =
  | GQL.SubmitStashBoxSceneDraftMutationVariables
  | GQL.SubmitStashBoxPerformerDraftMutationVariables;
type Query =
  | GQL.SubmitStashBoxSceneDraftMutation
  | GQL.SubmitStashBoxPerformerDraftMutation;

const isSceneDraft = (
  query: Query | null
): query is GQL.SubmitStashBoxSceneDraftMutation =>
  (query as GQL.SubmitStashBoxSceneDraftMutation).submitStashBoxSceneDraft !==
  undefined;

const getResponseId = (query: Query | null) =>
  isSceneDraft(query)
    ? query.submitStashBoxSceneDraft
    : query?.submitStashBoxPerformerDraft;

export const SubmitStashBoxDraft: React.FC<IProps> = ({
  show,
  boxes,
  entity,
  query,
  onHide,
}) => {
  const [submit, { data, error, loading }] = useMutation<Query, Variables>(
    query
  );
  const [selectedBox, setSelectedBox] = useState(0);

  const handleSubmit = () => {
    submit({
      variables: {
        input: {
          id: entity.id,
          stash_box_index: selectedBox,
        },
      },
    });
  };

  const handleSelectBox = (e: React.ChangeEvent<HTMLSelectElement>) =>
    setSelectedBox(Number.parseInt(e.currentTarget.value) ?? 0);

  console.log(data);

  return (
    <Modal
      icon="paper-plane"
      header="Submit to Stash-Box"
      isRunning={loading}
      show={show}
      accept={{
        onClick: onHide,
      }}
    >
      {data === undefined ? (
        <>
          <Form.Group className="form-row align-items-end">
            <Form.Label className="col-6">
              Selected Stash-Box endpoint:
            </Form.Label>
            <Form.Control
              as="select"
              onChange={handleSelectBox}
              className="col-6"
            >
              {boxes.map((box, i) => (
                <option value={i} key={`${box.endpoint}-${i}`}>
                  {box.name}
                </option>
              ))}
            </Form.Control>
          </Form.Group>
          <Button onClick={handleSubmit}>
            Submit {`"${entity.name ?? entity.title}"`}
          </Button>
        </>
      ) : (
        <>
          <h6>Submission successful</h6>
          <div>
            <a
              target="_blank"
              rel="noreferrer noopener"
              href={`${getStashboxBase(
                boxes[selectedBox].endpoint
              )}drafts/${getResponseId(data)}`}
            >
              Go to {boxes[selectedBox].name} to review draft.
            </a>
          </div>
        </>
      )}
      {error !== undefined && (
        <>
          <h6 className="mt-2">Submission failed</h6>
          <div>{error.message}</div>
        </>
      )}
    </Modal>
  );
};
