import React, { useCallback, useState } from "react";
import cx from "classnames";
import { Button } from "react-bootstrap";
import { uniq } from "lodash";
import { blobToBase64 } from "base64-blob";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, SuccessIcon } from "src/components/Shared";
import PerformerResult, { PerformerOperation } from "./PerformerResult";
import StudioResult, { StudioOperation } from "./StudioResult";
import { IStashBoxScene } from "./utils";
import {
  useCreateTag,
  useCreatePerformer,
  useCreateStudio,
  useUpdatePerformerStashID,
  useUpdateStudioStashID,
} from "./queries";

const getDurationStatus = (
  scene: IStashBoxScene,
  stashDuration: number | undefined | null
) => {
  const fingerprintDuration =
    scene.fingerprints.map((f) => f.duration)?.[0] ?? null;
  const sceneDuration = scene.duration || fingerprintDuration;
  if (!sceneDuration || !stashDuration) return "";
  const diff = Math.abs(sceneDuration - stashDuration);
  if (diff < 5) {
    return (
      <div className="font-weight-bold">
        <SuccessIcon className="mr-2" />
        Duration is a match
      </div>
    );
  }
  return <div>Duration off by {Math.floor(diff)}s</div>;
};

const getFingerprintStatus = (
  scene: IStashBoxScene,
  stashChecksum?: string
) => {
  if (scene.fingerprints.some((f) => f.hash === stashChecksum))
    return (
      <div className="font-weight-bold">
        <SuccessIcon className="mr-2" />
        Checksum is a match
      </div>
    );
};

