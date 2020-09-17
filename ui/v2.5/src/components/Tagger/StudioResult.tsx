import React, { useEffect, useState, Dispatch, SetStateAction } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import cx from "classnames";

import { SuccessIcon, Modal, StudioSelect } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { ValidTypes } from "src/components/Shared/Select";
import { IStashBoxStudio } from "./utils";

export interface IStudioOperation {
  create?: IStashBoxStudio;
  update?: GQL.StudioDataFragment | GQL.SlimStudioDataFragment;
  existing?: GQL.StudioDataFragment;
  skip?: boolean;
}

interface IStudioResultProps {
  studio: IStashBoxStudio | null;
  setStudio: Dispatch<SetStateAction<IStudioOperation | undefined>>;
}

const StudioResult: React.FC<IStudioResultProps> = ({ studio, setStudio }) => {
  const [selectedStudio, setSelectedStudio] = useState<string | null>();
  const [modalVisible, showModal] = useState(false);
  const [selectedSource, setSelectedSource] = useState<
    "create" | "existing" | "skip" | undefined
  >();
  const {
    data: stashIDData,
    loading: loadingStashID,
  } = GQL.useFindStudiosQuery({
    variables: {
      studio_filter: {
        stash_id: studio?.id ?? "",
      },
    },
  });
  const { data: searchData, loading: loadingSearch } = GQL.useFindStudiosQuery({
    variables: {
      filter: {
        q: `"${studio?.name ?? ""}"`,
      },
    },
  });

  useEffect(() => {
    if (stashIDData?.findStudios.studios?.[0]) {
      setStudio({
        existing: stashIDData.findStudios.studios[0],
      });
    } else if (searchData?.findStudios.studios?.[0]) {
      const result = searchData.findStudios.studios[0];
      setSelectedSource("existing");
      setSelectedStudio(result.id);
      setStudio({
        update: result,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stashIDData, searchData]);

  const handleStudioSelect = (newStudio: ValidTypes[]) => {
    if (newStudio.length) {
      setSelectedSource("existing");
      setSelectedStudio(newStudio[0].id);
      setStudio({
        update: newStudio[0] as GQL.SlimStudioDataFragment,
      });
    } else {
      setSelectedSource(undefined);
      setSelectedStudio(null);
    }
  };

  const handleStudioCreate = () => {
    if (!studio) return;
    setSelectedSource("create");
    setStudio({
      create: studio,
    });
    showModal(false);
  };

  const handleStudioSkip = () => {
    setSelectedSource("skip");
    setStudio({
      skip: true,
    });
  };

  if (loadingSearch || loadingStashID) return <div>Loading studio</div>;

  if (stashIDData?.findStudios.studios.length) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          Studio:
          <b className="ml-2">{studio?.name}</b>
        </div>
        <span className="ml-auto">
          <SuccessIcon className="mr-2" />
          Matched:
        </span>
        <b className="col-3 text-right">
          {stashIDData.findStudios.studios[0].name}
        </b>
      </div>
    );
  }

  return (
    <div className="row no-gutters align-items-center mt-2">
      <Modal
        show={modalVisible}
        accept={{ text: "Save", onClick: handleStudioCreate }}
        cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      >
        <div className="row">
          <strong className="col-2">Name:</strong>
          <span className="col-10">{studio?.name}</span>
        </div>
        <div className="row">
          <strong className="col-2">URL:</strong>
          <span className="col-10">{studio?.url ?? ""}</span>
        </div>
        <div className="row">
          <strong className="col-2">Logo:</strong>
          <span className="col-10">
            <img src={studio?.image ?? ""} alt="" />
          </span>
        </div>
      </Modal>

      <div className="entity-name">
        Studio:
        <b className="ml-2">{studio?.name}</b>
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
          onClick={() => handleStudioSkip()}
        >
          Skip
        </Button>
        <StudioSelect
          ids={selectedStudio ? [selectedStudio] : []}
          onSelect={handleStudioSelect}
          className={cx("studio-select", {
            "studio-select-active": selectedSource === "existing",
          })}
          isClearable={false}
        />
      </ButtonGroup>
    </div>
  );
};

export default StudioResult;
