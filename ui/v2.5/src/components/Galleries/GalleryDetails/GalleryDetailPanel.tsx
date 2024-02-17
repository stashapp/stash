import React, { useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
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
import {
  faChevronDown,
  faChevronUp
} from "@fortawesome/free-solid-svg-icons";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";

interface IGalleryDetailProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryDetailPanel: React.FC<IGalleryDetailProps> = ({
  gallery,
}) => {
  const intl = useIntl();
  
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
  }

  function maybeRenderDetails() {
    if (!gallery.details) return;
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
          <p className={`pre ${collapsed ? 'collapsed' : ''}`}>{gallery.details}</p>
        </div>
      </div>
    );
  }

  function maybeRenderTags() {
    if (gallery.tags.length === 0) return;
    const tags = gallery.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="gallery" />
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
      <div className="row details-performers">
        <div className="col-12">
          <h5>
            <FormattedMessage id="performers" />
          </h5>
          <div className="row image-performers">
            {cards}
          </div>
        </div>
      </div>
    );
  }

  const title = galleryTitle(gallery);

  return (
    <div className="col-xl-12 details-display">
      <div className="details-basic">
        <div className="row">
          <div className="gallery-header d-xl-none">
            <h3>
              <TruncatedText text={title} />
            </h3>
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="studio" />
            </h5>
            {gallery.studio?.name ? (
              <h6>
                <Link to={`/studios/${gallery.studio.id}`}>
                  <TruncatedText text={gallery.studio.name} />
                </Link>
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="date" />
            </h5>
            {gallery.date ? (
              <h6>
                <FormattedDate
                  value={gallery.date}
                  format="long"
                  timeZone="utc"
                />
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="scene_code" />
            </h5>
            {gallery.code ? (
              <h6>
                {gallery.code}
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="rating" />
            </h5>
            <RatingSystem value={gallery.rating100 ?? undefined} disabled />
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="photographer" />
            </h5>
            {gallery.photographer ? (
              <h6>
                {gallery.photographer}
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
        </div>
        {maybeRenderTags()}
        {maybeRenderDetails()}
        {maybeRenderPerformers()}
        <div className="row details-extra">
          <div className="col-12">
            <h6>
              <FormattedMessage id="created_at" />:{" "}
              {TextUtils.formatDateTime(intl, gallery.created_at)}
            </h6>
            <h6>
              <FormattedMessage id="updated_at" />:{" "}
              {TextUtils.formatDateTime(intl, gallery.updated_at)}
            </h6>
          </div>
        </div>
      </div>
    </div>
  );
};
