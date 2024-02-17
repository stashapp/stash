import React, { useMemo, useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { Link } from "react-router-dom";
import { FormattedDate, FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { objectTitle } from "src/core/files";
import {
  faChevronDown,
  faChevronUp
} from "@fortawesome/free-solid-svg-icons";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();
  const [collapsed, setCollapsed] = useState<boolean>(true);

  const file = useMemo(
    () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
    [props.scene]
  );

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
  }

  function maybeRenderDetails() {
    if (!props.scene.details || props.scene.details === "") return;
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
          <p className={`pre ${collapsed ? 'collapsed' : ''}`}>{props.scene.details}</p>
        </div>
      </div>
    );
  }

  function maybeRenderTags() {
    if (props.scene.tags.length === 0) return;
    const tags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
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
    if (props.scene.performers.length === 0) return;
    const performers = sortPerformers(props.scene.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));

    return (
      <div className="row details-performers">
        <div className="col-12">
          <h5>
            <FormattedMessage id="performers" />
          </h5>
          <div className="row scene-performers">
            {cards}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="col-xl-12 details-display">
      <div className="details-basic">
        <div className="row">
          <div className="scene-header d-xl-none">
            <h3>
              <TruncatedText text={objectTitle(props.scene)} />
            </h3>
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="studio" />
            </h5>
            {props.scene.studio?.name ? (
              <h6>
                <Link to={`/studios/${props.scene.studio.id}`}>
                  <TruncatedText text={props.scene.studio?.name} />
                </Link>
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="date" />
            </h5>
            {props.scene.date ? (
              <h6>
                <FormattedDate
                  value={props.scene.date}
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
            {props.scene.code ? (
              <h6>
                {props.scene.code}{" "}
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="rating" />
            </h5>
            <RatingSystem value={props.scene.rating100 ?? undefined} disabled />
          </div>
        </div>
        <div className="row">
          <div className="col-6">
            <h5>
              <FormattedMessage id="director" />
            </h5>
            {props.scene.director ? (
              <h6>
                <TruncatedText text={props.scene.director} />
              </h6>
            ) : (<h6>&nbsp;</h6>)}
          </div>
          <div className="col-6">
            <h5>
              <FormattedMessage id="resolution" />
            </h5>
            {file?.width && file?.height && (
              <h6>
                {TextUtils.resolution(file.width, file.height)}
              </h6>
            ) || (<h6>&nbsp;</h6>)}
          </div>
        </div>
      </div>
      {maybeRenderTags()}
      {maybeRenderDetails()}
      {maybeRenderPerformers()}
      <div className="row details-extra">
        <div className="col-12">
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>
        </div>
      </div>
    </div>
  );
};

export default SceneDetailPanel;
