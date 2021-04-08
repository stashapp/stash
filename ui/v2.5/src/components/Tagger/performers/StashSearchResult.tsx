import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { IStashBoxPerformer, filterPerformer } from "../utils";
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

  const handleSave = async (image: number, excludedFields: string[]) => {
    if (modalPerformer) {
      const performerData = filterPerformer(modalPerformer, excludedFields);
      setError({});
      setSaveState("Saving performer");
      setModalPerformer(undefined);
      console.log(excludedFields);

      const res = await updatePerformer({
        ...performerData,
        image: excludedFields.includes("image")
          ? undefined
          : modalPerformer.images[image],
        stash_ids: [{ stash_id: modalPerformer.stash_id, endpoint }],
        id: performer.id,
      });

      if (!res.data?.performerUpdate)
        setError({
          message: `Failed to save performer "${performer.name}"`,
          details: res?.errors?.[0].message,
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
          handlePerformerCreate={handleSave}
          icon="tags"
          header="Update Performer"
          excludedPerformerFields={excludedPerformerFields}
        />
      )}
      <div className="PerformerTagger-performer-search">{performers}</div>
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
          <strong className="col-4 mt-1 mr-2 text-right">{saveState}</strong>
        )}
      </div>
    </>
  );
};

export default StashSearchResult;
