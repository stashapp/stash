import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { faLink } from "@fortawesome/free-solid-svg-icons";
import { Form } from "react-bootstrap";
import { Tag, TagSelect } from "../../Tags/TagSelect";

export const CreateLinkTagDialog: React.FC<{
  tag: GQL.ScrapedTag;
  onClose: (result: {
    create?: GQL.TagCreateInput;
    update?: GQL.TagUpdateInput;
  }) => void;
  endpoint?: string;
}> = ({ tag, onClose, endpoint }) => {
  const intl = useIntl();

  const [createNew, setCreateNew] = useState(false);
  const [name, setName] = useState(tag.name);
  const [existingTag, setExistingTag] = useState<Tag | null>(null);
  const [addAsAlias, setAddAsAlias] = useState(false);

  const canAddAlias = (createNew && name !== tag.name) || !createNew;

  useEffect(() => {
    setAddAsAlias(canAddAlias);
  }, [canAddAlias]);

  function handleTagSave() {
    if (createNew) {
      const createInput: GQL.TagCreateInput = {
        name: name,
        aliases: addAsAlias ? [tag.name] : [],
        stash_ids:
          endpoint && tag.remote_site_id
            ? [{ endpoint: endpoint!, stash_id: tag.remote_site_id }]
            : undefined,
      };
      onClose({ create: createInput });
    } else if (existingTag) {
      const updateInput: GQL.TagUpdateInput = {
        id: existingTag.id,
        aliases: addAsAlias
          ? [...(existingTag.aliases || []), tag.name]
          : undefined,
        // add stash id if applicable
        stash_ids:
          endpoint && tag.remote_site_id
            ? [
                ...(existingTag.stash_ids || []),
                { endpoint: endpoint!, stash_id: tag.remote_site_id },
              ]
            : undefined,
      };
      onClose({ update: updateInput });
    }
  }

  return (
    <ModalComponent
      show={true}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: () => handleTagSave(),
      }}
      disabled={createNew ? name.trim() === "" : existingTag === null}
      cancel={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: () => {
          onClose({});
        },
      }}
      dialogClassName="create-link-tag-modal"
      icon={faLink}
      header={intl.formatMessage({ id: "component_tagger.verb_match_tag" })}
    >
      <Form>
        <Form.Check
          type="radio"
          id="create-new"
          label={intl.formatMessage({ id: "actions.create_new" })}
          checked={createNew}
          onChange={() => setCreateNew(true)}
        />

        <Form.Group className="ml-3 mt-2">
          <Form.Label>
            <FormattedMessage id="name" />
          </Form.Label>
          <Form.Control
            className="input-control"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={!createNew}
          />
        </Form.Group>

        <Form.Check
          type="radio"
          id="link-existing"
          label={intl.formatMessage({
            id: "component_tagger.verb_link_existing",
          })}
          checked={!createNew}
          onChange={() => setCreateNew(false)}
        />

        <Form.Group className="ml-3 mt-2">
          <TagSelect
            isMulti={false}
            values={existingTag ? [existingTag] : []}
            onSelect={(t) => setExistingTag(t.length > 0 ? t[0] : null)}
            isDisabled={createNew}
            menuPortalTarget={document.body}
          />
        </Form.Group>

        <Form.Group className="mt-3">
          <Form.Check
            type="checkbox"
            id="add-as-alias"
            label={intl.formatMessage({
              id: "component_tagger.verb_add_as_alias",
            })}
            checked={addAsAlias}
            onChange={() => setAddAsAlias(!addAsAlias)}
            disabled={!canAddAlias}
          />
        </Form.Group>
      </Form>
    </ModalComponent>
  );
};
