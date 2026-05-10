import { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ObjectListScrapeResult } from "./scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import { Performer } from "src/components/Performers/PerformerSelect";
import {
  useCreateScrapedPerformer,
  useLinkScrapedPerformer,
} from "./createObjects";
import { ScrapedPerformersRow } from "./ScrapedObjectsRow";
import { CreateLinkPerformerDialog } from "./CreateLinkPerformerDialog";
import { usePerformerCreate, usePerformerUpdate } from "src/core/StashService";
import { toastOperation, useToast } from "src/hooks/Toast";

export function useScrapedPerformers(
  existingPerformers: Performer[],
  scrapedPerformers?: GQL.Maybe<GQL.ScrapedPerformer[]>,
  endpoint?: string,
  ageFromDate?: string | null
) {
  const intl = useIntl();
  const Toast = useToast();

  const [performers, setPerformers] = useState<
    ObjectListScrapeResult<GQL.ScrapedPerformer>
  >(
    new ObjectListScrapeResult<GQL.ScrapedPerformer>(
      sortStoredIdObjects(
        existingPerformers.map((p) => ({
          stored_id: p.id,
          name: p.name,
        }))
      ),
      sortStoredIdObjects(scrapedPerformers ?? undefined)
    )
  );

  const [newPerformers, setNewPerformers] = useState<GQL.ScrapedPerformer[]>(
    scrapedPerformers?.filter((p) => !p.stored_id) ?? []
  );
  const [linkedPerformerIndex, setLinkedPerformerIndex] = useState<
    number | null
  >(null);

  const createNewPerformer = useCreateScrapedPerformer({
    scrapeResult: performers,
    setScrapeResult: setPerformers,
    newObjects: newPerformers,
    setNewObjects: setNewPerformers,
    endpoint,
  });

  const [createPerformer] = usePerformerCreate();
  const [updatePerformer] = usePerformerUpdate();

  const linkScrapedPerformer = useLinkScrapedPerformer({
    scrapeResult: performers,
    setScrapeResult: setPerformers,
    newObjects: newPerformers,
    setNewObjects: setNewPerformers,
  });

  async function handleLinkPerformerResult(performer: {
    create?: GQL.PerformerCreateInput;
    update?: GQL.PerformerUpdateInput;
  }) {
    if (performer.create) {
      await toastOperation(
        Toast,
        async () => {
          const result = await createPerformer({
            variables: { input: performer.create! },
          });

          if (result.data?.performerCreate && linkedPerformerIndex !== null) {
            linkScrapedPerformer(
              result.data.performerCreate.id,
              result.data.performerCreate.name,
              linkedPerformerIndex
            );
          }
        },
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "performer" }).toLocaleLowerCase(),
          }
        )
      )();
    } else if (performer.update) {
      await toastOperation(
        Toast,
        async () => {
          const result = await updatePerformer({
            variables: { input: performer.update! },
          });

          if (result.data?.performerUpdate && linkedPerformerIndex !== null) {
            linkScrapedPerformer(
              result.data.performerUpdate.id,
              result.data.performerUpdate.name,
              linkedPerformerIndex
            );
          }
        },
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "performer" }).toLocaleLowerCase(),
          }
        )
      )();
    }

    setLinkedPerformerIndex(null);
  }

  const linkedPerformer =
    linkedPerformerIndex !== null ? newPerformers[linkedPerformerIndex] : null;

  const linkDialog = linkedPerformer ? (
    <CreateLinkPerformerDialog
      performer={linkedPerformer}
      onClose={handleLinkPerformerResult}
      endpoint={endpoint}
    />
  ) : null;

  const scrapedPerformersRow = (
    <ScrapedPerformersRow
      field="performers"
      title={intl.formatMessage({ id: "performers" })}
      result={performers}
      onChange={(value) => setPerformers(value)}
      newObjects={newPerformers}
      onCreateNew={createNewPerformer}
      onLinkExisting={(p, index) => setLinkedPerformerIndex(index)}
      ageFromDate={ageFromDate}
    />
  );

  return {
    performers,
    newPerformers,
    linkDialog,
    scrapedPerformersRow,
  };
}
