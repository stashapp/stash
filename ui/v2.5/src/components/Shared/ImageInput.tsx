import React, { useState } from "react";
import {
  Button,
  Col,
  Form,
  OverlayTrigger,
  Popover,
  Row,
} from "react-bootstrap";
import { useIntl } from "react-intl";
import { ModalComponent } from "./Modal";
import { Icon } from "./Icon";
import { faFile, faLink } from "@fortawesome/free-solid-svg-icons";
import { PatchComponent } from "src/patch";

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onImageURL?: (url: string) => void;
  acceptSVG?: boolean;
}

function acceptExtensions(acceptSVG: boolean = false) {
  return `.jpg,.jpeg,.png,.webp,.gif${acceptSVG ? ",.svg" : ""}`;
}

export const ImageInput: React.FC<IImageInput> = PatchComponent(
  "ImageInput",
  ({ isEditing, text, onImageChange, onImageURL, acceptSVG = false }) => {
    const [isShowDialog, setIsShowDialog] = useState(false);
    const [url, setURL] = useState("");
    const intl = useIntl();

    if (!isEditing) return <div />;

    if (!onImageURL) {
      // just return the file input
      return (
        <Form.Label className="image-input">
          <Button variant="secondary">
            {text ?? intl.formatMessage({ id: "actions.browse_for_image" })}
          </Button>
          <Form.Control
            type="file"
            onChange={onImageChange}
            accept={acceptExtensions(acceptSVG)}
          />
        </Form.Label>
      );
    }

    function showDialog() {
      setURL("");
      setIsShowDialog(true);
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
        <ModalComponent
          show={!!isShowDialog}
          onHide={() => setIsShowDialog(false)}
          header={intl.formatMessage({ id: "dialogs.set_image_url_title" })}
          accept={{
            onClick: onConfirmURL,
            text: intl.formatMessage({ id: "actions.confirm" }),
          }}
        >
          <div className="dialog-content">
            <Form.Group controlId="url" as={Row}>
              <Form.Label column xs={3}>
                {intl.formatMessage({ id: "url" })}
              </Form.Label>
              <Col xs={9}>
                <Form.Control
                  className="text-input"
                  onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
                    setURL(event.currentTarget.value)
                  }
                  value={url}
                  placeholder={intl.formatMessage({ id: "url" })}
                />
              </Col>
            </Form.Group>
          </div>
        </ModalComponent>
      );
    }

    const popover = (
      <Popover id="set-image-popover">
        <Popover.Content>
          <>
            <div>
              <Form.Label className="image-input">
                <Button variant="secondary">
                  <Icon icon={faFile} className="fa-fw" />
                  <span>{intl.formatMessage({ id: "actions.from_file" })}</span>
                </Button>
                <Form.Control
                  type="file"
                  onChange={onImageChange}
                  accept={acceptExtensions(acceptSVG)}
                />
              </Form.Label>
            </div>
            <div>
              <Button className="minimal" onClick={showDialog}>
                <Icon icon={faLink} className="fa-fw" />
                <span>{intl.formatMessage({ id: "actions.from_url" })}</span>
              </Button>
            </div>
          </>
        </Popover.Content>
      </Popover>
    );

    return (
      <>
        {renderDialog()}
        <OverlayTrigger
          trigger="click"
          placement="top"
          overlay={popover}
          rootClose
        >
          <Button variant="secondary" className="mr-2">
            {text ?? intl.formatMessage({ id: "actions.set_image" })}
          </Button>
        </OverlayTrigger>
      </>
    );
  }
);
