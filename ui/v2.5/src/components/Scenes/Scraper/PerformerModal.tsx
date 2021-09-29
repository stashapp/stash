import React, { useState, useContext } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import cx from "classnames";
import { IconName } from "@fortawesome/fontawesome-svg-core";

import {
  LoadingIndicator,
  Icon,
  Modal,
  TruncatedText,
} from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { genderToString, stringToGender } from "src/utils/gender";
import { SceneScraperStateContext } from "./context";

interface IPerformerModalProps {
  performer: GQL.ScrapedScenePerformerDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  handlePerformerCreate: (input: GQL.PerformerCreateInput) => void;
  header: string;
  icon: IconName;
}

const PerformerModal: React.FC<IPerformerModalProps> = ({
  modalVisible,
  performer,
  handlePerformerCreate,
  closeModal,
  header,
  icon,
}) => {
  const { currentSource } = useContext(SceneScraperStateContext);
  const intl = useIntl();
  const [imageIndex, setImageIndex] = useState(0);
  const [imageState, setImageState] = useState<
    "loading" | "error" | "loaded" | "empty"
  >("empty");
  const [loadDict, setLoadDict] = useState<Record<number, boolean>>({});

  const images = performer.images ?? [];

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

  const renderField = (
    id: string,
    text: string | null | undefined,
    truncate: boolean = true
  ) =>
    text && (
      <div className="row no-gutters">
        <div className="col-5 performer-create-modal-field" key={id}>
          <strong>
            <FormattedMessage id={id} />:
          </strong>
        </div>
        {truncate ? (
          <TruncatedText className="col-7" text={text} />
        ) : (
          <span className="col-7">{text}</span>
        )}
      </div>
    );

  function onSave() {
    if (!performer.name) {
      throw new Error("performer name must set");
    }

    const performerData: GQL.PerformerCreateInput = {
      name: performer.name ?? "",
      aliases: performer.aliases,
      gender: stringToGender(performer.gender ?? undefined),
      birthdate: performer.birthdate,
      ethnicity: performer.ethnicity,
      eye_color: performer.eye_color,
      country: performer.country,
      height: performer.height,
      measurements: performer.measurements,
      fake_tits: performer.fake_tits,
      career_length: performer.career_length,
      tattoos: performer.tattoos,
      piercings: performer.piercings,
      url: performer.url,
      twitter: performer.twitter,
      instagram: performer.instagram,
      image: images.length > imageIndex ? images[imageIndex] : undefined,
      details: performer.details,
      death_date: performer.death_date,
      hair_color: performer.hair_color,
      weight: Number.parseFloat(performer.weight ?? "") ?? undefined,
    };

    if (Number.isNaN(performerData.weight ?? 0)) {
      performerData.weight = undefined;
    }

    if (performer.tags) {
      performerData.tag_ids = performer.tags
        .map((t) => t.stored_id)
        .filter((t) => t) as string[];
    }

    // stashid handling code
    const remoteSiteID = performer.remote_site_id;
    if (remoteSiteID && currentSource?.stashboxEndpoint) {
      performerData.stash_ids = [
        {
          endpoint: currentSource.stashboxEndpoint,
          stash_id: remoteSiteID,
        },
      ];
    }

    handlePerformerCreate(performerData);
  }

  const base = currentSource?.stashboxEndpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base
    ? `${base}performers/${performer.remote_site_id}`
    : undefined;

  return (
    <Modal
      show={modalVisible}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: onSave,
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
          {renderField(
            "gender",
            performer.gender ? genderToString(performer.gender) : ""
          )}
          {renderField("birthdate", performer.birthdate)}
          {renderField("death_date", performer.death_date)}
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
          {renderField("weight", performer.weight, false)}
          {renderField("details", performer.details)}
          {renderField("url", performer.url)}
          {renderField("twitter", performer.twitter)}
          {renderField("instagram", performer.instagram)}
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
