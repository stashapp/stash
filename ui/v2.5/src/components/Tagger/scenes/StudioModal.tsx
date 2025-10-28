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
import { excludeFields } from "src/utils/data";
import { ExternalLink } from "src/components/Shared/ExternalLink";

interface IStudioDetailsProps {
  studio: GQL.ScrapedSceneStudioDataFragment;
  link?: string;
  excluded: Record<string, boolean>;
  toggleField: (field: string) => void;
  isNew?: boolean;
}

const StudioDetails: React.FC<IStudioDetailsProps> = ({
  studio,
  link,
  excluded,
  toggleField,
  isNew = false,
}) => {
  function maybeRenderImage() {
    if (!studio.image) return;

    return (
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
            <img src={studio.image} alt="" />
          </div>
        </div>
      </div>
    );
  }

  function maybeRenderField(
    id: string,
    text: string | null | undefined,
    isSelectable: boolean = true
  ) {
    if (!text) return;

    return (
      <div className="row no-gutters">
        <div className="col-5 studio-create-modal-field" key={id}>
          {isSelectable && (
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
        <TruncatedText className="col-7" text={text} />
      </div>
    );
  }

  function maybeRenderStashBoxLink() {
    if (!link) return;

    return (
      <h6 className="mt-2">
        <ExternalLink href={link}>
          <FormattedMessage id="stashbox.source" />
          <Icon icon={faExternalLinkAlt} className="ml-2" />
        </ExternalLink>
      </h6>
    );
  }

  return (
    <div>
      {maybeRenderImage()}
      <div className="row">
        <div className="col-12">
          {maybeRenderField("name", studio.name, !isNew)}
          {maybeRenderField("url", studio.url)}
          {maybeRenderField("parent_studio", studio.parent?.name, false)}
          {maybeRenderStashBoxLink()}
        </div>
      </div>
    </div>
  );
};

interface IStudioModalProps {
  studio: GQL.ScrapedSceneStudioDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  handleStudioCreate: (
    input: GQL.StudioCreateInput,
    parent?: GQL.StudioCreateInput
  ) => void;
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

  const [parentExcluded, setParentExcluded] = useState<Record<string, boolean>>(
    excludedStudioFields.reduce(
      (dict, field) => ({ ...dict, [field]: true }),
      {}
    )
  );
  const toggleParentField = (name: string) =>
    setParentExcluded({
      ...parentExcluded,
      [name]: !parentExcluded[name],
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

    const studioData: GQL.StudioCreateInput = {
      name: studio.name,
      url: studio.url,
      image: studio.image,
      parent_id: studio.parent?.stored_id,
    };

    // stashid handling code
    const remoteSiteID = studio.remote_site_id;
    const timeNow = new Date().toISOString();
    if (remoteSiteID && endpoint) {
      studioData.stash_ids = [
        {
          endpoint,
          stash_id: remoteSiteID,
          updated_at: timeNow,
        },
      ];
    }

    // handle exclusions
    excludeFields(studioData, excluded);

    let parentData: GQL.StudioCreateInput | undefined = undefined;

    if (createParentStudio && sendParentStudio) {
      if (!studio.parent?.name) {
        throw new Error("parent studio name must set");
      }

      parentData = {
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
            updated_at: timeNow,
          },
        ];
      }

      // handle exclusions
      // Can't exclude parent studio name when creating a new one
      parentExcluded.name = false;
      excludeFields(parentData, parentExcluded);
    }

    handleStudioCreate(studioData, parentData);
  }

  const base = endpoint?.match(/https?:\/\/.*?\//)?.[0];
  const link = base ? `${base}studios/${studio.remote_site_id}` : undefined;
  const parentLink = base
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
    if (!createParentStudio || !studio.parent) {
      return;
    }

    return (
      <StudioDetails
        studio={studio.parent}
        excluded={parentExcluded}
        toggleField={(field) => toggleParentField(field)}
        link={parentLink}
        isNew
      />
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
      <StudioDetails
        studio={studio}
        excluded={excluded}
        toggleField={(field) => toggleField(field)}
        link={link}
      />

      {maybeRenderParentStudio()}
    </ModalComponent>
  );
};

export default StudioModal;
