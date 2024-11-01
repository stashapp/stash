import React from "react";
import { Redirect, Route, RouteComponentProps, Switch } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import Gallery from "./GalleryDetails/Gallery";
import GalleryCreate from "./GalleryDetails/GalleryCreate";
import { GalleryList } from "./GalleryList";
import { View } from "../List/views";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { useFindGalleryImageID } from "src/core/StashService";

interface IGalleryImageParams {
  id: string;
  index: string;
}

const GalleryImage: React.FC<RouteComponentProps<IGalleryImageParams>> = ({
  match,
}) => {
  const { id, index: indexStr } = match.params;

  let index = parseInt(indexStr);
  if (isNaN(index)) {
    index = 0;
  }

  const { data, loading, error } = useFindGalleryImageID(id, index);

  if (isNaN(index)) {
    return <Redirect to={`/galleries/${id}`} />;
  }

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findGallery)
    return <ErrorMessage error={`No gallery found with id ${id}.`} />;

  return <Redirect to={`/images/${data.findGallery.image.id}`} />;
};

const Galleries: React.FC = () => {
  return <GalleryList view={View.Galleries} />;
};

const GalleryRoutes: React.FC = () => {
  const titleProps = useTitleProps({ id: "galleries" });
  return (
    <>
      <Helmet {...titleProps} />
      <Switch>
        <Route exact path="/galleries" component={Galleries} />
        <Route exact path="/galleries/new" component={GalleryCreate} />
        <Route
          exact
          path="/galleries/:id/images/:index"
          component={GalleryImage}
        />
        <Route path="/galleries/:id/:tab?" component={Gallery} />
      </Switch>
    </>
  );
};

export default GalleryRoutes;
