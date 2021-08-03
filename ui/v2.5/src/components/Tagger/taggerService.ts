import * as GQL from "src/core/generated-graphql";
import { blobToBase64 } from "base64-blob";
import {
  useCreatePerformer,
  useCreateStudio,
  useUpdatePerformerStashID,
  useUpdateStudioStashID,
} from "./queries";
import { IPerformerOperations } from "./PerformerResult";
import { StudioOperation } from "./StudioResult";
import { IStashBoxScene } from "./utils";

export interface ITagSceneOptions {
  setCoverImage?: boolean;
  setTags?: boolean;
  tagOperation: string;
}

export function useTagScene(
  options: ITagSceneOptions,
  setSaveState: (state: string) => void,
  setError: (err: { message?: string; details?: string }) => void
) {
  const createStudio = useCreateStudio();
  const createPerformer = useCreatePerformer();
  const updatePerformerStashID = useUpdatePerformerStashID();
  const updateStudioStashID = useUpdateStudioStashID();
  const [updateScene] = GQL.useSceneUpdateMutation({
    onError: (e) => {
      const message =
        e.message === "invalid JPEG format: short Huffman data"
          ? "Failed to save scene due to corrupted cover image"
          : "Failed to save scene";
      setError({
        message,
        details: e.message,
      });
    },
  });

  const handleSave = async (
    stashScene: GQL.SlimSceneDataFragment,
    scene: IStashBoxScene,
    studio: StudioOperation | undefined,
    performers: IPerformerOperations,
    tagIDs: string[],
    excludedFields: string[],
    endpoint: string
  ) => {
    function resolveField<T>(field: string, stashField: T, remoteField: T) {
      if (excludedFields.includes(field)) {
        return stashField;
      }

      return remoteField;
    }

    setError({});
    let performerIDs = [];
    let studioID = null;

    if (studio) {
      if (studio.type === "create") {
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
      } else if (studio.type === "update") {
        setSaveState("Saving studio stashID");
        const res = await updateStudioStashID(studio.data, [
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
      } else if (studio.type === "existing") {
        studioID = studio.data.id;
      } else if (studio.type === "skip") {
        studioID = stashScene.studio?.id;
      }
    }

    setSaveState("Saving performers");
    let failed = false;
    performerIDs = await Promise.all(
      Object.keys(performers).map(async (stashID) => {
        const performer = performers[stashID];
        if (performer.type === "skip") return "Skip";

        let performerID = performer.data.id;

        if (performer.type === "create") {
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
            details: performer.data.details,
            death_date: performer.data.death_date,
            hair_color: performer.data.hair_color,
            weight: Number(performer.data.weight),
          };

          const res = await createPerformer(performerInput, stashID);
          if (!res?.data?.performerCreate) {
            setError({
              message: `Failed to save performer "${performerInput.name}"`,
              details: res?.errors?.[0].message,
            });
            failed = true;
            return null;
          }
          performerID = res.data?.performerCreate.id;
        }

        if (performer.type === "update") {
          const stashIDs = performer.data.stash_ids;
          await updatePerformerStashID(performer.data.id, [
            ...stashIDs,
            { stash_id: stashID, endpoint },
          ]);
        }

        return performerID;
      })
    );

    if (failed) {
      return setSaveState("");
    }

    setSaveState("Updating scene");
    const imgurl = scene.images[0];
    let imgData;
    if (imgurl && options.setCoverImage) {
      const img = await fetch(imgurl, {
        mode: "cors",
        cache: "no-store",
      });
      if (img.status === 200) {
        const blob = await img.blob();
        // Sanity check on image size since bad images will fail
        if (blob.size > 10000) imgData = await blobToBase64(blob);
      }
    }

    const performer_ids = performerIDs.filter(
      (id) => id !== "Skip"
    ) as string[];

    const sceneUpdateResult = await updateScene({
      variables: {
        input: {
          id: stashScene.id ?? "",
          title: resolveField("title", stashScene.title, scene.title),
          details: resolveField("details", stashScene.details, scene.details),
          date: resolveField("date", stashScene.date, scene.date),
          performer_ids:
            performer_ids.length === 0
              ? stashScene.performers.map((p) => p.id)
              : performer_ids,
          studio_id: studioID,
          cover_image: resolveField("cover_image", undefined, imgData),
          url: resolveField("url", stashScene.url, scene.url),
          tag_ids: tagIDs,
          stash_ids: [
            ...(stashScene?.stash_ids ?? []),
            {
              endpoint,
              stash_id: scene.stash_id,
            },
          ],
        },
      },
    });

    setSaveState("");
    return sceneUpdateResult?.data?.sceneUpdate;
  };

  return handleSave;
}
