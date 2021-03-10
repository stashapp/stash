import React, { useState } from "react";
import { Button, Col, Form, OverlayTrigger, Popover, Row } from "react-bootstrap";
import { FormUtils } from "src/utils";
import { Modal } from ".";
import Icon from "./Icon";

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onImageURL?: (url: string) => void;
  acceptSVG?: boolean;
}

export const ImageInput: React.FC<IImageInput> = ({
  isEditing,
  text,
  onImageChange,
  onImageURL,
  acceptSVG = false,
}) => {
  const [isShowDialog, setIsShowDialog] = useState(false);
  const [url, setURL] = useState("");

  if (!isEditing) return <div />;

  if (!onImageURL) {
    // just return the file input
    return (
      <Form.Label className="image-input">
        <Button variant="secondary">{text ?? "Browse for image..."}</Button>
        <Form.Control
          type="file"
          onChange={onImageChange}
          accept={`.jpg,.jpeg,.png${acceptSVG ? ",.svg" : ""}`}
        />
      </Form.Label>
    )
  }

  function onConfirmURL() {
    if (!onImageURL) {
      return;
    }

    setIsShowDialog(false);
    onImageURL(url);
  }

  function renderDialog() {
    return (
      <Modal
        show={!!isShowDialog}
        onHide={() => setIsShowDialog(false)}
        header="Image URL"
        accept={{ onClick: onConfirmURL, text: "Confirm" }}
      >
        <div className="dialog-content">
          <Form.Group controlId="url" as={Row}>
            <Form.Label column xs={3}>
              URL
            </Form.Label>
            <Col xs={9}>
              <Form.Control
                className="text-input"
                onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
                  setURL(event.currentTarget.value)
                }
                value={url}
                placeholder="URL"
              />
            </Col>
          </Form.Group>
        </div>
      </Modal>
    );
  }

  const popover = (
    <Popover id="set-image-popover">
      <Popover.Content>
        <>
          <div>
            <Form.Label className="image-input">
              <Button variant="secondary">
                <span className="fa-icon">
                  <Icon icon="file" />
                </span>
                <span>
                  From file...
                </span>
              </Button>
              <Form.Control
                type="file"
                onChange={onImageChange}
                accept={`.jpg,.jpeg,.png${acceptSVG ? ",.svg" : ""}`}
              />
            </Form.Label>
          </div>
          <div>
            <Button className="minimal" onClick={() => setIsShowDialog(true)}>
              <span className="fa-icon">
                <Icon icon="link" />
              </span>
              <span>From URL...</span>
            </Button>
          </div>
        </>
      </Popover.Content>
    </Popover>
  );

  return (
    <>
      {renderDialog()}
      <OverlayTrigger trigger="click" placement="top" overlay={popover} rootClose>
        <Button variant="secondary" className="mr-2">
        {text ?? "Set image..."}
        </Button>
      </OverlayTrigger>
    </>
  );
};
