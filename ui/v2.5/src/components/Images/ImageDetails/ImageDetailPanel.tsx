import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { GalleryLink, TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { sortPerformers } from "src/core/performers";
import { FormattedDate, FormattedMessage, useIntl } from "react-intl";
import { objectTitle } from "src/core/files";
import { PhotographerLink } from "src/components/Shared/Link";
interface IImageDetailProps {
  image: GQL.ImageDataFragment;
}

export const ImageDetailPanel: React.FC<IImageDetailProps> = (props) => {
  const intl = useIntl();

  const file = useMemo(
    () => (props.image.files.length > 0 ? props.image.files[0] : undefined),
    [props.image]
  );

  function renderDetails() {
    if (!props.image.details) return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{props.image.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.image.tags.length === 0) return;
    const tags = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="image" />
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
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.image.date ?? undefined}
      />
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
    const galleries = props.image.galleries.map((gallery) => (
      <GalleryLink key={gallery.id} gallery={gallery} />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.galleries"
            values={{ count: props.image.galleries.length }}
          />
        </h6>
        {galleries}
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
              <TruncatedText text={objectTitle(props.image)} />
            </h3>
          </div>
          {props.image.date ? (
            <h5>
              <FormattedDate
                value={props.image.date}
                format="long"
                timeZone="utc"
              />
            </h5>
          ) : undefined}
          {props.image.rating100 ? (
            <h6>
              <FormattedMessage id="rating" />:{" "}
              <RatingSystem value={props.image.rating100} disabled />
            </h6>
          ) : (
            ""
          )}

          {renderGalleries()}
          {file?.width && file?.height ? (
            <h6>
              <FormattedMessage id="resolution" />:{" "}
              {TextUtils.resolution(file.width, file.height)}
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
          {props.image.code && (
            <h6>
              <FormattedMessage id="scene_code" />: {props.image.code}{" "}
            </h6>
          )}
          {props.image.photographer && (
            <h6>
              <FormattedMessage id="photographer" />:{" "}
              <PhotographerLink
                photographer={props.image.photographer}
                linkType="image"
              />
            </h6>
          )}
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
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
        </div>
      </div>
    </>
  );
};
