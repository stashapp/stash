import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import cx from "classnames";

import { SuccessIcon, PerformerSelect } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { ValidTypes } from "src/components/Shared/Select";
import { sortImageURLs, IStashBoxPerformer } from "./utils";

import PerformerModal from "./PerformerModal";

export interface IPerformerOperation {
  create?: IStashBoxPerformer;
  update?: GQL.PerformerDataFragment | GQL.SlimPerformerDataFragment;
  existing?: GQL.PerformerDataFragment;
  skip?: boolean;
}

interface IPerformerResultProps {
  performer: IStashBoxPerformer;
  setPerformer: (data: IPerformerOperation) => void;
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
  const { data: stashData, loading: stashLoading } = GQL.useFindPerformersQuery(
    {
      variables: {
        performer_filter: {
          stash_id: performer.id,
        },
      },
    }
  );
  const { loading } = GQL.useFindPerformersQuery({
    variables: {
      filter: {
        q: `"${performer.name}"`,
      },
    },
    onCompleted: (data) => {
      const performerResult = data.findPerformers?.performers?.[0];
      if (performerResult) {
        setSelectedPerformer(performerResult.id);
        setSelectedSource("existing");
        setPerformer({
          update: performerResult,
        });
      }
    },
  });

  useEffect(() => {
    if (!stashData?.findPerformers.performers.length) return;

    setPerformer({
      existing: stashData.findPerformers.performers[0],
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stashData]);

  const handlePerformerSelect = (performers: ValidTypes[]) => {
    if (performers.length) {
      setSelectedSource("existing");
      setSelectedPerformer(performers[0].id);
      setPerformer({
        update: performers[0] as GQL.SlimPerformerDataFragment,
      });
    } else {
      setSelectedSource(undefined);
      setSelectedPerformer(null);
    }
  };

  const handlePerformerCreate = (imageIndex: number) => {
    // TODO
    /*
    const images = sortImageURLs(performer.images, "portrait");
    const imageURLs = images.length
      ? [
          {
            url: images[imageIndex].url,
            id: images[imageIndex].id,
            width: null,
            height: null,
          },
        ]
      : [];
     */
    setSelectedSource("create");
    setPerformer({
      create: {
        ...performer,
        images: [],
      },
    });
    showModal(false);
  };

  const handlePerformerSkip = () => {
    setSelectedSource("skip");
    setPerformer({
      skip: true,
    });
  };

  if (stashLoading || loading) return <div>Loading performer</div>;

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
