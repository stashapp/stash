import React, { useState } from "react";
import { useMutation, DocumentNode } from "@apollo/client";
import { Button, Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { getStashboxBase } from "src/utils";
import { FormattedMessage, useIntl } from "react-intl";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  show: boolean;
  entity: {
    name?: string | null;
    id: string;
    title?: string | null;
    stash_ids: { stash_id: string; endpoint: string }[];
  };
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
  const [selectedBoxIndex, setSelectedBoxIndex] = useState(0);
  const intl = useIntl();

  const handleSubmit = () => {
    submit({
      variables: {
        input: {
          id: entity.id,
          stash_box_index: selectedBoxIndex,
        },
      },
    });
  };

  const selectedBox =
    boxes.length > selectedBoxIndex ? boxes[selectedBoxIndex] : undefined;

  const handleSelectBox = (e: React.ChangeEvent<HTMLSelectElement>) =>
    setSelectedBoxIndex(Number.parseInt(e.currentTarget.value) ?? 0);

  if (!selectedBox) {
    return <></>;
  }

  // If the scene has an attached stash_id from that endpoint, the operation will be an update
  const isUpdate =
    entity.stash_ids.find((id) => id.endpoint === selectedBox.endpoint) !==
    undefined;

  return (
    <Modal
      icon={faPaperPlane}
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
          <div className="text-right">
            {isUpdate && (
              <span className="mr-2">
                <FormattedMessage
                  id="stashbox.submit_update"
                  values={{ endpoint_name: boxes[selectedBoxIndex].name }}
                />
              </span>
            )}
            <Button
              onClick={handleSubmit}
              variant={isUpdate ? "primary" : "success"}
            >
              <FormattedMessage
                id={`actions.${isUpdate ? "submit_update" : "submit"}`}
              />{" "}
            </Button>
          </div>
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
                boxes[selectedBoxIndex].endpoint
              )}drafts/${getResponseId(data)}`}
            >
              <FormattedMessage
                id="stashbox.go_review_draft"
                values={{ endpoint_name: boxes[selectedBoxIndex].name }}
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

export default SubmitStashBoxDraft;
