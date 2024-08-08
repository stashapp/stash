import React, { useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useToast } from "src/hooks/Toast";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import {
  ContainingGroupTable,
  IRelatedGroupEntry,
} from "./ContainingGroupTable";
import { ModalComponent } from "src/components/Shared/Modal";
import { useAddSubGroups } from "src/core/StashService";

interface IListOperationProps {
  containingGroup: GQL.GroupDataFragment;
  onClose: (applied: boolean) => void;
}

export const AddSubGroupsDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const [isUpdating, setIsUpdating] = useState(false);

  const addSubGroups = useAddSubGroups();

  const Toast = useToast();

  const [entries, setEntries] = useState<IRelatedGroupEntry[]>([]);

  const excludeIDs = useMemo(
    () => [
      ...props.containingGroup.containing_groups.map((m) => m.group.id),
      props.containingGroup.id,
    ],
    [props.containingGroup]
  );

  const onSave = async () => {
    setIsUpdating(true);
    try {
      // add the sub groups
      await addSubGroups(
        props.containingGroup.id,
        entries.map((m) => ({
          group_id: m.group.id,
          description: m.description,
        }))
      );

      const imageCount = entries.length;
      Toast.success(
        intl.formatMessage(
          { id: "toast.added_entity" },
          {
            count: imageCount,
            singularEntity: intl.formatMessage({ id: "group" }),
            pluralEntity: intl.formatMessage({ id: "groups" }),
          }
        )
      );

      props.onClose(true);
    } catch (err) {
      Toast.error(err);
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <ModalComponent
      show
      icon={faPlus}
      header={intl.formatMessage({ id: "actions.add_sub_groups" })}
      accept={{
        onClick: onSave,
        text: intl.formatMessage({ id: "actions.add" }),
      }}
      cancel={{
        onClick: () => props.onClose(false),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      isRunning={isUpdating}
    >
      <Form>
        <ContainingGroupTable
          value={entries}
          onUpdate={(input) => setEntries(input)}
          excludeIDs={excludeIDs}
        />
      </Form>
    </ModalComponent>
  );
};
