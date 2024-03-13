import React from "react";
import { Link } from "react-router-dom";
import { FormattedDate, FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { sortPerformers } from "src/core/performers";
import { galleryTitle } from "src/core/galleries";
import { PhotographerLink } from "src/components/Shared/Link";

interface IGalleryDetailProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryDetailPanel: React.FC<IGalleryDetailProps> = ({
  gallery,
}) => {
  const intl = useIntl();

  function renderDetails() {
    if (!gallery.details) return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{gallery.details}</p>
      </>
    );
  }

  function renderTags() {
    if (gallery.tags.length === 0) return;
    const tags = gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="gallery" />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: gallery.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (gallery.performers.length === 0) return;
    const performers = sortPerformers(gallery.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={gallery.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: gallery.performers.length }}
          />
        </h6>
        <div className="row justify-content-center gallery-performers">
          {cards}
        </div>
      </>
    );
  }

  // filename should use entire row if there is no studio
  const galleryDetailsWidth = gallery.studio ? "col-9" : "col-12";
  const title = galleryTitle(gallery);

  return (
    <>
      <div className="row">
        <div className={`${galleryDetailsWidth} col-xl-12 gallery-details`}>
          <h3 className="gallery-header d-xl-none">
            <TruncatedText text={title} />
          </h3>
          {gallery.date ? (
            <h5>
              <FormattedDate
                value={gallery.date}
                format="long"
                timeZone="utc"
              />
            </h5>
          ) : undefined}
          {gallery.rating100 ? (
            <h6>
              <FormattedMessage id="rating" />:{" "}
              <RatingSystem value={gallery.rating100} disabled />
            </h6>
          ) : (
            ""
          )}
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, gallery.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, gallery.updated_at)}{" "}
          </h6>
          {gallery.code && (
            <h6>
              <FormattedMessage id="scene_code" />: {gallery.code}{" "}
            </h6>
          )}
          {gallery.photographer && (
            <h6>
              <FormattedMessage id="photographer" />:{" "}
              <PhotographerLink
                photographer={gallery.photographer}
                linkType="gallery"
              />
            </h6>
          )}
        </div>
        {gallery.studio && (
          <div className="col-3 d-xl-none">
            <Link to={`/studios/${gallery.studio.id}`}>
              <img
                src={gallery.studio.image_path ?? ""}
                alt={`${gallery.studio.name} logo`}
                className="studio-logo float-right"
              />
            </Link>
          </div>
        )}
      </div>
      <div className="row">
        <div className="col-12">
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
        </div>
      </div>
    </>
  );
};
