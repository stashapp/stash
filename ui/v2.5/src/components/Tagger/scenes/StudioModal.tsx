import React, { useContext, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

import * as GQL from "src/core/generated-graphql";
import { useFindStudio } from "src/core/StashService";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { TaggerStateContext } from "../context";
import { faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";
import { Form } from "react-bootstrap";

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

  const [createParentStudio, setCreateParentStudio] = useState<boolean>(
    !!studio.parent
  );

  const parentStudioCreateText = () => {
    if (studio.parent && studio.parent.stored_id) {
      return "actions.assign_stashid_to_parent_studio";
    }
    return "actions.create_parent_studio";
  };

  function onSave() {
    if (!studio.name) {
      throw new Error("studio name must set");
    }

    const studioData: GQL.StudioCreateInput = {
      name: studio.name,
      url: studio.url,
      image: studio.image,
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

    if (createParentStudio) {
      if (!studio.parent?.name) {
        throw new Error("parent studio name must set");
      }

      const parentData: GQL.StudioCreateInput = {
        name: studio.parent?.name,
        url: studio.parent?.url,
        image: studio.parent?.image,
      };

      // stashid handling code
      const remoteSiteID = studio.parent?.remote_site_id;
      if (remoteSiteID && currentSource?.stashboxEndpoint) {
        parentData.stash_ids = [
          {
            endpoint: currentSource.stashboxEndpoint,
            stash_id: remoteSiteID,
          },
        ];
      }

      studioData.parent = parentData;
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

  function maybeRenderParentStudio() {
    // There is no parent studio
    if (!studio.parent) {
      return false;
    } else if (
      studio.parent &&
      studio.parent.stored_id &&
      studio.parent.remote_site_id
    ) {
      // The parent studio exists, need to check if it has a Stash ID.
      const queryResult = useFindStudio(studio.parent.stored_id);
      if (
        queryResult.data?.findStudio?.stash_ids?.length &&
        queryResult.data?.findStudio?.stash_ids?.length > 0
      ) {
        // It already has a Stash ID, so we can skip worrying about it
        return false;
      }
    }

    return (
      <div>
        <div className="mb-4">
          <Form.Check
            id="create-parent"
            checked={createParentStudio}
            label={intl.formatMessage({
              id: parentStudioCreateText(),
            })}
            onChange={() => setCreateParentStudio(!createParentStudio)}
          />
        </div>
        {maybeRenderParentStudioDetails()}
      </div>
    );
  }

  function maybeRenderParentStudioDetails() {
    if (!createParentStudio) {
      return;
    }

    return (
      <div>
        <div className="row mb-4">
          <img
            className="col-12 studio-card-image"
            src={studio.parent?.image ?? ""}
            alt=""
          />
        </div>

        <div className="row">
          <div className="col-12">
            {renderField("name", studio.parent?.name)}
            {renderField("url", studio.parent?.url)}
            {parent_link && (
              <h6 className="mt-2">
                <a href={parent_link} target="_blank" rel="noopener noreferrer">
                  Stash-Box Source
                  <Icon icon={faExternalLinkAlt} className="ml-2" />
                </a>
              </h6>
            )}
          </div>
        </div>
      </div>
    );
  }

  const base = currentSource?.stashboxEndpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base ? `${base}studios/${studio.remote_site_id}` : undefined;
  const parent_link = base
    ? `${base}studios/${studio.parent?.remote_site_id}`
    : undefined;

  return (
    <ModalComponent
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
      <div className="row mb-4">
        <img
          className="col-12 studio-card-image"
          src={studio.image ?? ""}
          alt=""
        />
      </div>

      <div className="row mb-4">
        <div className="col-12">
          {renderField("name", studio.name)}
          {renderField("url", studio.url)}
          {renderField("parent_studio", studio.parent?.name)}
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
      {maybeRenderParentStudio()}
    </ModalComponent >
  );
};

export default StudioModal;
