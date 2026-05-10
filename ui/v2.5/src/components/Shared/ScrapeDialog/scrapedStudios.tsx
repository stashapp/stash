import { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ObjectScrapeResult } from "./scrapeResult";
import { Studio } from "src/components/Studios/StudioSelect";
import { useCreateScrapedStudio, useLinkScrapedStudio } from "./createObjects";
import { ScrapedStudioRow } from "./ScrapedObjectsRow";
import { CreateLinkStudioDialog } from "./CreateLinkStudioDialog";
import { useStudioCreate, useStudioUpdate } from "src/core/StashService";
import { toastOperation, useToast } from "src/hooks/Toast";

export function useScrapedStudios(
  existingStudio: Studio | null,
  scrapedStudio?: GQL.Maybe<GQL.ScrapedStudio>,
  endpoint?: string
) {
  const intl = useIntl();
  const Toast = useToast();

  const [studio, setStudio] = useState<ObjectScrapeResult<GQL.ScrapedStudio>>(
    new ObjectScrapeResult<GQL.ScrapedStudio>(
      existingStudio
        ? {
            stored_id: existingStudio.id,
            name: existingStudio.name,
          }
        : undefined,
      scrapedStudio?.stored_id ? scrapedStudio : undefined
    )
  );

  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
    scrapedStudio && !scrapedStudio.stored_id ? scrapedStudio : undefined
  );
  const [linkedStudio, setLinkedStudio] = useState<GQL.ScrapedStudio | null>(
    null
  );

  const createNewStudio = useCreateScrapedStudio({
    scrapeResult: studio,
    setScrapeResult: setStudio,
    setNewObject: setNewStudio,
    endpoint,
  });

  const [createStudio] = useStudioCreate();
  const [updateStudio] = useStudioUpdate();

  const linkScrapedStudio = useLinkScrapedStudio({
    scrapeResult: studio,
    setScrapeResult: setStudio,
    newObject: newStudio,
    setNewObject: setNewStudio,
  });

  async function handleLinkStudioResult(studioResult: {
    create?: GQL.StudioCreateInput;
    update?: GQL.StudioUpdateInput;
  }) {
    if (studioResult.create) {
      await toastOperation(
        Toast,
        async () => {
          const result = await createStudio({
            variables: { input: studioResult.create! },
          });

          if (result.data?.studioCreate) {
            linkScrapedStudio(
              result.data.studioCreate.id,
              result.data.studioCreate.name
            );
          }
        },
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "studio" }).toLocaleLowerCase(),
          }
        )
      )();
    } else if (studioResult.update) {
      await toastOperation(
        Toast,
        async () => {
          const result = await updateStudio({
            variables: { input: studioResult.update! },
          });

          if (result.data?.studioUpdate) {
            linkScrapedStudio(
              result.data.studioUpdate.id,
              result.data.studioUpdate.name
            );
          }
        },
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "studio" }).toLocaleLowerCase(),
          }
        )
      )();
    }

    setLinkedStudio(null);
  }

  const linkDialog = linkedStudio ? (
    <CreateLinkStudioDialog
      studio={linkedStudio}
      onClose={handleLinkStudioResult}
      endpoint={endpoint}
    />
  ) : null;

  const scrapedStudioRow = (
    <ScrapedStudioRow
      field="studio"
      title={intl.formatMessage({ id: "studios" })}
      result={studio}
      onChange={(value) => setStudio(value)}
      newStudio={newStudio}
      onCreateNew={createNewStudio}
      onLinkExisting={(s) => setLinkedStudio(s)}
    />
  );

  return {
    studio,
    newStudio,
    linkDialog,
    scrapedStudioRow,
  };
}
