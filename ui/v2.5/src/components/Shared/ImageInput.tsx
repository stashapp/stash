import React, { useState } from "react";
import {
  Button,
  Col,
  Form,
  OverlayTrigger,
  Popover,
  Row,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { ModalComponent } from "./Modal";
import { Icon } from "./Icon";
import {
  faClipboard,
  faFile,
  faLink,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { PatchComponent } from "src/patch";
import ImageUtils from "src/utils/image";
import { useToast } from "src/hooks/Toast";

export interface IImageInputExtraAction {
  icon: IconDefinition;
  labelId: string;
  onClick: () => void | Promise<void>;
}

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onImageURL?: (url: string) => void;
  onReset?: () => void;
  acceptSVG?: boolean;
  extraActions?: IImageInputExtraAction[];
}

function acceptExtensions(acceptSVG: boolean = false) {
  return `.jpg,.jpeg,.png,.webp,.gif${acceptSVG ? ",.svg" : ""}`;
}

export const ImageInput: React.FC<IImageInput> = PatchComponent(
  "ImageInput",
  ({
    isEditing,
    text,
    onImageChange,
    onImageURL,
    onReset,
    acceptSVG = false,
    extraActions,
  }) => {
    const [isShowDialog, setIsShowDialog] = useState(false);
    const [url, setURL] = useState("");
    const intl = useIntl();
    const Toast = useToast();
    if (!isEditing) return <div />;

    if (!onImageURL) {
      // just return the file input
      return (
        <Form.Label className="image-input">
          <Button variant="secondary">
            {text ?? <FormattedMessage id="actions.browse_for_image" />}
          </Button>
          <Form.Control
            type="file"
            onChange={onImageChange}
            accept={acceptExtensions(acceptSVG)}
          />
        </Form.Label>
      );
    }

    async function onPasteClipboard() {
      try {
        const data = await ImageUtils.readClipboardImage();
        if (data && onImageURL) {
          onImageURL(data);
          Toast.success(
            intl.formatMessage({ id: "toast.clipboard_image_pasted" })
          );
        } else {
          Toast.error(intl.formatMessage({ id: "toast.clipboard_no_image" }));
        }
      } catch (e) {
        if (e instanceof DOMException && e.name === "NotAllowedError") {
          Toast.error(
            intl.formatMessage({ id: "toast.clipboard_access_denied" })
          );
        } else {
          Toast.error(e);
        }
      }
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
                <FormattedMessage id="url" />
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
          <div>
            <span className="image-input">
              <Button className="minimal">
                <Icon icon={faFile} className="fa-fw" />
                <span>
                  <FormattedMessage id="actions.from_file" />
                </span>
                <Form.Control
                  type="file"
                  onChange={onImageChange}
                  accept={acceptExtensions(acceptSVG)}
                />
              </Button>
            </span>
          </div>
          <div>
            <Button className="minimal" onClick={showDialog}>
              <Icon icon={faLink} className="fa-fw" />
              <span>
                <FormattedMessage id="actions.from_url" />
              </span>
            </Button>
          </div>
          {window.isSecureContext && (
            <div>
              <Button className="minimal" onClick={onPasteClipboard}>
                <Icon icon={faClipboard} className="fa-fw" />
                <span>
                  <FormattedMessage id="actions.from_clipboard" />
                </span>
              </Button>
            </div>
          )}
          {extraActions && extraActions.length > 0 && (
            <div className="set-image-menu-divider" />
          )}
          {extraActions?.map((action) => (
            <div key={action.labelId}>
              <Button className="minimal" onClick={action.onClick}>
                <Icon icon={action.icon} className="fa-fw" />
                <span>
                  <FormattedMessage id={action.labelId} />
                </span>
              </Button>
            </div>
          ))}
          {onReset && (
            <>
              <div className="set-image-menu-divider" />
              <div>
                <Button className="minimal" onClick={onReset}>
                  <Icon icon={faTrashAlt} className="fa-fw" />
                  <span>
                    <FormattedMessage id="actions.clear_image" />
                  </span>
                </Button>
              </div>
            </>
          )}
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
            {text ?? <FormattedMessage id="actions.set_image" />}
          </Button>
        </OverlayTrigger>
      </>
    );
  }
);
