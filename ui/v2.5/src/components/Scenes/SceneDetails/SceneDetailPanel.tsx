import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { GalleryDetailedLink, TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { DirectorLink } from "src/components/Shared/Link";
import {
  DetailItem,
  maybeRenderShowMoreLess,
} from "src/components/Shared/DetailItem";
import { useContainerDimensions } from "src/components/Shared/GridCard/GridCard";
import { MovieCard } from "src/components/Movies/MovieCard";
import { GalleryCard } from "src/components/Galleries/GalleryCard";
import { faCompress, faExpand } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { Button } from "react-bootstrap";
import GalleryViewer from "src/components/Galleries/GalleryViewer";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();

  const [collapsedDetails, setCollapsedDetails] = useState<boolean>(true);
  const [collapsedGalleries, setCollapsedGalleries] = useState<boolean>(true);
  const [collapsedPerformers, setCollapsedPerformers] = useState<boolean>(true);
  const [collapsedTags, setCollapsedTags] = useState<boolean>(true);

  const [viewingGallery, setViewingGallery] = useState<number>(-1);

  const [detailsRef, { height: detailsHeight }] = useContainerDimensions();
  const [galleriesRef, { height: galleriesHeight }] = useContainerDimensions();
  const [perfRef, { height: perfHeight }] = useContainerDimensions();
  const [tagRef, { height: tagHeight }] = useContainerDimensions();

  const details = useMemo(() => {
    const limit = 160;
    return props.scene.details?.length ? (
      <>
        <div
          className={`details ${
            collapsedDetails && detailsHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
        >
          <p className="pre" ref={detailsRef}>
            {props.scene.details}
          </p>
        </div>
        {maybeRenderShowMoreLess(
          detailsHeight,
          limit,
          detailsRef,
          setCollapsedDetails,
          collapsedDetails
        )}
      </>
    ) : undefined;
  }, [
    props.scene.details,
    detailsRef,
    detailsHeight,
    setCollapsedDetails,
    collapsedDetails,
  ]);

  const movies = useMemo(
    () =>
      props.scene.movies.map((sceneMovie) => (
        <MovieCard
          key={sceneMovie.movie.id}
          movie={sceneMovie.movie}
          sceneIndex={sceneMovie.scene_index ?? undefined}
          titleOnImage={true}
        />
      )),
    [props.scene.movies]
  );

  const tags = useMemo(() => {
    const limit = 160;
    if (props.scene.tags.length === 0) return;
    const sceneTags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <div
          className={`scene-tags ${
            collapsedTags && tagHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={tagRef}
        >
          {sceneTags}
        </div>
        {maybeRenderShowMoreLess(
          tagHeight,
          limit,
          tagRef,
          setCollapsedTags,
          collapsedTags
        )}
      </>
    );
  }, [props.scene.tags, tagRef, tagHeight, setCollapsedTags, collapsedTags]);

  const galleries = useMemo(() => {
    const limit = 210;
    const sceneGalleries = props.scene.galleries.map((gallery, i) => (
      <div key={i} className="gallery-card-container">
        <Button
          className="minimal viewer-button"
          onClick={() => setViewingGallery(i)}
        >
          <Icon icon={faExpand} />
        </Button>
        <GalleryCard key={gallery.id} gallery={gallery} titleOnImage={true} />
      </div>
    ));
    /* provides a slimmer options users can swap to via CSS to reduce tab height */
    const slimSceneGalleries = props.scene.galleries.map((gallery) => (
      <GalleryDetailedLink key={gallery.id} gallery={gallery} />
    ));
    return (
      <>
        <div
          className={`scene-galleries ${
            collapsedGalleries && galleriesHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={galleriesRef}
        >
          {sceneGalleries}
          {slimSceneGalleries}
          {}
        </div>
        {maybeRenderShowMoreLess(
          galleriesHeight,
          limit,
          galleriesRef,
          setCollapsedGalleries,
          collapsedGalleries
        )}
      </>
    );
  }, [
    props.scene.galleries,
    galleriesRef,
    galleriesHeight,
    setCollapsedGalleries,
    collapsedGalleries,
  ]);

  const performers = useMemo(() => {
    const limit = 365;
    const sorted = sortPerformers(props.scene.performers);
    const cards = sorted.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));

    return (
      <>
        <div
          className={`row justify-content-center scene-performers ${
            collapsedPerformers && perfHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={perfRef}
        >
          {cards}
        </div>
        {maybeRenderShowMoreLess(
          perfHeight,
          limit,
          perfRef,
          setCollapsedPerformers,
          collapsedPerformers
        )}
      </>
    );
  }, [
    props.scene.performers,
    props.scene.date,
    perfRef,
    perfHeight,
    collapsedPerformers,
    setCollapsedPerformers,
  ]);

  function maybeRenderGalleryViewer() {
    if (viewingGallery >= 0) {
      return (
        <div className="gallery-viewer-container">
          <Button
            className="minimal viewer-button"
            onClick={() => setViewingGallery(-1)}
          >
            <Icon icon={faCompress} />
          </Button>
          <GalleryViewer galleryId={props.scene.galleries[viewingGallery].id} />
        </div>
      );
    }
  }

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  return (
    <>
      {maybeRenderGalleryViewer()}
      <div
        id="scene-details-panel"
        className={`row ${viewingGallery >= 0 ? "d-none" : ""}`}
      >
        <div className={`${sceneDetailsWidth} col-12 scene-details`}>
          <div className="detail-group">
            <DetailItem
              id="studio"
              value={props.scene.studio?.name}
              fullWidth
            />
            <DetailItem id="scene_code" value={props.scene.code} fullWidth />
            <DetailItem
              id="director"
              value={
                props.scene.director ? (
                  <DirectorLink
                    director={props.scene.director}
                    linkType="scene"
                  />
                ) : undefined
              }
              fullWidth
            />
            <DetailItem id="details" value={details} />
            <DetailItem
              id="tags"
              heading={<FormattedMessage id="tags" />}
              value={props.scene.tags.length ? tags : undefined}
            />
            <DetailItem
              id="performers"
              heading={<FormattedMessage id="performers" />}
              value={props.scene.performers.length ? performers : undefined}
            />
            <DetailItem
              id="movies"
              value={movies.length ? movies : undefined}
            />
            <DetailItem
              id="galleries"
              heading={<FormattedMessage id="galleries" />}
              value={props.scene.galleries.length ? galleries : undefined}
            />
            <DetailItem
              id="created_at"
              value={TextUtils.formatDateTime(intl, props.scene.created_at)}
              fullWidth
            />
            <DetailItem
              id="updated_at"
              value={TextUtils.formatDateTime(intl, props.scene.updated_at)}
              fullWidth
            />
          </div>
        </div>
      </div>
    </>
  );
};

export default SceneDetailPanel;
