import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TagLink, TruncatedText } from "src/components/Shared";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { sortPerformers } from "src/core/performers";
import { FormattedMessage, useIntl } from "react-intl";

interface IImageDetailProps {
  image: GQL.ImageDataFragment;
}

export const ImageDetailPanel: React.FC<IImageDetailProps> = (props) => {
  const intl = useIntl();

  function renderTags() {
    if (props.image.tags.length === 0) return;
    const tags = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} tagType="image" />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: props.image.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.image.performers.length === 0) return;
    const performers = sortPerformers(props.image.performers);
    const cards = performers.map((performer) => (
      <PerformerCard key={performer.id} performer={performer} />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: props.image.performers.length }}
          />
        </h6>
        <div className="row justify-content-center image-performers">
          {cards}
        </div>
      </>
    );
  }

  function renderGalleries() {
    if (props.image.galleries.length === 0) return;
    const tags = props.image.galleries.map((gallery) => (
      <TagLink key={gallery.id} gallery={gallery} />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.galleries"
            values={{ count: props.image.galleries.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  // filename should use entire row if there is no studio
  const imageDetailsWidth = props.image.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${imageDetailsWidth} col-xl-12 image-details`}>
          <div className="image-header d-xl-none">
            <h3>
              <TruncatedText
                text={
                  props.image.title ??
                  TextUtils.fileNameFromPath(props.image.path)
                }
              />
            </h3>
          </div>
          {props.image.rating ? (
            <h6>
              <FormattedMessage id="rating" />:{" "}
              <RatingStars value={props.image.rating} />
            </h6>
          ) : (
            ""
          )}
          {renderGalleries()}
          {props.image.file.width && props.image.file.height ? (
            <h6>
              <FormattedMessage id="resolution" />:{" "}
              {TextUtils.resolution(
                props.image.file.width,
                props.image.file.height
              )}
            </h6>
          ) : (
            ""
          )}
          {
            <h6>
              {" "}
              <FormattedMessage id="created_at" />:{" "}
              {TextUtils.formatDateTime(intl, props.image.created_at)}{" "}
            </h6>
          }
          {
            <h6>
              <FormattedMessage id="updated_at" />:{" "}
              {TextUtils.formatDateTime(intl, props.image.updated_at)}{" "}
            </h6>
          }
        </div>
        {props.image.studio && (
          <div className="col-3 d-xl-none">
            <Link to={`/studios/${props.image.studio.id}`}>
              <img
                src={props.image.studio.image_path ?? ""}
                alt={`${props.image.studio.name} logo`}
                className="studio-logo float-right"
              />
            </Link>
          </div>
        )}
      </div>
      <div className="row">
        <div className="col-12">
          {renderTags()}
          {renderPerformers()}
        </div>
      </div>
    </>
  );
};
