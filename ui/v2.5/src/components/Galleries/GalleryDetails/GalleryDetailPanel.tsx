import React from "react";
import { Link } from "react-router-dom";
import { FormattedDate } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TagLink, TruncatedText } from "src/components/Shared";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";

interface IGalleryDetailProps {
  gallery: Partial<GQL.GalleryDataFragment>;
}

export const GalleryDetailPanel: React.FC<IGalleryDetailProps> = (props) => {
  function renderDetails() {
    if (!props.gallery.details || props.gallery.details === "") return;
    return (
      <>
        <h6>Details</h6>
        <p className="pre">{props.gallery.details}</p>
      </>
    );
  }

  function renderTags() {
    if (!props.gallery.tags || props.gallery.tags.length === 0) return;
    const tags = props.gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} tagType="gallery" />
    ));
    return (
      <>
        <h6>Tags</h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (!props.gallery.performers || props.gallery.performers.length === 0)
      return;
    const performers = sortPerformers(props.gallery.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.gallery.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>Performers</h6>
        <div className="row justify-content-center gallery-performers">
          {cards}
        </div>
      </>
    );
  }

  // filename should use entire row if there is no studio
  const galleryDetailsWidth = props.gallery.studio ? "col-9" : "col-12";
  const title =
    props.gallery.title ?? TextUtils.fileNameFromPath(props.gallery.path ?? "");

  return (
    <>
      <div className="row">
        <div className={`${galleryDetailsWidth} col-xl-12 gallery-details`}>
          <h3 className="gallery-header d-xl-none">
            <TruncatedText text={title} />
          </h3>
          {props.gallery.date ? (
            <h5>
              <FormattedDate
                value={props.gallery.date}
                format="long"
                timeZone="utc"
              />
            </h5>
          ) : undefined}
        </div>
        {props.gallery.studio && (
          <div className="col-3 d-xl-none">
            <Link to={`/studios/${props.gallery.studio.id}`}>
              <img
                src={props.gallery.studio.image_path ?? ""}
                alt={`${props.gallery.studio.name} logo`}
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
