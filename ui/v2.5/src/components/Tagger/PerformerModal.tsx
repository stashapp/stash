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

  const renderField = (name: string) => (
    <div className="col-6 performer-create-modal-field" key={name}>
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
  );

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
        <div className="col-6">
          <div className="row no-gutters">
            {renderField("name")}
            <TruncatedText className="col-6" text={performer.name} />
          </div>
          <div className="row no-gutters">
            {renderField("gender")}
            <TruncatedText
              className="col-6 text-capitalize"
              text={performer.gender && genderToString(performer.gender)}
            />
          </div>
          <div className="row no-gutters">
            {renderField("birthdate")}
            <TruncatedText
              className="col-6"
              text={performer.birthdate ?? "Unknown"}
            />
          </div>
          <div className="row no-gutters">
            {renderField("death_date")}
            <TruncatedText
              className="col-6"
              text={performer.death_date ?? "Unknown"}
            />
          </div>
          <div className="row no-gutters">
            {renderField("ethnicity")}
            <TruncatedText
              className="col-6 text-capitalize"
              text={performer.ethnicity}
            />
          </div>
          <div className="row no-gutters">
            {renderField("country")}
            <TruncatedText className="col-6" text={performer.country ?? ""} />
          </div>
          <div className="row no-gutters">
            {renderField("hair_color")}
            <TruncatedText
              className="col-6 text-capitalize"
              text={performer.hair_color}
            />
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Eye Color:</strong>
            {renderField("eye_color")}
            <TruncatedText
              className="col-6 text-capitalize"
              text={performer.eye_color}
            />
          </div>
          <div className="row no-gutters">
            {renderField("height")}
            <TruncatedText className="col-6" text={performer.height} />
          </div>
          <div className="row no-gutters">
            {renderField("weight")}
            <TruncatedText className="col-6" text={performer.weight} />
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Measurements:</strong>
            {renderField("measurements")}
            <TruncatedText className="col-6" text={performer.measurements} />
          </div>
          {performer?.gender !== GQL.GenderEnum.Male && (
            <div className="row no-gutters">
              {renderField("fake_tits")}
              <TruncatedText className="col-6" text={performer.fake_tits} />
            </div>
          )}
          <div className="row no-gutters">
            {renderField("career_length")}
            <TruncatedText className="col-6" text={performer.career_length} />
          </div>
          <div className="row no-gutters">
            {renderField("tattoos")}
            <TruncatedText className="col-6" text={performer.tattoos} />
          </div>
          <div className="row no-gutters ">
            {renderField("piercings")}
            <TruncatedText className="col-6" text={performer.piercings} />
          </div>
        </div>
        {images.length > 0 && (
          <div className="col-6 image-selection">
            <div className="performer-image">
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
            <div className="d-flex mt-2">
              <Button
                className="mr-auto"
                onClick={setPrev}
                disabled={images.length === 1}
              >
                <Icon icon="arrow-left" />
              </Button>
              <h5>
                Select performer image
                <br />
                {imageIndex + 1} of {images.length}
              </h5>
              <Button
                className="ml-auto"
                onClick={setNext}
                disabled={images.length === 1}
              >
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
