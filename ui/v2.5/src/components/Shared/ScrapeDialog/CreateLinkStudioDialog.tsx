import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { faLink } from "@fortawesome/free-solid-svg-icons";
import { Form } from "react-bootstrap";
import { Studio, StudioSelect } from "../../Studios/StudioSelect";

export const CreateLinkStudioDialog: React.FC<{
  studio: GQL.ScrapedStudio;
  onClose: (result: {
    create?: GQL.StudioCreateInput;
    update?: GQL.StudioUpdateInput;
  }) => void;
  endpoint?: string;
}> = ({ studio, onClose, endpoint }) => {
  const intl = useIntl();

  const [createNew, setCreateNew] = useState(false);
  const [name, setName] = useState(studio.name);
  const [existingStudio, setExistingStudio] = useState<Studio | null>(null);

  function handleStudioSave() {
    if (createNew) {
      const createInput: GQL.StudioCreateInput = {
        name: name ?? "",
        stash_ids:
          endpoint && studio.remote_site_id
            ? [{ endpoint: endpoint!, stash_id: studio.remote_site_id }]
            : undefined,
      };
      onClose({ create: createInput });
    } else if (existingStudio) {
      const updateInput: GQL.StudioUpdateInput = {
        id: existingStudio.id,
        stash_ids:
          endpoint && studio.remote_site_id
            ? [{ endpoint: endpoint!, stash_id: studio.remote_site_id }]
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
        onClick: () => handleStudioSave(),
      }}
      disabled={createNew ? (name?.trim() ?? "") === "" : existingStudio === null}
      cancel={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: () => {
          onClose({});
        },
      }}
      dialogClassName="create-link-studio-modal"
      icon={faLink}
      header={intl.formatMessage({ id: "component_tagger.verb_match_studio" })}
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
          <StudioSelect
            isMulti={false}
            values={existingStudio ? [existingStudio] : []}
            onSelect={(s) => setExistingStudio(s.length > 0 ? s[0] : null)}
            isDisabled={createNew}
            menuPortalTarget={document.body}
          />
        </Form.Group>
      </Form>
    </ModalComponent>
  );
};
