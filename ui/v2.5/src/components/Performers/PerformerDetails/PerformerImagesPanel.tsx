import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ImageList } from "src/components/Images/ImageList";
import { performerFilterHook } from "src/core/performers";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerImagesPanel: React.FC<IPerformerDetailsProps> = ({
  performer,
}) => {
  return <ImageList filterHook={performerFilterHook(performer)} />;
};
