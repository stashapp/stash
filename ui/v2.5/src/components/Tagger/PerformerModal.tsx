import React, { useState } from "react";
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
import { stringToGender } from "src/utils/gender";
import { getCountryByISO } from "src/utils/country";

interface IPerformerModalProps {
  performer: GQL.ScrapedScenePerformerDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  onSave: (input: GQL.PerformerCreateInput) => void;
  excludedPerformerFields?: string[];
  header: string;
  icon: IconName;
  create?: boolean;
  endpoint?: string;
}

const PerformerModal: React.FC<IPerformerModalProps> = ({
  modalVisible,
  performer,
  onSave,
  closeModal,
  excludedPerformerFields = [],
  header,
  icon,
  create = false,
  endpoint,
}) => {
  const intl = useIntl();

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
          <strong>
            <FormattedMessage id={name} />:
          </strong>
        </div>
        {truncate ? (
          <TruncatedText className="col-7" text={text} />
        ) : (
          <span className="col-7">{text}</span>
        )}
      </div>
    );

  const base = endpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base
    ? `${base}performers/${performer.remote_site_id}`
    : undefined;

  function onSaveClicked() {
    if (!performer.name) {
      throw new Error("performer name must set");
    }

    const performerData: GQL.PerformerCreateInput & {
      [index: string]: unknown;
    } = {
      name: performer.name ?? "",
      aliases: performer.aliases,
      gender: stringToGender(performer.gender ?? undefined, true),
      birthdate: performer.birthdate,
      ethnicity: performer.ethnicity,
      eye_color: performer.eye_color,
      country: getCountryByISO(performer.country),
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
    if (remoteSiteID && endpoint) {
      performerData.stash_ids = [
        {
          endpoint,
          stash_id: remoteSiteID,
        },
      ];
    }

    // handle exclusions
    Object.keys(performerData).forEach((k) => {
      if (excluded[k] || !performerData[k]) {
        performerData[k] = undefined;
      }
    });

    onSave(performerData);
  }

  return (
    <Modal
      show={modalVisible}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: onSaveClicked,
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
          {renderField("aliases", performer.aliases)}
          {renderField(
            "gender",
            performer.gender
              ? intl.formatMessage({ id: "gender_types." + performer.gender })
              : ""
          )}
          {renderField("birthdate", performer.birthdate)}
          {renderField("death_date", performer.death_date)}
          {renderField("ethnicity", performer.ethnicity)}
          {renderField("country", getCountryByISO(performer.country))}
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
