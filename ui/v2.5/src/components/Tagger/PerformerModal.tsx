import React, { useState } from 'react';
import { Button } from 'react-bootstrap';

import { Icon, Modal } from 'src/components/Shared';
import { BreastTypeEnum, GenderEnum } from 'src/definitions-box/globalTypes';
import {
  SearchScene_searchScene_performers_performer as StashPerformer,
} from 'src/definitions-box/SearchScene';
import { getCountryByISO } from 'src/utils/country';
import { formatBodyModification, formatMeasurements, sortImageURLs } from './utils';

interface IPerformerModalProps {
  performer: StashPerformer;
  modalVisible: boolean;
  showModal: (show: boolean) => void;
  handlePerformerCreate: (imageIndex: number) => void;
};

const genderDict = {
  [GenderEnum.FEMALE]: "Female",
  [GenderEnum.MALE]: "Male",
  [GenderEnum.TRANSGENDER_FEMALE]: "Transgender Female",
  [GenderEnum.TRANSGENDER_MALE]: "Transgender Male",
  [GenderEnum.INTERSEX]: "Intersex",
};

const PerformerModal: React.FC<IPerformerModalProps> = ({ modalVisible, performer, handlePerformerCreate, showModal }) => {
  const [imageIndex, setImageIndex] = useState(0);

  const images = sortImageURLs(performer.images, 'portrait');

  const setPrev = () => (
    setImageIndex(imageIndex === 0 ? images.length - 1 : imageIndex - 1)
  );
  const setNext = () => (
    setImageIndex(imageIndex === images.length - 1 ? 0 : imageIndex + 1)
  );

  return (
    <Modal
      show={modalVisible}
      accept={{ text: "Save", onClick: () => handlePerformerCreate(imageIndex) }}
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
            <span className="col-6 text-truncate">{ performer.name }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Gender:</strong>
            <span className="col-6 text-truncate text-capitalize">{ performer.gender && genderDict[performer.gender] }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Birthdate:</strong>
            <span className="col-6 text-truncate">{ performer.birthdate?.date ?? 'Unknown' }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Ethnicity:</strong>
            <span className="col-6 text-truncate text-capitalize">{ performer.ethnicity?.toLowerCase() }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Country:</strong>
            <span className="col-6 text-truncate">{ getCountryByISO(performer.country) ?? '' }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Eye Color:</strong>
            <span className="col-6 text-truncate text-capitalize">{ performer.eye_color?.toLowerCase() }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Height:</strong>
            <span className="col-6 text-truncate">{ performer.height }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Measurements:</strong>
            <span className="col-6 text-truncate">{ formatMeasurements(performer.measurements) }</span>
          </div>
          { performer?.gender !== GenderEnum.MALE && (
            <div className="row no-gutters">
              <strong className="col-6">Fake Tits:</strong>
              <span className="col-6 text-truncate">{ performer.breast_type === BreastTypeEnum.FAKE ? "Yes" : "No" }</span>
            </div>
          )}
          <div className="row no-gutters">
            <strong className="col-6">Career Length:</strong>
            <span className="col-6 text-truncate">{
              (performer.career_start_year) &&
              `${performer.career_start_year} - ${ performer.career_end_year ?? ''}`} </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Tattoos:</strong>
            <span className="col-6 text-truncate">{ formatBodyModification(performer.tattoos) }</span>
          </div>
          <div className="row no-gutters ">
            <strong className="col-6">Piercings:</strong>
            <span className="col-6 text-truncate">{ formatBodyModification(performer.piercings) }</span>
          </div>
        </div>
        { images.length > 0 && (
          <div className="col-6 image-selection">
            <div className="performer-image">
              <img src={images[imageIndex].url} alt='' />
            </div>
            <div className="d-flex mt-2">
              <Button className="mr-auto" onClick={setPrev} disabled={images.length === 1}>
                <Icon icon="arrow-left" />
              </Button>
              <h5>Select performer image<br />{imageIndex+1} of {images.length}</h5>
              <Button className="ml-auto" onClick={setNext} disabled={images.length === 1}>
                <Icon icon="arrow-right" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}

export default PerformerModal;
