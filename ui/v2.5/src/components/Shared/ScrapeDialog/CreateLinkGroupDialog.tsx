import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { faLink } from "@fortawesome/free-solid-svg-icons";
import { Form } from "react-bootstrap";
import { Group, GroupSelect } from "../../Groups/GroupSelect";

export const CreateLinkGroupDialog: React.FC<{
  group: GQL.ScrapedGroup;
  onClose: (result: {
    create?: GQL.GroupCreateInput;
    update?: GQL.GroupUpdateInput;
  }) => void;
  endpoint?: string;
}> = ({ group, onClose }) => {
  const intl = useIntl();

  const [createNew, setCreateNew] = useState(false);
  const [name, setName] = useState(group.name);
  const [existingGroup, setExistingGroup] = useState<Group | null>(null);

  function handleGroupSave() {
    if (createNew) {
      const createInput: GQL.GroupCreateInput = {
        name: name ?? "",
      };
      onClose({ create: createInput });
    } else if (existingGroup) {
      const updateInput: GQL.GroupUpdateInput = {
        id: existingGroup.id,
      };
      onClose({ update: updateInput });
    }
  }

  return (
    <ModalComponent
      show={true}
      accept={{
        text: intl.formatMessage({ id: "actions.save" }),
        onClick: () => handleGroupSave(),
      }}
      disabled={
        createNew ? (name?.trim() ?? "") === "" : existingGroup === null
      }
      cancel={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: () => {
          onClose({});
        },
      }}
      dialogClassName="create-link-group-modal"
      icon={faLink}
      header={intl.formatMessage({ id: "component_tagger.verb_match_group" })}
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
            value={name ?? ""}
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
          <GroupSelect
            isMulti={false}
            values={existingGroup ? [existingGroup] : []}
            onSelect={(g) => setExistingGroup(g.length > 0 ? g[0] : null)}
            isDisabled={createNew}
            menuPortalTarget={document.body}
          />
        </Form.Group>
      </Form>
    </ModalComponent>
  );
};
