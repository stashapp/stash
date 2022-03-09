import React, { useState, useContext } from "react";
import * as GQL from "src/core/generated-graphql";

import StudioModal from "./StudioModal";
import PerformerModal from "../PerformerModal";
import { TaggerStateContext } from "../context";
import { useIntl } from "react-intl";

type PerformerModalCallback = (toCreate?: GQL.PerformerCreateInput) => void;
type StudioModalCallback = (toCreate?: GQL.StudioCreateInput) => void;

export interface ISceneTaggerModalsContextState {
  createPerformerModal: (
    performer: GQL.ScrapedPerformerDataFragment,
    callback: (toCreate?: GQL.PerformerCreateInput) => void
  ) => void;
  createStudioModal: (
    studio: GQL.ScrapedSceneStudioDataFragment,
    callback: (toCreate?: GQL.StudioCreateInput) => void
  ) => void;
}

export const SceneTaggerModalsState = React.createContext<ISceneTaggerModalsContextState>(
  {
    createPerformerModal: () => {},
    createStudioModal: () => {},
  }
);

export const SceneTaggerModals: React.FC = ({ children }) => {
  const { currentSource } = useContext(TaggerStateContext);

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

  const intl = useIntl();

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

  const endpoint = currentSource?.stashboxEndpoint;

  return (
    <SceneTaggerModalsState.Provider
      value={{ createPerformerModal, createStudioModal }}
    >
      {performerToCreate && (
        <PerformerModal
          closeModal={handlePerformerCancel}
          modalVisible
          performer={performerToCreate}
          onSave={handlePerformerSave}
          icon="tags"
          header={intl.formatMessage(
            { id: "actions.create_entity" },
            { entityType: intl.formatMessage({ id: "performer" }) }
          )}
          endpoint={endpoint}
          create
        />
      )}
      {studioToCreate && (
        <StudioModal
          closeModal={handleStudioCancel}
          modalVisible
          studio={studioToCreate}
          handleStudioCreate={handleStudioSave}
          icon="tags"
          header={intl.formatMessage(
            { id: "actions.create_entity" },
            { entityType: intl.formatMessage({ id: "studio" }) }
          )}
        />
      )}
      {children}
    </SceneTaggerModalsState.Provider>
  );
};
