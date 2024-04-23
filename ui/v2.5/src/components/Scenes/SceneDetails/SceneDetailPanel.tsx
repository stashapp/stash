import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { DirectorLink } from "src/components/Shared/Link";
import { DetailItem } from "src/components/Shared/DetailItem";
import { faCaretDown, faCaretUp } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useContainerDimensions } from "src/components/Shared/GridCard/GridCard";
import { MovieCard } from "src/components/Movies/MovieCard";
import { GalleryCard } from "src/components/Galleries/GalleryCard";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();

  const [collapsedDetails, setCollapsedDetails] = useState<boolean>(true);
  const [collapsedGalleries, setCollapsedGalleries] = useState<boolean>(true);
  const [collapsedPerformers, setCollapsedPerformers] = useState<boolean>(true);
  const [collapsedTags, setCollapsedTags] = useState<boolean>(true);

  const [detailsRef, { height: detailsHeight }] = useContainerDimensions();
  const [galleriesRef, { height: galleriesHeight }] = useContainerDimensions();
  const [perfRef, { height: perfHeight }] = useContainerDimensions();
  const [tagRef, { height: tagHeight }] = useContainerDimensions();

  const details = useMemo(() => {
    return props.scene.details?.length ? (
      <>
        <div
          className={`details ${
            collapsedDetails ? "collapsed-detail" : "expanded-detail"
          }`}
          ref={detailsRef}
        >
          <p className="pre">{props.scene.details}</p>
        </div>
        {maybeRenderShowMoreLess(
          detailsHeight,
          160,
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
        />
      )),
    [props.scene.movies]
  );

  const tags = useMemo(() => {
    if (props.scene.tags.length === 0) return;
    const sceneTags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <div
          className={`scene-tags ${
            collapsedTags ? "collapsed-detail" : "expanded-detail"
          }`}
          ref={tagRef}
        >
          {sceneTags}
        </div>
        {maybeRenderShowMoreLess(
          tagHeight,
          160,
          setCollapsedTags,
          collapsedTags
        )}
      </>
    );
  }, [props.scene.tags, tagRef, tagHeight, setCollapsedTags, collapsedTags]);

  const galleries_v2 = useMemo(() => {
    const sceneGalleries = props.scene.galleries.map((gallery) => (
      <GalleryCard key={gallery.id} gallery={gallery} />
    ));
    return (
      <>
        <div
          className={`scene-galleries ${
            collapsedGalleries ? "collapsed-detail" : "expanded-detail"
          }`}
          ref={galleriesRef}
        >
          {sceneGalleries}
        </div>
        {maybeRenderShowMoreLess(
          galleriesHeight,
          240,
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

  function maybeRenderShowMoreLess(
    height: number,
    limit: number,
    setCollapsed: React.Dispatch<React.SetStateAction<boolean>>,
    collapsed: boolean
  ) {
    if (height < limit) {
      return;
    }
    return (
      <span
        className={`show-${collapsed ? "more" : "less"}`}
        onClick={() => setCollapsed(!collapsed)}
      >
        {collapsed ? "Show more" : "Show less"}
        <Icon className="fa-solid" icon={collapsed ? faCaretDown : faCaretUp} />
      </span>
    );
  }

  const performers = useMemo(() => {
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
            collapsedPerformers ? "collapsed-detail" : "expanded-detail"
          }`}
          ref={perfRef}
        >
          {cards}
        </div>
        {maybeRenderShowMoreLess(
          perfHeight,
          165,
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

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-12 scene-details`}>
          <div className="detail-group">
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
              value={props.scene.galleries.length ? galleries_v2 : undefined}
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
