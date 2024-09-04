import React, { useEffect, useState } from "react";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  mutateSubmitStashBoxPerformerDraft,
  mutateSubmitStashBoxSceneDraft,
} from "src/core/StashService";
import { ModalComponent } from "src/components/Shared/Modal";
import { getStashboxBase } from "src/utils/stashbox";
import { FormattedMessage, useIntl } from "react-intl";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { ExternalLink } from "../Shared/ExternalLink";

interface IProps {
  type: "scene" | "performer";
  entity: Pick<
    GQL.SceneDataFragment | GQL.PerformerDataFragment,
    "id" | "stash_ids"
  >;
  boxes: Pick<GQL.StashBox, "name" | "endpoint">[];
  show: boolean;
  onHide: () => void;
}

export const SubmitStashBoxDraft: React.FC<IProps> = ({
  type,
  boxes,
  entity,
  show,
  onHide,
}) => {
  const intl = useIntl();

  const [selectedBoxIndex, setSelectedBoxIndex] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();
  const [reviewUrl, setReviewUrl] = useState<string>();

  // this can be undefined, if e.g. boxes is empty
  // since we aren't using noUncheckedIndexedAccess, add undefined explicitly
  const selectedBox: (typeof boxes)[number] | undefined =
    boxes[selectedBoxIndex];

  // #4354: reset state when shown, or if any props change
  useEffect(() => {
    if (show) {
      setSelectedBoxIndex(0);
      setLoading(false);
      setError(undefined);
      setReviewUrl(undefined);
    }
  }, [show, type, boxes, entity]);

  async function doSubmit() {
    if (!selectedBox) return;

    const input = {
      id: entity.id,
      stash_box_endpoint: selectedBox.endpoint,
    };

    if (type === "scene") {
      const r = await mutateSubmitStashBoxSceneDraft(input);
      return r.data?.submitStashBoxSceneDraft;
    } else if (type === "performer") {
      const r = await mutateSubmitStashBoxPerformerDraft(input);
      return r.data?.submitStashBoxPerformerDraft;
    }
  }

  async function onSubmit() {
    if (!selectedBox) return;

    try {
      setLoading(true);
      const responseId = await doSubmit();

      const stashboxBase = getStashboxBase(selectedBox.endpoint);
      if (responseId) {
        setReviewUrl(`${stashboxBase}drafts/${responseId}`);
      } else {
        // if the mutation returned a null id but didn't error, then just link to the drafts page
        setReviewUrl(`${stashboxBase}drafts`);
      }
    } catch (e) {
      if (e instanceof Error && e.message) {
        setError(e.message);
      } else {
        setError(String(e));
      }
    } finally {
      setLoading(false);
    }
  }

  function renderContents() {
    if (error !== undefined) {
      return (
        <>
          <h6 className="mt-2">
            <FormattedMessage id="stashbox.submission_failed" />
          </h6>
          <div>{error}</div>
        </>
      );
    } else if (reviewUrl !== undefined) {
      return (
        <>
          <h6>
            <FormattedMessage id="stashbox.submission_successful" />
          </h6>
          <div>
            <ExternalLink href={reviewUrl}>
              <FormattedMessage
                id="stashbox.go_review_draft"
                values={{ endpoint_name: selectedBox?.name }}
              />
            </ExternalLink>
          </div>
        </>
      );
    } else {
      return (
        <Form.Group className="form-row align-items-end">
          <Form.Label className="col-6">
            <FormattedMessage id="stashbox.selected_stash_box" />:
          </Form.Label>
          <Form.Control
            as="select"
            onChange={(e) => setSelectedBoxIndex(Number(e.currentTarget.value))}
            value={selectedBoxIndex}
            className="col-6 input-control"
          >
            {boxes.map((box, i) => (
              <option value={i} key={`${box.endpoint}-${i}`}>
                {box.name}
              </option>
            ))}
          </Form.Control>
        </Form.Group>
      );
    }
  }

  function getFooterProps() {
    if (error !== undefined || reviewUrl !== undefined) {
      return {
        accept: {
          onClick: () => onHide(),
        },
      };
    }

    // If the scene has an attached stash_id from that endpoint, the operation will be an update
    const isUpdate =
      entity.stash_ids.find((id) => id.endpoint === selectedBox?.endpoint) !==
      undefined;

    return {
      footerButtons: isUpdate && !loading && (
        <span className="mr-2 align-middle">
          <FormattedMessage
            id="stashbox.submit_update"
            values={{ endpoint_name: selectedBox?.name }}
          />
        </span>
      ),
      accept: {
        onClick: () => onSubmit(),
        text: intl.formatMessage({
          id: isUpdate ? "actions.submit_update" : "actions.submit",
        }),
        variant: isUpdate ? "primary" : "success",
      },
      cancel: {
        onClick: () => onHide(),
        variant: "secondary",
      },
    };
  }

  return (
    <ModalComponent
      icon={faPaperPlane}
      header={intl.formatMessage({ id: "actions.submit_stash_box" })}
      isRunning={loading}
      show={show}
      onHide={onHide}
      disabled={!selectedBox}
      {...getFooterProps()}
    >
      {renderContents()}
    </ModalComponent>
  );
};

export default SubmitStashBoxDraft;
