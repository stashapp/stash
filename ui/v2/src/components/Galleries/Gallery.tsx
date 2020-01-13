import {
  Spinner,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { IBaseProps } from "../../models";
import { GalleryViewer } from "./GalleryViewer";

interface IProps extends IBaseProps {}

export const Gallery: FunctionComponent<IProps> = (props: IProps) => {
  const [gallery, setGallery] = useState<Partial<GQL.GalleryDataFragment>>({});
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindGallery(props.match.params.id);

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findGallery || !!error) { return; }
    setGallery(data.findGallery);
  }, [data, loading, error]);

  if (!data || !data.findGallery || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
  if (!!error) { return <>{error.message}</>; }
  return (
    <div style={{width: "75vw", margin: "0 auto"}}>
      <GalleryViewer gallery={gallery as any} />
    </div>
  );
};
