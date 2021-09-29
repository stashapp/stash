import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import PerformerModal from "./PerformerModal";
import StudioModal from "./StudioModal";

type PerformerModalCallback = (toCreate?: GQL.PerformerCreateInput) => void;
type StudioModalCallback = (toCreate?: GQL.StudioCreateInput) => void;

export interface ISceneScraperDialogsContextState {
  createPerformerModal: (
    performer: GQL.ScrapedPerformerDataFragment,
    callback: (toCreate?: GQL.PerformerCreateInput) => void
  ) => void;
  createStudioModal: (
    studio: GQL.ScrapedSceneStudioDataFragment,
    callback: (toCreate?: GQL.StudioCreateInput) => void
  ) => void;
}

export const SceneScraperDialogsState = React.createContext<ISceneScraperDialogsContextState>(
  {
    createPerformerModal: () => {},
    createStudioModal: () => {},
  }
);

export const SceneScraperModals: React.FC = ({ children }) => {
  const [performerToCreate, setPerformerToCreate] = useState<
    GQL.ScrapedPerformerDataFragment | undefined
  >();
  const [performerCallback, setPerformerCallback] = useState<
    PerformerModalCallback | undefined
  >();

  const [studioToCreate, setStudioToCreate] = useState<
    GQL.ScrapedSceneStudioDataFragment | undefined
  >();
  const [studioCallback, setStudioCallback] = useState<
    StudioModalCallback | undefined
  >();

  function handlePerformerSave(toCreate: GQL.PerformerCreateInput) {
    if (performerCallback) {
      performerCallback(toCreate);
    }

    setPerformerToCreate(undefined);
    setPerformerCallback(undefined);
  }

  function handlePerformerCancel() {
    if (performerCallback) {
      performerCallback();
    }

    setPerformerToCreate(undefined);
    setPerformerCallback(undefined);
  }

  function createPerformerModal(
    performer: GQL.ScrapedPerformerDataFragment,
    callback: PerformerModalCallback
  ) {
    setPerformerToCreate(performer);
    // can't set the function directly - needs to be via a wrapping function
    setPerformerCallback(() => callback);
  }

  function handleStudioSave(toCreate: GQL.StudioCreateInput) {
    if (studioCallback) {
      studioCallback(toCreate);
    }

    setStudioToCreate(undefined);
    setStudioCallback(undefined);
  }

  function handleStudioCancel() {
    if (studioCallback) {
      studioCallback();
    }

    setStudioToCreate(undefined);
    setStudioCallback(undefined);
  }

  function createStudioModal(
    studio: GQL.ScrapedSceneStudioDataFragment,
    callback: StudioModalCallback
  ) {
    setStudioToCreate(studio);
    // can't set the function directly - needs to be via a wrapping function
    setStudioCallback(() => callback);
  }

  return (
    <SceneScraperDialogsState.Provider
      value={{ createPerformerModal, createStudioModal }}
    >
      {performerToCreate && (
        <PerformerModal
          closeModal={handlePerformerCancel}
          modalVisible
          performer={performerToCreate}
          handlePerformerCreate={handlePerformerSave}
          icon="tags"
          header="Create Performer"
        />
      )}
      {studioToCreate && (
        <StudioModal
          closeModal={handleStudioCancel}
          modalVisible
          studio={studioToCreate}
          handleStudioCreate={handleStudioSave}
          icon="tags"
          header="Create Studio"
        />
      )}
      {children}
    </SceneScraperDialogsState.Provider>
  );
};
