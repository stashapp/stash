import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import {
  faCheck,
  faExternalLinkAlt,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { Button } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { excludeFields } from "src/utils/data";
import { ExternalLink } from "src/components/Shared/ExternalLink";

interface ITagModalProps {
  tag: GQL.ScrapedSceneTagDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  onSave: (input: GQL.TagCreateInput) => void;
  excludedTagFields?: string[];
  header: string;
  icon: IconDefinition;
  endpoint?: string;
}

const TagModal: React.FC<ITagModalProps> = ({
  modalVisible,
  tag,
  onSave,
  closeModal,
  excludedTagFields = [],
  header,
  icon,
  endpoint,
}) => {
  const intl = useIntl();

  const [excluded, setExcluded] = useState<Record<string, boolean>>(
    excludedTagFields.reduce((dict, field) => ({ ...dict, [field]: true }), {})
  );
  const toggleField = (name: string) =>
    setExcluded({
      ...excluded,
      [name]: !excluded[name],
    });

  function maybeRenderField(id: string, text: string | null | undefined) {
    if (!text) return;

    return (
      <div className="row no-gutters">
        <div className="col-5 studio-create-modal-field" key={id}>
          <Button
            onClick={() => toggleField(id)}
            variant="secondary"
            className={excluded[id] ? "text-muted" : "text-success"}
          >
            <Icon icon={excluded[id] ? faTimes : faCheck} />
          </Button>
          <strong>
            <FormattedMessage id={id} />:
          </strong>
        </div>
        <TruncatedText className="col-7" text={text} />
      </div>
    );
  }

  function maybeRenderStashBoxLink() {
    const base = endpoint?.match(/https?:\/\/.*?\//)?.[0];
    const link = base ? `${base}tags/${tag.remote_site_id}` : undefined;

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

  function handleSave() {
    if (!tag.name) {
      throw new Error("tag name must be set");
    }

    const tagData: GQL.TagCreateInput = {
      name: tag.name,
      description: tag.description ?? undefined,
      aliases: tag.alias_list?.filter((a) => a) ?? undefined,
    };

    // stashid handling code
    const remoteSiteID = tag.remote_site_id;
    if (remoteSiteID && endpoint) {
      tagData.stash_ids = [
        {
          endpoint,
          stash_id: remoteSiteID,
          updated_at: new Date().toISOString(),
        },
      ];
    }

    // handle exclusions
    excludeFields(tagData, excluded);

    onSave(tagData);
  }

  return (
    <ModalComponent
      show={modalVisible}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: handleSave,
      }}
      cancel={{ onClick: () => closeModal(), variant: "secondary" }}
      onHide={() => closeModal()}
      dialogClassName="studio-create-modal"
      icon={icon}
      header={header}
    >
      <div>
        <div className="row">
          <div className="col-12">
            {maybeRenderField("name", tag.name)}
            {maybeRenderField("description", tag.description)}
            {maybeRenderField("aliases", tag.alias_list?.join(", "))}
            {maybeRenderStashBoxLink()}
          </div>
        </div>
      </div>
    </ModalComponent>
  );
};

export default TagModal;
