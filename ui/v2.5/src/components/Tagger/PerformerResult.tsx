import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import cx from "classnames";

import { SuccessIcon, PerformerSelect } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { ValidTypes } from "src/components/Shared/Select";
import { IStashBoxPerformer } from "./utils";

import PerformerModal from "./PerformerModal";

export type PerformerOperation =
  | { type: 'create', data: IStashBoxPerformer }
  | { type: 'update', data: GQL.SlimPerformerDataFragment }
  | { type: 'existing', data: GQL.PerformerDataFragment }
  | { type: 'skip' };

interface IPerformerResultProps {
  performer: IStashBoxPerformer;
  setPerformer: (data: PerformerOperation) => void;
}

const PerformerResult: React.FC<IPerformerResultProps> = ({
  performer,
  setPerformer,
}) => {
  const [selectedPerformer, setSelectedPerformer] = useState<string | null>();
  const [selectedSource, setSelectedSource] = useState<
    "create" | "existing" | "skip" | undefined
  >();
  const [modalVisible, showModal] = useState(false);
  const { data: performerData } = GQL.useFindPerformerQuery({
    variables: { id: performer.id ?? '' },
    skip: !performer.id
    });
  const { data: stashData, loading: stashLoading } = GQL.useFindPerformersQuery(
    {
      variables: {
        performer_filter: {
          stash_id: performer.stash_id,
        },
      },
    }
  );

  useEffect(() => {
    if (stashData?.findPerformers.performers.length)
      setPerformer({
        type: 'existing',
        data: stashData.findPerformers.performers[0],
      });
    else if (performerData?.findPerformer) {
      setSelectedPerformer(performerData.findPerformer.id);
      setSelectedSource("existing");
      setPerformer({
        type: 'update',
        data: performerData.findPerformer,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stashData, performerData]);

  const handlePerformerSelect = (performers: ValidTypes[]) => {
    if (performers.length) {
      setSelectedSource("existing");
      setSelectedPerformer(performers[0].id);
      setPerformer({
        type: 'update',
        data: performers[0] as GQL.SlimPerformerDataFragment,
      });
    } else {
      setSelectedSource(undefined);
      setSelectedPerformer(null);
    }
  };

  const handlePerformerCreate = (imageIndex: number) => {
    const selectedImage = performer.images[imageIndex];
    const images = selectedImage ? [selectedImage] : [];
    setSelectedSource("create");
    setPerformer({
      type: 'create',
      data: {
        ...performer,
        images
      },
    });
    showModal(false);
  };

  const handlePerformerSkip = () => {
    setSelectedSource("skip");
    setPerformer({
      type: 'skip',
    });
  };

  if (stashLoading) return <div>Loading performer</div>;

  if (stashData?.findPerformers.performers?.[0]?.id) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          Performer:
          <b className="ml-2">{performer.name}</b>
        </div>
        <span className="ml-auto">
          <SuccessIcon />
          Matched:
        </span>
        <b className="col-3 text-right">
          {stashData.findPerformers.performers[0].name}
        </b>
      </div>
    );
  }
  return (
    <div className="row no-gutters align-items-center mt-2">
      <PerformerModal
        showModal={showModal}
        modalVisible={modalVisible}
        performer={performer}
        handlePerformerCreate={handlePerformerCreate}
      />
      <div className="entity-name">
        Performer:
        <b className="ml-2">{performer.name}</b>
      </div>
      <ButtonGroup>
        <Button
          variant={selectedSource === "create" ? "primary" : "secondary"}
          onClick={() => showModal(true)}
        >
          Create
        </Button>
        <Button
          variant={selectedSource === "skip" ? "primary" : "secondary"}
          onClick={() => handlePerformerSkip()}
        >
          Skip
        </Button>
        <PerformerSelect
          ids={selectedPerformer ? [selectedPerformer] : []}
          onSelect={handlePerformerSelect}
          className={cx("performer-select", {
            "performer-select-active": selectedSource === "existing",
          })}
          isClearable={false}
        />
      </ButtonGroup>
    </div>
  );
};

export default PerformerResult;
