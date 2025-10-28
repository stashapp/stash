import React, { useState, useContext } from "react";
import * as GQL from "src/core/generated-graphql";

import StudioModal from "./StudioModal";
import PerformerModal from "../PerformerModal";
import { TaggerStateContext } from "../context";
import { useIntl } from "react-intl";
import { faTags } from "@fortawesome/free-solid-svg-icons";

type PerformerModalCallback = (toCreate?: GQL.PerformerCreateInput) => void;
type StudioModalCallback = (
  toCreate?: GQL.StudioCreateInput,
  parentInput?: GQL.StudioCreateInput
) => void;

export interface ISceneTaggerModalsContextState {
  createPerformerModal: (
    performer: GQL.ScrapedPerformerDataFragment,
    callback: PerformerModalCallback
  ) => void;
  createStudioModal: (
    studio: GQL.ScrapedSceneStudioDataFragment,
    callback: StudioModalCallback
  ) => void;
}

export const SceneTaggerModalsState =
  React.createContext<ISceneTaggerModalsContextState>({
    createPerformerModal: () => {},
    createStudioModal: () => {},
  });

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

  function handleStudioSave(
    toCreate: GQL.StudioCreateInput,
    parentInput?: GQL.StudioCreateInput
  ) {
    if (studioCallback) {
      studioCallback(toCreate, parentInput);
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

  const endpoint = currentSource?.sourceInput.stash_box_endpoint ?? undefined;

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
          icon={faTags}
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
          icon={faTags}
          header={intl.formatMessage(
            { id: "actions.create_entity" },
            { entityType: intl.formatMessage({ id: "studio" }) }
          )}
          endpoint={endpoint}
        />
      )}
      {children}
    </SceneTaggerModalsState.Provider>
  );
};
