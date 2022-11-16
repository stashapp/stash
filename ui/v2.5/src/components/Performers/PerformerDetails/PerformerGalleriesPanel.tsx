import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { usePerformerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: GQL.PerformerDataFragment;
}

export const PerformerGalleriesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);
  return <GalleryList filterHook={filterHook} />;
};
