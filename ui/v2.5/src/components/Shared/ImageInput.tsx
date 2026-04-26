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
import {
  faCameraRotate,
  faCirclePlay,
  faClipboard,
  faFile,
  faLink,
  faPhotoFilm,
} from "@fortawesome/free-solid-svg-icons";
import { PatchComponent } from "src/patch";
import ImageUtils from "src/utils/image";
import { useToast } from "src/hooks/Toast";

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onImageURL?: (url: string) => void;
  onImageURLSource?: (source: "url" | "clipboard", value?: string) => void;
  onReset?: () => void;
  acceptSVG?: boolean;
  onGenerateDefault?: () => void;
  onGenerateCurrent?: () => void;
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
    onImageURLSource,
    onReset,
    acceptSVG = false,
    onGenerateDefault,
    onGenerateCurrent,
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

    async function onPasteClipboard() {
      try {
        const data = await ImageUtils.readClipboardImage();
        if (data && onImageURL) {
          onImageURL(data);
          onImageURLSource?.("clipboard");
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
      onImageURLSource?.("url", url);
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
            {window.isSecureContext && (
              <div>
                <Button className="minimal" onClick={onPasteClipboard}>
                  <Icon icon={faClipboard} className="fa-fw" />
                  <span>
                    {intl.formatMessage({ id: "actions.from_clipboard" })}
                  </span>
                </Button>
              </div>
            )}
            {(onGenerateDefault || onGenerateCurrent) && <hr className="my-2" />}
            {onGenerateDefault && (
              <div>
                <Button className="minimal" onClick={onGenerateDefault}>
                  <Icon icon={faPhotoFilm} className="fa-fw" />
                  <span>
                    {intl.formatMessage({ id: "actions.generate_thumb_default" })}
                  </span>
                </Button>
              </div>
            )}
            {onGenerateCurrent && (
              <div>
                <Button className="minimal" onClick={onGenerateCurrent}>
                  <Icon icon={faCameraRotate} className="fa-fw" />
                  <span>
                    {intl.formatMessage({
                      id: "actions.generate_thumb_from_current",
                    })}
                  </span>
                </Button>
              </div>
            )}
            {onReset && (
              <>
                <hr className="my-2" />
                <div>
                  <Button className="minimal" onClick={onReset}>
                    <Icon icon={faCirclePlay} className="fa-solid" />
                    <span>{intl.formatMessage({ id: "actions.clear_image" })}</span>
                  </Button>
                </div>
              </>
            )}
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
