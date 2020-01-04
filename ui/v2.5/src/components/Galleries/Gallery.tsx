import React, { useEffect, useState } from "react";
import { Spinner } from 'react-bootstrap';
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { IBaseProps } from "../../models";
import { GalleryViewer } from "./GalleryViewer";

interface IProps extends IBaseProps {}

export const Gallery: React.FC<IProps> = (props: IProps) => {
  const [gallery, setGallery] = useState<Partial<GQL.GalleryDataFragment>>({});
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindGallery(props.match.params.id);

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findGallery || !!error) { return; }
    setGallery(data.findGallery);
  }, [data]);

  if (!data || !data.findGallery || isLoading) { return <Spinner animation="border" variant="light" />; }
  if (!!error) { return <>{error.message}</>; }
  return (
    <div style={{width: "75vw", margin: "0 auto"}}>
      <GalleryViewer gallery={gallery as any} />
    </div>
  );
};
