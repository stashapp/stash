import { Button, Form, Col, Row } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { Icon } from "../Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import * as FormUtils from "src/utils/form";
import { useTagsMerge } from "src/core/StashService";
import { useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { faExchangeAlt, faSignInAlt } from "@fortawesome/free-solid-svg-icons";
import { Tag, TagSelect } from "./TagSelect";

interface ITagMergeModalProps {
  show: boolean;
  onClose: (mergedID?: string) => void;
  tags: Tag[];
}

export const TagMergeModal: React.FC<ITagMergeModalProps> = ({
  show,
  onClose,
  tags,
}) => {
  const [src, setSrc] = useState<Tag[]>([]);
  const [dest, setDest] = useState<Tag | null>(null);

  const [running, setRunning] = useState(false);

  const [mergeTags] = useTagsMerge();

  const intl = useIntl();
  const Toast = useToast();

  const title = intl.formatMessage({
    id: "actions.merge",
  });

  useEffect(() => {
    if (tags.length > 0) {
      setDest(tags[0]);
      setSrc(tags.slice(1));
    }
  }, [tags]);

  async function onMerge() {
    if (!dest) return;

    const source = src.map((s) => s.id);
    const destination = dest.id;

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
        onClose(dest.id);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return src.length > 0 && dest !== null;
  }

  function switchTags() {
    if (src.length && dest !== null) {
      const newDest = src[0];
      setSrc([...src.slice(1), dest]);
      setDest(newDest);
    }
  }

  return (
    <ModalComponent
      show={show}
      header={title}
      icon={faSignInAlt}
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
          <Form.Group controlId="source" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "dialogs.merge.source" }),
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
                onSelect={(items) => setSrc(items)}
                values={src}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
          <Form.Group
            controlId="switch"
            as={Row}
            className="justify-content-center"
          >
            <Button
              variant="secondary"
              onClick={() => switchTags()}
              disabled={!src.length || !dest}
              title={intl.formatMessage({ id: "actions.swap" })}
            >
              <Icon className="fa-fw" icon={faExchangeAlt} />
            </Button>
          </Form.Group>
          <Form.Group controlId="destination" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({
                id: "dialogs.merge.destination",
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
                onSelect={(items) => setDest(items[0])}
                values={dest ? [dest] : undefined}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </ModalComponent>
  );
};
