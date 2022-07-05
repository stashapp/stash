import React, { useContext } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

import * as GQL from "src/core/generated-graphql";
import { Icon, Modal, TruncatedText } from "src/components/Shared";
import { TaggerStateContext } from "../context";
import { faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";

interface IStudioModalProps {
  studio: GQL.ScrapedSceneStudioDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  handleStudioCreate: (input: GQL.StudioCreateInput) => void;
  header: string;
  icon: IconDefinition;
}

const StudioModal: React.FC<IStudioModalProps> = ({
  modalVisible,
  studio,
  handleStudioCreate,
  closeModal,
  header,
  icon,
}) => {
  const { currentSource } = useContext(TaggerStateContext);
  const intl = useIntl();

  function onSave() {
    if (!studio.name) {
      throw new Error("studio name must set");
    }

    const studioData: GQL.StudioCreateInput = {
      name: studio.name ?? "",
      url: studio.url,
    };

    // stashid handling code
    const remoteSiteID = studio.remote_site_id;
    if (remoteSiteID && currentSource?.stashboxEndpoint) {
      studioData.stash_ids = [
        {
          endpoint: currentSource.stashboxEndpoint,
          stash_id: remoteSiteID,
        },
      ];
    }

    handleStudioCreate(studioData);
  }

  const renderField = (
    id: string,
    text: string | null | undefined,
    truncate: boolean = true
  ) =>
    text && (
      <div className="row no-gutters">
        <div className="col-5 studio-create-modal-field" key={id}>
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

  const base = currentSource?.stashboxEndpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base ? `${base}studios/${studio.remote_site_id}` : undefined;

  return (
    <Modal
      show={modalVisible}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: onSave,
      }}
      onHide={() => closeModal()}
      cancel={{ onClick: () => closeModal(), variant: "secondary" }}
      icon={icon}
      header={header}
    >
      <div className="row">
        <div className="col-12">
          {renderField("name", studio.name)}
          {renderField("url", studio.url)}
          {link && (
            <h6 className="mt-2">
              <a href={link} target="_blank" rel="noopener noreferrer">
                Stash-Box Source
                <Icon icon={faExternalLinkAlt} className="ml-2" />
              </a>
            </h6>
          )}
        </div>
      </div>

      {/* TODO - add image */}
      {/* <div className="row">
        <strong className="col-2">Logo:</strong>
        <span className="col-10">
          <img src={studio?.image ?? ""} alt="" />
        </span>
      </div> */}
    </Modal>
  );
};

export default StudioModal;
