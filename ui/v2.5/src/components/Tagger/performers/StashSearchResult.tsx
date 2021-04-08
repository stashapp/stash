import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, SuccessIcon } from "src/components/Shared";
import { IStashBoxPerformer } from "../utils";
import { useUpdatePerformerStashID } from "../queries";
import PerformerModal from "../PerformerModal";

interface IStashSearchResultProps {
  performer: GQL.SlimPerformerDataFragment;
  stashboxPerformers: IStashBoxPerformer[];
  endpoint: string;
  onPerformerTagged: (
    performer: Pick<GQL.SlimPerformerDataFragment, "id"> &
      Partial<Omit<GQL.SlimPerformerDataFragment, "id">>
  ) => void;
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  performer,
  stashboxPerformers,
  endpoint,
  onPerformerTagged,
}) => {
  const [modalPerformer, setModalPerformer] = useState<
    IStashBoxPerformer | undefined
  >();
  const [selectedPerformer, setSelectedPerformer] = useState<
    IStashBoxPerformer | undefined
  >();
  const [performerImage, setPerformerImage] = useState<Record<string, number>>(
    {}
  );
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const updatePerformerStashID = useUpdatePerformerStashID();

  const handleSave = async () => {
    if (selectedPerformer) {
      setError({});
      setSaveState("Saving performer");

      const res = await updatePerformerStashID(performer.id, [
        ...performer.stash_ids,
        { stash_id: selectedPerformer.stash_id, endpoint },
      ]);

      if (!res.data?.performerUpdate)
        setError({
          message: `Failed to save stashID to performer "${performer.name}"`,
          details: res?.errors?.[0].message,
        });
      else onPerformerTagged(performer);
      setSaveState("");
    }
  };

  const performers = stashboxPerformers.map((p) => {
    const isActive = selectedPerformer?.stash_id === p.stash_id;

    return (
      <Button
        className="PerformerTagger-performer-search-item minimal col-6"
        variant="link"
        key={p.stash_id}
        onClick={() => setModalPerformer(p)}
      >
        <img
          src={p.images[performerImage[p.stash_id] ?? 0]}
          alt=""
          className="PerformerTagger-thumb"
        />
        <span className={isActive ? "font-weight-bold" : ""}>{p.name}</span>
        {isActive && <SuccessIcon />}
      </Button>
    );
  });

  const handlePerformerSelect = (image: number) => {
    setSelectedPerformer(modalPerformer);
    if (modalPerformer?.stash_id)
      setPerformerImage({
        ...performerImage,
        [modalPerformer.stash_id]: image,
      });
    setModalPerformer(undefined);
  };

  return (
    <>
      {modalPerformer && (
        <PerformerModal
          closeModal={() => setModalPerformer(undefined)}
          modalVisible={modalPerformer !== undefined}
          performer={modalPerformer}
          handlePerformerCreate={handlePerformerSelect}
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
        <Button onClick={handleSave} disabled={selectedPerformer === undefined}>
          {saveState ? <LoadingIndicator inline small message="" /> : "Save"}
        </Button>
      </div>
    </>
  );
};

export default StashSearchResult;
