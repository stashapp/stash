import React, { useMemo, useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
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
import { faChevronDown, faChevronUp } from "@fortawesome/free-solid-svg-icons";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";

interface IImageDetailProps {
  image: GQL.ImageDataFragment;
}

export const ImageDetailPanel: React.FC<IImageDetailProps> = (props) => {
  const intl = useIntl();

  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);

  const file = useMemo(
    () => (props.image.files.length > 0 ? props.image.files[0] : undefined),
    [props.image]
  );

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
  }

  function maybeRenderDetails() {
    if (!props.image.details) return;
    return (
      <div className="row details-description">
        <div className="col-12">
          <h5>
            <FormattedMessage id="details" />
            <Button
              className="minimal expand-collapse"
              onClick={() => setCollapsed(!collapsed)}
            >
              <Icon className="fa-fw" icon={getCollapseButtonIcon()} />
            </Button>
          </h5>
          <p className={`pre ${collapsed ? "collapsed" : ""}`}>
            {props.image.details}
          </p>
        </div>
      </div>
    );
  }

  function maybeRenderTags() {
    if (props.image.tags.length === 0) return;
    const tags = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="image" />
    ));
    return (
      <div className="row details-tags">
        <div className="col-12">
          <h5>
            <FormattedMessage id="tags" />
          </h5>
          {tags}
        </div>
      </div>
    );
  }

  function maybeRenderPerformers() {
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
      <div className="row details-performers">
        <div className="col-12">
          <h5>
            <FormattedMessage id="performers" />
          </h5>
          <div className="row image-performers">{cards}</div>
        </div>
      </div>
    );
  }

  function maybeRenderGalleries() {
    if (props.image.galleries.length === 0) return;
    const galleries = props.image.galleries.map((gallery) => (
      <GalleryLink key={gallery.id} gallery={gallery} />
    ));
    return (
      <div className="row details-galleries">
        <div className="col-12">
          <h5>
            <FormattedMessage
              id="countables.galleries"
              values={{ count: props.image.galleries.length }}
            />
          </h5>
          {galleries}
        </div>
      </div>
    );
  }

  return (
    <div className="col-xl-12 details-display">
      <div className="details-basic">
        <div className="row">
          <div className="image-header d-xl-none">
            <h3>
              <TruncatedText text={objectTitle(props.image)} />
            </h3>
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="studio" />
            </h5>
            {props.image.studio?.name ? (
              <h6>
                <Link to={`/studios/${props.image.studio.id}`}>
                  <TruncatedText text={props.image.studio.name} />
                </Link>
              </h6>
            ) : (
              <h6>&nbsp;</h6>
            )}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="date" />
            </h5>
            {props.image.date ? (
              <h6>
                <FormattedDate
                  value={props.image.date}
                  format="long"
                  timeZone="utc"
                />
              </h6>
            ) : (
              <h6>&nbsp;</h6>
            )}
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="scene_code" />
            </h5>
            {props.image.code ? <h6>{props.image.code}</h6> : <h6>&nbsp;</h6>}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="rating" />
            </h5>
            <RatingSystem value={props.image.rating100 ?? undefined} disabled />
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="photographer" />
            </h5>
            {props.image.photographer ? (
              <h6>{props.image.photographer}</h6>
            ) : (
              <h6>&nbsp;</h6>
            )}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="dimensions" />
            </h5>
            {file?.width && file?.height ? (
              <h6>
                {file.width}x{file.height}
              </h6>
            ) : (
              <h6>&nbsp</h6>
            )}
          </div>
        </div>
      </div>
      {maybeRenderGalleries()}
      {maybeRenderTags()}
      {maybeRenderDetails()}
      {maybeRenderPerformers()}
      <div className="row details-extra">
        <div className="col-12">
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.image.created_at)}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.image.updated_at)}
          </h6>
        </div>
      </div>
    </div>
  );
};
