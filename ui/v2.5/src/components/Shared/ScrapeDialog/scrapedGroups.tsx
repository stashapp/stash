import { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ObjectListScrapeResult } from "./scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import { Group } from "src/components/Groups/GroupSelect";
import {
  useCreateScrapedGroup,
  useLinkScrapedGroup,
} from "./createObjects";
import { ScrapedGroupsRow } from "./ScrapedObjectsRow";
import { CreateLinkGroupDialog } from "./CreateLinkGroupDialog";
import { useGroupCreate, useGroupUpdate } from "src/core/StashService";
import { toastOperation, useToast } from "src/hooks/Toast";

export function useScrapedGroups(
  existingGroups: Group[],
  scrapedGroups?: GQL.Maybe<GQL.ScrapedGroup[]>,
  endpoint?: string
) {
  const intl = useIntl();
  const Toast = useToast();

  const [groups, setGroups] = useState<ObjectListScrapeResult<GQL.ScrapedGroup>>(
    new ObjectListScrapeResult<GQL.ScrapedGroup>(
      sortStoredIdObjects(
        existingGroups.map((g) => ({
          stored_id: g.id,
          name: g.name,
        }))
      ),
      sortStoredIdObjects(scrapedGroups ?? undefined)
    )
  );

  const [newGroups, setNewGroups] = useState<GQL.ScrapedGroup[]>(
    scrapedGroups?.filter((g) => !g.stored_id) ?? []
  );
  const [linkedGroupIndex, setLinkedGroupIndex] = useState<number | null>(null);

  const createNewGroup = useCreateScrapedGroup({
    scrapeResult: groups,
    setScrapeResult: setGroups,
    newObjects: newGroups,
    setNewObjects: setNewGroups,
  });

  const [createGroup] = useGroupCreate();
  const [updateGroup] = useGroupUpdate();

  const linkScrapedGroup = useLinkScrapedGroup({
    scrapeResult: groups,
    setScrapeResult: setGroups,
    newObjects: newGroups,
    setNewObjects: setNewGroups,
  });

  async function handleLinkGroupResult(group: {
    create?: GQL.GroupCreateInput;
    update?: GQL.GroupUpdateInput;
  }) {
    if (group.create) {
      await toastOperation(
        Toast,
        async () => {
          const result = await createGroup({
            variables: { input: group.create! },
          });

          if (result.data?.groupCreate && linkedGroupIndex !== null) {
            linkScrapedGroup(
              result.data.groupCreate.id,
              result.data.groupCreate.name,
              linkedGroupIndex
            );
          }
        },
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "group" }).toLocaleLowerCase(),
          }
        )
      )();
    } else if (group.update) {
      await toastOperation(
        Toast,
        async () => {
          const result = await updateGroup({
            variables: { input: group.update! },
          });

          if (result.data?.groupUpdate && linkedGroupIndex !== null) {
            linkScrapedGroup(
              result.data.groupUpdate.id,
              result.data.groupUpdate.name,
              linkedGroupIndex
            );
          }
        },
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "group" }).toLocaleLowerCase(),
          }
        )
      )();
    }

    setLinkedGroupIndex(null);
  }

  const linkedGroup = linkedGroupIndex !== null ? newGroups[linkedGroupIndex] : null;

  const linkDialog = linkedGroup ? (
    <CreateLinkGroupDialog
      group={linkedGroup}
      onClose={handleLinkGroupResult}
      endpoint={endpoint}
    />
  ) : null;

  const scrapedGroupsRow = (
    <ScrapedGroupsRow
      field="groups"
      title={intl.formatMessage({ id: "groups" })}
      result={groups}
      onChange={(value) => setGroups(value)}
      newObjects={newGroups}
      onCreateNew={createNewGroup}
      onLinkExisting={(g, index) => setLinkedGroupIndex(index)}
    />
  );

  return {
    groups,
    newGroups,
    linkDialog,
    scrapedGroupsRow,
  };
}
