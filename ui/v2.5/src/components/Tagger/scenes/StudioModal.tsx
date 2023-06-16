import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import cx from "classnames";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

import * as GQL from "src/core/generated-graphql";
import { useFindStudio } from "src/core/StashService";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import {
  faCheck,
  faExternalLinkAlt,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { Button, Form } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared/TruncatedText";

interface IStudioModalProps {
  studio: GQL.ScrapedSceneStudioDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  handleStudioCreate: (input: GQL.StudioCreateInput) => void;
  excludedStudioFields?: string[];
  header: string;
  icon: IconDefinition;
  endpoint?: string;
}

const StudioModal: React.FC<IStudioModalProps> = ({
  modalVisible,
  studio,
  handleStudioCreate,
  closeModal,
  excludedStudioFields = [],
  header,
  icon,
  endpoint,
}) => {
  const intl = useIntl();

  const [excluded, setExcluded] = useState<Record<string, boolean>>(
    excludedStudioFields.reduce(
      (dict, field) => ({ ...dict, [field]: true }),
      {}
    )
  );
  const toggleField = (name: string) =>
    setExcluded({
      ...excluded,
      [name]: !excluded[name],
    });

  const [createParentStudio, setCreateParentStudio] = useState<boolean>(
    !!studio.parent
  );

  let sendParentStudio = true;
  // The parent studio exists, need to check if it has a Stash ID.
  const queryResult = useFindStudio(studio.parent?.stored_id ?? "");
  if (
    queryResult.data?.findStudio?.stash_ids?.length &&
    queryResult.data?.findStudio?.stash_ids?.length > 0
  ) {
    // It already has a Stash ID, so we can skip worrying about it
    sendParentStudio = false;
  }

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

    const studioData: GQL.StudioCreateInput & {
      [index: string]: unknown;
    } = {
      name: studio.name,
      url: studio.url,
      image: studio.image,
      parent_id: studio.parent?.stored_id,
    };

    // stashid handling code
    const remoteSiteID = studio.remote_site_id;
    if (remoteSiteID && endpoint) {
      studioData.stash_ids = [
        {
          endpoint,
          stash_id: remoteSiteID,
        },
      ];
    }

    // handle exclusions
    Object.keys(studioData).forEach((k) => {
      if (excluded[k] || !studioData[k]) {
        studioData[k] = undefined;
      }
    });

    if (createParentStudio) {
      if (!studio.parent?.name) {
        throw new Error("parent studio name must set");
      }

      const parentData: GQL.StudioCreateInput & {
        [index: string]: unknown;
      } = {
        name: studio.parent?.name,
        url: studio.parent?.url,
        image: studio.parent?.image,
      };

      // stashid handling code
      const parentRemoteSiteID = studio.parent?.remote_site_id;
      if (parentRemoteSiteID && endpoint) {
        parentData.stash_ids = [
          {
            endpoint,
            stash_id: parentRemoteSiteID,
          },
        ];
      }

      // handle exclusions
      Object.keys(parentData).forEach((k) => {
        // Can't exclude parent studio name when creating a new one
        if (k != "name" && (excluded[k] || !parentData[k])) {
          parentData[k] = undefined;
        }
      });

      // Hack to not send parent data when we want to ignore the existing parent studio
      studioData.parent = null;
      if (sendParentStudio) {
        studioData.parent = parentData;
      }
    }

    handleStudioCreate(studioData);
  }

  const renderField = (
    id: string,
    text: string | null | undefined,
    is_selectable: boolean = true,
    truncate: boolean = true
  ) =>
    text && (
      <div className="row no-gutters">
        <div className="col-5 studio-create-modal-field" key={id}>
          {is_selectable && (
            <Button
              onClick={() => toggleField(id)}
              variant="secondary"
              className={excluded[id] ? "text-muted" : "text-success"}
            >
              <Icon icon={excluded[id] ? faTimes : faCheck} />
            </Button>
          )}
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

  const base = endpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base ? `${base}studios/${studio.remote_site_id}` : undefined;
  const parent_link = base
    ? `${base}studios/${studio.parent?.remote_site_id}`
    : undefined;

  function maybeRenderParentStudio() {
    // There is no parent studio or it already has a Stash ID
    if (!studio.parent || !sendParentStudio) {
      return;
    }

    return (
      <div>
        <div className="mb-4 mt-4">
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
            {renderField("name", studio.parent?.name, false)}
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

  return (
    <ModalComponent
      show={modalVisible}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: onSave,
      }}
      cancel={{ onClick: () => closeModal(), variant: "secondary" }}
      onHide={() => closeModal()}
      dialogClassName="studio-create-modal"
      icon={icon}
      header={header}
    >
      <div className="row">
        <div className="col-12 image-selection">
          <div className="studio-image">
            <Button
              onClick={() => toggleField("image")}
              variant="secondary"
              className={cx(
                "studio-image-exclude",
                excluded.image ? "text-muted" : "text-success"
              )}
            >
              <Icon icon={excluded.image ? faTimes : faCheck} />
            </Button>
            <img src={studio.image ?? ""} alt="" />
          </div>
        </div>
      </div>

      <div className="row">
        <div className="col-12">
          {renderField("name", studio.name)}
          {renderField("url", studio.url)}
          {renderField("parent_studio", studio.parent?.name, false)}
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
    </ModalComponent>
  );
};

export default StudioModal;
