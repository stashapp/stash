import React, { useState } from "react";
import { useMutation, DocumentNode } from "@apollo/client";
import { Button, Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { getStashboxBase } from "src/utils";
import { FormattedMessage, useIntl } from "react-intl";

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
  const intl = useIntl();

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

  return (
    <Modal
      icon="paper-plane"
      header={intl.formatMessage({ id: "actions.submit_stash_box" })}
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
              <FormattedMessage id="stashbox.selected_stash_box" />:
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
            <FormattedMessage id="actions.submit" />{" "}
            {`"${entity.name ?? entity.title}"`}
          </Button>
        </>
      ) : (
        <>
          <h6>
            <FormattedMessage id="stashbox.submission_successful" />
          </h6>
          <div>
            <a
              target="_blank"
              rel="noreferrer noopener"
              href={`${getStashboxBase(
                boxes[selectedBox].endpoint
              )}drafts/${getResponseId(data)}`}
            >
              <FormattedMessage
                id="stashbox.go_review_draft"
                values={{ endpoint_name: boxes[selectedBox].name }}
              />
            </a>
          </div>
        </>
      )}
      {error !== undefined && (
        <>
          <h6 className="mt-2">
            <FormattedMessage id="stashbox.submission_failed" />
          </h6>
          <div>{error.message}</div>
        </>
      )}
    </Modal>
  );
};
