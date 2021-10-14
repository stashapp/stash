import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { IStashBoxPerformer } from "../utils";
import { useUpdatePerformer } from "../queries";
import PerformerModal from "../PerformerModal";

interface IStashSearchResultProps {
  performer: GQL.SlimPerformerDataFragment;
  stashboxPerformers: IStashBoxPerformer[];
  endpoint: string;
  onPerformerTagged: (
    performer: Pick<GQL.SlimPerformerDataFragment, "id"> &
      Partial<Omit<GQL.SlimPerformerDataFragment, "id">>
  ) => void;
  excludedPerformerFields: string[];
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  performer,
  stashboxPerformers,
  onPerformerTagged,
  excludedPerformerFields,
  endpoint,
}) => {
  const [modalPerformer, setModalPerformer] = useState<
    IStashBoxPerformer | undefined
  >();
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const updatePerformer = useUpdatePerformer();

  const handleSave = async (input: GQL.PerformerCreateInput) => {
    const performerData = modalPerformer;
    if (performerData?.id) {
      setError({});
      setSaveState("Saving performer");
      setModalPerformer(undefined);

      const updateData: GQL.PerformerUpdateInput = {
        id: performerData.id,
        ...input,
      };

      const res = await updatePerformer(updateData);

      if (!res?.data?.performerUpdate)
        setError({
          message: `Failed to save performer "${performer.name}"`,
          details:
            res?.errors?.[0].message ===
            "UNIQUE constraint failed: performers.checksum"
              ? "Name already exists"
              : res?.errors?.[0].message,
        });
      else onPerformerTagged(performer);
      setSaveState("");
    }
  };

  const performers = stashboxPerformers.map((p) => (
    <Button
      className="PerformerTagger-performer-search-item minimal col-6"
      variant="link"
      key={p.stash_id}
      onClick={() => setModalPerformer(p)}
    >
      <img src={p.images[0]} alt="" className="PerformerTagger-thumb" />
      <span>{p.name}</span>
    </Button>
  ));

  return (
    <>
      {modalPerformer && (
        <PerformerModal
          closeModal={() => setModalPerformer(undefined)}
          modalVisible={modalPerformer !== undefined}
          performer={modalPerformer}
          onSave={handleSave}
          icon="tags"
          header="Update Performer"
          excludedPerformerFields={excludedPerformerFields}
          endpoint={endpoint}
        />
      )}
      <div className="PerformerTagger-performer-search">{performers}</div>
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
