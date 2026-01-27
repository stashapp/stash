import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapedInputGroupRow,
  ScrapedTextAreaRow,
  ScrapedImageRow,
  ScrapedStringListRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialogRow";
import { ScrapeDialog } from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import { useIntl } from "react-intl";
import { uniq } from "lodash-es";
import { Performer } from "src/components/Performers/PerformerSelect";
import { sortStoredIdObjects } from "src/utils/data";
import {
  ObjectListScrapeResult,
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import {
  ScrapedGroupsRow,
  ScrapedPerformersRow,
  ScrapedStudioRow,
} from "src/components/Shared/ScrapeDialog/ScrapedObjectsRow";
import {
  useCreateScrapedGroup,
  useCreateScrapedPerformer,
  useCreateScrapedStudio,
} from "src/components/Shared/ScrapeDialog/createObjects";
import { Tag } from "src/components/Tags/TagSelect";
import { Studio } from "src/components/Studios/StudioSelect";
import { Group } from "src/components/Groups/GroupSelect";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";

interface ISceneScrapeDialogProps {
  scene: Partial<GQL.SceneUpdateInput>;
  sceneStudio: Studio | null;
  scenePerformers: Performer[];
  sceneTags: Tag[];
  sceneGroups: Group[];
  scraped: GQL.ScrapedScene;
  endpoint?: string;

  onClose: (scrapedScene?: GQL.ScrapedScene) => void;
}

export const SceneScrapeDialog: React.FC<ISceneScrapeDialogProps> = ({
  scene,
  sceneStudio,
  scenePerformers,
  sceneTags,
  sceneGroups,
  scraped,
  onClose,
  endpoint,
}) => {
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.title, scraped.title)
  );
  const [code, setCode] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.code, scraped.code)
  );

  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      scene.urls,
      scraped.urls
        ? uniq((scene.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );

  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.date, scraped.date)
  );
  const [director, setDirector] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.director, scraped.director)
  );
  const [studio, setStudio] = useState<ObjectScrapeResult<GQL.ScrapedStudio>>(
    new ObjectScrapeResult<GQL.ScrapedStudio>(
      sceneStudio
        ? {
            stored_id: sceneStudio.id,
            name: sceneStudio.name,
          }
        : undefined,
      scraped.studio?.stored_id ? scraped.studio : undefined
    )
  );
  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
    scraped.studio && !scraped.studio.stored_id ? scraped.studio : undefined
  );

  const [stashID, setStashID] = useState(
    new ScrapeResult<string>(
      scene.stash_ids?.find((s) => s.endpoint === endpoint)?.stash_id,
      scraped.remote_site_id
    )
  );

  const [performers, setPerformers] = useState<
    ObjectListScrapeResult<GQL.ScrapedPerformer>
  >(
    new ObjectListScrapeResult<GQL.ScrapedPerformer>(
      sortStoredIdObjects(
        scenePerformers.map((p) => ({
          stored_id: p.id,
          name: p.name,
        }))
      ),
      sortStoredIdObjects(scraped.performers ?? undefined)
    )
  );
  const [newPerformers, setNewPerformers] = useState<GQL.ScrapedPerformer[]>(
    scraped.performers?.filter((t) => !t.stored_id) ?? []
  );

  const [groups, setGroups] = useState<
    ObjectListScrapeResult<GQL.ScrapedGroup>
  >(
    new ObjectListScrapeResult<GQL.ScrapedGroup>(
      sortStoredIdObjects(
        sceneGroups.map((p) => ({
          stored_id: p.id,
          name: p.name,
        }))
      ),
      sortStoredIdObjects(scraped.groups ?? undefined)
    )
  );
  const [newGroups, setNewGroups] = useState<GQL.ScrapedGroup[]>(
    scraped.groups?.filter((t) => !t.stored_id) ?? []
  );

  const { tags, newTags, scrapedTagsRow, linkDialog } = useScrapedTags(
    sceneTags,
    scraped.tags,
    endpoint
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.details, scraped.details)
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.cover_image, scraped.image)
  );

  const createNewStudio = useCreateScrapedStudio({
    scrapeResult: studio,
    setScrapeResult: setStudio,
    setNewObject: setNewStudio,
    endpoint,
  });

  const createNewPerformer = useCreateScrapedPerformer({
    scrapeResult: performers,
    setScrapeResult: setPerformers,
    newObjects: newPerformers,
    setNewObjects: setNewPerformers,
    endpoint,
  });

  const createNewGroup = useCreateScrapedGroup({
    scrapeResult: groups,
    setScrapeResult: setGroups,
    newObjects: newGroups,
    setNewObjects: setNewGroups,
    endpoint,
  });

  const intl = useIntl();

  // don't show the dialog if nothing was scraped
  if (
    [
      title,
      code,
      urls,
      date,
      director,
      studio,
      performers,
      groups,
      tags,
      details,
      image,
      stashID,
    ].every((r) => !r.scraped) &&
    newTags.length === 0 &&
    newPerformers.length === 0 &&
    newGroups.length === 0 &&
    !newStudio
  ) {
    onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedSceneDataFragment {
    const newStudioValue = studio.getNewValue();

    return {
      title: title.getNewValue(),
      code: code.getNewValue(),
      urls: urls.getNewValue(),
      date: date.getNewValue(),
      director: director.getNewValue(),
      studio: newStudioValue,
      performers: performers.getNewValue(),
      groups: groups.getNewValue(),
      tags: tags.getNewValue(),
      details: details.getNewValue(),
      image: image.getNewValue(),
      remote_site_id: stashID.getNewValue(),
    };
  }

  function renderScrapeRows() {
    return (
      <>
        <ScrapedInputGroupRow
          field="title"
          title={intl.formatMessage({ id: "title" })}
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          field="code"
          title={intl.formatMessage({ id: "scene_code" })}
          result={code}
          onChange={(value) => setCode(value)}
        />
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        <ScrapedInputGroupRow
          field="date"
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapedInputGroupRow
          field="director"
          title={intl.formatMessage({ id: "director" })}
          result={director}
          onChange={(value) => setDirector(value)}
        />
        <ScrapedStudioRow
          field="studio"
          title={intl.formatMessage({ id: "studios" })}
          result={studio}
          onChange={(value) => setStudio(value)}
          newStudio={newStudio}
          onCreateNew={createNewStudio}
        />
        <ScrapedPerformersRow
          field="performers"
          title={intl.formatMessage({ id: "performers" })}
          result={performers}
          onChange={(value) => setPerformers(value)}
          newObjects={newPerformers}
          onCreateNew={createNewPerformer}
          ageFromDate={date.useNewValue ? date.newValue : date.originalValue}
        />
        <ScrapedGroupsRow
          field="groups"
          title={intl.formatMessage({ id: "groups" })}
          result={groups}
          onChange={(value) => setGroups(value)}
          newObjects={newGroups}
          onCreateNew={createNewGroup}
        />
        {scrapedTagsRow}
        <ScrapedTextAreaRow
          field="details"
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapedInputGroupRow
          field="stash_ids"
          title={intl.formatMessage({ id: "stash_id" })}
          result={stashID}
          locked
          onChange={(value) => setStashID(value)}
        />
        <ScrapedImageRow
          field="cover_image"
          title={intl.formatMessage({ id: "cover_image" })}
          className="scene-cover"
          result={image}
          onChange={(value) => setImage(value)}
        />
      </>
    );
  }

  if (linkDialog) {
    return linkDialog;
  }

  return (
    <ScrapeDialog
      title={intl.formatMessage(
        { id: "dialogs.scrape_entity_title" },
        { entity_type: intl.formatMessage({ id: "scene" }) }
      )}
      onClose={(apply) => {
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};

export default SceneScrapeDialog;
