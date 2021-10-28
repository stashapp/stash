import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { performerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: GQL.PerformerDataFragment;
}

export const PerformerGalleriesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  return <GalleryList filterHook={performerFilterHook(performer)} />;
};
