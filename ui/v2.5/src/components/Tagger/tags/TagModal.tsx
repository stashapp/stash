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
import { Button, Form } from "react-bootstrap";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { excludeFields } from "src/utils/data";
import { ExternalLink } from "src/components/Shared/ExternalLink";

interface ITagModalProps {
  tag: GQL.ScrapedSceneTagDataFragment;
  modalVisible: boolean;
  closeModal: () => void;
  onSave: (input: GQL.TagCreateInput, parentInput?: GQL.TagCreateInput) => void;
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

  const [createParentTag, setCreateParentTag] = useState<boolean>(
    !!tag.parent && !tag.parent.stored_id
  );

  // Check if a tag with the parent name already exists locally.
  // Categories don't have stash IDs, so stored_id may be null even when the
  // parent tag has already been created (e.g. by tagging a sibling tag first).
  const parentNameQuery = GQL.useFindTagsQuery({
    skip: !tag.parent || !!tag.parent.stored_id,
    variables: {
      tag_filter: {
        name: {
          value: tag.parent?.name ?? "",
          modifier: GQL.CriterionModifier.Equals,
        },
      },
      filter: { per_page: 1 },
    },
  });
  const existingParentId = parentNameQuery.data?.findTags.tags[0]?.id;

  // If the parent already exists locally, don't offer to create it
  const sendParentTag = !existingParentId;

  const [parentExcluded, setParentExcluded] = useState<Record<string, boolean>>(
    excludedTagFields.reduce((dict, field) => ({ ...dict, [field]: true }), {})
  );
  const toggleParentField = (name: string) =>
    setParentExcluded({
      ...parentExcluded,
      [name]: !parentExcluded[name],
    });

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
        <TruncatedText className="col-7" text={text} lineCount={3} />
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

  function maybeRenderParentField(
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
              onClick={() => toggleParentField(id)}
              variant="secondary"
              className={parentExcluded[id] ? "text-muted" : "text-success"}
            >
              <Icon icon={parentExcluded[id] ? faTimes : faCheck} />
            </Button>
          )}
          <strong>
            <FormattedMessage id={id} />:
          </strong>
        </div>
        <TruncatedText className="col-7" text={text} lineCount={3} />
      </div>
    );
  }

  function maybeRenderParentTagDetails() {
    if (!createParentTag || !tag.parent) {
      return;
    }

    return (
      <div>
        {maybeRenderParentField("name", tag.parent.name, false)}
        {maybeRenderParentField("description", tag.parent.description)}
      </div>
    );
  }

  function maybeRenderParentTag() {
    // No parent tag, or parent already exists locally
    if (!tag.parent || tag.parent.stored_id || !sendParentTag) {
      return;
    }

    return (
      <div>
        <div className="mb-4 mt-4">
          <Form.Check
            id="create-parent"
            checked={createParentTag}
            label={intl.formatMessage({
              id: "actions.create_parent_tag",
            })}
            onChange={() => setCreateParentTag(!createParentTag)}
          />
        </div>
        {maybeRenderParentTagDetails()}
      </div>
    );
  }

  function handleSave() {
    if (!tag.name) {
      throw new Error("tag name must be set");
    }

    const parentId = tag.parent?.stored_id ?? existingParentId;

    const tagData: GQL.TagCreateInput = {
      name: tag.name,
      description: tag.description ?? undefined,
      aliases: tag.alias_list?.filter((a) => a) ?? undefined,
      parent_ids: parentId ? [parentId] : undefined,
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

    let parentData: GQL.TagCreateInput | undefined = undefined;

    // Categories don't have stash IDs, so we only create new parent tags
    if (
      createParentTag &&
      sendParentTag &&
      tag.parent &&
      !tag.parent.stored_id
    ) {
      parentData = {
        name: tag.parent.name,
        description: tag.parent.description ?? undefined,
      };

      // handle exclusions
      // Can't exclude parent tag name when creating a new one
      parentExcluded.name = false;
      excludeFields(parentData, parentExcluded);
    }

    onSave(tagData, parentData);
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
            {maybeRenderField("parent_tags", tag.parent?.name, false)}
            {maybeRenderStashBoxLink()}
          </div>
        </div>
      </div>
      {maybeRenderParentTag()}
    </ModalComponent>
  );
};

export default TagModal;
