import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { sceneAgeFromDate } from "src/utils/scene";
import { TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { DirectorLink } from "src/components/Shared/Link";
import { CustomFields } from "src/components/Shared/CustomFields";
import { Icon } from "src/components/Shared/Icon";
import { faLock } from "@fortawesome/free-solid-svg-icons";

const lockedFieldMessageIDs: Record<GQL.SceneLockableField, string> = {
  [GQL.SceneLockableField.Title]: "title",
  [GQL.SceneLockableField.Code]: "scene_code",
  [GQL.SceneLockableField.Details]: "details",
  [GQL.SceneLockableField.Director]: "director",
  [GQL.SceneLockableField.Date]: "date",
  [GQL.SceneLockableField.ProductionDate]: "production_date",
  [GQL.SceneLockableField.Studio]: "studio",
  [GQL.SceneLockableField.Urls]: "urls",
};

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();

  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{props.scene.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.scene.tags.length === 0) return;
    const tags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: props.scene.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.scene.performers.length === 0) return;
    const performers = sortPerformers(props.scene.performers);

    const ageFromDate = sceneAgeFromDate(
      props.scene.production_date,
      props.scene.date
    );

    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={ageFromDate}
      />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: props.scene.performers.length }}
          />
        </h6>
        <div className="row justify-content-center scene-performers">
          {cards}
        </div>
      </>
    );
  }

  function renderLockedFields() {
    if (props.scene.locked_fields.length === 0) return;
    const names = props.scene.locked_fields.map((f) =>
      intl.formatMessage({ id: lockedFieldMessageIDs[f] })
    );
    return (
      <h6 className="scene-locked-fields">
        <Icon icon={faLock} /> <FormattedMessage id="locked" />:{" "}
        {intl.formatList(names)}
      </h6>
    );
  }

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-12 scene-details`}>
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>
          {props.scene.production_date && (
            <h6>
              <FormattedMessage id="production_date" />:{" "}
              {/* fuzzy, so that a year- or month-only production date isn't
                  rendered as the first of the month */}
              {TextUtils.formatFuzzyDate(intl, props.scene.production_date)}
            </h6>
          )}
          {props.scene.code && (
            <h6>
              <FormattedMessage id="scene_code" />: {props.scene.code}{" "}
            </h6>
          )}
          {props.scene.director && (
            <h6>
              <FormattedMessage id="director" />:{" "}
              <DirectorLink director={props.scene.director} linkType="scene" />
            </h6>
          )}
          {renderLockedFields()}
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
          <CustomFields values={props.scene.custom_fields} fullWidth />
        </div>
      </div>
    </>
  );
};

export default SceneDetailPanel;
