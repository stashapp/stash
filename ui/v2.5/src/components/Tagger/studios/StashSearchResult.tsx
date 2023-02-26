import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { useUpdateStudio } from "../queries";
import StudioModal from "../scenes/StudioModal";
import { faTags } from "@fortawesome/free-solid-svg-icons";

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
  const [modalStudio, setModalStudio] = useState<
    GQL.ScrapedStudioDataFragment | undefined
  >();
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const updateStudio = useUpdateStudio();

  const handleSave = async (input: GQL.StudioCreateInput) => {
    setError({});
    setSaveState("Saving studio");
    setModalStudio(undefined);

    const updateData: GQL.StudioUpdateInput = {
      ...input,
      id: studio.id,
    };

    const res = await updateStudio(updateData);

    if (!res?.data?.studioUpdate)
      setError({
        message: `Failed to save studio "${studio.name}"`,
        details:
          res?.errors?.[0].message ===
          "UNIQUE constraint failed: studios.checksum"
            ? "Name already exists"
            : res?.errors?.[0].message,
      });
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
      <img src={(p.image ?? [])[0]} alt="" className="StudioTagger-thumb" />
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
