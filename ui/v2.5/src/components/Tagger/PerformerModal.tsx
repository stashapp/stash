import React, { useState } from "react";
import { Button } from "react-bootstrap";
import cx from "classnames";
import { IconName } from "@fortawesome/fontawesome-svg-core";

import {
  LoadingIndicator,
  Icon,
  Modal,
  TruncatedText,
} from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { genderToString } from "src/core/StashService";
import { TextUtils } from "src/utils";
import { IStashBoxPerformer } from "./utils";

interface IPerformerModalProps {
  performer: IStashBoxPerformer;
  modalVisible: boolean;
  closeModal: () => void;
  handlePerformerCreate: (imageIndex: number, excludedFields: string[]) => void;
  excludedPerformerFields?: string[];
  header: string;
  icon: IconName;
  create?: boolean;
  endpoint: string;
}

const PerformerModal: React.FC<IPerformerModalProps> = ({
  modalVisible,
  performer,
  handlePerformerCreate,
  closeModal,
  excludedPerformerFields = [],
  header,
  icon,
  create = false,
  endpoint,
}) => {
  const [imageIndex, setImageIndex] = useState(0);
  const [imageState, setImageState] = useState<
    "loading" | "error" | "loaded" | "empty"
  >("empty");
  const [loadDict, setLoadDict] = useState<Record<number, boolean>>({});
  const [excluded, setExcluded] = useState<Record<string, boolean>>(
    excludedPerformerFields.reduce(
      (dict, field) => ({ ...dict, [field]: true }),
      {}
    )
  );

  const { images } = performer;

  const changeImage = (index: number) => {
    setImageIndex(index);
    if (!loadDict[index]) setImageState("loading");
  };
  const setPrev = () =>
    changeImage(imageIndex === 0 ? images.length - 1 : imageIndex - 1);
  const setNext = () =>
    changeImage(imageIndex === images.length - 1 ? 0 : imageIndex + 1);

  const handleLoad = (index: number) => {
    setLoadDict({
      ...loadDict,
      [index]: true,
    });
    setImageState("loaded");
  };
  const handleError = () => setImageState("error");

  const toggleField = (name: string) =>
    setExcluded({
      ...excluded,
      [name]: !excluded[name],
    });

  const renderField = (
    name: string,
    text: string | null | undefined,
    truncate: boolean = true
  ) =>
    text && (
      <div className="row no-gutters">
        <div className="col-5 performer-create-modal-field" key={name}>
          {!create && (
            <Button
              onClick={() => toggleField(name)}
              variant="secondary"
              className={excluded[name] ? "text-muted" : "text-success"}
            >
              <Icon icon={excluded[name] ? "times" : "check"} />
            </Button>
          )}
          <strong>{TextUtils.capitalize(name)}:</strong>
        </div>
        {truncate ? (
          <TruncatedText className="col-7" text={text} />
        ) : (
          <span className="col-7">{text}</span>
        )}
      </div>
    );

  const base = endpoint.match(/https?:\/\/.*?\//)?.[0];
  const link = base ? `${base}performers/${performer.stash_id}` : undefined;

  return (
    <Modal
      show={modalVisible}
      accept={{
        text: "Save",
        onClick: () =>
          handlePerformerCreate(
            imageIndex,
            create ? [] : Object.keys(excluded).filter((key) => excluded[key])
          ),
      }}
      cancel={{ onClick: () => closeModal(), variant: "secondary" }}
      onHide={() => closeModal()}
      dialogClassName="performer-create-modal"
      icon={icon}
      header={header}
    >
      <div className="row">
        <div className="col-7">
          {renderField("name", performer.name)}
          {renderField("gender", genderToString(performer.gender))}
          {renderField("birthdate", performer.birthdate ?? "Unknown")}
          {renderField("death_date", performer.death_date ?? "Unknown")}
          {renderField("ethnicity", performer.ethnicity)}
          {renderField("country", performer.country)}
          {renderField("hair_color", performer.hair_color)}
          {renderField("eye_color", performer.eye_color)}
          {renderField("height", performer.height)}
          {renderField("weight", performer.weight)}
          {renderField("measurements", performer.measurements)}
          {performer?.gender !== GQL.GenderEnum.Male &&
            renderField("fake_tits", performer.fake_tits)}
          {renderField("career_length", performer.career_length)}
          {renderField("tattoos", performer.tattoos, false)}
          {renderField("piercings", performer.piercings, false)}
          {link && (
            <h6 className="mt-2">
              <a href={link} target="_blank" rel="noopener noreferrer">
                Stash-Box Source
                <Icon icon="external-link-alt" className="ml-2" />
              </a>
            </h6>
          )}
        </div>
        {images.length > 0 && (
          <div className="col-5 image-selection">
            <div className="performer-image">
              {!create && (
                <Button
                  onClick={() => toggleField("image")}
                  variant="secondary"
                  className={cx(
                    "performer-image-exclude",
                    excluded.image ? "text-muted" : "text-success"
                  )}
                >
                  <Icon icon={excluded.image ? "times" : "check"} />
                </Button>
              )}
              <img
                src={images[imageIndex]}
                className={cx({ "d-none": imageState !== "loaded" })}
                alt=""
                onLoad={() => handleLoad(imageIndex)}
                onError={handleError}
              />
              {imageState === "loading" && (
                <LoadingIndicator message="Loading image..." />
              )}
              {imageState === "error" && (
                <div className="h-100 d-flex justify-content-center align-items-center">
                  <b>Error loading image.</b>
                </div>
              )}
            </div>
            <div className="d-flex mt-3">
              <Button onClick={setPrev} disabled={images.length === 1}>
                <Icon icon="arrow-left" />
              </Button>
              <h5 className="flex-grow-1">
                Select performer image
                <br />
                {imageIndex + 1} of {images.length}
              </h5>
              <Button onClick={setNext} disabled={images.length === 1}>
                <Icon icon="arrow-right" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default PerformerModal;
