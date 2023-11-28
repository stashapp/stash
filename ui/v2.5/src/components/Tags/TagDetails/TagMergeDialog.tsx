import { Form, Col, Row } from "react-bootstrap";
import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { TagSelect } from "src/components/Shared/Select";
import * as FormUtils from "src/utils/form";
import { useTagsMerge } from "src/core/StashService";
import { useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { useHistory } from "react-router-dom";
import { faSignInAlt, faSignOutAlt } from "@fortawesome/free-solid-svg-icons";

interface ITagMergeModalProps {
  show: boolean;
  onClose: () => void;
  tag: Pick<GQL.Tag, "id">;
  mergeType: "from" | "into";
}

export const TagMergeModal: React.FC<ITagMergeModalProps> = ({
  show,
  onClose,
  tag,
  mergeType,
}) => {
  const [srcIds, setSrcIds] = useState<string[]>([]);
  const [destId, setDestId] = useState<string | null>(null);
  const [running, setRunning] = useState(false);

  const [mergeTags] = useTagsMerge();

  const intl = useIntl();
  const Toast = useToast();
  const history = useHistory();

  const title = intl.formatMessage({
    id: mergeType === "from" ? "actions.merge_from" : "actions.merge_into",
  });

  async function onMerge() {
    const source = mergeType === "from" ? srcIds : [tag.id];
    const destination = mergeType === "from" ? tag.id : destId;

    if (!destination) return;

    try {
      setRunning(true);
      const result = await mergeTags({
        variables: {
          source,
          destination,
        },
      });
      if (result.data?.tagsMerge) {
        Toast.success(intl.formatMessage({ id: "toast.merged_tags" }));
        onClose();
        history.push(`/tags/${destination}`);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return (
      (mergeType === "from" && srcIds.length > 0) ||
      (mergeType === "into" && destId)
    );
  }

  return (
    <ModalComponent
      show={show}
      header={title}
      icon={mergeType === "from" ? faSignInAlt : faSignOutAlt}
      accept={{
        text: intl.formatMessage({ id: "actions.merge" }),
        onClick: () => onMerge(),
      }}
      disabled={!canMerge()}
      cancel={{
        variant: "secondary",
        onClick: () => onClose(),
      }}
      isRunning={running}
    >
      <div className="form-container row px-3">
        <div className="col-12 col-lg-6 col-xl-12">
          {mergeType === "from" && (
            <Form.Group controlId="source" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "dialogs.merge_tags.source" }),
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
                  onSelect={(items) => setSrcIds(items.map((item) => item.id))}
                  ids={srcIds}
                  excludeIds={tag?.id ? [tag.id] : []}
                />
              </Col>
            </Form.Group>
          )}
          {mergeType === "into" && (
            <Form.Group controlId="destination" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({
                  id: "dialogs.merge_tags.destination",
                }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <TagSelect
                  isMulti={false}
                  creatable={false}
                  onSelect={(items) => setDestId(items[0]?.id)}
                  ids={destId ? [destId] : undefined}
                  excludeIds={tag?.id ? [tag.id] : []}
                />
              </Col>
            </Form.Group>
          )}
        </div>
      </div>
    </ModalComponent>
  );
};