interface IStashSearchResultProps {
  scene: IStashBoxScene;
  stashScene: GQL.SlimSceneDataFragment;
  isActive: boolean;
  setActive: () => void;
  showMales: boolean;
  setScene: (scene: GQL.SlimSceneDataFragment) => void;
  setCoverImage: boolean;
  tagOperation: string;
  endpoint: string;
  queueFingerprintSubmission: (sceneId: string, endpoint: string) => void;
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  scene,
  stashScene,
  isActive,
  setActive,
  showMales,
  setScene,
  setCoverImage,
  tagOperation,
  endpoint,
  queueFingerprintSubmission
}) => {
  const [studio, setStudio] = useState<StudioOperation>();
  const [performers, setPerformers] = useState<
    Record<string, PerformerOperation>
  >({});
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const createStudio = useCreateStudio();
  const createPerformer = useCreatePerformer();
  const createTag = useCreateTag();
  const updatePerformerStashID = useUpdatePerformerStashID();
  const updateStudioStashID = useUpdateStudioStashID();
  const [updateScene] = GQL.useSceneUpdateMutation({
    onError: (errors) => errors,
  });
  const { data: allTags } = GQL.useAllTagsForFilterQuery();

  const setPerformer = useCallback(
    (performerData: PerformerOperation, performerID: string) =>
      setPerformers({ ...performers, [performerID]: performerData }),
    [performers]
  );

  const handleSave = async () => {
    setError({});
    let performerIDs = [];
    let studioID = null

    if (!studio)
      return;

    if (studio.type === 'create') {
      setSaveState("Creating studio");
      const newStudio = {
        name: studio.data.name,
        stash_ids: [
          {
            endpoint,
            stash_id: scene.studio.stash_id,
          },
        ],
        url: studio.data.url,
      };
      const studioCreateResult = await createStudio(
        newStudio,
        scene.studio.stash_id
      );

      if (!studioCreateResult?.data?.studioCreate) {
        setError({
          message: `Failed to save studio "${newStudio.name}"`,
          details: studioCreateResult?.errors?.[0].message,
        });
        return setSaveState("");
      }
      studioID = studioCreateResult.data.studioCreate.id;
    }
    else if (studio.type === 'update') {
      setSaveState("Saving studio stashID");
      const res = await updateStudioStashID(studio.data.id, [
        ...studio.data.stash_ids,
        { stash_id: scene.studio.stash_id, endpoint },
      ]);
      if (!res?.data?.studioUpdate) {
        setError({
          message: `Failed to save stashID to studio "${studio.data.name}"`,
          details: res?.errors?.[0].message,
        });
        return setSaveState("");
      }
      studioID = res.data.studioUpdate.id;
    }
    else if (studio.type === 'existing') {
      studioID = studio.data.id;
    }

    setSaveState("Saving performers");
    performerIDs = await Promise.all(
      Object.keys(performers).map(async (stashID) => {
        const performer = performers[stashID];
        if (performer.type ===  'skip') return "Skip";

        let performerID = performer.data.id;

        if (performer.type === 'create') {
          const imgurl = performer.data.images[0];
          let imgData = null;
          if (imgurl) {
            const img = await fetch(imgurl, {
              mode: "cors",
              cache: "no-store",
            });
            if (img.status === 200) {
              const blob = await img.blob();
              imgData = await blobToBase64(blob);
            }
          }

          const performerInput = {
            name: performer.data.name,
            gender: performer.data.gender,
            country: performer.data.country,
            height: performer.data.height,
            ethnicity: performer.data.ethnicity,
            birthdate: performer.data.birthdate,
            eye_color: performer.data.eye_color,
            fake_tits: performer.data.fake_tits,
            measurements: performer.data.measurements,
            career_length: performer.data.career_length,
            tattoos: performer.data.tattoos,
            piercings: performer.data.piercings,
            twitter: performer.data.twitter,
            instagram: performer.data.instagram,
            image: imgData,
            stash_ids: [
              {
                endpoint,
                stash_id: stashID,
              },
            ],
          };

          const res = await createPerformer(performerInput, stashID);
          if (!res?.data?.performerCreate) {
            setError({
              message: `Failed to save performer "${performerInput.name}"`,
              details: res?.errors?.[0].message,
            });
            return null;
          }
          performerID = res.data?.performerCreate.id;
        }

        if (performer.type === 'update') {
          const stashIDs = performer.data.stash_ids;
          await updatePerformerStashID(performer.data.id, [
            ...stashIDs,
            { stash_id: stashID, endpoint },
          ]);
        }

        return performerID;
      })
    );

    if (!performerIDs.some((id) => !id)) {
      setSaveState("Updating scene");
      const imgurl = scene.images[0];
      let imgData = null;
      if (imgurl && setCoverImage) {
        const img = await fetch(imgurl, {
          mode: "cors",
          cache: "no-store",
        });
        if (img.status === 200) {
          const blob = await img.blob();
          imgData = await blobToBase64(blob);
        }
      }

      const tagIDs: string[] =
        tagOperation === "merge"
          ? stashScene?.tags?.map((t) => t.id) ?? []
          : [];
      const tags = scene.tags ?? [];
      if (tags.length > 0) {
        const tagDict: Record<string, string> = (allTags?.allTagsSlim ?? [])
          .filter((t) => t.name)
          .reduce((dict, t) => ({ ...dict, [t.name.toLowerCase()]: t.id }), {});
        const newTags: string[] = [];
        tags.forEach((tag) => {
          if (tagDict[tag.name.toLowerCase()])
            tagIDs.push(tagDict[tag.name.toLowerCase()]);
          else newTags.push(tag.name);
        });

        const createdTags = await Promise.all(
          newTags.map((tag) => createTag(tag))
        );
        createdTags.forEach((createdTag) => {
          if (createdTag?.data?.tagCreate?.id)
            tagIDs.push(createdTag.data.tagCreate.id);
        });
      }

      const sceneUpdateResult = await updateScene({
        variables: {
          id: stashScene.id ?? "",
          title: scene.title,
          details: scene.details,
          date: scene.date,
          performer_ids: performerIDs.filter((id) => id !== "Skip") as string[],
          studio_id: studioID,
          cover_image: imgData,
          url: scene.url,
          ...(tagIDs ? { tag_ids: uniq(tagIDs) } : {}),
          stash_ids: [
            ...(stashScene?.stash_ids ?? []),
            {
              endpoint,
              stash_id: scene.stash_id,
            },
          ],
        },
      });

      if (!sceneUpdateResult?.data?.sceneUpdate) {
        setError({
          message: "Failed to save scene",
          details: sceneUpdateResult?.errors?.[0].message,
        });
      } else if (sceneUpdateResult.data?.sceneUpdate)
        setScene(sceneUpdateResult.data.sceneUpdate);

      queueFingerprintSubmission(stashScene.id, endpoint);
    }

    setSaveState("");
  };

  const classname = cx("row no-gutters mt-2 search-result", {
    "selected-result": isActive,
  });

  const sceneTitle = scene.url ? (
    <a
      href={scene.url}
      target="_blank"
      rel="noopener noreferrer"
      className="scene-link"
    >
      {scene?.title}
    </a>
  ) : (
    <span>{scene?.title}</span>
  );

  const saveEnabled =
    Object.keys(performers ?? []).length ===
      scene.performers.filter((p) => p.gender !== "MALE" || showMales).length &&
    Object.keys(performers ?? []).every(
      (id) => performers?.[id].type
    ) &&
    saveState === "";

  return (
    // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions
    <li
      className={classname}
      key={scene.stash_id}
      onClick={() => !isActive && setActive()}
    >
      <div className="col-6">
        <div className="row">
          <img
            src={scene.images[0]}
            alt=""
            className="align-self-center scene-image"
          />
          <div className="d-flex flex-column justify-content-center scene-metadata">
            <h4 className="text-truncate" title={scene?.title ?? ""}>
              {sceneTitle}
            </h4>
            <h5>
              {scene?.studio?.name} • {scene?.date}
            </h5>
            <div>
              Performers: {scene?.performers?.map((p) => p.name).join(", ")}
            </div>
            {getDurationStatus(scene, stashScene.file?.duration)}
            {getFingerprintStatus(
              scene,
              stashScene.checksum ?? stashScene.oshash ?? undefined
            )}
          </div>
        </div>
      </div>
      {isActive && (
        <div className="col-6">
          <StudioResult studio={scene.studio} setStudio={setStudio} />
          {scene.performers
            .filter((p) => p.gender !== "MALE" || showMales)
            .map((performer) => (
              <PerformerResult
                performer={performer}
                setPerformer={(data: PerformerOperation) =>
                  setPerformer(data, performer.stash_id)
                }
                key={`${scene.stash_id}${performer.stash_id}`}
              />
            ))}
          <div className="row no-gutters mt-2 align-items-center justify-content-end">
            {error.message && (
              <strong className="mt-1 mr-2 text-danger text-right">
                <abbr title={error.details} className="mr-2">
                  Error:
                </abbr>
                {error.message}
              </strong>
            )}
            {saveState && (
              <strong className="col-4 mt-1 mr-2 text-right">
                {saveState}
              </strong>
            )}
            <Button onClick={handleSave} disabled={!saveEnabled}>
              {saveState ? (
                <LoadingIndicator inline small message="" />
              ) : (
                "Save"
              )}
            </Button>
          </div>
        </div>
      )}
    </li>
  );
};

export default StashSearchResult;
