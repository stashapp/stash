import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { useUpdateStudio } from "../queries";
import StudioModal from "../scenes/StudioModal";
import { faTags } from "@fortawesome/free-solid-svg-icons";
import { useStudioCreate } from "src/core/StashService";
import { useIntl } from "react-intl";
import { apolloError } from "src/utils";
import { mergeStudioStashIDs } from "../utils";

interface IStashSearchResultProps {
  studio: GQL.SlimStudioDataFragment;
  stashboxStudios: GQL.ScrapedStudioDataFragment[];
  endpoint: string;
  onStudioTagged: (
    studio: Pick<GQL.SlimStudioDataFragment, "id"> &
      Partial<Omit<GQL.SlimStudioDataFragment, "id">>
  ) => void;
  excludedStudioFields: string[];
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  studio,
  stashboxStudios,
  onStudioTagged,
  excludedStudioFields,
  endpoint,
}) => {
  const intl = useIntl();

  const [modalStudio, setModalStudio] =
    useState<GQL.ScrapedStudioDataFragment>();
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const [createStudio] = useStudioCreate();
  const updateStudio = useUpdateStudio();

  function handleSaveError(name: string, message: string) {
    setError({
      message: intl.formatMessage(
        { id: "studio_tagger.failed_to_save_studio" },
        { studio: name }
      ),
      details:
        message === "UNIQUE constraint failed: studios.name"
          ? "Name already exists"
          : message,
    });
  }

  const handleSave = async (
    input: GQL.StudioCreateInput,
    parentInput?: GQL.StudioCreateInput
  ) => {
    setError({});
    setModalStudio(undefined);

    if (parentInput) {
      setSaveState("Saving parent studio");

      try {
        // if parent id is set, then update the existing studio
        if (input.parent_id) {
          const parentUpdateData: GQL.StudioUpdateInput = {
            ...parentInput,
            id: input.parent_id,
          };

          parentUpdateData.stash_ids = await mergeStudioStashIDs(
            input.parent_id,
            parentInput.stash_ids ?? []
          );

          await updateStudio(parentUpdateData);
        } else {
          const parentRes = await createStudio({
            variables: { input: parentInput },
          });
          input.parent_id = parentRes.data?.studioCreate?.id;
        }
      } catch (e) {
        handleSaveError(parentInput.name, apolloError(e));
      }
    }

    setSaveState("Saving studio");
    const updateData: GQL.StudioUpdateInput = {
      ...input,
      id: studio.id,
    };

    updateData.stash_ids = await mergeStudioStashIDs(
      studio.id,
      input.stash_ids ?? []
    );

    const res = await updateStudio(updateData);

    if (!res?.data?.studioUpdate)
      handleSaveError(studio.name, res?.errors?.[0]?.message ?? "");
    else onStudioTagged(studio);
    setSaveState("");
  };

  const studios = stashboxStudios.map((p) => (
    <Button
      className="StudioTagger-studio-search-item minimal col-6"
      variant="link"
      key={p.remote_site_id}
      onClick={() => setModalStudio(p)}
    >
      <img
        loading="lazy"
        src={(p.image ?? [])[0]}
        alt=""
        className="StudioTagger-thumb"
      />
      <span>{p.name}</span>
    </Button>
  ));

  return (
    <>
      {modalStudio && (
        <StudioModal
          closeModal={() => setModalStudio(undefined)}
          modalVisible={modalStudio !== undefined}
          studio={modalStudio}
          handleStudioCreate={handleSave}
          icon={faTags}
          header="Update Studio"
          excludedStudioFields={excludedStudioFields}
          endpoint={endpoint}
        />
      )}
      <div className="StudioTagger-studio-search">{studios}</div>
      <div className="row no-gutters mt-2 align-items-center justify-content-end">
        {error.message && (
          <div className="text-right text-danger mt-1">
            <strong>
              <span className="mr-2">Error:</span>
              {error.message}
            </strong>
            <div>{error.details}</div>
          </div>
        )}
        {saveState && (
          <strong className="col-4 mt-1 mr-2 text-right">{saveState}</strong>
        )}
      </div>
    </>
  );
};

export default StashSearchResult;
