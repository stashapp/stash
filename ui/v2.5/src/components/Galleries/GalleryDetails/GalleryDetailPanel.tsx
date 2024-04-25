import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { PhotographerLink } from "src/components/Shared/Link";
import { faCaretDown, faCaretUp } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useContainerDimensions } from "src/components/Shared/GridCard/GridCard";
import { DetailItem } from "src/components/Shared/DetailItem";

interface IGalleryDetailProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryDetailPanel: React.FC<IGalleryDetailProps> = ({
  gallery,
}) => {
  const intl = useIntl();

  const [collapsedDetails, setCollapsedDetails] = useState<boolean>(true);
  const [collapsedPerformers, setCollapsedPerformers] = useState<boolean>(true);
  const [collapsedTags, setCollapsedTags] = useState<boolean>(true);

  const [detailsRef, { height: detailsHeight }] = useContainerDimensions();
  const [perfRef, { height: perfHeight }] = useContainerDimensions();
  const [tagRef, { height: tagHeight }] = useContainerDimensions();

  const details = useMemo(() => {
    const limit = 160;
    return gallery.details?.length ? (
      <>
        <div
          className={`details ${
            collapsedDetails && detailsHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
        >
          <p className="pre" ref={detailsRef}>
            {gallery.details}
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
    gallery.details,
    detailsRef,
    detailsHeight,
    setCollapsedDetails,
    collapsedDetails,
  ]);

  const tags = useMemo(() => {
    const limit = 160;
    if (gallery.tags.length === 0) return;
    const galleryTags = gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <div
          className={`gallery-tags ${
            collapsedTags && tagHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={tagRef}
        >
          {galleryTags}
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
  }, [gallery.tags, tagRef, tagHeight, setCollapsedTags, collapsedTags]);

  function maybeRenderShowMoreLess(
    height: number,
    limit: number,
    ref: React.MutableRefObject<HTMLDivElement | null>,
    setCollapsed: React.Dispatch<React.SetStateAction<boolean>>,
    collapsed: boolean
  ) {
    if (height < limit) {
      return;
    }
    return (
      <span
        className={`show-${collapsed ? "more" : "less"}`}
        onClick={() => {
          const container = ref.current;
          if (container == null) {
            return;
          }
          if (container.style.maxHeight) {
            container.style.maxHeight = "";
          } else {
            container.style.maxHeight = container.scrollHeight + "px";
          }
          setCollapsed(!collapsed);
        }}
      >
        {collapsed ? "Show more" : "Show less"}
        <Icon className="fa-solid" icon={collapsed ? faCaretDown : faCaretUp} />
      </span>
    );
  }

  const performers = useMemo(() => {
    const limit = 365;
    const sorted = sortPerformers(gallery.performers);
    const cards = sorted.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={gallery.date ?? undefined}
      />
    ));

    return (
      <>
        <div
          className={`row justify-content-center gallery-performers ${
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
    gallery.performers,
    gallery.date,
    perfRef,
    perfHeight,
    collapsedPerformers,
    setCollapsedPerformers,
  ]);

  // filename should use entire row if there is no studio
  const galleryDetailsWidth = gallery.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${galleryDetailsWidth} col-12 gallery-details`}>
          <div className="detail-group">
            <DetailItem id="scene_code" value={gallery.code} fullWidth />
            <DetailItem
              id="director"
              value={
                gallery.photographer ? (
                  <PhotographerLink
                    photographer={gallery.photographer}
                    linkType="gallery"
                  />
                ) : undefined
              }
              fullWidth
            />
            <DetailItem id="details" value={details} />
            <DetailItem
              id="tags"
              heading={<FormattedMessage id="tags" />}
              value={gallery.tags.length ? tags : undefined}
            />
            <DetailItem
              id="performers"
              heading={<FormattedMessage id="performers" />}
              value={gallery.performers.length ? performers : undefined}
            />
            <DetailItem
              id="created_at"
              value={TextUtils.formatDateTime(intl, gallery.created_at)}
              fullWidth
            />
            <DetailItem
              id="updated_at"
              value={TextUtils.formatDateTime(intl, gallery.updated_at)}
              fullWidth
            />
          </div>
        </div>
      </div>
    </>
  );
};
