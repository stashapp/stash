import React, { useState } from "react";
import { Button } from "react-bootstrap";
import cx from "classnames";

import { LoadingIndicator, Icon, Modal } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { genderToString } from "src/core/StashService";
import { IStashBoxPerformer } from "./utils";

interface IPerformerModalProps {
  performer: IStashBoxPerformer;
  modalVisible: boolean;
  showModal: (show: boolean) => void;
  handlePerformerCreate: (imageIndex: number) => void;
}

const PerformerModal: React.FC<IPerformerModalProps> = ({
  modalVisible,
  performer,
  handlePerformerCreate,
  showModal,
}) => {
  const [imageIndex, setImageIndex] = useState(0);
  const [imageState, setImageState] = useState<
    "loading" | "error" | "loaded" | "empty"
  >("empty");
  const [loadDict, setLoadDict] = useState<Record<number, boolean>>({});

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

  return (
    <Modal
      show={modalVisible}
      accept={{
        text: "Save",
        onClick: () => handlePerformerCreate(imageIndex),
      }}
      cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      onHide={() => showModal(false)}
      dialogClassName="performer-create-modal"
    >
      <div className="row">
        <div className="col-6">
          <div className="row no-gutters mb-4">
            <strong>Performer information</strong>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Name:</strong>
            <span className="col-6 text-truncate">{performer.name}</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Gender:</strong>
            <span className="col-6 text-truncate text-capitalize">
              {performer.gender && genderToString(performer.gender)}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Birthdate:</strong>
            <span className="col-6 text-truncate">
              {performer.birthdate ?? "Unknown"}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Ethnicity:</strong>
            <span className="col-6 text-truncate text-capitalize">
              {performer.ethnicity}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Country:</strong>
            <span className="col-6 text-truncate">
              {performer.country ?? ""}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Eye Color:</strong>
            <span className="col-6 text-truncate text-capitalize">
              {performer.eye_color}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Height:</strong>
            <span className="col-6 text-truncate">{performer.height}</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Measurements:</strong>
            <span className="col-6 text-truncate">
              {performer.measurements}
            </span>
          </div>
          {performer?.gender !== GQL.GenderEnum.Male && (
            <div className="row no-gutters">
              <strong className="col-6">Fake Tits:</strong>
              <span className="col-6 text-truncate">{performer.fake_tits}</span>
            </div>
          )}
          <div className="row no-gutters">
            <strong className="col-6">Career Length:</strong>
            <span className="col-6 text-truncate">
              {performer.career_length}
            </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Tattoos:</strong>
            <span className="col-6 text-truncate">{performer.tattoos}</span>
          </div>
          <div className="row no-gutters ">
            <strong className="col-6">Piercings:</strong>
            <span className="col-6 text-truncate">{performer.piercings}</span>
          </div>
        </div>
        {images.length > 0 && (
          <div className="col-6 image-selection">
            <div className="performer-image">
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
